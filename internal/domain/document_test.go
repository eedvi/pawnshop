package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDocument_TableName(t *testing.T) {
	assert.Equal(t, "documents", Document{}.TableName())
}

func TestAuditLog_TableName(t *testing.T) {
	assert.Equal(t, "audit_logs", AuditLog{}.TableName())
}

func TestSetting_TableName(t *testing.T) {
	assert.Equal(t, "settings", Setting{}.TableName())
}

func TestRefreshToken_TableName(t *testing.T) {
	assert.Equal(t, "refresh_tokens", RefreshToken{}.TableName())
}

func TestRefreshToken_IsExpired_Expired(t *testing.T) {
	rt := &RefreshToken{ExpiresAt: time.Now().Add(-1 * time.Hour)}
	assert.True(t, rt.IsExpired())
}

func TestRefreshToken_IsExpired_NotExpired(t *testing.T) {
	rt := &RefreshToken{ExpiresAt: time.Now().Add(1 * time.Hour)}
	assert.False(t, rt.IsExpired())
}

func TestRefreshToken_IsRevoked_Revoked(t *testing.T) {
	now := time.Now()
	rt := &RefreshToken{RevokedAt: &now}
	assert.True(t, rt.IsRevoked())
}

func TestRefreshToken_IsRevoked_NotRevoked(t *testing.T) {
	rt := &RefreshToken{RevokedAt: nil}
	assert.False(t, rt.IsRevoked())
}

func TestRefreshToken_IsValid_ValidToken(t *testing.T) {
	rt := &RefreshToken{
		ExpiresAt: time.Now().Add(1 * time.Hour),
		RevokedAt: nil,
	}
	assert.True(t, rt.IsValid())
}

func TestRefreshToken_IsValid_ExpiredToken(t *testing.T) {
	rt := &RefreshToken{
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		RevokedAt: nil,
	}
	assert.False(t, rt.IsValid())
}

func TestRefreshToken_IsValid_RevokedToken(t *testing.T) {
	now := time.Now()
	rt := &RefreshToken{
		ExpiresAt: time.Now().Add(1 * time.Hour),
		RevokedAt: &now,
	}
	assert.False(t, rt.IsValid())
}

func TestRefreshToken_IsValid_ExpiredAndRevoked(t *testing.T) {
	now := time.Now()
	rt := &RefreshToken{
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		RevokedAt: &now,
	}
	assert.False(t, rt.IsValid())
}
