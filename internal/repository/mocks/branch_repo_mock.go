package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// MockBranchRepository is a mock implementation of BranchRepository
type MockBranchRepository struct {
	mock.Mock
}

func (m *MockBranchRepository) GetByID(ctx context.Context, id int64) (*domain.Branch, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Branch), args.Error(1)
}

func (m *MockBranchRepository) GetByCode(ctx context.Context, code string) (*domain.Branch, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Branch), args.Error(1)
}

func (m *MockBranchRepository) List(ctx context.Context, params repository.PaginationParams) (*repository.PaginatedResult[domain.Branch], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[domain.Branch]), args.Error(1)
}

func (m *MockBranchRepository) Create(ctx context.Context, branch *domain.Branch) error {
	args := m.Called(ctx, branch)
	return args.Error(0)
}

func (m *MockBranchRepository) Update(ctx context.Context, branch *domain.Branch) error {
	args := m.Called(ctx, branch)
	return args.Error(0)
}

func (m *MockBranchRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
