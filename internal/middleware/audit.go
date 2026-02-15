package middleware

import (
	"github.com/gofiber/fiber/v2"
	"pawnshop/internal/service"
)

// AuditMiddleware handles audit logging
type AuditMiddleware struct {
	auditService *service.AuditService
}

// NewAuditMiddleware creates a new AuditMiddleware
func NewAuditMiddleware(auditService *service.AuditService) *AuditMiddleware {
	return &AuditMiddleware{auditService: auditService}
}

// LogAction logs an action to the audit log
func (m *AuditMiddleware) LogAction(action, entityType string, entityID *int64, oldValues, newValues interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Execute the handler first
		err := c.Next()

		// Only log if the request was successful (2xx status)
		status := c.Response().StatusCode()
		if status >= 200 && status < 300 {
			user := GetUser(c)
			var userID, branchID *int64
			if user != nil {
				userID = &user.ID
				branchID = user.BranchID
			}

			// Log the action asynchronously to not block the response
			go func() {
				m.auditService.LogAction(
					c.Context(),
					branchID,
					userID,
					action,
					entityType,
					entityID,
					oldValues,
					newValues,
					c.IP(),
					c.Get("User-Agent"),
				)
			}()
		}

		return err
	}
}

// AuditLogger is a helper to log audit events from handlers
type AuditLogger struct {
	auditService *service.AuditService
}

// NewAuditLogger creates a new AuditLogger
func NewAuditLogger(auditService *service.AuditService) *AuditLogger {
	return &AuditLogger{auditService: auditService}
}

// LogCreate logs a create action
func (l *AuditLogger) LogCreate(c *fiber.Ctx, entityType string, entityID int64, newValues interface{}) {
	user := GetUser(c)
	var userID, branchID *int64
	if user != nil {
		userID = &user.ID
		branchID = user.BranchID
	}

	go l.auditService.LogCreate(
		c.Context(),
		branchID,
		userID,
		entityType,
		entityID,
		newValues,
		c.IP(),
		c.Get("User-Agent"),
	)
}

// LogUpdate logs an update action
func (l *AuditLogger) LogUpdate(c *fiber.Ctx, entityType string, entityID int64, oldValues, newValues interface{}) {
	user := GetUser(c)
	var userID, branchID *int64
	if user != nil {
		userID = &user.ID
		branchID = user.BranchID
	}

	go l.auditService.LogUpdate(
		c.Context(),
		branchID,
		userID,
		entityType,
		entityID,
		oldValues,
		newValues,
		c.IP(),
		c.Get("User-Agent"),
	)
}

// LogDelete logs a delete action
func (l *AuditLogger) LogDelete(c *fiber.Ctx, entityType string, entityID int64, oldValues interface{}) {
	user := GetUser(c)
	var userID, branchID *int64
	if user != nil {
		userID = &user.ID
		branchID = user.BranchID
	}

	go l.auditService.LogDelete(
		c.Context(),
		branchID,
		userID,
		entityType,
		entityID,
		oldValues,
		c.IP(),
		c.Get("User-Agent"),
	)
}

// LogLogin logs a login action
func (l *AuditLogger) LogLogin(c *fiber.Ctx, userID int64, branchID *int64) {
	go l.auditService.LogLogin(
		c.Context(),
		branchID,
		&userID,
		c.IP(),
		c.Get("User-Agent"),
	)
}

// LogLogout logs a logout action
func (l *AuditLogger) LogLogout(c *fiber.Ctx, userID int64, branchID *int64) {
	go l.auditService.LogLogout(
		c.Context(),
		branchID,
		&userID,
		c.IP(),
		c.Get("User-Agent"),
	)
}
