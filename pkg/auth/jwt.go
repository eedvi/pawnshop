package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType represents the type of token
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// JWTClaims represents the custom JWT claims
type JWTClaims struct {
	UserID      int64    `json:"user_id"`
	Email       string   `json:"email"`
	RoleID      int64    `json:"role_id"`
	BranchID    *int64   `json:"branch_id,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	TokenType   string   `json:"token_type"`
	jwt.RegisteredClaims
}

// JWTConfig contains JWT configuration
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

// JWTManager handles JWT operations
type JWTManager struct {
	config JWTConfig
}

// NewJWTManager creates a new JWTManager
func NewJWTManager(config JWTConfig) *JWTManager {
	return &JWTManager{config: config}
}

// GenerateAccessToken generates a new access token
func (m *JWTManager) GenerateAccessToken(claims JWTClaims) (string, error) {
	now := time.Now()
	claims.TokenType = string(AccessToken)
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(m.config.AccessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    m.config.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.Secret))
}

// GenerateRefreshToken generates a new refresh token
func (m *JWTManager) GenerateRefreshToken(userID int64) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.config.RefreshTokenTTL)

	claims := JWTClaims{
		UserID:    userID,
		TokenType: string(RefreshToken),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.config.Secret))
	return tokenString, expiresAt, err
}

// ValidateToken validates and parses a JWT token
func (m *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateAccessToken validates an access token
func (m *JWTManager) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != string(AccessToken) {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (m *JWTManager) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != string(RefreshToken) {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

// TokenPair represents an access and refresh token pair
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// GenerateTokenPair generates both access and refresh tokens
func (m *JWTManager) GenerateTokenPair(claims JWTClaims) (*TokenPair, error) {
	accessToken, err := m.GenerateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	refreshToken, expiresAt, err := m.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
	}, nil
}
