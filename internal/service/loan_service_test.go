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

func adultBirthDate() *time.Time {
	t := time.Now().AddDate(-30, 0, 0)
	return &t
}

func setupLoanService() (*LoanService, *mocks.MockLoanRepository, *mocks.MockItemRepository, *mocks.MockCustomerRepository, *mocks.MockPaymentRepository) {
	loanRepo := new(mocks.MockLoanRepository)
	itemRepo := new(mocks.MockItemRepository)
	customerRepo := new(mocks.MockCustomerRepository)
	paymentRepo := new(mocks.MockPaymentRepository)
	service := NewLoanService(loanRepo, itemRepo, customerRepo, paymentRepo)
	return service, loanRepo, itemRepo, customerRepo, paymentRepo
}

// --- Create tests ---

func TestLoanService_Create_Success(t *testing.T) {
	service, loanRepo, itemRepo, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		IsActive:  true,
		BirthDate: adultBirthDate(),
	}
	item := &domain.Item{
		ID:        1,
		Name:      "iPhone 15",
		Status:    domain.ItemStatusAvailable,
		LoanValue: 1000,
	}

	tx := new(mocks.MockTransaction)

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	loanRepo.On("GenerateNumber", ctx).Return("LN-000001", nil)
	loanRepo.On("BeginTx", ctx).Return(tx, nil)
	loanRepo.On("CreateTx", ctx, tx, mock.AnythingOfType("*domain.Loan")).Return(nil)
	itemRepo.On("UpdateStatus", ctx, int64(1), domain.ItemStatusCollateral).Return(nil)
	tx.On("Commit").Return(nil)
	tx.On("Rollback").Return(nil)
	customerRepo.On("UpdateCreditInfo", ctx, int64(1), mock.AnythingOfType("repository.CustomerCreditUpdate")).Return(nil)

	input := CreateLoanInput{
		CustomerID:      1,
		ItemID:          1,
		BranchID:        1,
		LoanAmount:      800,
		InterestRate:    10,
		LoanTermDays:    30,
		PaymentPlanType: "single",
		CreatedBy:       1,
	}

	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LN-000001", result.LoanNumber)
	assert.Equal(t, domain.LoanStatusActive, result.Status)
	assert.Equal(t, 800.0, result.LoanAmount)
	assert.Equal(t, 80.0, result.InterestAmount) // 800 * 10/100
	assert.Equal(t, 880.0, result.TotalAmount)
	assert.NotNil(t, result.Customer)
	assert.NotNil(t, result.Item)
	loanRepo.AssertExpectations(t)
	tx.AssertExpectations(t)
}

func TestLoanService_Create_WithInstallments(t *testing.T) {
	service, loanRepo, itemRepo, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:        1,
		IsActive:  true,
		BirthDate: adultBirthDate(),
	}
	item := &domain.Item{
		ID:        1,
		Status:    domain.ItemStatusAvailable,
		LoanValue: 1000,
	}

	tx := new(mocks.MockTransaction)

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	loanRepo.On("GenerateNumber", ctx).Return("LN-000002", nil)
	loanRepo.On("BeginTx", ctx).Return(tx, nil)
	loanRepo.On("CreateTx", ctx, tx, mock.AnythingOfType("*domain.Loan")).Return(nil)
	itemRepo.On("UpdateStatus", ctx, int64(1), domain.ItemStatusCollateral).Return(nil)
	loanRepo.On("CreateInstallments", ctx, mock.AnythingOfType("[]*domain.LoanInstallment")).Return(nil)
	tx.On("Commit").Return(nil)
	tx.On("Rollback").Return(nil)
	customerRepo.On("UpdateCreditInfo", ctx, int64(1), mock.AnythingOfType("repository.CustomerCreditUpdate")).Return(nil)

	input := CreateLoanInput{
		CustomerID:           1,
		ItemID:               1,
		BranchID:             1,
		LoanAmount:           600,
		InterestRate:         12,
		LoanTermDays:         90,
		PaymentPlanType:      "installments",
		NumberOfInstallments: 3,
		CreatedBy:            1,
	}

	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	loanRepo.AssertExpectations(t)
}

func TestLoanService_Create_CustomerNotFound(t *testing.T) {
	service, _, _, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	customerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreateLoanInput{CustomerID: 999, ItemID: 1}

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "customer not found", err.Error())
}

func TestLoanService_Create_CustomerCannotTakeLoan(t *testing.T) {
	service, _, _, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:        1,
		IsBlocked: true,
	}
	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)

	input := CreateLoanInput{CustomerID: 1, ItemID: 1}

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "customer cannot take loans", err.Error())
}

