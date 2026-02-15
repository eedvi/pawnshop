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

func setupSettingService() (*SettingService, *mocks.MockSettingRepository) {
	settingRepo := new(mocks.MockSettingRepository)
	service := NewSettingService(settingRepo)
	return service, settingRepo
}

func TestSettingService_Get_Success(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	setting := &domain.Setting{Key: "app_name", Value: "PawnShop"}

	settingRepo.On("Get", ctx, "app_name", (*int64)(nil)).Return(setting, nil)

	result, err := service.Get(ctx, "app_name", nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "app_name", result.Key)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_Get_NotFound(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "missing", (*int64)(nil)).Return(nil, errors.New("not found"))

	result, err := service.Get(ctx, "missing", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "setting not found", err.Error())
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetAll_Success(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settings := []*domain.Setting{
		{Key: "key1", Value: "val1"},
		{Key: "key2", Value: "val2"},
	}

	settingRepo.On("GetAll", ctx, (*int64)(nil)).Return(settings, nil)

	result, err := service.GetAll(ctx, nil)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetMerged_GlobalOnly(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settings := []*domain.Setting{
		{Key: "key1", Value: "global_val", BranchID: nil},
		{Key: "key2", Value: "global_val2", BranchID: nil},
	}

	settingRepo.On("GetAll", ctx, (*int64)(nil)).Return(settings, nil)

	result, err := service.GetMerged(ctx, nil)

	assert.NoError(t, err)
	assert.Equal(t, "global_val", result["key1"])
	assert.Equal(t, "global_val2", result["key2"])
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetMerged_BranchOverridesGlobal(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	branchID := int64(1)
	settings := []*domain.Setting{
		{Key: "key1", Value: "global_val", BranchID: nil},
		{Key: "key1", Value: "branch_val", BranchID: &branchID},
		{Key: "key2", Value: "global_only", BranchID: nil},
	}

	settingRepo.On("GetAll", ctx, &branchID).Return(settings, nil)

	result, err := service.GetMerged(ctx, &branchID)

	assert.NoError(t, err)
	assert.Equal(t, "branch_val", result["key1"])
	assert.Equal(t, "global_only", result["key2"])
	settingRepo.AssertExpectations(t)
}

func TestSettingService_Set_Success(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Set", ctx, mock.AnythingOfType("*domain.Setting")).Return(nil)

	input := SetSettingInput{Key: "app_name", Value: "PawnShop", Description: "Application name"}
	result, err := service.Set(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "app_name", result.Key)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_Set_Error(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Set", ctx, mock.AnythingOfType("*domain.Setting")).Return(errors.New("db error"))

	input := SetSettingInput{Key: "app_name", Value: "PawnShop"}
	result, err := service.Set(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to save setting", err.Error())
	settingRepo.AssertExpectations(t)
}

func TestSettingService_SetMultiple_Success(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Set", ctx, mock.AnythingOfType("*domain.Setting")).Return(nil).Times(2)

	inputs := []SetSettingInput{
		{Key: "key1", Value: "val1"},
		{Key: "key2", Value: "val2"},
	}
	err := service.SetMultiple(ctx, inputs)

	assert.NoError(t, err)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_SetMultiple_ErrorOnSecond(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Set", ctx, mock.AnythingOfType("*domain.Setting")).Return(nil).Once()
	settingRepo.On("Set", ctx, mock.AnythingOfType("*domain.Setting")).Return(errors.New("db error")).Once()

	inputs := []SetSettingInput{
		{Key: "key1", Value: "val1"},
		{Key: "key2", Value: "val2"},
	}
	err := service.SetMultiple(ctx, inputs)

	assert.Error(t, err)
	assert.Equal(t, "failed to save settings", err.Error())
	settingRepo.AssertExpectations(t)
}

func TestSettingService_Delete_Success(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Delete", ctx, "key1", (*int64)(nil)).Return(nil)

	err := service.Delete(ctx, "key1", nil)

	assert.NoError(t, err)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetString_Success(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "name", (*int64)(nil)).Return(&domain.Setting{Key: "name", Value: "PawnShop"}, nil)

	result := service.GetString(ctx, "name", nil, "default")

	assert.Equal(t, "PawnShop", result)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetString_NotFound_ReturnsDefault(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "missing", (*int64)(nil)).Return(nil, errors.New("not found"))

	result := service.GetString(ctx, "missing", nil, "default_val")

	assert.Equal(t, "default_val", result)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetString_WrongType_ReturnsDefault(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "num", (*int64)(nil)).Return(&domain.Setting{Key: "num", Value: 42}, nil)

	result := service.GetString(ctx, "num", nil, "default_val")

	assert.Equal(t, "default_val", result)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetInt_Success(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	// JSON numbers are float64
	settingRepo.On("Get", ctx, "count", (*int64)(nil)).Return(&domain.Setting{Key: "count", Value: float64(42)}, nil)

	result := service.GetInt(ctx, "count", nil, 0)

	assert.Equal(t, 42, result)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetInt_FromInt(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "count", (*int64)(nil)).Return(&domain.Setting{Key: "count", Value: 42}, nil)

	result := service.GetInt(ctx, "count", nil, 0)

	assert.Equal(t, 42, result)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetInt_FromInt64(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "count", (*int64)(nil)).Return(&domain.Setting{Key: "count", Value: int64(42)}, nil)

	result := service.GetInt(ctx, "count", nil, 0)

	assert.Equal(t, 42, result)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetInt_NotFound_ReturnsDefault(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "missing", (*int64)(nil)).Return(nil, errors.New("not found"))

	result := service.GetInt(ctx, "missing", nil, 99)

	assert.Equal(t, 99, result)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetInt_WrongType_ReturnsDefault(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "str", (*int64)(nil)).Return(&domain.Setting{Key: "str", Value: "not a number"}, nil)

	result := service.GetInt(ctx, "str", nil, 99)

	assert.Equal(t, 99, result)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetFloat_Success(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "rate", (*int64)(nil)).Return(&domain.Setting{Key: "rate", Value: 3.14}, nil)

	result := service.GetFloat(ctx, "rate", nil, 0.0)

	assert.Equal(t, 3.14, result)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetFloat_NotFound_ReturnsDefault(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "missing", (*int64)(nil)).Return(nil, errors.New("not found"))

	result := service.GetFloat(ctx, "missing", nil, 1.5)

	assert.Equal(t, 1.5, result)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetBool_Success(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "enabled", (*int64)(nil)).Return(&domain.Setting{Key: "enabled", Value: true}, nil)

	result := service.GetBool(ctx, "enabled", nil, false)

	assert.True(t, result)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetBool_NotFound_ReturnsDefault(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "missing", (*int64)(nil)).Return(nil, errors.New("not found"))

	result := service.GetBool(ctx, "missing", nil, true)

	assert.True(t, result)
	settingRepo.AssertExpectations(t)
}

func TestSettingService_GetBool_WrongType_ReturnsDefault(t *testing.T) {
	service, settingRepo := setupSettingService()
	ctx := context.Background()

	settingRepo.On("Get", ctx, "str", (*int64)(nil)).Return(&domain.Setting{Key: "str", Value: "not bool"}, nil)

	result := service.GetBool(ctx, "str", nil, false)

	assert.False(t, result)
	settingRepo.AssertExpectations(t)
}
