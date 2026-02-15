package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/internal/repository/mocks"
)

func setupSaleService() (*SaleService, *mocks.MockSaleRepository, *mocks.MockItemRepository, *mocks.MockCustomerRepository, *mocks.MockBranchRepository) {
	saleRepo := new(mocks.MockSaleRepository)
	itemRepo := new(mocks.MockItemRepository)
	customerRepo := new(mocks.MockCustomerRepository)
	branchRepo := new(mocks.MockBranchRepository)
	service := NewSaleService(saleRepo, itemRepo, customerRepo, branchRepo)
	return service, saleRepo, itemRepo, customerRepo, branchRepo
}

func TestSaleService_Create_Success(t *testing.T) {
	service, saleRepo, itemRepo, _, branchRepo := setupSaleService()
	ctx := context.Background()

	salePrice := 500.0
	item := &domain.Item{
		ID:       1,
		BranchID: 1,
		Status:   domain.ItemStatusForSale,
		SalePrice: &salePrice,
	}
	updatedItem := &domain.Item{
		ID:       1,
		BranchID: 1,
		Status:   domain.ItemStatusSold,
		SalePrice: &salePrice,
	}

	branchRepo.On("GetByID", ctx, int64(1)).Return(&domain.Branch{ID: 1}, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil).Once()
	saleRepo.On("GenerateNumber", ctx).Return("SALE-001", nil)
	saleRepo.On("Create", ctx, mock.AnythingOfType("*domain.Sale")).Return(nil)
	itemRepo.On("UpdateStatus", ctx, int64(1), domain.ItemStatusSold).Return(nil)
	itemRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.ItemHistory")).Return(nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(updatedItem, nil).Once()

	input := CreateSaleInput{
		BranchID:      1,
		ItemID:        1,
		SaleType:      "direct",
		PaymentMethod: "cash",
		CreatedBy:     10,
	}
	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Sale)
	assert.Equal(t, 500.0, result.Sale.SalePrice)
	assert.Equal(t, 500.0, result.Sale.FinalPrice)
	saleRepo.AssertExpectations(t)
	itemRepo.AssertExpectations(t)
}

func TestSaleService_Create_InvalidBranch(t *testing.T) {
	service, _, _, _, branchRepo := setupSaleService()
	ctx := context.Background()

	branchRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreateSaleInput{BranchID: 999, ItemID: 1, SaleType: "direct", PaymentMethod: "cash"}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid branch", err.Error())
	branchRepo.AssertExpectations(t)
}

func TestSaleService_Create_ItemNotFound(t *testing.T) {
	service, _, itemRepo, _, branchRepo := setupSaleService()
	ctx := context.Background()

	branchRepo.On("GetByID", ctx, int64(1)).Return(&domain.Branch{ID: 1}, nil)
	itemRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreateSaleInput{BranchID: 1, ItemID: 999, SaleType: "direct", PaymentMethod: "cash"}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "item not found", err.Error())
}

func TestSaleService_Create_ItemNotForSale(t *testing.T) {
	service, _, itemRepo, _, branchRepo := setupSaleService()
	ctx := context.Background()

	item := &domain.Item{ID: 1, BranchID: 1, Status: domain.ItemStatusPawned}

	branchRepo.On("GetByID", ctx, int64(1)).Return(&domain.Branch{ID: 1}, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	input := CreateSaleInput{BranchID: 1, ItemID: 1, SaleType: "direct", PaymentMethod: "cash"}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "item is not available for sale", err.Error())
}

func TestSaleService_Create_ItemWrongBranch(t *testing.T) {
	service, _, itemRepo, _, branchRepo := setupSaleService()
	ctx := context.Background()

	item := &domain.Item{ID: 1, BranchID: 2, Status: domain.ItemStatusForSale}

	branchRepo.On("GetByID", ctx, int64(1)).Return(&domain.Branch{ID: 1}, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	input := CreateSaleInput{BranchID: 1, ItemID: 1, SaleType: "direct", PaymentMethod: "cash"}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "item does not belong to this branch", err.Error())
}

