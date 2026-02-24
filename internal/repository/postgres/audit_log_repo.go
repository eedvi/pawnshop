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
		INSERT INTO audit_logs (branch_id, user_id, action, entity_type, entity_id, description, old_values, new_values, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
		NullStringPtr(log.Description),
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
		where += fmt.Sprintf(" AND al.branch_id = $%d", argNum)
		args = append(args, *params.BranchID)
		argNum++
	}

	if params.UserID != nil {
		where += fmt.Sprintf(" AND al.user_id = $%d", argNum)
		args = append(args, *params.UserID)
		argNum++
	}

	if params.Action != "" {
		where += fmt.Sprintf(" AND al.action = $%d", argNum)
		args = append(args, params.Action)
		argNum++
	}

	if params.EntityType != "" {
		where += fmt.Sprintf(" AND al.entity_type = $%d", argNum)
		args = append(args, params.EntityType)
		argNum++
	}

	if params.EntityID != nil {
		where += fmt.Sprintf(" AND al.entity_id = $%d", argNum)
		args = append(args, *params.EntityID)
		argNum++
	}

	if params.DateFrom != nil {
		where += fmt.Sprintf(" AND al.created_at >= $%d", argNum)
		args = append(args, *params.DateFrom)
		argNum++
	}

	if params.DateTo != nil {
		where += fmt.Sprintf(" AND al.created_at <= $%d", argNum)
		args = append(args, *params.DateTo)
		argNum++
	}

	// Count total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs al %s", where)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Get data
	offset := (params.Page - 1) * params.PerPage
	query := fmt.Sprintf(`
		SELECT
			al.id, al.branch_id, al.user_id, al.action, al.entity_type, al.entity_id, al.description,
			al.old_values, al.new_values, al.ip_address, al.user_agent, al.created_at,
			COALESCE(u.first_name || ' ' || u.last_name, '') as user_name,
			COALESCE(b.name, '') as branch_name
		FROM audit_logs al
		LEFT JOIN users u ON u.id = al.user_id AND u.deleted_at IS NULL
		LEFT JOIN branches b ON b.id = al.branch_id AND b.deleted_at IS NULL
		%s
		ORDER BY al.created_at DESC
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
		var ipAddress, userAgent, userName, branchName, description sql.NullString

		err := rows.Scan(
			&log.ID,
			&branchID,
			&userID,
			&log.Action,
			&log.EntityType,
			&entityID,
			&description,
			&oldValues,
			&newValues,
			&ipAddress,
			&userAgent,
			&log.CreatedAt,
			&userName,
			&branchName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		log.BranchID = Int64Ptr(branchID)
		log.UserID = Int64Ptr(userID)
		log.EntityID = Int64Ptr(entityID)
		log.Description = StringPtrVal(description)
		log.IPAddress = StringPtr(ipAddress)
		log.UserAgent = StringPtr(userAgent)
		log.UserName = StringPtrVal(userName)
		log.BranchName = StringPtrVal(branchName)

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

// AuditStats represents audit statistics
type AuditStats struct {
	TotalActions     int                 `json:"total_actions"`
	ActionsByType    map[string]int      `json:"actions_by_type"`
	ActionsByEntity  map[string]int      `json:"actions_by_entity"`
	ActiveUsers      int                 `json:"active_users"`
	TopUsers         []TopUserStat       `json:"top_users"`
	RecentCritical   []domain.AuditLog   `json:"recent_critical"`
}

// TopUserStat represents a user's activity count
type TopUserStat struct {
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	Count     int    `json:"count"`
}

// GetStats retrieves audit statistics
func (r *AuditLogRepository) GetStats(ctx context.Context, params repository.AuditLogListParams) (interface{}, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if params.BranchID != nil {
		where += fmt.Sprintf(" AND branch_id = $%d", argNum)
		args = append(args, *params.BranchID)
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

	stats := &AuditStats{
		ActionsByType:   make(map[string]int),
		ActionsByEntity: make(map[string]int),
		TopUsers:        []TopUserStat{},
		RecentCritical:  []domain.AuditLog{},
	}

	// Get total actions
	totalQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs %s", where)
	if err := r.db.QueryRowContext(ctx, totalQuery, args...).Scan(&stats.TotalActions); err != nil {
		return nil, fmt.Errorf("failed to count total actions: %w", err)
	}

	// Get actions by type
	actionQuery := fmt.Sprintf(`
		SELECT action, COUNT(*) as count
		FROM audit_logs %s
		GROUP BY action
		ORDER BY count DESC
	`, where)
	rows, err := r.db.QueryContext(ctx, actionQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get actions by type: %w", err)
	}
	for rows.Next() {
		var action string
		var count int
		if err := rows.Scan(&action, &count); err != nil {
			rows.Close()
			return nil, err
		}
		stats.ActionsByType[action] = count
	}
	rows.Close()

	// Get actions by entity
	entityQuery := fmt.Sprintf(`
		SELECT entity_type, COUNT(*) as count
		FROM audit_logs %s
		GROUP BY entity_type
		ORDER BY count DESC
	`, where)
	rows, err = r.db.QueryContext(ctx, entityQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get actions by entity: %w", err)
	}
	for rows.Next() {
		var entityType string
		var count int
		if err := rows.Scan(&entityType, &count); err != nil {
			rows.Close()
			return nil, err
		}
		stats.ActionsByEntity[entityType] = count
	}
	rows.Close()

	// Get active users count
	activeUsersQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT user_id) FROM audit_logs %s AND user_id IS NOT NULL
	`, where)
	if err := r.db.QueryRowContext(ctx, activeUsersQuery, args...).Scan(&stats.ActiveUsers); err != nil {
		return nil, fmt.Errorf("failed to count active users: %w", err)
	}

	// Get top 5 users by activity
	topUsersQuery := fmt.Sprintf(`
		SELECT a.user_id, COALESCE(u.first_name || ' ' || u.last_name, 'Usuario Desconocido') as user_name, COUNT(*) as count
		FROM audit_logs a
		LEFT JOIN users u ON u.id = a.user_id
		%s AND a.user_id IS NOT NULL
		GROUP BY a.user_id, u.first_name, u.last_name
		ORDER BY count DESC
		LIMIT 5
	`, where)
	rows, err = r.db.QueryContext(ctx, topUsersQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get top users: %w", err)
	}
	for rows.Next() {
		var userStat TopUserStat
		if err := rows.Scan(&userStat.UserID, &userStat.UserName, &userStat.Count); err != nil {
			rows.Close()
			return nil, err
		}
		stats.TopUsers = append(stats.TopUsers, userStat)
	}
	rows.Close()

	// Get recent critical actions (delete, reject)
	criticalQuery := fmt.Sprintf(`
		SELECT id, branch_id, user_id, action, entity_type, entity_id, old_values, new_values, ip_address, user_agent, created_at
		FROM audit_logs
		%s AND action IN ('delete', 'reject')
		ORDER BY created_at DESC
		LIMIT 10
	`, where)
	rows, err = r.db.QueryContext(ctx, criticalQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent critical actions: %w", err)
	}
	defer rows.Close()

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
			return nil, err
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

		stats.RecentCritical = append(stats.RecentCritical, log)
	}

	return stats, nil
}
