package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
)

// MockCashRegisterRepository is a mock implementation of CashRegisterRepository
type MockCashRegisterRepository struct {
	mock.Mock
}

func (m *MockCashRegisterRepository) GetByID(ctx context.Context, id int64) (*domain.CashRegister, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CashRegister), args.Error(1)
}

func (m *MockCashRegisterRepository) List(ctx context.Context, branchID int64) ([]*domain.CashRegister, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CashRegister), args.Error(1)
}

func (m *MockCashRegisterRepository) Create(ctx context.Context, register *domain.CashRegister) error {
	args := m.Called(ctx, register)
	return args.Error(0)
}

func (m *MockCashRegisterRepository) Update(ctx context.Context, register *domain.CashRegister) error {
	args := m.Called(ctx, register)
	return args.Error(0)
}
