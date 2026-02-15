package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"pawnshop/pkg/response"
)

// RateLimitConfig contains rate limit configuration
type RateLimitConfig struct {
	Max        int           // Maximum number of requests
	Window     time.Duration // Time window for rate limiting
	KeyFunc    func(*fiber.Ctx) string // Function to get the key for rate limiting
}

// DefaultRateLimitConfig returns the default rate limit configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Max:    100, // 100 requests
		Window: time.Minute, // per minute
		KeyFunc: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}
}

// LoginRateLimitConfig returns rate limit config for login endpoint
func LoginRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Max:    5, // 5 attempts
		Window: 15 * time.Minute, // per 15 minutes
		KeyFunc: func(c *fiber.Ctx) string {
			return "login:" + c.IP()
		},
	}
}

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	config  RateLimitConfig
	storage sync.Map
}

type rateLimitEntry struct {
	count     int
	expiresAt time.Time
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{config: config}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := rl.config.KeyFunc(c)
		now := time.Now()

		// Get or create entry
		val, _ := rl.storage.LoadOrStore(key, &rateLimitEntry{
			count:     0,
			expiresAt: now.Add(rl.config.Window),
		})
		entry := val.(*rateLimitEntry)

		// Reset if expired
		if now.After(entry.expiresAt) {
			entry.count = 0
			entry.expiresAt = now.Add(rl.config.Window)
		}

		// Increment count
		entry.count++

		// Set rate limit headers
		remaining := rl.config.Max - entry.count
		if remaining < 0 {
			remaining = 0
		}
		c.Set("X-RateLimit-Limit", strconv.Itoa(rl.config.Max))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Set("X-RateLimit-Reset", entry.expiresAt.Format(time.RFC3339))

		// Check if limit exceeded
		if entry.count > rl.config.Max {
			retryAfter := int(entry.expiresAt.Sub(now).Seconds())
			c.Set("Retry-After", strconv.Itoa(retryAfter))
			return response.TooManyRequests(c, "")
		}

		return c.Next()
	}
}

// cleanup periodically removes expired entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		now := time.Now()
		rl.storage.Range(func(key, value interface{}) bool {
			entry := value.(*rateLimitEntry)
			if now.After(entry.expiresAt) {
				rl.storage.Delete(key)
			}
			return true
		})
	}
}
