package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository/mocks"
	"pawnshop/pkg/auth"
)

func setupAuthService() (*AuthService, *mocks.MockUserRepository, *mocks.MockRoleRepository, *mocks.MockRefreshTokenRepository) {
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	refreshTokenRepo := new(mocks.MockRefreshTokenRepository)

	jwtManager := auth.NewJWTManager(auth.JWTConfig{
		Secret:          "test-secret-key-for-testing-only",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 168 * time.Hour,
		Issuer:          "pawnshop-test",
	})
	passwordManager := auth.NewPasswordManager()

	service := NewAuthService(userRepo, roleRepo, refreshTokenRepo, jwtManager, passwordManager)

	return service, userRepo, roleRepo, refreshTokenRepo
}

// --- Login tests ---

func TestAuthService_Login_Success(t *testing.T) {
	service, userRepo, roleRepo, refreshTokenRepo := setupAuthService()
	ctx := context.Background()

	pm := auth.NewPasswordManager()
	passwordHash, _ := pm.HashPassword("testpassword123")

	branchID := int64(1)
	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		FirstName:    "Test",
		LastName:     "User",
		RoleID:       1,
		BranchID:     &branchID,
		IsActive:     true,
	}

	role := &domain.Role{
		ID:          1,
		Name:        "admin",
		DisplayName: "Administrator",
		Permissions: []byte(`["*"]`),
	}

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
	roleRepo.On("GetByID", ctx, int64(1)).Return(role, nil)
	refreshTokenRepo.On("Create", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)
	userRepo.On("UpdateLastLogin", ctx, int64(1), "127.0.0.1").Return(nil)

	input := LoginInput{
		Email:    "test@example.com",
		Password: "testpassword123",
	}
	result, err := service.Login(ctx, input, "127.0.0.1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.User.Email)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, "Bearer", result.TokenType)

	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
	refreshTokenRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidEmail(t *testing.T) {
	service, userRepo, _, _ := setupAuthService()
	ctx := context.Background()

	userRepo.On("GetByEmail", ctx, "invalid@example.com").Return(nil, errors.New("not found"))

	input := LoginInput{
		Email:    "invalid@example.com",
		Password: "password123",
	}
	result, err := service.Login(ctx, input, "127.0.0.1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid credentials", err.Error())

	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	service, userRepo, _, _ := setupAuthService()
	ctx := context.Background()

	pm := auth.NewPasswordManager()
	passwordHash, _ := pm.HashPassword("correctpassword")

	branchID := int64(1)
	user := &domain.User{
		ID:                  1,
		Email:               "test@example.com",
		PasswordHash:        passwordHash,
		RoleID:              1,
		BranchID:            &branchID,
		IsActive:            true,
		FailedLoginAttempts: 0,
	}

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
	userRepo.On("IncrementFailedLogins", ctx, int64(1)).Return(nil)

	input := LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	result, err := service.Login(ctx, input, "127.0.0.1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid credentials", err.Error())

	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_AccountLocked(t *testing.T) {
	service, userRepo, _, _ := setupAuthService()
	ctx := context.Background()

	lockTime := time.Now().Add(15 * time.Minute)
	branchID := int64(1)
	user := &domain.User{
		ID:          1,
		Email:       "test@example.com",
		RoleID:      1,
		BranchID:    &branchID,
		IsActive:    true,
		LockedUntil: &lockTime,
	}

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)

	input := LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	}
	result, err := service.Login(ctx, input, "127.0.0.1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "account is locked", err.Error())

	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_AccountInactive(t *testing.T) {
	service, userRepo, _, _ := setupAuthService()
	ctx := context.Background()

	branchID := int64(1)
	user := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		RoleID:   1,
		BranchID: &branchID,
		IsActive: false,
	}

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)

	input := LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	}
	result, err := service.Login(ctx, input, "127.0.0.1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "account is inactive", err.Error())

	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_LockAccountAfterFailedAttempts(t *testing.T) {
	service, userRepo, _, _ := setupAuthService()
	ctx := context.Background()

	pm := auth.NewPasswordManager()
	passwordHash, _ := pm.HashPassword("correctpassword")

	branchID := int64(1)
	user := &domain.User{
		ID:                  1,
		Email:               "test@example.com",
		PasswordHash:        passwordHash,
		RoleID:              1,
		BranchID:            &branchID,
		IsActive:            true,
		FailedLoginAttempts: 4, // One more failure locks
	}

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
	userRepo.On("IncrementFailedLogins", ctx, int64(1)).Return(nil)
	userRepo.On("LockUser", ctx, int64(1), mock.AnythingOfType("*int64")).Return(nil)

	input := LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	result, err := service.Login(ctx, input, "127.0.0.1")

	assert.Error(t, err)
	assert.Nil(t, result)
	userRepo.AssertExpectations(t)
}

// --- Refresh tests ---

func TestAuthService_Refresh_Success(t *testing.T) {
	service, userRepo, roleRepo, refreshTokenRepo := setupAuthService()
	ctx := context.Background()

	// First login to get a valid refresh token
	pm := auth.NewPasswordManager()
	passwordHash, _ := pm.HashPassword("testpassword123")

	branchID := int64(1)
	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		FirstName:    "Test",
		LastName:     "User",
		RoleID:       1,
		BranchID:     &branchID,
		IsActive:     true,
	}

	role := &domain.Role{
		ID:          1,
		Name:        "admin",
		Permissions: []byte(`["*"]`),
	}

	// Login to get valid tokens
	userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
	roleRepo.On("GetByID", ctx, int64(1)).Return(role, nil)
	refreshTokenRepo.On("Create", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)
	userRepo.On("UpdateLastLogin", ctx, int64(1), "127.0.0.1").Return(nil)

	loginInput := LoginInput{Email: "test@example.com", Password: "testpassword123"}
	loginResult, err := service.Login(ctx, loginInput, "127.0.0.1")
	assert.NoError(t, err)

	// Now test refresh
	storedToken := &domain.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenHash: auth.HashToken(loginResult.RefreshToken),
		IPAddress: "127.0.0.1",
		ExpiresAt: time.Now().Add(168 * time.Hour),
	}

	refreshTokenRepo.On("GetByHash", ctx, mock.AnythingOfType("string")).Return(storedToken, nil)
	userRepo.On("GetByID", ctx, int64(1)).Return(user, nil)
	refreshTokenRepo.On("Revoke", ctx, int64(1)).Return(nil)

	refreshInput := RefreshInput{RefreshToken: loginResult.RefreshToken}
	refreshResult, err := service.Refresh(ctx, refreshInput)

	assert.NoError(t, err)
	assert.NotNil(t, refreshResult)
	assert.NotEmpty(t, refreshResult.AccessToken)
	assert.NotEmpty(t, refreshResult.RefreshToken)
	assert.Equal(t, "Bearer", refreshResult.TokenType)
}

