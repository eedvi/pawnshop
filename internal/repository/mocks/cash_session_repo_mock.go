package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// MockCashSessionRepository is a mock implementation of CashSessionRepository
type MockCashSessionRepository struct {
	mock.Mock
}

func (m *MockCashSessionRepository) GetByID(ctx context.Context, id int64) (*domain.CashSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CashSession), args.Error(1)
}

func (m *MockCashSessionRepository) GetOpenSession(ctx context.Context, userID int64) (*domain.CashSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CashSession), args.Error(1)
}

func (m *MockCashSessionRepository) GetOpenSessionByRegister(ctx context.Context, registerID int64) (*domain.CashSession, error) {
	args := m.Called(ctx, registerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CashSession), args.Error(1)
}

func (m *MockCashSessionRepository) List(ctx context.Context, params repository.CashSessionListParams) (*repository.PaginatedResult[domain.CashSession], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[domain.CashSession]), args.Error(1)
}

func (m *MockCashSessionRepository) Create(ctx context.Context, session *domain.CashSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockCashSessionRepository) Update(ctx context.Context, session *domain.CashSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockCashSessionRepository) Close(ctx context.Context, id int64, closingData repository.CashSessionCloseData) error {
	args := m.Called(ctx, id, closingData)
	return args.Error(0)
}
