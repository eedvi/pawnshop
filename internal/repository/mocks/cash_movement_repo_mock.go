package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// MockCashMovementRepository is a mock implementation of CashMovementRepository
type MockCashMovementRepository struct {
	mock.Mock
}

func (m *MockCashMovementRepository) GetByID(ctx context.Context, id int64) (*domain.CashMovement, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CashMovement), args.Error(1)
}

func (m *MockCashMovementRepository) List(ctx context.Context, params repository.CashMovementListParams) (*repository.PaginatedResult[domain.CashMovement], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[domain.CashMovement]), args.Error(1)
}

func (m *MockCashMovementRepository) ListBySession(ctx context.Context, sessionID int64) ([]*domain.CashMovement, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CashMovement), args.Error(1)
}

func (m *MockCashMovementRepository) Create(ctx context.Context, movement *domain.CashMovement) error {
	args := m.Called(ctx, movement)
	return args.Error(0)
}
