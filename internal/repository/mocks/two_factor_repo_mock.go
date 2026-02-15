package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
)

// MockTwoFactorRepository is a mock implementation of TwoFactorRepository
type MockTwoFactorRepository struct {
	mock.Mock
}

func (m *MockTwoFactorRepository) CreateBackupCodes(ctx context.Context, userID int64, codeHashes []string) error {
	args := m.Called(ctx, userID, codeHashes)
	return args.Error(0)
}

func (m *MockTwoFactorRepository) GetBackupCodeByHash(ctx context.Context, userID int64, codeHash string) (*domain.TwoFactorBackupCode, error) {
	args := m.Called(ctx, userID, codeHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TwoFactorBackupCode), args.Error(1)
}

func (m *MockTwoFactorRepository) MarkBackupCodeUsed(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTwoFactorRepository) GetUnusedBackupCodesCount(ctx context.Context, userID int64) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockTwoFactorRepository) DeleteBackupCodes(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockTwoFactorRepository) CreateChallenge(ctx context.Context, challenge *domain.TwoFactorChallenge) error {
	args := m.Called(ctx, challenge)
	return args.Error(0)
}

func (m *MockTwoFactorRepository) GetChallengeByToken(ctx context.Context, token string) (*domain.TwoFactorChallenge, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TwoFactorChallenge), args.Error(1)
}

func (m *MockTwoFactorRepository) MarkChallengeVerified(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTwoFactorRepository) DeleteExpiredChallenges(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTwoFactorRepository) DeleteChallengesByUser(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockTwoFactorRepository) Enable2FA(ctx context.Context, userID int64, secret string) error {
	args := m.Called(ctx, userID, secret)
	return args.Error(0)
}

func (m *MockTwoFactorRepository) Confirm2FA(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockTwoFactorRepository) Disable2FA(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockTwoFactorRepository) Get2FASecret(ctx context.Context, userID int64) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockTwoFactorRepository) Is2FAEnabled(ctx context.Context, userID int64) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}
