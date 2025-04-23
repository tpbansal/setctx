package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"safectx/internal/contextfilter"
	"safectx/internal/policy"
	"safectx/pkg/schema"
)

// WithValidation ensures that the incoming request is valid.
func WithValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example validation step (you can extend this)
		var req schema.MCPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		// Pass the validated request into the context
		ctx := context.WithValue(r.Context(), "MCPRequest", req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithPolicy checks the policy for the request
func WithPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the request from context
		req, ok := r.Context().Value("MCPRequest").(schema.MCPRequest)
		if !ok {
			http.Error(w, "Request not found", http.StatusInternalServerError)
			return
		}

		// Evaluate the policy
		allowed, err := policy.Evaluate(&req)
		if err != nil || !allowed {
			http.Error(w, "Request denied by policy", http.StatusForbidden)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// WithRedaction redacts sensitive information from the request
func WithRedaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the request from context
		req, ok := r.Context().Value("MCPRequest").(schema.MCPRequest)
		if !ok {
			http.Error(w, "Request not found", http.StatusInternalServerError)
			return
		}

		// Redact sensitive data
		contextfilter.Redact(&req)

		// Pass the redacted request into the context
		ctx := context.WithValue(r.Context(), "MCPRequest", req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithLogging logs the request
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// WithRecovery ensures the server does not crash on errors.
func WithRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
