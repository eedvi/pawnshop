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

func setupBranchService() (*BranchService, *mocks.MockBranchRepository) {
	branchRepo := new(mocks.MockBranchRepository)
	service := NewBranchService(branchRepo)
	return service, branchRepo
}

func TestBranchService_Create_Success(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	// Setup expectations
	branchRepo.On("GetByCode", ctx, "TEST01").Return(nil, errors.New("not found"))
	branchRepo.On("Create", ctx, mock.AnythingOfType("*domain.Branch")).Return(nil)

	// Execute
	input := CreateBranchInput{
		Name:    "Test Branch",
		Code:    "TEST01",
		Address: "123 Test St",
		Phone:   "555-1234",
		Email:   "test@example.com",
	}
	result, err := service.Create(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Branch", result.Name)
	assert.Equal(t, "TEST01", result.Code)
	assert.True(t, result.IsActive)

	branchRepo.AssertExpectations(t)
}

func TestBranchService_Create_DuplicateCode(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	existingBranch := &domain.Branch{ID: 1, Code: "TEST01"}

	// Setup expectations
	branchRepo.On("GetByCode", ctx, "TEST01").Return(existingBranch, nil)

	// Execute
	input := CreateBranchInput{
		Name: "Test Branch",
		Code: "TEST01",
	}
	result, err := service.Create(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "branch with this code already exists", err.Error())

	branchRepo.AssertExpectations(t)
}

func TestBranchService_Update_Success(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	existingBranch := &domain.Branch{
		ID:       1,
		Name:     "Old Name",
		Code:     "OLD01",
		IsActive: true,
	}

	// Setup expectations
	branchRepo.On("GetByID", ctx, int64(1)).Return(existingBranch, nil)
	branchRepo.On("Update", ctx, mock.AnythingOfType("*domain.Branch")).Return(nil)

	// Execute
	input := UpdateBranchInput{
		Name:    "New Name",
		Address: "New Address",
	}
	result, err := service.Update(ctx, 1, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.Name)

	branchRepo.AssertExpectations(t)
}

func TestBranchService_Update_NotFound(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	// Setup expectations
	branchRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	// Execute
	input := UpdateBranchInput{
		Name: "New Name",
	}
	result, err := service.Update(ctx, 999, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "branch not found", err.Error())

	branchRepo.AssertExpectations(t)
}

func TestBranchService_GetByID_Success(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branch := &domain.Branch{
		ID:   1,
		Name: "Test Branch",
		Code: "TEST01",
	}

	// Setup expectations
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)

	// Execute
	result, err := service.GetByID(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Branch", result.Name)

	branchRepo.AssertExpectations(t)
}

func TestBranchService_List_Success(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branches := []domain.Branch{
		{ID: 1, Name: "Branch 1", Code: "B001"},
		{ID: 2, Name: "Branch 2", Code: "B002"},
	}

	paginatedResult := &repository.PaginatedResult[domain.Branch]{
		Data:       branches,
		Total:      2,
		Page:       1,
		PerPage:    10,
		TotalPages: 1,
	}

	// Setup expectations
	branchRepo.On("List", ctx, mock.AnythingOfType("repository.PaginationParams")).Return(paginatedResult, nil)

	// Execute
	params := repository.PaginationParams{
		Page:    1,
		PerPage: 10,
	}
	result, err := service.List(ctx, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)

	branchRepo.AssertExpectations(t)
}

func TestBranchService_Delete_Success(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Test Branch"}

	// Setup expectations
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	branchRepo.On("Delete", ctx, int64(1)).Return(nil)

	// Execute
	err := service.Delete(ctx, 1)

	// Assert
	assert.NoError(t, err)
	branchRepo.AssertExpectations(t)
}

func TestBranchService_Activate_Success(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branch := &domain.Branch{
		ID:       1,
		Name:     "Test Branch",
		IsActive: false,
	}

	// Setup expectations
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	branchRepo.On("Update", ctx, mock.AnythingOfType("*domain.Branch")).Return(nil)

	// Execute
	err := service.Activate(ctx, 1)

	// Assert
	assert.NoError(t, err)
	branchRepo.AssertExpectations(t)
}

func TestBranchService_Deactivate_Success(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branch := &domain.Branch{
		ID:       1,
		Name:     "Test Branch",
		IsActive: true,
	}

	// Setup expectations
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	branchRepo.On("Update", ctx, mock.AnythingOfType("*domain.Branch")).Return(nil)

	// Execute
	err := service.Deactivate(ctx, 1)

	// Assert
	assert.NoError(t, err)
	branchRepo.AssertExpectations(t)
}

// --- Missing coverage tests ---

func TestBranchService_GetByCode_Success(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Test Branch", Code: "TEST01"}
	branchRepo.On("GetByCode", ctx, "TEST01").Return(branch, nil)

	result, err := service.GetByCode(ctx, "TEST01")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "TEST01", result.Code)
	branchRepo.AssertExpectations(t)
}

