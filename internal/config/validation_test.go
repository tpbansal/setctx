package config

import (
	"os"
	"testing"
	"time"
)

func TestValidateAuthConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *AuthConfig
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: &AuthConfig{
				JWT: JWTConfig{
					PublicKeyPath: "testdata/public.key",
					Issuer:        "test-issuer",
					Audience:      "test-audience",
					TokenExpiry:   1 * time.Hour,
				},
				Session: SessionConfig{
					SessionTimeout:     24 * time.Hour,
					MaxSessionsPerUser: 5,
					SessionStoreType:   "memory",
				},
			},
			wantErr: false,
		},
		{
			name: "missing JWT issuer",
			config: &AuthConfig{
				JWT: JWTConfig{
					PublicKeyPath: "testdata/public.key",
					Audience:      "test-audience",
					TokenExpiry:   1 * time.Hour,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid session timeout",
			config: &AuthConfig{
				JWT: JWTConfig{
					PublicKeyPath: "testdata/public.key",
					Issuer:        "test-issuer",
					Audience:      "test-audience",
					TokenExpiry:   1 * time.Hour,
				},
				Session: SessionConfig{
					SessionTimeout:     0,
					MaxSessionsPerUser: 5,
					SessionStoreType:   "memory",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid session store type",
			config: &AuthConfig{
				JWT: JWTConfig{
					PublicKeyPath: "testdata/public.key",
					Issuer:        "test-issuer",
					Audience:      "test-audience",
					TokenExpiry:   1 * time.Hour,
				},
				Session: SessionConfig{
					SessionTimeout:     24 * time.Hour,
					MaxSessionsPerUser: 5,
					SessionStoreType:   "invalid",
				},
			},
			wantErr: true,
		},
		{
			name: "missing Redis configuration",
			config: &AuthConfig{
				JWT: JWTConfig{
					PublicKeyPath: "testdata/public.key",
					Issuer:        "test-issuer",
					Audience:      "test-audience",
					TokenExpiry:   1 * time.Hour,
				},
				Session: SessionConfig{
					SessionTimeout:     24 * time.Hour,
					MaxSessionsPerUser: 5,
					SessionStoreType:   "redis",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid OIDC configuration",
			config: &AuthConfig{
				JWT: JWTConfig{
					PublicKeyPath: "testdata/public.key",
					Issuer:        "test-issuer",
					Audience:      "test-audience",
					TokenExpiry:   1 * time.Hour,
				},
				Session: SessionConfig{
					SessionTimeout:     24 * time.Hour,
					MaxSessionsPerUser: 5,
					SessionStoreType:   "memory",
				},
				OIDC: OIDCConfig{
					Enabled: true,
					// Missing required fields
				},
			},
			wantErr: true,
		},
	}

	// Create test public key file
	if err := os.MkdirAll("testdata", 0755); err != nil {
		t.Fatalf("failed to create testdata directory: %v", err)
	}
	if err := os.WriteFile("testdata/public.key", []byte("test public key"), 0644); err != nil {
		t.Fatalf("failed to create test public key file: %v", err)
	}
	defer os.RemoveAll("testdata")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuthConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAuthConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				t.Logf("Validation error: %v", err)
			}
		})
	}
}

func TestValidateRedisConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *RedisConfig
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: &RedisConfig{
				Addr:     "localhost:6379",
				Password: "password",
				DB:       0,
			},
			wantErr: false,
		},
		{
			name: "missing address",
			config: &RedisConfig{
				Password: "password",
				DB:       0,
			},
			wantErr: true,
		},
		{
			name: "invalid DB number",
			config: &RedisConfig{
				Addr:     "localhost:6379",
				Password: "password",
				DB:       -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRedisConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRedisConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
