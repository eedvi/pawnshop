package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"strings"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/pkg/auth"
	"github.com/pquerna/otp/totp"
)

var (
	ErrTwoFactorAlreadyEnabled  = errors.New("2FA is already enabled")
	ErrTwoFactorNotEnabled      = errors.New("2FA is not enabled")
	ErrTwoFactorNotSetup        = errors.New("2FA is not set up")
	ErrInvalidTOTPCode          = errors.New("invalid TOTP code")
	ErrInvalidBackupCode        = errors.New("invalid backup code")
	ErrChallengeNotFound        = errors.New("challenge not found")
	ErrChallengeExpired         = errors.New("challenge expired")
	ErrChallengeAlreadyVerified = errors.New("challenge already verified")
)

const (
	// ChallengeExpiration is the duration for which a 2FA challenge is valid
	ChallengeExpiration = 5 * time.Minute
	// BackupCodesCount is the number of backup codes generated
	BackupCodesCount = 10
)

// TwoFactorService defines the interface for 2FA operations
type TwoFactorService interface {
	// Setup generates a new 2FA secret and returns setup data
	Setup(ctx context.Context, userID int64, email string) (*domain.TwoFactorSetup, error)

	// Enable enables 2FA for a user after verifying the TOTP code
	Enable(ctx context.Context, userID int64, code string) error

	// Disable disables 2FA for a user
	Disable(ctx context.Context, userID int64, password string) error

	// GetStatus returns the 2FA status for a user
	GetStatus(ctx context.Context, userID int64) (*domain.TwoFactorStatus, error)

	// CreateChallenge creates a 2FA challenge for login
	CreateChallenge(ctx context.Context, userID int64, ipAddress, userAgent string) (*domain.TwoFactorChallenge, error)

	// VerifyChallenge verifies a 2FA challenge with TOTP code
	VerifyChallenge(ctx context.Context, token, code string) (*domain.TwoFactorChallenge, error)

	// VerifyChallengeWithBackup verifies a 2FA challenge with backup code
	VerifyChallengeWithBackup(ctx context.Context, token, backupCode string) (*domain.TwoFactorChallenge, error)

	// RegenerateBackupCodes regenerates backup codes for a user
	RegenerateBackupCodes(ctx context.Context, userID int64) ([]string, error)

	// GetBackupCodesCount returns the number of unused backup codes
	GetBackupCodesCount(ctx context.Context, userID int64) (int, error)

	// ValidateTOTP validates a TOTP code for a user
	ValidateTOTP(ctx context.Context, userID int64, code string) (bool, error)
}

type twoFactorService struct {
	twoFactorRepo   repository.TwoFactorRepository
	userRepo        repository.UserRepository
	passwordManager *auth.PasswordManager
	issuer          string
}

// NewTwoFactorService creates a new 2FA service
func NewTwoFactorService(
	twoFactorRepo repository.TwoFactorRepository,
	userRepo repository.UserRepository,
	passwordManager *auth.PasswordManager,
	issuer string,
) TwoFactorService {
	return &twoFactorService{
		twoFactorRepo:   twoFactorRepo,
		userRepo:        userRepo,
		passwordManager: passwordManager,
		issuer:          issuer,
	}
}

func (s *twoFactorService) Setup(ctx context.Context, userID int64, email string) (*domain.TwoFactorSetup, error) {
	// Check if 2FA is already enabled
	enabled, err := s.twoFactorRepo.Is2FAEnabled(ctx, userID)
	if err != nil {
		return nil, err
	}
	if enabled {
		return nil, ErrTwoFactorAlreadyEnabled
	}

	// Generate new secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// Store secret (not yet enabled)
	if err := s.twoFactorRepo.Enable2FA(ctx, userID, key.Secret()); err != nil {
		return nil, err
	}

	// Generate backup codes
	backupCodes, codeHashes, err := s.generateBackupCodes()
	if err != nil {
		return nil, err
	}

	// Store backup code hashes
	if err := s.twoFactorRepo.CreateBackupCodes(ctx, userID, codeHashes); err != nil {
		return nil, err
	}

	return &domain.TwoFactorSetup{
		Secret:      key.Secret(),
		QRCodeURL:   key.URL(),
		BackupCodes: backupCodes,
	}, nil
}

func (s *twoFactorService) Enable(ctx context.Context, userID int64, code string) error {
	// Get the secret
	secret, err := s.twoFactorRepo.Get2FASecret(ctx, userID)
	if err != nil {
		return err
	}
	if secret == "" {
		return ErrTwoFactorNotSetup
	}

	// Verify the code
	if !totp.Validate(code, secret) {
		return ErrInvalidTOTPCode
	}

	// Confirm 2FA
	return s.twoFactorRepo.Confirm2FA(ctx, userID)
}

func (s *twoFactorService) Disable(ctx context.Context, userID int64, password string) error {
	// Get user to verify password
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Verify password
	valid, err := s.passwordManager.VerifyPassword(password, user.PasswordHash)
	if err != nil || !valid {
		return errors.New("invalid password")
	}

	// Disable 2FA
	return s.twoFactorRepo.Disable2FA(ctx, userID)
}

