package service

import (
	"context"
	"fmt"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/pkg/cache"
)

// CachedSettingService wraps SettingService with caching
type CachedSettingService struct {
	*SettingService
	cache *cache.Cache
}

// NewCachedSettingService creates a new CachedSettingService
func NewCachedSettingService(settingRepo repository.SettingRepository, c *cache.Cache) *CachedSettingService {
	return &CachedSettingService{
		SettingService: NewSettingService(settingRepo),
		cache:          c,
	}
}

func (s *CachedSettingService) settingKey(key string, branchID *int64) string {
	if branchID != nil {
		return fmt.Sprintf("settings:%s:branch:%d", key, *branchID)
	}
	return fmt.Sprintf("settings:%s:global", key)
}

func (s *CachedSettingService) allSettingsKey(branchID *int64) string {
	if branchID != nil {
		return fmt.Sprintf("settings:all:branch:%d", *branchID)
	}
	return "settings:all:global"
}

// Get retrieves a setting by key with caching
func (s *CachedSettingService) Get(ctx context.Context, key string, branchID *int64) (*domain.Setting, error) {
	if s.cache == nil {
		return s.SettingService.Get(ctx, key, branchID)
	}

	cacheKey := s.settingKey(key, branchID)
	var setting domain.Setting

	if err := s.cache.Get(ctx, cacheKey, &setting); err == nil {
		return &setting, nil
	}

	// Cache miss - fetch from DB
	result, err := s.SettingService.Get(ctx, key, branchID)
	if err != nil {
		return nil, err
	}

	// Store in cache
	_ = s.cache.Set(ctx, cacheKey, result, cache.SettingsTTL)

	return result, nil
}

// GetAll retrieves all settings with caching
func (s *CachedSettingService) GetAll(ctx context.Context, branchID *int64) ([]*domain.Setting, error) {
	if s.cache == nil {
		return s.SettingService.GetAll(ctx, branchID)
	}

	cacheKey := s.allSettingsKey(branchID)
	var settings []*domain.Setting

	if err := s.cache.Get(ctx, cacheKey, &settings); err == nil {
		return settings, nil
	}

	// Cache miss - fetch from DB
	result, err := s.SettingService.GetAll(ctx, branchID)
	if err != nil {
		return nil, err
	}

	// Store in cache
	_ = s.cache.Set(ctx, cacheKey, result, cache.SettingsTTL)

	return result, nil
}

// GetMerged retrieves all settings merged with caching
func (s *CachedSettingService) GetMerged(ctx context.Context, branchID *int64) (map[string]interface{}, error) {
	if s.cache == nil {
		return s.SettingService.GetMerged(ctx, branchID)
	}

	cacheKey := fmt.Sprintf("settings:merged:branch:%v", branchID)
	var merged map[string]interface{}

	if err := s.cache.Get(ctx, cacheKey, &merged); err == nil {
		return merged, nil
	}

	// Cache miss - fetch from DB
	result, err := s.SettingService.GetMerged(ctx, branchID)
	if err != nil {
		return nil, err
	}

	// Store in cache
	_ = s.cache.Set(ctx, cacheKey, result, cache.SettingsTTL)

	return result, nil
}

// Set creates or updates a setting and invalidates cache
func (s *CachedSettingService) Set(ctx context.Context, input SetSettingInput) (*domain.Setting, error) {
	result, err := s.SettingService.Set(ctx, input)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	if s.cache != nil {
		_ = s.cache.Delete(ctx, s.settingKey(input.Key, input.BranchID))
		_ = s.cache.Delete(ctx, s.allSettingsKey(input.BranchID))
		_ = s.cache.DeleteByPattern(ctx, "settings:merged:*")
	}

	return result, nil
}

// SetMultiple sets multiple settings and invalidates cache
func (s *CachedSettingService) SetMultiple(ctx context.Context, settings []SetSettingInput) error {
	err := s.SettingService.SetMultiple(ctx, settings)
	if err != nil {
		return err
	}

	// Invalidate all settings cache
	if s.cache != nil {
		_ = s.cache.DeleteByPattern(ctx, "settings:*")
	}

	return nil
}

// Delete deletes a setting and invalidates cache
func (s *CachedSettingService) Delete(ctx context.Context, key string, branchID *int64) error {
	err := s.SettingService.Delete(ctx, key, branchID)
	if err != nil {
		return err
	}

	// Invalidate cache
	if s.cache != nil {
		_ = s.cache.Delete(ctx, s.settingKey(key, branchID))
		_ = s.cache.Delete(ctx, s.allSettingsKey(branchID))
		_ = s.cache.DeleteByPattern(ctx, "settings:merged:*")
	}

	return nil
}
