package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoggingMiddleware(t *testing.T) {
	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware chain
	middleware := Chain(LoggingMiddleware())

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Serve the request
	middleware(handler).ServeHTTP(rr, req)

	// Check response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		requests       int
		rate           float64
		capacity       float64
		window         time.Duration
		expectedStatus int
	}{
		{
			name:           "Within rate limit",
			requests:       5,
			rate:           10.0,
			capacity:       10.0,
			window:         time.Second,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Exceed rate limit",
			requests:       15,
			rate:           10.0,
			capacity:       10.0,
			window:         time.Second,
			expectedStatus: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create rate limiter
			limiter := NewRateLimiter(tt.rate, tt.capacity, tt.window)

			// Create middleware chain
			middleware := Chain(RateLimitMiddleware(limiter))

			// Create test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Make multiple requests
			for i := 0; i < tt.requests; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				rr := httptest.NewRecorder()

				// Serve the request
				middleware(handler).ServeHTTP(rr, req)

				// Check the last request's status
				if i == tt.requests-1 {
					if status := rr.Code; status != tt.expectedStatus {
						t.Errorf("handler returned wrong status code: got %v want %v",
							status, tt.expectedStatus)
					}
				}
			}
		})
	}
}

func TestMiddlewareChain(t *testing.T) {
	// Create test handlers that set headers
	logging := LoggingMiddleware()
	rateLimiter := NewRateLimiter(10.0, 10.0, time.Second)
	rateLimit := RateLimitMiddleware(rateLimiter)

	// Create final handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware chain
	middleware := Chain(logging, rateLimit)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Serve the request
	middleware(handler).ServeHTTP(rr, req)

	// Check response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
