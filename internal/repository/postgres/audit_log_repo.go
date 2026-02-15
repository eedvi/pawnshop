package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// AuditLogRepository implements repository.AuditLogRepository
type AuditLogRepository struct {
	db *DB
}

// NewAuditLogRepository creates a new AuditLogRepository
func NewAuditLogRepository(db *DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

// Create creates a new audit log entry
func (r *AuditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (branch_id, user_id, action, entity_type, entity_id, old_values, new_values, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	var oldValuesJSON, newValuesJSON []byte
	var err error

	if log.OldValues != nil {
		oldValuesJSON, err = json.Marshal(log.OldValues)
		if err != nil {
			return fmt.Errorf("failed to marshal old values: %w", err)
		}
	}
	if log.NewValues != nil {
		newValuesJSON, err = json.Marshal(log.NewValues)
		if err != nil {
			return fmt.Errorf("failed to marshal new values: %w", err)
		}
	}

	err = r.db.QueryRowContext(ctx, query,
		NullInt64(log.BranchID),
		NullInt64(log.UserID),
		log.Action,
		log.EntityType,
		NullInt64(log.EntityID),
		oldValuesJSON,
		newValuesJSON,
		NullString(log.IPAddress),
		NullString(log.UserAgent),
	).Scan(&log.ID, &log.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// List retrieves audit logs with filters
func (r *AuditLogRepository) List(ctx context.Context, params repository.AuditLogListParams) (*repository.PaginatedResult[domain.AuditLog], error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 50
	}

	// Build WHERE clause
	where := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if params.BranchID != nil {
		where += fmt.Sprintf(" AND branch_id = $%d", argNum)
		args = append(args, *params.BranchID)
		argNum++
	}

	if params.UserID != nil {
		where += fmt.Sprintf(" AND user_id = $%d", argNum)
		args = append(args, *params.UserID)
		argNum++
	}

	if params.Action != "" {
		where += fmt.Sprintf(" AND action = $%d", argNum)
		args = append(args, params.Action)
		argNum++
	}

	if params.EntityType != "" {
		where += fmt.Sprintf(" AND entity_type = $%d", argNum)
		args = append(args, params.EntityType)
		argNum++
	}

	if params.EntityID != nil {
		where += fmt.Sprintf(" AND entity_id = $%d", argNum)
		args = append(args, *params.EntityID)
		argNum++
	}

	if params.DateFrom != nil {
		where += fmt.Sprintf(" AND created_at >= $%d", argNum)
		args = append(args, *params.DateFrom)
		argNum++
	}

	if params.DateTo != nil {
		where += fmt.Sprintf(" AND created_at <= $%d", argNum)
		args = append(args, *params.DateTo)
		argNum++
	}

	// Count total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs %s", where)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Get data
	offset := (params.Page - 1) * params.PerPage
	query := fmt.Sprintf(`
		SELECT id, branch_id, user_id, action, entity_type, entity_id, old_values, new_values, ip_address, user_agent, created_at
		FROM audit_logs
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argNum, argNum+1)

	args = append(args, params.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	logs := []domain.AuditLog{}
	for rows.Next() {
		var log domain.AuditLog
		var branchID, userID, entityID sql.NullInt64
		var oldValues, newValues []byte
		var ipAddress, userAgent sql.NullString

		err := rows.Scan(
			&log.ID,
			&branchID,
			&userID,
			&log.Action,
			&log.EntityType,
			&entityID,
			&oldValues,
			&newValues,
			&ipAddress,
			&userAgent,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		log.BranchID = Int64Ptr(branchID)
		log.UserID = Int64Ptr(userID)
		log.EntityID = Int64Ptr(entityID)
		log.IPAddress = StringPtr(ipAddress)
		log.UserAgent = StringPtr(userAgent)

		if oldValues != nil {
			var v interface{}
			if err := json.Unmarshal(oldValues, &v); err == nil {
				log.OldValues = v
			}
		}
		if newValues != nil {
			var v interface{}
			if err := json.Unmarshal(newValues, &v); err == nil {
				log.NewValues = v
			}
		}

		logs = append(logs, log)
	}

	totalPages := total / params.PerPage
	if total%params.PerPage > 0 {
		totalPages++
	}

	return &repository.PaginatedResult[domain.AuditLog]{
		Data:       logs,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}
