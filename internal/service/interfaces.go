package service

import (
	"context"

	"pawnshop/internal/domain"
)

// RoleServiceInterface defines the interface for role service operations
type RoleServiceInterface interface {
	Create(ctx context.Context, input CreateRoleInput) (*domain.Role, error)
	Update(ctx context.Context, id int64, input UpdateRoleInput) (*domain.Role, error)
	GetByID(ctx context.Context, id int64) (*domain.Role, error)
	GetByName(ctx context.Context, name string) (*domain.Role, error)
	List(ctx context.Context) ([]*domain.Role, error)
	Delete(ctx context.Context, id int64) error
	GetAvailablePermissions() []string
}

// SettingServiceInterface defines the interface for setting service operations
type SettingServiceInterface interface {
	Get(ctx context.Context, key string, branchID *int64) (*domain.Setting, error)
	GetAll(ctx context.Context, branchID *int64) ([]*domain.Setting, error)
	GetMerged(ctx context.Context, branchID *int64) (map[string]interface{}, error)
	Set(ctx context.Context, input SetSettingInput) (*domain.Setting, error)
	SetMultiple(ctx context.Context, settings []SetSettingInput) error
	Delete(ctx context.Context, key string, branchID *int64) error
	GetString(ctx context.Context, key string, branchID *int64, defaultValue string) string
	GetInt(ctx context.Context, key string, branchID *int64, defaultValue int) int
	GetFloat(ctx context.Context, key string, branchID *int64, defaultValue float64) float64
	GetBool(ctx context.Context, key string, branchID *int64, defaultValue bool) bool
}

// Ensure concrete types implement interfaces
var _ RoleServiceInterface = (*RoleService)(nil)
var _ RoleServiceInterface = (*CachedRoleService)(nil)
var _ SettingServiceInterface = (*SettingService)(nil)
var _ SettingServiceInterface = (*CachedSettingService)(nil)
