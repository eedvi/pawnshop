package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// BranchRepository implements repository.BranchRepository
type BranchRepository struct {
	db *DB
}

// NewBranchRepository creates a new BranchRepository
func NewBranchRepository(db *DB) *BranchRepository {
	return &BranchRepository{db: db}
}

// GetByID retrieves a branch by ID
func (r *BranchRepository) GetByID(ctx context.Context, id int64) (*domain.Branch, error) {
	query := `
		SELECT id, name, code, address, phone, email, is_active, timezone, currency,
			   default_interest_rate, default_loan_term_days, default_grace_period,
			   created_at, updated_at, deleted_at
		FROM branches
		WHERE id = $1 AND deleted_at IS NULL
	`

	branch := &domain.Branch{}
	var address, phone, email sql.NullString
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&branch.ID, &branch.Name, &branch.Code, &address, &phone, &email,
		&branch.IsActive, &branch.Timezone, &branch.Currency,
		&branch.DefaultInterestRate, &branch.DefaultLoanTermDays, &branch.DefaultGracePeriod,
		&branch.CreatedAt, &branch.UpdatedAt, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("branch not found")
		}
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	branch.Address = StringPtr(address)
	branch.Phone = StringPtr(phone)
	branch.Email = StringPtr(email)
	branch.DeletedAt = TimePtr(deletedAt)

	return branch, nil
}

// GetByCode retrieves a branch by code
func (r *BranchRepository) GetByCode(ctx context.Context, code string) (*domain.Branch, error) {
	query := `
		SELECT id, name, code, address, phone, email, is_active, timezone, currency,
			   default_interest_rate, default_loan_term_days, default_grace_period,
			   created_at, updated_at, deleted_at
		FROM branches
		WHERE code = $1 AND deleted_at IS NULL
	`

	branch := &domain.Branch{}
	var address, phone, email sql.NullString
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&branch.ID, &branch.Name, &branch.Code, &address, &phone, &email,
		&branch.IsActive, &branch.Timezone, &branch.Currency,
		&branch.DefaultInterestRate, &branch.DefaultLoanTermDays, &branch.DefaultGracePeriod,
		&branch.CreatedAt, &branch.UpdatedAt, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("branch not found")
		}
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	branch.Address = StringPtr(address)
	branch.Phone = StringPtr(phone)
	branch.Email = StringPtr(email)
	branch.DeletedAt = TimePtr(deletedAt)

	return branch, nil
}

// List retrieves branches with pagination
func (r *BranchRepository) List(ctx context.Context, params repository.PaginationParams) (*repository.PaginatedResult[domain.Branch], error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 20
	}

	// Count total
	var total int
	countQuery := `SELECT COUNT(*) FROM branches WHERE deleted_at IS NULL`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count branches: %w", err)
	}

	// Get data
	orderBy := "created_at"
	if params.OrderBy != "" {
		orderBy = params.OrderBy
	}
	order := "DESC"
	if params.Order == "asc" {
		order = "ASC"
	}

	offset := (params.Page - 1) * params.PerPage
	query := fmt.Sprintf(`
		SELECT id, name, code, address, phone, email, is_active, timezone, currency,
			   default_interest_rate, default_loan_term_days, default_grace_period,
			   created_at, updated_at, deleted_at
		FROM branches
		WHERE deleted_at IS NULL
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, orderBy, order)

	rows, err := r.db.QueryContext(ctx, query, params.PerPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}
	defer rows.Close()

	branches := []domain.Branch{}
	for rows.Next() {
		var branch domain.Branch
		var address, phone, email sql.NullString
		var deletedAt sql.NullTime

		err := rows.Scan(
			&branch.ID, &branch.Name, &branch.Code, &address, &phone, &email,
			&branch.IsActive, &branch.Timezone, &branch.Currency,
			&branch.DefaultInterestRate, &branch.DefaultLoanTermDays, &branch.DefaultGracePeriod,
			&branch.CreatedAt, &branch.UpdatedAt, &deletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan branch: %w", err)
		}

		branch.Address = StringPtr(address)
		branch.Phone = StringPtr(phone)
		branch.Email = StringPtr(email)
		branch.DeletedAt = TimePtr(deletedAt)

		branches = append(branches, branch)
	}

	totalPages := total / params.PerPage
	if total%params.PerPage > 0 {
		totalPages++
	}

	return &repository.PaginatedResult[domain.Branch]{
		Data:       branches,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

// Create creates a new branch
func (r *BranchRepository) Create(ctx context.Context, branch *domain.Branch) error {
	query := `
		INSERT INTO branches (name, code, address, phone, email, is_active, timezone, currency,
							  default_interest_rate, default_loan_term_days, default_grace_period)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		branch.Name, branch.Code, NullString(branch.Address), NullString(branch.Phone),
		NullString(branch.Email), branch.IsActive, branch.Timezone, branch.Currency,
		branch.DefaultInterestRate, branch.DefaultLoanTermDays, branch.DefaultGracePeriod,
	).Scan(&branch.ID, &branch.CreatedAt, &branch.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

// Update updates an existing branch
func (r *BranchRepository) Update(ctx context.Context, branch *domain.Branch) error {
	query := `
		UPDATE branches SET
			name = $2, code = $3, address = $4, phone = $5, email = $6,
			is_active = $7, timezone = $8, currency = $9,
			default_interest_rate = $10, default_loan_term_days = $11, default_grace_period = $12,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		branch.ID, branch.Name, branch.Code, NullString(branch.Address),
		NullString(branch.Phone), NullString(branch.Email), branch.IsActive,
		branch.Timezone, branch.Currency, branch.DefaultInterestRate,
		branch.DefaultLoanTermDays, branch.DefaultGracePeriod,
	)
	if err != nil {
		return fmt.Errorf("failed to update branch: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("branch not found")
	}

	return nil
}

// Delete soft deletes a branch
func (r *BranchRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE branches SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("branch not found")
	}

	return nil
}
