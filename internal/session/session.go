package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"

	"safectx/internal/config"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
	ErrMaxSessions     = errors.New("maximum number of sessions reached")
)

// Session represents a user session
type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
	Data      map[string]interface{}
}

// SessionStore defines the interface for session storage
type SessionStore interface {
	// Get retrieves a session by ID
	Get(ctx context.Context, id string) (*Session, error)

	// Create creates a new session
	Create(ctx context.Context, session *Session) error

	// Update updates an existing session
	Update(ctx context.Context, session *Session) error

	// Delete deletes a session
	Delete(ctx context.Context, id string) error

	// ListByUserID lists all sessions for a user
	ListByUserID(ctx context.Context, userID string) ([]*Session, error)
}

// Manager handles session operations
type Manager struct {
	store  SessionStore
	config *config.SessionConfig
	mu     sync.RWMutex
}

// NewManager creates a new session manager
func NewManager(store SessionStore, config *config.SessionConfig) *Manager {
	return &Manager{
		store:  store,
		config: config,
	}
}

// CreateSession creates a new session for a user
func (m *Manager) CreateSession(ctx context.Context, userID string) (*Session, error) {
	// Check if user has reached max sessions
	sessions, err := m.store.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	if len(sessions) >= m.config.MaxSessionsPerUser {
		return nil, ErrMaxSessions
	}

	// Generate session ID
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Create new session
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(m.config.SessionTimeout),
		Data:      make(map[string]interface{}),
	}

	if err := m.store.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// GetSession retrieves a session by ID
func (m *Manager) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	session, err := m.store.Get(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		_ = m.store.Delete(ctx, sessionID)
		return nil, ErrSessionExpired
	}

	return session, nil
}

// UpdateSession updates a session
func (m *Manager) UpdateSession(ctx context.Context, session *Session) error {
	return m.store.Update(ctx, session)
}

// DeleteSession deletes a session
func (m *Manager) DeleteSession(ctx context.Context, sessionID string) error {
	return m.store.Delete(ctx, sessionID)
}

// ListUserSessions lists all sessions for a user
func (m *Manager) ListUserSessions(ctx context.Context, userID string) ([]*Session, error) {
	return m.store.ListByUserID(ctx, userID)
}

// RefreshSession extends the session expiration time
func (m *Manager) RefreshSession(ctx context.Context, sessionID string) error {
	session, err := m.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(m.config.SessionTimeout)
	return m.UpdateSession(ctx, session)
}

// generateSessionID generates a random session ID
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
