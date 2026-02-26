package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/internal/repository/postgres"
)

// CashService handles cash/POS business logic
type CashService struct {
	registerRepo repository.CashRegisterRepository
	sessionRepo  repository.CashSessionRepository
	movementRepo repository.CashMovementRepository
	branchRepo   repository.BranchRepository
}

// NewCashService creates a new CashService
func NewCashService(
	registerRepo repository.CashRegisterRepository,
	sessionRepo repository.CashSessionRepository,
	movementRepo repository.CashMovementRepository,
	branchRepo repository.BranchRepository,
) *CashService {
	return &CashService{
		registerRepo: registerRepo,
		sessionRepo:  sessionRepo,
		movementRepo: movementRepo,
		branchRepo:   branchRepo,
	}
}

// === Cash Register Methods ===

// CreateRegisterInput represents create register request data
type CreateRegisterInput struct {
	BranchID    int64   `json:"branch_id" validate:"required"`
	Name        string  `json:"name" validate:"required,min=2"`
	Description *string `json:"description"`
}

// CreateRegister creates a new cash register
func (s *CashService) CreateRegister(ctx context.Context, input CreateRegisterInput) (*domain.CashRegister, error) {
	// Validate branch
	_, err := s.branchRepo.GetByID(ctx, input.BranchID)
	if err != nil {
		return nil, errors.New("invalid branch")
	}

	register := &domain.CashRegister{
		BranchID:    input.BranchID,
		Name:        input.Name,
		Description: input.Description,
		IsActive:    true,
	}

	if err := s.registerRepo.Create(ctx, register); err != nil {
		return nil, fmt.Errorf("failed to create cash register: %w", err)
	}

	return register, nil
}

// GetRegister retrieves a cash register by ID
func (s *CashService) GetRegister(ctx context.Context, id int64) (*domain.CashRegister, error) {
	register, err := s.registerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("cash register not found")
	}
	return register, nil
}

// ListRegisters retrieves all cash registers for a branch
func (s *CashService) ListRegisters(ctx context.Context, branchID int64) ([]*domain.CashRegister, error) {
	return s.registerRepo.List(ctx, branchID)
}

