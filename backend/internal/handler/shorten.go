package handler

import (
	"context"
)

type ShortenHandler struct {
}

func NewShortenHandler() *ShortenHandler {
	return &ShortenHandler{}
}

func (s *ShortenHandler) Handle(ctx context.Context, i *ShortenRequest) (*ShortenResponse, error) {
	return &ShortenResponse{}, nil
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}
