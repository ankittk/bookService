package server

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ankittk/bookService/internal/config"
	"github.com/ankittk/bookService/internal/middleware"
	"github.com/ankittk/bookService/pkg/logger"
	bs "github.com/ankittk/bookService/proto/gen"
)

type GRPCServer struct {
	server *grpc.Server
	addr   string
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer(cfg *config.Config) *GRPCServer {
	// Create a new gRPC server with a custom interceptor
	// The interceptor is used for logging, authentication, etc.
	// In this case, we are using a middleware that limits the rate of requests
	// to the server. The rate limit is set to 100 requests per 5 minutes.
	// The middleware is implemented in the middleware package.
	// The interceptor is a function that takes a context, a request, and a response writer
	// and returns an error. It is used to intercept the request before it reaches the handler.

	s := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.NewMiddleware().Interceptor),
	)

	// Register the BookServiceServer with the gRPC server
	bs.RegisterBookServiceServer(s, NewBookServiceServer())

	// Register reflection service on the gRPC server
	// it allows client to discover the services for tools like grpcurl
	reflection.Register(s)

	return &GRPCServer{
		server: s,
		addr:   fmt.Sprintf(":%s", cfg.GRPCPort),
	}
}

// Start starts the gRPC server
func (g *GRPCServer) Start(ctx context.Context) {
	lis, err := net.Listen("tcp", g.addr)
	if err != nil {
		logger.Error(ctx, errors.Wrap(err, "failed to listen"))
	}

	logger.Info(ctx, fmt.Sprintf("Starting gRPC server at: %s", g.addr))
	if err := g.server.Serve(lis); err != nil {
		logger.Error(ctx, errors.Wrap(err, "failed to serve gRPC server"))
	}
}

// Stop gracefully stops the gRPC server
func (g *GRPCServer) Stop(ctx context.Context) {
	g.server.GracefulStop()
	logger.Info(ctx, "gRPC server stopped gracefully")
}
