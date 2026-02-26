package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"pawnshop/pkg/logger"
)

// RequestLoggingConfig configures request logging behavior
type RequestLoggingConfig struct {
	// LogRequestBody enables request body logging (development only)
	LogRequestBody bool

	// LogResponseBody enables response body logging (development only)
	LogResponseBody bool

	// MaxBodySize is the maximum body size to log (bytes)
	MaxBodySize int

	// SanitizeSensitive enables automatic sanitization of sensitive data
	SanitizeSensitive bool

	// ExcludePaths paths to exclude from detailed logging
	ExcludePaths []string

	// SampleRate for production (0.0 to 1.0, where 1.0 = log all)
	SampleRate float64
}

// DefaultRequestLoggingConfig returns safe defaults for production
func DefaultRequestLoggingConfig() RequestLoggingConfig {
	return RequestLoggingConfig{
		LogRequestBody:    false, // NEVER in production
		LogResponseBody:   false, // NEVER in production
		MaxBodySize:       1024,  // 1KB max
		SanitizeSensitive: true,  // ALWAYS sanitize
		ExcludePaths: []string{
			"/health",
			"/metrics",
			"/swagger.json",
		},
		SampleRate: 1.0, // Log all by default
	}
}

// DevelopmentRequestLoggingConfig returns config for development
func DevelopmentRequestLoggingConfig() RequestLoggingConfig {
	return RequestLoggingConfig{
		LogRequestBody:    true,  // OK in development
		LogResponseBody:   true,  // OK in development
		MaxBodySize:       10240, // 10KB
		SanitizeSensitive: true,  // ALWAYS sanitize
		ExcludePaths:      []string{"/health", "/metrics"},
		SampleRate:        1.0,
	}
}

// RequestLoggingMiddleware creates a middleware for detailed request logging
func RequestLoggingMiddleware(config RequestLoggingConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if path should be excluded
		path := c.Path()
		for _, excludePath := range config.ExcludePaths {
			if strings.HasPrefix(path, excludePath) {
				return c.Next()
			}
		}

		start := time.Now()
		requestID := c.Get("X-Request-ID")

		// Get user info
		var userID int64
		if user := GetUser(c); user != nil {
			userID = user.ID
		}

		// Build request log event
		logEvent := log.Info().
			Str("type", "http_request_start").
			Str("request_id", requestID).
			Str("method", c.Method()).
			Str("path", path).
			Str("protocol", c.Protocol()).
			Str("client_ip", c.IP()).
			Int("content_length", len(c.Body()))

		if userID > 0 {
			logEvent.Int64("user_id", userID)
		}

		// Log request headers (sanitized)
		logEvent.Interface("headers", sanitizeHeaders(c.GetReqHeaders()))

		// Log request body if enabled (development only)
		if config.LogRequestBody && len(c.Body()) > 0 && len(c.Body()) <= config.MaxBodySize {
			bodyStr := string(c.Body())
			if config.SanitizeSensitive {
				bodyStr = logger.SanitizeJSON(bodyStr)
			}
			logEvent.Str("request_body", bodyStr)
		} else if len(c.Body()) > 0 {
			logEvent.Int("request_body_size", len(c.Body()))
			logEvent.Str("content_type", c.Get("Content-Type"))
		}

		logEvent.Msg("Request started")

		// Process request
		err := c.Next()

		// Get response body if enabled
		var responseBody []byte
		if config.LogResponseBody {
			responseBody = c.Response().Body()
		}

		// Log response
		duration := time.Since(start)
		statusCode := c.Response().StatusCode()

		respLogEvent := log.Info()
		if statusCode >= 400 && statusCode < 500 {
			respLogEvent = log.Warn()
		} else if statusCode >= 500 {
			respLogEvent = log.Error()
		}

		respLogEvent.
			Str("type", "http_request_complete").
			Str("request_id", requestID).
			Str("method", c.Method()).
			Str("path", path).
			Int("status_code", statusCode).
			Dur("duration", duration).
			Float64("duration_ms", float64(duration.Milliseconds())).
			Int("response_size", len(c.Response().Body()))

		if userID > 0 {
			respLogEvent.Int64("user_id", userID)
		}

		// Log response body if enabled and not too large
		if config.LogResponseBody && len(responseBody) > 0 && len(responseBody) <= config.MaxBodySize {
			bodyStr := string(responseBody)
			if config.SanitizeSensitive {
				bodyStr = logger.SanitizeJSON(bodyStr)
			}
			respLogEvent.Str("response_body", bodyStr)
		}

		if err != nil {
			respLogEvent.Err(err)
		}

		respLogEvent.Msg("Request completed")

		return err
	}
}

// sanitizeHeaders removes sensitive headers
func sanitizeHeaders(headers map[string][]string) map[string]interface{} {
	sanitized := make(map[string]interface{})

	// Sensitive headers to redact
	sensitiveHeaders := map[string]bool{
		"authorization":  true,
		"cookie":         true,
		"set-cookie":     true,
		"x-api-key":      true,
		"x-auth-token":   true,
		"proxy-authorization": true,
	}

	for key, values := range headers {
		lowerKey := strings.ToLower(key)

		if sensitiveHeaders[lowerKey] {
			sanitized[key] = "***REDACTED***"
		} else if len(values) == 1 {
			sanitized[key] = values[0]
		} else {
			sanitized[key] = values
		}
	}

	return sanitized
}
