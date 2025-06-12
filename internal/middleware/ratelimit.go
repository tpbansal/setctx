package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu          sync.Mutex
	rate        float64        // tokens per second
	capacity    float64        // maximum tokens
	tokens      float64        // current tokens
	lastUpdated time.Time      // last token update time
	ipLimits    map[string]int // per-IP request counts
	window      time.Duration  // time window for IP limits
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate float64, capacity float64, window time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:        rate,
		capacity:    capacity,
		tokens:      capacity,
		lastUpdated: time.Now(),
		ipLimits:    make(map[string]int),
		window:      window,
	}
}

// RateLimitMiddleware creates a middleware that rate limits requests
func RateLimitMiddleware(limiter *RateLimiter) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			ip := r.RemoteAddr

			// Check if the request should be rate limited
			if !limiter.Allow(ip) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Process the request
			next.ServeHTTP(w, r)
		})
	}
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Update tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(rl.lastUpdated).Seconds()
	rl.tokens = min(rl.capacity, rl.tokens+elapsed*rl.rate)
	rl.lastUpdated = now

	// Check global rate limit
	if rl.tokens < 1 {
		return false
	}

	// Check IP-based rate limit
	count := rl.ipLimits[ip]
	if count >= 100 { // Example limit: 100 requests per window
		return false
	}

	// Update counters
	rl.tokens--
	rl.ipLimits[ip]++

	// Clean up old IP entries
	go rl.cleanup()

	return true
}

// cleanup removes old IP entries
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Remove entries older than the window
	for ip := range rl.ipLimits {
		delete(rl.ipLimits, ip)
	}
}

// min returns the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
