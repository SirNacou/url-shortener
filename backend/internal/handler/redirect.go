package handler

import (
	"context"
	"errors"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/repository"

	"github.com/danielgtaylor/huma/v2"
)

type RedirectHandler struct {
	repo *repository.URLRepository
}

func NewRedirectHandler(repo *repository.URLRepository) *RedirectHandler {
	return &RedirectHandler{repo: repo}
}

func (r *RedirectHandler) Handle(ctx context.Context, req *RedirectRequest) (*RedirectResponse, error) {
	record, err := r.repo.Find(ctx, req.Code)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, huma.Error404NotFound("short URL not found")
		}
		return nil, huma.Error500InternalServerError("failed to fetch target URL")
	}
	now := time.Now()
	if record.ExpiresAt.Before(now) {
		return nil, huma.Error404NotFound("short URL has expired")
	}

	return &RedirectResponse{
		Location: record.TargetURL,
	}, nil
}

type RedirectRequest struct {
	Code string `path:"code" minLength:"1" maxLength:"30" pattern:"^[a-zA-Z0-9_-]+$" doc:"Short code slug"`
}

type RedirectResponse struct {
	Location string `header:"Location" doc:"Target destination URL"`
}
