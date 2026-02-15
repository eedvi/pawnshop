package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/internal/repository/mocks"
)

func setupAuditService() (*AuditService, *mocks.MockAuditLogRepository) {
	auditRepo := new(mocks.MockAuditLogRepository)
	service := NewAuditService(auditRepo)
	return service, auditRepo
}

func TestAuditService_LogAction_Success(t *testing.T) {
	service, auditRepo := setupAuditService()
	ctx := context.Background()

	branchID := int64(1)
	userID := int64(10)
	entityID := int64(100)

	auditRepo.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil)

	err := service.LogAction(ctx, &branchID, &userID, "create", "customer", &entityID, nil, map[string]string{"name": "John"}, "192.168.1.1", "Mozilla/5.0")

	assert.NoError(t, err)
	auditRepo.AssertExpectations(t)
}

func TestAuditService_LogAction_Error(t *testing.T) {
	service, auditRepo := setupAuditService()
	ctx := context.Background()

	auditRepo.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(errors.New("db error"))

	err := service.LogAction(ctx, nil, nil, "create", "customer", nil, nil, nil, "", "")

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	auditRepo.AssertExpectations(t)
}

func TestAuditService_List_Success(t *testing.T) {
	service, auditRepo := setupAuditService()
	ctx := context.Background()

	params := repository.AuditLogListParams{
		PaginationParams: repository.PaginationParams{Page: 1, PerPage: 10},
	}

	result := &repository.PaginatedResult[domain.AuditLog]{
		Data:       []domain.AuditLog{{Action: "create"}},
		Total:      1,
		Page:       1,
		PerPage:    10,
		TotalPages: 1,
	}

	auditRepo.On("List", ctx, params).Return(result, nil)

	res, err := service.List(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 1)
	auditRepo.AssertExpectations(t)
}

func TestAuditService_LogCreate(t *testing.T) {
	service, auditRepo := setupAuditService()
	ctx := context.Background()

	branchID := int64(1)
	userID := int64(10)

	auditRepo.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil)

	err := service.LogCreate(ctx, &branchID, &userID, "customer", 100, map[string]string{"name": "John"}, "192.168.1.1", "Agent")

	assert.NoError(t, err)
	auditRepo.AssertExpectations(t)

	// Verify the log was created with correct action
	call := auditRepo.Calls[0]
	log := call.Arguments.Get(1).(*domain.AuditLog)
	assert.Equal(t, "create", log.Action)
	assert.Equal(t, "customer", log.EntityType)
	assert.Nil(t, log.OldValues)
	assert.NotNil(t, log.NewValues)
}

func TestAuditService_LogUpdate(t *testing.T) {
	service, auditRepo := setupAuditService()
	ctx := context.Background()

	branchID := int64(1)
	userID := int64(10)

	auditRepo.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil)

	err := service.LogUpdate(ctx, &branchID, &userID, "customer", 100, map[string]string{"name": "Old"}, map[string]string{"name": "New"}, "192.168.1.1", "Agent")

	assert.NoError(t, err)
	auditRepo.AssertExpectations(t)

	call := auditRepo.Calls[0]
	log := call.Arguments.Get(1).(*domain.AuditLog)
	assert.Equal(t, "update", log.Action)
	assert.NotNil(t, log.OldValues)
	assert.NotNil(t, log.NewValues)
}

func TestAuditService_LogDelete(t *testing.T) {
	service, auditRepo := setupAuditService()
	ctx := context.Background()

	branchID := int64(1)
	userID := int64(10)

	auditRepo.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil)

	err := service.LogDelete(ctx, &branchID, &userID, "customer", 100, map[string]string{"name": "John"}, "192.168.1.1", "Agent")

	assert.NoError(t, err)
	auditRepo.AssertExpectations(t)

	call := auditRepo.Calls[0]
	log := call.Arguments.Get(1).(*domain.AuditLog)
	assert.Equal(t, "delete", log.Action)
	assert.NotNil(t, log.OldValues)
	assert.Nil(t, log.NewValues)
}

func TestAuditService_LogLogin(t *testing.T) {
	service, auditRepo := setupAuditService()
	ctx := context.Background()

	branchID := int64(1)
	userID := int64(10)

	auditRepo.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil)

	err := service.LogLogin(ctx, &branchID, &userID, "192.168.1.1", "Agent")

	assert.NoError(t, err)
	auditRepo.AssertExpectations(t)

	call := auditRepo.Calls[0]
	log := call.Arguments.Get(1).(*domain.AuditLog)
	assert.Equal(t, "login", log.Action)
	assert.Equal(t, "user", log.EntityType)
}

func TestAuditService_LogLogout(t *testing.T) {
	service, auditRepo := setupAuditService()
	ctx := context.Background()

	branchID := int64(1)
	userID := int64(10)

	auditRepo.On("Create", ctx, mock.AnythingOfType("*domain.AuditLog")).Return(nil)

	err := service.LogLogout(ctx, &branchID, &userID, "192.168.1.1", "Agent")

	assert.NoError(t, err)
	auditRepo.AssertExpectations(t)

	call := auditRepo.Calls[0]
	log := call.Arguments.Get(1).(*domain.AuditLog)
	assert.Equal(t, "logout", log.Action)
	assert.Equal(t, "user", log.EntityType)
}