func (s *twoFactorService) GetStatus(ctx context.Context, userID int64) (*domain.TwoFactorStatus, error) {
	enabled, err := s.twoFactorRepo.Is2FAEnabled(ctx, userID)
	if err != nil {
		return nil, err
	}

	status := &domain.TwoFactorStatus{
		Enabled: enabled,
	}

	// Get confirmed_at if enabled
	if enabled {
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if user != nil && user.TwoFactorConfirmedAt != nil {
			status.ConfirmedAt = user.TwoFactorConfirmedAt
		}
	}

	return status, nil
}

func (s *twoFactorService) CreateChallenge(ctx context.Context, userID int64, ipAddress, userAgent string) (*domain.TwoFactorChallenge, error) {
	// Generate challenge token
	token, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	challenge := &domain.TwoFactorChallenge{
		UserID:         userID,
		ChallengeToken: token,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		ExpiresAt:      time.Now().Add(ChallengeExpiration),
	}

	if err := s.twoFactorRepo.CreateChallenge(ctx, challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}

func (s *twoFactorService) VerifyChallenge(ctx context.Context, token, code string) (*domain.TwoFactorChallenge, error) {
	// Get challenge
	challenge, err := s.twoFactorRepo.GetChallengeByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, ErrChallengeNotFound
	}

	if challenge.IsExpired() {
		return nil, ErrChallengeExpired
	}
	if challenge.IsVerified() {
		return nil, ErrChallengeAlreadyVerified
	}

	// Get secret
	secret, err := s.twoFactorRepo.Get2FASecret(ctx, challenge.UserID)
	if err != nil {
		return nil, err
	}

	// Verify TOTP code
	if !totp.Validate(code, secret) {
		return nil, ErrInvalidTOTPCode
	}

	// Mark as verified
	if err := s.twoFactorRepo.MarkChallengeVerified(ctx, challenge.ID); err != nil {
		return nil, err
	}

	now := time.Now()
	challenge.VerifiedAt = &now

	return challenge, nil
}

func (s *twoFactorService) VerifyChallengeWithBackup(ctx context.Context, token, backupCode string) (*domain.TwoFactorChallenge, error) {
	// Get challenge
	challenge, err := s.twoFactorRepo.GetChallengeByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, ErrChallengeNotFound
	}

	if challenge.IsExpired() {
		return nil, ErrChallengeExpired
	}
	if challenge.IsVerified() {
		return nil, ErrChallengeAlreadyVerified
	}

	// Hash the backup code
	codeHash, err := s.passwordManager.HashPassword(normalizeBackupCode(backupCode))
	if err != nil {
		return nil, err
	}

	// Find and verify backup code
	code, err := s.twoFactorRepo.GetBackupCodeByHash(ctx, challenge.UserID, codeHash)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, ErrInvalidBackupCode
	}

	// Mark backup code as used
	if err := s.twoFactorRepo.MarkBackupCodeUsed(ctx, code.ID); err != nil {
		return nil, err
	}

	// Mark challenge as verified
	if err := s.twoFactorRepo.MarkChallengeVerified(ctx, challenge.ID); err != nil {
		return nil, err
	}

	now := time.Now()
	challenge.VerifiedAt = &now

	return challenge, nil
}

func (s *twoFactorService) RegenerateBackupCodes(ctx context.Context, userID int64) ([]string, error) {
	// Check if 2FA is enabled
	enabled, err := s.twoFactorRepo.Is2FAEnabled(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, ErrTwoFactorNotEnabled
	}

	// Generate new backup codes
	backupCodes, codeHashes, err := s.generateBackupCodes()
	if err != nil {
		return nil, err
	}

	// Store new backup codes (replaces old ones)
	if err := s.twoFactorRepo.CreateBackupCodes(ctx, userID, codeHashes); err != nil {
		return nil, err
	}

	return backupCodes, nil
}

func (s *twoFactorService) GetBackupCodesCount(ctx context.Context, userID int64) (int, error) {
	return s.twoFactorRepo.GetUnusedBackupCodesCount(ctx, userID)
}

func (s *twoFactorService) ValidateTOTP(ctx context.Context, userID int64, code string) (bool, error) {
	secret, err := s.twoFactorRepo.Get2FASecret(ctx, userID)
	if err != nil {
		return false, err
	}
	if secret == "" {
		return false, ErrTwoFactorNotSetup
	}

	return totp.Validate(code, secret), nil
}

// Helper functions

func (s *twoFactorService) generateBackupCodes() ([]string, []string, error) {
	codes := make([]string, BackupCodesCount)
	hashes := make([]string, BackupCodesCount)

	for i := 0; i < BackupCodesCount; i++ {
		code, err := generateBackupCode()
		if err != nil {
			return nil, nil, err
		}
		codes[i] = code
		hash, err := s.passwordManager.HashPassword(normalizeBackupCode(code))
		if err != nil {
			return nil, nil, err
		}
		hashes[i] = hash
	}

	return codes, hashes, nil
}

func generateBackupCode() (string, error) {
	bytes := make([]byte, 7) // 7 bytes = 56 bits, encodes to at least 12 base32 chars
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Format: XXXXX-XXXXX
	code := strings.ToUpper(base32.StdEncoding.EncodeToString(bytes)[:10])
	return code[:5] + "-" + code[5:], nil
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(bytes)[:length], nil
}

func normalizeBackupCode(code string) string {
	// Remove dashes and convert to uppercase
	return strings.ToUpper(strings.ReplaceAll(code, "-", ""))
}
