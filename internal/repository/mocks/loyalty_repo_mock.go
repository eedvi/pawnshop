package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
)

// MockLoyaltyRepository is a mock implementation of LoyaltyRepository
type MockLoyaltyRepository struct {
	mock.Mock
}

func (m *MockLoyaltyRepository) CreateHistory(ctx context.Context, history *domain.LoyaltyPointsHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockLoyaltyRepository) GetHistoryByCustomer(ctx context.Context, customerID int64, page, pageSize int) ([]*domain.LoyaltyPointsHistory, int64, error) {
	args := m.Called(ctx, customerID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.LoyaltyPointsHistory), args.Get(1).(int64), args.Error(2)
}

func (m *MockLoyaltyRepository) GetHistoryByReference(ctx context.Context, referenceType string, referenceID int64) ([]*domain.LoyaltyPointsHistory, error) {
	args := m.Called(ctx, referenceType, referenceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.LoyaltyPointsHistory), args.Error(1)
}

func (m *MockLoyaltyRepository) GetTotalPointsEarned(ctx context.Context, customerID int64) (int, error) {
	args := m.Called(ctx, customerID)
	return args.Int(0), args.Error(1)
}

func (m *MockLoyaltyRepository) GetTotalPointsRedeemed(ctx context.Context, customerID int64) (int, error) {
	args := m.Called(ctx, customerID)
	return args.Int(0), args.Error(1)
}