func TestAuthService_Refresh_InvalidToken(t *testing.T) {
	service, _, _, _ := setupAuthService()
	ctx := context.Background()

	input := RefreshInput{RefreshToken: "invalid-token"}

	result, err := service.Refresh(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid refresh token", err.Error())
}

func TestAuthService_Refresh_TokenNotInDatabase(t *testing.T) {
	service, userRepo, roleRepo, refreshTokenRepo := setupAuthService()
	ctx := context.Background()

	// Login to get a valid JWT refresh token
	pm := auth.NewPasswordManager()
	passwordHash, _ := pm.HashPassword("testpassword123")

	branchID := int64(1)
	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		FirstName:    "Test",
		LastName:     "User",
		RoleID:       1,
		BranchID:     &branchID,
		IsActive:     true,
	}

	role := &domain.Role{
		ID:          1,
		Name:        "admin",
		Permissions: []byte(`["*"]`),
	}

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
	roleRepo.On("GetByID", ctx, int64(1)).Return(role, nil)
	refreshTokenRepo.On("Create", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)
	userRepo.On("UpdateLastLogin", ctx, int64(1), "127.0.0.1").Return(nil)

	loginInput := LoginInput{Email: "test@example.com", Password: "testpassword123"}
	loginResult, _ := service.Login(ctx, loginInput, "127.0.0.1")

	// Token not found in database
	refreshTokenRepo.On("GetByHash", ctx, mock.AnythingOfType("string")).Return(nil, errors.New("not found"))

	input := RefreshInput{RefreshToken: loginResult.RefreshToken}
	result, err := service.Refresh(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "refresh token not found", err.Error())
}

func TestAuthService_Refresh_TokenRevoked(t *testing.T) {
	service, userRepo, roleRepo, refreshTokenRepo := setupAuthService()
	ctx := context.Background()

	pm := auth.NewPasswordManager()
	passwordHash, _ := pm.HashPassword("testpassword123")

	branchID := int64(1)
	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		FirstName:    "Test",
		LastName:     "User",
		RoleID:       1,
		BranchID:     &branchID,
		IsActive:     true,
	}

	role := &domain.Role{
		ID:          1,
		Name:        "admin",
		Permissions: []byte(`["*"]`),
	}

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
	roleRepo.On("GetByID", ctx, int64(1)).Return(role, nil)
	refreshTokenRepo.On("Create", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)
	userRepo.On("UpdateLastLogin", ctx, int64(1), "127.0.0.1").Return(nil)

	loginInput := LoginInput{Email: "test@example.com", Password: "testpassword123"}
	loginResult, _ := service.Login(ctx, loginInput, "127.0.0.1")

	// Token is revoked (already used)
	now := time.Now()
	storedToken := &domain.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenHash: auth.HashToken(loginResult.RefreshToken),
		RevokedAt: &now, // Already revoked
		ExpiresAt: time.Now().Add(168 * time.Hour),
	}

	refreshTokenRepo.On("GetByHash", ctx, mock.AnythingOfType("string")).Return(storedToken, nil)

	input := RefreshInput{RefreshToken: loginResult.RefreshToken}
	result, err := service.Refresh(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "refresh token is invalid or expired", err.Error())
}

