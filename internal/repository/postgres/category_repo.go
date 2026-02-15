package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// CategoryRepository implements repository.CategoryRepository
type CategoryRepository struct {
	db *DB
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	query := `
		SELECT id, parent_id, name, slug, description, icon,
			   default_interest_rate, min_loan_amount, max_loan_amount,
			   loan_to_value_ratio, sort_order, is_active,
			   created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	return r.scanCategory(r.db.QueryRowContext(ctx, query, id))
}

// GetBySlug retrieves a category by slug
func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	query := `
		SELECT id, parent_id, name, slug, description, icon,
			   default_interest_rate, min_loan_amount, max_loan_amount,
			   loan_to_value_ratio, sort_order, is_active,
			   created_at, updated_at
		FROM categories
		WHERE slug = $1
	`

	return r.scanCategory(r.db.QueryRowContext(ctx, query, slug))
}

// List retrieves categories with filters
func (r *CategoryRepository) List(ctx context.Context, params repository.CategoryListParams) ([]*domain.Category, error) {
	query := `
		SELECT id, parent_id, name, slug, description, icon,
			   default_interest_rate, min_loan_amount, max_loan_amount,
			   loan_to_value_ratio, sort_order, is_active,
			   created_at, updated_at
		FROM categories
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if params.ParentID != nil {
		argCount++
		query += fmt.Sprintf(" AND parent_id = $%d", argCount)
		args = append(args, *params.ParentID)
	}

	if params.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *params.IsActive)
	}

	query += " ORDER BY sort_order ASC, name ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	categories := []*domain.Category{}
	for rows.Next() {
		category, err := r.scanCategoryRow(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// ListWithChildren retrieves all categories with their children
func (r *CategoryRepository) ListWithChildren(ctx context.Context) ([]*domain.Category, error) {
	query := `
		SELECT id, parent_id, name, slug, description, icon,
			   default_interest_rate, min_loan_amount, max_loan_amount,
			   loan_to_value_ratio, sort_order, is_active,
			   created_at, updated_at
		FROM categories
		WHERE is_active = true
		ORDER BY parent_id NULLS FIRST, sort_order ASC, name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	categoriesMap := make(map[int64]*domain.Category)
	var rootCategories []*domain.Category

	for rows.Next() {
		category, err := r.scanCategoryRow(rows)
		if err != nil {
			return nil, err
		}
		categoriesMap[category.ID] = category

		if category.ParentID == nil {
			rootCategories = append(rootCategories, category)
		}
	}

	// Build tree structure
	for _, category := range categoriesMap {
		if category.ParentID != nil {
			if parent, ok := categoriesMap[*category.ParentID]; ok {
				parent.Children = append(parent.Children, category)
			}
		}
	}

	return rootCategories, nil
}

// Create creates a new category
func (r *CategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	query := `
		INSERT INTO categories (
			parent_id, name, slug, description, icon,
			default_interest_rate, min_loan_amount, max_loan_amount,
			loan_to_value_ratio, sort_order, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		NullInt64(category.ParentID), category.Name, category.Slug, NullStringPtr(category.Description),
		NullStringPtr(category.Icon), category.DefaultInterestRate,
		NullFloat64(category.MinLoanAmount), NullFloat64(category.MaxLoanAmount),
		category.LoanToValueRatio, category.SortOrder, category.IsActive,
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

// Update updates an existing category
func (r *CategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	query := `
		UPDATE categories SET
			parent_id = $2, name = $3, slug = $4, description = $5, icon = $6,
			default_interest_rate = $7, min_loan_amount = $8, max_loan_amount = $9,
			loan_to_value_ratio = $10, sort_order = $11, is_active = $12,
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		category.ID, NullInt64(category.ParentID), category.Name, category.Slug,
		NullStringPtr(category.Description), NullStringPtr(category.Icon),
		category.DefaultInterestRate, NullFloat64(category.MinLoanAmount),
		NullFloat64(category.MaxLoanAmount), category.LoanToValueRatio,
		category.SortOrder, category.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

// Delete deletes a category
func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	// Check if category has children
	var childCount int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM categories WHERE parent_id = $1", id).Scan(&childCount)
	if err != nil {
		return fmt.Errorf("failed to check for children: %w", err)
	}
	if childCount > 0 {
		return fmt.Errorf("cannot delete category with children")
	}

	// Check if category has items
	var itemCount int
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM items WHERE category_id = $1 AND deleted_at IS NULL", id).Scan(&itemCount)
	if err != nil {
		return fmt.Errorf("failed to check for items: %w", err)
	}
	if itemCount > 0 {
		return fmt.Errorf("cannot delete category with items")
	}

	result, err := r.db.ExecContext(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

// Helper functions
func (r *CategoryRepository) scanCategory(row *sql.Row) (*domain.Category, error) {
	category := &domain.Category{}
	var parentID sql.NullInt64
	var description, icon sql.NullString
	var minLoanAmount, maxLoanAmount sql.NullFloat64

	err := row.Scan(
		&category.ID, &parentID, &category.Name, &category.Slug,
		&description, &icon, &category.DefaultInterestRate,
		&minLoanAmount, &maxLoanAmount, &category.LoanToValueRatio,
		&category.SortOrder, &category.IsActive,
		&category.CreatedAt, &category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	category.ParentID = Int64Ptr(parentID)
	category.Description = StringPtrVal(description)
	category.Icon = StringPtrVal(icon)
	category.MinLoanAmount = Float64Ptr(minLoanAmount)
	category.MaxLoanAmount = Float64Ptr(maxLoanAmount)

	return category, nil
}

func (r *CategoryRepository) scanCategoryRow(rows *sql.Rows) (*domain.Category, error) {
	category := &domain.Category{}
	var parentID sql.NullInt64
	var description, icon sql.NullString
	var minLoanAmount, maxLoanAmount sql.NullFloat64

	err := rows.Scan(
		&category.ID, &parentID, &category.Name, &category.Slug,
		&description, &icon, &category.DefaultInterestRate,
		&minLoanAmount, &maxLoanAmount, &category.LoanToValueRatio,
		&category.SortOrder, &category.IsActive,
		&category.CreatedAt, &category.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan category: %w", err)
	}

	category.ParentID = Int64Ptr(parentID)
	category.Description = StringPtrVal(description)
	category.Icon = StringPtrVal(icon)
	category.MinLoanAmount = Float64Ptr(minLoanAmount)
	category.MaxLoanAmount = Float64Ptr(maxLoanAmount)

	return category, nil
}
