package mocks

import "github.com/stretchr/testify/mock"

// MockTransaction is a mock implementation of repository.Transaction
type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Rollback() error {
	args := m.Called()
	return args.Error(0)
}
