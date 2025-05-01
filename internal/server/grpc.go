package server

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ankittk/bookService/internal/config"
	log "github.com/ankittk/bookService/pkg/logger"
	bs "github.com/ankittk/bookService/proto/gen"
)

type GRPCServer struct {
	server *grpc.Server
	addr   string
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer(cfg *config.Config) *GRPCServer {
	// Create a gRPC server
	s := grpc.NewServer()

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
		log.NewLogger().Error(ctx, errors.Wrap(err, "failed to listen"))
	}

	log.NewLogger().Info(ctx, fmt.Sprintf("Starting gRPC server at: %s", g.addr))
	if err := g.server.Serve(lis); err != nil {
		log.NewLogger().Error(ctx, errors.Wrap(err, "failed to serve gRPC server"))
	}
}

// Stop gracefully stops the gRPC server
func (g *GRPCServer) Stop(ctx context.Context) {
	g.server.GracefulStop()
	log.NewLogger().Info(ctx, "gRPC server stopped gracefully")
}
