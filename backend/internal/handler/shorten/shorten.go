package shorten

import (
	"context"
	"errors"
	"fmt"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/repository"

	"github.com/danielgtaylor/huma/v2"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

var DefaultBackoff = wait.Backoff{
	Duration: 100 * time.Millisecond,
	Factor:   1,
	Jitter:   0.1,
	Steps:    3,
	Cap:      0,
}

type ShortenHandler struct {
	repo              *repository.URLRepository
	configuredBaseURL string
}

func NewShortenHandler(repo *repository.URLRepository) *ShortenHandler {
	return &ShortenHandler{repo: repo}
}

func (h *ShortenHandler) Handle(ctx context.Context, req *ShortenRequest) (*ShortenResponse, error) {

	if req.Body.CustomCode != "" {
		item := domain.NewURL(req.Body.CustomCode, req.Body.URL, req.Body.ExpiresInDays)
		err := h.repo.Save(ctx, item)
		if err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				return nil, huma.Error409Conflict("custom code already in use")
			}
			return nil, huma.Error500InternalServerError("failed to save custom URL")
		}

		return buildResponse(item, req.resolveBaseURL("")), nil
	}

	res := new(ShortenResponse)
	err := retry.RetryOnConflict(DefaultBackoff, func() error {
		code, err := GenerateCode(6)
		if err != nil {
			return err
		}

		record := domain.NewURL(code, req.Body.URL, req.Body.ExpiresInDays)

		err = h.repo.Save(ctx, record)
		if err != nil {
			if !errors.Is(err, domain.ErrAlreadyExists) {
				return fmt.Errorf("database error")
			}
			return err
		}

		res = buildResponse(record, req.resolveBaseURL(""))
		return nil
	})

	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return res, nil
}

type ShortenRequest struct {
	XForwardedProto string `header:"x-forwarded-proto" doc:"Protocol from CloudFront/ALB (http or https)"`
	XForwardedHost  string `header:"x-forwarded-host" doc:"Host from CloudFront/ALB"`
	Host            string `header:"host" doc:"Fallback host header"`

	Body struct {
		URL           string `json:"url" format:"uri"`
		CustomCode    string `json:"custom_code,omitempty" minLength:"3" maxLength:"30" pattern:"^[a-zA-Z0-9_-]+$" doc:"Optional custom slug"`
		ExpiresInDays int    `json:"expires_in_days,omitempty" minimum:"1" maximum:"365" default:"30" doc:"Duration before expiration"`
	}
}

type ShortenResponse struct {
	Body ShortenResponseBody
}
type ShortenResponseBody struct {
	ShortURL  string    `json:"short_url" example:"https://s.apps.nacou.dev/s/xyz123"`
	Code      string    `json:"code" example:"xyz123"`
	ExpiresAt time.Time `json:"expires_at" example:"2022-01-01T00:00:00Z"`
}

func buildResponse(record *domain.URL, baseURL string) *ShortenResponse {
	return &ShortenResponse{
		Body: ShortenResponseBody{
			Code:      record.Code,
			ShortURL:  fmt.Sprintf("%s/s/%s", baseURL, record.Code),
			ExpiresAt: record.ExpiresAt,
		},
	}
}
