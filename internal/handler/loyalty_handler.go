package handler

import (
	"strconv"

	"pawnshop/internal/middleware"
	"pawnshop/internal/service"
	"github.com/gofiber/fiber/v2"
)

type LoyaltyHandler struct {
	loyaltyService service.LoyaltyService
}

func NewLoyaltyHandler(loyaltyService service.LoyaltyService) *LoyaltyHandler {
	return &LoyaltyHandler{
		loyaltyService: loyaltyService,
	}
}

// EnrollCustomer enrolls a customer in the loyalty program
// @Summary Enroll customer in loyalty program
// @Tags Loyalty
// @Accept json
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/customers/{customer_id}/loyalty/enroll [post]
func (h *LoyaltyHandler) EnrollCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.ParseInt(c.Params("customer_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID format",
		})
	}

	if err := h.loyaltyService.EnrollCustomer(c.Context(), customerID); err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Customer enrolled in loyalty program successfully",
	})
}

// GetCustomerLoyalty gets loyalty information for a customer
// @Summary Get customer loyalty info
// @Tags Loyalty
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Success 200 {object} service.CustomerLoyaltyInfo
// @Router /api/v1/customers/{customer_id}/loyalty [get]
func (h *LoyaltyHandler) GetCustomerLoyalty(c *fiber.Ctx) error {
	customerID, err := strconv.ParseInt(c.Params("customer_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID format",
		})
	}

	info, err := h.loyaltyService.GetCustomerLoyalty(c.Context(), customerID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(info)
}

// AddPoints adds loyalty points to a customer
// @Summary Add loyalty points
// @Tags Loyalty
// @Accept json
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Param request body addPointsRequest true "Add points request"
// @Success 200 {object} domain.LoyaltyPointsHistory
// @Router /api/v1/customers/{customer_id}/loyalty/points [post]
func (h *LoyaltyHandler) AddPoints(c *fiber.Ctx) error {
	customerID, err := strconv.ParseInt(c.Params("customer_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID format",
		})
	}

	var req addPointsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	// Get user ID from context
	userID := middleware.GetUserID(c)

	history, err := h.loyaltyService.AddPoints(c.Context(), service.AddPointsRequest{
		CustomerID:    customerID,
		Points:        req.Points,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		Description:   req.Description,
		BranchID:      req.BranchID,
		CreatedBy:     &userID,
	})
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(history)
}

type addPointsRequest struct {
	Points        int    `json:"points" validate:"required,gt=0"`
	ReferenceType string `json:"reference_type"`
	ReferenceID   *int64 `json:"reference_id"`
	Description   string `json:"description"`
	BranchID      *int64 `json:"branch_id"`
}

// RedeemPoints redeems loyalty points for a customer
// @Summary Redeem loyalty points
// @Tags Loyalty
// @Accept json
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Param request body redeemPointsRequest true "Redeem points request"
// @Success 200 {object} domain.LoyaltyPointsHistory
// @Router /api/v1/customers/{customer_id}/loyalty/redeem [post]
func (h *LoyaltyHandler) RedeemPoints(c *fiber.Ctx) error {
	customerID, err := strconv.ParseInt(c.Params("customer_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID format",
		})
	}

	var req redeemPointsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	// Get user ID from context
	userID := middleware.GetUserID(c)

	history, err := h.loyaltyService.RedeemPoints(c.Context(), service.RedeemPointsRequest{
		CustomerID:    customerID,
		Points:        req.Points,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		Description:   req.Description,
		BranchID:      req.BranchID,
		CreatedBy:     &userID,
	})
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(history)
}

type redeemPointsRequest struct {
	Points        int    `json:"points" validate:"required,gt=0"`
	ReferenceType string `json:"reference_type"`
	ReferenceID   *int64 `json:"reference_id"`
	Description   string `json:"description"`
	BranchID      *int64 `json:"branch_id"`
}

// GetPointsHistory gets loyalty points history for a customer
// @Summary Get loyalty points history
// @Tags Loyalty
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} paginatedLoyaltyHistoryResponse
// @Router /api/v1/customers/{customer_id}/loyalty/history [get]
func (h *LoyaltyHandler) GetPointsHistory(c *fiber.Ctx) error {
	customerID, err := strconv.ParseInt(c.Params("customer_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID format",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	history, total, err := h.loyaltyService.GetPointsHistory(c.Context(), customerID, page, pageSize)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"data":      history,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

type paginatedLoyaltyHistoryResponse struct {
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// CalculateDiscount calculates the loyalty discount for an amount
// @Summary Calculate loyalty discount
// @Tags Loyalty
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Param amount query number true "Amount to calculate discount for"
// @Success 200 {object} map[string]float64
// @Router /api/v1/customers/{customer_id}/loyalty/discount [get]
func (h *LoyaltyHandler) CalculateDiscount(c *fiber.Ctx) error {
	customerID, err := strconv.ParseInt(c.Params("customer_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID format",
		})
	}

	amount, err := strconv.ParseFloat(c.Query("amount", "0"), 64)
	if err != nil || amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid amount",
		})
	}

	discount, err := h.loyaltyService.CalculateDiscount(c.Context(), customerID, amount)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"amount":   amount,
		"discount": discount,
		"total":    amount - discount,
	})
}

// RegisterRoutes registers loyalty routes
func (h *LoyaltyHandler) RegisterRoutes(apiRouter fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	loyalty := apiRouter.Group("/customers/:customer_id/loyalty")
	loyalty.Use(authMiddleware.Authenticate())

	loyalty.Get("/", authMiddleware.RequirePermission("customers:read"), h.GetCustomerLoyalty)
	loyalty.Post("/enroll", authMiddleware.RequirePermission("customers:update"), h.EnrollCustomer)
	loyalty.Post("/points", authMiddleware.RequirePermission("loyalty:manage"), h.AddPoints)
	loyalty.Post("/redeem", authMiddleware.RequirePermission("loyalty:manage"), h.RedeemPoints)
	loyalty.Get("/history", authMiddleware.RequirePermission("customers:read"), h.GetPointsHistory)
	loyalty.Get("/discount", authMiddleware.RequirePermission("customers:read"), h.CalculateDiscount)
}
