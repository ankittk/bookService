package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ankittk/bookService/pkg/logger"
)

// RateLimiter is a simple rate limiter that allows a certain number of requests
type RateLimiter struct {
	limit     int
	counter   int
	mu        sync.Mutex
	resetTime time.Time
}

// NewRateLimiter creates a new RateLimiter with the specified limit
func NewRateLimiter(limit int, resetInterval time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:     limit,
		counter:   0,
		resetTime: time.Now().Add(resetInterval),
	}
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Reset the counter if the reset time has passed
	if time.Now().After(r.resetTime) {
		r.counter = 0
		r.resetTime = time.Now().Add(time.Minute * 5) // Reset interval
	}

	// Check if the request is allowed
	if r.counter < r.limit {
		r.counter++
		return true
	}
	return false
}

// Middleware struct that holds the rate limiter and implements the gRPC interceptor
type Middleware struct {
	rateLimiter *RateLimiter
}

// NewMiddleware creates a new Middleware instance with the specified rate limit
func NewMiddleware() *Middleware {
	return &Middleware{
		rateLimiter: NewRateLimiter(100, time.Minute*5), // 100 requests per 5 minutes
	}
}

// Interceptor checks the rate limit and proceeds with the handler if allowed
func (m *Middleware) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	// Log the request metadata (for debugging or monitoring)
	md, _ := metadata.FromIncomingContext(ctx)
	for key, values := range md {
		logger.Info(ctx, fmt.Sprintf("Metadata: %s=%v", key, values))
	}
	// Step 1: Rate limiting check
	if !m.rateLimiter.Allow() {
		return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
	}

	// Step 2: Authentication check (if API key is provided in metadata)
	apiKey := md.Get("api-key") // Use lowercase for consistency, since gRPC metadata is case-insensitive
	if len(apiKey) == 0 || apiKey[0] != "expected_api_key" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid API key")
	}

	// Step 3: Call the actual handler if rate limit is not exceeded
	resp, err = handler(ctx, req)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Error handling request: %v", err))
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	// Step 4: Return the response
	return resp, err
}

// HTTPMiddleware checks the rate limit, authentication, and logs the request
func (m *Middleware) HTTPMiddleware(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(ctx, fmt.Sprintf("Request received: method=%s, headers=%v", r.Method, r.Header))

		// Step 1: Rate limiting check
		if !m.rateLimiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Step 2: Authentication check (if API-Key is provided in headers)
		apiKey := r.Header.Get("Api-Key") // canonical MIME-style casing
		if apiKey != "expected_api_key" {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		for k, v := range r.Header {
			logger.Info(ctx, fmt.Sprintf("Header %q: %v", k, v))
		}

		// Step 3: Attach api-key to gRPC metadata
		md := metadata.Pairs("api-key", apiKey)
		ctx = metadata.NewOutgoingContext(ctx, md)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
