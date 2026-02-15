package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository/mocks"
)

func setupLoyaltyService() (LoyaltyService, *mocks.MockCustomerRepository, *mocks.MockLoyaltyRepository) {
	customerRepo := new(mocks.MockCustomerRepository)
	loyaltyRepo := new(mocks.MockLoyaltyRepository)
	service := NewLoyaltyService(customerRepo, loyaltyRepo)
	return service, customerRepo, loyaltyRepo
}

func TestLoyaltyService_EnrollCustomer_Success(t *testing.T) {
	service, customerRepo, _ := setupLoyaltyService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:                 1,
		FirstName:          "John",
		LastName:           "Doe",
		LoyaltyEnrolledAt:  nil, // Not enrolled
		LoyaltyTier:        "",
		LoyaltyPoints:      0,
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	customerRepo.On("Update", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)

	err := service.EnrollCustomer(ctx, 1)

	assert.NoError(t, err)
	customerRepo.AssertExpectations(t)
}

func TestLoyaltyService_EnrollCustomer_AlreadyEnrolled(t *testing.T) {
	service, customerRepo, _ := setupLoyaltyService()
	ctx := context.Background()

	enrolledAt := time.Now()
	customer := &domain.Customer{
		ID:                1,
		LoyaltyEnrolledAt: &enrolledAt, // Already enrolled
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)

	err := service.EnrollCustomer(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, ErrAlreadyEnrolled, err)
	customerRepo.AssertExpectations(t)
}

func TestLoyaltyService_EnrollCustomer_CustomerNotFound(t *testing.T) {
	service, customerRepo, _ := setupLoyaltyService()
	ctx := context.Background()

	customerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := service.EnrollCustomer(ctx, 999)

	assert.Error(t, err)
	customerRepo.AssertExpectations(t)
}

func TestLoyaltyService_GetCustomerLoyalty_Success(t *testing.T) {
	service, customerRepo, _ := setupLoyaltyService()
	ctx := context.Background()

	enrolledAt := time.Now()
	customer := &domain.Customer{
		ID:                1,
		FirstName:         "John",
		LoyaltyEnrolledAt: &enrolledAt,
		LoyaltyTier:       domain.LoyaltyTierSilver,
		LoyaltyPoints:     1500,
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)

	info, err := service.GetCustomerLoyalty(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, int64(1), info.CustomerID)
	assert.Equal(t, 1500, info.Points)
	assert.Equal(t, domain.LoyaltyTierSilver, info.Tier)
	assert.True(t, info.IsEnrolled)
	assert.Equal(t, domain.LoyaltyTierGold, info.NextTier)
	customerRepo.AssertExpectations(t)
}

func TestLoyaltyService_GetCustomerLoyalty_NotEnrolled(t *testing.T) {
	service, customerRepo, _ := setupLoyaltyService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:                1,
		LoyaltyEnrolledAt: nil,
		LoyaltyTier:       "",
		LoyaltyPoints:     0,
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)

	info, err := service.GetCustomerLoyalty(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.False(t, info.IsEnrolled)
	assert.Equal(t, 0, info.Points)
	customerRepo.AssertExpectations(t)
}

func TestLoyaltyService_AddPoints_Success(t *testing.T) {
	service, customerRepo, loyaltyRepo := setupLoyaltyService()
	ctx := context.Background()

	enrolledAt := time.Now()
	customer := &domain.Customer{
		ID:                1,
		LoyaltyEnrolledAt: &enrolledAt,
		LoyaltyTier:       domain.LoyaltyTierStandard,
		LoyaltyPoints:     500,
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	customerRepo.On("Update", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)
	loyaltyRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.LoyaltyPointsHistory")).Return(nil)

	req := AddPointsRequest{
		CustomerID:  1,
		Points:      100,
		Description: "Test points",
	}
	history, err := service.AddPoints(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, history)
	assert.Equal(t, 100, history.PointsChange)
	assert.Equal(t, 600, history.PointsBalance)
	customerRepo.AssertExpectations(t)
	loyaltyRepo.AssertExpectations(t)
}

func TestLoyaltyService_AddPoints_NotEnrolled(t *testing.T) {
	service, customerRepo, _ := setupLoyaltyService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:                1,
		LoyaltyEnrolledAt: nil, // Not enrolled
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)

	req := AddPointsRequest{
		CustomerID: 1,
		Points:     100,
	}
	history, err := service.AddPoints(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, history)
	assert.Equal(t, ErrLoyaltyNotEnrolled, err)
	customerRepo.AssertExpectations(t)
}

