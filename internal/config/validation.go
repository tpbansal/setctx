package config

import (
	"fmt"
	"os"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid configuration: %s: %s", e.Field, e.Message)
}

// ValidateAuthConfig validates the authentication configuration
func ValidateAuthConfig(cfg *AuthConfig) error {
	if err := validateJWTConfig(&cfg.JWT); err != nil {
		return err
	}

	if err := validateSessionConfig(&cfg.Session); err != nil {
		return err
	}

	if cfg.OIDC.Enabled {
		if err := validateOIDCConfig(&cfg.OIDC); err != nil {
			return err
		}
	}

	if cfg.SAML.Enabled {
		if err := validateSAMLConfig(&cfg.SAML); err != nil {
			return err
		}
	}

	return nil
}

// validateJWTConfig validates JWT configuration
func validateJWTConfig(cfg *JWTConfig) error {
	if cfg.PublicKeyPath == "" && cfg.JWKSURL == "" {
		return &ValidationError{
			Field:   "jwt",
			Message: "either publicKeyPath or jwksUrl must be specified",
		}
	}

	if cfg.PublicKeyPath != "" {
		if _, err := os.Stat(cfg.PublicKeyPath); os.IsNotExist(err) {
			return &ValidationError{
				Field:   "jwt.publicKeyPath",
				Message: fmt.Sprintf("public key file not found: %s", cfg.PublicKeyPath),
			}
		}
	}

	if cfg.Issuer == "" {
		return &ValidationError{
			Field:   "jwt.issuer",
			Message: "issuer must be specified",
		}
	}

	if cfg.Audience == "" {
		return &ValidationError{
			Field:   "jwt.audience",
			Message: "audience must be specified",
		}
	}

	if cfg.TokenExpiry <= 0 {
		return &ValidationError{
			Field:   "jwt.tokenExpiry",
			Message: "token expiry must be greater than 0",
		}
	}

	return nil
}

// validateSessionConfig validates session configuration
func validateSessionConfig(cfg *SessionConfig) error {
	if cfg.SessionTimeout <= 0 {
		return &ValidationError{
			Field:   "session.sessionTimeout",
			Message: "session timeout must be greater than 0",
		}
	}

	if cfg.MaxSessionsPerUser <= 0 {
		return &ValidationError{
			Field:   "session.maxSessionsPerUser",
			Message: "max sessions per user must be greater than 0",
		}
	}

	switch cfg.SessionStoreType {
	case "memory", "redis":
		// Valid store types
	default:
		return &ValidationError{
			Field:   "session.sessionStoreType",
			Message: "invalid session store type, must be 'memory' or 'redis'",
		}
	}

	if cfg.SessionStoreType == "redis" {
		if err := validateRedisConfig(&cfg.Redis); err != nil {
			return err
		}
	}

	return nil
}

// validateRedisConfig validates Redis configuration
func validateRedisConfig(cfg *RedisConfig) error {
	if cfg.Addr == "" {
		return &ValidationError{
			Field:   "session.redis.addr",
			Message: "Redis address must be specified",
		}
	}

	if cfg.DB < 0 {
		return &ValidationError{
			Field:   "session.redis.db",
			Message: "Redis database number must be non-negative",
		}
	}

	return nil
}

// validateOIDCConfig validates OIDC configuration
func validateOIDCConfig(cfg *OIDCConfig) error {
	if cfg.ProviderURL == "" {
		return &ValidationError{
			Field:   "oidc.providerUrl",
			Message: "OIDC provider URL must be specified",
		}
	}

	if cfg.ClientID == "" {
		return &ValidationError{
			Field:   "oidc.clientId",
			Message: "OIDC client ID must be specified",
		}
	}

	if cfg.ClientSecret == "" {
		return &ValidationError{
			Field:   "oidc.clientSecret",
			Message: "OIDC client secret must be specified",
		}
	}

	if cfg.RedirectURL == "" {
		return &ValidationError{
			Field:   "oidc.redirectUrl",
			Message: "OIDC redirect URL must be specified",
		}
	}

	if len(cfg.Scopes) == 0 {
		return &ValidationError{
			Field:   "oidc.scopes",
			Message: "at least one OIDC scope must be specified",
		}
	}

	return nil
}

// validateSAMLConfig validates SAML configuration
func validateSAMLConfig(cfg *SAMLConfig) error {
	if cfg.IDPMetadataURL == "" {
		return &ValidationError{
			Field:   "saml.idpMetadataUrl",
			Message: "SAML IDP metadata URL must be specified",
		}
	}

	if cfg.EntityID == "" {
		return &ValidationError{
			Field:   "saml.entityId",
			Message: "SAML entity ID must be specified",
		}
	}

	if cfg.ACSURL == "" {
		return &ValidationError{
			Field:   "saml.acsUrl",
			Message: "SAML ACS URL must be specified",
		}
	}

	return nil
}
