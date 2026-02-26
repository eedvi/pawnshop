package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
	"pawnshop/pkg/logger"
)

// LoggingDB wraps sql.DB to add query logging
type LoggingDB struct {
	*DB
	slowQueryThreshold time.Duration
	logAllQueries      bool
}

// NewLoggingDB creates a new LoggingDB with the specified slow query threshold
func NewLoggingDB(db *DB, slowQueryThreshold time.Duration, logAllQueries bool) *LoggingDB {
	return &LoggingDB{
		DB:                 db,
		slowQueryThreshold: slowQueryThreshold,
		logAllQueries:      logAllQueries,
	}
}

// QueryContext executes a query and logs if it's slow or if logging all queries
func (db *LoggingDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := db.DB.QueryContext(ctx, query, args...)
	duration := time.Since(start)

	db.logQuery(ctx, query, duration, err, "query")

	return rows, err
}

// QueryRowContext executes a query that returns a single row
func (db *LoggingDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := db.DB.QueryRowContext(ctx, query, args...)
	duration := time.Since(start)

	db.logQuery(ctx, query, duration, nil, "query_row")

	return row
}

// ExecContext executes a query that doesn't return rows
func (db *LoggingDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := db.DB.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	db.logQuery(ctx, query, duration, err, "exec")

	return result, err
}

// logQuery logs the query if it's slow or if logging all queries
func (db *LoggingDB) logQuery(ctx context.Context, query string, duration time.Duration, err error, queryType string) {
	requestID := logger.GetRequestID(ctx)
	userID := logger.GetUserID(ctx)

	// Sanitizar query
	sanitizedQuery := logger.SanitizeSQL(query)

	// Log errores siempre
	if err != nil {
		logEvent := log.Error().
			Err(err).
			Str("query", sanitizedQuery).
			Int64("duration_ms", duration.Milliseconds()).
			Str("query_type", queryType)

		if requestID != "" {
			logEvent.Str("request_id", requestID)
		}
		if userID != 0 {
			logEvent.Int64("user_id", userID)
		}

		logEvent.Msg("Database query error")
		return
	}

	// Log queries lentas
	if duration > db.slowQueryThreshold {
		logEvent := log.Warn().
			Str("query", sanitizedQuery).
			Int64("duration_ms", duration.Milliseconds()).
			Dur("threshold", db.slowQueryThreshold).
			Str("query_type", queryType)

		if requestID != "" {
			logEvent.Str("request_id", requestID)
		}
		if userID != 0 {
			logEvent.Int64("user_id", userID)
		}

		logEvent.Msg("Slow database query detected")
		return
	}

	// Log todas las queries si está habilitado (modo debug)
	if db.logAllQueries {
		logEvent := log.Debug().
			Str("query", sanitizedQuery).
			Int64("duration_ms", duration.Milliseconds()).
			Str("query_type", queryType)

		if requestID != "" {
			logEvent.Str("request_id", requestID)
		}
		if userID != 0 {
			logEvent.Int64("user_id", userID)
		}

		logEvent.Msg("Database query executed")
	}
}

// BeginTx starts a transaction with logging
func (db *LoggingDB) BeginTx(ctx context.Context) (*Tx, error) {
	start := time.Now()
	tx, err := db.DB.BeginTx(ctx)
	duration := time.Since(start)

	requestID := logger.GetRequestID(ctx)
	logEvent := log.Debug().
		Int64("duration_ms", duration.Milliseconds())

	if requestID != "" {
		logEvent.Str("request_id", requestID)
	}

	if err != nil {
		logEvent.Err(err).Msg("Failed to begin transaction")
	} else {
		logEvent.Msg("Transaction started")
	}

	return tx, err
}
