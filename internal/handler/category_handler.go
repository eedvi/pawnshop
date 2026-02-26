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

// CategoryHandler handles category endpoints
type CategoryHandler struct {
	categoryService *service.CategoryService
	auditLogger     *middleware.AuditLogger
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(categoryService *service.CategoryService, auditLogger *middleware.AuditLogger) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService, auditLogger: auditLogger}
}

// Create handles category creation
func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var input service.CreateCategoryInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	category, err := h.categoryService.Create(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil {
		description := fmt.Sprintf("Categoría '%s' creada", input.Name)
		h.auditLogger.LogCreateWithDescription(c, "category", category.ID, description, fiber.Map{
			"name":        input.Name,
			"slug":        category.Slug,
			"description": input.Description,
			"parent_id":   input.ParentID,
		})
	}

	return response.Created(c, category)
}

// GetByID handles getting a category by ID
func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID format")
	}

	category, err := h.categoryService.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Category not found")
	}

	return response.OK(c, category)
}

// GetBySlug handles getting a category by slug
func (h *CategoryHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	category, err := h.categoryService.GetBySlug(c.Context(), slug)
	if err != nil {
		return response.NotFound(c, "Category not found")
	}

	return response.OK(c, category)
}

// List handles listing categories
func (h *CategoryHandler) List(c *fiber.Ctx) error {
	params := repository.CategoryListParams{}

	// Parse optional filters
	if parentID := c.Query("parent_id"); parentID != "" {
		id, _ := strconv.ParseInt(parentID, 10, 64)
		params.ParentID = &id
	}
	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		params.IsActive = &active
	}

	categories, err := h.categoryService.List(c.Context(), params)
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.OK(c, categories)
}

// ListTree handles listing categories as a tree structure
func (h *CategoryHandler) ListTree(c *fiber.Ctx) error {
	categories, err := h.categoryService.ListWithChildren(c.Context())
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.OK(c, categories)
}

// Update handles category update
func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	var input service.UpdateCategoryInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	// Get original category for audit
	originalCategory, _ := h.categoryService.GetByID(c.Context(), id)

	category, err := h.categoryService.Update(c.Context(), id, input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && originalCategory != nil {
		description := fmt.Sprintf("Categoría '%s' actualizada", category.Name)
		h.auditLogger.LogUpdateWithDescription(c, "category", id, description,
			fiber.Map{
				"name":        originalCategory.Name,
				"slug":        originalCategory.Slug,
				"description": originalCategory.Description,
				"parent_id":   originalCategory.ParentID,
			},
			fiber.Map{
				"name":        category.Name,
				"slug":        category.Slug,
				"description": category.Description,
				"parent_id":   category.ParentID,
			})
	}

	return response.OK(c, category)
}

// Delete handles category deletion
func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	// Get original category for audit
	originalCategory, _ := h.categoryService.GetByID(c.Context(), id)

	if err := h.categoryService.Delete(c.Context(), id); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && originalCategory != nil {
		description := fmt.Sprintf("Categoría '%s' eliminada", originalCategory.Name)
		h.auditLogger.LogDeleteWithDescription(c, "category", id, description, fiber.Map{
			"name":        originalCategory.Name,
			"slug":        originalCategory.Slug,
			"description": originalCategory.Description,
			"parent_id":   originalCategory.ParentID,
		})
	}

	return response.NoContent(c)
}

// RegisterRoutes registers category routes
func (h *CategoryHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	categories := app.Group("/categories")
	categories.Use(authMiddleware.Authenticate())

	categories.Get("/", authMiddleware.RequirePermission("categories.read"), h.List)
	categories.Get("/tree", authMiddleware.RequirePermission("categories.read"), h.ListTree)
	categories.Post("/", authMiddleware.RequirePermission("categories.create"), h.Create)
	categories.Get("/slug/:slug", authMiddleware.RequirePermission("categories.read"), h.GetBySlug)
	categories.Get("/:id", authMiddleware.RequirePermission("categories.read"), h.GetByID)
	categories.Put("/:id", authMiddleware.RequirePermission("categories.update"), h.Update)
	categories.Delete("/:id", authMiddleware.RequirePermission("categories.delete"), h.Delete)
}
