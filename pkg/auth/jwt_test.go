package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupJWTManager() *JWTManager {
	return NewJWTManager(JWTConfig{
		Secret:          "test-secret-key-for-testing-purposes-only",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test-issuer",
	})
}

func TestNewJWTManager(t *testing.T) {
	config := JWTConfig{
		Secret:          "secret",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test",
	}

	manager := NewJWTManager(config)
	assert.NotNil(t, manager)
}

func TestJWTManager_GenerateAccessToken(t *testing.T) {
	manager := setupJWTManager()
	branchID := int64(1)

	claims := JWTClaims{
		UserID:      123,
		Email:       "test@example.com",
		RoleID:      1,
		BranchID:    &branchID,
		Permissions: []string{"read", "write"},
	}

	token, err := manager.GenerateAccessToken(claims)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	manager := setupJWTManager()

	token, expiresAt, err := manager.GenerateRefreshToken(123)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))
}

func TestJWTManager_ValidateToken_Success(t *testing.T) {
	manager := setupJWTManager()
	branchID := int64(1)

	claims := JWTClaims{
		UserID:      123,
		Email:       "test@example.com",
		RoleID:      1,
		BranchID:    &branchID,
		Permissions: []string{"read", "write"},
	}

	token, err := manager.GenerateAccessToken(claims)
	require.NoError(t, err)

	parsedClaims, err := manager.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, int64(123), parsedClaims.UserID)
	assert.Equal(t, "test@example.com", parsedClaims.Email)
	assert.Equal(t, int64(1), parsedClaims.RoleID)
	assert.Equal(t, int64(1), *parsedClaims.BranchID)
	assert.Equal(t, []string{"read", "write"}, parsedClaims.Permissions)
	assert.Equal(t, string(AccessToken), parsedClaims.TokenType)
}

func TestJWTManager_ValidateToken_InvalidToken(t *testing.T) {
	manager := setupJWTManager()

	_, err := manager.ValidateToken("invalid.token.here")
	assert.Error(t, err)
}

func TestJWTManager_ValidateToken_WrongSecret(t *testing.T) {
	manager1 := setupJWTManager()
	manager2 := NewJWTManager(JWTConfig{
		Secret:          "different-secret-key",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test",
	})

	claims := JWTClaims{
		UserID: 123,
		Email:  "test@example.com",
		RoleID: 1,
	}

	token, err := manager1.GenerateAccessToken(claims)
	require.NoError(t, err)

	// Try to validate with different secret
	_, err = manager2.ValidateToken(token)
	assert.Error(t, err)
}

