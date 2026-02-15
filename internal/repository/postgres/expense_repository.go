package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

type expenseCategoryRepository struct {
	db *DB
}

// NewExpenseCategoryRepository creates a new expense category repository
func NewExpenseCategoryRepository(db *DB) repository.ExpenseCategoryRepository {
	return &expenseCategoryRepository{db: db}
}

func (r *expenseCategoryRepository) Create(ctx context.Context, category *domain.ExpenseCategory) error {
	query := `
		INSERT INTO expense_categories (name, code, description, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		category.Name,
		category.Code,
		category.Description,
		category.IsActive,
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
}

func (r *expenseCategoryRepository) GetByID(ctx context.Context, id int64) (*domain.ExpenseCategory, error) {
	query := `
		SELECT id, name, code, description, is_active, created_at, updated_at
		FROM expense_categories
		WHERE id = $1`

	category := &domain.ExpenseCategory{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Code,
		&category.Description,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *expenseCategoryRepository) GetByCode(ctx context.Context, code string) (*domain.ExpenseCategory, error) {
	query := `
		SELECT id, name, code, description, is_active, created_at, updated_at
		FROM expense_categories
		WHERE code = $1`

	category := &domain.ExpenseCategory{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&category.ID,
		&category.Name,
		&category.Code,
		&category.Description,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *expenseCategoryRepository) Update(ctx context.Context, category *domain.ExpenseCategory) error {
	query := `
		UPDATE expense_categories SET
			name = $2,
			code = $3,
			description = $4,
			is_active = $5,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRowContext(ctx, query,
		category.ID,
		category.Name,
		category.Code,
		category.Description,
		category.IsActive,
	).Scan(&category.UpdatedAt)
}

