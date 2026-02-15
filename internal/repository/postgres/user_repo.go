package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// UserRepository implements repository.UserRepository
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT
			id, branch_id, role_id, email, password_hash,
			first_name, last_name, phone, avatar_url,
			is_active, email_verified,
			failed_login_attempts, locked_until, password_changed_at,
			last_login_at, last_login_ip,
			two_factor_enabled, two_factor_secret,
			created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := &domain.User{}
	var lockedUntil, passwordChangedAt, lastLoginAt, deletedAt sql.NullTime
	var branchIDNull sql.NullInt64
	var phone, avatarURL, lastLoginIP, twoFactorSecret sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &branchIDNull, &user.RoleID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &phone, &avatarURL,
		&user.IsActive, &user.EmailVerified,
		&user.FailedLoginAttempts, &lockedUntil, &passwordChangedAt,
		&lastLoginAt, &lastLoginIP,
		&user.TwoFactorEnabled, &twoFactorSecret,
		&user.CreatedAt, &user.UpdatedAt, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.BranchID = Int64Ptr(branchIDNull)
	user.Phone = StringPtr(phone)
	user.AvatarURL = StringPtr(avatarURL)
	user.LockedUntil = TimePtr(lockedUntil)
	user.PasswordChangedAt = TimePtr(passwordChangedAt)
	user.LastLoginAt = TimePtr(lastLoginAt)
	user.LastLoginIP = StringPtr(lastLoginIP)
	user.TwoFactorSecret = StringPtr(twoFactorSecret)
	user.DeletedAt = TimePtr(deletedAt)

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT
			id, branch_id, role_id, email, password_hash,
			first_name, last_name, phone, avatar_url,
			is_active, email_verified,
			failed_login_attempts, locked_until, password_changed_at,
			last_login_at, last_login_ip,
			two_factor_enabled, two_factor_secret,
			created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	user := &domain.User{}
	var branchIDNull sql.NullInt64
	var lockedUntil, passwordChangedAt, lastLoginAt, deletedAt sql.NullTime
	var phone, avatarURL, lastLoginIP, twoFactorSecret sql.NullString

	err := r.db.QueryRowContext(ctx, query, strings.ToLower(email)).Scan(
		&user.ID, &branchIDNull, &user.RoleID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &phone, &avatarURL,
		&user.IsActive, &user.EmailVerified,
		&user.FailedLoginAttempts, &lockedUntil, &passwordChangedAt,
		&lastLoginAt, &lastLoginIP,
		&user.TwoFactorEnabled, &twoFactorSecret,
		&user.CreatedAt, &user.UpdatedAt, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.BranchID = Int64Ptr(branchIDNull)
	user.Phone = StringPtr(phone)
	user.AvatarURL = StringPtr(avatarURL)
	user.LockedUntil = TimePtr(lockedUntil)
	user.PasswordChangedAt = TimePtr(passwordChangedAt)
	user.LastLoginAt = TimePtr(lastLoginAt)
	user.LastLoginIP = StringPtr(lastLoginIP)
	user.TwoFactorSecret = StringPtr(twoFactorSecret)
	user.DeletedAt = TimePtr(deletedAt)

	return user, nil
}

