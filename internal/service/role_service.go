package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// RoleService handles role business logic
type RoleService struct {
	roleRepo repository.RoleRepository
}

// NewRoleService creates a new RoleService
func NewRoleService(roleRepo repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

// CreateRoleInput represents create role request data
type CreateRoleInput struct {
	Name        string   `json:"name" validate:"required,min=2,max=50"`
	DisplayName string   `json:"display_name" validate:"required,min=2"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions" validate:"required,min=1"`
}

// UpdateRoleInput represents update role request data
type UpdateRoleInput struct {
	Name        string   `json:"name" validate:"omitempty,min=2,max=50"`
	DisplayName string   `json:"display_name" validate:"omitempty,min=2"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// Create creates a new role
func (s *RoleService) Create(ctx context.Context, input CreateRoleInput) (*domain.Role, error) {
	// Check for duplicate name
	existing, _ := s.roleRepo.GetByName(ctx, input.Name)
	if existing != nil {
		return nil, errors.New("role with this name already exists")
	}

	// Convert permissions to JSON
	permissionsJSON, err := json.Marshal(input.Permissions)
	if err != nil {
		return nil, errors.New("invalid permissions format")
	}

	role := &domain.Role{
		Name:        input.Name,
		DisplayName: input.DisplayName,
		Description: input.Description,
		Permissions: permissionsJSON,
		IsSystem:    false,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

// Update updates an existing role
func (s *RoleService) Update(ctx context.Context, id int64, input UpdateRoleInput) (*domain.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Check if it's a system role
	if role.IsSystem {
		return nil, errors.New("cannot update system role")
	}

	// Update fields
	if input.Name != "" {
		// Check for duplicate name
		existing, _ := s.roleRepo.GetByName(ctx, input.Name)
		if existing != nil && existing.ID != id {
			return nil, errors.New("role with this name already exists")
		}
		role.Name = input.Name
	}
	if input.DisplayName != "" {
		role.DisplayName = input.DisplayName
	}
	role.Description = input.Description

	if len(input.Permissions) > 0 {
		permissionsJSON, err := json.Marshal(input.Permissions)
		if err != nil {
			return nil, errors.New("invalid permissions format")
		}
		role.Permissions = permissionsJSON
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return role, nil
}

// GetByID retrieves a role by ID
func (s *RoleService) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("role not found")
	}
	return role, nil
}

// GetByName retrieves a role by name
func (s *RoleService) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	role, err := s.roleRepo.GetByName(ctx, name)
	if err != nil {
		return nil, errors.New("role not found")
	}
	return role, nil
}

// List retrieves all roles
func (s *RoleService) List(ctx context.Context) ([]*domain.Role, error) {
	return s.roleRepo.List(ctx)
}

// Delete deletes a role
func (s *RoleService) Delete(ctx context.Context, id int64) error {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("role not found")
	}

	if role.IsSystem {
		return errors.New("cannot delete system role")
	}

	return s.roleRepo.Delete(ctx, id)
}

// GetAvailablePermissions returns a list of all available permissions
func (s *RoleService) GetAvailablePermissions() []string {
	return []string{
		// Users
		"users.read",
		"users.create",
		"users.update",
		"users.delete",
		// Customers
		"customers.read",
		"customers.create",
		"customers.update",
		"customers.delete",
		// Items
		"items.read",
		"items.create",
		"items.update",
		"items.delete",
		"items.appraise",
		// Loans
		"loans.read",
		"loans.create",
		"loans.update",
		"loans.delete",
		"loans.approve",
		"loans.extend",
		"loans.default",
		// Payments
		"payments.read",
		"payments.create",
		"payments.void",
		// Sales
		"sales.read",
		"sales.create",
		"sales.update",
		"sales.delete",
		"sales.refund",
		// Categories
		"categories.read",
		"categories.create",
		"categories.update",
		"categories.delete",
		// Branches
		"branches.read",
		"branches.create",
		"branches.update",
		"branches.delete",
		// Roles
		"roles.read",
		"roles.create",
		"roles.update",
		"roles.delete",
		// Cash
		"cash.read",
		"cash.manage_registers",
		"cash.manage_sessions",
		"cash.manage_movements",
		// Reports
		"reports.read",
		"reports.export",
		// Settings
		"settings.read",
		"settings.update",
		// Audit
		"audit.read",
	}
}
