package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"pawnshop/internal/domain"
	"pawnshop/internal/middleware"
	"pawnshop/internal/repository"
	"pawnshop/internal/service"
	"pawnshop/pkg/response"
	"pawnshop/pkg/validator"
)

// ItemHandler handles item endpoints
type ItemHandler struct {
	itemService *service.ItemService
	auditLogger *middleware.AuditLogger
}

// NewItemHandler creates a new ItemHandler
func NewItemHandler(itemService *service.ItemService, auditLogger *middleware.AuditLogger) *ItemHandler {
	return &ItemHandler{itemService: itemService, auditLogger: auditLogger}
}

// Create handles item creation
func (h *ItemHandler) Create(c *fiber.Ctx) error {
	var input service.CreateItemInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	// Set branch and user from context
	user := middleware.GetUser(c)
	if input.BranchID == 0 && user.BranchID != nil {
		input.BranchID = *user.BranchID
	}
	input.CreatedBy = user.ID

	// Validate
	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	item, err := h.itemService.Create(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil {
		description := fmt.Sprintf("Artículo '%s' (SKU: %s) creado con avalúo de Q%.2f", input.Name, item.SKU, input.AppraisedValue)
		h.auditLogger.LogCreateWithDescription(c, "item", item.ID, description, fiber.Map{
			"sku":              item.SKU,
			"name":             input.Name,
			"category_id":      input.CategoryID,
			"appraised_value":  input.AppraisedValue,
			"acquisition_type": input.AcquisitionType,
		})
	}

	return response.Created(c, item)
}

// GetByID handles getting an item by ID
func (h *ItemHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid item ID format")
	}

	item, err := h.itemService.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Item not found")
	}

	return response.OK(c, item)
}

// GetBySKU handles getting an item by SKU
func (h *ItemHandler) GetBySKU(c *fiber.Ctx) error {
	sku := c.Params("sku")

	item, err := h.itemService.GetBySKU(c.Context(), sku)
	if err != nil {
		return response.NotFound(c, "Item not found")
	}

	return response.OK(c, item)
}

// List handles listing items
func (h *ItemHandler) List(c *fiber.Ctx) error {
	user := middleware.GetUser(c)

	params := repository.ItemListParams{
		PaginationParams: repository.PaginationParams{
			Page:    c.QueryInt("page", 1),
			PerPage: c.QueryInt("per_page", 20),
			OrderBy: c.Query("order_by", "created_at"),
			Order:   c.Query("order", "desc"),
		},
		Search: c.Query("search"),
	}

	// Filter by user's branch if not admin
	if user.BranchID != nil {
		params.BranchID = *user.BranchID
	} else if branchID := c.Query("branch_id"); branchID != "" {
		id, _ := strconv.ParseInt(branchID, 10, 64)
		params.BranchID = id
	}

	// Parse optional filters
	if categoryID := c.Query("category_id"); categoryID != "" {
		id, _ := strconv.ParseInt(categoryID, 10, 64)
		params.CategoryID = &id
	}
	if customerID := c.Query("customer_id"); customerID != "" {
		id, _ := strconv.ParseInt(customerID, 10, 64)
		params.CustomerID = &id
	}
	if status := c.Query("status"); status != "" {
		s := domain.ItemStatus(status)
		params.Status = &s
	}

	result, err := h.itemService.List(c.Context(), params)
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.Paginated(c, result.Data, result.Page, result.PerPage, result.Total)
}

// Update handles item update
func (h *ItemHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid item ID")
	}

	var input service.UpdateItemInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	input.UpdatedBy = middleware.GetUser(c).ID

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	// Get original item for audit
	originalItem, _ := h.itemService.GetByID(c.Context(), id)

	item, err := h.itemService.Update(c.Context(), id, input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && originalItem != nil {
		description := fmt.Sprintf("Artículo '%s' (SKU: %s) actualizado", item.Name, item.SKU)

		oldValues := fiber.Map{
			"name":            originalItem.Name,
			"appraised_value": originalItem.AppraisedValue,
			"loan_value":      originalItem.LoanValue,
		}

		newValues := fiber.Map{
			"name":            input.Name,
			"appraised_value": input.AppraisedValue,
			"loan_value":      input.LoanValue,
		}

		h.auditLogger.LogUpdateWithDescription(c, "item", id, description, oldValues, newValues)
	}

	return response.OK(c, item)
}

// Delete handles item deletion
func (h *ItemHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid item ID")
	}

	// Get item before deleting for audit
	item, _ := h.itemService.GetByID(c.Context(), id)

	user := middleware.GetUser(c)
	if err := h.itemService.Delete(c.Context(), id, user.ID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && item != nil {
		description := fmt.Sprintf("Artículo '%s' (SKU: %s) eliminado", item.Name, item.SKU)
		h.auditLogger.LogDeleteWithDescription(c, "item", id, description, fiber.Map{
			"sku":              item.SKU,
			"name":             item.Name,
			"appraised_value":  item.AppraisedValue,
			"status":           item.Status,
			"acquisition_type": item.AcquisitionType,
		})
	}

	return response.NoContent(c)
}

