package main

import (
	"context"
	"log/slog"

	"github.com/ankittk/bookService/pkg/logger"
)

var log *logger.Logger

func init() {
	// Initialize the logger
	log = logger.NewLogger()
}

func main() {
	// Create a new context with dynamic key-value pairs
	contextData := map[string]interface{}{
		"header_id":  "12345",
		"user_id":    "67890",
		"request_id": "abc123",
	}

	// Create the context with a map that contains dynamic key-value pairs
	ctx := context.WithValue(context.Background(), "context_data", contextData)

	log.Add(slog.Int("status_code", 404)).Error(ctx, "404 Not Found: User missing")
	log.Add(slog.String("user_id", "67890")).Info(ctx, "User found")
	log.Info(context.Background(), "Hello World")
	log.Error(ctx, "Error occurred")
	log.Debug(ctx, "Debug message")
}
