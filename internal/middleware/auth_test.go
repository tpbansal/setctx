package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "No auth header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid auth header",
			authHeader:     "Invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid auth header",
			authHeader:     "Bearer valid-token",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock authenticator
			auth := &mockAuthenticator{
				shouldAuthenticate: tt.authHeader == "Bearer valid-token",
			}

			// Create middleware chain
			middleware := Chain(AuthMiddleware(auth))

			// Create test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Get user from context
				user, ok := GetUserFromContext(r)
				if !ok {
					t.Error("User not found in context")
					return
				}

				// Check user details
				if user.ID != "test-user" {
					t.Errorf("Unexpected user ID: got %v want %v", user.ID, "test-user")
				}
			})

			// Create test request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			// Serve the request
			middleware(handler).ServeHTTP(rr, req)

			// Check response
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

// mockAuthenticator is a test implementation of the Authenticator interface
type mockAuthenticator struct {
	shouldAuthenticate bool
}

func (m *mockAuthenticator) Authenticate(r *http.Request) (*User, error) {
	if !m.shouldAuthenticate {
		return nil, ErrNoAuthHeader
	}

	return &User{
		ID:         "test-user",
		Email:      "test@example.com",
		Name:       "Test User",
		Roles:      []string{"user"},
		AuthMethod: "mock",
		Claims:     make(map[string]interface{}),
	}, nil
}
