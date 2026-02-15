package service

import (
	"context"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/pkg/cache"
)

// CachedRoleService wraps RoleService with caching
type CachedRoleService struct {
	*RoleService
	cache *cache.Cache
}

// NewCachedRoleService creates a new CachedRoleService
func NewCachedRoleService(roleRepo repository.RoleRepository, c *cache.Cache) *CachedRoleService {
	return &CachedRoleService{
		RoleService: NewRoleService(roleRepo),
		cache:       c,
	}
}

// GetByID retrieves a role by ID with caching
func (s *CachedRoleService) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	if s.cache == nil {
		return s.RoleService.GetByID(ctx, id)
	}

	cacheKey := cache.RoleKey(id)
	var role domain.Role

	if err := s.cache.Get(ctx, cacheKey, &role); err == nil {
		return &role, nil
	}

	// Cache miss - fetch from DB
	result, err := s.RoleService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache
	_ = s.cache.Set(ctx, cacheKey, result, cache.RolesTTL)

	return result, nil
}

// GetByName retrieves a role by name with caching
func (s *CachedRoleService) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	if s.cache == nil {
		return s.RoleService.GetByName(ctx, name)
	}

	cacheKey := "roles:name:" + name
	var role domain.Role

	if err := s.cache.Get(ctx, cacheKey, &role); err == nil {
		return &role, nil
	}

	// Cache miss - fetch from DB
	result, err := s.RoleService.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// Store in cache
	_ = s.cache.Set(ctx, cacheKey, result, cache.RolesTTL)

	return result, nil
}

// List retrieves all roles with caching
func (s *CachedRoleService) List(ctx context.Context) ([]*domain.Role, error) {
	if s.cache == nil {
		return s.RoleService.List(ctx)
	}

	cacheKey := cache.RolesAllKey
	var roles []*domain.Role

	if err := s.cache.Get(ctx, cacheKey, &roles); err == nil {
		return roles, nil
	}

	// Cache miss - fetch from DB
	result, err := s.RoleService.List(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache
	_ = s.cache.Set(ctx, cacheKey, result, cache.RolesTTL)

	return result, nil
}

// Create creates a new role and invalidates cache
func (s *CachedRoleService) Create(ctx context.Context, input CreateRoleInput) (*domain.Role, error) {
	result, err := s.RoleService.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	if s.cache != nil {
		_ = s.cache.Delete(ctx, cache.RolesAllKey)
	}

	return result, nil
}

// Update updates a role and invalidates cache
func (s *CachedRoleService) Update(ctx context.Context, id int64, input UpdateRoleInput) (*domain.Role, error) {
	result, err := s.RoleService.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	if s.cache != nil {
		_ = s.cache.Delete(ctx, cache.RoleKey(id))
		_ = s.cache.Delete(ctx, cache.RolesAllKey)
		_ = s.cache.DeleteByPattern(ctx, "roles:name:*")
		// Invalidate user permissions that might use this role
		_ = s.cache.DeleteByPattern(ctx, "users:*:permissions")
	}

	return result, nil
}

// Delete deletes a role and invalidates cache
func (s *CachedRoleService) Delete(ctx context.Context, id int64) error {
	err := s.RoleService.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache
	if s.cache != nil {
		_ = s.cache.Delete(ctx, cache.RoleKey(id))
		_ = s.cache.Delete(ctx, cache.RolesAllKey)
		_ = s.cache.DeleteByPattern(ctx, "roles:name:*")
		// Invalidate user permissions that might use this role
		_ = s.cache.DeleteByPattern(ctx, "users:*:permissions")
	}

	return nil
}
