package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// OIDCAuthenticator implements OIDC authentication
type OIDCAuthenticator struct {
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	clientID     string
	clientSecret string
	redirectURL  string
	oauth2Config *oauth2.Config
}

// NewOIDCAuthenticator creates a new OIDC authenticator
func NewOIDCAuthenticator(issuerURL, clientID, clientSecret, redirectURL string) (*OIDCAuthenticator, error) {
	provider, err := oidc.NewProvider(context.Background(), issuerURL)
	if err != nil {
		return nil, err
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &OIDCAuthenticator{
		provider:     provider,
		verifier:     provider.Verifier(&oidc.Config{ClientID: clientID}),
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		oauth2Config: config,
	}, nil
}

// Authenticate implements the Authenticator interface for OIDC
func (a *OIDCAuthenticator) Authenticate(r *http.Request) (*User, error) {
	// Get the authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, ErrNoAuthHeader
	}

	// Check if it's a Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, ErrInvalidAuthHeader
	}

	// Verify the token
	token, err := a.verifier.Verify(r.Context(), parts[1])
	if err != nil {
		return nil, err
	}

	// Extract claims
	var claims struct {
		Sub           string   `json:"sub"`
		Email         string   `json:"email"`
		Name          string   `json:"name"`
		Groups        []string `json:"groups"`
		EmailVerified bool     `json:"email_verified"`
	}

	if err := token.Claims(&claims); err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		ID:         claims.Sub,
		Email:      claims.Email,
		Name:       claims.Name,
		Roles:      claims.Groups,
		AuthMethod: "oidc",
		Claims:     make(map[string]interface{}),
	}

	// Add all claims to the user's claims map
	if err := token.Claims(&user.Claims); err != nil {
		return nil, err
	}

	return user, nil
}

// LoginURL returns the URL to redirect to for login
func (a *OIDCAuthenticator) LoginURL(state string) string {
	return a.oauth2Config.AuthCodeURL(state)
}

// ExchangeCode exchanges an authorization code for tokens
func (a *OIDCAuthenticator) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return a.oauth2Config.Exchange(ctx, code)
}

// Error definitions
var (
	ErrNoAuthHeader      = &AuthError{"no authorization header"}
	ErrInvalidAuthHeader = &AuthError{"invalid authorization header"}
)

// AuthError represents an authentication error
type AuthError struct {
	message string
}

func (e *AuthError) Error() string {
	return e.message
}
