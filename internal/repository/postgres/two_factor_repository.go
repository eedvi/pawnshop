package postgres

import (
	"context"
	"database/sql"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

type twoFactorRepository struct {
	db *DB
}

// NewTwoFactorRepository creates a new 2FA repository
func NewTwoFactorRepository(db *DB) repository.TwoFactorRepository {
	return &twoFactorRepository{db: db}
}

// Backup codes implementation

func (r *twoFactorRepository) CreateBackupCodes(ctx context.Context, userID int64, codeHashes []string) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing backup codes
	_, err = tx.ExecContext(ctx, "DELETE FROM two_factor_backup_codes WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	// Insert new backup codes
	query := `INSERT INTO two_factor_backup_codes (user_id, code_hash) VALUES ($1, $2)`
	for _, hash := range codeHashes {
		_, err = tx.ExecContext(ctx, query, userID, hash)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *twoFactorRepository) GetBackupCodeByHash(ctx context.Context, userID int64, codeHash string) (*domain.TwoFactorBackupCode, error) {
	query := `
		SELECT id, user_id, code_hash, used_at, created_at
		FROM two_factor_backup_codes
		WHERE user_id = $1 AND code_hash = $2 AND used_at IS NULL`

	code := &domain.TwoFactorBackupCode{}
	err := r.db.QueryRowContext(ctx, query, userID, codeHash).Scan(
		&code.ID,
		&code.UserID,
		&code.CodeHash,
		&code.UsedAt,
		&code.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return code, nil
}

func (r *twoFactorRepository) MarkBackupCodeUsed(ctx context.Context, id int64) error {
	query := `UPDATE two_factor_backup_codes SET used_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *twoFactorRepository) GetUnusedBackupCodesCount(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM two_factor_backup_codes WHERE user_id = $1 AND used_at IS NULL`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *twoFactorRepository) DeleteBackupCodes(ctx context.Context, userID int64) error {
	query := `DELETE FROM two_factor_backup_codes WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// Challenges implementation

func (r *twoFactorRepository) CreateChallenge(ctx context.Context, challenge *domain.TwoFactorChallenge) error {
	query := `
		INSERT INTO two_factor_challenges (user_id, challenge_token, ip_address, user_agent, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query,
		challenge.UserID,
		challenge.ChallengeToken,
		challenge.IPAddress,
		challenge.UserAgent,
		challenge.ExpiresAt,
	).Scan(&challenge.ID, &challenge.CreatedAt)
}

func (r *twoFactorRepository) GetChallengeByToken(ctx context.Context, token string) (*domain.TwoFactorChallenge, error) {
	query := `
		SELECT id, user_id, challenge_token, ip_address, user_agent, expires_at, verified_at, created_at
		FROM two_factor_challenges
		WHERE challenge_token = $1`

	challenge := &domain.TwoFactorChallenge{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&challenge.ID,
		&challenge.UserID,
		&challenge.ChallengeToken,
		&challenge.IPAddress,
		&challenge.UserAgent,
		&challenge.ExpiresAt,
		&challenge.VerifiedAt,
		&challenge.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return challenge, nil
}

func (r *twoFactorRepository) MarkChallengeVerified(ctx context.Context, id int64) error {
	query := `UPDATE two_factor_challenges SET verified_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *twoFactorRepository) DeleteExpiredChallenges(ctx context.Context) (int64, error) {
	query := `DELETE FROM two_factor_challenges WHERE expires_at < NOW()`
	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *twoFactorRepository) DeleteChallengesByUser(ctx context.Context, userID int64) error {
	query := `DELETE FROM two_factor_challenges WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// User 2FA settings implementation

func (r *twoFactorRepository) Enable2FA(ctx context.Context, userID int64, secret string) error {
	query := `
		UPDATE users SET
			two_factor_secret = $2,
			updated_at = NOW()
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID, secret)
	return err
}

func (r *twoFactorRepository) Confirm2FA(ctx context.Context, userID int64) error {
	query := `
		UPDATE users SET
			two_factor_enabled = true,
			two_factor_confirmed_at = NOW(),
			updated_at = NOW()
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *twoFactorRepository) Disable2FA(ctx context.Context, userID int64) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Disable 2FA on user
	query := `
		UPDATE users SET
			two_factor_enabled = false,
			two_factor_secret = NULL,
			two_factor_confirmed_at = NULL,
			updated_at = NOW()
		WHERE id = $1`
	_, err = tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	// Delete backup codes
	_, err = tx.ExecContext(ctx, "DELETE FROM two_factor_backup_codes WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	// Delete challenges
	_, err = tx.ExecContext(ctx, "DELETE FROM two_factor_challenges WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *twoFactorRepository) Get2FASecret(ctx context.Context, userID int64) (string, error) {
	query := `SELECT COALESCE(two_factor_secret, '') FROM users WHERE id = $1`
	var secret string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&secret)
	return secret, err
}

func (r *twoFactorRepository) Is2FAEnabled(ctx context.Context, userID int64) (bool, error) {
	query := `SELECT two_factor_enabled FROM users WHERE id = $1`
	var enabled bool
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&enabled)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return enabled, err
}

// CleanupExpiredChallenges removes expired challenges (for scheduled job)
func (r *twoFactorRepository) CleanupExpiredChallenges(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `DELETE FROM two_factor_challenges WHERE expires_at < $1`
	result, err := r.db.ExecContext(ctx, query, time.Now().Add(-olderThan))
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
