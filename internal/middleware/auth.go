package middleware

import (
	"context"
	"net/http"
)

// Authenticator defines the interface for authentication methods
type Authenticator interface {
	// Authenticate validates the request and returns the authenticated user
	Authenticate(r *http.Request) (*User, error)
}

// User represents an authenticated user
type User struct {
	ID         string
	Email      string
	Name       string
	Roles      []string
	Claims     map[string]interface{}
	AuthMethod string // "oidc" or "saml"
}

// AuthMiddleware creates a middleware that authenticates requests
func AuthMiddleware(auth Authenticator) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authenticate the request
			user, err := auth.Authenticate(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add user to request context
			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves the authenticated user from the request context
func GetUserFromContext(r *http.Request) (*User, bool) {
	user, ok := r.Context().Value(userContextKey).(*User)
	return user, ok
}

type userContextKeyType struct{}

var userContextKey = userContextKeyType{}
