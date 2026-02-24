package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
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

// safeGo runs a function in a goroutine with panic recovery
func safeGo(fn func(), operation string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Interface("panic", r).
					Str("operation", operation).
					Msg("Panic recovered in audit logging goroutine")
			}
		}()
		fn()
	}()
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
			safeGo(func() {
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
			}, "audit.LogAction")
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

	safeGo(func() {
		l.auditService.LogCreate(
			c.Context(),
			branchID,
			userID,
			entityType,
			entityID,
			newValues,
			c.IP(),
			c.Get("User-Agent"),
		)
	}, "audit.LogCreate")
}

// LogUpdate logs an update action
func (l *AuditLogger) LogUpdate(c *fiber.Ctx, entityType string, entityID int64, oldValues, newValues interface{}) {
	user := GetUser(c)
	var userID, branchID *int64
	if user != nil {
		userID = &user.ID
		branchID = user.BranchID
	}

	safeGo(func() {
		l.auditService.LogUpdate(
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
	}, "audit.LogUpdate")
}

// LogDelete logs a delete action
func (l *AuditLogger) LogDelete(c *fiber.Ctx, entityType string, entityID int64, oldValues interface{}) {
	user := GetUser(c)
	var userID, branchID *int64
	if user != nil {
		userID = &user.ID
		branchID = user.BranchID
	}

	safeGo(func() {
		l.auditService.LogDelete(
			c.Context(),
			branchID,
			userID,
			entityType,
			entityID,
			oldValues,
			c.IP(),
			c.Get("User-Agent"),
		)
	}, "audit.LogDelete")
}

// LogLogin logs a login action
func (l *AuditLogger) LogLogin(c *fiber.Ctx, userID int64, branchID *int64) {
	safeGo(func() {
		l.auditService.LogLogin(
			c.Context(),
			branchID,
			&userID,
			c.IP(),
			c.Get("User-Agent"),
		)
	}, "audit.LogLogin")
}

// LogLogout logs a logout action
func (l *AuditLogger) LogLogout(c *fiber.Ctx, userID int64, branchID *int64) {
	safeGo(func() {
		l.auditService.LogLogout(
			c.Context(),
			branchID,
			&userID,
			c.IP(),
			c.Get("User-Agent"),
		)
	}, "audit.LogLogout")
}

// LogCreateWithDescription logs a create action with description
func (l *AuditLogger) LogCreateWithDescription(c *fiber.Ctx, entityType string, entityID int64, description string, newValues interface{}) {
	user := GetUser(c)
	var userID, branchID *int64
	if user != nil {
		userID = &user.ID
		branchID = user.BranchID
	}

	safeGo(func() {
		l.auditService.LogActionWithDescription(
			c.Context(),
			branchID,
			userID,
			"create",
			entityType,
			&entityID,
			description,
			nil,
			newValues,
			c.IP(),
			c.Get("User-Agent"),
		)
	}, "audit.LogCreateWithDescription")
}

// LogUpdateWithDescription logs an update action with description
func (l *AuditLogger) LogUpdateWithDescription(c *fiber.Ctx, entityType string, entityID int64, description string, oldValues, newValues interface{}) {
	user := GetUser(c)
	var userID, branchID *int64
	if user != nil {
		userID = &user.ID
		branchID = user.BranchID
	}

	safeGo(func() {
		l.auditService.LogActionWithDescription(
			c.Context(),
			branchID,
			userID,
			"update",
			entityType,
			&entityID,
			description,
			oldValues,
			newValues,
			c.IP(),
			c.Get("User-Agent"),
		)
	}, "audit.LogUpdateWithDescription")
}

// LogDeleteWithDescription logs a delete action with description
func (l *AuditLogger) LogDeleteWithDescription(c *fiber.Ctx, entityType string, entityID int64, description string, oldValues interface{}) {
	user := GetUser(c)
	var userID, branchID *int64
	if user != nil {
		userID = &user.ID
		branchID = user.BranchID
	}

	safeGo(func() {
		l.auditService.LogActionWithDescription(
			c.Context(),
			branchID,
			userID,
			"delete",
			entityType,
			&entityID,
			description,
			oldValues,
			nil,
			c.IP(),
			c.Get("User-Agent"),
		)
	}, "audit.LogDeleteWithDescription")
}

// LogCustomAction logs a custom action with description
func (l *AuditLogger) LogCustomAction(c *fiber.Ctx, action, entityType string, entityID int64, description string, oldValues, newValues interface{}) {
	user := GetUser(c)
	var userID, branchID *int64
	if user != nil {
		userID = &user.ID
		branchID = user.BranchID
	}

	safeGo(func() {
		l.auditService.LogActionWithDescription(
			c.Context(),
			branchID,
			userID,
			action,
			entityType,
			&entityID,
			description,
			oldValues,
			newValues,
			c.IP(),
			c.Get("User-Agent"),
		)
	}, "audit.LogCustomAction")
}
