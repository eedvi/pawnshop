package service

import (
	"context"
	"errors"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

var (
	ErrExpenseAlreadyApproved = errors.New("expense is already approved")
)

// ExpenseService defines the interface for expense operations
type ExpenseService interface {
	// CreateCategory creates a new expense category
	CreateCategory(ctx context.Context, req CreateExpenseCategoryRequest) (*domain.ExpenseCategory, error)

	// GetCategoryByID retrieves a category by ID
	GetCategoryByID(ctx context.Context, id int64) (*domain.ExpenseCategory, error)

	// UpdateCategory updates an expense category
	UpdateCategory(ctx context.Context, id int64, req UpdateExpenseCategoryRequest) (*domain.ExpenseCategory, error)

	// ListCategories retrieves all expense categories
	ListCategories(ctx context.Context, includeInactive bool) ([]*domain.ExpenseCategory, error)

	// Create creates a new expense
	Create(ctx context.Context, req CreateExpenseRequest) (*domain.Expense, error)

	// GetByID retrieves an expense by ID
	GetByID(ctx context.Context, id int64) (*domain.Expense, error)

	// Update updates an expense
	Update(ctx context.Context, id int64, req UpdateExpenseRequest) (*domain.Expense, error)

	// Delete deletes an expense
	Delete(ctx context.Context, id int64) error

	// List retrieves expenses with filtering
	List(ctx context.Context, filter repository.ExpenseFilter) ([]*domain.Expense, int64, error)

	// ListByBranch retrieves expenses for a branch
	ListByBranch(ctx context.Context, branchID int64, filter repository.ExpenseFilter) ([]*domain.Expense, int64, error)

	// Approve approves an expense
	Approve(ctx context.Context, id int64, approvedBy int64) (*domain.Expense, error)

	// GetTotalByBranchAndDate retrieves total expenses for a branch on a date
	GetTotalByBranchAndDate(ctx context.Context, branchID int64, date time.Time) (float64, error)
}

type expenseService struct {
	expenseRepo  repository.ExpenseRepository
	categoryRepo repository.ExpenseCategoryRepository
	branchRepo   repository.BranchRepository
}

// NewExpenseService creates a new expense service
func NewExpenseService(
	expenseRepo repository.ExpenseRepository,
	categoryRepo repository.ExpenseCategoryRepository,
	branchRepo repository.BranchRepository,
) ExpenseService {
	return &expenseService{
		expenseRepo:  expenseRepo,
		categoryRepo: categoryRepo,
		branchRepo:   branchRepo,
	}
}

// CreateExpenseCategoryRequest represents a request to create an expense category
type CreateExpenseCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description"`
}

// UpdateExpenseCategoryRequest represents a request to update an expense category
type UpdateExpenseCategoryRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// CreateExpenseRequest represents a request to create an expense
type CreateExpenseRequest struct {
	BranchID      int64     `json:"branch_id" validate:"required"`
	CategoryID    *int64    `json:"category_id"`
	Description   string    `json:"description" validate:"required"`
	Amount        float64   `json:"amount" validate:"required,gt=0"`
	ExpenseDate   time.Time `json:"expense_date" validate:"required"`
	PaymentMethod string    `json:"payment_method" validate:"required"`
	ReceiptNumber string    `json:"receipt_number"`
	Vendor        string    `json:"vendor"`
	CreatedBy     int64     `json:"created_by" validate:"required"`
}

// UpdateExpenseRequest represents a request to update an expense
type UpdateExpenseRequest struct {
	CategoryID    *int64    `json:"category_id"`
	Description   string    `json:"description"`
	Amount        float64   `json:"amount" validate:"gt=0"`
	ExpenseDate   time.Time `json:"expense_date"`
	PaymentMethod string    `json:"payment_method"`
	ReceiptNumber string    `json:"receipt_number"`
	Vendor        string    `json:"vendor"`
}

