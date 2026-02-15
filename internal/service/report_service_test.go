package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/internal/repository/mocks"
)

func setupReportService() (*ReportService, *mocks.MockLoanRepository, *mocks.MockPaymentRepository, *mocks.MockSaleRepository, *mocks.MockCustomerRepository, *mocks.MockItemRepository) {
	loanRepo := new(mocks.MockLoanRepository)
	paymentRepo := new(mocks.MockPaymentRepository)
	saleRepo := new(mocks.MockSaleRepository)
	customerRepo := new(mocks.MockCustomerRepository)
	itemRepo := new(mocks.MockItemRepository)
	service := NewReportService(loanRepo, paymentRepo, saleRepo, customerRepo, itemRepo, nil)
	return service, loanRepo, paymentRepo, saleRepo, customerRepo, itemRepo
}

func TestReportService_GetDashboardStats_Success(t *testing.T) {
	service, loanRepo, paymentRepo, saleRepo, customerRepo, itemRepo := setupReportService()
	ctx := context.Background()

	// Active loans
	loanRepo.On("List", ctx, mock.AnythingOfType("repository.LoanListParams")).Return(&repository.PaginatedResult[domain.Loan]{
		Data:  []domain.Loan{},
		Total: 0,
	}, nil)

	// Payments
	paymentRepo.On("List", ctx, mock.AnythingOfType("repository.PaymentListParams")).Return(&repository.PaginatedResult[domain.Payment]{
		Data:  []domain.Payment{},
		Total: 0,
	}, nil)

	// Sales
	saleRepo.On("List", ctx, mock.AnythingOfType("repository.SaleListParams")).Return(&repository.PaginatedResult[domain.Sale]{
		Data:  []domain.Sale{},
		Total: 0,
	}, nil)

	// Items
	itemRepo.On("List", ctx, mock.AnythingOfType("repository.ItemListParams")).Return(&repository.PaginatedResult[domain.Item]{
		Data:  []domain.Item{},
		Total: 0,
	}, nil)

	// Customers
	customerRepo.On("List", ctx, mock.AnythingOfType("repository.CustomerListParams")).Return(&repository.PaginatedResult[domain.Customer]{
		Data:  []domain.Customer{},
		Total: 5,
	}, nil)

	result, err := service.GetDashboardStats(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 5, result.TotalCustomers)
}

func TestReportService_GetDashboardStats_WithData(t *testing.T) {
	service, loanRepo, paymentRepo, saleRepo, customerRepo, itemRepo := setupReportService()
	ctx := context.Background()

	// Active loans query, overdue loans query, all loans query
	loanRepo.On("List", ctx, mock.AnythingOfType("repository.LoanListParams")).Return(&repository.PaginatedResult[domain.Loan]{
		Data:  []domain.Loan{},
		Total: 3,
	}, nil)

	// Today's payments
	paymentRepo.On("List", ctx, mock.AnythingOfType("repository.PaymentListParams")).Return(&repository.PaginatedResult[domain.Payment]{
		Data: []domain.Payment{
			{Amount: 100.0, Status: domain.PaymentStatusCompleted},
			{Amount: 200.0, Status: domain.PaymentStatusCompleted},
		},
		Total: 2,
	}, nil)

	// Today's sales
	saleRepo.On("List", ctx, mock.AnythingOfType("repository.SaleListParams")).Return(&repository.PaginatedResult[domain.Sale]{
		Data: []domain.Sale{
			{FinalPrice: 500.0, Status: domain.SaleStatusCompleted},
		},
		Total: 1,
	}, nil)

	// Items
	itemRepo.On("List", ctx, mock.AnythingOfType("repository.ItemListParams")).Return(&repository.PaginatedResult[domain.Item]{
		Data: []domain.Item{
			{Status: domain.ItemStatusPawned, AppraisedValue: 1000.0},
			{Status: domain.ItemStatusForSale, AppraisedValue: 500.0},
		},
		Total: 2,
	}, nil)

	// Customers
	customerRepo.On("List", ctx, mock.AnythingOfType("repository.CustomerListParams")).Return(&repository.PaginatedResult[domain.Customer]{
		Total: 10,
	}, nil)

	result, err := service.GetDashboardStats(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.PaymentsToday)
	assert.Equal(t, 300.0, result.PaymentAmount)
	assert.Equal(t, 1, result.SalesToday)
	assert.Equal(t, 500.0, result.SalesAmountToday)
	assert.Equal(t, 1, result.ItemsInPawn)
	assert.Equal(t, 1, result.ItemsForSale)
	assert.Equal(t, 1500.0, result.InventoryValue)
	assert.Equal(t, 10, result.TotalCustomers)
}

