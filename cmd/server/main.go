package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/ankittk/bookService/internal/config"
	"github.com/ankittk/bookService/internal/server"
	"github.com/ankittk/bookService/pkg/logger"
)

func init() {
	// Load environment variables from .env file
	_ = godotenv.Load(".env")
}

func main() {
	appCfg := config.NewDefaultConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcServer := server.NewGRPCServer(appCfg)
	httpServer := server.NewHTTPServer(ctx, appCfg)

	// Using a goroutine allows the main application to remain responsive to the shutdown signal
	// while the server continues to serve grpc requests.
	// If the server runs in the main thread, it will block further execution.
	go grpcServer.Start(ctx)
	go httpServer.Start(ctx)

	// Handle shutdown
	stop := make(chan os.Signal, 1)
	// Use os/signal to listen for interrupt signals (like Ctrl+C)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for the interrupt signal to gracefully shut down the server
	<-stop

	logger.Info(ctx, "Shutting down servers...")

	// When shutting down, it's essential to give enough time for the server to finish ongoing requests
	// before shutting down completely.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcServer.Stop(ctx)
	httpServer.Stop(shutdownCtx)

	logger.Info(shutdownCtx, "Servers shut down gracefully")
}
