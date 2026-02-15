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

func setupPaymentService() (*PaymentService, *mocks.MockPaymentRepository, *mocks.MockLoanRepository, *mocks.MockCustomerRepository) {
	paymentRepo := new(mocks.MockPaymentRepository)
	loanRepo := new(mocks.MockLoanRepository)
	customerRepo := new(mocks.MockCustomerRepository)
	service := NewPaymentService(paymentRepo, loanRepo, customerRepo)
	return service, paymentRepo, loanRepo, customerRepo
}

// --- Create tests ---

func TestPaymentService_Create_Success_PartialPayment(t *testing.T) {
	service, paymentRepo, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:                 1,
		CustomerID:         10,
		Status:             domain.LoanStatusActive,
		PrincipalRemaining: 800,
		InterestRemaining:  100,
		LateFeeAmount:      20,
		AmountPaid:         0,
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)
	paymentRepo.On("GenerateNumber", ctx).Return("PAY-000001", nil)
	paymentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)
	loanRepo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)

	input := CreatePaymentInput{
		LoanID:        1,
		Amount:        50,
		PaymentMethod: "cash",
		BranchID:      1,
		CreatedBy:     1,
	}

	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsFullyPaid)
	// 50 applied: 20 to late fees, 30 to interest
	assert.Equal(t, "PAY-000001", result.Payment.PaymentNumber)
	assert.Equal(t, 20.0, result.Payment.LateFeeAmount)
	assert.Equal(t, 30.0, result.Payment.InterestAmount)
	assert.Equal(t, 0.0, result.Payment.PrincipalAmount)
	paymentRepo.AssertExpectations(t)
	loanRepo.AssertExpectations(t)
}

func TestPaymentService_Create_Success_FullPayment(t *testing.T) {
	service, paymentRepo, loanRepo, customerRepo := setupPaymentService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:                 1,
		CustomerID:         10,
		Status:             domain.LoanStatusActive,
		PrincipalRemaining: 100,
		InterestRemaining:  20,
		LateFeeAmount:      0,
		AmountPaid:         880,
	}

	customer := &domain.Customer{ID: 10, TotalPaid: 500}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)
	paymentRepo.On("GenerateNumber", ctx).Return("PAY-000002", nil)
	paymentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)
	loanRepo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)
	customerRepo.On("GetByID", ctx, int64(10)).Return(customer, nil)
	customerRepo.On("UpdateCreditInfo", ctx, int64(10), mock.AnythingOfType("repository.CustomerCreditUpdate")).Return(nil)

	input := CreatePaymentInput{
		LoanID:        1,
		Amount:        120, // Exact remaining: 100 principal + 20 interest
		PaymentMethod: "card",
		BranchID:      1,
		CreatedBy:     1,
	}

	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsFullyPaid)
	assert.Equal(t, 0.0, result.RemainingBalance)
	assert.Equal(t, domain.LoanStatusPaid, result.Loan.Status)
	assert.NotNil(t, result.Loan.PaidDate)
	customerRepo.AssertExpectations(t)
}

func TestPaymentService_Create_AllocatesLateFeeFirst(t *testing.T) {
	service, paymentRepo, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:                 1,
		CustomerID:         10,
		Status:             domain.LoanStatusOverdue,
		PrincipalRemaining: 500,
		InterestRemaining:  100,
		LateFeeAmount:      50,
		AmountPaid:         0,
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)
	paymentRepo.On("GenerateNumber", ctx).Return("PAY-000003", nil)
	paymentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)
	loanRepo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)

	input := CreatePaymentInput{
		LoanID:        1,
		Amount:        200,
		PaymentMethod: "cash",
		BranchID:      1,
		CreatedBy:     1,
	}

	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// 200 applied: 50 late fees, 100 interest, 50 principal
	assert.Equal(t, 50.0, result.Payment.LateFeeAmount)
	assert.Equal(t, 100.0, result.Payment.InterestAmount)
	assert.Equal(t, 50.0, result.Payment.PrincipalAmount)
	assert.False(t, result.IsFullyPaid)
}

func TestPaymentService_Create_LoanAlreadyPaid(t *testing.T) {
	service, _, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:     1,
		Status: domain.LoanStatusPaid,
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)

	input := CreatePaymentInput{LoanID: 1, Amount: 100, PaymentMethod: "cash"}

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "loan is already fully paid", err.Error())
}

func TestPaymentService_Create_LoanConfiscated(t *testing.T) {
	service, _, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:     1,
		Status: domain.LoanStatusConfiscated,
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)

	input := CreatePaymentInput{LoanID: 1, Amount: 100, PaymentMethod: "cash"}

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "loan has been confiscated", err.Error())
}

func TestPaymentService_Create_LoanNotFound(t *testing.T) {
	service, _, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	loanRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreatePaymentInput{LoanID: 999, Amount: 100, PaymentMethod: "cash"}

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "loan not found", err.Error())
}

