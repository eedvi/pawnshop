package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"pawnshop/internal/domain"
)

// SettingRepository implements repository.SettingRepository
type SettingRepository struct {
	db *DB
}

// NewSettingRepository creates a new SettingRepository
func NewSettingRepository(db *DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// Get retrieves a setting by key and optional branch ID
func (r *SettingRepository) Get(ctx context.Context, key string, branchID *int64) (*domain.Setting, error) {
	var query string
	var args []interface{}

	if branchID != nil {
		// Try branch-specific setting first, then fall back to global
		query = `
			SELECT id, key, value, description, branch_id, created_at, updated_at
			FROM settings
			WHERE key = $1 AND (branch_id = $2 OR branch_id IS NULL)
			ORDER BY branch_id DESC NULLS LAST
			LIMIT 1
		`
		args = []interface{}{key, *branchID}
	} else {
		// Only get global setting
		query = `
			SELECT id, key, value, description, branch_id, created_at, updated_at
			FROM settings
			WHERE key = $1 AND branch_id IS NULL
		`
		args = []interface{}{key}
	}

	setting := &domain.Setting{}
	var valueJSON []byte
	var description sql.NullString
	var settingBranchID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&setting.ID,
		&setting.Key,
		&valueJSON,
		&description,
		&settingBranchID,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("setting not found")
		}
		return nil, fmt.Errorf("failed to get setting: %w", err)
	}

	// Parse value from JSON
	if err := json.Unmarshal(valueJSON, &setting.Value); err != nil {
		setting.Value = string(valueJSON) // Fall back to string
	}

	setting.Description = StringPtr(description)
	setting.BranchID = Int64Ptr(settingBranchID)

	return setting, nil
}

// GetAll retrieves all settings for a branch (including global settings)
func (r *SettingRepository) GetAll(ctx context.Context, branchID *int64) ([]*domain.Setting, error) {
	var query string
	var args []interface{}

	if branchID != nil {
		// Get branch-specific and global settings
		query = `
			SELECT id, key, value, description, branch_id, created_at, updated_at
			FROM settings
			WHERE branch_id = $1 OR branch_id IS NULL
			ORDER BY branch_id NULLS FIRST, key
		`
		args = []interface{}{*branchID}
	} else {
		// Only get global settings
		query = `
			SELECT id, key, value, description, branch_id, created_at, updated_at
			FROM settings
			WHERE branch_id IS NULL
			ORDER BY key
		`
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list settings: %w", err)
	}
	defer rows.Close()

	settings := []*domain.Setting{}
	for rows.Next() {
		setting := &domain.Setting{}
		var valueJSON []byte
		var description sql.NullString
		var settingBranchID sql.NullInt64

		err := rows.Scan(
			&setting.ID,
			&setting.Key,
			&valueJSON,
			&description,
			&settingBranchID,
			&setting.CreatedAt,
			&setting.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}

		// Parse value from JSON
		if err := json.Unmarshal(valueJSON, &setting.Value); err != nil {
			setting.Value = string(valueJSON)
		}

		setting.Description = StringPtr(description)
		setting.BranchID = Int64Ptr(settingBranchID)
		settings = append(settings, setting)
	}

	return settings, nil
}

// Set creates or updates a setting
func (r *SettingRepository) Set(ctx context.Context, setting *domain.Setting) error {
	// Convert value to JSON
	valueJSON, err := json.Marshal(setting.Value)
	if err != nil {
		return fmt.Errorf("failed to marshal setting value: %w", err)
	}

	query := `
		INSERT INTO settings (key, value, description, branch_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (key, COALESCE(branch_id, 0))
		DO UPDATE SET value = $2, description = $3, updated_at = NOW()
		RETURNING id, created_at, updated_at
	`

	err = r.db.QueryRowContext(ctx, query,
		setting.Key,
		valueJSON,
		NullString(setting.Description),
		NullInt64(setting.BranchID),
	).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to set setting: %w", err)
	}

	return nil
}

// Delete deletes a setting
func (r *SettingRepository) Delete(ctx context.Context, key string, branchID *int64) error {
	var query string
	var args []interface{}

	if branchID != nil {
		query = `DELETE FROM settings WHERE key = $1 AND branch_id = $2`
		args = []interface{}{key, *branchID}
	} else {
		query = `DELETE FROM settings WHERE key = $1 AND branch_id IS NULL`
		args = []interface{}{key}
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("setting not found")
	}

	return nil
}
