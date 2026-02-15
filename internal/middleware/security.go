package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders returns a middleware that sets security headers
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Set("X-Frame-Options", "DENY")

		// XSS protection (legacy browsers)
		c.Set("X-XSS-Protection", "1; mode=block")

		// HSTS (HTTP Strict Transport Security)
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content Security Policy
		c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'")

		// Referrer Policy
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		return c.Next()
	}
}

// RequestID returns a middleware that ensures a request ID is present
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = c.GetRespHeader("X-Request-ID")
		}
		if requestID != "" {
			c.Set("X-Request-ID", requestID)
		}
		return c.Next()
	}
}
