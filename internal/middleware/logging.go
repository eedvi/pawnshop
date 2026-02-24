package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"pawnshop/pkg/logger"
)

// LoggingMiddleware handles request logging
type LoggingMiddleware struct {
	logger zerolog.Logger
}

// NewLoggingMiddleware creates a new LoggingMiddleware
func NewLoggingMiddleware(logger zerolog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

// Logger returns a middleware that logs requests
func (m *LoggingMiddleware) Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Generate request ID if not present
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("X-Request-ID", requestID)

		// Inject request ID into context for services
		ctx := logger.WithRequestID(c.Context(), requestID)
		c.SetUserContext(ctx)

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get user info if authenticated
		var userID int64
		if user := GetUser(c); user != nil {
			userID = user.ID
		}

		// Build log event
		event := m.logger.Info()
		if c.Response().StatusCode() >= 400 {
			event = m.logger.Warn()
		}
		if c.Response().StatusCode() >= 500 {
			event = m.logger.Error()
		}

		// Log the request
		event.
			Str("request_id", requestID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Int("status", c.Response().StatusCode()).
			Dur("latency", latency).
			Int("bytes", len(c.Response().Body())).
			Str("user_agent", c.Get("User-Agent"))

		if userID > 0 {
			event.Int64("user_id", userID)
		}

		if err != nil {
			event.Err(err)
		}

		event.Msg("request")

		return err
	}
}

// Recovery returns a middleware that recovers from panics
func (m *LoggingMiddleware) Recovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				requestID := c.Get("X-Request-ID")

				m.logger.Error().
					Str("request_id", requestID).
					Str("method", c.Method()).
					Str("path", c.Path()).
					Interface("panic", r).
					Msg("panic recovered")

				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"error": fiber.Map{
						"code":    "INTERNAL_ERROR",
						"message": "An internal server error occurred",
					},
				})
			}
		}()

		return c.Next()
	}
}
