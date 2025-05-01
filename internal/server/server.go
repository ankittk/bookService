package server

import (
	"context"

	bs "github.com/ankittk/bookService/proto/gen"
)

type Server struct {
	bs.UnimplementedBookServiceServer
}

func NewBookServiceServer() *Server {
	return &Server{}
}

func (s *Server) Check(ctx context.Context, req *bs.HealthCheckRequest) (*bs.HealthCheckResponse, error) {
	return &bs.HealthCheckResponse{
		Status:  "UP",
		Message: "Service is running",
	}, nil
}

func (s *Server) SearchBooks(ctx context.Context, req *bs.SearchBooksRequest) (*bs.SearchBooksResponse, error) {
	return nil, nil
}

func (s *Server) GetBookDetails(ctx context.Context, req *bs.GetBookDetailsRequest) (*bs.GetBookDetailsResponse, error) {
	return nil, nil
}

func (s *Server) AddBook(ctx context.Context, req *bs.AddBookRequest) (*bs.AddBookResponse, error) {
	return nil, nil
}

func (s *Server) UpdateBook(ctx context.Context, req *bs.UpdateBookRequest) (*bs.UpdateBookResponse, error) {
	return nil, nil
}

func (s *Server) DeleteBook(ctx context.Context, req *bs.DeleteBookRequest) (*bs.DeleteBookResponse, error) {
	return nil, nil
}
