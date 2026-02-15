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
	"pawnshop/pkg/auth"
)

func setupUserService() (*UserService, *mocks.MockUserRepository, *mocks.MockRoleRepository, *mocks.MockBranchRepository) {
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	branchRepo := new(mocks.MockBranchRepository)
	passwordManager := auth.NewPasswordManager()
	service := NewUserService(userRepo, roleRepo, branchRepo, passwordManager)
	return service, userRepo, roleRepo, branchRepo
}

func TestUserService_Create_Success(t *testing.T) {
	service, userRepo, roleRepo, _ := setupUserService()
	ctx := context.Background()

	role := &domain.Role{ID: 1, Name: "admin"}

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, errors.New("not found"))
	roleRepo.On("GetByID", ctx, int64(1)).Return(role, nil)
	userRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	input := CreateUserInput{
		Email:     "test@example.com",
		Password:  "StrongPass1!",
		FirstName: "John",
		LastName:  "Doe",
		RoleID:    1,
		IsActive:  true,
	}
	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
}

func TestUserService_Create_DuplicateEmail(t *testing.T) {
	service, userRepo, _, _ := setupUserService()
	ctx := context.Background()

	existing := &domain.User{ID: 1, Email: "test@example.com"}
	userRepo.On("GetByEmail", ctx, "test@example.com").Return(existing, nil)

	input := CreateUserInput{Email: "test@example.com", Password: "StrongPass1!", FirstName: "John", LastName: "Doe", RoleID: 1}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "email is already registered", err.Error())
	userRepo.AssertExpectations(t)
}

func TestUserService_Create_InvalidRole(t *testing.T) {
	service, userRepo, roleRepo, _ := setupUserService()
	ctx := context.Background()

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, errors.New("not found"))
	roleRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreateUserInput{Email: "test@example.com", Password: "StrongPass1!", FirstName: "John", LastName: "Doe", RoleID: 999}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid role", err.Error())
	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
}

func TestUserService_Create_InvalidBranch(t *testing.T) {
	service, userRepo, roleRepo, branchRepo := setupUserService()
	ctx := context.Background()

	role := &domain.Role{ID: 1, Name: "admin"}

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, errors.New("not found"))
	roleRepo.On("GetByID", ctx, int64(1)).Return(role, nil)

	branchID := int64(999)
	branchRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreateUserInput{Email: "test@example.com", Password: "StrongPass1!", FirstName: "John", LastName: "Doe", RoleID: 1, BranchID: &branchID}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid branch", err.Error())
	branchRepo.AssertExpectations(t)
}

func TestUserService_Create_WeakPassword(t *testing.T) {
	service, userRepo, roleRepo, _ := setupUserService()
	ctx := context.Background()

	role := &domain.Role{ID: 1, Name: "admin"}

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, errors.New("not found"))
	roleRepo.On("GetByID", ctx, int64(1)).Return(role, nil)

	input := CreateUserInput{Email: "test@example.com", Password: "weak", FirstName: "John", LastName: "Doe", RoleID: 1}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	userRepo.AssertExpectations(t)
}

func TestUserService_Update_Success(t *testing.T) {
	service, userRepo, roleRepo, _ := setupUserService()
	ctx := context.Background()

	existing := &domain.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		RoleID:    1,
		IsActive:  true,
	}

	userRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
	roleRepo.On("GetByID", ctx, int64(1)).Return(&domain.Role{ID: 1, Name: "admin"}, nil)

	input := UpdateUserInput{FirstName: "Jane", LastName: "Smith"}
	result, err := service.Update(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Jane", result.FirstName)
	assert.Equal(t, "Smith", result.LastName)
	userRepo.AssertExpectations(t)
}

