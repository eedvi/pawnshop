package service

import (
	"context"
	"errors"
	"fmt"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// SettingService handles settings business logic
type SettingService struct {
	settingRepo repository.SettingRepository
}

// NewSettingService creates a new SettingService
func NewSettingService(settingRepo repository.SettingRepository) *SettingService {
	return &SettingService{settingRepo: settingRepo}
}

// Get retrieves a setting by key
func (s *SettingService) Get(ctx context.Context, key string, branchID *int64) (*domain.Setting, error) {
	setting, err := s.settingRepo.Get(ctx, key, branchID)
	if err != nil {
		return nil, errors.New("setting not found")
	}
	return setting, nil
}

// GetAll retrieves all settings
func (s *SettingService) GetAll(ctx context.Context, branchID *int64) ([]*domain.Setting, error) {
	return s.settingRepo.GetAll(ctx, branchID)
}

// GetMerged retrieves all settings merged (branch settings override global)
func (s *SettingService) GetMerged(ctx context.Context, branchID *int64) (map[string]interface{}, error) {
	settings, err := s.settingRepo.GetAll(ctx, branchID)
	if err != nil {
		return nil, err
	}

	// Create merged map - branch settings override global
	merged := make(map[string]interface{})
	for _, setting := range settings {
		// If branchID matches or it's a global setting and not already set
		if setting.BranchID != nil && branchID != nil && *setting.BranchID == *branchID {
			merged[setting.Key] = setting.Value
		} else if setting.BranchID == nil {
			if _, exists := merged[setting.Key]; !exists {
				merged[setting.Key] = setting.Value
			}
		}
	}

	return merged, nil
}

// SetSettingInput represents set setting request data
type SetSettingInput struct {
	Key         string      `json:"key" validate:"required,min=1"`
	Value       interface{} `json:"value" validate:"required"`
	Description string      `json:"description"`
	BranchID    *int64      `json:"branch_id"`
}

// Set creates or updates a setting
func (s *SettingService) Set(ctx context.Context, input SetSettingInput) (*domain.Setting, error) {
	setting := &domain.Setting{
		Key:         input.Key,
		Value:       input.Value,
		Description: input.Description,
		BranchID:    input.BranchID,
	}

	if err := s.settingRepo.Set(ctx, setting); err != nil {
		return nil, fmt.Errorf("failed to save setting: %w", err)
	}

	return setting, nil
}

// SetMultiple sets multiple settings at once
func (s *SettingService) SetMultiple(ctx context.Context, settings []SetSettingInput) error {
	for _, input := range settings {
		setting := &domain.Setting{
			Key:         input.Key,
			Value:       input.Value,
			Description: input.Description,
			BranchID:    input.BranchID,
		}

		if err := s.settingRepo.Set(ctx, setting); err != nil {
			return fmt.Errorf("failed to save settings: %w", err)
		}
	}

	return nil
}

// Delete deletes a setting
func (s *SettingService) Delete(ctx context.Context, key string, branchID *int64) error {
	return s.settingRepo.Delete(ctx, key, branchID)
}

// GetString retrieves a setting value as string
func (s *SettingService) GetString(ctx context.Context, key string, branchID *int64, defaultValue string) string {
	setting, err := s.settingRepo.Get(ctx, key, branchID)
	if err != nil {
		return defaultValue
	}

	if str, ok := setting.Value.(string); ok {
		return str
	}
	return defaultValue
}

// GetInt retrieves a setting value as int
func (s *SettingService) GetInt(ctx context.Context, key string, branchID *int64, defaultValue int) int {
	setting, err := s.settingRepo.Get(ctx, key, branchID)
	if err != nil {
		return defaultValue
	}

	switch v := setting.Value.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	}
	return defaultValue
}

// GetFloat retrieves a setting value as float
func (s *SettingService) GetFloat(ctx context.Context, key string, branchID *int64, defaultValue float64) float64 {
	setting, err := s.settingRepo.Get(ctx, key, branchID)
	if err != nil {
		return defaultValue
	}

	if f, ok := setting.Value.(float64); ok {
		return f
	}
	return defaultValue
}

// GetBool retrieves a setting value as bool
func (s *SettingService) GetBool(ctx context.Context, key string, branchID *int64, defaultValue bool) bool {
	setting, err := s.settingRepo.Get(ctx, key, branchID)
	if err != nil {
		return defaultValue
	}

	if b, ok := setting.Value.(bool); ok {
		return b
	}
	return defaultValue
}
