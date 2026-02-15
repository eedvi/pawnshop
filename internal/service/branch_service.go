package service

import (
	"context"
	"errors"
	"fmt"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// BranchService handles branch business logic
type BranchService struct {
	branchRepo repository.BranchRepository
}

// NewBranchService creates a new BranchService
func NewBranchService(branchRepo repository.BranchRepository) *BranchService {
	return &BranchService{branchRepo: branchRepo}
}

// CreateBranchInput represents create branch request data
type CreateBranchInput struct {
	Code                string  `json:"code" validate:"required,min=2,max=10"`
	Name                string  `json:"name" validate:"required,min=2"`
	Address             string  `json:"address"`
	Phone               string  `json:"phone"`
	Email               string  `json:"email" validate:"omitempty,email"`
	Timezone            string  `json:"timezone"`
	Currency            string  `json:"currency"`
	DefaultInterestRate float64 `json:"default_interest_rate" validate:"gte=0,lte=100"`
	DefaultLoanTermDays int     `json:"default_loan_term_days" validate:"gte=1"`
	DefaultGracePeriod  int     `json:"default_grace_period" validate:"gte=0"`
}

// Create creates a new branch
func (s *BranchService) Create(ctx context.Context, input CreateBranchInput) (*domain.Branch, error) {
	// Check for duplicate code
	existing, _ := s.branchRepo.GetByCode(ctx, input.Code)
	if existing != nil {
		return nil, errors.New("branch with this code already exists")
	}

	// Set defaults
	if input.Timezone == "" {
		input.Timezone = "America/Mexico_City"
	}
	if input.Currency == "" {
		input.Currency = "MXN"
	}
	if input.DefaultLoanTermDays == 0 {
		input.DefaultLoanTermDays = 30
	}

	branch := &domain.Branch{
		Code:                input.Code,
		Name:                input.Name,
		Address:             input.Address,
		Phone:               input.Phone,
		Email:               input.Email,
		IsActive:            true,
		Timezone:            input.Timezone,
		Currency:            input.Currency,
		DefaultInterestRate: input.DefaultInterestRate,
		DefaultLoanTermDays: input.DefaultLoanTermDays,
		DefaultGracePeriod:  input.DefaultGracePeriod,
	}

	if err := s.branchRepo.Create(ctx, branch); err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	return branch, nil
}

// UpdateBranchInput represents update branch request data
type UpdateBranchInput struct {
	Name                string   `json:"name" validate:"omitempty,min=2"`
	Address             string   `json:"address"`
	Phone               string   `json:"phone"`
	Email               string   `json:"email" validate:"omitempty,email"`
	Timezone            string   `json:"timezone"`
	Currency            string   `json:"currency"`
	IsActive            *bool    `json:"is_active"`
	DefaultInterestRate *float64 `json:"default_interest_rate" validate:"omitempty,gte=0,lte=100"`
	DefaultLoanTermDays *int     `json:"default_loan_term_days" validate:"omitempty,gte=1"`
	DefaultGracePeriod  *int     `json:"default_grace_period" validate:"omitempty,gte=0"`
}

// Update updates an existing branch
func (s *BranchService) Update(ctx context.Context, id int64, input UpdateBranchInput) (*domain.Branch, error) {
	branch, err := s.branchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	// Update fields
	if input.Name != "" {
		branch.Name = input.Name
	}
	if input.Address != "" {
		branch.Address = input.Address
	}
	if input.Phone != "" {
		branch.Phone = input.Phone
	}
	if input.Email != "" {
		branch.Email = input.Email
	}
	if input.Timezone != "" {
		branch.Timezone = input.Timezone
	}
	if input.Currency != "" {
		branch.Currency = input.Currency
	}
	if input.IsActive != nil {
		branch.IsActive = *input.IsActive
	}
	if input.DefaultInterestRate != nil {
		branch.DefaultInterestRate = *input.DefaultInterestRate
	}
	if input.DefaultLoanTermDays != nil {
		branch.DefaultLoanTermDays = *input.DefaultLoanTermDays
	}
	if input.DefaultGracePeriod != nil {
		branch.DefaultGracePeriod = *input.DefaultGracePeriod
	}

	if err := s.branchRepo.Update(ctx, branch); err != nil {
		return nil, fmt.Errorf("failed to update branch: %w", err)
	}

	return branch, nil
}

// GetByID retrieves a branch by ID
func (s *BranchService) GetByID(ctx context.Context, id int64) (*domain.Branch, error) {
	branch, err := s.branchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("branch not found")
	}
	return branch, nil
}

// GetByCode retrieves a branch by code
func (s *BranchService) GetByCode(ctx context.Context, code string) (*domain.Branch, error) {
	branch, err := s.branchRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, errors.New("branch not found")
	}
	return branch, nil
}

// List retrieves branches with pagination
func (s *BranchService) List(ctx context.Context, params repository.PaginationParams) (*repository.PaginatedResult[domain.Branch], error) {
	return s.branchRepo.List(ctx, params)
}

// Delete soft deletes a branch
func (s *BranchService) Delete(ctx context.Context, id int64) error {
	_, err := s.branchRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("branch not found")
	}

	return s.branchRepo.Delete(ctx, id)
}

// Activate activates a branch
func (s *BranchService) Activate(ctx context.Context, id int64) error {
	branch, err := s.branchRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("branch not found")
	}

	branch.IsActive = true
	return s.branchRepo.Update(ctx, branch)
}

// Deactivate deactivates a branch
func (s *BranchService) Deactivate(ctx context.Context, id int64) error {
	branch, err := s.branchRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("branch not found")
	}

	branch.IsActive = false
	return s.branchRepo.Update(ctx, branch)
}