func TestSaleService_Create_BlockedCustomer(t *testing.T) {
	service, _, itemRepo, customerRepo, branchRepo := setupSaleService()
	ctx := context.Background()

	salePrice := 500.0
	item := &domain.Item{ID: 1, BranchID: 1, Status: domain.ItemStatusForSale, SalePrice: &salePrice}
	customer := &domain.Customer{ID: 1, IsBlocked: true}

	branchRepo.On("GetByID", ctx, int64(1)).Return(&domain.Branch{ID: 1}, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	customerID := int64(1)
	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)

	input := CreateSaleInput{BranchID: 1, ItemID: 1, CustomerID: &customerID, SaleType: "direct", PaymentMethod: "cash"}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "customer is blocked", err.Error())
}

func TestSaleService_Create_NoSalePrice(t *testing.T) {
	service, _, itemRepo, _, branchRepo := setupSaleService()
	ctx := context.Background()

	item := &domain.Item{ID: 1, BranchID: 1, Status: domain.ItemStatusForSale, SalePrice: nil}

	branchRepo.On("GetByID", ctx, int64(1)).Return(&domain.Branch{ID: 1}, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	input := CreateSaleInput{BranchID: 1, ItemID: 1, SaleType: "direct", PaymentMethod: "cash"}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "item does not have a valid sale price", err.Error())
}

func TestSaleService_Create_DiscountExceedsPrice(t *testing.T) {
	service, _, itemRepo, _, branchRepo := setupSaleService()
	ctx := context.Background()

	salePrice := 500.0
	item := &domain.Item{ID: 1, BranchID: 1, Status: domain.ItemStatusForSale, SalePrice: &salePrice}

	branchRepo.On("GetByID", ctx, int64(1)).Return(&domain.Branch{ID: 1}, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	input := CreateSaleInput{BranchID: 1, ItemID: 1, SaleType: "direct", PaymentMethod: "cash", DiscountAmount: 600}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "discount cannot exceed sale price", err.Error())
}

func TestSaleService_Create_WithDiscount(t *testing.T) {
	service, saleRepo, itemRepo, _, branchRepo := setupSaleService()
	ctx := context.Background()

	salePrice := 500.0
	item := &domain.Item{ID: 1, BranchID: 1, Status: domain.ItemStatusForSale, SalePrice: &salePrice}

	branchRepo.On("GetByID", ctx, int64(1)).Return(&domain.Branch{ID: 1}, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	saleRepo.On("GenerateNumber", ctx).Return("SALE-002", nil)
	saleRepo.On("Create", ctx, mock.AnythingOfType("*domain.Sale")).Return(nil)
	itemRepo.On("UpdateStatus", ctx, int64(1), domain.ItemStatusSold).Return(nil)
	itemRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.ItemHistory")).Return(nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	input := CreateSaleInput{BranchID: 1, ItemID: 1, SaleType: "direct", PaymentMethod: "cash", DiscountAmount: 50, CreatedBy: 10}
	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 450.0, result.Sale.FinalPrice)
	assert.Equal(t, 50.0, result.Sale.DiscountAmount)
}

func TestSaleService_GetByID_Success(t *testing.T) {
	service, saleRepo, itemRepo, _, _ := setupSaleService()
	ctx := context.Background()

	sale := &domain.Sale{ID: 1, ItemID: 10, SaleNumber: "SALE-001"}
	item := &domain.Item{ID: 10}

	saleRepo.On("GetByID", ctx, int64(1)).Return(sale, nil)
	itemRepo.On("GetByID", ctx, int64(10)).Return(item, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "SALE-001", result.SaleNumber)
	assert.NotNil(t, result.Item)
	saleRepo.AssertExpectations(t)
}

