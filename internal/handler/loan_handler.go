package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"pawnshop/internal/domain"
	"pawnshop/internal/middleware"
	"pawnshop/internal/repository"
	"pawnshop/internal/service"
	"pawnshop/pkg/logger"
	"pawnshop/pkg/response"
	"pawnshop/pkg/validator"
)

// LoanHandler handles loan endpoints
type LoanHandler struct {
	loanService *service.LoanService
	auditLogger *middleware.AuditLogger
	logger      zerolog.Logger
}

// NewLoanHandler creates a new LoanHandler
func NewLoanHandler(loanService *service.LoanService, auditLogger *middleware.AuditLogger, logger zerolog.Logger) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
		auditLogger: auditLogger,
		logger:      logger.With().Str("handler", "loan").Logger(),
	}
}

// Create handles loan creation
func (h *LoanHandler) Create(c *fiber.Ctx) error {
	log := logger.FromContext(c.UserContext(), h.logger)

	var input service.CreateLoanInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Msg("Failed to parse loan request body")
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
		log.Warn().Interface("validation_errors", errors).Msg("Loan validation failed")
		return response.ValidationError(c, errors)
	}

	// Service layer handles detailed logging
	loan, err := h.loanService.Create(c.Context(), input)
	if err != nil {
		// Service already logged the error
		return response.BadRequest(c, err.Error())
	}

	log.Info().
		Int64("loan_id", loan.ID).
		Str("loan_number", loan.LoanNumber).
		Int64("customer_id", input.CustomerID).
		Int64("item_id", input.ItemID).
		Float64("loan_amount", input.LoanAmount).
		Msg("Loan created successfully at handler level")

	// Audit log
	if h.auditLogger != nil {
		description := fmt.Sprintf("Préstamo #%s creado por Q%.2f a %d días", loan.LoanNumber, input.LoanAmount, input.LoanTermDays)
		h.auditLogger.LogCreateWithDescription(c, "loan", loan.ID, description, fiber.Map{
			"loan_number":   loan.LoanNumber,
			"customer_id":   input.CustomerID,
			"item_id":       input.ItemID,
			"loan_amount":   input.LoanAmount,
			"interest_rate": input.InterestRate,
			"term_days":     input.LoanTermDays,
			"due_date":      loan.DueDate,
		})
	}

	return response.Created(c, loan)
}

// Calculate handles loan calculation preview (without creating)
func (h *LoanHandler) Calculate(c *fiber.Ctx) error {
	var input service.CreateLoanInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	// Set branch from context if not provided
	user := middleware.GetUser(c)
	if input.BranchID == 0 && user.BranchID != nil {
		input.BranchID = *user.BranchID
	}

	result, err := h.loanService.Calculate(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, result)
}

// GetByID handles getting a loan by ID
func (h *LoanHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid loan ID format")
	}

	loan, err := h.loanService.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Loan not found")
	}

	return response.OK(c, loan)
}

// GetByNumber handles getting a loan by number
func (h *LoanHandler) GetByNumber(c *fiber.Ctx) error {
	loanNumber := c.Params("number")

	loan, err := h.loanService.GetByNumber(c.Context(), loanNumber)
	if err != nil {
		return response.NotFound(c, "Loan not found")
	}

	return response.OK(c, loan)
}

// List handles listing loans
func (h *LoanHandler) List(c *fiber.Ctx) error {
	user := middleware.GetUser(c)

	params := repository.LoanListParams{
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
	if customerID := c.Query("customer_id"); customerID != "" {
		id, _ := strconv.ParseInt(customerID, 10, 64)
		params.CustomerID = &id
	}
	if itemID := c.Query("item_id"); itemID != "" {
		id, _ := strconv.ParseInt(itemID, 10, 64)
		params.ItemID = &id
	}
	if status := c.Query("status"); status != "" {
		s := domain.LoanStatus(status)
		params.Status = &s
	}
	if dueBefore := c.Query("due_before"); dueBefore != "" {
		params.DueBefore = &dueBefore
	}
	if dueAfter := c.Query("due_after"); dueAfter != "" {
		params.DueAfter = &dueAfter
	}

	result, err := h.loanService.List(c.Context(), params)
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.Paginated(c, result.Data, result.Page, result.PerPage, result.Total)
}

// GetPayments handles getting payments for a loan
func (h *LoanHandler) GetPayments(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid loan ID")
	}

	payments, err := h.loanService.GetPayments(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Loan payments not found")
	}

	return response.OK(c, payments)
}

