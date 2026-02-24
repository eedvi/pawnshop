package logger

import (
	"context"

	"github.com/rs/zerolog"
)

type contextKey string

const (
	// RequestIDKey is the context key for request ID
	RequestIDKey contextKey = "request_id"
	// UserIDKey is the context key for user ID
	UserIDKey contextKey = "user_id"
)

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID retrieves the user ID from context
func GetUserID(ctx context.Context) int64 {
	if userID, ok := ctx.Value(UserIDKey).(int64); ok {
		return userID
	}
	return 0
}

// FromContext creates a logger with context values
func FromContext(ctx context.Context, logger zerolog.Logger) zerolog.Logger {
	contextLogger := logger

	if requestID := GetRequestID(ctx); requestID != "" {
		contextLogger = contextLogger.With().Str("request_id", requestID).Logger()
	}

	if userID := GetUserID(ctx); userID != 0 {
		contextLogger = contextLogger.With().Int64("user_id", userID).Logger()
	}

	return contextLogger
}
