package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"pawnshop/internal/middleware"
	"pawnshop/internal/repository"
	"pawnshop/internal/service"
	"pawnshop/pkg/response"
	"pawnshop/pkg/validator"
)

// BranchHandler handles branch endpoints
type BranchHandler struct {
	branchService *service.BranchService
	auditLogger   *middleware.AuditLogger
}

// NewBranchHandler creates a new BranchHandler
func NewBranchHandler(branchService *service.BranchService, auditLogger *middleware.AuditLogger) *BranchHandler {
	return &BranchHandler{branchService: branchService, auditLogger: auditLogger}
}

// Create handles branch creation
func (h *BranchHandler) Create(c *fiber.Ctx) error {
	var input service.CreateBranchInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	branch, err := h.branchService.Create(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil {
		description := fmt.Sprintf("Sucursal '%s' (código: %s) creada", input.Name, input.Code)
		h.auditLogger.LogCreateWithDescription(c, "branch", branch.ID, description, fiber.Map{
			"name":    input.Name,
			"code":    input.Code,
			"address": input.Address,
			"phone":   input.Phone,
		})
	}

	return response.Created(c, branch)
}

// GetByID handles getting a branch by ID
func (h *BranchHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid branch ID format")
	}

	branch, err := h.branchService.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Branch not found")
	}

	return response.OK(c, branch)
}

// GetByCode handles getting a branch by code
func (h *BranchHandler) GetByCode(c *fiber.Ctx) error {
	code := c.Params("code")

	branch, err := h.branchService.GetByCode(c.Context(), code)
	if err != nil {
		return response.NotFound(c, "Branch not found")
	}

	return response.OK(c, branch)
}

// List handles listing branches
func (h *BranchHandler) List(c *fiber.Ctx) error {
	params := repository.PaginationParams{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
		OrderBy: c.Query("order_by", "name"),
		Order:   c.Query("order", "asc"),
	}

	result, err := h.branchService.List(c.Context(), params)
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.Paginated(c, result.Data, result.Page, result.PerPage, result.Total)
}

// Update handles branch update
func (h *BranchHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid branch ID format")
	}

	var input service.UpdateBranchInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	// Get original branch for audit
	originalBranch, _ := h.branchService.GetByID(c.Context(), id)

	branch, err := h.branchService.Update(c.Context(), id, input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && originalBranch != nil {
		description := fmt.Sprintf("Sucursal '%s' actualizada", branch.Name)
		h.auditLogger.LogUpdateWithDescription(c, "branch", id, description,
			fiber.Map{
				"name":    originalBranch.Name,
				"code":    originalBranch.Code,
				"address": originalBranch.Address,
				"phone":   originalBranch.Phone,
			},
			fiber.Map{
				"name":    branch.Name,
				"code":    branch.Code,
				"address": branch.Address,
				"phone":   branch.Phone,
			})
	}

	return response.OK(c, branch)
}

// Delete handles branch deletion
func (h *BranchHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid branch ID format")
	}

	// Get original branch for audit
	originalBranch, _ := h.branchService.GetByID(c.Context(), id)

	if err := h.branchService.Delete(c.Context(), id); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && originalBranch != nil {
		description := fmt.Sprintf("Sucursal '%s' (código: %s) eliminada", originalBranch.Name, originalBranch.Code)
		h.auditLogger.LogDeleteWithDescription(c, "branch", id, description, fiber.Map{
			"name":    originalBranch.Name,
			"code":    originalBranch.Code,
			"address": originalBranch.Address,
			"phone":   originalBranch.Phone,
		})
	}

	return response.NoContent(c)
}

// Activate handles branch activation
func (h *BranchHandler) Activate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid branch ID format")
	}

	// Get branch before activation for audit
	branch, _ := h.branchService.GetByID(c.Context(), id)

	if err := h.branchService.Activate(c.Context(), id); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && branch != nil {
		description := fmt.Sprintf("Sucursal '%s' activada", branch.Name)
		h.auditLogger.LogCustomAction(c, "activate", "branch", id, description,
			fiber.Map{"is_active": branch.IsActive},
			fiber.Map{"is_active": true})
	}

	return response.OK(c, fiber.Map{"message": "Branch activated successfully"})
}

// Deactivate handles branch deactivation
func (h *BranchHandler) Deactivate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid branch ID format")
	}

	// Get branch before deactivation for audit
	branch, _ := h.branchService.GetByID(c.Context(), id)

	if err := h.branchService.Deactivate(c.Context(), id); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && branch != nil {
		description := fmt.Sprintf("Sucursal '%s' desactivada", branch.Name)
		h.auditLogger.LogCustomAction(c, "deactivate", "branch", id, description,
			fiber.Map{"is_active": branch.IsActive},
			fiber.Map{"is_active": false})
	}

	return response.OK(c, fiber.Map{"message": "Branch deactivated successfully"})
}

// RegisterRoutes registers branch routes
func (h *BranchHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	branches := app.Group("/branches")
	branches.Use(authMiddleware.Authenticate())

	branches.Get("/", authMiddleware.RequirePermission("branches.read"), h.List)
	branches.Post("/", authMiddleware.RequirePermission("branches.create"), h.Create)
	branches.Get("/code/:code", authMiddleware.RequirePermission("branches.read"), h.GetByCode)
	branches.Get("/:id", authMiddleware.RequirePermission("branches.read"), h.GetByID)
	branches.Put("/:id", authMiddleware.RequirePermission("branches.update"), h.Update)
	branches.Delete("/:id", authMiddleware.RequirePermission("branches.delete"), h.Delete)
	branches.Post("/:id/activate", authMiddleware.RequirePermission("branches.update"), h.Activate)
	branches.Post("/:id/deactivate", authMiddleware.RequirePermission("branches.update"), h.Deactivate)
}
