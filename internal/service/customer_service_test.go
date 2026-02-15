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

func setupCustomerService() (*CustomerService, *mocks.MockCustomerRepository, *mocks.MockBranchRepository) {
	customerRepo := new(mocks.MockCustomerRepository)
	branchRepo := new(mocks.MockBranchRepository)
	service := NewCustomerService(customerRepo, branchRepo)
	return service, customerRepo, branchRepo
}

func TestCustomerService_Create_Success(t *testing.T) {
	service, customerRepo, branchRepo := setupCustomerService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main", IsActive: true}

	// Setup expectations
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	customerRepo.On("GetByIdentity", ctx, int64(1), "DPI", "1234567890").Return(nil, errors.New("not found"))
	customerRepo.On("Create", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)

	// Execute
	input := CreateCustomerInput{
		BranchID:       1,
		FirstName:      "John",
		LastName:       "Doe",
		IdentityType:   "DPI",
		IdentityNumber: "1234567890",
		Phone:          "555-1234",
		Email:          "john@example.com",
	}
	result, err := service.Create(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	assert.True(t, result.IsActive)
	assert.False(t, result.IsBlocked)

	customerRepo.AssertExpectations(t)
	branchRepo.AssertExpectations(t)
}

func TestCustomerService_Create_DuplicateIdentity(t *testing.T) {
	service, customerRepo, branchRepo := setupCustomerService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main", IsActive: true}
	existingCustomer := &domain.Customer{
		ID:             1,
		IdentityNumber: "1234567890",
	}

	// Setup expectations
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	customerRepo.On("GetByIdentity", ctx, int64(1), "DPI", "1234567890").Return(existingCustomer, nil)

	// Execute
	input := CreateCustomerInput{
		BranchID:       1,
		FirstName:      "John",
		LastName:       "Doe",
		IdentityType:   "DPI",
		IdentityNumber: "1234567890",
	}
	result, err := service.Create(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "customer with this identity already exists", err.Error())

	customerRepo.AssertExpectations(t)
	branchRepo.AssertExpectations(t)
}

func TestCustomerService_Update_Success(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	existingCustomer := &domain.Customer{
		ID:             1,
		BranchID:       1,
		FirstName:      "John",
		LastName:       "Doe",
		IdentityType:   "DPI",
		IdentityNumber: "1234567890",
		IsActive:       true,
	}

	// Setup expectations
	customerRepo.On("GetByID", ctx, int64(1)).Return(existingCustomer, nil)
	customerRepo.On("Update", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)

	// Execute
	input := UpdateCustomerInput{
		FirstName: "Jane",
		LastName:  "Smith",
		Phone:     "555-5678",
	}
	result, err := service.Update(ctx, 1, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Jane", result.FirstName)
	assert.Equal(t, "Smith", result.LastName)

	customerRepo.AssertExpectations(t)
}

func TestCustomerService_Update_NotFound(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	// Setup expectations
	customerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	// Execute
	input := UpdateCustomerInput{
		FirstName: "Jane",
	}
	result, err := service.Update(ctx, 999, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "customer not found", err.Error())

	customerRepo.AssertExpectations(t)
}

func TestCustomerService_GetByID_Success(t *testing.T) {
	service, customerRepo, branchRepo := setupCustomerService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:        1,
		BranchID:  1,
		FirstName: "John",
		LastName:  "Doe",
	}

	branch := &domain.Branch{
		ID:   1,
		Name: "Main Branch",
	}

	// Setup expectations
	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)

	// Execute
	result, err := service.GetByID(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "John", result.FirstName)
	assert.NotNil(t, result.Branch)

	customerRepo.AssertExpectations(t)
	branchRepo.AssertExpectations(t)
}

func TestCustomerService_List_Success(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	customers := []domain.Customer{
		{ID: 1, FirstName: "John", LastName: "Doe"},
		{ID: 2, FirstName: "Jane", LastName: "Smith"},
	}

	paginatedResult := &repository.PaginatedResult[domain.Customer]{
		Data:       customers,
		Total:      2,
		Page:       1,
		PerPage:    10,
		TotalPages: 1,
	}

	// Setup expectations
	customerRepo.On("List", ctx, mock.AnythingOfType("repository.CustomerListParams")).Return(paginatedResult, nil)

	// Execute
	params := repository.CustomerListParams{
		BranchID: 1,
		PaginationParams: repository.PaginationParams{
			Page:    1,
			PerPage: 10,
		},
	}
	result, err := service.List(ctx, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, 2, result.Total)

	customerRepo.AssertExpectations(t)
}