func TestLoanService_Create_ItemNotFound(t *testing.T) {
	service, _, itemRepo, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	customer := &domain.Customer{ID: 1, IsActive: true, BirthDate: adultBirthDate()}
	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	itemRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreateLoanInput{CustomerID: 1, ItemID: 999}

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "item not found", err.Error())
}

func TestLoanService_Create_ItemNotAvailable(t *testing.T) {
	service, _, itemRepo, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	customer := &domain.Customer{ID: 1, IsActive: true, BirthDate: adultBirthDate()}
	item := &domain.Item{ID: 1, Status: domain.ItemStatusPawned}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	input := CreateLoanInput{CustomerID: 1, ItemID: 1}

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "item is not available for loan", err.Error())
}

func TestLoanService_Create_LoanExceedsItemValue(t *testing.T) {
	service, _, itemRepo, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	customer := &domain.Customer{ID: 1, IsActive: true, BirthDate: adultBirthDate()}
	item := &domain.Item{ID: 1, Status: domain.ItemStatusAvailable, LoanValue: 500}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	input := CreateLoanInput{CustomerID: 1, ItemID: 1, LoanAmount: 1000}

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "loan amount cannot exceed item loan value", err.Error())
}

func TestLoanService_Create_GenerateNumberError(t *testing.T) {
	service, loanRepo, itemRepo, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	customer := &domain.Customer{ID: 1, IsActive: true, BirthDate: adultBirthDate()}
	item := &domain.Item{ID: 1, Status: domain.ItemStatusAvailable, LoanValue: 1000}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	loanRepo.On("GenerateNumber", ctx).Return("", errors.New("db error"))

	input := CreateLoanInput{CustomerID: 1, ItemID: 1, LoanAmount: 500, PaymentPlanType: "single"}

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to generate loan number", err.Error())
}

func TestLoanService_Create_BeginTxError(t *testing.T) {
	service, loanRepo, itemRepo, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	customer := &domain.Customer{ID: 1, IsActive: true, BirthDate: adultBirthDate()}
	item := &domain.Item{ID: 1, Status: domain.ItemStatusAvailable, LoanValue: 1000}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	loanRepo.On("GenerateNumber", ctx).Return("LN-000001", nil)
	loanRepo.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	input := CreateLoanInput{CustomerID: 1, ItemID: 1, LoanAmount: 500, InterestRate: 10, LoanTermDays: 30, PaymentPlanType: "single"}

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to start transaction", err.Error())
}

func TestLoanService_Create_WithMinimumPayment(t *testing.T) {
	service, loanRepo, itemRepo, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	customer := &domain.Customer{ID: 1, IsActive: true, BirthDate: adultBirthDate()}
	item := &domain.Item{ID: 1, Status: domain.ItemStatusAvailable, LoanValue: 1000}
	tx := new(mocks.MockTransaction)

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	loanRepo.On("GenerateNumber", ctx).Return("LN-000003", nil)
	loanRepo.On("BeginTx", ctx).Return(tx, nil)
	loanRepo.On("CreateTx", ctx, tx, mock.AnythingOfType("*domain.Loan")).Return(nil)
	itemRepo.On("UpdateStatus", ctx, int64(1), domain.ItemStatusCollateral).Return(nil)
	tx.On("Commit").Return(nil)
	tx.On("Rollback").Return(nil)
	customerRepo.On("UpdateCreditInfo", ctx, int64(1), mock.AnythingOfType("repository.CustomerCreditUpdate")).Return(nil)

	input := CreateLoanInput{
		CustomerID:             1,
		ItemID:                 1,
		BranchID:               1,
		LoanAmount:             500,
		InterestRate:           10,
		LoanTermDays:           60,
		PaymentPlanType:        "minimum_payment",
		RequiresMinimumPayment: true,
		MinimumPaymentAmount:   100,
		CreatedBy:              1,
	}

	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.MinimumPaymentAmount)
	assert.NotNil(t, result.NextPaymentDueDate)
}

// --- GetByID tests ---

func TestLoanService_GetByID_Success(t *testing.T) {
	service, loanRepo, itemRepo, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:         1,
		LoanNumber: "LN-000001",
		Status:     domain.LoanStatusActive,
		CustomerID: 10,
		ItemID:     20,
	}

	customer := &domain.Customer{ID: 10, FirstName: "John", LastName: "Doe"}
	item := &domain.Item{ID: 20, Name: "iPhone 15"}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)
	customerRepo.On("GetByID", ctx, int64(10)).Return(customer, nil)
	itemRepo.On("GetByID", ctx, int64(20)).Return(item, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LN-000001", result.LoanNumber)
}

