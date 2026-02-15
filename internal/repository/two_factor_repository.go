package repository

import (
	"context"

	"pawnshop/internal/domain"
)

// TwoFactorRepository defines the interface for 2FA operations
type TwoFactorRepository interface {
	// Backup codes
	CreateBackupCodes(ctx context.Context, userID int64, codeHashes []string) error
	GetBackupCodeByHash(ctx context.Context, userID int64, codeHash string) (*domain.TwoFactorBackupCode, error)
	MarkBackupCodeUsed(ctx context.Context, id int64) error
	GetUnusedBackupCodesCount(ctx context.Context, userID int64) (int, error)
	DeleteBackupCodes(ctx context.Context, userID int64) error

	// Challenges
	CreateChallenge(ctx context.Context, challenge *domain.TwoFactorChallenge) error
	GetChallengeByToken(ctx context.Context, token string) (*domain.TwoFactorChallenge, error)
	MarkChallengeVerified(ctx context.Context, id int64) error
	DeleteExpiredChallenges(ctx context.Context) (int64, error)
	DeleteChallengesByUser(ctx context.Context, userID int64) error

	// User 2FA settings (these update the users table)
	Enable2FA(ctx context.Context, userID int64, secret string) error
	Confirm2FA(ctx context.Context, userID int64) error
	Disable2FA(ctx context.Context, userID int64) error
	Get2FASecret(ctx context.Context, userID int64) (string, error)
	Is2FAEnabled(ctx context.Context, userID int64) (bool, error)
}
