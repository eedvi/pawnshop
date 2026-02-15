package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// MockTransferRepository is a mock implementation of TransferRepository
type MockTransferRepository struct {
	mock.Mock
}

func (m *MockTransferRepository) Create(ctx context.Context, transfer *domain.ItemTransfer) error {
	args := m.Called(ctx, transfer)
	return args.Error(0)
}

func (m *MockTransferRepository) GetByID(ctx context.Context, id int64) (*domain.ItemTransfer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ItemTransfer), args.Error(1)
}

func (m *MockTransferRepository) GetByNumber(ctx context.Context, number string) (*domain.ItemTransfer, error) {
	args := m.Called(ctx, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ItemTransfer), args.Error(1)
}

func (m *MockTransferRepository) Update(ctx context.Context, transfer *domain.ItemTransfer) error {
	args := m.Called(ctx, transfer)
	return args.Error(0)
}

func (m *MockTransferRepository) List(ctx context.Context, filter repository.TransferFilter) ([]*domain.ItemTransfer, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.ItemTransfer), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransferRepository) ListByBranch(ctx context.Context, branchID int64, filter repository.TransferFilter) ([]*domain.ItemTransfer, int64, error) {
	args := m.Called(ctx, branchID, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.ItemTransfer), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransferRepository) ListByItem(ctx context.Context, itemID int64) ([]*domain.ItemTransfer, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ItemTransfer), args.Error(1)
}

func (m *MockTransferRepository) GetPendingForBranch(ctx context.Context, branchID int64) ([]*domain.ItemTransfer, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ItemTransfer), args.Error(1)
}

func (m *MockTransferRepository) GetInTransitForBranch(ctx context.Context, branchID int64) ([]*domain.ItemTransfer, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ItemTransfer), args.Error(1)
}

func (m *MockTransferRepository) GenerateTransferNumber(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