func TestCustomerService_Delete_Success(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:        1,
		FirstName: "John",
	}

	// Setup expectations
	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	customerRepo.On("Delete", ctx, int64(1)).Return(nil)

	// Execute
	err := service.Delete(ctx, 1)

	// Assert
	assert.NoError(t, err)
	customerRepo.AssertExpectations(t)
}

func TestCustomerService_Block_Success(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:        1,
		FirstName: "John",
		IsBlocked: false,
	}

	// Setup expectations
	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	customerRepo.On("Update", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)

	// Execute
	input := BlockCustomerInput{
		CustomerID: 1,
		Reason:     "Too many defaults",
	}
	err := service.Block(ctx, input)

	// Assert
	assert.NoError(t, err)
	customerRepo.AssertExpectations(t)
}

func TestCustomerService_Unblock_Success(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	customer := &domain.Customer{
		ID:        1,
		FirstName: "John",
		IsBlocked: true,
	}

	// Setup expectations
	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	customerRepo.On("Update", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)

	// Execute
	err := service.Unblock(ctx, 1)

	// Assert
	assert.NoError(t, err)
	customerRepo.AssertExpectations(t)
}

// --- Missing coverage tests ---

func TestCustomerService_Create_InvalidBranch(t *testing.T) {
	service, _, branchRepo := setupCustomerService()
	ctx := context.Background()

	branchRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreateCustomerInput{
		BranchID:       999,
		FirstName:      "John",
		LastName:       "Doe",
		IdentityType:   "DPI",
		IdentityNumber: "1234567890",
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid branch", err.Error())
	branchRepo.AssertExpectations(t)
}

func TestCustomerService_Create_UnderageBirthDate(t *testing.T) {
	service, customerRepo, branchRepo := setupCustomerService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main"}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	customerRepo.On("GetByIdentity", ctx, int64(1), "DPI", "1234567890").Return(nil, errors.New("not found"))

	birthDate := time.Now().AddDate(-16, 0, 0)
	input := CreateCustomerInput{
		BranchID:       1,
		FirstName:      "Young",
		LastName:       "Person",
		IdentityType:   "DPI",
		IdentityNumber: "1234567890",
		Phone:          "555-1234",
		BirthDate:      &birthDate,
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "customer must be at least 18 years old", err.Error())
}

func TestCustomerService_Create_AdultBirthDate(t *testing.T) {
	service, customerRepo, branchRepo := setupCustomerService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main"}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	customerRepo.On("GetByIdentity", ctx, int64(1), "DPI", "9999999999").Return(nil, errors.New("not found"))
	customerRepo.On("Create", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)

	birthDate := time.Now().AddDate(-25, 0, 0)
	input := CreateCustomerInput{
		BranchID:       1,
		FirstName:      "Adult",
		LastName:       "Person",
		IdentityType:   "DPI",
		IdentityNumber: "9999999999",
		Phone:          "555-1234",
		BirthDate:      &birthDate,
	}
	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, &birthDate, result.BirthDate)
}

func TestCustomerService_Create_RepoError(t *testing.T) {
	service, customerRepo, branchRepo := setupCustomerService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main"}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	customerRepo.On("GetByIdentity", ctx, int64(1), "DPI", "5555555555").Return(nil, errors.New("not found"))
	customerRepo.On("Create", ctx, mock.AnythingOfType("*domain.Customer")).Return(errors.New("db error"))

	input := CreateCustomerInput{
		BranchID:       1,
		FirstName:      "John",
		LastName:       "Doe",
		IdentityType:   "DPI",
		IdentityNumber: "5555555555",
		Phone:          "555-1234",
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to create customer", err.Error())
}

func TestCustomerService_Update_WithBirthDateUnderage(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	existing := &domain.Customer{ID: 1, FirstName: "John", IsActive: true}
	customerRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)

	birthDate := time.Now().AddDate(-15, 0, 0)
	input := UpdateCustomerInput{
		BirthDate: &birthDate,
	}
	result, err := service.Update(ctx, 1, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "customer must be at least 18 years old", err.Error())
}