func TestJWTManager_ValidateToken_ExpiredToken(t *testing.T) {
	// Create manager with very short TTL
	manager := NewJWTManager(JWTConfig{
		Secret:          "test-secret",
		AccessTokenTTL:  -1 * time.Hour, // Already expired
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test",
	})

	claims := JWTClaims{
		UserID: 123,
		Email:  "test@example.com",
		RoleID: 1,
	}

	token, err := manager.GenerateAccessToken(claims)
	require.NoError(t, err)

	_, err = manager.ValidateToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestJWTManager_ValidateAccessToken_Success(t *testing.T) {
	manager := setupJWTManager()

	claims := JWTClaims{
		UserID: 123,
		Email:  "test@example.com",
		RoleID: 1,
	}

	token, err := manager.GenerateAccessToken(claims)
	require.NoError(t, err)

	parsedClaims, err := manager.ValidateAccessToken(token)
	require.NoError(t, err)
	assert.Equal(t, int64(123), parsedClaims.UserID)
	assert.Equal(t, string(AccessToken), parsedClaims.TokenType)
}

func TestJWTManager_ValidateAccessToken_WrongTokenType(t *testing.T) {
	manager := setupJWTManager()

	// Generate refresh token
	token, _, err := manager.GenerateRefreshToken(123)
	require.NoError(t, err)

	// Try to validate as access token
	_, err = manager.ValidateAccessToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token type")
}

func TestJWTManager_ValidateRefreshToken_Success(t *testing.T) {
	manager := setupJWTManager()

	token, _, err := manager.GenerateRefreshToken(123)
	require.NoError(t, err)

	parsedClaims, err := manager.ValidateRefreshToken(token)
	require.NoError(t, err)
	assert.Equal(t, int64(123), parsedClaims.UserID)
	assert.Equal(t, string(RefreshToken), parsedClaims.TokenType)
}

func TestJWTManager_ValidateRefreshToken_WrongTokenType(t *testing.T) {
	manager := setupJWTManager()

	claims := JWTClaims{
		UserID: 123,
		Email:  "test@example.com",
		RoleID: 1,
	}

	// Generate access token
	token, err := manager.GenerateAccessToken(claims)
	require.NoError(t, err)

	// Try to validate as refresh token
	_, err = manager.ValidateRefreshToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token type")
}

func TestJWTManager_GenerateTokenPair(t *testing.T) {
	manager := setupJWTManager()
	branchID := int64(1)

	claims := JWTClaims{
		UserID:      123,
		Email:       "test@example.com",
		RoleID:      1,
		BranchID:    &branchID,
		Permissions: []string{"read", "write"},
	}

	pair, err := manager.GenerateTokenPair(claims)
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.True(t, pair.ExpiresAt.After(time.Now()))
	assert.Equal(t, "Bearer", pair.TokenType)

	// Validate access token
	accessClaims, err := manager.ValidateAccessToken(pair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, int64(123), accessClaims.UserID)

	// Validate refresh token
	refreshClaims, err := manager.ValidateRefreshToken(pair.RefreshToken)
	require.NoError(t, err)
	assert.Equal(t, int64(123), refreshClaims.UserID)
}

func TestJWTManager_GenerateAccessToken_WithNilBranchID(t *testing.T) {
	manager := setupJWTManager()

	claims := JWTClaims{
		UserID:   123,
		Email:    "test@example.com",
		RoleID:   1,
		BranchID: nil, // No branch
	}

	token, err := manager.GenerateAccessToken(claims)
	require.NoError(t, err)

	parsedClaims, err := manager.ValidateToken(token)
	require.NoError(t, err)
	assert.Nil(t, parsedClaims.BranchID)
}

func TestJWTManager_ValidateToken_UnexpectedSigningMethod(t *testing.T) {
	manager := setupJWTManager()

	// Create a token with a different signing method (RS256)
	claims := jwt.MapClaims{
		"user_id": 123,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	// Create an unsigned token and manually set it as valid
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	_, err := manager.ValidateToken(tokenString)
	assert.Error(t, err)
}

func TestJWTManager_TokenClaims_Issuer(t *testing.T) {
	manager := setupJWTManager()

	claims := JWTClaims{
		UserID: 123,
		Email:  "test@example.com",
		RoleID: 1,
	}

	token, err := manager.GenerateAccessToken(claims)
	require.NoError(t, err)

	parsedClaims, err := manager.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, "test-issuer", parsedClaims.Issuer)
}

func TestJWTManager_TokenClaims_Timestamps(t *testing.T) {
	manager := setupJWTManager()
	beforeGenerate := time.Now().Add(-time.Second)

	claims := JWTClaims{
		UserID: 123,
		Email:  "test@example.com",
		RoleID: 1,
	}

	token, err := manager.GenerateAccessToken(claims)
	require.NoError(t, err)

	afterGenerate := time.Now().Add(time.Second)

	parsedClaims, err := manager.ValidateToken(token)
	require.NoError(t, err)

	// IssuedAt should be within the time range
	assert.True(t, parsedClaims.IssuedAt.Time.After(beforeGenerate))
	assert.True(t, parsedClaims.IssuedAt.Time.Before(afterGenerate))

	// ExpiresAt should be after IssuedAt
	assert.True(t, parsedClaims.ExpiresAt.Time.After(parsedClaims.IssuedAt.Time))
}

func TestJWTManager_RefreshToken_ExpiresAt(t *testing.T) {
	manager := setupJWTManager()
	beforeGenerate := time.Now()

	token, expiresAt, err := manager.GenerateRefreshToken(123)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// ExpiresAt should be approximately 7 days from now
	expectedExpiry := beforeGenerate.Add(7 * 24 * time.Hour)
	assert.WithinDuration(t, expectedExpiry, expiresAt, time.Second)
}

func TestJWTManager_EmptyPermissions(t *testing.T) {
	manager := setupJWTManager()

	claims := JWTClaims{
		UserID:      123,
		Email:       "test@example.com",
		RoleID:      1,
		Permissions: []string{}, // Empty permissions
	}

	token, err := manager.GenerateAccessToken(claims)
	require.NoError(t, err)

	parsedClaims, err := manager.ValidateToken(token)
	require.NoError(t, err)
	assert.Empty(t, parsedClaims.Permissions)
}

func TestTokenType_Constants(t *testing.T) {
	assert.Equal(t, TokenType("access"), AccessToken)
	assert.Equal(t, TokenType("refresh"), RefreshToken)
}
