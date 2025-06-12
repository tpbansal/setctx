package config

import (
	"crypto/rsa"
	"time"
)

// AuthConfig holds all authentication-related configuration
type AuthConfig struct {
	// JWT settings
	JWT JWTConfig `yaml:"jwt"`

	// Session settings
	Session SessionConfig `yaml:"session"`

	// OIDC settings
	OIDC OIDCConfig `yaml:"oidc"`

	// SAML settings
	SAML SAMLConfig `yaml:"saml"`
}

// JWTConfig holds JWT-specific configuration
type JWTConfig struct {
	// PublicKey is the RSA public key used to verify tokens
	PublicKey *rsa.PublicKey `yaml:"-"`

	// PublicKeyPath is the path to the public key file
	PublicKeyPath string `yaml:"publicKeyPath"`

	// JWKSURL is the URL to fetch public keys from
	JWKSURL string `yaml:"jwksUrl"`

	// Issuer is the expected token issuer
	Issuer string `yaml:"issuer"`

	// Audience is the expected token audience
	Audience string `yaml:"audience"`

	// TokenExpiry is the duration for which tokens are valid
	TokenExpiry time.Duration `yaml:"tokenExpiry"`
}

// SessionConfig holds session-specific configuration
type SessionConfig struct {
	// SessionTimeout is the duration after which a session expires
	SessionTimeout time.Duration `yaml:"sessionTimeout"`

	// MaxSessionsPerUser is the maximum number of concurrent sessions allowed per user
	MaxSessionsPerUser int `yaml:"maxSessionsPerUser"`

	// SessionStoreType defines where sessions are stored (memory, redis, etc.)
	SessionStoreType string `yaml:"sessionStoreType"`

	// Redis settings (if using Redis as session store)
	Redis RedisConfig `yaml:"redis"`
}

// RedisConfig holds Redis-specific configuration
type RedisConfig struct {
	// Addr is the Redis server address
	Addr string `yaml:"addr"`

	// Password is the Redis password
	Password string `yaml:"password"`

	// DB is the Redis database number
	DB int `yaml:"db"`
}

// OIDCConfig holds OpenID Connect configuration
type OIDCConfig struct {
	// Enabled indicates if OIDC authentication is enabled
	Enabled bool `yaml:"enabled"`

	// ProviderURL is the OIDC provider URL
	ProviderURL string `yaml:"providerUrl"`

	// ClientID is the OIDC client ID
	ClientID string `yaml:"clientId"`

	// ClientSecret is the OIDC client secret
	ClientSecret string `yaml:"clientSecret"`

	// RedirectURL is the OIDC redirect URL
	RedirectURL string `yaml:"redirectUrl"`

	// Scopes are the OIDC scopes to request
	Scopes []string `yaml:"scopes"`
}

// SAMLConfig holds SAML-specific configuration
type SAMLConfig struct {
	// Enabled indicates if SAML authentication is enabled
	Enabled bool `yaml:"enabled"`

	// IDPMetadataURL is the URL to fetch IDP metadata from
	IDPMetadataURL string `yaml:"idpMetadataUrl"`

	// EntityID is the SAML entity ID
	EntityID string `yaml:"entityId"`

	// ACSURL is the Assertion Consumer Service URL
	ACSURL string `yaml:"acsUrl"`

	// SLOURL is the Single Logout URL
	SLOURL string `yaml:"sloUrl"`
}

// DefaultAuthConfig returns a default configuration
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWT: JWTConfig{
			TokenExpiry: 1 * time.Hour,
		},
		Session: SessionConfig{
			SessionTimeout:     24 * time.Hour,
			MaxSessionsPerUser: 5,
			SessionStoreType:   "memory",
		},
		OIDC: OIDCConfig{
			Enabled: false,
			Scopes:  []string{"openid", "profile", "email"},
		},
		SAML: SAMLConfig{
			Enabled: false,
		},
	}
}
