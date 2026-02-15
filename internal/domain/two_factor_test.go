package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTwoFactorBackupCode_IsUsed_Used(t *testing.T) {
	now := time.Now()
	c := &TwoFactorBackupCode{UsedAt: &now}
	assert.True(t, c.IsUsed())
}

func TestTwoFactorBackupCode_IsUsed_NotUsed(t *testing.T) {
	c := &TwoFactorBackupCode{UsedAt: nil}
	assert.False(t, c.IsUsed())
}

func TestTwoFactorChallenge_IsExpired_Expired(t *testing.T) {
	c := &TwoFactorChallenge{ExpiresAt: time.Now().Add(-1 * time.Hour)}
	assert.True(t, c.IsExpired())
}

func TestTwoFactorChallenge_IsExpired_NotExpired(t *testing.T) {
	c := &TwoFactorChallenge{ExpiresAt: time.Now().Add(1 * time.Hour)}
	assert.False(t, c.IsExpired())
}

func TestTwoFactorChallenge_IsVerified_Verified(t *testing.T) {
	now := time.Now()
	c := &TwoFactorChallenge{VerifiedAt: &now}
	assert.True(t, c.IsVerified())
}

func TestTwoFactorChallenge_IsVerified_NotVerified(t *testing.T) {
	c := &TwoFactorChallenge{VerifiedAt: nil}
	assert.False(t, c.IsVerified())
}

func TestTwoFactorChallenge_CanVerify_Valid(t *testing.T) {
	c := &TwoFactorChallenge{
		ExpiresAt:  time.Now().Add(1 * time.Hour),
		VerifiedAt: nil,
	}
	assert.True(t, c.CanVerify())
}

func TestTwoFactorChallenge_CanVerify_Expired(t *testing.T) {
	c := &TwoFactorChallenge{
		ExpiresAt:  time.Now().Add(-1 * time.Hour),
		VerifiedAt: nil,
	}
	assert.False(t, c.CanVerify())
}

func TestTwoFactorChallenge_CanVerify_AlreadyVerified(t *testing.T) {
	now := time.Now()
	c := &TwoFactorChallenge{
		ExpiresAt:  time.Now().Add(1 * time.Hour),
		VerifiedAt: &now,
	}
	assert.False(t, c.CanVerify())
}

func TestTwoFactorChallenge_CanVerify_ExpiredAndVerified(t *testing.T) {
	now := time.Now()
	c := &TwoFactorChallenge{
		ExpiresAt:  time.Now().Add(-1 * time.Hour),
		VerifiedAt: &now,
	}
	assert.False(t, c.CanVerify())
}
