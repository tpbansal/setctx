package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"safectx/internal/config"
	"safectx/internal/session"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Create default configuration
	cfg := config.DefaultAuthConfig()

	// Configure JWT settings
	cfg.JWT.PublicKeyPath = "testdata/public.key"
	cfg.JWT.Issuer = "example-issuer"
	cfg.JWT.Audience = "example-audience"
	cfg.JWT.TokenExpiry = 1 * time.Hour

	// Configure session settings
	cfg.Session.SessionTimeout = 24 * time.Hour
	cfg.Session.MaxSessionsPerUser = 5
	cfg.Session.SessionStoreType = "redis"
	cfg.Session.Redis.Addr = "localhost:6379"

	// Validate configuration
	if err := config.ValidateAuthConfig(cfg); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Session.Redis.Addr,
		Password: cfg.Session.Redis.Password,
		DB:       cfg.Session.Redis.DB,
	})

	// Initialize session store and manager
	sessionStore := session.NewRedisStore(redisClient, "example")
	sessionManager := session.NewManager(sessionStore, &cfg.Session)

	// Create HTTP handlers
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// In a real application, you would validate credentials here
		userID := r.FormValue("user_id")
		if userID == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		// Create a new session
		session, err := sessionManager.CreateSession(r.Context(), userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create session: %v", err), http.StatusInternalServerError)
			return
		}

		// Store some session data
		session.Data["last_login"] = time.Now()
		if err := sessionManager.UpdateSession(r.Context(), session); err != nil {
			http.Error(w, fmt.Sprintf("failed to update session: %v", err), http.StatusInternalServerError)
			return
		}

		// Return session ID to client
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"session_id": session.ID,
		})
	})

	http.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("X-Session-ID")
		if sessionID == "" {
			http.Error(w, "X-Session-ID header is required", http.StatusBadRequest)
			return
		}

		// Get session
		session, err := sessionManager.GetSession(r.Context(), sessionID)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid session: %v", err), http.StatusUnauthorized)
			return
		}

		// Refresh session
		if err := sessionManager.RefreshSession(r.Context(), sessionID); err != nil {
			http.Error(w, fmt.Sprintf("failed to refresh session: %v", err), http.StatusInternalServerError)
			return
		}

		// Return session info
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(session)
	})

	http.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		userID := r.FormValue("user_id")
		if userID == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		// List user's sessions
		sessions, err := sessionManager.ListUserSessions(r.Context(), userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to list sessions: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sessions)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("X-Session-ID")
		if sessionID == "" {
			http.Error(w, "X-Session-ID header is required", http.StatusBadRequest)
			return
		}

		// Delete session
		if err := sessionManager.DeleteSession(r.Context(), sessionID); err != nil {
			http.Error(w, fmt.Sprintf("failed to delete session: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// Start server
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