// UpdateRegisterInput represents update register request data
type UpdateRegisterInput struct {
	Name        string  `json:"name" validate:"omitempty,min=2"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

// UpdateRegister updates a cash register
func (s *CashService) UpdateRegister(ctx context.Context, id int64, input UpdateRegisterInput) (*domain.CashRegister, error) {
	register, err := s.registerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("cash register not found")
	}

	if input.Name != "" {
		register.Name = input.Name
	}
	if input.Description != nil {
		register.Description = input.Description
	}
	if input.IsActive != nil {
		register.IsActive = *input.IsActive
	}

	if err := s.registerRepo.Update(ctx, register); err != nil {
		return nil, fmt.Errorf("failed to update cash register: %w", err)
	}

	return register, nil
}

// === Cash Session Methods ===

// OpenSessionInput represents open session request data
type OpenSessionInput struct {
	BranchID         int64   `json:"branch_id" validate:"required"`
	CashRegisterID   int64   `json:"cash_register_id" validate:"required"`
	UserID           int64   `json:"-"`
	OpeningAmount    float64 `json:"opening_amount" validate:"gte=0"`
	OpeningNotes     *string `json:"opening_notes"`
}

// OpenSession opens a new cash session
func (s *CashService) OpenSession(ctx context.Context, input OpenSessionInput) (*domain.CashSession, error) {
	// Validate register
	register, err := s.registerRepo.GetByID(ctx, input.CashRegisterID)
	if err != nil {
		return nil, errors.New("cash register not found")
	}

	if !register.IsActive {
		return nil, errors.New("cash register is not active")
	}

	if register.BranchID != input.BranchID {
		return nil, errors.New("cash register does not belong to this branch")
	}

	// Check if user already has an open session
	existingSession, _ := s.sessionRepo.GetOpenSession(ctx, input.UserID)
	if existingSession != nil {
		return nil, errors.New("user already has an open cash session")
	}

	// Check if register already has an open session
	registerSession, _ := s.sessionRepo.GetOpenSessionByRegister(ctx, input.CashRegisterID)
	if registerSession != nil {
		return nil, errors.New("register already has an open session")
	}

	session := &domain.CashSession{
		BranchID:       input.BranchID,
		CashRegisterID: input.CashRegisterID,
		UserID:         input.UserID,
		Status:         domain.CashSessionStatusOpen,
		OpeningAmount:  input.OpeningAmount,
		OpenedAt:       time.Now(),
		OpeningNotes:   input.OpeningNotes,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create cash session: %w", err)
	}

	return session, nil
}

// GetSession retrieves a cash session by ID
func (s *CashService) GetSession(ctx context.Context, id int64) (*domain.CashSession, error) {
	session, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("cash session not found")
	}

	// Load movements
	session.Movements, _ = s.movementRepo.ListBySession(ctx, id)

	return session, nil
}

// GetCurrentSession retrieves the current open session for a user
func (s *CashService) GetCurrentSession(ctx context.Context, userID int64) (*domain.CashSession, error) {
	session, err := s.sessionRepo.GetOpenSession(ctx, userID)
	if err != nil {
		return nil, errors.New("no open cash session found")
	}

	// Load movements
	session.Movements, _ = s.movementRepo.ListBySession(ctx, session.ID)

	return session, nil
}

// ListSessions retrieves cash sessions with filters
func (s *CashService) ListSessions(ctx context.Context, params repository.CashSessionListParams) (*repository.PaginatedResult[domain.CashSession], error) {
	return s.sessionRepo.List(ctx, params)
}

// CloseSessionInput represents close session request data
type CloseSessionInput struct {
	SessionID     int64   `json:"session_id" validate:"required"`
	ClosingAmount float64 `json:"closing_amount" validate:"gte=0"`
	ClosingNotes  *string `json:"closing_notes"`
	ClosedBy      int64   `json:"-"`
}

// CloseSession closes a cash session
func (s *CashService) CloseSession(ctx context.Context, input CloseSessionInput) (*domain.CashSession, error) {
	session, err := s.sessionRepo.GetByID(ctx, input.SessionID)
	if err != nil {
		return nil, errors.New("cash session not found")
	}

	if session.Status != domain.CashSessionStatusOpen {
		return nil, errors.New("cash session is not open")
	}

	// Calculate expected amount based on movements
	movementRepo, ok := s.movementRepo.(*postgres.CashMovementRepository)
	if !ok {
		return nil, errors.New("invalid movement repository")
	}

	summary, err := movementRepo.GetSessionSummary(ctx, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session summary: %w", err)
	}

	expectedAmount := session.OpeningAmount + summary.CashIncome - summary.CashExpense
	difference := input.ClosingAmount - expectedAmount

	closingNotes := ""
	if input.ClosingNotes != nil {
		closingNotes = *input.ClosingNotes
	}

	err = s.sessionRepo.Close(ctx, input.SessionID, repository.CashSessionCloseData{
		ClosingAmount:  input.ClosingAmount,
		ExpectedAmount: expectedAmount,
		Difference:     difference,
		ClosedBy:       input.ClosedBy,
		ClosingNotes:   closingNotes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to close cash session: %w", err)
	}

	// Reload session
	return s.sessionRepo.GetByID(ctx, input.SessionID)
}

// GetSessionSummary retrieves summary for a session
func (s *CashService) GetSessionSummary(ctx context.Context, sessionID int64) (*CashSessionSummaryResult, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, errors.New("cash session not found")
	}

	movementRepo, ok := s.movementRepo.(*postgres.CashMovementRepository)
	if !ok {
		return nil, errors.New("invalid movement repository")
	}

	summary, err := movementRepo.GetSessionSummary(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session summary: %w", err)
	}

	balance, err := movementRepo.GetSessionBalance(ctx, sessionID)
	if err != nil {
		balance = session.OpeningAmount
	}

	expectedCash := session.OpeningAmount + summary.CashIncome - summary.CashExpense

	return &CashSessionSummaryResult{
		Session:        session,
		TotalIncome:    summary.TotalIncome,
		TotalExpense:   summary.TotalExpense,
		CashIncome:     summary.CashIncome,
		CashExpense:    summary.CashExpense,
		TotalMovements: summary.TotalMovements,
		CurrentBalance: balance,
		ExpectedCash:   expectedCash,
	}, nil
}

// CashSessionSummaryResult contains session summary data
type CashSessionSummaryResult struct {
	Session        *domain.CashSession `json:"session"`
	TotalIncome    float64             `json:"total_income"`
	TotalExpense   float64             `json:"total_expense"`
	CashIncome     float64             `json:"cash_income"`
	CashExpense    float64             `json:"cash_expense"`
	TotalMovements int                 `json:"total_movements"`
	CurrentBalance float64             `json:"current_balance"`
	ExpectedCash   float64             `json:"expected_cash"`
}

// === Cash Movement Methods ===

// CreateMovementInput represents create movement request data
type CreateMovementInput struct {
	SessionID     int64   `json:"session_id" validate:"required"`
	MovementType  string  `json:"movement_type" validate:"required,oneof=income expense"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" validate:"required,oneof=cash card transfer check other"`
	ReferenceType *string `json:"reference_type"`
	ReferenceID   *int64  `json:"reference_id"`
	Description   string  `json:"description" validate:"required"`
	CreatedBy     int64   `json:"-"`
}