func TestSaleService_GetByID_WithCustomer(t *testing.T) {
	service, saleRepo, itemRepo, customerRepo, _ := setupSaleService()
	ctx := context.Background()

	customerID := int64(5)
	sale := &domain.Sale{ID: 1, ItemID: 10, CustomerID: &customerID}
	item := &domain.Item{ID: 10}
	customer := &domain.Customer{ID: 5}

	saleRepo.On("GetByID", ctx, int64(1)).Return(sale, nil)
	itemRepo.On("GetByID", ctx, int64(10)).Return(item, nil)
	customerRepo.On("GetByID", ctx, int64(5)).Return(customer, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result.Customer)
}

func TestSaleService_GetByID_NotFound(t *testing.T) {
	service, saleRepo, _, _, _ := setupSaleService()
	ctx := context.Background()

	saleRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "sale not found", err.Error())
}

func TestSaleService_GetByNumber_Success(t *testing.T) {
	service, saleRepo, _, _, _ := setupSaleService()
	ctx := context.Background()

	sale := &domain.Sale{ID: 1, SaleNumber: "SALE-001"}
	saleRepo.On("GetByNumber", ctx, "SALE-001").Return(sale, nil)

	result, err := service.GetByNumber(ctx, "SALE-001")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "SALE-001", result.SaleNumber)
	saleRepo.AssertExpectations(t)
}

func TestSaleService_GetByNumber_NotFound(t *testing.T) {
	service, saleRepo, _, _, _ := setupSaleService()
	ctx := context.Background()

	saleRepo.On("GetByNumber", ctx, "MISSING").Return(nil, errors.New("not found"))

	result, err := service.GetByNumber(ctx, "MISSING")

	assert.Error(t, err)
	assert.Nil(t, result)
	saleRepo.AssertExpectations(t)
}

func TestSaleService_List_Success(t *testing.T) {
	service, saleRepo, _, _, _ := setupSaleService()
	ctx := context.Background()

	params := repository.SaleListParams{
		BranchID:         1,
		PaginationParams: repository.PaginationParams{Page: 1, PerPage: 10},
	}

	result := &repository.PaginatedResult[domain.Sale]{
		Data:       []domain.Sale{{ID: 1}, {ID: 2}},
		Total:      2,
		Page:       1,
		PerPage:    10,
		TotalPages: 1,
	}

	saleRepo.On("List", ctx, params).Return(result, nil)

	res, err := service.List(ctx, params)

	assert.NoError(t, err)
	assert.Len(t, res.Data, 2)
	saleRepo.AssertExpectations(t)
}

func TestSaleService_Refund_Full_Success(t *testing.T) {
	service, saleRepo, itemRepo, _, _ := setupSaleService()
	ctx := context.Background()

	sale := &domain.Sale{ID: 1, ItemID: 10, FinalPrice: 500.0, Status: domain.SaleStatusCompleted, SaleNumber: "SALE-001"}
	item := &domain.Item{ID: 10, Status: domain.ItemStatusSold}

	saleRepo.On("GetByID", ctx, int64(1)).Return(sale, nil)
	itemRepo.On("GetByID", ctx, int64(10)).Return(item, nil)
	saleRepo.On("Update", ctx, mock.AnythingOfType("*domain.Sale")).Return(nil)
	itemRepo.On("UpdateStatus", ctx, int64(10), domain.ItemStatusAvailable).Return(nil)
	itemRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.ItemHistory")).Return(nil)

	input := RefundSaleInput{SaleID: 1, RefundAmount: 500.0, Reason: "Customer request", RefundedBy: 10}
	result, err := service.Refund(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.SaleStatusRefunded, result.Status)
	assert.NotNil(t, result.RefundAmount)
	assert.Equal(t, 500.0, *result.RefundAmount)
	saleRepo.AssertExpectations(t)
	itemRepo.AssertExpectations(t)
}

