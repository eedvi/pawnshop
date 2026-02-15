package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"pawnshop/internal/middleware"
	"pawnshop/internal/repository"
	"pawnshop/internal/service"
	"pawnshop/pkg/response"
	"pawnshop/pkg/validator"
)

// CustomerHandler handles customer endpoints
type CustomerHandler struct {
	customerService *service.CustomerService
	auditLogger     *middleware.AuditLogger
}

// NewCustomerHandler creates a new CustomerHandler
func NewCustomerHandler(customerService *service.CustomerService, auditLogger *middleware.AuditLogger) *CustomerHandler {
	return &CustomerHandler{customerService: customerService, auditLogger: auditLogger}
}

// Create handles customer creation
func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	var input service.CreateCustomerInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	// Set branch from user if not provided
	user := middleware.GetUser(c)
	if input.BranchID == 0 && user.BranchID != nil {
		input.BranchID = *user.BranchID
	}
	input.CreatedBy = user.ID

	// Validate
	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	customer, err := h.customerService.Create(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, customer)
}

// GetByID handles getting a customer by ID
func (h *CustomerHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID format")
	}

	customer, err := h.customerService.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Customer not found")
	}

	return response.OK(c, customer)
}

// List handles listing customers
func (h *CustomerHandler) List(c *fiber.Ctx) error {
	user := middleware.GetUser(c)

	params := repository.CustomerListParams{
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
	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		params.IsActive = &active
	}
	if isBlocked := c.Query("is_blocked"); isBlocked != "" {
		blocked := isBlocked == "true"
		params.IsBlocked = &blocked
	}

	result, err := h.customerService.List(c.Context(), params)
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.Paginated(c, result.Data, result.Page, result.PerPage, result.Total)
}

// Update handles customer update
func (h *CustomerHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID format")
	}

	var input service.UpdateCustomerInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	customer, err := h.customerService.Update(c.Context(), id, input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, customer)
}

// Delete handles customer deletion
func (h *CustomerHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID format")
	}

	if err := h.customerService.Delete(c.Context(), id); err != nil {
		return response.NotFound(c, "Customer not found")
	}

	return response.NoContent(c)
}

// Block handles blocking a customer
func (h *CustomerHandler) Block(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID format")
	}

	var input struct {
		Reason string `json:"reason" validate:"required"`
	}
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	err = h.customerService.Block(c.Context(), service.BlockCustomerInput{
		CustomerID: id,
		Reason:     input.Reason,
	})
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{"message": "Customer blocked successfully"})
}

// Unblock handles unblocking a customer
func (h *CustomerHandler) Unblock(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID format")
	}

	if err := h.customerService.Unblock(c.Context(), id); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{"message": "Customer unblocked successfully"})
}

// RegisterRoutes registers customer routes
func (h *CustomerHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	customers := app.Group("/customers")
	customers.Use(authMiddleware.Authenticate())

	customers.Get("/", authMiddleware.RequirePermission("customers.read"), h.List)
	customers.Post("/", authMiddleware.RequirePermission("customers.create"), h.Create)
	customers.Get("/:id", authMiddleware.RequirePermission("customers.read"), h.GetByID)
	customers.Put("/:id", authMiddleware.RequirePermission("customers.update"), h.Update)
	customers.Delete("/:id", authMiddleware.RequirePermission("customers.delete"), h.Delete)
	customers.Post("/:id/block", authMiddleware.RequirePermission("customers.update"), h.Block)
	customers.Post("/:id/unblock", authMiddleware.RequirePermission("customers.update"), h.Unblock)
}