func TestLoanService_GetByID_NotFound(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	loanRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// --- GetByNumber tests ---

func TestLoanService_GetByNumber_Success(t *testing.T) {
	service, loanRepo, itemRepo, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:         1,
		LoanNumber: "LN-000001",
		CustomerID: 10,
		ItemID:     20,
	}

	customer := &domain.Customer{ID: 10, FirstName: "John", LastName: "Doe"}
	item := &domain.Item{ID: 20, Name: "iPhone 15"}

	loanRepo.On("GetByNumber", ctx, "LN-000001").Return(loan, nil)
	customerRepo.On("GetByID", ctx, int64(10)).Return(customer, nil)
	itemRepo.On("GetByID", ctx, int64(20)).Return(item, nil)

	result, err := service.GetByNumber(ctx, "LN-000001")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LN-000001", result.LoanNumber)
}

func TestLoanService_GetByNumber_NotFound(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	loanRepo.On("GetByNumber", ctx, "NONEXISTENT").Return(nil, errors.New("not found"))

	result, err := service.GetByNumber(ctx, "NONEXISTENT")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// --- List tests ---

func TestLoanService_List_Success(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	loans := []domain.Loan{
		{ID: 1, LoanNumber: "LN-000001"},
		{ID: 2, LoanNumber: "LN-000002"},
	}

	paginatedResult := &repository.PaginatedResult[domain.Loan]{
		Data:       loans,
		Total:      2,
		Page:       1,
		PerPage:    10,
		TotalPages: 1,
	}

	loanRepo.On("List", ctx, mock.AnythingOfType("repository.LoanListParams")).Return(paginatedResult, nil)

	params := repository.LoanListParams{BranchID: 1}
	result, err := service.List(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
}

// --- GetPayments tests ---

func TestLoanService_GetPayments_Success(t *testing.T) {
	service, _, _, _, paymentRepo := setupLoanService()
	ctx := context.Background()

	payments := []*domain.Payment{
		{ID: 1, Amount: 100},
		{ID: 2, Amount: 200},
	}

	paymentRepo.On("ListByLoan", ctx, int64(1)).Return(payments, nil)

	result, err := service.GetPayments(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	paymentRepo.AssertExpectations(t)
}

func TestLoanService_GetPayments_Error(t *testing.T) {
	service, _, _, _, paymentRepo := setupLoanService()
	ctx := context.Background()

	paymentRepo.On("ListByLoan", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetPayments(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// --- GetInstallments tests ---

func TestLoanService_GetInstallments_Success(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	installments := []*domain.LoanInstallment{
		{ID: 1, LoanID: 1, InstallmentNumber: 1},
		{ID: 2, LoanID: 1, InstallmentNumber: 2},
	}

	loanRepo.On("GetInstallments", ctx, int64(1)).Return(installments, nil)

	result, err := service.GetInstallments(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// --- Renew tests ---

func TestLoanService_Renew_Success(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:                 1,
		LoanNumber:         "LN-000001",
		BranchID:           1,
		CustomerID:         1,
		ItemID:             1,
		LoanAmount:         1000,
		InterestRate:       10,
		InterestAmount:     100,
		PrincipalRemaining: 500,
		InterestRemaining:  50,
		Status:             domain.LoanStatusActive,
		PaymentPlanType:    domain.PaymentPlanType("single"),
		GracePeriodDays:    5,
		LateFeeRate:        2,
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)
	loanRepo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)
	loanRepo.On("GenerateNumber", ctx).Return("LN-000002", nil)
	loanRepo.On("Create", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)

	input := RenewLoanInput{
		LoanID:      1,
		NewTermDays: 30,
		UpdatedBy:   1,
	}

	result, err := service.Renew(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LN-000002", result.LoanNumber)
	assert.Equal(t, domain.LoanStatusActive, result.Status)
	assert.Equal(t, 500.0, result.LoanAmount)               // Remaining principal
	assert.Equal(t, loan.InterestRate, result.InterestRate)  // Same rate since NewInterestRate=0
	assert.Equal(t, int64(1), *result.RenewedFromID)
	assert.Equal(t, 1, result.RenewalCount)
	loanRepo.AssertExpectations(t)
}

func TestLoanService_Renew_WithNewInterestRate(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:                 1,
		BranchID:           1,
		CustomerID:         1,
		ItemID:             1,
		InterestRate:       10,
		PrincipalRemaining: 500,
		Status:             domain.LoanStatusOverdue,
		PaymentPlanType:    domain.PaymentPlanType("single"),
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)
	loanRepo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)
	loanRepo.On("GenerateNumber", ctx).Return("LN-000003", nil)
	loanRepo.On("Create", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)

	input := RenewLoanInput{
		LoanID:          1,
		NewTermDays:     60,
		NewInterestRate: 15,
		UpdatedBy:       1,
	}

	result, err := service.Renew(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 15.0, result.InterestRate)
	assert.Equal(t, 75.0, result.InterestAmount) // 500 * 15/100
}

func TestLoanService_Renew_NotFound(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	loanRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := RenewLoanInput{LoanID: 999}

	result, err := service.Renew(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "loan not found", err.Error())
}

func TestLoanService_Renew_InvalidStatus(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:     1,
		Status: domain.LoanStatusPaid,
	}
	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)

	input := RenewLoanInput{LoanID: 1}

	result, err := service.Renew(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "only active or overdue loans can be renewed", err.Error())
}

func TestLoanService_Renew_InterestNotPaid(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:                1,
		Status:            domain.LoanStatusActive,
		InterestRemaining: 50,
	}
	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)

	input := RenewLoanInput{LoanID: 1, PayInterest: true}

	result, err := service.Renew(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "interest must be paid before renewal", err.Error())
}

// --- Confiscate tests ---

func TestLoanService_Confiscate_Success(t *testing.T) {
	service, loanRepo, itemRepo, customerRepo, _ := setupLoanService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:                 1,
		ItemID:             10,
		CustomerID:         20,
		Status:             domain.LoanStatusOverdue,
		PrincipalRemaining: 500.00,
		InterestRemaining:  50.00,
		LateFeeAmount:      10.00,
	}

	customer := &domain.Customer{
		ID:             20,
		FirstName:      "John",
		LastName:       "Doe",
		TotalDefaulted: 0,
	}

	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)
	loanRepo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)
	itemRepo.On("UpdateStatus", ctx, int64(10), domain.ItemStatusConfiscated).Return(nil)
	customerRepo.On("GetByID", ctx, int64(20)).Return(customer, nil)
	customerRepo.On("UpdateCreditInfo", ctx, int64(20), mock.AnythingOfType("repository.CustomerCreditUpdate")).Return(nil)

	err := service.Confiscate(ctx, 1, 1, "Loan overdue")

	assert.NoError(t, err)
}

