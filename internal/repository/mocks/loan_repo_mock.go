package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// MockLoanRepository is a mock implementation of LoanRepository
type MockLoanRepository struct {
	mock.Mock
}

func (m *MockLoanRepository) GetByID(ctx context.Context, id int64) (*domain.Loan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Loan), args.Error(1)
}

func (m *MockLoanRepository) GetByNumber(ctx context.Context, loanNumber string) (*domain.Loan, error) {
	args := m.Called(ctx, loanNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Loan), args.Error(1)
}

func (m *MockLoanRepository) List(ctx context.Context, params repository.LoanListParams) (*repository.PaginatedResult[domain.Loan], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[domain.Loan]), args.Error(1)
}

func (m *MockLoanRepository) Create(ctx context.Context, loan *domain.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepository) Update(ctx context.Context, loan *domain.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepository) GenerateNumber(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockLoanRepository) GetOverdueLoans(ctx context.Context, branchID int64) ([]*domain.Loan, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Loan), args.Error(1)
}

func (m *MockLoanRepository) UpdateStatus(ctx context.Context, id int64, status domain.LoanStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockLoanRepository) BeginTx(ctx context.Context) (repository.Transaction, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(repository.Transaction), args.Error(1)
}

func (m *MockLoanRepository) CreateTx(ctx context.Context, tx repository.Transaction, loan *domain.Loan) error {
	args := m.Called(ctx, tx, loan)
	return args.Error(0)
}

func (m *MockLoanRepository) CreateInstallments(ctx context.Context, installments []*domain.LoanInstallment) error {
	args := m.Called(ctx, installments)
	return args.Error(0)
}

func (m *MockLoanRepository) CreateInstallmentsTx(ctx context.Context, tx repository.Transaction, installments []*domain.LoanInstallment) error {
	args := m.Called(ctx, tx, installments)
	return args.Error(0)
}

func (m *MockLoanRepository) GetInstallments(ctx context.Context, loanID int64) ([]*domain.LoanInstallment, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.LoanInstallment), args.Error(1)
}

func (m *MockLoanRepository) UpdateInstallment(ctx context.Context, installment *domain.LoanInstallment) error {
	args := m.Called(ctx, installment)
	return args.Error(0)
}
