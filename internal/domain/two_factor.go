package domain

import "time"

// TwoFactorBackupCode represents a backup code for 2FA
type TwoFactorBackupCode struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	CodeHash  string     `json:"-"` // Never expose hash
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// IsUsed checks if the backup code has been used
func (c *TwoFactorBackupCode) IsUsed() bool {
	return c.UsedAt != nil
}

// TwoFactorChallenge represents a 2FA challenge during login
type TwoFactorChallenge struct {
	ID             int64      `json:"id"`
	UserID         int64      `json:"user_id"`
	ChallengeToken string     `json:"challenge_token"`
	IPAddress      string     `json:"ip_address,omitempty"`
	UserAgent      string     `json:"user_agent,omitempty"`
	ExpiresAt      time.Time  `json:"expires_at"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// IsExpired checks if the challenge has expired
func (c *TwoFactorChallenge) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// IsVerified checks if the challenge has been verified
func (c *TwoFactorChallenge) IsVerified() bool {
	return c.VerifiedAt != nil
}

// CanVerify checks if the challenge can still be verified
func (c *TwoFactorChallenge) CanVerify() bool {
	return !c.IsExpired() && !c.IsVerified()
}

// TwoFactorSetup represents the setup data for enabling 2FA
type TwoFactorSetup struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

// TwoFactorStatus represents the 2FA status for a user
type TwoFactorStatus struct {
	Enabled     bool       `json:"enabled"`
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
}
