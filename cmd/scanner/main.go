package main

import (
    "log"
    "net/http"
    "github.com/stevan1008/scanner-service-api-go/internal/implement"
	"github.com/stevan1008/scanner-service-api-go/internal/api"
)

func main() {
    zapAPIURL := "http://zap:8090"
    zapAPIKey := "fookey"
    zapScanner := implement.NewZAPAdapter(zapAPIURL, zapAPIKey)
    handler := api.NewHandler(zapScanner)
    mux := http.NewServeMux()
    handler.RegisterRoutes(mux)
    port := "8080"
    log.Printf("Starting scanner service on port %s\n", port)
    // Run the server
    if err := http.ListenAndServe(":"+port, mux); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}