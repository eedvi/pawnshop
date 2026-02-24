package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/pkg/auth"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo         repository.UserRepository
	roleRepo         repository.RoleRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtManager       *auth.JWTManager
	passwordManager  *auth.PasswordManager
	logger           zerolog.Logger
}

// NewAuthService creates a new AuthService
func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtManager *auth.JWTManager,
	passwordManager *auth.PasswordManager,
	logger zerolog.Logger,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		roleRepo:         roleRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
		passwordManager:  passwordManager,
		logger:           logger.With().Str("service", "auth").Logger(),
	}
}

// LoginInput represents login request data
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginOutput represents login response data
type LoginOutput struct {
	User         *domain.UserPublic `json:"user"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	ExpiresAt    time.Time          `json:"expires_at"`
	TokenType    string             `json:"token_type"`
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, input LoginInput, ip string) (*LoginOutput, error) {
	s.logger.Info().
		Str("email", input.Email).
		Str("ip", ip).
		Msg("Login attempt")

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		s.logger.Warn().
			Str("email", input.Email).
			Str("ip", ip).
			Msg("Login failed: user not found")
		return nil, errors.New("invalid credentials")
	}

	// Check if user can login
	if !user.CanLogin() {
		if user.IsLocked() {
			s.logger.Warn().
				Int64("user_id", user.ID).
				Str("email", input.Email).
				Str("ip", ip).
				Msg("Login failed: account locked")
			return nil, errors.New("account is locked")
		}
		s.logger.Warn().
			Int64("user_id", user.ID).
			Str("email", input.Email).
			Str("ip", ip).
			Msg("Login failed: account inactive")
		return nil, errors.New("account is inactive")
	}

	// Verify password (NEVER log the password)
	valid, err := s.passwordManager.VerifyPassword(input.Password, user.PasswordHash)
	if err != nil || !valid {
		// Increment failed login attempts
		s.userRepo.IncrementFailedLogins(ctx, user.ID)

		// Lock account after 5 failed attempts
		if user.FailedLoginAttempts >= 4 {
			lockDuration := int64(15) // 15 minutes
			s.userRepo.LockUser(ctx, user.ID, &lockDuration)
			s.logger.Warn().
				Int64("user_id", user.ID).
				Str("email", input.Email).
				Str("ip", ip).
				Int("failed_attempts", user.FailedLoginAttempts+1).
				Msg("Account locked due to repeated failed login attempts")
		} else {
			s.logger.Warn().
				Int64("user_id", user.ID).
				Str("email", input.Email).
				Str("ip", ip).
				Int("failed_attempts", user.FailedLoginAttempts+1).
				Msg("Login failed: invalid password")
		}

		return nil, errors.New("invalid credentials")
	}

	// Load role
	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Int64("user_id", user.ID).
			Int64("role_id", user.RoleID).
			Msg("Failed to load user role")
		return nil, fmt.Errorf("failed to load user role: %w", err)
	}
	user.Role = role

	// Get permissions
	permissions, _ := role.GetPermissions()

	// Generate token pair
	claims := auth.JWTClaims{
		UserID:      user.ID,
		Email:       user.Email,
		RoleID:      user.RoleID,
		BranchID:    user.BranchID,
		Permissions: permissions,
	}

	tokenPair, err := s.jwtManager.GenerateTokenPair(claims)
	if err != nil {
		s.logger.Error().
			Err(err).
			Int64("user_id", user.ID).
			Msg("Failed to generate JWT tokens")
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token
	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: auth.HashToken(tokenPair.RefreshToken),
		IPAddress: ip,
		ExpiresAt: tokenPair.ExpiresAt,
	}
	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		s.logger.Error().
			Err(err).
			Int64("user_id", user.ID).
			Msg("Failed to store refresh token")
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Update last login
	s.userRepo.UpdateLastLogin(ctx, user.ID, ip)

	s.logger.Info().
		Int64("user_id", user.ID).
		Str("email", user.Email).
		Str("role", role.Name).
		Str("ip", ip).
		Msg("Login successful")

	return &LoginOutput{
		User:         user.ToPublic(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		TokenType:    tokenPair.TokenType,
	}, nil
}

// RefreshInput represents refresh token request data
type RefreshInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Refresh generates new tokens using a refresh token
func (s *AuthService) Refresh(ctx context.Context, input RefreshInput) (*LoginOutput, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if token exists in database
	tokenHash := auth.HashToken(input.RefreshToken)
	storedToken, err := s.refreshTokenRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("refresh token not found")
	}

	// Check if token is valid
	if !storedToken.IsValid() {
		return nil, errors.New("refresh token is invalid or expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user can login
	if !user.CanLogin() {
		return nil, errors.New("account is inactive or locked")
	}

	// Load role
	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user role: %w", err)
	}
	user.Role = role

	// Revoke old token
	s.refreshTokenRepo.Revoke(ctx, storedToken.ID)

	// Get permissions
	permissions, _ := role.GetPermissions()

	// Generate new token pair
	newClaims := auth.JWTClaims{
		UserID:      user.ID,
		Email:       user.Email,
		RoleID:      user.RoleID,
		BranchID:    user.BranchID,
		Permissions: permissions,
	}

	tokenPair, err := s.jwtManager.GenerateTokenPair(newClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store new refresh token
	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: auth.HashToken(tokenPair.RefreshToken),
		IPAddress: storedToken.IPAddress,
		ExpiresAt: tokenPair.ExpiresAt,
	}
	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &LoginOutput{
		User:         user.ToPublic(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		TokenType:    tokenPair.TokenType,
	}, nil
}

// Logout invalidates all refresh tokens for a user
func (s *AuthService) Logout(ctx context.Context, userID int64) error {
	s.logger.Info().Int64("user_id", userID).Msg("Logout initiated")

	err := s.refreshTokenRepo.RevokeAllForUser(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("Failed to revoke refresh tokens during logout")
		return err
	}

	s.logger.Info().Int64("user_id", userID).Msg("Logout successful")
	return nil
}

// ChangePasswordInput represents change password request data
type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID int64, input ChangePasswordInput) error {
	s.logger.Info().Int64("user_id", userID).Msg("Password change requested")

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("User not found during password change")
		return errors.New("user not found")
	}

	// Verify current password (NEVER log passwords)
	valid, err := s.passwordManager.VerifyPassword(input.CurrentPassword, user.PasswordHash)
	if err != nil || !valid {
		s.logger.Warn().Int64("user_id", userID).Msg("Password change failed: current password incorrect")
		return errors.New("current password is incorrect")
	}

	// Validate new password strength
	if err := s.passwordManager.ValidatePasswordStrength(input.NewPassword); err != nil {
		s.logger.Warn().Int64("user_id", userID).Err(err).Msg("Password change failed: weak password")
		return err
	}

	// Hash new password
	newHash, err := s.passwordManager.HashPassword(input.NewPassword)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("Failed to hash new password")
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, userID, newHash); err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("Failed to update password in database")
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all refresh tokens
	if err := s.refreshTokenRepo.RevokeAllForUser(ctx, userID); err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("Failed to revoke refresh tokens after password change")
	}

	s.logger.Info().Int64("user_id", userID).Msg("Password changed successfully")
	return nil
}
