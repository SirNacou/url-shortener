package handler

import (
	"context"
)

type ListHandler struct{}

func NewListHandler() *ListHandler {
	return &ListHandler{}
}

func (h *ListHandler) Handle(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	return &ListResponse{}, nil
}

type ListRequest struct{}
type ListResponse struct{}