// UpdateStatus handles item status update
func (h *ItemHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid item ID")
	}

	var input service.UpdateStatusInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	input.UpdatedBy = middleware.GetUser(c).ID

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	if err := h.itemService.UpdateStatus(c.Context(), id, input); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{"message": "Item status updated successfully"})
}

// MarkForSale handles marking an item for sale
func (h *ItemHandler) MarkForSale(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid item ID")
	}

	var input struct {
		SalePrice float64 `json:"sale_price" validate:"required,gt=0"`
	}
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	// Get item before marking for sale for audit
	item, _ := h.itemService.GetByID(c.Context(), id)

	user := middleware.GetUser(c)
	if err := h.itemService.MarkForSale(c.Context(), id, input.SalePrice, user.ID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && item != nil {
		description := fmt.Sprintf("Artículo '%s' (SKU: %s) marcado para venta en Q%.2f", item.Name, item.SKU, input.SalePrice)
		h.auditLogger.LogCustomAction(c, "mark_for_sale", "item", id, description,
			fiber.Map{
				"status":     item.Status,
				"sale_price": item.SalePrice,
			},
			fiber.Map{
				"status":     "for_sale",
				"sale_price": input.SalePrice,
			})
	}

	return response.OK(c, fiber.Map{"message": "Item marked for sale successfully"})
}

// GetForSale handles getting items available for sale
func (h *ItemHandler) GetForSale(c *fiber.Ctx) error {
	user := middleware.GetUser(c)

	var branchID int64
	if user.BranchID != nil {
		branchID = *user.BranchID
	} else if bid := c.Query("branch_id"); bid != "" {
		branchID, _ = strconv.ParseInt(bid, 10, 64)
	}

	if branchID == 0 {
		return response.BadRequest(c, "Branch ID is required")
	}

	items, err := h.itemService.GetAvailableForSale(c.Context(), branchID)
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.OK(c, items)
}

// MarkAsDelivered handles marking an item as physically delivered to customer
func (h *ItemHandler) MarkAsDelivered(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid item ID")
	}

	var input struct {
		Notes string `json:"notes"`
	}
	if err := c.BodyParser(&input); err != nil {
		// Allow empty body
		input.Notes = ""
	}

	user := middleware.GetUser(c)
	deliveryInput := service.MarkAsDeliveredInput{
		ItemID:    id,
		Notes:     input.Notes,
		UpdatedBy: user.ID,
	}

	item, err := h.itemService.MarkAsDelivered(c.Context(), deliveryInput)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Log delivery action
	if h.auditLogger != nil {
		h.auditLogger.LogUpdate(c, "item", id, fiber.Map{
			"delivered_at": nil,
		}, fiber.Map{
			"delivered_at": item.DeliveredAt,
			"notes":        input.Notes,
		})
	}

	return response.OK(c, fiber.Map{
		"message": "Item marked as delivered successfully",
		"item":    item,
	})
}

// GetPendingDeliveries handles getting items pending physical delivery
func (h *ItemHandler) GetPendingDeliveries(c *fiber.Ctx) error {
	user := middleware.GetUser(c)

	var branchID int64
	if user.BranchID != nil {
		branchID = *user.BranchID
	} else if bid := c.Query("branch_id"); bid != "" {
		branchID, _ = strconv.ParseInt(bid, 10, 64)
	}

	if branchID == 0 {
		return response.BadRequest(c, "Branch ID is required")
	}

	items, err := h.itemService.GetPendingDeliveries(c.Context(), branchID)
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.OK(c, items)
}

// RegisterRoutes registers item routes
func (h *ItemHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	items := app.Group("/items")
	items.Use(authMiddleware.Authenticate())

	items.Get("/", authMiddleware.RequirePermission("items.read"), h.List)
	items.Post("/", authMiddleware.RequirePermission("items.create"), h.Create)
	items.Get("/for-sale", authMiddleware.RequirePermission("items.read"), h.GetForSale)
	items.Get("/pending-deliveries", authMiddleware.RequirePermission("items.read"), h.GetPendingDeliveries)
	items.Get("/sku/:sku", authMiddleware.RequirePermission("items.read"), h.GetBySKU)
	items.Get("/:id", authMiddleware.RequirePermission("items.read"), h.GetByID)
	items.Put("/:id", authMiddleware.RequirePermission("items.update"), h.Update)
	items.Delete("/:id", authMiddleware.RequirePermission("items.delete"), h.Delete)
	items.Post("/:id/status", authMiddleware.RequirePermission("items.update"), h.UpdateStatus)
	items.Post("/:id/mark-for-sale", authMiddleware.RequirePermission("items.update"), h.MarkForSale)
	items.Post("/:id/mark-as-delivered", authMiddleware.RequirePermission("items.update"), h.MarkAsDelivered)
}