func TestCustomerService_Update_AllFields(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	existing := &domain.Customer{ID: 1, FirstName: "John", IsActive: true}
	customerRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	customerRepo.On("Update", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)

	birthDate := time.Now().AddDate(-30, 0, 0)
	income := 5000.0
	creditLimit := 10000.0
	isActive := false
	input := UpdateCustomerInput{
		FirstName:     "Updated",
		LastName:      "Name",
		BirthDate:     &birthDate,
		Gender:        "male",
		Phone:         "555-9999",
		MonthlyIncome: &income,
		CreditLimit:   &creditLimit,
		IsActive:      &isActive,
	}
	result, err := service.Update(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated", result.FirstName)
	assert.Equal(t, "Name", result.LastName)
	assert.Equal(t, "male", result.Gender)
	assert.Equal(t, 5000.0, result.MonthlyIncome)
	assert.Equal(t, 10000.0, result.CreditLimit)
	assert.False(t, result.IsActive)
}

func TestCustomerService_Update_RepoError(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	existing := &domain.Customer{ID: 1, FirstName: "John", IsActive: true}
	customerRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	customerRepo.On("Update", ctx, mock.AnythingOfType("*domain.Customer")).Return(errors.New("db error"))

	input := UpdateCustomerInput{FirstName: "Jane"}
	result, err := service.Update(ctx, 1, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to update customer", err.Error())
}

func TestCustomerService_GetByID_NotFound(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	customerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "customer not found", err.Error())
}

func TestCustomerService_Delete_NotFound(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	customerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := service.Delete(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, "customer not found", err.Error())
}

func TestCustomerService_Block_NotFound(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	customerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := BlockCustomerInput{CustomerID: 999, Reason: "Bad behavior"}
	err := service.Block(ctx, input)

	assert.Error(t, err)
	assert.Equal(t, "customer not found", err.Error())
}

func TestCustomerService_Unblock_NotFound(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	customerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := service.Unblock(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, "customer not found", err.Error())
}

func TestCustomerService_UpdateCreditScore_Success(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	score := 75
	customerRepo.On("UpdateCreditInfo", ctx, int64(1), repository.CustomerCreditUpdate{
		CreditScore: &score,
	}).Return(nil)

	err := service.UpdateCreditScore(ctx, 1, 75)

	assert.NoError(t, err)
	customerRepo.AssertExpectations(t)
}

func TestCustomerService_UpdateCreditScore_TooLow(t *testing.T) {
	service, _, _ := setupCustomerService()
	ctx := context.Background()

	err := service.UpdateCreditScore(ctx, 1, -1)

	assert.Error(t, err)
	assert.Equal(t, "credit score must be between 0 and 100", err.Error())
}

func TestCustomerService_UpdateCreditScore_TooHigh(t *testing.T) {
	service, _, _ := setupCustomerService()
	ctx := context.Background()

	err := service.UpdateCreditScore(ctx, 1, 101)

	assert.Error(t, err)
	assert.Equal(t, "credit score must be between 0 and 100", err.Error())
}

func TestCustomerService_UpdateCreditScore_BoundaryValues(t *testing.T) {
	service, customerRepo, _ := setupCustomerService()
	ctx := context.Background()

	score0 := 0
	customerRepo.On("UpdateCreditInfo", ctx, int64(1), repository.CustomerCreditUpdate{
		CreditScore: &score0,
	}).Return(nil)
	err := service.UpdateCreditScore(ctx, 1, 0)
	assert.NoError(t, err)

	score100 := 100
	customerRepo.On("UpdateCreditInfo", ctx, int64(2), repository.CustomerCreditUpdate{
		CreditScore: &score100,
	}).Return(nil)
	err = service.UpdateCreditScore(ctx, 2, 100)
	assert.NoError(t, err)
}

func TestCalculateAge(t *testing.T) {
	birthDate := time.Now().AddDate(-30, 0, 0)
	age := calculateAge(birthDate)
	assert.Equal(t, 30, age)

	birthDate = time.Now().AddDate(-10, 0, 0)
	age = calculateAge(birthDate)
	assert.Equal(t, 10, age)

	birthDate = time.Now().AddDate(-20, 0, 1) // birthday is tomorrow
	age = calculateAge(birthDate)
	assert.Equal(t, 19, age)
}
