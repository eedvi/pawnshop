package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"pawnshop/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB wraps the database connection
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(cfg *config.DatabaseConfig) (*DB, error) {
	db, err := sql.Open("pgx", cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// BeginTx starts a new transaction
func (db *DB) BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

// Tx wraps a database transaction
type Tx struct {
	*sql.Tx
}

// Commit commits the transaction
func (tx *Tx) Commit() error {
	return tx.Tx.Commit()
}

// Rollback rolls back the transaction
func (tx *Tx) Rollback() error {
	return tx.Tx.Rollback()
}

// Querier interface for both DB and Tx
type Querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// NullInt64 converts *int64 to sql.NullInt64
func NullInt64(v *int64) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *v, Valid: true}
}

// NullString converts string to sql.NullString
func NullString(v string) sql.NullString {
	if v == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: v, Valid: true}
}

// NullFloat64 converts *float64 to sql.NullFloat64
func NullFloat64(v *float64) sql.NullFloat64 {
	if v == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *v, Valid: true}
}

// NullTime converts *time.Time to sql.NullTime
func NullTime(v *time.Time) sql.NullTime {
	if v == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *v, Valid: true}
}

// Int64Ptr converts sql.NullInt64 to *int64
func Int64Ptr(v sql.NullInt64) *int64 {
	if !v.Valid {
		return nil
	}
	return &v.Int64
}

// StringPtr converts sql.NullString to string (returns empty string if null)
func StringPtr(v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	return v.String
}

// StringPtrVal converts sql.NullString to *string (returns nil if null)
func StringPtrVal(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}

// NullStringPtr converts *string to sql.NullString
func NullStringPtr(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *v, Valid: true}
}

// Float64Ptr converts sql.NullFloat64 to *float64
func Float64Ptr(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}
	return &v.Float64
}

// TimePtr converts sql.NullTime to *time.Time
func TimePtr(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	return &v.Time
}

// IntPtr converts sql.NullInt64 to *int
func IntPtr(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}
	i := int(v.Int64)
	return &i
}

// NullIntPtr converts *int to sql.NullInt64
func NullIntPtr(v *int) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*v), Valid: true}
}
