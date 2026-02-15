package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"pawnshop/internal/domain"
	"pawnshop/internal/middleware"
	"pawnshop/internal/repository"
	"pawnshop/internal/service"
	"pawnshop/pkg/response"
	"pawnshop/pkg/validator"
)

// CashHandler handles cash/POS endpoints
type CashHandler struct {
	cashService *service.CashService
}

// NewCashHandler creates a new CashHandler
func NewCashHandler(cashService *service.CashService) *CashHandler {
	return &CashHandler{cashService: cashService}
}

// === Cash Register Endpoints ===

// CreateRegister handles cash register creation
func (h *CashHandler) CreateRegister(c *fiber.Ctx) error {
	var input service.CreateRegisterInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	// Set branch from user if not provided
	user := middleware.GetUser(c)
	if input.BranchID == 0 && user.BranchID != nil {
		input.BranchID = *user.BranchID
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	register, err := h.cashService.CreateRegister(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, register)
}

// GetRegister handles getting a cash register by ID
func (h *CashHandler) GetRegister(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid register ID format")
	}

	register, err := h.cashService.GetRegister(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Register not found")
	}

	return response.OK(c, register)
}

// ListRegisters handles listing cash registers
func (h *CashHandler) ListRegisters(c *fiber.Ctx) error {
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

	registers, err := h.cashService.ListRegisters(c.Context(), branchID)
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.OK(c, registers)
}

// UpdateRegister handles cash register update
func (h *CashHandler) UpdateRegister(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid register ID")
	}

	var input service.UpdateRegisterInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	register, err := h.cashService.UpdateRegister(c.Context(), id, input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, register)
}

// === Cash Session Endpoints ===

// OpenSession handles opening a cash session
func (h *CashHandler) OpenSession(c *fiber.Ctx) error {
	var input service.OpenSessionInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	user := middleware.GetUser(c)
	if input.BranchID == 0 && user.BranchID != nil {
		input.BranchID = *user.BranchID
	}
	input.UserID = user.ID

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	session, err := h.cashService.OpenSession(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, session)
}

// GetSession handles getting a cash session by ID
func (h *CashHandler) GetSession(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid session ID format")
	}

	session, err := h.cashService.GetSession(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Session not found")
	}

	return response.OK(c, session)
}

// GetCurrentSession handles getting the current open session for the user
func (h *CashHandler) GetCurrentSession(c *fiber.Ctx) error {
	user := middleware.GetUser(c)

	session, err := h.cashService.GetCurrentSession(c.Context(), user.ID)
	if err != nil {
		return response.NotFound(c, "Current session not found")
	}

	return response.OK(c, session)
}

// ListSessions handles listing cash sessions
func (h *CashHandler) ListSessions(c *fiber.Ctx) error {
	user := middleware.GetUser(c)

	params := repository.CashSessionListParams{
		PaginationParams: repository.PaginationParams{
			Page:    c.QueryInt("page", 1),
			PerPage: c.QueryInt("per_page", 20),
			OrderBy: c.Query("order_by", "opened_at"),
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
	if userID := c.Query("user_id"); userID != "" {
		id, _ := strconv.ParseInt(userID, 10, 64)
		params.UserID = &id
	}
	if registerID := c.Query("register_id"); registerID != "" {
		id, _ := strconv.ParseInt(registerID, 10, 64)
		params.RegisterID = &id
	}
	if status := c.Query("status"); status != "" {
		s := domain.CashSessionStatus(status)
		params.Status = &s
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		params.DateFrom = &dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		params.DateTo = &dateTo
	}

	result, err := h.cashService.ListSessions(c.Context(), params)
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.Paginated(c, result.Data, result.Page, result.PerPage, result.Total)
}

// CloseSession handles closing a cash session
func (h *CashHandler) CloseSession(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid session ID")
	}

	var input struct {
		ClosingAmount float64 `json:"closing_amount" validate:"gte=0"`
		ClosingNotes  *string `json:"closing_notes"`
	}
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	user := middleware.GetUser(c)
	session, err := h.cashService.CloseSession(c.Context(), service.CloseSessionInput{
		SessionID:     id,
		ClosingAmount: input.ClosingAmount,
		ClosingNotes:  input.ClosingNotes,
		ClosedBy:      user.ID,
	})
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, session)
}