func TestSaleService_Refund_Partial_Success(t *testing.T) {
	service, saleRepo, itemRepo, _, _ := setupSaleService()
	ctx := context.Background()

	sale := &domain.Sale{ID: 1, ItemID: 10, FinalPrice: 500.0, Status: domain.SaleStatusCompleted}
	item := &domain.Item{ID: 10, Status: domain.ItemStatusSold}

	saleRepo.On("GetByID", ctx, int64(1)).Return(sale, nil)
	itemRepo.On("GetByID", ctx, int64(10)).Return(item, nil)
	saleRepo.On("Update", ctx, mock.AnythingOfType("*domain.Sale")).Return(nil)

	input := RefundSaleInput{SaleID: 1, RefundAmount: 200.0, Reason: "Partial refund", RefundedBy: 10}
	result, err := service.Refund(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.SaleStatusPartialRefund, result.Status)
}

func TestSaleService_Refund_NotCompleted(t *testing.T) {
	service, saleRepo, _, _, _ := setupSaleService()
	ctx := context.Background()

	sale := &domain.Sale{ID: 1, Status: domain.SaleStatusRefunded}
	saleRepo.On("GetByID", ctx, int64(1)).Return(sale, nil)

	input := RefundSaleInput{SaleID: 1, RefundAmount: 100.0, Reason: "test"}
	result, err := service.Refund(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "only completed sales can be refunded", err.Error())
}

func TestSaleService_Refund_AmountExceedsPrice(t *testing.T) {
	service, saleRepo, _, _, _ := setupSaleService()
	ctx := context.Background()

	sale := &domain.Sale{ID: 1, FinalPrice: 500.0, Status: domain.SaleStatusCompleted}
	saleRepo.On("GetByID", ctx, int64(1)).Return(sale, nil)

	input := RefundSaleInput{SaleID: 1, RefundAmount: 600.0, Reason: "test"}
	result, err := service.Refund(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "refund amount cannot exceed sale price", err.Error())
}

func TestSaleService_Refund_ItemStatusChanged(t *testing.T) {
	service, saleRepo, itemRepo, _, _ := setupSaleService()
	ctx := context.Background()

	sale := &domain.Sale{ID: 1, ItemID: 10, FinalPrice: 500.0, Status: domain.SaleStatusCompleted}
	item := &domain.Item{ID: 10, Status: domain.ItemStatusForSale}

	saleRepo.On("GetByID", ctx, int64(1)).Return(sale, nil)
	itemRepo.On("GetByID", ctx, int64(10)).Return(item, nil)

	input := RefundSaleInput{SaleID: 1, RefundAmount: 500.0, Reason: "test"}
	result, err := service.Refund(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "item status has changed, cannot process refund", err.Error())
}

func TestSaleService_GetSalesSummary_Success(t *testing.T) {
	service, saleRepo, _, _, _ := setupSaleService()
	ctx := context.Background()

	sales := []domain.Sale{
		{FinalPrice: 200.0, DiscountAmount: 10.0, PaymentMethod: "cash", Status: domain.SaleStatusCompleted},
		{FinalPrice: 300.0, DiscountAmount: 20.0, PaymentMethod: "card", Status: domain.SaleStatusCompleted},
	}

	saleRepo.On("List", ctx, mock.AnythingOfType("repository.SaleListParams")).Return(&repository.PaginatedResult[domain.Sale]{
		Data:  sales,
		Total: 2,
	}, nil)

	result, err := service.GetSalesSummary(ctx, 1, "2025-01-01", "2025-01-31")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.TotalSales)
	assert.Equal(t, 500.0, result.TotalRevenue)
	assert.Equal(t, 30.0, result.TotalDiscount)
	assert.Equal(t, 200.0, result.ByMethod["cash"])
	assert.Equal(t, 300.0, result.ByMethod["card"])
	saleRepo.AssertExpectations(t)
}
