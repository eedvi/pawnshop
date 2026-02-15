package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// MockSaleRepository is a mock implementation of SaleRepository
type MockSaleRepository struct {
	mock.Mock
}

func (m *MockSaleRepository) GetByID(ctx context.Context, id int64) (*domain.Sale, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Sale), args.Error(1)
}

func (m *MockSaleRepository) GetByNumber(ctx context.Context, saleNumber string) (*domain.Sale, error) {
	args := m.Called(ctx, saleNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Sale), args.Error(1)
}

func (m *MockSaleRepository) List(ctx context.Context, params repository.SaleListParams) (*repository.PaginatedResult[domain.Sale], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[domain.Sale]), args.Error(1)
}

func (m *MockSaleRepository) Create(ctx context.Context, sale *domain.Sale) error {
	args := m.Called(ctx, sale)
	return args.Error(0)
}

func (m *MockSaleRepository) Update(ctx context.Context, sale *domain.Sale) error {
	args := m.Called(ctx, sale)
	return args.Error(0)
}

func (m *MockSaleRepository) GenerateNumber(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
