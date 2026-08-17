package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"url-shortener/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("URL Shortener API", "1.0.0"))

	// Pass nil for DB client since we only need route metadata
	srv := handler.NewServer(nil, api)
	srv.Register()

	data, err := json.MarshalIndent(api.OpenAPI(), "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal openapi: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}
