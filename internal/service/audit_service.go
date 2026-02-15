package service

import (
	"context"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// AuditService handles audit log business logic
type AuditService struct {
	auditRepo repository.AuditLogRepository
}

// NewAuditService creates a new AuditService
func NewAuditService(auditRepo repository.AuditLogRepository) *AuditService {
	return &AuditService{auditRepo: auditRepo}
}

// LogAction creates an audit log entry
func (s *AuditService) LogAction(ctx context.Context, branchID, userID *int64, action, entityType string, entityID *int64, oldValues, newValues interface{}, ipAddress, userAgent string) error {
	log := &domain.AuditLog{
		BranchID:   branchID,
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		OldValues:  oldValues,
		NewValues:  newValues,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	return s.auditRepo.Create(ctx, log)
}

// List retrieves audit logs with filters
func (s *AuditService) List(ctx context.Context, params repository.AuditLogListParams) (*repository.PaginatedResult[domain.AuditLog], error) {
	return s.auditRepo.List(ctx, params)
}

// LogCreate logs a create action
func (s *AuditService) LogCreate(ctx context.Context, branchID, userID *int64, entityType string, entityID int64, newValues interface{}, ipAddress, userAgent string) error {
	return s.LogAction(ctx, branchID, userID, "create", entityType, &entityID, nil, newValues, ipAddress, userAgent)
}

// LogUpdate logs an update action
func (s *AuditService) LogUpdate(ctx context.Context, branchID, userID *int64, entityType string, entityID int64, oldValues, newValues interface{}, ipAddress, userAgent string) error {
	return s.LogAction(ctx, branchID, userID, "update", entityType, &entityID, oldValues, newValues, ipAddress, userAgent)
}

// LogDelete logs a delete action
func (s *AuditService) LogDelete(ctx context.Context, branchID, userID *int64, entityType string, entityID int64, oldValues interface{}, ipAddress, userAgent string) error {
	return s.LogAction(ctx, branchID, userID, "delete", entityType, &entityID, oldValues, nil, ipAddress, userAgent)
}

// LogLogin logs a login action
func (s *AuditService) LogLogin(ctx context.Context, branchID, userID *int64, ipAddress, userAgent string) error {
	return s.LogAction(ctx, branchID, userID, "login", "user", userID, nil, nil, ipAddress, userAgent)
}

// LogLogout logs a logout action
func (s *AuditService) LogLogout(ctx context.Context, branchID, userID *int64, ipAddress, userAgent string) error {
	return s.LogAction(ctx, branchID, userID, "logout", "user", userID, nil, nil, ipAddress, userAgent)
}
