package handler

import (
	"github.com/gofiber/fiber/v2"
	"pawnshop/internal/middleware"
	"pawnshop/internal/repository"
	"pawnshop/internal/service"
	"pawnshop/pkg/response"
)

// AuditHandler handles audit log endpoints
type AuditHandler struct {
	auditService *service.AuditService
}

// NewAuditHandler creates a new AuditHandler
func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

// List retrieves audit logs with filters
func (h *AuditHandler) List(c *fiber.Ctx) error {
	params := repository.AuditLogListParams{
		PaginationParams: repository.PaginationParams{
			Page:    c.QueryInt("page", 1),
			PerPage: c.QueryInt("per_page", 50),
			OrderBy: c.Query("order_by", "created_at"),
			Order:   c.Query("order", "desc"),
		},
	}

	// Parse optional filters
	if branchID := c.QueryInt("branch_id", 0); branchID > 0 {
		id := int64(branchID)
		params.BranchID = &id
	}
	if userID := c.QueryInt("user_id", 0); userID > 0 {
		id := int64(userID)
		params.UserID = &id
	}
	if entityID := c.QueryInt("entity_id", 0); entityID > 0 {
		id := int64(entityID)
		params.EntityID = &id
	}
	params.Action = c.Query("action")
	params.EntityType = c.Query("entity_type")

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		params.DateFrom = &dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		params.DateTo = &dateTo
	}

	result, err := h.auditService.List(c.Context(), params)
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.Paginated(c, result.Data, result.Page, result.PerPage, result.Total)
}

// GetStats retrieves audit statistics
func (h *AuditHandler) GetStats(c *fiber.Ctx) error {
	params := repository.AuditLogListParams{}

	// Parse optional filters
	if branchID := c.QueryInt("branch_id", 0); branchID > 0 {
		id := int64(branchID)
		params.BranchID = &id
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		params.DateFrom = &dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		params.DateTo = &dateTo
	}

	stats, err := h.auditService.GetStats(c.Context(), params)
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.OK(c, stats)
}

// RegisterRoutes registers audit log routes
func (h *AuditHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	audit := app.Group("/audit")
	audit.Use(authMiddleware.Authenticate())

	audit.Get("/", authMiddleware.RequirePermission("audit.read"), h.List)
	audit.Get("/stats", authMiddleware.RequirePermission("audit.read"), h.GetStats)
}
