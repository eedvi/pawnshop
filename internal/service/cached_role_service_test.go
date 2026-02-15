package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository/mocks"
)

func setupCachedRoleService() (*CachedRoleService, *mocks.MockRoleRepository) {
	roleRepo := new(mocks.MockRoleRepository)
	service := NewCachedRoleService(roleRepo, nil)
	return service, roleRepo
}

func TestCachedRoleService_GetByID_NilCache_Success(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	role := &domain.Role{ID: 1, Name: "admin"}
	roleRepo.On("GetByID", ctx, int64(1)).Return(role, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "admin", result.Name)
	roleRepo.AssertExpectations(t)
}

func TestCachedRoleService_GetByID_NilCache_NotFound(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	roleRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	roleRepo.AssertExpectations(t)
}

func TestCachedRoleService_GetByName_NilCache_Success(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	role := &domain.Role{ID: 1, Name: "admin"}
	roleRepo.On("GetByName", ctx, "admin").Return(role, nil)

	result, err := service.GetByName(ctx, "admin")

	assert.NoError(t, err)
	assert.Equal(t, "admin", result.Name)
	roleRepo.AssertExpectations(t)
}

func TestCachedRoleService_GetByName_NilCache_NotFound(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	roleRepo.On("GetByName", ctx, "missing").Return(nil, errors.New("not found"))

	result, err := service.GetByName(ctx, "missing")

	assert.Error(t, err)
	assert.Nil(t, result)
	roleRepo.AssertExpectations(t)
}

func TestCachedRoleService_List_NilCache_Success(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	roles := []*domain.Role{
		{ID: 1, Name: "admin"},
		{ID: 2, Name: "cashier"},
	}
	roleRepo.On("List", ctx).Return(roles, nil)

	result, err := service.List(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	roleRepo.AssertExpectations(t)
}

func TestCachedRoleService_List_NilCache_Error(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	roleRepo.On("List", ctx).Return(nil, errors.New("db error"))

	result, err := service.List(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	roleRepo.AssertExpectations(t)
}

func TestCachedRoleService_Create_NilCache_Success(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	roleRepo.On("GetByName", ctx, "editor").Return(nil, errors.New("not found"))
	roleRepo.On("Create", ctx, mock.AnythingOfType("*domain.Role")).Return(nil)

	input := CreateRoleInput{
		Name:        "editor",
		DisplayName: "Editor",
		Permissions: []string{"users.read"},
	}
	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "editor", result.Name)
	roleRepo.AssertExpectations(t)
}

func TestCachedRoleService_Create_NilCache_DuplicateName(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	existing := &domain.Role{ID: 1, Name: "admin"}
	roleRepo.On("GetByName", ctx, "admin").Return(existing, nil)

	input := CreateRoleInput{
		Name:        "admin",
		DisplayName: "Admin",
		Permissions: []string{"*"},
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "role with this name already exists", err.Error())
	roleRepo.AssertExpectations(t)
}

func TestCachedRoleService_Update_NilCache_Success(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	existing := &domain.Role{
		ID:          1,
		Name:        "editor",
		Permissions: json.RawMessage(`["users.read"]`),
	}
	roleRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	roleRepo.On("Update", ctx, mock.AnythingOfType("*domain.Role")).Return(nil)

	input := UpdateRoleInput{DisplayName: "Updated Editor"}
	result, err := service.Update(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Editor", result.DisplayName)
	roleRepo.AssertExpectations(t)
}

func TestCachedRoleService_Update_NilCache_SystemRole(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	existing := &domain.Role{ID: 1, Name: "super_admin", IsSystem: true}
	roleRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)

	input := UpdateRoleInput{DisplayName: "Hacked"}
	result, err := service.Update(ctx, 1, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "cannot update system role", err.Error())
	roleRepo.AssertExpectations(t)
}

func TestCachedRoleService_Delete_NilCache_Success(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	existing := &domain.Role{ID: 1, Name: "editor", IsSystem: false}
	roleRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	roleRepo.On("Delete", ctx, int64(1)).Return(nil)

	err := service.Delete(ctx, 1)

	assert.NoError(t, err)
	roleRepo.AssertExpectations(t)
}

func TestCachedRoleService_Delete_NilCache_SystemRole(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	existing := &domain.Role{ID: 1, Name: "super_admin", IsSystem: true}
	roleRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)

	err := service.Delete(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, "cannot delete system role", err.Error())
	roleRepo.AssertExpectations(t)
}

func TestCachedRoleService_Delete_NilCache_NotFound(t *testing.T) {
	service, roleRepo := setupCachedRoleService()
	ctx := context.Background()

	roleRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := service.Delete(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, "role not found", err.Error())
	roleRepo.AssertExpectations(t)
}
