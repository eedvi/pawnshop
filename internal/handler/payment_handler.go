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

// PaymentHandler handles payment endpoints
type PaymentHandler struct {
	paymentService *service.PaymentService
	auditLogger    *middleware.AuditLogger
	logger         zerolog.Logger
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(paymentService *service.PaymentService, auditLogger *middleware.AuditLogger, logger zerolog.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		auditLogger:    auditLogger,
		logger:         logger.With().Str("handler", "payment").Logger(),
	}
}

// Create handles payment creation
func (h *PaymentHandler) Create(c *fiber.Ctx) error {
	log := logger.FromContext(c.UserContext(), h.logger)

	var input service.CreatePaymentInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Msg("Failed to parse payment request body")
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	// Set branch and user from context
	user := middleware.GetUser(c)
	if user.BranchID != nil {
		input.BranchID = *user.BranchID
	}
	input.CreatedBy = user.ID

	// Validate
	if errors := validator.Validate(&input); errors != nil {
		log.Warn().Interface("validation_errors", errors).Msg("Payment validation failed")
		return response.ValidationError(c, errors)
	}

	// Service layer handles detailed logging
	result, err := h.paymentService.Create(c.Context(), input)
	if err != nil {
		// Service already logged the error
		return response.BadRequest(c, err.Error())
	}

	log.Info().
		Int64("payment_id", result.Payment.ID).
		Str("payment_number", result.Payment.PaymentNumber).
		Int64("loan_id", input.LoanID).
		Float64("amount", input.Amount).
		Bool("fully_paid", result.IsFullyPaid).
		Msg("Payment created successfully at handler level")

	// Audit log
	if h.auditLogger != nil {
		description := fmt.Sprintf("Pago #%s creado por Q%.2f en préstamo #%d", result.Payment.PaymentNumber, input.Amount, input.LoanID)
		if result.IsFullyPaid {
			description += " (préstamo totalmente pagado)"
		}
		h.auditLogger.LogCreateWithDescription(c, "payment", result.Payment.ID, description, fiber.Map{
			"payment_number": result.Payment.PaymentNumber,
			"loan_id":        input.LoanID,
			"amount":         input.Amount,
			"payment_method": input.PaymentMethod,
			"fully_paid":     result.IsFullyPaid,
		})
	}

	return response.Created(c, result)
}

// GetByID handles getting a payment by ID
func (h *PaymentHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid payment ID format")
	}

	payment, err := h.paymentService.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Payment not found")
	}

	return response.OK(c, payment)
}

// List handles listing payments
func (h *PaymentHandler) List(c *fiber.Ctx) error {
	user := middleware.GetUser(c)

	params := repository.PaymentListParams{
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
	if loanID := c.Query("loan_id"); loanID != "" {
		id, _ := strconv.ParseInt(loanID, 10, 64)
		params.LoanID = &id
	}
	if status := c.Query("status"); status != "" {
		s := domain.PaymentStatus(status)
		params.Status = &s
	}
	if method := c.Query("method"); method != "" {
		m := domain.PaymentMethod(method)
		params.Method = &m
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		params.DateFrom = &dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		params.DateTo = &dateTo
	}

	result, err := h.paymentService.List(c.Context(), params)
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.Paginated(c, result.Data, result.Page, result.PerPage, result.Total)
}

// Reverse handles payment reversal
func (h *PaymentHandler) Reverse(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid payment ID")
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

	user := middleware.GetUser(c)

	// Get payment before reversing for audit
	originalPayment, _ := h.paymentService.GetByID(c.Context(), id)

	payment, err := h.paymentService.Reverse(c.Context(), service.ReversePaymentInput{
		PaymentID:  id,
		Reason:     input.Reason,
		ReversedBy: user.ID,
	})
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Audit log
	if h.auditLogger != nil && originalPayment != nil {
		description := fmt.Sprintf("Pago #%s reversado por Q%.2f. Razón: %s", originalPayment.PaymentNumber, originalPayment.Amount, input.Reason)
		h.auditLogger.LogCustomAction(c, "reverse", "payment", id, description,
			fiber.Map{
				"status":         originalPayment.Status,
				"amount":         originalPayment.Amount,
				"payment_number": originalPayment.PaymentNumber,
			},
			fiber.Map{
				"status":      payment.Status,
				"reversed_at": payment.ReversedAt,
				"reason":      input.Reason,
			})
	}

	return response.OK(c, payment)
}

// CalculatePayoff handles calculating loan payoff
func (h *PaymentHandler) CalculatePayoff(c *fiber.Ctx) error {
	loanID, err := strconv.ParseInt(c.Query("loan_id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid loan ID")
	}

	amount, err := h.paymentService.CalculatePayoff(c.Context(), loanID)
	if err != nil {
		return response.NotFound(c, "Loan not found")
	}

	return response.OK(c, fiber.Map{
		"loan_id":       loanID,
		"payoff_amount": amount,
	})
}

// CalculateMinimum handles calculating minimum payment
func (h *PaymentHandler) CalculateMinimum(c *fiber.Ctx) error {
	loanID, err := strconv.ParseInt(c.Query("loan_id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid loan ID")
	}

	amount, err := h.paymentService.CalculateMinimumPayment(c.Context(), loanID)
	if err != nil {
		return response.NotFound(c, "Loan not found")
	}

	return response.OK(c, fiber.Map{
		"loan_id":         loanID,
		"minimum_payment": amount,
	})
}

// RegisterRoutes registers payment routes
func (h *PaymentHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	payments := app.Group("/payments")
	payments.Use(authMiddleware.Authenticate())

	payments.Get("/", authMiddleware.RequirePermission("payments.read"), h.List)
	payments.Post("/", authMiddleware.RequirePermission("payments.create"), h.Create)
	payments.Get("/calculate-payoff", authMiddleware.RequirePermission("payments.read"), h.CalculatePayoff)
	payments.Get("/calculate-minimum", authMiddleware.RequirePermission("payments.read"), h.CalculateMinimum)
	payments.Get("/:id", authMiddleware.RequirePermission("payments.read"), h.GetByID)
	payments.Post("/:id/reverse", authMiddleware.RequirePermission("payments.update"), h.Reverse)
}