func TestPaymentService_Create_GenerateNumberError(t *testing.T) {
	service, paymentRepo, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:                 1,
		CustomerID:         10,
		Status:             domain.LoanStatusActive,
		PrincipalRemaining: 100,
		InterestRemaining:  0,
		LateFeeAmount:      0,
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)
	paymentRepo.On("GenerateNumber", ctx).Return("", errors.New("db error"))

	input := CreatePaymentInput{LoanID: 1, Amount: 50, PaymentMethod: "cash"}

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to generate payment number", err.Error())
}

func TestPaymentService_Create_PaymentRepoError(t *testing.T) {
	service, paymentRepo, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:                 1,
		CustomerID:         10,
		Status:             domain.LoanStatusActive,
		PrincipalRemaining: 100,
		InterestRemaining:  0,
		LateFeeAmount:      0,
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)
	paymentRepo.On("GenerateNumber", ctx).Return("PAY-000001", nil)
	paymentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(errors.New("db error"))

	input := CreatePaymentInput{LoanID: 1, Amount: 50, PaymentMethod: "cash", BranchID: 1, CreatedBy: 1}

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to create payment", err.Error())
}

func TestPaymentService_Create_OverpaymentAppliedCorrectly(t *testing.T) {
	service, paymentRepo, loanRepo, customerRepo := setupPaymentService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:                 1,
		CustomerID:         10,
		Status:             domain.LoanStatusActive,
		PrincipalRemaining: 100,
		InterestRemaining:  50,
		LateFeeAmount:      0,
		AmountPaid:         850,
	}

	customer := &domain.Customer{ID: 10, TotalPaid: 0}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)
	paymentRepo.On("GenerateNumber", ctx).Return("PAY-000004", nil)
	paymentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)
	loanRepo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)
	customerRepo.On("GetByID", ctx, int64(10)).Return(customer, nil)
	customerRepo.On("UpdateCreditInfo", ctx, int64(10), mock.AnythingOfType("repository.CustomerCreditUpdate")).Return(nil)

	// Overpay by 50 (200 total but only 150 remaining)
	input := CreatePaymentInput{
		LoanID:        1,
		Amount:        200,
		PaymentMethod: "cash",
		BranchID:      1,
		CreatedBy:     1,
	}

	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsFullyPaid)
	assert.Equal(t, 50.0, result.Payment.InterestAmount)
	assert.Equal(t, 100.0, result.Payment.PrincipalAmount)
}

// --- GetByID tests ---

func TestPaymentService_GetByID_Success(t *testing.T) {
	service, paymentRepo, loanRepo, customerRepo := setupPaymentService()
	ctx := context.Background()

	payment := &domain.Payment{
		ID:            1,
		PaymentNumber: "PAY-000001",
		Amount:        100.00,
		LoanID:        10,
		CustomerID:    20,
	}

	loan := &domain.Loan{ID: 10, LoanNumber: "LN-000001"}
	customer := &domain.Customer{ID: 20, FirstName: "John", LastName: "Doe"}

	paymentRepo.On("GetByID", ctx, int64(1)).Return(payment, nil)
	loanRepo.On("GetByID", ctx, int64(10)).Return(loan, nil)
	customerRepo.On("GetByID", ctx, int64(20)).Return(customer, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "PAY-000001", result.PaymentNumber)
}

