package handler

import (
	"context"

	"github.com/google/uuid"
)

type GetHandler struct{}

func NewGetHandler() *GetHandler {
	return &GetHandler{}
}

func (h *GetHandler) Handle(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	return &GetResponse{}, nil
}

type GetRequest struct {
	ID uuid.UUID `path:"id"`
}

type GetResponse struct{}
