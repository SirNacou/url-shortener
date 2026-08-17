package main

import (
	"context"
	"log"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/handler"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(cfg.AWSConfig)

	srv := handler.NewServer(dynamodbClient, cfg.TableName)

	mux := http.NewServeMux()
	mux.Handle("GET /api/urls", srv.List())
	mux.Handle("GET /api/urls/{id}", srv.Get())

	log.Printf("Server listening on port %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		panic(err)
	}
}
