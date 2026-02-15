package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"pawnshop/internal/middleware"
	"pawnshop/internal/service"
	"pawnshop/pkg/response"
	"pawnshop/pkg/validator"
)

// RoleHandler handles role endpoints
type RoleHandler struct {
	roleService service.RoleServiceInterface
	auditLogger *middleware.AuditLogger
}

// NewRoleHandler creates a new RoleHandler
func NewRoleHandler(roleService service.RoleServiceInterface, auditLogger *middleware.AuditLogger) *RoleHandler {
	return &RoleHandler{roleService: roleService, auditLogger: auditLogger}
}

// Create handles role creation
func (h *RoleHandler) Create(c *fiber.Ctx) error {
	var input service.CreateRoleInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	role, err := h.roleService.Create(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, role)
}

// GetByID handles getting a role by ID
func (h *RoleHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID format")
	}

	role, err := h.roleService.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Role not found")
	}

	return response.OK(c, role)
}

// GetByName handles getting a role by name
func (h *RoleHandler) GetByName(c *fiber.Ctx) error {
	name := c.Params("name")

	role, err := h.roleService.GetByName(c.Context(), name)
	if err != nil {
		return response.NotFound(c, "Role not found")
	}

	return response.OK(c, role)
}

// List handles listing roles
func (h *RoleHandler) List(c *fiber.Ctx) error {
	roles, err := h.roleService.List(c.Context())
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.OK(c, roles)
}

// Update handles role update
func (h *RoleHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID")
	}

	var input service.UpdateRoleInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	role, err := h.roleService.Update(c.Context(), id, input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, role)
}

// Delete handles role deletion
func (h *RoleHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID")
	}

	if err := h.roleService.Delete(c.Context(), id); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.NoContent(c)
}

// GetPermissions handles getting available permissions
func (h *RoleHandler) GetPermissions(c *fiber.Ctx) error {
	permissions := h.roleService.GetAvailablePermissions()
	return response.OK(c, permissions)
}

// RegisterRoutes registers role routes
func (h *RoleHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	roles := app.Group("/roles")
	roles.Use(authMiddleware.Authenticate())

	roles.Get("/", authMiddleware.RequirePermission("roles.read"), h.List)
	roles.Get("/permissions", authMiddleware.RequirePermission("roles.read"), h.GetPermissions)
	roles.Post("/", authMiddleware.RequirePermission("roles.create"), h.Create)
	roles.Get("/name/:name", authMiddleware.RequirePermission("roles.read"), h.GetByName)
	roles.Get("/:id", authMiddleware.RequirePermission("roles.read"), h.GetByID)
	roles.Put("/:id", authMiddleware.RequirePermission("roles.update"), h.Update)
	roles.Delete("/:id", authMiddleware.RequirePermission("roles.delete"), h.Delete)
}
