package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// MockItemRepository is a mock implementation of ItemRepository
type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) GetByID(ctx context.Context, id int64) (*domain.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockItemRepository) GetBySKU(ctx context.Context, sku string) (*domain.Item, error) {
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockItemRepository) List(ctx context.Context, params repository.ItemListParams) (*repository.PaginatedResult[domain.Item], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[domain.Item]), args.Error(1)
}

func (m *MockItemRepository) Create(ctx context.Context, item *domain.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemRepository) Update(ctx context.Context, item *domain.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemRepository) UpdateStatus(ctx context.Context, id int64, status domain.ItemStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockItemRepository) GenerateSKU(ctx context.Context, branchID int64) (string, error) {
	args := m.Called(ctx, branchID)
	return args.String(0), args.Error(1)
}

func (m *MockItemRepository) CreateHistory(ctx context.Context, history *domain.ItemHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}
