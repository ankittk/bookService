package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ankittk/bookService/internal/config"
	"github.com/ankittk/bookService/internal/middleware"
	"github.com/ankittk/bookService/pkg/logger"
	bs "github.com/ankittk/bookService/proto/gen"
)

type HTTPServer struct {
	server *http.Server
	addr   string
}

// NewHTTPServer creates a new HTTP server instance
func NewHTTPServer(ctx context.Context, cfg *config.Config) *HTTPServer {
	// Create a new gRPC-Gateway mux
	// the mux is an HTTP request multiplexer that matches incoming HTTP requests to the appropriate gRPC service
	// it maps HTTP requests to gRPC methods based on annotation option (google.api.http)
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
	)

	// grpc.DialOption is used to configure the gRPC client side behavior when connecting to a gRPC server
	// since grpc-gateway is acting as a client to the gRPC server, we need to specify the connection options
	// so we are connecting to the gRPC server using insecure credentials
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register the BookServiceServer with the gRPC-Gateway mux
	// this takes any HTTP requests on these endpoints and forward them to the gRPC server
	err := bs.RegisterBookServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%s", cfg.GRPCPort), opts)
	if err != nil {
		logger.Error(ctx, errors.Wrap(err, "failed to register BookServiceHandler"))
		return nil
	}

	// middleware is used to add custom behavior to the HTTP server
	// in this case, we are using a middleware that limits the rate of requests to the server
	// the rate limit is set to 100 requests per 5 minutes

	mw := middleware.NewMiddleware()

	// wrap the gRPC-Gateway original mux with the middleware
	handler := mw.HTTPMiddleware(ctx, mux)

	// Create a new HTTP server which performs rate limiting before forwarding the request to the gRPC server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler: handler,
	}

	logger.Info(ctx, fmt.Sprintf("Starting HTTP server at: %s", srv.Addr))
	return &HTTPServer{
		server: srv,
	}
}

// Start starts the HTTP server
func (h *HTTPServer) Start(ctx context.Context) {
	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, errors.Wrap(err, "failed to start HTTP server"))
	}
}

// Stop gracefully stops the HTTP server
func (h *HTTPServer) Stop(ctx context.Context) {
	if err := h.server.Shutdown(ctx); err != nil {
		logger.Error(ctx, errors.Wrap(err, "failed to shutdown HTTP server"))
	}
}

// In server/server.go (HTTP server)
func customHeaderMatcher(key string) (string, bool) {
	if strings.EqualFold(key, "API-KEY") {
		logger.Info(context.Background(), "Mapping API-KEY to api-key")
		return "api-key", true
	}
	return runtime.DefaultHeaderMatcher(key)
}