func (r *expenseCategoryRepository) List(ctx context.Context, includeInactive bool) ([]*domain.ExpenseCategory, error) {
	query := `
		SELECT id, name, code, description, is_active, created_at, updated_at
		FROM expense_categories`

	if !includeInactive {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.ExpenseCategory
	for rows.Next() {
		category := &domain.ExpenseCategory{}
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Code,
			&category.Description,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, rows.Err()
}

// Expense Repository
type expenseRepository struct {
	db *DB
}

// NewExpenseRepository creates a new expense repository
func NewExpenseRepository(db *DB) repository.ExpenseRepository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
	query := `
		INSERT INTO expenses (
			expense_number, branch_id, category_id, description, amount,
			expense_date, payment_method, receipt_number, vendor, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		expense.ExpenseNumber,
		expense.BranchID,
		expense.CategoryID,
		expense.Description,
		expense.Amount,
		expense.ExpenseDate,
		expense.PaymentMethod,
		expense.ReceiptNumber,
		expense.Vendor,
		expense.CreatedBy,
	).Scan(&expense.ID, &expense.CreatedAt, &expense.UpdatedAt)
}

func (r *expenseRepository) GetByID(ctx context.Context, id int64) (*domain.Expense, error) {
	query := `
		SELECT id, expense_number, branch_id, category_id, description, amount,
			   expense_date, payment_method, receipt_number, vendor,
			   approved_by, approved_at, created_by, created_at, updated_at
		FROM expenses
		WHERE id = $1`

	expense := &domain.Expense{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&expense.ID,
		&expense.ExpenseNumber,
		&expense.BranchID,
		&expense.CategoryID,
		&expense.Description,
		&expense.Amount,
		&expense.ExpenseDate,
		&expense.PaymentMethod,
		&expense.ReceiptNumber,
		&expense.Vendor,
		&expense.ApprovedBy,
		&expense.ApprovedAt,
		&expense.CreatedBy,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (r *expenseRepository) GetByNumber(ctx context.Context, number string) (*domain.Expense, error) {
	query := `
		SELECT id, expense_number, branch_id, category_id, description, amount,
			   expense_date, payment_method, receipt_number, vendor,
			   approved_by, approved_at, created_by, created_at, updated_at
		FROM expenses
		WHERE expense_number = $1`

	expense := &domain.Expense{}
	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&expense.ID,
		&expense.ExpenseNumber,
		&expense.BranchID,
		&expense.CategoryID,
		&expense.Description,
		&expense.Amount,
		&expense.ExpenseDate,
		&expense.PaymentMethod,
		&expense.ReceiptNumber,
		&expense.Vendor,
		&expense.ApprovedBy,
		&expense.ApprovedAt,
		&expense.CreatedBy,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (r *expenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	query := `
		UPDATE expenses SET
			category_id = $2,
			description = $3,
			amount = $4,
			expense_date = $5,
			payment_method = $6,
			receipt_number = $7,
			vendor = $8,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRowContext(ctx, query,
		expense.ID,
		expense.CategoryID,
		expense.Description,
		expense.Amount,
		expense.ExpenseDate,
		expense.PaymentMethod,
		expense.ReceiptNumber,
		expense.Vendor,
	).Scan(&expense.UpdatedAt)
}

func (r *expenseRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM expenses WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *expenseRepository) List(ctx context.Context, filter repository.ExpenseFilter) ([]*domain.Expense, int64, error) {
	var conditions []string
	var args []interface{}
	argPos := 1

	if filter.BranchID != nil {
		conditions = append(conditions, fmt.Sprintf("branch_id = $%d", argPos))
		args = append(args, *filter.BranchID)
		argPos++
	}
	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argPos))
		args = append(args, *filter.CategoryID)
		argPos++
	}
	if filter.IsApproved != nil {
		if *filter.IsApproved {
			conditions = append(conditions, "approved_by IS NOT NULL")
		} else {
			conditions = append(conditions, "approved_by IS NULL")
		}
	}
	if filter.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("expense_date >= $%d", argPos))
		args = append(args, *filter.DateFrom)
		argPos++
	}
	if filter.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("expense_date <= $%d", argPos))
		args = append(args, *filter.DateTo)
		argPos++
	}
	if filter.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("amount >= $%d", argPos))
		args = append(args, *filter.MinAmount)
		argPos++
	}
	if filter.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("amount <= $%d", argPos))
		args = append(args, *filter.MaxAmount)
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM expenses %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT id, expense_number, branch_id, category_id, description, amount,
			   expense_date, payment_method, receipt_number, vendor,
			   approved_by, approved_at, created_by, created_at, updated_at
		FROM expenses
		%s
		ORDER BY expense_date DESC
		LIMIT $%d OFFSET $%d`, whereClause, argPos, argPos+1)

	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []*domain.Expense
	for rows.Next() {
		expense := &domain.Expense{}
		if err := rows.Scan(
			&expense.ID,
			&expense.ExpenseNumber,
			&expense.BranchID,
			&expense.CategoryID,
			&expense.Description,
			&expense.Amount,
			&expense.ExpenseDate,
			&expense.PaymentMethod,
			&expense.ReceiptNumber,
			&expense.Vendor,
			&expense.ApprovedBy,
			&expense.ApprovedAt,
			&expense.CreatedBy,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, total, rows.Err()
}

func (r *expenseRepository) ListByBranch(ctx context.Context, branchID int64, filter repository.ExpenseFilter) ([]*domain.Expense, int64, error) {
	filter.BranchID = &branchID
	return r.List(ctx, filter)
}

func (r *expenseRepository) Approve(ctx context.Context, id int64, approvedBy int64) error {
	query := `
		UPDATE expenses SET
			approved_by = $2,
			approved_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND approved_by IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, approvedBy)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("expense already approved or not found")
	}
	return nil
}

func (r *expenseRepository) GenerateExpenseNumber(ctx context.Context) (string, error) {
	now := time.Now()
	prefix := fmt.Sprintf("EXP-%s-", now.Format("20060102"))

	query := `
		SELECT COUNT(*) + 1 FROM expenses
		WHERE expense_number LIKE $1`

	var seq int
	if err := r.db.QueryRowContext(ctx, query, prefix+"%").Scan(&seq); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%04d", prefix, seq), nil
}

func (r *expenseRepository) GetTotalByBranchAndDate(ctx context.Context, branchID int64, date time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE branch_id = $1 AND DATE(expense_date) = DATE($2)`

	var total float64
	if err := r.db.QueryRowContext(ctx, query, branchID, date).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
