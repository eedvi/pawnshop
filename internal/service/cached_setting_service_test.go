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

func setupCachedSettingService() (*CachedSettingService, *mocks.MockSettingRepository) {
	settingRepo := new(mocks.MockSettingRepository)
	service := NewCachedSettingService(settingRepo, nil)
	return service, settingRepo
}

func TestCachedSettingService_Get_NilCache_Success(t *testing.T) {
	service, settingRepo := setupCachedSettingService()
	ctx := context.Background()

	setting := &domain.Setting{Key: "app_name", Value: "PawnShop"}
	settingRepo.On("Get", ctx, "app_name", (*int64)(nil)).Return(setting, nil)

	result, err := service.Get(ctx, "app_name", nil)

	assert.NoError(t, err)
	assert.Equal(t, "app_name", result.Key)
	settingRepo.AssertExpectations(t)
}

func TestCachedSettingService_Get_NilCache_WithBranch(t *testing.T) {
	service, settingRepo := setupCachedSettingService()
	ctx := context.Background()

	branchID := int64(1)
	setting := &domain.Setting{Key: "tax_rate", Value: 0.12, BranchID: &branchID}
	settingRepo.On("Get", ctx, "tax_rate", &branchID).Return(setting, nil)

	result, err := service.Get(ctx, "tax_rate", &branchID)

	assert.NoError(t, err)
	assert.Equal(t, "tax_rate", result.Key)
	settingRepo.AssertExpectations(t)
}

func TestCachedSettingService_Get_NilCache_NotFound(t *testing.T) {
	service, settingRepo := setupCachedSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "missing", (*int64)(nil)).Return(nil, errors.New("not found"))

	result, err := service.Get(ctx, "missing", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	settingRepo.AssertExpectations(t)
}

func TestCachedSettingService_GetAll_NilCache_Success(t *testing.T) {
	service, settingRepo := setupCachedSettingService()
	ctx := context.Background()

	settings := []*domain.Setting{
		{Key: "app_name", Value: "PawnShop"},
		{Key: "currency", Value: "GTQ"},
	}
	settingRepo.On("GetAll", ctx, (*int64)(nil)).Return(settings, nil)

	result, err := service.GetAll(ctx, nil)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	settingRepo.AssertExpectations(t)
}

func TestCachedSettingService_GetAll_NilCache_WithBranch(t *testing.T) {
	service, settingRepo := setupCachedSettingService()
	ctx := context.Background()

	branchID := int64(1)
	settings := []*domain.Setting{
		{Key: "tax_rate", Value: 0.12, BranchID: &branchID},
	}
	settingRepo.On("GetAll", ctx, &branchID).Return(settings, nil)

	result, err := service.GetAll(ctx, &branchID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	settingRepo.AssertExpectations(t)
}

func TestCachedSettingService_GetMerged_NilCache_Success(t *testing.T) {
	service, settingRepo := setupCachedSettingService()
	ctx := context.Background()

	branchID := int64(1)

	// GetMerged calls GetAll once with the branchID; the repo returns both global and branch settings
	allSettings := []*domain.Setting{
		{Key: "app_name", Value: "PawnShop"},
		{Key: "currency", Value: "GTQ"},
		{Key: "currency", Value: "USD", BranchID: &branchID},
	}
	settingRepo.On("GetAll", ctx, &branchID).Return(allSettings, nil)

	result, err := service.GetMerged(ctx, &branchID)

	assert.NoError(t, err)
	assert.Equal(t, "PawnShop", result["app_name"])
	assert.Equal(t, "USD", result["currency"]) // Branch overrides global
	settingRepo.AssertExpectations(t)
}

func TestCachedSettingService_GetMerged_NilCache_GlobalOnly(t *testing.T) {
	service, settingRepo := setupCachedSettingService()
	ctx := context.Background()

	globalSettings := []*domain.Setting{
		{Key: "app_name", Value: "PawnShop"},
	}
	settingRepo.On("GetAll", ctx, (*int64)(nil)).Return(globalSettings, nil)

	result, err := service.GetMerged(ctx, nil)

	assert.NoError(t, err)
	assert.Equal(t, "PawnShop", result["app_name"])
	settingRepo.AssertExpectations(t)
}

func TestCachedSettingService_Set_NilCache_Success(t *testing.T) {
	service, settingRepo := setupCachedSettingService()
	ctx := context.Background()

	settingRepo.On("Set", ctx, mock.AnythingOfType("*domain.Setting")).Return(nil)

	input := SetSettingInput{Key: "app_name", Value: "NewPawnShop"}
	result, err := service.Set(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "app_name", result.Key)
	settingRepo.AssertExpectations(t)
}

func TestCachedSettingService_SetMultiple_NilCache_Success(t *testing.T) {
	service, settingRepo := setupCachedSettingService()
	ctx := context.Background()

	settingRepo.On("Set", ctx, mock.AnythingOfType("*domain.Setting")).Return(nil)

	settings := []SetSettingInput{
		{Key: "key1", Value: "val1"},
		{Key: "key2", Value: "val2"},
	}
	err := service.SetMultiple(ctx, settings)

	assert.NoError(t, err)
	settingRepo.AssertExpectations(t)
}

func TestCachedSettingService_Delete_NilCache_Success(t *testing.T) {
	service, settingRepo := setupCachedSettingService()
	ctx := context.Background()

	settingRepo.On("Delete", ctx, "old_key", (*int64)(nil)).Return(nil)

	err := service.Delete(ctx, "old_key", nil)

	assert.NoError(t, err)
	settingRepo.AssertExpectations(t)
}

func TestCachedSettingService_Delete_NilCache_WithBranch(t *testing.T) {
	service, settingRepo := setupCachedSettingService()
	ctx := context.Background()

	branchID := int64(1)
	settingRepo.On("Delete", ctx, "branch_key", &branchID).Return(nil)

	err := service.Delete(ctx, "branch_key", &branchID)

	assert.NoError(t, err)
	settingRepo.AssertExpectations(t)
}

func TestCachedSettingService_Delete_NilCache_Error(t *testing.T) {
	service, settingRepo := setupCachedSettingService()
	ctx := context.Background()

	settingRepo.On("Delete", ctx, "missing", (*int64)(nil)).Return(errors.New("not found"))

	err := service.Delete(ctx, "missing", nil)

	assert.Error(t, err)
	settingRepo.AssertExpectations(t)
}
