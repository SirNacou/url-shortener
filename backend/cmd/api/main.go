package main

import (
	"context"
	"log"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/handler"
	"url-shortener/internal/logger"
	"url-shortener/internal/middleware"
	"url-shortener/internal/repository"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func main() {
	ctx := context.Background()

	logger.Init()

	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(cfg.AWSConfig)
	repo := repository.NewURLRepository(dynamodbClient, cfg.TableName)

	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("URL Shortener API", "1.0.0"))
	api.UseMiddleware(middleware.Logger())

	srv := handler.NewServer(repo, api)
	srv.Register()

	log.Printf("Server listening on port %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		panic(err)
	}
}
