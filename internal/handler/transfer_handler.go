package handler

import (
	"strconv"

	"pawnshop/internal/middleware"
	"pawnshop/internal/repository"
	"pawnshop/internal/service"
	"github.com/gofiber/fiber/v2"
)

type TransferHandler struct {
	transferService service.TransferService
}

func NewTransferHandler(transferService service.TransferService) *TransferHandler {
	return &TransferHandler{
		transferService: transferService,
	}
}

// Create creates a new item transfer
// @Summary Create a new item transfer
// @Tags Transfers
// @Accept json
// @Produce json
// @Param transfer body service.CreateTransferRequest true "Transfer data"
// @Success 201 {object} domain.ItemTransfer
// @Router /api/v1/transfers [post]
func (h *TransferHandler) Create(c *fiber.Ctx) error {
	var req service.CreateTransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	// Get user ID from context
	userID := c.Locals("userID").(int64)
	req.RequestedBy = userID

	transfer, err := h.transferService.Create(c.Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(transfer)
}

// GetByID retrieves a transfer by ID
// @Summary Get a transfer by ID
// @Tags Transfers
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} domain.ItemTransfer
// @Router /api/v1/transfers/{id} [get]
func (h *TransferHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transfer ID format",
		})
	}

	transfer, err := h.transferService.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(transfer)
}

// GetByNumber retrieves a transfer by number
// @Summary Get a transfer by number
// @Tags Transfers
// @Produce json
// @Param number path string true "Transfer Number"
// @Success 200 {object} domain.ItemTransfer
// @Router /api/v1/transfers/number/{number} [get]
func (h *TransferHandler) GetByNumber(c *fiber.Ctx) error {
	number := c.Params("number")

	transfer, err := h.transferService.GetByNumber(c.Context(), number)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(transfer)
}

// List retrieves transfers with filtering
// @Summary List transfers
// @Tags Transfers
// @Produce json
// @Param status query string false "Filter by status"
// @Param from_branch_id query int false "Filter by source branch"
// @Param to_branch_id query int false "Filter by destination branch"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/transfers [get]
func (h *TransferHandler) List(c *fiber.Ctx) error {
	filter := repository.TransferFilter{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 20),
	}

	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}
	if fromBranchID := c.QueryInt("from_branch_id"); fromBranchID > 0 {
		id := int64(fromBranchID)
		filter.FromBranchID = &id
	}
	if toBranchID := c.QueryInt("to_branch_id"); toBranchID > 0 {
		id := int64(toBranchID)
		filter.ToBranchID = &id
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filter.DateFrom = &dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		filter.DateTo = &dateTo
	}

	transfers, total, err := h.transferService.List(c.Context(), filter)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"data":  transfers,
		"total": total,
		"page":  filter.Page,
		"size":  filter.PageSize,
	})
}

// ListByBranch retrieves transfers for a branch
// @Summary List transfers for a branch
// @Tags Transfers
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Param status query string false "Filter by status"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/branches/{branch_id}/transfers [get]
func (h *TransferHandler) ListByBranch(c *fiber.Ctx) error {
	branchID, err := strconv.ParseInt(c.Params("branch_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid branch ID format",
		})
	}

	filter := repository.TransferFilter{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 20),
	}

	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	transfers, total, err := h.transferService.ListByBranch(c.Context(), branchID, filter)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"data":  transfers,
		"total": total,
		"page":  filter.Page,
		"size":  filter.PageSize,
	})
}

// Approve approves a pending transfer
// @Summary Approve a transfer
// @Tags Transfers
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Param body body object{notes string} false "Approval notes"
// @Success 200 {object} domain.ItemTransfer
// @Router /api/v1/transfers/{id}/approve [post]
func (h *TransferHandler) Approve(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transfer ID format",
		})
	}

	var body struct {
		Notes string `json:"notes"`
	}
	c.BodyParser(&body)

	userID := c.Locals("userID").(int64)

	transfer, err := h.transferService.Approve(c.Context(), id, userID, body.Notes)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(transfer)
}