func TestBranchService_GetByCode_NotFound(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branchRepo.On("GetByCode", ctx, "NOTEXIST").Return(nil, errors.New("not found"))

	result, err := service.GetByCode(ctx, "NOTEXIST")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "branch not found", err.Error())
	branchRepo.AssertExpectations(t)
}

func TestBranchService_GetByID_NotFound(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branchRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "branch not found", err.Error())
	branchRepo.AssertExpectations(t)
}

func TestBranchService_Delete_NotFound(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branchRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := service.Delete(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, "branch not found", err.Error())
	branchRepo.AssertExpectations(t)
}

func TestBranchService_Activate_NotFound(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branchRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := service.Activate(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, "branch not found", err.Error())
	branchRepo.AssertExpectations(t)
}

func TestBranchService_Deactivate_NotFound(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branchRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := service.Deactivate(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, "branch not found", err.Error())
	branchRepo.AssertExpectations(t)
}

func TestBranchService_Create_WithDefaults(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branchRepo.On("GetByCode", ctx, "DEF01").Return(nil, errors.New("not found"))
	branchRepo.On("Create", ctx, mock.AnythingOfType("*domain.Branch")).Return(nil)

	input := CreateBranchInput{
		Name: "Default Branch",
		Code: "DEF01",
	}
	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "America/Mexico_City", result.Timezone)
	assert.Equal(t, "MXN", result.Currency)
	assert.Equal(t, 30, result.DefaultLoanTermDays)
	branchRepo.AssertExpectations(t)
}

func TestBranchService_Create_RepoError(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	branchRepo.On("GetByCode", ctx, "ERR01").Return(nil, errors.New("not found"))
	branchRepo.On("Create", ctx, mock.AnythingOfType("*domain.Branch")).Return(errors.New("db error"))

	input := CreateBranchInput{
		Name: "Error Branch",
		Code: "ERR01",
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to create branch", err.Error())
	branchRepo.AssertExpectations(t)
}

func TestBranchService_Update_AllFields(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	existingBranch := &domain.Branch{
		ID:                  1,
		Name:                "Old",
		Code:                "OLD01",
		IsActive:            true,
		DefaultInterestRate: 5.0,
		DefaultLoanTermDays: 30,
		DefaultGracePeriod:  3,
	}

	branchRepo.On("GetByID", ctx, int64(1)).Return(existingBranch, nil)
	branchRepo.On("Update", ctx, mock.AnythingOfType("*domain.Branch")).Return(nil)

	isActive := false
	interestRate := 10.0
	loanDays := 60
	gracePeriod := 5
	input := UpdateBranchInput{
		Name:                "New Name",
		Address:             "New Address",
		Phone:               "555-0000",
		Email:               "new@example.com",
		Timezone:            "America/Guatemala",
		Currency:            "GTQ",
		IsActive:            &isActive,
		DefaultInterestRate: &interestRate,
		DefaultLoanTermDays: &loanDays,
		DefaultGracePeriod:  &gracePeriod,
	}
	result, err := service.Update(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, "New Address", result.Address)
	assert.Equal(t, "555-0000", result.Phone)
	assert.Equal(t, "new@example.com", result.Email)
	assert.Equal(t, "America/Guatemala", result.Timezone)
	assert.Equal(t, "GTQ", result.Currency)
	assert.False(t, result.IsActive)
	assert.Equal(t, 10.0, result.DefaultInterestRate)
	assert.Equal(t, 60, result.DefaultLoanTermDays)
	assert.Equal(t, 5, result.DefaultGracePeriod)
	branchRepo.AssertExpectations(t)
}

func TestBranchService_Update_RepoError(t *testing.T) {
	service, branchRepo := setupBranchService()
	ctx := context.Background()

	existing := &domain.Branch{ID: 1, Name: "Test"}
	branchRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	branchRepo.On("Update", ctx, mock.AnythingOfType("*domain.Branch")).Return(errors.New("db error"))

	input := UpdateBranchInput{Name: "Updated"}
	result, err := service.Update(ctx, 1, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to update branch", err.Error())
	branchRepo.AssertExpectations(t)
}
