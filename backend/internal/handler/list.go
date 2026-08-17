package handler

import (
	"context"
	"fmt"
	"time"
	"url-shortener/internal/handler/utils"
	"url-shortener/internal/repository"

	"github.com/danielgtaylor/huma/v2"
)

type ListHandler struct {
	repo *repository.URLRepository
}

func NewListHandler(repo *repository.URLRepository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	urls, cursor, err := h.repo.List(ctx, int32(req.Limit), req.Cursor)
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to fetch URLs: %v", err))
	}

	baseURL := req.ResolveBaseURL("")

	items := make([]URLItem, len(urls))
	for i, r := range urls {
		items[i] = URLItem{
			Code:          r.Code,
			ShortURL:      fmt.Sprintf("%s/s/%s", baseURL, r.Code),
			TargetURL:     r.TargetURL,
			Clicks:        r.Clicks,
			CreatedAt:     r.CreateAt,
			ExpiresAt:     r.ExpiresAt,
			LastClickedAt: r.LastClickedAt,
		}
	}

	return &ListResponse{
		Body: ResponseBody{
			Items:      items,
			NextCursor: cursor,
		},
	}, nil
}

type ListRequest struct {
	Limit  int    `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Number of items to return"`
	Cursor string `query:"cursor" doc:"Pagination token from previous response"`

	utils.URLResolver
}
type ListResponse struct {
	Body ResponseBody
}

type ResponseBody struct {
	Items      []URLItem `json:"items"`
	NextCursor string    `json:"next_cursor,omitempty" example:"eyJwayI6IkNPREUjeHl6MTIzIn0="`
}
type URLItem struct {
	Code          string     `json:"code" example:"xyz123"`
	ShortURL      string     `json:"short_url" example:"https://nacou.dev/xyz123"`
	TargetURL     string     `json:"target_url" example:"https://example.com/very/long/url"`
	Clicks        int64      `json:"clicks" example:"42"`
	CreatedAt     time.Time  `json:"created_at" example:"2022-01-01T00:00:00Z"`
	ExpiresAt     time.Time  `json:"expires_at" example:"2022-01-01T00:00:00Z"`
	LastClickedAt *time.Time `json:"last_clicked_at,omitempty" example:"2022-01-01T00:00:00Z"`
}