// GetSessionSummary handles getting session summary
func (h *CashHandler) GetSessionSummary(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid session ID")
	}

	summary, err := h.cashService.GetSessionSummary(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Session summary not found")
	}

	return response.OK(c, summary)
}

// === Cash Movement Endpoints ===

// CreateMovement handles cash movement creation
func (h *CashHandler) CreateMovement(c *fiber.Ctx) error {
	var input service.CreateMovementInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	input.CreatedBy = middleware.GetUser(c).ID

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	movement, err := h.cashService.CreateMovement(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, movement)
}

// GetMovement handles getting a cash movement by ID
func (h *CashHandler) GetMovement(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid movement ID format")
	}

	movement, err := h.cashService.GetMovement(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Movement not found")
	}

	return response.OK(c, movement)
}

// ListMovements handles listing cash movements
func (h *CashHandler) ListMovements(c *fiber.Ctx) error {
	user := middleware.GetUser(c)

	params := repository.CashMovementListParams{
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
	if sessionID := c.Query("session_id"); sessionID != "" {
		id, _ := strconv.ParseInt(sessionID, 10, 64)
		params.SessionID = &id
	}
	if movementType := c.Query("type"); movementType != "" {
		t := domain.CashMovementType(movementType)
		params.MovementType = &t
	}
	if method := c.Query("method"); method != "" {
		m := domain.PaymentMethod(method)
		params.PaymentMethod = &m
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		params.DateFrom = &dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		params.DateTo = &dateTo
	}

	result, err := h.cashService.ListMovements(c.Context(), params)
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.Paginated(c, result.Data, result.Page, result.PerPage, result.Total)
}

// ListSessionMovements handles listing movements for a session
func (h *CashHandler) ListSessionMovements(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid session ID")
	}

	movements, err := h.cashService.ListSessionMovements(c.Context(), id)
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.OK(c, movements)
}

// RegisterRoutes registers cash/POS routes
func (h *CashHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	cash := app.Group("/cash")
	cash.Use(authMiddleware.Authenticate())

	// Cash registers
	registers := cash.Group("/registers")
	registers.Get("/", authMiddleware.RequirePermission("cash.read"), h.ListRegisters)
	registers.Post("/", authMiddleware.RequirePermission("cash.create"), h.CreateRegister)
	registers.Get("/:id", authMiddleware.RequirePermission("cash.read"), h.GetRegister)
	registers.Put("/:id", authMiddleware.RequirePermission("cash.update"), h.UpdateRegister)

	// Cash sessions
	sessions := cash.Group("/sessions")
	sessions.Get("/", authMiddleware.RequirePermission("cash.read"), h.ListSessions)
	sessions.Post("/", authMiddleware.RequirePermission("cash.create"), h.OpenSession)
	sessions.Get("/current", authMiddleware.RequirePermission("cash.read"), h.GetCurrentSession)
	sessions.Get("/:id", authMiddleware.RequirePermission("cash.read"), h.GetSession)
	sessions.Post("/:id/close", authMiddleware.RequirePermission("cash.update"), h.CloseSession)
	sessions.Get("/:id/summary", authMiddleware.RequirePermission("cash.read"), h.GetSessionSummary)
	sessions.Get("/:id/movements", authMiddleware.RequirePermission("cash.read"), h.ListSessionMovements)

	// Cash movements
	movements := cash.Group("/movements")
	movements.Get("/", authMiddleware.RequirePermission("cash.read"), h.ListMovements)
	movements.Post("/", authMiddleware.RequirePermission("cash.create"), h.CreateMovement)
	movements.Get("/:id", authMiddleware.RequirePermission("cash.read"), h.GetMovement)
}
