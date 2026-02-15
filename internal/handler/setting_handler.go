package handler

import (
	"github.com/gofiber/fiber/v2"
	"pawnshop/internal/middleware"
	"pawnshop/internal/service"
	"pawnshop/pkg/response"
	"pawnshop/pkg/validator"
)

// SettingHandler handles setting endpoints
type SettingHandler struct {
	settingService service.SettingServiceInterface
	auditLogger    *middleware.AuditLogger
}

// NewSettingHandler creates a new SettingHandler
func NewSettingHandler(settingService service.SettingServiceInterface, auditLogger *middleware.AuditLogger) *SettingHandler {
	return &SettingHandler{settingService: settingService, auditLogger: auditLogger}
}

// Get retrieves a setting by key
func (h *SettingHandler) Get(c *fiber.Ctx) error {
	key := c.Params("key")
	branchID := getBranchIDFromQuery(c)

	setting, err := h.settingService.Get(c.Context(), key, branchID)
	if err != nil {
		return response.NotFound(c, "Setting not found")
	}

	return response.OK(c, setting)
}

// List retrieves all settings
func (h *SettingHandler) List(c *fiber.Ctx) error {
	branchID := getBranchIDFromQuery(c)

	settings, err := h.settingService.GetAll(c.Context(), branchID)
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.OK(c, settings)
}

// GetMerged retrieves all settings merged (branch overrides global)
func (h *SettingHandler) GetMerged(c *fiber.Ctx) error {
	branchID := getBranchIDFromQuery(c)

	settings, err := h.settingService.GetMerged(c.Context(), branchID)
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.OK(c, settings)
}

// Set creates or updates a setting
func (h *SettingHandler) Set(c *fiber.Ctx) error {
	var input service.SetSettingInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	setting, err := h.settingService.Set(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, setting)
}

// SetMultiple sets multiple settings at once
func (h *SettingHandler) SetMultiple(c *fiber.Ctx) error {
	var inputs []service.SetSettingInput
	if err := c.BodyParser(&inputs); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if len(inputs) == 0 {
		return response.BadRequest(c, "No settings provided")
	}

	for _, input := range inputs {
		if errors := validator.Validate(&input); errors != nil {
			return response.ValidationError(c, errors)
		}
	}

	if err := h.settingService.SetMultiple(c.Context(), inputs); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{"message": "Settings updated successfully"})
}

// Delete deletes a setting
func (h *SettingHandler) Delete(c *fiber.Ctx) error {
	key := c.Params("key")
	branchID := getBranchIDFromQuery(c)

	if err := h.settingService.Delete(c.Context(), key, branchID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.NoContent(c)
}

// RegisterRoutes registers setting routes
func (h *SettingHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	settings := app.Group("/settings")
	settings.Use(authMiddleware.Authenticate())

	settings.Get("/", authMiddleware.RequirePermission("settings.read"), h.List)
	settings.Get("/merged", authMiddleware.RequirePermission("settings.read"), h.GetMerged)
	settings.Post("/", authMiddleware.RequirePermission("settings.update"), h.Set)
	settings.Post("/bulk", authMiddleware.RequirePermission("settings.update"), h.SetMultiple)
	settings.Get("/:key", authMiddleware.RequirePermission("settings.read"), h.Get)
	settings.Delete("/:key", authMiddleware.RequirePermission("settings.update"), h.Delete)
}

// Helper function to get branch ID from query
func getBranchIDFromQuery(c *fiber.Ctx) *int64 {
	branchID := c.QueryInt("branch_id", 0)
	if branchID > 0 {
		id := int64(branchID)
		return &id
	}
	return nil
}
