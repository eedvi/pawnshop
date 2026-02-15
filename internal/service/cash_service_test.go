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

func setupCashService() (*CashService, *mocks.MockCashRegisterRepository, *mocks.MockCashSessionRepository, *mocks.MockCashMovementRepository, *mocks.MockBranchRepository) {
	registerRepo := new(mocks.MockCashRegisterRepository)
	sessionRepo := new(mocks.MockCashSessionRepository)
	movementRepo := new(mocks.MockCashMovementRepository)
	branchRepo := new(mocks.MockBranchRepository)
	service := NewCashService(registerRepo, sessionRepo, movementRepo, branchRepo)
	return service, registerRepo, sessionRepo, movementRepo, branchRepo
}

// === Cash Register Tests ===

func TestCashService_CreateRegister_Success(t *testing.T) {
	service, registerRepo, _, _, branchRepo := setupCashService()
	ctx := context.Background()

	branchRepo.On("GetByID", ctx, int64(1)).Return(&domain.Branch{ID: 1}, nil)
	registerRepo.On("Create", ctx, mock.AnythingOfType("*domain.CashRegister")).Return(nil)

	input := CreateRegisterInput{BranchID: 1, Name: "Register 1"}
	result, err := service.CreateRegister(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Register 1", result.Name)
	assert.True(t, result.IsActive)
	assert.Equal(t, int64(1), result.BranchID)
	branchRepo.AssertExpectations(t)
	registerRepo.AssertExpectations(t)
}

func TestCashService_CreateRegister_InvalidBranch(t *testing.T) {
	service, _, _, _, branchRepo := setupCashService()
	ctx := context.Background()

	branchRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreateRegisterInput{BranchID: 999, Name: "Register 1"}
	result, err := service.CreateRegister(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid branch", err.Error())
	branchRepo.AssertExpectations(t)
}

func TestCashService_CreateRegister_CreateFails(t *testing.T) {
	service, registerRepo, _, _, branchRepo := setupCashService()
	ctx := context.Background()

	branchRepo.On("GetByID", ctx, int64(1)).Return(&domain.Branch{ID: 1}, nil)
	registerRepo.On("Create", ctx, mock.AnythingOfType("*domain.CashRegister")).Return(errors.New("db error"))

	input := CreateRegisterInput{BranchID: 1, Name: "Register 1"}
	result, err := service.CreateRegister(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to create cash register", err.Error())
}

func TestCashService_GetRegister_Success(t *testing.T) {
	service, registerRepo, _, _, _ := setupCashService()
	ctx := context.Background()

	register := &domain.CashRegister{ID: 1, Name: "Register 1"}
	registerRepo.On("GetByID", ctx, int64(1)).Return(register, nil)

	result, err := service.GetRegister(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Register 1", result.Name)
	registerRepo.AssertExpectations(t)
}

func TestCashService_GetRegister_NotFound(t *testing.T) {
	service, registerRepo, _, _, _ := setupCashService()
	ctx := context.Background()

	registerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetRegister(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "cash register not found", err.Error())
	registerRepo.AssertExpectations(t)
}

func TestCashService_ListRegisters_Success(t *testing.T) {
	service, registerRepo, _, _, _ := setupCashService()
	ctx := context.Background()

	registers := []*domain.CashRegister{
		{ID: 1, Name: "Register 1"},
		{ID: 2, Name: "Register 2"},
	}
	registerRepo.On("List", ctx, int64(1)).Return(registers, nil)

	result, err := service.ListRegisters(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	registerRepo.AssertExpectations(t)
}

func TestCashService_UpdateRegister_Success(t *testing.T) {
	service, registerRepo, _, _, _ := setupCashService()
	ctx := context.Background()

	existing := &domain.CashRegister{ID: 1, Name: "Old Name", IsActive: true}
	registerRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	registerRepo.On("Update", ctx, mock.AnythingOfType("*domain.CashRegister")).Return(nil)

	input := UpdateRegisterInput{Name: "New Name"}
	result, err := service.UpdateRegister(ctx, 1, input)

	assert.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)
	registerRepo.AssertExpectations(t)
}

func TestCashService_UpdateRegister_NotFound(t *testing.T) {
	service, registerRepo, _, _, _ := setupCashService()
	ctx := context.Background()

	registerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := UpdateRegisterInput{Name: "New Name"}
	result, err := service.UpdateRegister(ctx, 999, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "cash register not found", err.Error())
	registerRepo.AssertExpectations(t)
}

func TestCashService_UpdateRegister_SetInactive(t *testing.T) {
	service, registerRepo, _, _, _ := setupCashService()
	ctx := context.Background()

	existing := &domain.CashRegister{ID: 1, Name: "Register 1", IsActive: true}
	registerRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	registerRepo.On("Update", ctx, mock.AnythingOfType("*domain.CashRegister")).Return(nil)

	isActive := false
	input := UpdateRegisterInput{IsActive: &isActive}
	result, err := service.UpdateRegister(ctx, 1, input)

	assert.NoError(t, err)
	assert.False(t, result.IsActive)
	registerRepo.AssertExpectations(t)
}

func TestCashService_UpdateRegister_WithDescription(t *testing.T) {
	service, registerRepo, _, _, _ := setupCashService()
	ctx := context.Background()

	existing := &domain.CashRegister{ID: 1, Name: "Register 1"}
	registerRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	registerRepo.On("Update", ctx, mock.AnythingOfType("*domain.CashRegister")).Return(nil)

	desc := "Main register"
	input := UpdateRegisterInput{Description: &desc}
	result, err := service.UpdateRegister(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result.Description)
	assert.Equal(t, "Main register", *result.Description)
}

// === Cash Session Tests ===

func TestCashService_OpenSession_Success(t *testing.T) {
	service, registerRepo, sessionRepo, _, _ := setupCashService()
	ctx := context.Background()

	register := &domain.CashRegister{ID: 1, BranchID: 1, IsActive: true}
	registerRepo.On("GetByID", ctx, int64(1)).Return(register, nil)
	sessionRepo.On("GetOpenSession", ctx, int64(10)).Return(nil, errors.New("none"))
	sessionRepo.On("GetOpenSessionByRegister", ctx, int64(1)).Return(nil, errors.New("none"))
	sessionRepo.On("Create", ctx, mock.AnythingOfType("*domain.CashSession")).Return(nil)

	input := OpenSessionInput{BranchID: 1, RegisterID: 1, UserID: 10, OpeningAmount: 1000.0}
	result, err := service.OpenSession(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.CashSessionStatusOpen, result.Status)
	assert.Equal(t, 1000.0, result.OpeningAmount)
	registerRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
}

func TestCashService_OpenSession_RegisterNotFound(t *testing.T) {
	service, registerRepo, _, _, _ := setupCashService()
	ctx := context.Background()

	registerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := OpenSessionInput{BranchID: 1, RegisterID: 999, UserID: 10}
	result, err := service.OpenSession(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "cash register not found", err.Error())
}

func TestCashService_OpenSession_InactiveRegister(t *testing.T) {
	service, registerRepo, _, _, _ := setupCashService()
	ctx := context.Background()

	register := &domain.CashRegister{ID: 1, BranchID: 1, IsActive: false}
	registerRepo.On("GetByID", ctx, int64(1)).Return(register, nil)

	input := OpenSessionInput{BranchID: 1, RegisterID: 1, UserID: 10}
	result, err := service.OpenSession(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "cash register is not active", err.Error())
}

func TestCashService_OpenSession_WrongBranch(t *testing.T) {
	service, registerRepo, _, _, _ := setupCashService()
	ctx := context.Background()

	register := &domain.CashRegister{ID: 1, BranchID: 2, IsActive: true}
	registerRepo.On("GetByID", ctx, int64(1)).Return(register, nil)

	input := OpenSessionInput{BranchID: 1, RegisterID: 1, UserID: 10}
	result, err := service.OpenSession(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "cash register does not belong to this branch", err.Error())
}

func TestCashService_OpenSession_UserAlreadyHasSession(t *testing.T) {
	service, registerRepo, sessionRepo, _, _ := setupCashService()
	ctx := context.Background()

	register := &domain.CashRegister{ID: 1, BranchID: 1, IsActive: true}
	existingSession := &domain.CashSession{ID: 99}

	registerRepo.On("GetByID", ctx, int64(1)).Return(register, nil)
	sessionRepo.On("GetOpenSession", ctx, int64(10)).Return(existingSession, nil)

	input := OpenSessionInput{BranchID: 1, RegisterID: 1, UserID: 10}
	result, err := service.OpenSession(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user already has an open cash session", err.Error())
}

func TestCashService_OpenSession_RegisterAlreadyHasSession(t *testing.T) {
	service, registerRepo, sessionRepo, _, _ := setupCashService()
	ctx := context.Background()

	register := &domain.CashRegister{ID: 1, BranchID: 1, IsActive: true}
	existingSession := &domain.CashSession{ID: 99}

	registerRepo.On("GetByID", ctx, int64(1)).Return(register, nil)
	sessionRepo.On("GetOpenSession", ctx, int64(10)).Return(nil, errors.New("none"))
	sessionRepo.On("GetOpenSessionByRegister", ctx, int64(1)).Return(existingSession, nil)

	input := OpenSessionInput{BranchID: 1, RegisterID: 1, UserID: 10}
	result, err := service.OpenSession(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "register already has an open session", err.Error())
}

func TestCashService_GetSession_Success(t *testing.T) {
	service, _, sessionRepo, movementRepo, _ := setupCashService()
	ctx := context.Background()

	session := &domain.CashSession{ID: 1, BranchID: 1}
	movements := []*domain.CashMovement{{ID: 1}, {ID: 2}}

	sessionRepo.On("GetByID", ctx, int64(1)).Return(session, nil)
	movementRepo.On("ListBySession", ctx, int64(1)).Return(movements, nil)

	result, err := service.GetSession(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Movements, 2)
	sessionRepo.AssertExpectations(t)
	movementRepo.AssertExpectations(t)
}

func TestCashService_GetSession_NotFound(t *testing.T) {
	service, _, sessionRepo, _, _ := setupCashService()
	ctx := context.Background()

	sessionRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetSession(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "cash session not found", err.Error())
}

func TestCashService_GetCurrentSession_Success(t *testing.T) {
	service, _, sessionRepo, movementRepo, _ := setupCashService()
	ctx := context.Background()

	session := &domain.CashSession{ID: 1}
	movements := []*domain.CashMovement{{ID: 1}}

	sessionRepo.On("GetOpenSession", ctx, int64(10)).Return(session, nil)
	movementRepo.On("ListBySession", ctx, int64(1)).Return(movements, nil)

	result, err := service.GetCurrentSession(ctx, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Movements, 1)
}

func TestCashService_GetCurrentSession_NoOpenSession(t *testing.T) {
	service, _, sessionRepo, _, _ := setupCashService()
	ctx := context.Background()

	sessionRepo.On("GetOpenSession", ctx, int64(10)).Return(nil, errors.New("none"))

	result, err := service.GetCurrentSession(ctx, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "no open cash session found", err.Error())
}

func TestCashService_ListSessions_Success(t *testing.T) {
	service, _, sessionRepo, _, _ := setupCashService()
	ctx := context.Background()

	params := repository.CashSessionListParams{
		BranchID:         1,
		PaginationParams: repository.PaginationParams{Page: 1, PerPage: 10},
	}

	result := &repository.PaginatedResult[domain.CashSession]{
		Data:  []domain.CashSession{{ID: 1}},
		Total: 1,
	}

	sessionRepo.On("List", ctx, params).Return(result, nil)

	res, err := service.ListSessions(ctx, params)

	assert.NoError(t, err)
	assert.Len(t, res.Data, 1)
	sessionRepo.AssertExpectations(t)
}

// === Cash Movement Tests ===

func TestCashService_GetMovement_Success(t *testing.T) {
	service, _, _, movementRepo, _ := setupCashService()
	ctx := context.Background()

	movement := &domain.CashMovement{ID: 1, Amount: 100.0, Description: "Payment"}
	movementRepo.On("GetByID", ctx, int64(1)).Return(movement, nil)

	result, err := service.GetMovement(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 100.0, result.Amount)
	movementRepo.AssertExpectations(t)
}

func TestCashService_GetMovement_NotFound(t *testing.T) {
	service, _, _, movementRepo, _ := setupCashService()
	ctx := context.Background()

	movementRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetMovement(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "cash movement not found", err.Error())
}

func TestCashService_ListMovements_Success(t *testing.T) {
	service, _, _, movementRepo, _ := setupCashService()
	ctx := context.Background()

	params := repository.CashMovementListParams{
		BranchID:         1,
		PaginationParams: repository.PaginationParams{Page: 1, PerPage: 10},
	}

	result := &repository.PaginatedResult[domain.CashMovement]{
		Data:  []domain.CashMovement{{ID: 1}, {ID: 2}},
		Total: 2,
	}

	movementRepo.On("List", ctx, params).Return(result, nil)

	res, err := service.ListMovements(ctx, params)

	assert.NoError(t, err)
	assert.Len(t, res.Data, 2)
	movementRepo.AssertExpectations(t)
}

func TestCashService_ListSessionMovements_Success(t *testing.T) {
	service, _, _, movementRepo, _ := setupCashService()
	ctx := context.Background()

	movements := []*domain.CashMovement{{ID: 1}, {ID: 2}, {ID: 3}}
	movementRepo.On("ListBySession", ctx, int64(1)).Return(movements, nil)

	result, err := service.ListSessionMovements(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	movementRepo.AssertExpectations(t)
}
