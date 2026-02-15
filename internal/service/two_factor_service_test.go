package service

import (
	"context"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository/mocks"
	"pawnshop/pkg/auth"
)

func setupTwoFactorService() (TwoFactorService, *mocks.MockTwoFactorRepository, *mocks.MockUserRepository) {
	twoFactorRepo := new(mocks.MockTwoFactorRepository)
	userRepo := new(mocks.MockUserRepository)
	passwordManager := auth.NewPasswordManager()
	service := NewTwoFactorService(twoFactorRepo, userRepo, passwordManager, "TestApp")
	return service, twoFactorRepo, userRepo
}

func TestTwoFactorService_Setup_Success(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	twoFactorRepo.On("Is2FAEnabled", ctx, int64(1)).Return(false, nil)
	twoFactorRepo.On("Enable2FA", ctx, int64(1), mock.AnythingOfType("string")).Return(nil)
	twoFactorRepo.On("CreateBackupCodes", ctx, int64(1), mock.AnythingOfType("[]string")).Return(nil)

	result, err := service.Setup(ctx, 1, "test@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Secret)
	assert.NotEmpty(t, result.QRCodeURL)
	assert.Len(t, result.BackupCodes, BackupCodesCount)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_Setup_AlreadyEnabled(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	twoFactorRepo.On("Is2FAEnabled", ctx, int64(1)).Return(true, nil)

	result, err := service.Setup(ctx, 1, "test@example.com")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTwoFactorAlreadyEnabled, err)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_Enable_Success(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	// Generate a real secret for testing
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "TestApp",
		AccountName: "test@example.com",
	})
	secret := key.Secret()
	validCode, _ := totp.GenerateCode(secret, time.Now())

	twoFactorRepo.On("Get2FASecret", ctx, int64(1)).Return(secret, nil)
	twoFactorRepo.On("Confirm2FA", ctx, int64(1)).Return(nil)

	err := service.Enable(ctx, 1, validCode)

	assert.NoError(t, err)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_Enable_NotSetup(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	twoFactorRepo.On("Get2FASecret", ctx, int64(1)).Return("", nil)

	err := service.Enable(ctx, 1, "123456")

	assert.Error(t, err)
	assert.Equal(t, ErrTwoFactorNotSetup, err)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_Enable_InvalidCode(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	// Generate a real secret
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "TestApp",
		AccountName: "test@example.com",
	})
	secret := key.Secret()

	twoFactorRepo.On("Get2FASecret", ctx, int64(1)).Return(secret, nil)

	err := service.Enable(ctx, 1, "000000") // Invalid code

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTOTPCode, err)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_Disable_Success(t *testing.T) {
	service, twoFactorRepo, userRepo := setupTwoFactorService()
	ctx := context.Background()

	passwordManager := auth.NewPasswordManager()
	passwordHash, _ := passwordManager.HashPassword("password123")

	user := &domain.User{
		ID:           1,
		PasswordHash: passwordHash,
	}

	userRepo.On("GetByID", ctx, int64(1)).Return(user, nil)
	twoFactorRepo.On("Disable2FA", ctx, int64(1)).Return(nil)

	err := service.Disable(ctx, 1, "password123")

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_Disable_InvalidPassword(t *testing.T) {
	service, _, userRepo := setupTwoFactorService()
	ctx := context.Background()

	passwordManager := auth.NewPasswordManager()
	passwordHash, _ := passwordManager.HashPassword("password123")

	user := &domain.User{
		ID:           1,
		PasswordHash: passwordHash,
	}

	userRepo.On("GetByID", ctx, int64(1)).Return(user, nil)

	err := service.Disable(ctx, 1, "wrongpassword")

	assert.Error(t, err)
	assert.Equal(t, "invalid password", err.Error())
	userRepo.AssertExpectations(t)
}

func TestTwoFactorService_Disable_UserNotFound(t *testing.T) {
	service, _, userRepo := setupTwoFactorService()
	ctx := context.Background()

	userRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	err := service.Disable(ctx, 999, "password123")

	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	userRepo.AssertExpectations(t)
}

func TestTwoFactorService_GetStatus_Enabled(t *testing.T) {
	service, twoFactorRepo, userRepo := setupTwoFactorService()
	ctx := context.Background()

	confirmedAt := time.Now()
	user := &domain.User{
		ID:                    1,
		TwoFactorConfirmedAt: &confirmedAt,
	}

	twoFactorRepo.On("Is2FAEnabled", ctx, int64(1)).Return(true, nil)
	userRepo.On("GetByID", ctx, int64(1)).Return(user, nil)

	status, err := service.GetStatus(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.True(t, status.Enabled)
	assert.NotNil(t, status.ConfirmedAt)
	twoFactorRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestTwoFactorService_GetStatus_NotEnabled(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	twoFactorRepo.On("Is2FAEnabled", ctx, int64(1)).Return(false, nil)

	status, err := service.GetStatus(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.False(t, status.Enabled)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_CreateChallenge_Success(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	twoFactorRepo.On("CreateChallenge", ctx, mock.AnythingOfType("*domain.TwoFactorChallenge")).Return(nil)

	challenge, err := service.CreateChallenge(ctx, 1, "192.168.1.1", "Mozilla/5.0")

	assert.NoError(t, err)
	assert.NotNil(t, challenge)
	assert.Equal(t, int64(1), challenge.UserID)
	assert.NotEmpty(t, challenge.ChallengeToken)
	assert.Equal(t, "192.168.1.1", challenge.IPAddress)
	assert.Equal(t, "Mozilla/5.0", challenge.UserAgent)
	assert.False(t, challenge.IsExpired())
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_VerifyChallenge_Success(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	// Generate a real secret
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "TestApp",
		AccountName: "test@example.com",
	})
	secret := key.Secret()
	validCode, _ := totp.GenerateCode(secret, time.Now())

	challenge := &domain.TwoFactorChallenge{
		ID:             1,
		UserID:         1,
		ChallengeToken: "test-token",
		ExpiresAt:      time.Now().Add(5 * time.Minute),
	}

	twoFactorRepo.On("GetChallengeByToken", ctx, "test-token").Return(challenge, nil)
	twoFactorRepo.On("Get2FASecret", ctx, int64(1)).Return(secret, nil)
	twoFactorRepo.On("MarkChallengeVerified", ctx, int64(1)).Return(nil)

	result, err := service.VerifyChallenge(ctx, "test-token", validCode)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.VerifiedAt)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_VerifyChallenge_NotFound(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	twoFactorRepo.On("GetChallengeByToken", ctx, "invalid-token").Return(nil, nil)

	result, err := service.VerifyChallenge(ctx, "invalid-token", "123456")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrChallengeNotFound, err)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_VerifyChallenge_Expired(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	challenge := &domain.TwoFactorChallenge{
		ID:             1,
		UserID:         1,
		ChallengeToken: "test-token",
		ExpiresAt:      time.Now().Add(-1 * time.Minute), // Expired
	}

	twoFactorRepo.On("GetChallengeByToken", ctx, "test-token").Return(challenge, nil)

	result, err := service.VerifyChallenge(ctx, "test-token", "123456")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrChallengeExpired, err)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_VerifyChallenge_AlreadyVerified(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	verifiedAt := time.Now()
	challenge := &domain.TwoFactorChallenge{
		ID:             1,
		UserID:         1,
		ChallengeToken: "test-token",
		ExpiresAt:      time.Now().Add(5 * time.Minute),
		VerifiedAt:     &verifiedAt, // Already verified
	}

	twoFactorRepo.On("GetChallengeByToken", ctx, "test-token").Return(challenge, nil)

	result, err := service.VerifyChallenge(ctx, "test-token", "123456")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrChallengeAlreadyVerified, err)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_VerifyChallenge_InvalidCode(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	// Generate a real secret
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "TestApp",
		AccountName: "test@example.com",
	})
	secret := key.Secret()

	challenge := &domain.TwoFactorChallenge{
		ID:             1,
		UserID:         1,
		ChallengeToken: "test-token",
		ExpiresAt:      time.Now().Add(5 * time.Minute),
	}

	twoFactorRepo.On("GetChallengeByToken", ctx, "test-token").Return(challenge, nil)
	twoFactorRepo.On("Get2FASecret", ctx, int64(1)).Return(secret, nil)

	result, err := service.VerifyChallenge(ctx, "test-token", "000000")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrInvalidTOTPCode, err)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_RegenerateBackupCodes_Success(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	twoFactorRepo.On("Is2FAEnabled", ctx, int64(1)).Return(true, nil)
	twoFactorRepo.On("CreateBackupCodes", ctx, int64(1), mock.AnythingOfType("[]string")).Return(nil)

	codes, err := service.RegenerateBackupCodes(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, codes, BackupCodesCount)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_RegenerateBackupCodes_NotEnabled(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	twoFactorRepo.On("Is2FAEnabled", ctx, int64(1)).Return(false, nil)

	codes, err := service.RegenerateBackupCodes(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, codes)
	assert.Equal(t, ErrTwoFactorNotEnabled, err)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_GetBackupCodesCount_Success(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	twoFactorRepo.On("GetUnusedBackupCodesCount", ctx, int64(1)).Return(8, nil)

	count, err := service.GetBackupCodesCount(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 8, count)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_ValidateTOTP_Success(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	// Generate a real secret
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "TestApp",
		AccountName: "test@example.com",
	})
	secret := key.Secret()
	validCode, _ := totp.GenerateCode(secret, time.Now())

	twoFactorRepo.On("Get2FASecret", ctx, int64(1)).Return(secret, nil)

	valid, err := service.ValidateTOTP(ctx, 1, validCode)

	assert.NoError(t, err)
	assert.True(t, valid)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_ValidateTOTP_InvalidCode(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	// Generate a real secret
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "TestApp",
		AccountName: "test@example.com",
	})
	secret := key.Secret()

	twoFactorRepo.On("Get2FASecret", ctx, int64(1)).Return(secret, nil)

	valid, err := service.ValidateTOTP(ctx, 1, "000000")

	assert.NoError(t, err)
	assert.False(t, valid)
	twoFactorRepo.AssertExpectations(t)
}

func TestTwoFactorService_ValidateTOTP_NotSetup(t *testing.T) {
	service, twoFactorRepo, _ := setupTwoFactorService()
	ctx := context.Background()

	twoFactorRepo.On("Get2FASecret", ctx, int64(1)).Return("", nil)

	valid, err := service.ValidateTOTP(ctx, 1, "123456")

	assert.Error(t, err)
	assert.False(t, valid)
	assert.Equal(t, ErrTwoFactorNotSetup, err)
	twoFactorRepo.AssertExpectations(t)
}