// CreateMovement creates a new cash movement
func (s *CashService) CreateMovement(ctx context.Context, input CreateMovementInput) (*domain.CashMovement, error) {
	// Get session
	session, err := s.sessionRepo.GetByID(ctx, input.SessionID)
	if err != nil {
		return nil, errors.New("cash session not found")
	}

	if session.Status != domain.CashSessionStatusOpen {
		return nil, errors.New("cash session is not open")
	}

	// Get current balance
	movementRepo, ok := s.movementRepo.(*postgres.CashMovementRepository)
	if !ok {
		return nil, errors.New("invalid movement repository")
	}

	currentBalance, err := movementRepo.GetSessionBalance(ctx, session.ID)
	if err != nil {
		currentBalance = session.OpeningAmount
	}

	// Calculate new balance
	var newBalance float64
	movementType := domain.CashMovementType(input.MovementType)
	if movementType == domain.CashMovementTypeIncome {
		newBalance = currentBalance + input.Amount
	} else {
		newBalance = currentBalance - input.Amount
		// For cash expenses, check sufficient balance
		if input.PaymentMethod == "cash" && newBalance < 0 {
			return nil, errors.New("insufficient cash balance")
		}
	}

	movement := &domain.CashMovement{
		BranchID:      session.BranchID,
		SessionID:     input.SessionID,
		MovementType:  movementType,
		Amount:        input.Amount,
		PaymentMethod: domain.PaymentMethod(input.PaymentMethod),
		ReferenceType: input.ReferenceType,
		ReferenceID:   input.ReferenceID,
		Description:   input.Description,
		BalanceAfter:  newBalance,
		CreatedBy:     input.CreatedBy,
	}

	if err := s.movementRepo.Create(ctx, movement); err != nil {
		return nil, fmt.Errorf("failed to create cash movement: %w", err)
	}

	return movement, nil
}

// GetMovement retrieves a cash movement by ID
func (s *CashService) GetMovement(ctx context.Context, id int64) (*domain.CashMovement, error) {
	movement, err := s.movementRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("cash movement not found")
	}
	return movement, nil
}

// ListMovements retrieves cash movements with filters
func (s *CashService) ListMovements(ctx context.Context, params repository.CashMovementListParams) (*repository.PaginatedResult[domain.CashMovement], error) {
	return s.movementRepo.List(ctx, params)
}

// ListSessionMovements retrieves all movements for a session
func (s *CashService) ListSessionMovements(ctx context.Context, sessionID int64) ([]*domain.CashMovement, error) {
	return s.movementRepo.ListBySession(ctx, sessionID)
}

// RecordPaymentMovement records a movement from a payment
func (s *CashService) RecordPaymentMovement(ctx context.Context, sessionID int64, paymentID int64, amount float64, method string, createdBy int64) error {
	refType := "payment"
	_, err := s.CreateMovement(ctx, CreateMovementInput{
		SessionID:     sessionID,
		MovementType:  "income",
		Amount:        amount,
		PaymentMethod: method,
		ReferenceType: &refType,
		ReferenceID:   &paymentID,
		Description:   "Payment received",
		CreatedBy:     createdBy,
	})
	return err
}

// RecordSaleMovement records a movement from a sale
func (s *CashService) RecordSaleMovement(ctx context.Context, sessionID int64, saleID int64, amount float64, method string, createdBy int64) error {
	refType := "sale"
	_, err := s.CreateMovement(ctx, CreateMovementInput{
		SessionID:     sessionID,
		MovementType:  "income",
		Amount:        amount,
		PaymentMethod: method,
		ReferenceType: &refType,
		ReferenceID:   &saleID,
		Description:   "Sale completed",
		CreatedBy:     createdBy,
	})
	return err
}

// RecordLoanDisbursement records a movement from a loan disbursement
func (s *CashService) RecordLoanDisbursement(ctx context.Context, sessionID int64, loanID int64, amount float64, createdBy int64) error {
	refType := "loan"
	_, err := s.CreateMovement(ctx, CreateMovementInput{
		SessionID:     sessionID,
		MovementType:  "expense",
		Amount:        amount,
		PaymentMethod: "cash",
		ReferenceType: &refType,
		ReferenceID:   &loanID,
		Description:   "Loan disbursement",
		CreatedBy:     createdBy,
	})
	return err
}

// RecordRefundMovement records a movement from a refund
func (s *CashService) RecordRefundMovement(ctx context.Context, sessionID int64, saleID int64, amount float64, method string, createdBy int64) error {
	refType := "sale_refund"
	_, err := s.CreateMovement(ctx, CreateMovementInput{
		SessionID:     sessionID,
		MovementType:  "expense",
		Amount:        amount,
		PaymentMethod: method,
		ReferenceType: &refType,
		ReferenceID:   &saleID,
		Description:   "Sale refund",
		CreatedBy:     createdBy,
	})
	return err
}