func TestLoyaltyService_AddPoints_TierUpgrade(t *testing.T) {
	service, customerRepo, loyaltyRepo := setupLoyaltyService()
	ctx := context.Background()

	enrolledAt := time.Now()
	customer := &domain.Customer{
		ID:                1,
		LoyaltyEnrolledAt: &enrolledAt,
		LoyaltyTier:       domain.LoyaltyTierStandard,
		LoyaltyPoints:     900, // Close to Silver threshold (1000)
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	customerRepo.On("Update", ctx, mock.MatchedBy(func(c *domain.Customer) bool {
		return c.LoyaltyPoints == 1100 && c.LoyaltyTier == domain.LoyaltyTierSilver
	})).Return(nil)
	loyaltyRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.LoyaltyPointsHistory")).Return(nil)

	req := AddPointsRequest{
		CustomerID: 1,
		Points:     200, // Will push to 1100, triggering Silver tier
	}
	history, err := service.AddPoints(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, history)
	assert.Equal(t, 1100, history.PointsBalance)
	customerRepo.AssertExpectations(t)
	loyaltyRepo.AssertExpectations(t)
}

func TestLoyaltyService_RedeemPoints_Success(t *testing.T) {
	service, customerRepo, loyaltyRepo := setupLoyaltyService()
	ctx := context.Background()

	enrolledAt := time.Now()
	customer := &domain.Customer{
		ID:                1,
		LoyaltyEnrolledAt: &enrolledAt,
		LoyaltyTier:       domain.LoyaltyTierSilver,
		LoyaltyPoints:     1500,
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	customerRepo.On("Update", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)
	loyaltyRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.LoyaltyPointsHistory")).Return(nil)

	req := RedeemPointsRequest{
		CustomerID:  1,
		Points:      500,
		Description: "Discount redemption",
	}
	history, err := service.RedeemPoints(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, history)
	assert.Equal(t, -500, history.PointsChange) // Negative for redemption
	assert.Equal(t, 1000, history.PointsBalance)
	customerRepo.AssertExpectations(t)
	loyaltyRepo.AssertExpectations(t)
}

func TestLoyaltyService_RedeemPoints_InsufficientPoints(t *testing.T) {
	service, customerRepo, _ := setupLoyaltyService()
	ctx := context.Background()

	enrolledAt := time.Now()
	customer := &domain.Customer{
		ID:                1,
		LoyaltyEnrolledAt: &enrolledAt,
		LoyaltyPoints:     100,
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)

	req := RedeemPointsRequest{
		CustomerID: 1,
		Points:     500, // More than available
	}
	history, err := service.RedeemPoints(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, history)
	assert.Equal(t, ErrInsufficientPoints, err)
	customerRepo.AssertExpectations(t)
}

func TestLoyaltyService_RedeemPoints_NotEnrolled(t *testing.T) {
	service, customerRepo, _ := setupLoyaltyService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:                1,
		LoyaltyEnrolledAt: nil,
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)

	req := RedeemPointsRequest{
		CustomerID: 1,
		Points:     100,
	}
	history, err := service.RedeemPoints(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, history)
	assert.Equal(t, ErrLoyaltyNotEnrolled, err)
	customerRepo.AssertExpectations(t)
}

func TestLoyaltyService_GetPointsHistory_Success(t *testing.T) {
	service, _, loyaltyRepo := setupLoyaltyService()
	ctx := context.Background()

	history := []*domain.LoyaltyPointsHistory{
		{ID: 1, CustomerID: 1, PointsChange: 100, PointsBalance: 100},
		{ID: 2, CustomerID: 1, PointsChange: 50, PointsBalance: 150},
	}

	loyaltyRepo.On("GetHistoryByCustomer", ctx, int64(1), 1, 10).Return(history, int64(2), nil)

	result, total, err := service.GetPointsHistory(ctx, 1, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	loyaltyRepo.AssertExpectations(t)
}

func TestLoyaltyService_AwardLoanPoints_Success(t *testing.T) {
	service, customerRepo, loyaltyRepo := setupLoyaltyService()
	ctx := context.Background()

	enrolledAt := time.Now()
	customer := &domain.Customer{
		ID:                1,
		LoyaltyEnrolledAt: &enrolledAt,
		LoyaltyTier:       domain.LoyaltyTierStandard,
		LoyaltyPoints:     0,
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	customerRepo.On("Update", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)
	loyaltyRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.LoyaltyPointsHistory")).Return(nil)

	err := service.AwardLoanPoints(ctx, 1, 100, 500.00, nil) // $500 loan = 500 points

	assert.NoError(t, err)
	customerRepo.AssertExpectations(t)
	loyaltyRepo.AssertExpectations(t)
}

func TestLoyaltyService_AwardPaymentPoints_Success(t *testing.T) {
	service, customerRepo, loyaltyRepo := setupLoyaltyService()
	ctx := context.Background()

	enrolledAt := time.Now()
	customer := &domain.Customer{
		ID:                1,
		LoyaltyEnrolledAt: &enrolledAt,
		LoyaltyTier:       domain.LoyaltyTierStandard,
		LoyaltyPoints:     100,
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	customerRepo.On("Update", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)
	loyaltyRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.LoyaltyPointsHistory")).Return(nil)

	err := service.AwardPaymentPoints(ctx, 1, 50, 100.00, nil) // $100 payment = 200 points (2x)

	assert.NoError(t, err)
	customerRepo.AssertExpectations(t)
	loyaltyRepo.AssertExpectations(t)
}

func TestLoyaltyService_CalculateDiscount_Success(t *testing.T) {
	service, customerRepo, _ := setupLoyaltyService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:          1,
		LoyaltyTier: domain.LoyaltyTierGold, // 5% discount
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)

	discount, err := service.CalculateDiscount(ctx, 1, 100.00)

	assert.NoError(t, err)
	assert.Equal(t, 5.0, discount) // 5% of $100
	customerRepo.AssertExpectations(t)
}

func TestLoyaltyService_CalculateDiscount_NoDiscount(t *testing.T) {
	service, customerRepo, _ := setupLoyaltyService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:          1,
		LoyaltyTier: domain.LoyaltyTierStandard, // 0% discount
	}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)

	discount, err := service.CalculateDiscount(ctx, 1, 100.00)

	assert.NoError(t, err)
	assert.Equal(t, 0.0, discount)
	customerRepo.AssertExpectations(t)
}
