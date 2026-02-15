package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPasswordManager(t *testing.T) {
	pm := NewPasswordManager()
	assert.NotNil(t, pm)
}

func TestPasswordManager_HashPassword(t *testing.T) {
	pm := NewPasswordManager()

	hash, err := pm.HashPassword("TestPassword123!")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Verify hash format: $argon2id$v=...
	assert.True(t, strings.HasPrefix(hash, "$argon2id$v="))
	parts := strings.Split(hash, "$")
	assert.Equal(t, 6, len(parts))
}

func TestPasswordManager_HashPassword_DifferentSalts(t *testing.T) {
	pm := NewPasswordManager()
	password := "SamePassword123!"

	hash1, err := pm.HashPassword(password)
	require.NoError(t, err)

	hash2, err := pm.HashPassword(password)
	require.NoError(t, err)

	// Same password should produce different hashes due to different salts
	assert.NotEqual(t, hash1, hash2)
}

func TestPasswordManager_VerifyPassword_Success(t *testing.T) {
	pm := NewPasswordManager()
	password := "SecurePassword123!"

	hash, err := pm.HashPassword(password)
	require.NoError(t, err)

	valid, err := pm.VerifyPassword(password, hash)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestPasswordManager_VerifyPassword_WrongPassword(t *testing.T) {
	pm := NewPasswordManager()
	password := "CorrectPassword123!"
	wrongPassword := "WrongPassword123!"

	hash, err := pm.HashPassword(password)
	require.NoError(t, err)

	valid, err := pm.VerifyPassword(wrongPassword, hash)
	require.NoError(t, err)
	assert.False(t, valid)
}

func TestPasswordManager_VerifyPassword_InvalidHashFormat(t *testing.T) {
	pm := NewPasswordManager()

	tests := []struct {
		name    string
		hash    string
		wantErr string
	}{
		{
			name:    "not enough parts",
			hash:    "$argon2id$v=19$m=65536",
			wantErr: "invalid hash format",
		},
		{
			name:    "unsupported algorithm",
			hash:    "$bcrypt$v=19$m=65536,t=3,p=2$salt$hash",
			wantErr: "unsupported algorithm",
		},
		{
			name:    "invalid version",
			hash:    "$argon2id$invalid$m=65536,t=3,p=2$salt$hash",
			wantErr: "invalid version",
		},
		{
			name:    "invalid parameters",
			hash:    "$argon2id$v=19$invalid$salt$hash",
			wantErr: "invalid parameters",
		},
		{
			name:    "invalid salt base64",
			hash:    "$argon2id$v=19$m=65536,t=3,p=2$invalid!!!$hash",
			wantErr: "invalid salt",
		},
		{
			name:    "invalid hash base64",
			hash:    "$argon2id$v=19$m=65536,t=3,p=2$dGVzdHNhbHQ$invalid!!!",
			wantErr: "invalid hash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := pm.VerifyPassword("password", tt.hash)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			assert.False(t, valid)
		})
	}
}

func TestPasswordManager_ValidatePasswordStrength(t *testing.T) {
	pm := NewPasswordManager()

	tests := []struct {
		name     string
		password string
		wantErr  string
	}{
		{
			name:     "valid password",
			password: "SecurePass123!",
			wantErr:  "",
		},
		{
			name:     "too short",
			password: "Abc1!",
			wantErr:  "at least 8 characters",
		},
		{
			name:     "no uppercase",
			password: "securepass123!",
			wantErr:  "uppercase letter",
		},
		{
			name:     "no lowercase",
			password: "SECUREPASS123!",
			wantErr:  "lowercase letter",
		},
		{
			name:     "no digit",
			password: "SecurePassword!",
			wantErr:  "digit",
		},
		{
			name:     "no special character",
			password: "SecurePass123",
			wantErr:  "special character",
		},
		{
			name:     "valid with different special chars",
			password: "Test@Pass#123",
			wantErr:  "",
		},
		{
			name:     "exactly 8 characters valid",
			password: "Abcdef1!",
			wantErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pm.ValidatePasswordStrength(tt.password)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestGenerateRandomToken(t *testing.T) {
	// Test generating tokens of different lengths
	lengths := []int{16, 32, 64}

	for _, length := range lengths {
		t.Run("length_"+string(rune('0'+length/10))+string(rune('0'+length%10)), func(t *testing.T) {
			token, err := GenerateRandomToken(length)
			require.NoError(t, err)
			assert.NotEmpty(t, token)
		})
	}
}

func TestGenerateRandomToken_Uniqueness(t *testing.T) {
	tokens := make(map[string]bool)

	// Generate multiple tokens and ensure they're unique
	for i := 0; i < 100; i++ {
		token, err := GenerateRandomToken(32)
		require.NoError(t, err)
		assert.False(t, tokens[token], "token should be unique")
		tokens[token] = true
	}
}

func TestHashToken(t *testing.T) {
	token := "test-token-123"

	hash := HashToken(token)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, token, hash)

	// Same token should produce same hash (deterministic)
	hash2 := HashToken(token)
	assert.Equal(t, hash, hash2)

	// Different tokens should produce different hashes
	hash3 := HashToken("different-token")
	assert.NotEqual(t, hash, hash3)
}

func TestPasswordManager_VerifyPassword_EmptyPassword(t *testing.T) {
	pm := NewPasswordManager()

	// Hash a normal password
	hash, err := pm.HashPassword("ValidPassword123!")
	require.NoError(t, err)

	// Verify with empty password should fail
	valid, err := pm.VerifyPassword("", hash)
	require.NoError(t, err)
	assert.False(t, valid)
}

func TestPasswordManager_HashPassword_EmptyPassword(t *testing.T) {
	pm := NewPasswordManager()

	// Empty password should still hash (validation is separate)
	hash, err := pm.HashPassword("")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// And verify should work
	valid, err := pm.VerifyPassword("", hash)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestPasswordManager_HashPassword_LongPassword(t *testing.T) {
	pm := NewPasswordManager()

	// Test with a very long password
	longPassword := strings.Repeat("SecurePass123!", 100)

	hash, err := pm.HashPassword(longPassword)
	require.NoError(t, err)

	valid, err := pm.VerifyPassword(longPassword, hash)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestPasswordManager_HashPassword_UnicodePassword(t *testing.T) {
	pm := NewPasswordManager()

	// Test with unicode characters
	password := "Contraseña123!日本語"

	hash, err := pm.HashPassword(password)
	require.NoError(t, err)

	valid, err := pm.VerifyPassword(password, hash)
	require.NoError(t, err)
	assert.True(t, valid)
}