func TestLoanService_Confiscate_InvalidStatus(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	loan := &domain.Loan{
		ID:     1,
		Status: domain.LoanStatusActive,
	}
	loanRepo.On("GetByID", ctx, int64(1)).Return(loan, nil)

	err := service.Confiscate(ctx, 1, 1, "test")

	assert.Error(t, err)
	assert.Equal(t, "only defaulted or overdue loans can be confiscated", err.Error())
}

func TestLoanService_Confiscate_NotFound(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	loanRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := service.Confiscate(ctx, 999, 1, "test")

	assert.Error(t, err)
	assert.Equal(t, "loan not found", err.Error())
}

// --- GetOverdueLoans tests ---

func TestLoanService_GetOverdueLoans_Success(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	dueDate := time.Now().Add(-24 * time.Hour)
	loans := []*domain.Loan{
		{ID: 1, LoanNumber: "LN-000001", DueDate: dueDate, Status: domain.LoanStatusActive},
	}

	loanRepo.On("GetOverdueLoans", ctx, int64(1)).Return(loans, nil)

	result, err := service.GetOverdueLoans(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// --- UpdateOverdueStatus tests ---

func TestLoanService_UpdateOverdueStatus_Success(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	// Loan past due but in grace period
	graceLoan := &domain.Loan{
		ID:              1,
		Status:          domain.LoanStatusActive,
		DueDate:         time.Now().Add(-2 * 24 * time.Hour),
		GracePeriodDays: 7,
	}

	// Loan past due and past grace period
	defaultedLoan := &domain.Loan{
		ID:              2,
		Status:          domain.LoanStatusOverdue,
		DueDate:         time.Now().Add(-30 * 24 * time.Hour),
		GracePeriodDays: 7,
	}

	loans := []*domain.Loan{graceLoan, defaultedLoan}

	loanRepo.On("GetOverdueLoans", ctx, int64(1)).Return(loans, nil)
	loanRepo.On("Update", ctx, mock.AnythingOfType("*domain.Loan")).Return(nil)

	err := service.UpdateOverdueStatus(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, domain.LoanStatusOverdue, graceLoan.Status)
	assert.Equal(t, domain.LoanStatusDefaulted, defaultedLoan.Status)
	loanRepo.AssertExpectations(t)
}

func TestLoanService_UpdateOverdueStatus_NoLoans(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	loanRepo.On("GetOverdueLoans", ctx, int64(1)).Return([]*domain.Loan{}, nil)

	err := service.UpdateOverdueStatus(ctx, 1)

	assert.NoError(t, err)
}

func TestLoanService_UpdateOverdueStatus_Error(t *testing.T) {
	service, loanRepo, _, _, _ := setupLoanService()
	ctx := context.Background()

	loanRepo.On("GetOverdueLoans", ctx, int64(1)).Return(nil, errors.New("db error"))

	err := service.UpdateOverdueStatus(ctx, 1)

	assert.Error(t, err)
}