// GetInstallments handles getting installments for a loan
func (h *LoanHandler) GetInstallments(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid loan ID")
	}

	installments, err := h.loanService.GetInstallments(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Loan installments not found")
	}

	return response.OK(c, installments)
}

// Renew handles loan renewal
func (h *LoanHandler) Renew(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid loan ID")
	}

	var input service.RenewLoanInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	input.LoanID = id
	input.UpdatedBy = middleware.GetUser(c).ID

	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	// Get original loan for audit
	originalLoan, _ := h.loanService.GetByID(c.Context(), id)

	loan, err := h.loanService.Renew(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && originalLoan != nil {
		description := fmt.Sprintf("Préstamo #%s renovado por %d días", originalLoan.LoanNumber, input.NewTermDays)
		if input.PayInterest {
			description += " con pago de interés"
		}
		h.auditLogger.LogCustomAction(c, "renew", "loan", id, description,
			fiber.Map{
				"old_due_date": originalLoan.DueDate,
				"old_status":   originalLoan.Status,
			},
			fiber.Map{
				"new_due_date":      loan.DueDate,
				"new_term_days":     input.NewTermDays,
				"pay_interest":      input.PayInterest,
				"new_interest_rate": input.NewInterestRate,
			})
	}

	return response.OK(c, loan)
}

// Confiscate handles loan confiscation
func (h *LoanHandler) Confiscate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid loan ID")
	}

	var input struct {
		Notes string `json:"notes"`
	}
	c.BodyParser(&input)

	// Get loan before confiscating for audit
	originalLoan, _ := h.loanService.GetByID(c.Context(), id)

	user := middleware.GetUser(c)
	err = h.loanService.Confiscate(c.Context(), id, user.ID, input.Notes)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && originalLoan != nil {
		description := fmt.Sprintf("Préstamo #%s confiscado. Artículo pasó a propiedad de la casa de empeño", originalLoan.LoanNumber)
		if input.Notes != "" {
			description += fmt.Sprintf(". Notas: %s", input.Notes)
		}
		h.auditLogger.LogCustomAction(c, "confiscate", "loan", id, description,
			fiber.Map{
				"status":      originalLoan.Status,
				"item_id":     originalLoan.ItemID,
				"loan_number": originalLoan.LoanNumber,
			},
			fiber.Map{
				"status":         domain.LoanStatusConfiscated,
				"confiscated_at": "now",
				"notes":          input.Notes,
			})
	}

	return response.OK(c, fiber.Map{"message": "Loan confiscated successfully"})
}

// GetOverdue handles getting overdue loans
func (h *LoanHandler) GetOverdue(c *fiber.Ctx) error {
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

	loans, err := h.loanService.GetOverdueLoans(c.Context(), branchID)
	if err != nil {
		return response.InternalErrorWithErr(c, err)
	}

	return response.OK(c, loans)
}

// RegisterRoutes registers loan routes
func (h *LoanHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	loans := app.Group("/loans")
	loans.Use(authMiddleware.Authenticate())

	loans.Get("/", authMiddleware.RequirePermission("loans.read"), h.List)
	loans.Post("/", authMiddleware.RequirePermission("loans.create"), h.Create)
	loans.Post("/calculate", authMiddleware.RequirePermission("loans.read"), h.Calculate)
	loans.Get("/overdue", authMiddleware.RequirePermission("loans.read"), h.GetOverdue)
	loans.Get("/number/:number", authMiddleware.RequirePermission("loans.read"), h.GetByNumber)
	loans.Get("/:id", authMiddleware.RequirePermission("loans.read"), h.GetByID)
	loans.Get("/:id/payments", authMiddleware.RequirePermission("loans.read"), h.GetPayments)
	loans.Get("/:id/installments", authMiddleware.RequirePermission("loans.read"), h.GetInstallments)
	loans.Post("/:id/renew", authMiddleware.RequirePermission("loans.update"), h.Renew)
	loans.Post("/:id/confiscate", authMiddleware.RequirePermission("loans.update"), h.Confiscate)
}
