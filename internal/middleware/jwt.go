package middleware

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthenticator handles JWT token validation and claim extraction
type JWTAuthenticator struct {
	// PublicKey is the RSA public key used to verify tokens
	PublicKey *rsa.PublicKey
	// JWKSURL is the URL to fetch public keys from (optional)
	JWKSURL string
	// Issuer is the expected token issuer
	Issuer string
	// Audience is the expected token audience
	Audience string
	// CustomValidators are additional claim validators
	CustomValidators []func(claims jwt.MapClaims) error
}

// NewJWTAuthenticator creates a new JWT authenticator
func NewJWTAuthenticator(publicKey *rsa.PublicKey, issuer, audience string) *JWTAuthenticator {
	return &JWTAuthenticator{
		PublicKey: publicKey,
		Issuer:    issuer,
		Audience:  audience,
	}
}

// WithJWKSURL sets the JWKS URL for dynamic key fetching
func (a *JWTAuthenticator) WithJWKSURL(url string) *JWTAuthenticator {
	a.JWKSURL = url
	return a
}

// WithCustomValidator adds a custom claim validator
func (a *JWTAuthenticator) WithCustomValidator(validator func(claims jwt.MapClaims) error) *JWTAuthenticator {
	a.CustomValidators = append(a.CustomValidators, validator)
	return a
}

// Authenticate validates the JWT token and extracts claims
func (a *JWTAuthenticator) Authenticate(r *http.Request) (map[string]interface{}, error) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	tokenString := parts[1]

	// Parse and validate token
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(a.Issuer),
		jwt.WithAudience(a.Audience),
	)

	var keyFunc jwt.Keyfunc
	if a.JWKSURL != "" {
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			// TODO: Implement JWKS key fetching
			return a.PublicKey, nil
		}
	} else {
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			return a.PublicKey, nil
		}
	}

	token, err := parser.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Run custom validators
	for _, validator := range a.CustomValidators {
		if err := validator(claims); err != nil {
			return nil, fmt.Errorf("custom validation failed: %w", err)
		}
	}

	// Convert claims to map[string]interface{}
	result := make(map[string]interface{})
	for k, v := range claims {
		result[k] = v
	}

	return result, nil
}

// ExtractClaim extracts a specific claim from the token
func (a *JWTAuthenticator) ExtractClaim(r *http.Request, claimName string) (interface{}, error) {
	claims, err := a.Authenticate(r)
	if err != nil {
		return nil, err
	}

	value, ok := claims[claimName]
	if !ok {
		return nil, fmt.Errorf("claim %s not found", claimName)
	}

	return value, nil
}

// ExtractStringClaim extracts a string claim from the token
func (a *JWTAuthenticator) ExtractStringClaim(r *http.Request, claimName string) (string, error) {
	value, err := a.ExtractClaim(r, claimName)
	if err != nil {
		return "", err
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("claim %s is not a string", claimName)
	}

	return str, nil
}

// ExtractTimeClaim extracts a time claim from the token
func (a *JWTAuthenticator) ExtractTimeClaim(r *http.Request, claimName string) (time.Time, error) {
	value, err := a.ExtractClaim(r, claimName)
	if err != nil {
		return time.Time{}, err
	}

	switch v := value.(type) {
	case float64:
		return time.Unix(int64(v), 0), nil
	case int64:
		return time.Unix(v, 0), nil
	default:
		return time.Time{}, fmt.Errorf("claim %s is not a valid timestamp", claimName)
	}
}

// Common claim validators
func ValidateExpiration(claims jwt.MapClaims) error {
	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("exp claim missing or invalid")
	}

	if time.Unix(int64(exp), 0).Before(time.Now()) {
		return errors.New("token expired")
	}

	return nil
}

func ValidateIssuer(expectedIssuer string) func(claims jwt.MapClaims) error {
	return func(claims jwt.MapClaims) error {
		iss, ok := claims["iss"].(string)
		if !ok {
			return errors.New("iss claim missing or invalid")
		}

		if iss != expectedIssuer {
			return fmt.Errorf("invalid issuer: got %s, want %s", iss, expectedIssuer)
		}

		return nil
	}
}

func ValidateAudience(expectedAudience string) func(claims jwt.MapClaims) error {
	return func(claims jwt.MapClaims) error {
		aud, ok := claims["aud"].(string)
		if !ok {
			return errors.New("aud claim missing or invalid")
		}

		if aud != expectedAudience {
			return fmt.Errorf("invalid audience: got %s, want %s", aud, expectedAudience)
		}

		return nil
	}
}