func TestPaymentService_GetByID_NotFound(t *testing.T) {
	service, paymentRepo, _, _ := setupPaymentService()
	ctx := context.Background()

	paymentRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// --- List tests ---

func TestPaymentService_List_Success(t *testing.T) {
	service, paymentRepo, _, _ := setupPaymentService()
	ctx := context.Background()

	payments := []domain.Payment{
		{ID: 1, PaymentNumber: "PAY-000001", Amount: 100.00},
		{ID: 2, PaymentNumber: "PAY-000002", Amount: 150.00},
	}

	paginatedResult := &repository.PaginatedResult[domain.Payment]{
		Data:       payments,
		Total:      2,
		Page:       1,
		PerPage:    10,
		TotalPages: 1,
	}

	paymentRepo.On("List", ctx, mock.AnythingOfType("repository.PaymentListParams")).Return(paginatedResult, nil)

	params := repository.PaymentListParams{BranchID: 1}
	result, err := service.List(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
}

// --- Reverse tests ---

func TestPaymentService_Reverse_Success(t *testing.T) {
	service, paymentRepo, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	payment := &domain.Payment{
		ID:               1,
		LoanID:           10,
		Amount:           200,
		PrincipalAmount:  100,
		InterestAmount:   80,
		LateFeeAmount:    20,
		Status:           domain.PaymentStatusCompleted,
	}

	loan := &domain.Loan{
		ID:                 10,
		Status:             domain.LoanStatusActive,
		PrincipalRemaining: 400,
		InterestRemaining:  0,
		LateFeeAmount:      0,
		AmountPaid:         600,
	}

	paymentRepo.On("GetByID", ctx, int64(1)).Return(payment, nil)
	loanRepo.On("GetByID", ctx, int64(10)).Return(loan, nil)
	loanRepo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)
	paymentRepo.On("Update", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)

	input := ReversePaymentInput{
		PaymentID:  1,
		Reason:     "Customer dispute",
		ReversedBy: 1,
	}

	result, err := service.Reverse(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.PaymentStatusReversed, result.Status)
	assert.NotNil(t, result.ReversedAt)
	assert.Equal(t, "Customer dispute", result.ReversalReason)

	// Verify loan balances were restored
	assert.Equal(t, 500.0, loan.PrincipalRemaining)  // 400 + 100
	assert.Equal(t, 80.0, loan.InterestRemaining)     // 0 + 80
	assert.Equal(t, 20.0, loan.LateFeeAmount)         // 0 + 20
	assert.Equal(t, 400.0, loan.AmountPaid)            // 600 - 200
	loanRepo.AssertExpectations(t)
	paymentRepo.AssertExpectations(t)
}

func TestPaymentService_Reverse_ReactivatesPaidLoan(t *testing.T) {
	service, paymentRepo, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	now := time.Now()
	payment := &domain.Payment{
		ID:               1,
		LoanID:           10,
		Amount:           100,
		PrincipalAmount:  80,
		InterestAmount:   20,
		LateFeeAmount:    0,
		Status:           domain.PaymentStatusCompleted,
	}

	loan := &domain.Loan{
		ID:                 10,
		Status:             domain.LoanStatusPaid,
		PaidDate:           &now,
		PrincipalRemaining: 0,
		InterestRemaining:  0,
		LateFeeAmount:      0,
		AmountPaid:         1000,
	}

	paymentRepo.On("GetByID", ctx, int64(1)).Return(payment, nil)
	loanRepo.On("GetByID", ctx, int64(10)).Return(loan, nil)
	loanRepo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)
	paymentRepo.On("Update", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)

	input := ReversePaymentInput{
		PaymentID:  1,
		Reason:     "Error",
		ReversedBy: 1,
	}

	result, err := service.Reverse(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.LoanStatusActive, loan.Status)
	assert.Nil(t, loan.PaidDate)
}

func TestPaymentService_Reverse_NotFound(t *testing.T) {
	service, paymentRepo, _, _ := setupPaymentService()
	ctx := context.Background()

	paymentRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := ReversePaymentInput{PaymentID: 999, Reason: "test"}

	result, err := service.Reverse(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "payment not found", err.Error())
}

func TestPaymentService_Reverse_AlreadyReversed(t *testing.T) {
	service, paymentRepo, _, _ := setupPaymentService()
	ctx := context.Background()

	now := time.Now()
	payment := &domain.Payment{
		ID:         1,
		Status:     domain.PaymentStatusReversed,
		ReversedAt: &now,
	}

	paymentRepo.On("GetByID", ctx, int64(1)).Return(payment, nil)

	input := ReversePaymentInput{PaymentID: 1, Reason: "test"}

	result, err := service.Reverse(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "payment cannot be reversed", err.Error())
}

// --- CalculatePayoff tests ---

func TestPaymentService_CalculatePayoff_Success(t *testing.T) {
	service, _, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:                 1,
		PrincipalRemaining: 800.00,
		InterestRemaining:  50.00,
		LateFeeAmount:      10.00,
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)

	result, err := service.CalculatePayoff(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 860.00, result)
}

func TestPaymentService_CalculatePayoff_NotFound(t *testing.T) {
	service, _, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	loanRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.CalculatePayoff(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, 0.0, result)
}

// --- CalculateMinimumPayment tests ---

func TestPaymentService_CalculateMinimumPayment_Success(t *testing.T) {
	service, _, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	minPayment := 100.00
	loan := &domain.Loan{
		ID:                     1,
		PrincipalRemaining:     800.00,
		InterestRemaining:      50.00,
		LateFeeAmount:          10.00,
		RequiresMinimumPayment: true,
		MinimumPaymentAmount:   &minPayment,
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)

	result, err := service.CalculateMinimumPayment(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 110.00, result) // MinPayment + LateFee
}

func TestPaymentService_CalculateMinimumPayment_NoMinimum(t *testing.T) {
	service, _, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:                     1,
		PrincipalRemaining:     200,
		InterestRemaining:      30,
		LateFeeAmount:          0,
		RequiresMinimumPayment: false,
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)

	result, err := service.CalculateMinimumPayment(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 230.0, result) // Full remaining balance
}

func TestPaymentService_CalculateMinimumPayment_BalanceLessThanMinimum(t *testing.T) {
	service, _, loanRepo, _ := setupPaymentService()
	ctx := context.Background()

	minPayment := 500.0
	loan := &domain.Loan{
		ID:                     1,
		PrincipalRemaining:     50,
		InterestRemaining:      10,
		LateFeeAmount:          0,
		RequiresMinimumPayment: true,
		MinimumPaymentAmount:   &minPayment,
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)

	result, err := service.CalculateMinimumPayment(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 60.0, result) // Remaining balance < MinPayment, return balance
}
