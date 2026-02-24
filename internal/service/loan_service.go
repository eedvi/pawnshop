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

// LoanService handles loan business logic
type LoanService struct {
	loanRepo     repository.LoanRepository
	itemRepo     repository.ItemRepository
	customerRepo repository.CustomerRepository
	paymentRepo  repository.PaymentRepository
	logger       zerolog.Logger
}

// NewLoanService creates a new LoanService
func NewLoanService(
	loanRepo repository.LoanRepository,
	itemRepo repository.ItemRepository,
	customerRepo repository.CustomerRepository,
	paymentRepo repository.PaymentRepository,
	logger zerolog.Logger,
) *LoanService {
	return &LoanService{
		loanRepo:     loanRepo,
		itemRepo:     itemRepo,
		customerRepo: customerRepo,
		paymentRepo:  paymentRepo,
		logger:       logger.With().Str("service", "loan").Logger(),
	}
}

// CreateLoanInput represents create loan request data
type CreateLoanInput struct {
	CustomerID             int64   `json:"customer_id" validate:"required"`
	ItemID                 int64   `json:"item_id" validate:"required"`
	BranchID               int64   `json:"branch_id" validate:"required"`
	LoanAmount             float64 `json:"loan_amount" validate:"required,gt=0"`
	InterestRate           float64 `json:"interest_rate" validate:"required,gte=0,lte=100"`
	LoanTermDays           int     `json:"loan_term_days" validate:"required,gt=0"`
	PaymentPlanType        string  `json:"payment_plan_type" validate:"required,oneof=single minimum_payment installments"`
	RequiresMinimumPayment bool    `json:"requires_minimum_payment"`
	MinimumPaymentAmount   float64 `json:"minimum_payment_amount" validate:"gte=0"`
	GracePeriodDays        int     `json:"grace_period_days" validate:"gte=0,lte=30"`
	NumberOfInstallments   int     `json:"number_of_installments" validate:"gte=0"`
	LateFeeRate            float64 `json:"late_fee_rate" validate:"gte=0"`
	Notes                  string  `json:"notes"`
	CreatedBy              int64   `json:"-"`
}