func TestUserService_Update_NotFound(t *testing.T) {
	service, userRepo, _, _ := setupUserService()
	ctx := context.Background()

	userRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := UpdateUserInput{FirstName: "Jane"}
	result, err := service.Update(ctx, 999, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user not found", err.Error())
	userRepo.AssertExpectations(t)
}

func TestUserService_Update_ChangeRole(t *testing.T) {
	service, userRepo, roleRepo, _ := setupUserService()
	ctx := context.Background()

	existing := &domain.User{ID: 1, RoleID: 1, FirstName: "John", LastName: "Doe"}
	newRole := &domain.Role{ID: 2, Name: "manager"}

	userRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	roleRepo.On("GetByID", ctx, int64(2)).Return(newRole, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	newRoleID := int64(2)
	input := UpdateUserInput{RoleID: &newRoleID}
	result, err := service.Update(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.RoleID)
	roleRepo.AssertExpectations(t)
}

func TestUserService_Update_InvalidRole(t *testing.T) {
	service, userRepo, roleRepo, _ := setupUserService()
	ctx := context.Background()

	existing := &domain.User{ID: 1, RoleID: 1, FirstName: "John", LastName: "Doe"}

	userRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	roleRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	invalidRole := int64(999)
	input := UpdateUserInput{RoleID: &invalidRole}
	result, err := service.Update(ctx, 1, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid role", err.Error())
	roleRepo.AssertExpectations(t)
}

func TestUserService_Update_ChangeBranch(t *testing.T) {
	service, userRepo, roleRepo, branchRepo := setupUserService()
	ctx := context.Background()

	existing := &domain.User{ID: 1, RoleID: 1, FirstName: "John", LastName: "Doe"}

	userRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	branchRepo.On("GetByID", ctx, int64(2)).Return(&domain.Branch{ID: 2}, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
	roleRepo.On("GetByID", ctx, int64(1)).Return(&domain.Role{ID: 1, Name: "admin"}, nil)

	newBranch := int64(2)
	input := UpdateUserInput{BranchID: &newBranch}
	result, err := service.Update(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	branchRepo.AssertExpectations(t)
}

func TestUserService_GetByID_Success(t *testing.T) {
	service, userRepo, roleRepo, _ := setupUserService()
	ctx := context.Background()

	user := &domain.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe", RoleID: 1}
	role := &domain.Role{ID: 1, Name: "admin"}

	userRepo.On("GetByID", ctx, int64(1)).Return(user, nil)
	roleRepo.On("GetByID", ctx, int64(1)).Return(role, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	userRepo.AssertExpectations(t)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	service, userRepo, _, _ := setupUserService()
	ctx := context.Background()

	userRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user not found", err.Error())
	userRepo.AssertExpectations(t)
}

func TestUserService_List_Success(t *testing.T) {
	service, userRepo, roleRepo, _ := setupUserService()
	ctx := context.Background()

	users := []domain.User{
		{ID: 1, Email: "user1@example.com", FirstName: "User", LastName: "One", RoleID: 1},
		{ID: 2, Email: "user2@example.com", FirstName: "User", LastName: "Two", RoleID: 1},
	}

	paginatedResult := &repository.PaginatedResult[domain.User]{
		Data:       users,
		Total:      2,
		Page:       1,
		PerPage:    10,
		TotalPages: 1,
	}

	params := repository.UserListParams{
		PaginationParams: repository.PaginationParams{Page: 1, PerPage: 10},
	}

	userRepo.On("List", ctx, params).Return(paginatedResult, nil)
	roleRepo.On("GetByID", ctx, int64(1)).Return(&domain.Role{ID: 1, Name: "admin"}, nil)

	result, err := service.List(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, 2, result.Total)
	userRepo.AssertExpectations(t)
}

func TestUserService_Delete_Success(t *testing.T) {
	service, userRepo, _, _ := setupUserService()
	ctx := context.Background()

	user := &domain.User{ID: 1}
	userRepo.On("GetByID", ctx, int64(1)).Return(user, nil)
	userRepo.On("Delete", ctx, int64(1)).Return(nil)

	err := service.Delete(ctx, 1)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestUserService_Delete_NotFound(t *testing.T) {
	service, userRepo, _, _ := setupUserService()
	ctx := context.Background()

	userRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := service.Delete(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	userRepo.AssertExpectations(t)
}

func TestUserService_ResetPassword_Success(t *testing.T) {
	service, userRepo, _, _ := setupUserService()
	ctx := context.Background()

	userRepo.On("UpdatePassword", ctx, int64(1), mock.AnythingOfType("string")).Return(nil)

	err := service.ResetPassword(ctx, 1, "NewStrongPass1!@")

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestUserService_ResetPassword_WeakPassword(t *testing.T) {
	service, _, _, _ := setupUserService()
	ctx := context.Background()

	err := service.ResetPassword(ctx, 1, "weak")

	assert.Error(t, err)
}