// Ship marks a transfer as shipped
// @Summary Ship a transfer
// @Tags Transfers
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} domain.ItemTransfer
// @Router /api/v1/transfers/{id}/ship [post]
func (h *TransferHandler) Ship(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transfer ID format",
		})
	}

	userID := c.Locals("userID").(int64)

	transfer, err := h.transferService.Ship(c.Context(), id, userID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(transfer)
}

// Receive marks a transfer as received
// @Summary Receive a transfer
// @Tags Transfers
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Param body body object{notes string} false "Receipt notes"
// @Success 200 {object} domain.ItemTransfer
// @Router /api/v1/transfers/{id}/receive [post]
func (h *TransferHandler) Receive(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transfer ID format",
		})
	}

	var body struct {
		Notes string `json:"notes"`
	}
	c.BodyParser(&body)

	userID := c.Locals("userID").(int64)

	transfer, err := h.transferService.Receive(c.Context(), id, userID, body.Notes)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(transfer)
}

// Cancel cancels a transfer
// @Summary Cancel a transfer
// @Tags Transfers
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Param body body object{reason string} true "Cancellation reason"
// @Success 200 {object} domain.ItemTransfer
// @Router /api/v1/transfers/{id}/cancel [post]
func (h *TransferHandler) Cancel(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transfer ID format",
		})
	}

	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID := c.Locals("userID").(int64)

	transfer, err := h.transferService.Cancel(c.Context(), id, userID, body.Reason)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(transfer)
}

// GetPendingForBranch retrieves pending transfers for a branch
// @Summary Get pending transfers for a branch
// @Tags Transfers
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Success 200 {array} domain.ItemTransfer
// @Router /api/v1/branches/{branch_id}/transfers/pending [get]
func (h *TransferHandler) GetPendingForBranch(c *fiber.Ctx) error {
	branchID, err := strconv.ParseInt(c.Params("branch_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid branch ID format",
		})
	}

	transfers, err := h.transferService.GetPendingForBranch(c.Context(), branchID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(transfers)
}

// GetInTransitForBranch retrieves in-transit transfers for a branch
// @Summary Get in-transit transfers for a branch
// @Tags Transfers
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Success 200 {array} domain.ItemTransfer
// @Router /api/v1/branches/{branch_id}/transfers/in-transit [get]
func (h *TransferHandler) GetInTransitForBranch(c *fiber.Ctx) error {
	branchID, err := strconv.ParseInt(c.Params("branch_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid branch ID format",
		})
	}

	transfers, err := h.transferService.GetInTransitForBranch(c.Context(), branchID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(transfers)
}

// RegisterRoutes registers transfer routes
func (h *TransferHandler) RegisterRoutes(router fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	transfers := router.Group("/transfers")
	transfers.Use(authMiddleware.Authenticate())

	transfers.Post("/", authMiddleware.RequirePermission("transfers:create"), h.Create)
	transfers.Get("/", authMiddleware.RequirePermission("transfers:read"), h.List)
	transfers.Get("/:id", authMiddleware.RequirePermission("transfers:read"), h.GetByID)
	transfers.Get("/number/:number", authMiddleware.RequirePermission("transfers:read"), h.GetByNumber)
	transfers.Post("/:id/approve", authMiddleware.RequirePermission("transfers:approve"), h.Approve)
	transfers.Post("/:id/ship", authMiddleware.RequirePermission("transfers:ship"), h.Ship)
	transfers.Post("/:id/receive", authMiddleware.RequirePermission("transfers:receive"), h.Receive)
	transfers.Post("/:id/cancel", authMiddleware.RequirePermission("transfers:cancel"), h.Cancel)

	// Branch-specific transfer routes
	branchTransfers := router.Group("/branches/:branch_id/transfers")
	branchTransfers.Use(authMiddleware.Authenticate())
	branchTransfers.Get("/", authMiddleware.RequirePermission("transfers:read"), h.ListByBranch)
	branchTransfers.Get("/pending", authMiddleware.RequirePermission("transfers:read"), h.GetPendingForBranch)
	branchTransfers.Get("/in-transit", authMiddleware.RequirePermission("transfers:read"), h.GetInTransitForBranch)
}