func TestReportService_GetLoanReport_Success(t *testing.T) {
	service, loanRepo, _, _, _, _ := setupReportService()
	ctx := context.Background()

	loans := []domain.Loan{
		{
			LoanAmount:         1000.0,
			InterestAmount:     100.0,
			TotalAmount:        1100.0,
			PrincipalRemaining: 500.0,
			InterestRemaining:  50.0,
			LateFeeAmount:      10.0,
			Status:             domain.LoanStatusActive,
		},
		{
			LoanAmount:         2000.0,
			InterestAmount:     200.0,
			TotalAmount:        2200.0,
			PrincipalRemaining: 0,
			InterestRemaining:  0,
			LateFeeAmount:      0,
			Status:             domain.LoanStatusPaid,
		},
	}

	loanRepo.On("List", ctx, mock.AnythingOfType("repository.LoanListParams")).Return(&repository.PaginatedResult[domain.Loan]{
		Data:  loans,
		Total: 2,
	}, nil)

	result, err := service.GetLoanReport(ctx, 1, "2025-01-01", "2025-01-31")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.TotalLoans)
	assert.Equal(t, 3000.0, result.TotalAmount)
	assert.Equal(t, 300.0, result.TotalInterest)
	assert.Equal(t, 560.0, result.TotalOutstanding) // 500 + 50 + 10 from active loan
	assert.Equal(t, 1, result.ByStatus["active"])
	assert.Equal(t, 1, result.ByStatus["paid"])
	assert.Len(t, result.RecentLoans, 2)
}

func TestReportService_GetLoanReport_MoreThan10Loans(t *testing.T) {
	service, loanRepo, _, _, _, _ := setupReportService()
	ctx := context.Background()

	loans := make([]domain.Loan, 15)
	for i := range loans {
		loans[i] = domain.Loan{LoanAmount: 100.0, InterestAmount: 10.0, TotalAmount: 110.0, Status: domain.LoanStatusPaid}
	}

	loanRepo.On("List", ctx, mock.AnythingOfType("repository.LoanListParams")).Return(&repository.PaginatedResult[domain.Loan]{
		Data:  loans,
		Total: 15,
	}, nil)

	result, err := service.GetLoanReport(ctx, 1, "2025-01-01", "2025-12-31")

	assert.NoError(t, err)
	assert.Equal(t, 15, result.TotalLoans)
	assert.Len(t, result.RecentLoans, 10)
}

func TestReportService_GetLoanReport_Error(t *testing.T) {
	service, loanRepo, _, _, _, _ := setupReportService()
	ctx := context.Background()

	loanRepo.On("List", ctx, mock.AnythingOfType("repository.LoanListParams")).Return(nil, errors.New("db error"))

	result, err := service.GetLoanReport(ctx, 1, "2025-01-01", "2025-01-31")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestReportService_GetPaymentReport_Success(t *testing.T) {
	service, _, paymentRepo, _, _, _ := setupReportService()
	ctx := context.Background()

	payments := []domain.Payment{
		{
			Amount:          500.0,
			PrincipalAmount: 400.0,
			InterestAmount:  80.0,
			LateFeeAmount:   20.0,
			PaymentMethod:   "cash",
			Status:          domain.PaymentStatusCompleted,
		},
		{
			Amount:          300.0,
			PrincipalAmount: 250.0,
			InterestAmount:  50.0,
			LateFeeAmount:   0,
			PaymentMethod:   "card",
			Status:          domain.PaymentStatusCompleted,
		},
		{
			Amount:        100.0,
			PaymentMethod: "cash",
			Status:        domain.PaymentStatusPending,
		},
	}

	paymentRepo.On("List", ctx, mock.AnythingOfType("repository.PaymentListParams")).Return(&repository.PaginatedResult[domain.Payment]{
		Data:  payments,
		Total: 3,
	}, nil)

	result, err := service.GetPaymentReport(ctx, 1, "2025-01-01", "2025-01-31")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, result.TotalPayments)
	assert.Equal(t, 800.0, result.TotalAmount)
	assert.Equal(t, 650.0, result.TotalPrincipal)
	assert.Equal(t, 130.0, result.TotalInterest)
	assert.Equal(t, 20.0, result.TotalLateFees)
	assert.Equal(t, 1, result.ByMethod["cash"])
	assert.Equal(t, 1, result.ByMethod["card"])
	assert.Equal(t, 500.0, result.ByMethodAmount["cash"])
	assert.Equal(t, 300.0, result.ByMethodAmount["card"])
}