func TestAuthService_Refresh_UserInactive(t *testing.T) {
	service, userRepo, roleRepo, refreshTokenRepo := setupAuthService()
	ctx := context.Background()

	pm := auth.NewPasswordManager()
	passwordHash, _ := pm.HashPassword("testpassword123")

	branchID := int64(1)
	activeUser := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		FirstName:    "Test",
		LastName:     "User",
		RoleID:       1,
		BranchID:     &branchID,
		IsActive:     true,
	}

	role := &domain.Role{
		ID:          1,
		Name:        "admin",
		Permissions: []byte(`["*"]`),
	}

	userRepo.On("GetByEmail", ctx, "test@example.com").Return(activeUser, nil)
	roleRepo.On("GetByID", ctx, int64(1)).Return(role, nil)
	refreshTokenRepo.On("Create", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)
	userRepo.On("UpdateLastLogin", ctx, int64(1), "127.0.0.1").Return(nil)

	loginInput := LoginInput{Email: "test@example.com", Password: "testpassword123"}
	loginResult, _ := service.Login(ctx, loginInput, "127.0.0.1")

	storedToken := &domain.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenHash: auth.HashToken(loginResult.RefreshToken),
		ExpiresAt: time.Now().Add(168 * time.Hour),
	}

	inactiveUser := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		IsActive: false,
	}

	refreshTokenRepo.On("GetByHash", ctx, mock.AnythingOfType("string")).Return(storedToken, nil)
	userRepo.On("GetByID", ctx, int64(1)).Return(inactiveUser, nil)

	input := RefreshInput{RefreshToken: loginResult.RefreshToken}
	result, err := service.Refresh(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "account is inactive or locked", err.Error())
}

// --- Logout tests ---

func TestAuthService_Logout_Success(t *testing.T) {
	service, _, _, refreshTokenRepo := setupAuthService()
	ctx := context.Background()

	refreshTokenRepo.On("RevokeAllForUser", ctx, int64(1)).Return(nil)

	err := service.Logout(ctx, 1)

	assert.NoError(t, err)
	refreshTokenRepo.AssertExpectations(t)
}

// --- ChangePassword tests ---

func TestAuthService_ChangePassword_Success(t *testing.T) {
	service, userRepo, _, refreshTokenRepo := setupAuthService()
	ctx := context.Background()

	pm := auth.NewPasswordManager()
	passwordHash, _ := pm.HashPassword("oldpassword123")

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	userRepo.On("GetByID", ctx, int64(1)).Return(user, nil)
	userRepo.On("UpdatePassword", ctx, int64(1), mock.AnythingOfType("string")).Return(nil)
	refreshTokenRepo.On("RevokeAllForUser", ctx, int64(1)).Return(nil)

	input := ChangePasswordInput{
		CurrentPassword: "oldpassword123",
		NewPassword:     "NewPassword123!",
	}
	err := service.ChangePassword(ctx, 1, input)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	refreshTokenRepo.AssertExpectations(t)
}

func TestAuthService_ChangePassword_WrongCurrentPassword(t *testing.T) {
	service, userRepo, _, _ := setupAuthService()
	ctx := context.Background()

	pm := auth.NewPasswordManager()
	passwordHash, _ := pm.HashPassword("correctpassword")

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	userRepo.On("GetByID", ctx, int64(1)).Return(user, nil)

	input := ChangePasswordInput{
		CurrentPassword: "wrongpassword",
		NewPassword:     "NewPassword123!",
	}
	err := service.ChangePassword(ctx, 1, input)

	assert.Error(t, err)
	assert.Equal(t, "current password is incorrect", err.Error())
	userRepo.AssertExpectations(t)
}

func TestAuthService_ChangePassword_WeakNewPassword(t *testing.T) {
	service, userRepo, _, _ := setupAuthService()
	ctx := context.Background()

	pm := auth.NewPasswordManager()
	passwordHash, _ := pm.HashPassword("oldpassword123")

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	userRepo.On("GetByID", ctx, int64(1)).Return(user, nil)

	input := ChangePasswordInput{
		CurrentPassword: "oldpassword123",
		NewPassword:     "weakpassword",
	}
	err := service.ChangePassword(ctx, 1, input)

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestAuthService_ChangePassword_UserNotFound(t *testing.T) {
	service, userRepo, _, _ := setupAuthService()
	ctx := context.Background()

	userRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := ChangePasswordInput{
		CurrentPassword: "old",
		NewPassword:     "new",
	}
	err := service.ChangePassword(ctx, 999, input)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}
