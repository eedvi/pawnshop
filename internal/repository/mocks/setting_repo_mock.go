package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
)

// MockSettingRepository is a mock implementation of SettingRepository
type MockSettingRepository struct {
	mock.Mock
}

func (m *MockSettingRepository) Get(ctx context.Context, key string, branchID *int64) (*domain.Setting, error) {
	args := m.Called(ctx, key, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Setting), args.Error(1)
}

func (m *MockSettingRepository) GetAll(ctx context.Context, branchID *int64) ([]*domain.Setting, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Setting), args.Error(1)
}

func (m *MockSettingRepository) Set(ctx context.Context, setting *domain.Setting) error {
	args := m.Called(ctx, setting)
	return args.Error(0)
}

func (m *MockSettingRepository) Delete(ctx context.Context, key string, branchID *int64) error {
	args := m.Called(ctx, key, branchID)
	return args.Error(0)
}
