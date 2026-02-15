package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/zulfikawr/vault/internal/errors"
)

type RateLimiter struct {
	limit    int           // requests per interval
	interval time.Duration // interval
	clients  sync.Map
}

type client struct {
	mu         sync.Mutex
	tokens     int
	lastRefill time.Time
}

func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limit:    limit,
		interval: interval,
	}

	// Cleanup routine to remove old clients
	go func() {
		for range time.Tick(time.Minute) {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *RateLimiter) cleanup() {
	rl.clients.Range(func(key, value any) bool {
		c := value.(*client)
		c.mu.Lock()
		defer c.mu.Unlock()

		if time.Since(c.lastRefill) > rl.interval*2 {
			rl.clients.Delete(key)
		}
		return true
	})
}

func (rl *RateLimiter) Allow(ip string) bool {
	if rl.limit <= 0 {
		return true
	}

	// Fast path check if exists
	val, ok := rl.clients.Load(ip)
	if !ok {
		// Create new client
		newClient := &client{
			tokens:     rl.limit,
			lastRefill: time.Now(),
		}
		val, _ = rl.clients.LoadOrStore(ip, newClient)
	}

	c := val.(*client)
	c.mu.Lock()
	defer c.mu.Unlock()

	// Refill logic
	now := time.Now()
	if now.Sub(c.lastRefill) >= rl.interval {
		c.tokens = rl.limit
		c.lastRefill = now
	}

	if c.tokens > 0 {
		c.tokens--
		return true
	}

	return false
}

func RateLimitMiddleware(limit int) func(http.Handler) http.Handler {
	rl := NewRateLimiter(limit, time.Minute)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			// Handle X-Forwarded-For if behind proxy
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			}

			if !rl.Allow(ip) {
				errors.SendError(w, errors.NewError(http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Too many requests"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
