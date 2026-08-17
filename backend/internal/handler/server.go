package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/danielgtaylor/huma/v2"
)

type Server struct {
	db        *dynamodb.Client
	tableName string
	api       huma.API
}

func NewServer(db *dynamodb.Client, tableName string, api huma.API) *Server {
	return &Server{db: db, tableName: tableName, api: api}
}

func (s *Server) Register() {
	huma.Register(s.api, huma.Operation{
		Method:        "GET",
		Path:          "/healthz",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct{}) (*struct{}, error) {
		return nil, nil
	})

	huma.Register(s.api, huma.Operation{
		Method: "GET",
		Path:   "/s/{code}",
	}, NewRedirectHandler().Handle)

	huma.Register(s.api, huma.Operation{
		Method: "POST",
		Path:   "/api/shorten",
	}, NewShortenHandler().Handle)

	huma.Register(s.api, huma.Operation{
		Method: "GET",
		Path:   "/api/urls",
	}, NewListHandler().Handle)

	huma.Register(s.api, huma.Operation{
		Method: "GET",
		Path:   "/api/urls/{id}",
	}, func(ctx context.Context, i *struct{}) (*struct{}, error) {
		return nil, nil
	})
}