func (s *expenseService) CreateCategory(ctx context.Context, req CreateExpenseCategoryRequest) (*domain.ExpenseCategory, error) {
	// Check if code already exists
	existing, err := s.categoryRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("category code already exists")
	}

	category := &domain.ExpenseCategory{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *expenseService) GetCategoryByID(ctx context.Context, id int64) (*domain.ExpenseCategory, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *expenseService) UpdateCategory(ctx context.Context, id int64, req UpdateExpenseCategoryRequest) (*domain.ExpenseCategory, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Code != "" {
		category.Code = req.Code
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *expenseService) ListCategories(ctx context.Context, includeInactive bool) ([]*domain.ExpenseCategory, error) {
	return s.categoryRepo.List(ctx, includeInactive)
}

func (s *expenseService) Create(ctx context.Context, req CreateExpenseRequest) (*domain.Expense, error) {
	// Validate branch exists
	branch, err := s.branchRepo.GetByID(ctx, req.BranchID)
	if err != nil || branch == nil {
		return nil, ErrBranchNotFound
	}

	// Validate category if provided
	if req.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil || category == nil {
			return nil, ErrCategoryNotFound
		}
	}

	// Generate expense number
	expenseNumber, err := s.expenseRepo.GenerateExpenseNumber(ctx)
	if err != nil {
		return nil, err
	}

	expense := &domain.Expense{
		ExpenseNumber: expenseNumber,
		BranchID:      req.BranchID,
		CategoryID:    req.CategoryID,
		Description:   req.Description,
		Amount:        req.Amount,
		ExpenseDate:   req.ExpenseDate,
		PaymentMethod: req.PaymentMethod,
		ReceiptNumber: req.ReceiptNumber,
		Vendor:        req.Vendor,
		CreatedBy:     &req.CreatedBy,
	}

	if err := s.expenseRepo.Create(ctx, expense); err != nil {
		return nil, err
	}

	return expense, nil
}

func (s *expenseService) GetByID(ctx context.Context, id int64) (*domain.Expense, error) {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if expense == nil {
		return nil, ErrExpenseNotFound
	}

	// Load related entities
	if expense.CategoryID != nil {
		category, _ := s.categoryRepo.GetByID(ctx, *expense.CategoryID)
		expense.Category = category
	}

	branch, _ := s.branchRepo.GetByID(ctx, expense.BranchID)
	expense.Branch = branch

	return expense, nil
}

func (s *expenseService) Update(ctx context.Context, id int64, req UpdateExpenseRequest) (*domain.Expense, error) {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if expense == nil {
		return nil, ErrExpenseNotFound
	}

	// Cannot update approved expenses
	if expense.IsApproved() {
		return nil, ErrExpenseAlreadyApproved
	}

	// Validate category if changing
	if req.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil || category == nil {
			return nil, ErrCategoryNotFound
		}
		expense.CategoryID = req.CategoryID
	}

	if req.Description != "" {
		expense.Description = req.Description
	}
	if req.Amount > 0 {
		expense.Amount = req.Amount
	}
	if !req.ExpenseDate.IsZero() {
		expense.ExpenseDate = req.ExpenseDate
	}
	if req.PaymentMethod != "" {
		expense.PaymentMethod = req.PaymentMethod
	}
	if req.ReceiptNumber != "" {
		expense.ReceiptNumber = req.ReceiptNumber
	}
	if req.Vendor != "" {
		expense.Vendor = req.Vendor
	}

	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return nil, err
	}

	return expense, nil
}

func (s *expenseService) Delete(ctx context.Context, id int64) error {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if expense == nil {
		return ErrExpenseNotFound
	}

	// Cannot delete approved expenses
	if expense.IsApproved() {
		return ErrExpenseAlreadyApproved
	}

	return s.expenseRepo.Delete(ctx, id)
}

func (s *expenseService) List(ctx context.Context, filter repository.ExpenseFilter) ([]*domain.Expense, int64, error) {
	return s.expenseRepo.List(ctx, filter)
}

func (s *expenseService) ListByBranch(ctx context.Context, branchID int64, filter repository.ExpenseFilter) ([]*domain.Expense, int64, error) {
	return s.expenseRepo.ListByBranch(ctx, branchID, filter)
}

func (s *expenseService) Approve(ctx context.Context, id int64, approvedBy int64) (*domain.Expense, error) {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if expense == nil {
		return nil, ErrExpenseNotFound
	}

	if expense.IsApproved() {
		return nil, ErrExpenseAlreadyApproved
	}

	if err := s.expenseRepo.Approve(ctx, id, approvedBy); err != nil {
		return nil, err
	}

	// Reload expense
	return s.GetByID(ctx, id)
}

func (s *expenseService) GetTotalByBranchAndDate(ctx context.Context, branchID int64, date time.Time) (float64, error) {
	return s.expenseRepo.GetTotalByBranchAndDate(ctx, branchID, date)
}
