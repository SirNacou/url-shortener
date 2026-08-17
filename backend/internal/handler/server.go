package handler

import (
	"context"
	"net/http"
	"url-shortener/internal/handler/shorten"
	"url-shortener/internal/repository"

	"github.com/danielgtaylor/huma/v2"
)

type Server struct {
	repo *repository.URLRepository
	api  huma.API
}

func NewServer(repo *repository.URLRepository, api huma.API) *Server {
	return &Server{repo: repo, api: api}
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
		Method:        "GET",
		Path:          "/s/{code}",
		DefaultStatus: http.StatusTemporaryRedirect,
	}, NewRedirectHandler(s.repo).Handle)

	huma.Register(s.api, huma.Operation{
		Method: "POST",
		Path:   "/api/shorten",
	}, shorten.NewShortenHandler(s.repo).Handle)

	huma.Register(s.api, huma.Operation{
		Method: "GET",
		Path:   "/api/urls",
	}, NewListHandler().Handle)

	huma.Register(s.api, huma.Operation{
		Method: "GET",
		Path:   "/api/urls/{id}",
	}, NewGetHandler().Handle)
}
