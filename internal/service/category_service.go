package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// CategoryService handles category business logic
type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new CategoryService
func NewCategoryService(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// CreateCategoryInput represents create category request data
type CreateCategoryInput struct {
	Name                string   `json:"name" validate:"required,min=2"`
	ParentID            *int64   `json:"parent_id"`
	Description         *string  `json:"description"`
	Icon                *string  `json:"icon"`
	DefaultInterestRate float64  `json:"default_interest_rate" validate:"gte=0,lte=100"`
	MinLoanAmount       *float64 `json:"min_loan_amount" validate:"omitempty,gte=0"`
	MaxLoanAmount       *float64 `json:"max_loan_amount" validate:"omitempty,gte=0"`
	LoanToValueRatio    float64  `json:"loan_to_value_ratio" validate:"gte=0,lte=1"`
	SortOrder           int      `json:"sort_order"`
}

// Create creates a new category
func (s *CategoryService) Create(ctx context.Context, input CreateCategoryInput) (*domain.Category, error) {
	// Generate slug from name
	slug := generateSlug(input.Name)

	// Check for duplicate slug
	existing, _ := s.categoryRepo.GetBySlug(ctx, slug)
	if existing != nil {
		return nil, errors.New("category with this name already exists")
	}

	// Validate parent if provided
	if input.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(ctx, *input.ParentID)
		if err != nil {
			return nil, errors.New("parent category not found")
		}
		// Prevent deep nesting (max 2 levels)
		if parent.ParentID != nil {
			return nil, errors.New("categories can only be nested 2 levels deep")
		}
	}

	// Validate min/max loan amounts
	if input.MinLoanAmount != nil && input.MaxLoanAmount != nil {
		if *input.MinLoanAmount > *input.MaxLoanAmount {
			return nil, errors.New("min loan amount cannot exceed max loan amount")
		}
	}

	category := &domain.Category{
		Name:                input.Name,
		Slug:                slug,
		ParentID:            input.ParentID,
		Description:         input.Description,
		Icon:                input.Icon,
		DefaultInterestRate: input.DefaultInterestRate,
		MinLoanAmount:       input.MinLoanAmount,
		MaxLoanAmount:       input.MaxLoanAmount,
		LoanToValueRatio:    input.LoanToValueRatio,
		SortOrder:           input.SortOrder,
		IsActive:            true,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// UpdateCategoryInput represents update category request data
type UpdateCategoryInput struct {
	Name                string   `json:"name" validate:"omitempty,min=2"`
	Description         *string  `json:"description"`
	Icon                *string  `json:"icon"`
	DefaultInterestRate *float64 `json:"default_interest_rate" validate:"omitempty,gte=0,lte=100"`
	MinLoanAmount       *float64 `json:"min_loan_amount" validate:"omitempty,gte=0"`
	MaxLoanAmount       *float64 `json:"max_loan_amount" validate:"omitempty,gte=0"`
	LoanToValueRatio    *float64 `json:"loan_to_value_ratio" validate:"omitempty,gte=0,lte=1"`
	SortOrder           *int     `json:"sort_order"`
	IsActive            *bool    `json:"is_active"`
}

// Update updates an existing category
func (s *CategoryService) Update(ctx context.Context, id int64, input UpdateCategoryInput) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Update fields
	if input.Name != "" {
		category.Name = input.Name
		category.Slug = generateSlug(input.Name)
	}
	if input.Description != nil {
		category.Description = input.Description
	}
	if input.Icon != nil {
		category.Icon = input.Icon
	}
	if input.DefaultInterestRate != nil {
		category.DefaultInterestRate = *input.DefaultInterestRate
	}
	if input.MinLoanAmount != nil {
		category.MinLoanAmount = input.MinLoanAmount
	}
	if input.MaxLoanAmount != nil {
		category.MaxLoanAmount = input.MaxLoanAmount
	}
	if input.LoanToValueRatio != nil {
		category.LoanToValueRatio = *input.LoanToValueRatio
	}
	if input.SortOrder != nil {
		category.SortOrder = *input.SortOrder
	}
	if input.IsActive != nil {
		category.IsActive = *input.IsActive
	}

	// Validate min/max loan amounts
	if category.MinLoanAmount != nil && category.MaxLoanAmount != nil {
		if *category.MinLoanAmount > *category.MaxLoanAmount {
			return nil, errors.New("min loan amount cannot exceed max loan amount")
		}
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// GetByID retrieves a category by ID
func (s *CategoryService) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

// GetBySlug retrieves a category by slug
func (s *CategoryService) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	category, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

// List retrieves categories with filters
func (s *CategoryService) List(ctx context.Context, params repository.CategoryListParams) ([]*domain.Category, error) {
	return s.categoryRepo.List(ctx, params)
}

// ListWithChildren retrieves all categories as a tree structure
func (s *CategoryService) ListWithChildren(ctx context.Context) ([]*domain.Category, error) {
	return s.categoryRepo.ListWithChildren(ctx)
}

// Delete deletes a category
func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	_, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("category not found")
	}

	return s.categoryRepo.Delete(ctx, id)
}

// Helper function to generate slug from name
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters
	reg := regexp.MustCompile("[^a-z0-9-]")
	slug = reg.ReplaceAllString(slug, "")
	// Remove multiple hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")
	return slug
}
