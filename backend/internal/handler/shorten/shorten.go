package shorten

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/handler/utils"
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
	// 1. Create a structured logger with common request context
	log := slog.With(
		"url", req.Body.URL,
		"custom_code", req.Body.CustomCode,
		"expires_in_days", req.Body.ExpiresInDays,
	)

	if req.Body.CustomCode != "" {
		log.Info("attempting to create custom short URL")

		item := domain.NewURL(req.Body.CustomCode, req.Body.URL, req.Body.ExpiresInDays)
		err := h.repo.Save(ctx, item)
		if err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				log.Warn("custom code conflict")
				return nil, huma.Error409Conflict("custom code already in use")
			}
			log.Error("failed to save custom URL", "error", err)
			return nil, huma.Error500InternalServerError("failed to save custom URL")
		}

		log.Info("custom URL created successfully", "code", item.Code)
		return buildResponse(item, req.ResolveBaseURL("")), nil
	}

	// 2. Logic for random code generation
	res := new(ShortenResponse)
	attempts := 0

	err := retry.RetryOnConflict(DefaultBackoff, func() error {
		attempts++
		code, err := GenerateCode(6)
		if err != nil {
			log.Error("failed to generate random code", "error", err)
			return err
		}

		record := domain.NewURL(code, req.Body.URL, req.Body.ExpiresInDays)

		err = h.repo.Save(ctx, record)
		if err != nil {
			if errors.Is(err, domain.ErrAlreadyExists) {
				log.Warn("collision detected, retrying", "attempt", attempts, "code", code)
				return err
			}
			log.Error("database save failed", "error", err)
			return fmt.Errorf("database error")
		}

		res = buildResponse(record, req.ResolveBaseURL(""))
		return nil
	})

	if err != nil {
		log.Error("failed to generate short URL after retries", "error", err, "total_attempts", attempts)
		return nil, huma.Error500InternalServerError(err.Error())
	}

	log.Info("random URL created successfully", "code", res.Body.Code, "attempts", attempts)
	return res, nil
}

type ShortenRequest struct {
	utils.URLResolver

	Body struct {
		URL           string `json:"url" format:"uri"`
		CustomCode    string `json:"custom_code,omitempty" minLength:"3" maxLength:"30" pattern:"^[a-zA-Z0-9_-]+$" doc:"Optional custom slug"`
		ExpiresInDays int    `json:"expires_in_days" required:"true" minimum:"1" maximum:"365" default:"30" doc:"Duration before expiration"`
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
