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

// SaleHandler handles sale endpoints
type SaleHandler struct {
	saleService *service.SaleService
	auditLogger *middleware.AuditLogger
}

// NewSaleHandler creates a new SaleHandler
func NewSaleHandler(saleService *service.SaleService, auditLogger *middleware.AuditLogger) *SaleHandler {
	return &SaleHandler{saleService: saleService, auditLogger: auditLogger}
}

// Create handles sale creation
func (h *SaleHandler) Create(c *fiber.Ctx) error {
	var input service.CreateSaleInput
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

	result, err := h.saleService.Create(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil {
		description := fmt.Sprintf("Venta #%s creada por Q%.2f (artículo #%d)", result.Sale.SaleNumber, result.Sale.FinalPrice, input.ItemID)
		h.auditLogger.LogCreateWithDescription(c, "sale", result.Sale.ID, description, fiber.Map{
			"sale_number":    result.Sale.SaleNumber,
			"final_price":    result.Sale.FinalPrice,
			"sale_price":     result.Sale.SalePrice,
			"discount":       result.Sale.DiscountAmount,
			"payment_method": input.PaymentMethod,
			"item_id":        input.ItemID,
			"customer_id":    input.CustomerID,
		})
	}

	return response.Created(c, result)
}

// GetByID handles getting a sale by ID
func (h *SaleHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid sale ID format")
	}

	sale, err := h.saleService.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Sale not found")
	}

	return response.OK(c, sale)
}

// GetByNumber handles getting a sale by number
func (h *SaleHandler) GetByNumber(c *fiber.Ctx) error {
	saleNumber := c.Params("number")

	sale, err := h.saleService.GetByNumber(c.Context(), saleNumber)
	if err != nil {
		return response.NotFound(c, "Sale not found")
	}

	return response.OK(c, sale)
}

// List handles listing sales
func (h *SaleHandler) List(c *fiber.Ctx) error {
	user := middleware.GetUser(c)

	params := repository.SaleListParams{
		PaginationParams: repository.PaginationParams{
			Page:    c.QueryInt("page", 1),
			PerPage: c.QueryInt("per_page", 20),
			OrderBy: c.Query("order_by", "created_at"),
			Order:   c.Query("order", "desc"),
		},
	}

	// Filter by user's branch if not admin
	if user.BranchID != nil {
		params.BranchID = *user.BranchID
	} else if branchID := c.Query("branch_id"); branchID != "" {
		id, _ := strconv.ParseInt(branchID, 10, 64)
		params.BranchID = id
	}

	// Parse optional filters
	if customerID := c.Query("customer_id"); customerID != "" {
		id, _ := strconv.ParseInt(customerID, 10, 64)
		params.CustomerID = &id
	}
	if itemID := c.Query("item_id"); itemID != "" {
		id, _ := strconv.ParseInt(itemID, 10, 64)
		params.ItemID = &id
	}
	if status := c.Query("status"); status != "" {
		s := domain.SaleStatus(status)
		params.Status = &s
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		params.DateFrom = &dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		params.DateTo = &dateTo
	}

	result, err := h.saleService.List(c.Context(), params)
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.Paginated(c, result.Data, result.Page, result.PerPage, result.Total)
}

// Refund handles sale refund
func (h *SaleHandler) Refund(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid sale ID")
	}

	var input struct {
		RefundAmount float64 `json:"refund_amount" validate:"required,gt=0"`
		Reason       string  `json:"reason" validate:"required"`
	}
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	// Get sale before refund for audit
	originalSale, _ := h.saleService.GetByID(c.Context(), id)

	user := middleware.GetUser(c)
	sale, err := h.saleService.Refund(c.Context(), service.RefundSaleInput{
		SaleID:       id,
		RefundAmount: input.RefundAmount,
		Reason:       input.Reason,
		RefundedBy:   user.ID,
	})
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && originalSale != nil {
		description := fmt.Sprintf("Venta #%s reembolsada por Q%.2f. Razón: %s", originalSale.SaleNumber, input.RefundAmount, input.Reason)
		h.auditLogger.LogCustomAction(c, "refund", "sale", id, description,
			fiber.Map{
				"final_price": originalSale.FinalPrice,
				"status":      originalSale.Status,
			},
			fiber.Map{
				"refund_amount": input.RefundAmount,
				"refund_reason": input.Reason,
				"refunded_at":   "now",
			})
	}

	return response.OK(c, sale)
}

// GetSummary handles getting sales summary
func (h *SaleHandler) GetSummary(c *fiber.Ctx) error {
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

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	if dateFrom == "" || dateTo == "" {
		return response.BadRequest(c, "date_from and date_to are required")
	}

	summary, err := h.saleService.GetSalesSummary(c.Context(), branchID, dateFrom, dateTo)
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.OK(c, summary)
}

// RegisterRoutes registers sale routes
func (h *SaleHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	sales := app.Group("/sales")
	sales.Use(authMiddleware.Authenticate())

	sales.Get("/", authMiddleware.RequirePermission("sales.read"), h.List)
	sales.Post("/", authMiddleware.RequirePermission("sales.create"), h.Create)
	sales.Get("/summary", authMiddleware.RequirePermission("sales.read"), h.GetSummary)
	sales.Get("/number/:number", authMiddleware.RequirePermission("sales.read"), h.GetByNumber)
	sales.Get("/:id", authMiddleware.RequirePermission("sales.read"), h.GetByID)
	sales.Post("/:id/refund", authMiddleware.RequirePermission("sales.update"), h.Refund)
}