func TestReportService_GetPaymentReport_Error(t *testing.T) {
	service, _, paymentRepo, _, _, _ := setupReportService()
	ctx := context.Background()

	paymentRepo.On("List", ctx, mock.AnythingOfType("repository.PaymentListParams")).Return(nil, errors.New("db error"))

	result, err := service.GetPaymentReport(ctx, 1, "2025-01-01", "2025-01-31")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestReportService_GetSalesReport_Success(t *testing.T) {
	service, _, _, saleRepo, _, _ := setupReportService()
	ctx := context.Background()

	sales := []domain.Sale{
		{
			SalePrice:      600.0,
			DiscountAmount: 50.0,
			FinalPrice:     550.0,
			PaymentMethod:  "cash",
			Status:         domain.SaleStatusCompleted,
		},
		{
			SalePrice:      400.0,
			DiscountAmount: 0,
			FinalPrice:     400.0,
			PaymentMethod:  "card",
			Status:         domain.SaleStatusCompleted,
		},
		{
			Status: domain.SaleStatusRefunded,
		},
	}

	saleRepo.On("List", ctx, mock.AnythingOfType("repository.SaleListParams")).Return(&repository.PaginatedResult[domain.Sale]{
		Data:  sales,
		Total: 3,
	}, nil)

	result, err := service.GetSalesReport(ctx, 1, "2025-01-01", "2025-01-31")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, result.TotalSales)
	assert.Equal(t, 1000.0, result.TotalAmount)
	assert.Equal(t, 50.0, result.TotalDiscounts)
	assert.Equal(t, 950.0, result.NetAmount)
	assert.Equal(t, 2, result.ByStatus["completed"])
	assert.Equal(t, 1, result.ByStatus["refunded"])
	assert.Equal(t, 1, result.ByMethod["cash"])
	assert.Equal(t, 1, result.ByMethod["card"])
}

func TestReportService_GetSalesReport_Error(t *testing.T) {
	service, _, _, saleRepo, _, _ := setupReportService()
	ctx := context.Background()

	saleRepo.On("List", ctx, mock.AnythingOfType("repository.SaleListParams")).Return(nil, errors.New("db error"))

	result, err := service.GetSalesReport(ctx, 1, "2025-01-01", "2025-01-31")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestReportService_GetOverdueReport_Success(t *testing.T) {
	service, loanRepo, _, _, _, _ := setupReportService()
	ctx := context.Background()

	now := time.Now()
	overdueLoans := []*domain.Loan{
		{
			ID:                 1,
			LoanAmount:         1000.0,
			PrincipalRemaining: 700.0,
			InterestRemaining:  50.0,
			LateFeeAmount:      50.0,
			Status:             domain.LoanStatusOverdue,
			DueDate:            now.AddDate(0, 0, -10),
			GracePeriodDays:    15,
		},
		{
			ID:                 2,
			LoanAmount:         500.0,
			PrincipalRemaining: 500.0,
			InterestRemaining:  25.0,
			LateFeeAmount:      25.0,
			Status:             domain.LoanStatusOverdue,
			DueDate:            now.AddDate(0, 0, -14),
			GracePeriodDays:    15,
		},
	}

	loanRepo.On("GetOverdueLoans", ctx, int64(1)).Return(overdueLoans, nil)

	// Active loans for approaching due
	loanRepo.On("List", ctx, mock.AnythingOfType("repository.LoanListParams")).Return(&repository.PaginatedResult[domain.Loan]{
		Data: []domain.Loan{
			{DueDate: now.AddDate(0, 0, 3), Status: domain.LoanStatusActive},
		},
		Total: 1,
	}, nil)

	result, err := service.GetOverdueReport(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.TotalOverdue)
	assert.Len(t, result.OverdueLoans, 2)
	assert.Equal(t, 75.0, result.TotalLateFees)
	assert.Len(t, result.ApproachingDue, 1)
}

func TestReportService_GetOverdueReport_Error(t *testing.T) {
	service, loanRepo, _, _, _, _ := setupReportService()
	ctx := context.Background()

	loanRepo.On("GetOverdueLoans", ctx, int64(1)).Return(nil, errors.New("db error"))

	result, err := service.GetOverdueReport(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestReportService_GetOverdueReport_NoOverdue(t *testing.T) {
	service, loanRepo, _, _, _, _ := setupReportService()
	ctx := context.Background()

	loanRepo.On("GetOverdueLoans", ctx, int64(1)).Return([]*domain.Loan{}, nil)
	loanRepo.On("List", ctx, mock.AnythingOfType("repository.LoanListParams")).Return(&repository.PaginatedResult[domain.Loan]{
		Data:  []domain.Loan{},
		Total: 0,
	}, nil)

	result, err := service.GetOverdueReport(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 0, result.TotalOverdue)
	assert.Empty(t, result.OverdueLoans)
	assert.Empty(t, result.ApproachingDue)
}
