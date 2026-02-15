package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// MockExpenseRepository is a mock implementation of ExpenseRepository
type MockExpenseRepository struct {
	mock.Mock
}

func (m *MockExpenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) GetByID(ctx context.Context, id int64) (*domain.Expense, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Expense), args.Error(1)
}

func (m *MockExpenseRepository) GetByNumber(ctx context.Context, number string) (*domain.Expense, error) {
	args := m.Called(ctx, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Expense), args.Error(1)
}

func (m *MockExpenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockExpenseRepository) List(ctx context.Context, filter repository.ExpenseFilter) ([]*domain.Expense, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Expense), args.Get(1).(int64), args.Error(2)
}

func (m *MockExpenseRepository) ListByBranch(ctx context.Context, branchID int64, filter repository.ExpenseFilter) ([]*domain.Expense, int64, error) {
	args := m.Called(ctx, branchID, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Expense), args.Get(1).(int64), args.Error(2)
}

func (m *MockExpenseRepository) Approve(ctx context.Context, id int64, approvedBy int64) error {
	args := m.Called(ctx, id, approvedBy)
	return args.Error(0)
}

func (m *MockExpenseRepository) GenerateExpenseNumber(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockExpenseRepository) GetTotalByBranchAndDate(ctx context.Context, branchID int64, date time.Time) (float64, error) {
	args := m.Called(ctx, branchID, date)
	return args.Get(0).(float64), args.Error(1)
}

// MockExpenseCategoryRepository is a mock implementation of ExpenseCategoryRepository
type MockExpenseCategoryRepository struct {
	mock.Mock
}

func (m *MockExpenseCategoryRepository) Create(ctx context.Context, category *domain.ExpenseCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockExpenseCategoryRepository) GetByID(ctx context.Context, id int64) (*domain.ExpenseCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ExpenseCategory), args.Error(1)
}

func (m *MockExpenseCategoryRepository) GetByCode(ctx context.Context, code string) (*domain.ExpenseCategory, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ExpenseCategory), args.Error(1)
}

func (m *MockExpenseCategoryRepository) Update(ctx context.Context, category *domain.ExpenseCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockExpenseCategoryRepository) List(ctx context.Context, includeInactive bool) ([]*domain.ExpenseCategory, error) {
	args := m.Called(ctx, includeInactive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ExpenseCategory), args.Error(1)
}
