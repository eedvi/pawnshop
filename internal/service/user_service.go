package service

import (
	"context"
	"errors"
	"fmt"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/pkg/auth"
)

// UserService handles user business logic
type UserService struct {
	userRepo        repository.UserRepository
	roleRepo        repository.RoleRepository
	branchRepo      repository.BranchRepository
	passwordManager *auth.PasswordManager
}

// NewUserService creates a new UserService
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	branchRepo repository.BranchRepository,
	passwordManager *auth.PasswordManager,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		roleRepo:        roleRepo,
		branchRepo:      branchRepo,
		passwordManager: passwordManager,
	}
}

// CreateUserInput represents create user request data
type CreateUserInput struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Phone     string `json:"phone"`
	RoleID    int64  `json:"role_id" validate:"required"`
	BranchID  *int64 `json:"branch_id"`
	IsActive  bool   `json:"is_active"`
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*domain.UserPublic, error) {
	// Check if email is already taken
	existing, _ := s.userRepo.GetByEmail(ctx, input.Email)
	if existing != nil {
		return nil, errors.New("email is already registered")
	}

	// Validate role exists
	role, err := s.roleRepo.GetByID(ctx, input.RoleID)
	if err != nil {
		return nil, errors.New("invalid role")
	}

	// Validate branch exists if provided
	if input.BranchID != nil {
		_, err := s.branchRepo.GetByID(ctx, *input.BranchID)
		if err != nil {
			return nil, errors.New("invalid branch")
		}
	}

	// Validate password strength
	if err := s.passwordManager.ValidatePasswordStrength(input.Password); err != nil {
		return nil, err
	}

	// Hash password
	passwordHash, err := s.passwordManager.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Email:        input.Email,
		PasswordHash: passwordHash,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Phone:        input.Phone,
		RoleID:       input.RoleID,
		BranchID:     input.BranchID,
		IsActive:     input.IsActive,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.Role = role

	return user.ToPublic(), nil
}

// UpdateUserInput represents update user request data
type UpdateUserInput struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2"`
	LastName  string `json:"last_name" validate:"omitempty,min=2"`
	Phone     string `json:"phone"`
	RoleID    *int64 `json:"role_id"`
	BranchID  *int64 `json:"branch_id"`
	IsActive  *bool  `json:"is_active"`
}

// Update updates an existing user
func (s *UserService) Update(ctx context.Context, id int64, input UpdateUserInput) (*domain.UserPublic, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}
	if input.RoleID != nil {
		// Validate role exists
		role, err := s.roleRepo.GetByID(ctx, *input.RoleID)
		if err != nil {
			return nil, errors.New("invalid role")
		}
		user.RoleID = *input.RoleID
		user.Role = role
	}
	if input.BranchID != nil {
		// Validate branch exists
		_, err := s.branchRepo.GetByID(ctx, *input.BranchID)
		if err != nil {
			return nil, errors.New("invalid branch")
		}
		user.BranchID = input.BranchID
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	// Save changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Load role if not loaded
	if user.Role == nil {
		user.Role, _ = s.roleRepo.GetByID(ctx, user.RoleID)
	}

	return user.ToPublic(), nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id int64) (*domain.UserPublic, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Load role
	user.Role, _ = s.roleRepo.GetByID(ctx, user.RoleID)

	// Load branch if applicable
	if user.BranchID != nil {
		user.Branch, _ = s.branchRepo.GetByID(ctx, *user.BranchID)
	}

	return user.ToPublic(), nil
}

// List retrieves users with pagination and filters
func (s *UserService) List(ctx context.Context, params repository.UserListParams) (*repository.PaginatedResult[domain.UserPublic], error) {
	result, err := s.userRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	// Convert to public format
	publicUsers := make([]domain.UserPublic, len(result.Data))
	for i, user := range result.Data {
		user.Role, _ = s.roleRepo.GetByID(ctx, user.RoleID)
		publicUsers[i] = *user.ToPublic()
	}

	return &repository.PaginatedResult[domain.UserPublic]{
		Data:       publicUsers,
		Total:      result.Total,
		Page:       result.Page,
		PerPage:    result.PerPage,
		TotalPages: result.TotalPages,
	}, nil
}

// Delete soft deletes a user
func (s *UserService) Delete(ctx context.Context, id int64) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(ctx, id)
}

// ResetPassword resets a user's password (admin action)
func (s *UserService) ResetPassword(ctx context.Context, id int64, newPassword string) error {
	// Validate password strength
	if err := s.passwordManager.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	// Hash password
	passwordHash, err := s.passwordManager.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.userRepo.UpdatePassword(ctx, id, passwordHash)
}
