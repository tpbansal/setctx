package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore implements SessionStore using Redis
type RedisStore struct {
	client *redis.Client
	prefix string
}

// NewRedisStore creates a new Redis-based session store
func NewRedisStore(client *redis.Client, prefix string) *RedisStore {
	return &RedisStore{
		client: client,
		prefix: prefix,
	}
}

// Get retrieves a session by ID
func (s *RedisStore) Get(ctx context.Context, id string) (*Session, error) {
	key := s.getKey(id)
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// Create creates a new session
func (s *RedisStore) Create(ctx context.Context, session *Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	key := s.getKey(session.ID)
	expiration := time.Until(session.ExpiresAt)

	// Store session data
	if err := s.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return err
	}

	// Add to user's session list
	userKey := s.getUserKey(session.UserID)
	return s.client.SAdd(ctx, userKey, session.ID).Err()
}

// Update updates an existing session
func (s *RedisStore) Update(ctx context.Context, session *Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	key := s.getKey(session.ID)
	expiration := time.Until(session.ExpiresAt)

	return s.client.Set(ctx, key, data, expiration).Err()
}

// Delete deletes a session
func (s *RedisStore) Delete(ctx context.Context, id string) error {
	// Get session to remove from user's session list
	session, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	key := s.getKey(id)
	userKey := s.getUserKey(session.UserID)

	// Delete session data and remove from user's session list
	pipe := s.client.Pipeline()
	pipe.Del(ctx, key)
	pipe.SRem(ctx, userKey, id)

	_, err = pipe.Exec(ctx)
	return err
}

// ListByUserID lists all sessions for a user
func (s *RedisStore) ListByUserID(ctx context.Context, userID string) ([]*Session, error) {
	userKey := s.getUserKey(userID)
	sessionIDs, err := s.client.SMembers(ctx, userKey).Result()
	if err != nil {
		return nil, err
	}

	var sessions []*Session
	for _, id := range sessionIDs {
		session, err := s.Get(ctx, id)
		if err != nil {
			if err == ErrSessionNotFound {
				// Clean up stale session ID
				_ = s.client.SRem(ctx, userKey, id)
				continue
			}
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// getKey returns the Redis key for a session
func (s *RedisStore) getKey(id string) string {
	return s.prefix + ":session:" + id
}

// getUserKey returns the Redis key for a user's session list
func (s *RedisStore) getUserKey(userID string) string {
	return s.prefix + ":user:" + userID + ":sessions"
}
