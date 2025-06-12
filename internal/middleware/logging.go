package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware creates a middleware that logs request details
func LoggingMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a custom response writer to capture the status code
			lrw := &loggingResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Process the request
			next.ServeHTTP(lrw, r)

			// Log the request details
			duration := time.Since(start)
			log.Printf(
				"%s %s %s %d %v",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				lrw.statusCode,
				duration,
			)
		})
	}
}

// loggingResponseWriter wraps http.ResponseWriter to capture the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before writing it
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
