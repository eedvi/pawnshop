package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository/mocks"
)

func setupRoleService() (*RoleService, *mocks.MockRoleRepository) {
	roleRepo := new(mocks.MockRoleRepository)
	service := NewRoleService(roleRepo)
	return service, roleRepo
}

func TestRoleService_Create_Success(t *testing.T) {
	service, roleRepo := setupRoleService()
	ctx := context.Background()

	// Setup expectations
	roleRepo.On("GetByName", ctx, "test-role").Return(nil, errors.New("not found"))
	roleRepo.On("Create", ctx, mock.AnythingOfType("*domain.Role")).Return(nil)

	// Execute
	input := CreateRoleInput{
		Name:        "test-role",
		DisplayName: "Test Role",
		Description: "A test role",
		Permissions: []string{"users.read", "users.create"},
	}
	result, err := service.Create(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-role", result.Name)
	assert.Equal(t, "Test Role", result.DisplayName)
	assert.False(t, result.IsSystem)

	roleRepo.AssertExpectations(t)
}

func TestRoleService_Create_DuplicateName(t *testing.T) {
	service, roleRepo := setupRoleService()
	ctx := context.Background()

	existingRole := &domain.Role{
		ID:   1,
		Name: "existing-role",
	}

	// Setup expectations
	roleRepo.On("GetByName", ctx, "existing-role").Return(existingRole, nil)

	// Execute
	input := CreateRoleInput{
		Name:        "existing-role",
		DisplayName: "Existing Role",
		Permissions: []string{"users.read"},
	}
	result, err := service.Create(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "role with this name already exists", err.Error())

	roleRepo.AssertExpectations(t)
}

func TestRoleService_Update_Success(t *testing.T) {
	service, roleRepo := setupRoleService()
	ctx := context.Background()

	existingRole := &domain.Role{
		ID:          1,
		Name:        "old-name",
		DisplayName: "Old Name",
		IsSystem:    false,
	}

	// Setup expectations
	roleRepo.On("GetByID", ctx, int64(1)).Return(existingRole, nil)
	roleRepo.On("GetByName", ctx, "new-name").Return(nil, errors.New("not found"))
	roleRepo.On("Update", ctx, mock.AnythingOfType("*domain.Role")).Return(nil)

	// Execute
	input := UpdateRoleInput{
		Name:        "new-name",
		DisplayName: "New Name",
		Permissions: []string{"users.read", "users.update"},
	}
	result, err := service.Update(ctx, 1, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new-name", result.Name)
	assert.Equal(t, "New Name", result.DisplayName)

	roleRepo.AssertExpectations(t)
}

func TestRoleService_Update_SystemRole(t *testing.T) {
	service, roleRepo := setupRoleService()
	ctx := context.Background()

	systemRole := &domain.Role{
		ID:          1,
		Name:        "super_admin",
		DisplayName: "Super Admin",
		IsSystem:    true,
	}

	// Setup expectations
	roleRepo.On("GetByID", ctx, int64(1)).Return(systemRole, nil)

	// Execute
	input := UpdateRoleInput{
		Name: "new-name",
	}
	result, err := service.Update(ctx, 1, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "cannot update system role", err.Error())

	roleRepo.AssertExpectations(t)
}

func TestRoleService_Update_NotFound(t *testing.T) {
	service, roleRepo := setupRoleService()
	ctx := context.Background()

	// Setup expectations
	roleRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	// Execute
	input := UpdateRoleInput{
		Name: "new-name",
	}
	result, err := service.Update(ctx, 999, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "role not found", err.Error())

	roleRepo.AssertExpectations(t)
}

func TestRoleService_GetByID_Success(t *testing.T) {
	service, roleRepo := setupRoleService()
	ctx := context.Background()

	role := &domain.Role{
		ID:          1,
		Name:        "admin",
		DisplayName: "Administrator",
	}

	// Setup expectations
	roleRepo.On("GetByID", ctx, int64(1)).Return(role, nil)

	// Execute
	result, err := service.GetByID(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "admin", result.Name)

	roleRepo.AssertExpectations(t)
}

func TestRoleService_GetByID_NotFound(t *testing.T) {
	service, roleRepo := setupRoleService()
	ctx := context.Background()

	// Setup expectations
	roleRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	// Execute
	result, err := service.GetByID(ctx, 999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "role not found", err.Error())

	roleRepo.AssertExpectations(t)
}

func TestRoleService_List_Success(t *testing.T) {
	service, roleRepo := setupRoleService()
	ctx := context.Background()

	roles := []*domain.Role{
		{ID: 1, Name: "admin"},
		{ID: 2, Name: "manager"},
	}

	// Setup expectations
	roleRepo.On("List", ctx).Return(roles, nil)

	// Execute
	result, err := service.List(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	roleRepo.AssertExpectations(t)
}

func TestRoleService_Delete_Success(t *testing.T) {
	service, roleRepo := setupRoleService()
	ctx := context.Background()

	role := &domain.Role{
		ID:       1,
		Name:     "custom-role",
		IsSystem: false,
	}

	// Setup expectations
	roleRepo.On("GetByID", ctx, int64(1)).Return(role, nil)
	roleRepo.On("Delete", ctx, int64(1)).Return(nil)

	// Execute
	err := service.Delete(ctx, 1)

	// Assert
	assert.NoError(t, err)
	roleRepo.AssertExpectations(t)
}

func TestRoleService_Delete_SystemRole(t *testing.T) {
	service, roleRepo := setupRoleService()
	ctx := context.Background()

	systemRole := &domain.Role{
		ID:       1,
		Name:     "super_admin",
		IsSystem: true,
	}

	// Setup expectations
	roleRepo.On("GetByID", ctx, int64(1)).Return(systemRole, nil)

	// Execute
	err := service.Delete(ctx, 1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "cannot delete system role", err.Error())

	roleRepo.AssertExpectations(t)
}

func TestRoleService_GetAvailablePermissions(t *testing.T) {
	service, _ := setupRoleService()

	// Execute
	permissions := service.GetAvailablePermissions()

	// Assert
	assert.NotEmpty(t, permissions)
	assert.Contains(t, permissions, "users.read")
	assert.Contains(t, permissions, "loans.create")
	assert.Contains(t, permissions, "payments.void")
}
