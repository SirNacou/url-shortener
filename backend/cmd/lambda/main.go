package main

import (
	"log"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/handler"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /api/urls", handler.List())
	mux.Handle("GET /api/urls/{id}", handler.Get())

	log.Printf("Server listening on port %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		panic(err)
	}
}
