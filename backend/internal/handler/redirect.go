package handler

import "context"

type RedirectHandler struct {
}

func NewRedirectHandler() *RedirectHandler {
	return &RedirectHandler{}
}

func (r *RedirectHandler) Handle(ctx context.Context, i *RedirectRequest) (*RedirectResponse, error) {
	return &RedirectResponse{}, nil
}

type RedirectRequest struct {
	Code string `path:"code"`
}

type RedirectResponse struct{}
