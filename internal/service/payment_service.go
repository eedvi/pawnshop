package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// PaymentService handles payment business logic
type PaymentService struct {
	paymentRepo  repository.PaymentRepository
	loanRepo     repository.LoanRepository
	customerRepo repository.CustomerRepository
	itemRepo     repository.ItemRepository
	logger       zerolog.Logger
}

// NewPaymentService creates a new PaymentService
func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	loanRepo repository.LoanRepository,
	customerRepo repository.CustomerRepository,
	itemRepo repository.ItemRepository,
	logger zerolog.Logger,
) *PaymentService {
	return &PaymentService{
		paymentRepo:  paymentRepo,
		loanRepo:     loanRepo,
		customerRepo: customerRepo,
		itemRepo:     itemRepo,
		logger:       logger.With().Str("service", "payment").Logger(),
	}
}

// CreatePaymentInput represents create payment request data
type CreatePaymentInput struct {
	LoanID          int64   `json:"loan_id" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	PaymentMethod   string  `json:"payment_method" validate:"required,oneof=cash card transfer check other"`
	ReferenceNumber string  `json:"reference_number"`
	Notes           string  `json:"notes"`
	CashSessionID   *int64  `json:"cash_session_id"`
	BranchID        int64   `json:"-"`
	CreatedBy       int64   `json:"-"`
}

// PaymentResult contains the result of a payment
type PaymentResult struct {
	Payment          *domain.Payment `json:"payment"`
	Loan             *domain.Loan    `json:"loan"`
	IsFullyPaid      bool            `json:"is_fully_paid"`
	RemainingBalance float64         `json:"remaining_balance"`
}

// Create creates a new payment and applies it to the loan
func (s *PaymentService) Create(ctx context.Context, input CreatePaymentInput) (*PaymentResult, error) {
	s.logger.Info().
		Int64("loan_id", input.LoanID).
		Float64("amount", input.Amount).
		Str("payment_method", input.PaymentMethod).
		Int64("created_by", input.CreatedBy).
		Msg("Processing payment")

	// Get loan
	loan, err := s.loanRepo.GetByID(ctx, input.LoanID)
	if err != nil {
		s.logger.Error().Err(err).Int64("loan_id", input.LoanID).Msg("Loan not found")
		return nil, errors.New("loan not found")
	}

	// Validate loan can receive payments
	if loan.Status == domain.LoanStatusPaid {
		s.logger.Warn().Int64("loan_id", input.LoanID).Msg("Payment rejected: loan already fully paid")
		return nil, errors.New("loan is already fully paid")
	}
	if loan.Status == domain.LoanStatusConfiscated {
		s.logger.Warn().Int64("loan_id", input.LoanID).Msg("Payment rejected: loan confiscated")
		return nil, errors.New("loan has been confiscated")
	}

	// Calculate total amount owed (prevent overpayment)
	totalOwed := loan.PrincipalRemaining + loan.InterestRemaining + loan.LateFeeAmount
	if input.Amount > totalOwed {
		s.logger.Warn().
			Int64("loan_id", input.LoanID).
			Float64("payment_amount", input.Amount).
			Float64("total_owed", totalOwed).
			Msg("Payment rejected: amount exceeds total owed")
		return nil, fmt.Errorf("payment amount (Q%.2f) exceeds total owed (Q%.2f)", input.Amount, totalOwed)
	}

	// Calculate how to apply the payment
	// Order: Late fees -> Interest -> Principal
	remainingPayment := input.Amount
	lateFeePayment := 0.0
	interestPayment := 0.0
	principalPayment := 0.0

	// Apply to late fees first
	if loan.LateFeeAmount > 0 && remainingPayment > 0 {
		if remainingPayment >= loan.LateFeeAmount {
			lateFeePayment = loan.LateFeeAmount
			remainingPayment -= lateFeePayment
		} else {
			lateFeePayment = remainingPayment
			remainingPayment = 0
		}
	}

	// Apply to interest
	if loan.InterestRemaining > 0 && remainingPayment > 0 {
		if remainingPayment >= loan.InterestRemaining {
			interestPayment = loan.InterestRemaining
			remainingPayment -= interestPayment
		} else {
			interestPayment = remainingPayment
			remainingPayment = 0
		}
	}

	// Apply to principal
	if loan.PrincipalRemaining > 0 && remainingPayment > 0 {
		if remainingPayment >= loan.PrincipalRemaining {
			principalPayment = loan.PrincipalRemaining
			remainingPayment -= principalPayment
		} else {
			principalPayment = remainingPayment
			remainingPayment = 0
		}
	}

	// Update loan balances
	loan.LateFeeAmount -= lateFeePayment
	loan.InterestRemaining -= interestPayment
	loan.PrincipalRemaining -= principalPayment
	loan.AmountPaid += input.Amount
	loan.UpdatedBy = &input.CreatedBy

	// Check if loan is fully paid
	isFullyPaid := loan.PrincipalRemaining == 0 && loan.InterestRemaining == 0 && loan.LateFeeAmount == 0
	if isFullyPaid {
		loan.Status = domain.LoanStatusPaid
		now := time.Now()
		loan.PaidDate = &now
	}

	// Generate payment number
	paymentNumber, err := s.paymentRepo.GenerateNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate payment number: %w", err)
	}

	// Create payment record
	payment := &domain.Payment{
		PaymentNumber:        paymentNumber,
		BranchID:             input.BranchID,
		LoanID:               loan.ID,
		CustomerID:           loan.CustomerID,
		Amount:               input.Amount,
		PrincipalAmount:      principalPayment,
		InterestAmount:       interestPayment,
		LateFeeAmount:        lateFeePayment,
		PaymentMethod:        domain.PaymentMethod(input.PaymentMethod),
		ReferenceNumber:      input.ReferenceNumber,
		Status:               domain.PaymentStatusCompleted,
		PaymentDate:          time.Now(),
		LoanBalanceAfter:     loan.PrincipalRemaining,
		InterestBalanceAfter: loan.InterestRemaining,
		Notes:                input.Notes,
		CashSessionID:        input.CashSessionID,
		CreatedBy:            input.CreatedBy,
	}

	// Save payment
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Update loan
	if err := s.loanRepo.Update(ctx, loan); err != nil {
		return nil, fmt.Errorf("failed to update loan: %w", err)
	}

	// Update item status to available if loan is fully paid
	if isFullyPaid {
		if err := s.itemRepo.UpdateStatus(ctx, loan.ItemID, domain.ItemStatusAvailable); err != nil {
			s.logger.Error().Err(err).Int64("item_id", loan.ItemID).Msg("Failed to update item status to available")
			// Don't fail the payment, but log the error
		} else {
			s.logger.Info().Int64("item_id", loan.ItemID).Msg("Item returned to customer (status: available)")
		}
	}

	// Apply payment to installments if loan has installment payment plan
	if loan.PaymentPlanType == "installments" {
		if err := s.applyPaymentToInstallments(ctx, loan, input.Amount); err != nil {
			// Log error but don't fail the payment
			// Payment has already been recorded
		}
	}

	// Update customer total_paid stats with every payment
	customer, _ := s.customerRepo.GetByID(ctx, loan.CustomerID)
	if customer != nil {
		totalPaid := customer.TotalPaid + input.Amount
		s.customerRepo.UpdateCreditInfo(ctx, customer.ID, repository.CustomerCreditUpdate{
			TotalPaid: &totalPaid,
		})
	}

	s.logger.Info().
		Int64("payment_id", payment.ID).
		Str("payment_number", payment.PaymentNumber).
		Int64("loan_id", loan.ID).
		Float64("amount", input.Amount).
		Float64("principal_paid", principalPayment).
		Float64("interest_paid", interestPayment).
		Float64("late_fee_paid", lateFeePayment).
		Float64("remaining_balance", loan.RemainingBalance()).
		Bool("fully_paid", isFullyPaid).
		Msg("Payment processed successfully")

	return &PaymentResult{
		Payment:          payment,
		Loan:             loan,
		IsFullyPaid:      isFullyPaid,
		RemainingBalance: loan.RemainingBalance(),
	}, nil
}

// GetByID retrieves a payment by ID
func (s *PaymentService) GetByID(ctx context.Context, id int64) (*domain.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load loan info
	payment.Loan, _ = s.loanRepo.GetByID(ctx, payment.LoanID)
	payment.Customer, _ = s.customerRepo.GetByID(ctx, payment.CustomerID)

	return payment, nil
}

// List retrieves payments with pagination and filters
func (s *PaymentService) List(ctx context.Context, params repository.PaymentListParams) (*repository.PaginatedResult[domain.Payment], error) {
	// Repository now loads relations via JOIN - no N+1 queries
	return s.paymentRepo.List(ctx, params)
}

// ReversePaymentInput represents reverse payment request data
type ReversePaymentInput struct {
	PaymentID  int64  `json:"payment_id" validate:"required"`
	Reason     string `json:"reason" validate:"required"`
	ReversedBy int64  `json:"-"`
}

// Reverse reverses a payment
func (s *PaymentService) Reverse(ctx context.Context, input ReversePaymentInput) (*domain.Payment, error) {
	// Get payment
	payment, err := s.paymentRepo.GetByID(ctx, input.PaymentID)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	// Validate payment can be reversed
	if !payment.CanBeReversed() {
		return nil, errors.New("payment cannot be reversed")
	}

	// Get loan
	loan, err := s.loanRepo.GetByID(ctx, payment.LoanID)
	if err != nil {
		return nil, errors.New("loan not found")
	}

	// Check if loan was paid off before reversal
	wasPaid := loan.Status == domain.LoanStatusPaid

	// If loan was paid off, reactivate it
	if wasPaid {
		loan.Status = domain.LoanStatusActive
		loan.PaidDate = nil
	}

	// Reverse the payment amounts
	loan.PrincipalRemaining += payment.PrincipalAmount
	loan.InterestRemaining += payment.InterestAmount
	loan.LateFeeAmount += payment.LateFeeAmount
	loan.AmountPaid -= payment.Amount
	loan.UpdatedBy = &input.ReversedBy

	// Update loan
	if err := s.loanRepo.Update(ctx, loan); err != nil {
		return nil, fmt.Errorf("failed to update loan: %w", err)
	}

	// Update item status back to collateral if loan was paid and is being reactivated
	if wasPaid {
		if err := s.itemRepo.UpdateStatus(ctx, loan.ItemID, domain.ItemStatusCollateral); err != nil {
			s.logger.Error().Err(err).Int64("item_id", loan.ItemID).Msg("Failed to update item status to collateral")
			// Don't fail the reversal, but log the error
		} else {
			s.logger.Info().Int64("item_id", loan.ItemID).Msg("Item returned to collateral status due to payment reversal")
		}
	}

	// Reverse payment from installments if loan has installment payment plan
	if loan.PaymentPlanType == "installments" {
		if err := s.reversePaymentFromInstallments(ctx, loan, payment.Amount); err != nil {
			// Log error but don't fail the reversal
		}
	}

	// Update customer total_paid stats (subtract reversed amount)
	customer, _ := s.customerRepo.GetByID(ctx, payment.CustomerID)
	if customer != nil {
		totalPaid := customer.TotalPaid - payment.Amount
		if totalPaid < 0 {
			totalPaid = 0
		}
		s.customerRepo.UpdateCreditInfo(ctx, customer.ID, repository.CustomerCreditUpdate{
			TotalPaid: &totalPaid,
		})
	}

	// Mark payment as reversed
	now := time.Now()
	payment.Status = domain.PaymentStatusReversed
	payment.ReversedAt = &now
	payment.ReversedBy = &input.ReversedBy
	payment.ReversalReason = input.Reason

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	return payment, nil
}

// CalculatePayoff calculates the payoff amount for a loan
func (s *PaymentService) CalculatePayoff(ctx context.Context, loanID int64) (float64, error) {
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return 0, errors.New("loan not found")
	}

	return loan.RemainingBalance(), nil
}

// CalculatePayoffDetailed calculates the payoff amount and returns loan details
func (s *PaymentService) CalculatePayoffDetailed(ctx context.Context, loanID int64) (float64, *domain.Loan, error) {
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return 0, nil, errors.New("loan not found")
	}

	return loan.RemainingBalance(), loan, nil
}

// CalculateMinimumPayment calculates the minimum payment due for a loan
func (s *PaymentService) CalculateMinimumPayment(ctx context.Context, loanID int64) (float64, error) {
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return 0, errors.New("loan not found")
	}

	if !loan.RequiresMinimumPayment || loan.MinimumPaymentAmount == nil {
		return loan.RemainingBalance(), nil
	}

	// Minimum is the lesser of the minimum payment amount or remaining balance
	minimumPayment := *loan.MinimumPaymentAmount
	if loan.RemainingBalance() < minimumPayment {
		return loan.RemainingBalance(), nil
	}

	// Add any late fees
	return minimumPayment + loan.LateFeeAmount, nil
}

// CalculateMinimumPaymentDetailed calculates the minimum payment and returns loan details
func (s *PaymentService) CalculateMinimumPaymentDetailed(ctx context.Context, loanID int64) (float64, *domain.Loan, error) {
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return 0, nil, errors.New("loan not found")
	}

	if !loan.RequiresMinimumPayment || loan.MinimumPaymentAmount == nil {
		return loan.RemainingBalance(), loan, nil
	}

	// Minimum is the lesser of the minimum payment amount or remaining balance
	minimumPayment := *loan.MinimumPaymentAmount
	if loan.RemainingBalance() < minimumPayment {
		return loan.RemainingBalance(), loan, nil
	}

	// Add any late fees
	return minimumPayment + loan.LateFeeAmount, loan, nil
}

// applyPaymentToInstallments applies a payment amount to loan installments
func (s *PaymentService) applyPaymentToInstallments(ctx context.Context, loan *domain.Loan, paymentAmount float64) error {
	// Get all installments for the loan
	installments, err := s.loanRepo.GetInstallments(ctx, loan.ID)
	if err != nil {
		return fmt.Errorf("failed to get installments: %w", err)
	}

	if len(installments) == 0 {
		return nil // No installments to apply payment to
	}

	remainingPayment := paymentAmount

	// Apply payment to installments in order (unpaid first)
	for _, installment := range installments {
		if remainingPayment <= 0 {
			break
		}

		// Skip already fully paid installments
		if installment.IsPaid {
			continue
		}

		// Calculate remaining balance for this installment
		installmentBalance := installment.TotalAmount - installment.AmountPaid

		// Determine how much to apply to this installment
		paymentToApply := remainingPayment
		if paymentToApply > installmentBalance {
			paymentToApply = installmentBalance
		}

		// Update installment
		installment.AmountPaid += paymentToApply
		remainingPayment -= paymentToApply

		// Mark as paid if fully paid
		if installment.AmountPaid >= installment.TotalAmount {
			installment.IsPaid = true
			now := time.Now()
			installment.PaidDate = &now
		}

		// Update installment in database
		if err := s.loanRepo.UpdateInstallment(ctx, installment); err != nil {
			return fmt.Errorf("failed to update installment: %w", err)
		}
	}

	return nil
}

// reversePaymentFromInstallments reverses a payment from loan installments
func (s *PaymentService) reversePaymentFromInstallments(ctx context.Context, loan *domain.Loan, paymentAmount float64) error {
	// Get all installments for the loan
	installments, err := s.loanRepo.GetInstallments(ctx, loan.ID)
	if err != nil {
		return fmt.Errorf("failed to get installments: %w", err)
	}

	if len(installments) == 0 {
		return nil // No installments to reverse payment from
	}

	remainingReversal := paymentAmount

	// Reverse payment from installments in reverse order (last paid first)
	for i := len(installments) - 1; i >= 0; i-- {
		installment := installments[i]

		if remainingReversal <= 0 {
			break
		}

		// Skip installments with no payments
		if installment.AmountPaid <= 0 {
			continue
		}

		// Determine how much to reverse from this installment
		reversalAmount := remainingReversal
		if reversalAmount > installment.AmountPaid {
			reversalAmount = installment.AmountPaid
		}

		// Update installment
		installment.AmountPaid -= reversalAmount
		remainingReversal -= reversalAmount

		// Mark as unpaid if no longer fully paid
		if installment.AmountPaid < installment.TotalAmount {
			installment.IsPaid = false
			installment.PaidDate = nil
		}

		// Update installment in database
		if err := s.loanRepo.UpdateInstallment(ctx, installment); err != nil {
			return fmt.Errorf("failed to update installment: %w", err)
		}
	}

	return nil
}
