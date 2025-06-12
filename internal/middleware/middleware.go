package middleware

import "net/http"

// Middleware is a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// Chain creates a middleware chain that will be executed in order
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// RequestContext holds request-specific data that can be passed through middleware
type RequestContext struct {
	RequestID string
	UserID    string
	IP        string
	// Add other context fields as needed
}

// GetRequestContext retrieves the request context from the request
func GetRequestContext(r *http.Request) *RequestContext {
	if ctx := r.Context().Value(contextKey); ctx != nil {
		return ctx.(*RequestContext)
	}
	return &RequestContext{}
}

// contextKey is used to store the request context in the request context
type contextKey struct{}

var contextKeyValue = contextKey{}