// Create creates a new loan
func (s *LoanService) Create(ctx context.Context, input CreateLoanInput) (*domain.Loan, error) {
	s.logger.Info().
		Int64("customer_id", input.CustomerID).
		Int64("item_id", input.ItemID).
		Float64("loan_amount", input.LoanAmount).
		Float64("interest_rate", input.InterestRate).
		Int("loan_term_days", input.LoanTermDays).
		Str("payment_plan_type", input.PaymentPlanType).
		Int64("created_by", input.CreatedBy).
		Msg("Creating new loan")

	// Validate customer exists and can take loan
	customer, err := s.customerRepo.GetByID(ctx, input.CustomerID)
	if err != nil {
		s.logger.Error().Err(err).Int64("customer_id", input.CustomerID).Msg("Customer not found")
		return nil, errors.New("customer not found")
	}
	if !customer.CanTakeLoan() {
		s.logger.Warn().
			Int64("customer_id", input.CustomerID).
			Bool("is_active", customer.IsActive).
			Bool("is_blocked", customer.IsBlocked).
			Msg("Loan rejected: customer cannot take loans")
		return nil, errors.New("customer cannot take loans")
	}

	// Validate item exists and is available
	item, err := s.itemRepo.GetByID(ctx, input.ItemID)
	if err != nil {
		s.logger.Error().Err(err).Int64("item_id", input.ItemID).Msg("Item not found")
		return nil, errors.New("item not found")
	}
	if !item.IsAvailable() {
		s.logger.Warn().
			Int64("item_id", input.ItemID).
			Str("status", string(item.Status)).
			Msg("Loan rejected: item not available")
		return nil, errors.New("item is not available for loan")
	}

	// Validate loan amount doesn't exceed item loan value
	if input.LoanAmount > item.LoanValue {
		s.logger.Warn().
			Int64("item_id", input.ItemID).
			Float64("requested_amount", input.LoanAmount).
			Float64("max_loan_value", item.LoanValue).
			Msg("Loan rejected: amount exceeds item loan value")
		return nil, errors.New("loan amount cannot exceed item loan value")
	}

	// Calculate interest
	interestAmount := input.LoanAmount * (input.InterestRate / 100)
	totalAmount := input.LoanAmount + interestAmount

	// Generate loan number
	loanNumber, err := s.loanRepo.GenerateNumber(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to generate loan number")
		return nil, fmt.Errorf("failed to generate loan number: %w", err)
	}

	// Calculate due date and minimum payment info
	startDate := time.Now()
	var dueDate time.Time
	var loanTermDays int

	// For installment payment plans, due date is the last installment date
	if input.PaymentPlanType == "installments" && input.NumberOfInstallments > 0 {
		dueDate = startDate.AddDate(0, input.NumberOfInstallments, 0)
		// Calculate actual term in days based on installments
		loanTermDays = int(dueDate.Sub(startDate).Hours() / 24)
	} else {
		dueDate = startDate.AddDate(0, 0, input.LoanTermDays)
		loanTermDays = input.LoanTermDays
	}

	var minimumPaymentAmount *float64
	var nextPaymentDueDate *time.Time
	if input.RequiresMinimumPayment && input.MinimumPaymentAmount > 0 {
		minimumPaymentAmount = &input.MinimumPaymentAmount
		next := startDate.AddDate(0, 1, 0) // Monthly payment
		nextPaymentDueDate = &next
	}

	// Create loan
	loan := &domain.Loan{
		LoanNumber:             loanNumber,
		BranchID:               input.BranchID,
		CustomerID:             input.CustomerID,
		ItemID:                 input.ItemID,
		LoanAmount:             input.LoanAmount,
		InterestRate:           input.InterestRate,
		InterestAmount:         interestAmount,
		PrincipalRemaining:     input.LoanAmount,
		InterestRemaining:      interestAmount,
		TotalAmount:            totalAmount,
		LateFeeRate:            input.LateFeeRate,
		StartDate:              startDate,
		DueDate:                dueDate,
		PaymentPlanType:        domain.PaymentPlanType(input.PaymentPlanType),
		LoanTermDays:           loanTermDays,
		RequiresMinimumPayment: input.RequiresMinimumPayment,
		MinimumPaymentAmount:   minimumPaymentAmount,
		NextPaymentDueDate:     nextPaymentDueDate,
		GracePeriodDays:        input.GracePeriodDays,
		Status:                 domain.LoanStatusActive,
		Notes:                  input.Notes,
		CreatedBy:              input.CreatedBy,
	}

	// Start transaction
	tx, err := s.loanRepo.BeginTx(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to start transaction")
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create loan
	if err := s.loanRepo.CreateTx(ctx, tx, loan); err != nil {
		s.logger.Error().Err(err).Str("loan_number", loanNumber).Msg("Failed to create loan")
		return nil, fmt.Errorf("failed to create loan: %w", err)
	}

	// Update item status to collateral
	if err := s.itemRepo.UpdateStatus(ctx, item.ID, domain.ItemStatusCollateral); err != nil {
		s.logger.Error().Err(err).Int64("item_id", item.ID).Msg("Failed to update item status")
		return nil, fmt.Errorf("failed to update item status: %w", err)
	}

	// Create installments if applicable
	if input.PaymentPlanType == "installments" && input.NumberOfInstallments > 0 {
		installments := s.calculateInstallments(loan, input.NumberOfInstallments)
		if err := s.loanRepo.CreateInstallmentsTx(ctx, tx, installments); err != nil {
			s.logger.Error().Err(err).
				Str("loan_number", loanNumber).
				Int("num_installments", input.NumberOfInstallments).
				Msg("Failed to create installments")
			return nil, fmt.Errorf("failed to create installments: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error().Err(err).Str("loan_number", loanNumber).Msg("Failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Update customer stats
	totalLoans := customer.TotalLoans + 1
	s.customerRepo.UpdateCreditInfo(ctx, customer.ID, repository.CustomerCreditUpdate{
		TotalLoans: &totalLoans,
	})

	// Load relations
	loan.Customer = customer
	loan.Item = item

	s.logger.Info().
		Int64("loan_id", loan.ID).
		Str("loan_number", loan.LoanNumber).
		Int64("customer_id", input.CustomerID).
		Int64("item_id", input.ItemID).
		Float64("loan_amount", input.LoanAmount).
		Float64("interest_amount", interestAmount).
		Float64("total_amount", totalAmount).
		Str("due_date", dueDate.Format("2006-01-02")).
		Msg("Loan created successfully")

	return loan, nil
}

// calculateInstallments calculates installments for a loan
func (s *LoanService) calculateInstallments(loan *domain.Loan, numInstallments int) []*domain.LoanInstallment {
	installments := make([]*domain.LoanInstallment, numInstallments)

	principalPerInstallment := loan.LoanAmount / float64(numInstallments)
	interestPerInstallment := loan.InterestAmount / float64(numInstallments)
	totalPerInstallment := principalPerInstallment + interestPerInstallment

	for i := 0; i < numInstallments; i++ {
		dueDate := loan.StartDate.AddDate(0, i+1, 0)
		installments[i] = &domain.LoanInstallment{
			LoanID:            loan.ID,
			InstallmentNumber: i + 1,
			DueDate:           dueDate,
			PrincipalAmount:   principalPerInstallment,
			InterestAmount:    interestPerInstallment,
			TotalAmount:       totalPerInstallment,
		}
	}

	return installments
}

// LoanCalculation represents the result of a loan calculation
type LoanCalculation struct {
	LoanAmount        float64                   `json:"loan_amount"`
	InterestRate      float64                   `json:"interest_rate"`
	InterestAmount    float64                   `json:"interest_amount"`
	TotalAmount       float64                   `json:"total_amount"`
	InstallmentAmount float64                   `json:"installment_amount,omitempty"`
	Installments      []*domain.LoanInstallment `json:"installments,omitempty"`
}

// Calculate calculates loan terms without creating the loan (preview)
func (s *LoanService) Calculate(ctx context.Context, input CreateLoanInput) (*LoanCalculation, error) {
	// Validate item exists and check loan value
	item, err := s.itemRepo.GetByID(ctx, input.ItemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	// Validate loan amount doesn't exceed item loan value
	if input.LoanAmount > item.LoanValue {
		return nil, fmt.Errorf("loan amount cannot exceed item loan value (max: %.2f)", item.LoanValue)
	}

	// Calculate interest
	interestAmount := input.LoanAmount * (input.InterestRate / 100)
	totalAmount := input.LoanAmount + interestAmount

	result := &LoanCalculation{
		LoanAmount:     input.LoanAmount,
		InterestRate:   input.InterestRate,
		InterestAmount: interestAmount,
		TotalAmount:    totalAmount,
	}

	// Calculate installments if applicable
	if input.PaymentPlanType == "installments" && input.NumberOfInstallments > 0 {
		result.InstallmentAmount = totalAmount / float64(input.NumberOfInstallments)

		// Create preview installments (without loan ID)
		startDate := time.Now()
		installments := make([]*domain.LoanInstallment, input.NumberOfInstallments)
		principalPerInstallment := input.LoanAmount / float64(input.NumberOfInstallments)
		interestPerInstallment := interestAmount / float64(input.NumberOfInstallments)

		for i := 0; i < input.NumberOfInstallments; i++ {
			dueDate := startDate.AddDate(0, i+1, 0)
			installments[i] = &domain.LoanInstallment{
				InstallmentNumber: i + 1,
				DueDate:           dueDate,
				PrincipalAmount:   principalPerInstallment,
				InterestAmount:    interestPerInstallment,
				TotalAmount:       principalPerInstallment + interestPerInstallment,
			}
		}
		result.Installments = installments
	}

	return result, nil
}

// GetByID retrieves a loan by ID
func (s *LoanService) GetByID(ctx context.Context, id int64) (*domain.Loan, error) {
	loan, err := s.loanRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load relations
	loan.Customer, _ = s.customerRepo.GetByID(ctx, loan.CustomerID)
	loan.Item, _ = s.itemRepo.GetByID(ctx, loan.ItemID)

	return loan, nil
}

// GetByNumber retrieves a loan by number
func (s *LoanService) GetByNumber(ctx context.Context, loanNumber string) (*domain.Loan, error) {
	loan, err := s.loanRepo.GetByNumber(ctx, loanNumber)
	if err != nil {
		return nil, err
	}

	// Load relations
	loan.Customer, _ = s.customerRepo.GetByID(ctx, loan.CustomerID)
	loan.Item, _ = s.itemRepo.GetByID(ctx, loan.ItemID)

	return loan, nil
}

// List retrieves loans with pagination and filters
func (s *LoanService) List(ctx context.Context, params repository.LoanListParams) (*repository.PaginatedResult[domain.Loan], error) {
	// Repository now loads relations via JOIN - no N+1 queries
	return s.loanRepo.List(ctx, params)
}

// GetPayments retrieves all payments for a loan
func (s *LoanService) GetPayments(ctx context.Context, loanID int64) ([]*domain.Payment, error) {
	return s.paymentRepo.ListByLoan(ctx, loanID)
}

// GetInstallments retrieves installments for a loan
func (s *LoanService) GetInstallments(ctx context.Context, loanID int64) ([]*domain.LoanInstallment, error) {
	return s.loanRepo.GetInstallments(ctx, loanID)
}

// RenewLoanInput represents renew loan request data
type RenewLoanInput struct {
	LoanID          int64   `json:"loan_id" validate:"required"`
	NewTermDays     int     `json:"new_term_days" validate:"required,gt=0"`
	NewInterestRate float64 `json:"new_interest_rate" validate:"gte=0"`
	PayInterest     bool    `json:"pay_interest"`
	UpdatedBy       int64   `json:"-"`
}

// Renew renews an existing loan
func (s *LoanService) Renew(ctx context.Context, input RenewLoanInput) (*domain.Loan, error) {
	// Get original loan
	loan, err := s.loanRepo.GetByID(ctx, input.LoanID)
	if err != nil {
		return nil, errors.New("loan not found")
	}

	if loan.Status != domain.LoanStatusActive && loan.Status != domain.LoanStatusOverdue {
		return nil, errors.New("only active or overdue loans can be renewed")
	}

	// If paying interest, interest must be fully paid
	if input.PayInterest && loan.InterestRemaining > 0 {
		return nil, errors.New("interest must be paid before renewal")
	}

	// Mark old loan as renewed
	loan.Status = domain.LoanStatusRenewed
	loan.UpdatedBy = &input.UpdatedBy
	if err := s.loanRepo.Update(ctx, loan); err != nil {
		return nil, fmt.Errorf("failed to update original loan: %w", err)
	}

	// Calculate new interest
	interestRate := input.NewInterestRate
	if interestRate == 0 {
		interestRate = loan.InterestRate
	}
	newInterestAmount := loan.PrincipalRemaining * (interestRate / 100)

	// Generate new loan number
	loanNumber, err := s.loanRepo.GenerateNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate loan number: %w", err)
	}

	// Create new loan
	newLoan := &domain.Loan{
		LoanNumber:             loanNumber,
		BranchID:               loan.BranchID,
		CustomerID:             loan.CustomerID,
		ItemID:                 loan.ItemID,
		LoanAmount:             loan.PrincipalRemaining,
		InterestRate:           interestRate,
		InterestAmount:         newInterestAmount,
		PrincipalRemaining:     loan.PrincipalRemaining,
		InterestRemaining:      newInterestAmount,
		TotalAmount:            loan.PrincipalRemaining + newInterestAmount,
		LateFeeRate:            loan.LateFeeRate,
		StartDate:              time.Now(),
		DueDate:                time.Now().AddDate(0, 0, input.NewTermDays),
		PaymentPlanType:        loan.PaymentPlanType,
		LoanTermDays:           input.NewTermDays,
		RequiresMinimumPayment: loan.RequiresMinimumPayment,
		MinimumPaymentAmount:   loan.MinimumPaymentAmount,
		GracePeriodDays:        loan.GracePeriodDays,
		Status:                 domain.LoanStatusActive,
		RenewedFromID:          &loan.ID,
		RenewalCount:           loan.RenewalCount + 1,
		CreatedBy:              input.UpdatedBy,
	}

	if err := s.loanRepo.Create(ctx, newLoan); err != nil {
		return nil, fmt.Errorf("failed to create renewed loan: %w", err)
	}

	return newLoan, nil
}

// Confiscate marks a loan as confiscated and updates the item status
func (s *LoanService) Confiscate(ctx context.Context, loanID int64, updatedBy int64, notes string) error {
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return errors.New("loan not found")
	}

	if loan.Status != domain.LoanStatusDefaulted && loan.Status != domain.LoanStatusOverdue {
		return errors.New("only defaulted or overdue loans can be confiscated")
	}

	// Update loan status
	now := time.Now()
	loan.Status = domain.LoanStatusConfiscated
	loan.ConfiscatedDate = &now
	loan.Notes = notes
	loan.UpdatedBy = &updatedBy

	if err := s.loanRepo.Update(ctx, loan); err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
	}

	// Update item status
	if err := s.itemRepo.UpdateStatus(ctx, loan.ItemID, domain.ItemStatusConfiscated); err != nil {
		return fmt.Errorf("failed to update item status: %w", err)
	}

	// Update customer defaulted amount
	customer, _ := s.customerRepo.GetByID(ctx, loan.CustomerID)
	if customer != nil {
		totalDefaulted := customer.TotalDefaulted + loan.RemainingBalance()
		s.customerRepo.UpdateCreditInfo(ctx, customer.ID, repository.CustomerCreditUpdate{
			TotalDefaulted: &totalDefaulted,
		})
	}

	return nil
}

// GetOverdueLoans retrieves overdue loans for a branch
func (s *LoanService) GetOverdueLoans(ctx context.Context, branchID int64) ([]*domain.Loan, error) {
	loans, err := s.loanRepo.GetOverdueLoans(ctx, branchID)
	if err != nil {
		return nil, err
	}

	// Load relations for each loan
	for i := range loans {
		loans[i].Customer, _ = s.customerRepo.GetByID(ctx, loans[i].CustomerID)
		loans[i].Item, _ = s.itemRepo.GetByID(ctx, loans[i].ItemID)
	}

	return loans, nil
}

// UpdateOverdueStatus updates the status of overdue loans
func (s *LoanService) UpdateOverdueStatus(ctx context.Context, branchID int64) error {
	loans, err := s.GetOverdueLoans(ctx, branchID)
	if err != nil {
		return err
	}

	for _, loan := range loans {
		daysOverdue := loan.CalculateDaysOverdue()
		loan.DaysOverdue = daysOverdue

		if loan.IsInGracePeriod() {
			loan.Status = domain.LoanStatusOverdue
		} else {
			loan.Status = domain.LoanStatusDefaulted
		}

		s.loanRepo.Update(ctx, loan)
	}

	return nil
}