// List retrieves users with pagination and filters
func (r *UserRepository) List(ctx context.Context, params repository.UserListParams) (*repository.PaginatedResult[domain.User], error) {
	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 20
	}

	// Build query
	baseQuery := `FROM users WHERE deleted_at IS NULL`
	args := []interface{}{}
	argCount := 0

	if params.BranchID != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND branch_id = $%d", argCount)
		args = append(args, *params.BranchID)
	}

	if params.RoleID != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND role_id = $%d", argCount)
		args = append(args, *params.RoleID)
	}

	if params.IsActive != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *params.IsActive)
	}

	if params.Search != "" {
		argCount++
		baseQuery += fmt.Sprintf(" AND (first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+params.Search+"%")
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Get data with pagination
	orderBy := "created_at"
	if params.OrderBy != "" {
		orderBy = params.OrderBy
	}
	order := "DESC"
	if params.Order == "asc" {
		order = "ASC"
	}

	offset := (params.Page - 1) * params.PerPage
	dataQuery := fmt.Sprintf(`
		SELECT
			id, branch_id, role_id, email, password_hash,
			first_name, last_name, phone, avatar_url,
			is_active, email_verified,
			failed_login_attempts, locked_until, password_changed_at,
			last_login_at, last_login_ip,
			two_factor_enabled, two_factor_secret,
			created_at, updated_at, deleted_at
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		baseQuery, orderBy, order, argCount+1, argCount+2,
	)
	args = append(args, params.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var user domain.User
		var branchIDNull sql.NullInt64
		var lockedUntil, passwordChangedAt, lastLoginAt, deletedAt sql.NullTime
		var phone, avatarURL, lastLoginIP, twoFactorSecret sql.NullString

		err := rows.Scan(
			&user.ID, &branchIDNull, &user.RoleID, &user.Email, &user.PasswordHash,
			&user.FirstName, &user.LastName, &phone, &avatarURL,
			&user.IsActive, &user.EmailVerified,
			&user.FailedLoginAttempts, &lockedUntil, &passwordChangedAt,
			&lastLoginAt, &lastLoginIP,
			&user.TwoFactorEnabled, &twoFactorSecret,
			&user.CreatedAt, &user.UpdatedAt, &deletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		user.BranchID = Int64Ptr(branchIDNull)
		user.Phone = StringPtr(phone)
		user.AvatarURL = StringPtr(avatarURL)
		user.LockedUntil = TimePtr(lockedUntil)
		user.PasswordChangedAt = TimePtr(passwordChangedAt)
		user.LastLoginAt = TimePtr(lastLoginAt)
		user.LastLoginIP = StringPtr(lastLoginIP)
		user.TwoFactorSecret = StringPtr(twoFactorSecret)
		user.DeletedAt = TimePtr(deletedAt)

		users = append(users, user)
	}

	totalPages := total / params.PerPage
	if total%params.PerPage > 0 {
		totalPages++
	}

	return &repository.PaginatedResult[domain.User]{
		Data:       users,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			branch_id, role_id, email, password_hash,
			first_name, last_name, phone, avatar_url,
			is_active, email_verified
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		NullInt64(user.BranchID),
		user.RoleID,
		strings.ToLower(user.Email),
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		NullString(user.Phone),
		NullString(user.AvatarURL),
		user.IsActive,
		user.EmailVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET
			branch_id = $2,
			role_id = $3,
			email = $4,
			first_name = $5,
			last_name = $6,
			phone = $7,
			avatar_url = $8,
			is_active = $9,
			email_verified = $10,
			two_factor_enabled = $11,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		NullInt64(user.BranchID),
		user.RoleID,
		strings.ToLower(user.Email),
		user.FirstName,
		user.LastName,
		NullString(user.Phone),
		NullString(user.AvatarURL),
		user.IsActive,
		user.EmailVerified,
		user.TwoFactorEnabled,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdatePassword updates the user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	query := `
		UPDATE users SET
			password_hash = $2,
			password_changed_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateLastLogin updates the last login info
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id int64, ip string) error {
	query := `
		UPDATE users SET
			last_login_at = NOW(),
			last_login_ip = $2,
			failed_login_attempts = 0,
			locked_until = NULL,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, id, ip)
	return err
}

// IncrementFailedLogins increments failed login attempts
func (r *UserRepository) IncrementFailedLogins(ctx context.Context, id int64) error {
	query := `
		UPDATE users SET
			failed_login_attempts = failed_login_attempts + 1,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ResetFailedLogins resets failed login attempts
func (r *UserRepository) ResetFailedLogins(ctx context.Context, id int64) error {
	query := `
		UPDATE users SET
			failed_login_attempts = 0,
			locked_until = NULL,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// LockUser locks the user until a specified time
func (r *UserRepository) LockUser(ctx context.Context, id int64, until *int64) error {
	var lockedUntil *time.Time
	if until != nil {
		t := time.Now().Add(time.Duration(*until) * time.Minute)
		lockedUntil = &t
	}

	query := `
		UPDATE users SET
			locked_until = $2,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, id, NullTime(lockedUntil))
	return err
}
