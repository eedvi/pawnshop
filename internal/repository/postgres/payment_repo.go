package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// PaymentRepository implements repository.PaymentRepository
type PaymentRepository struct {
	db *DB
}

// NewPaymentRepository creates a new PaymentRepository
func NewPaymentRepository(db *DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// GetByID retrieves a payment by ID
func (r *PaymentRepository) GetByID(ctx context.Context, id int64) (*domain.Payment, error) {
	query := `
		SELECT id, payment_number, branch_id, loan_id, customer_id,
			   amount, principal_amount, interest_amount, late_fee_amount,
			   payment_method, reference_number, status, payment_date,
			   loan_balance_after, interest_balance_after,
			   reversed_at, reversed_by, reversal_reason, notes, cash_session_id,
			   created_by, created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	return r.scanPayment(r.db.QueryRowContext(ctx, query, id))
}

// GetByNumber retrieves a payment by number
func (r *PaymentRepository) GetByNumber(ctx context.Context, paymentNumber string) (*domain.Payment, error) {
	query := `
		SELECT id, payment_number, branch_id, loan_id, customer_id,
			   amount, principal_amount, interest_amount, late_fee_amount,
			   payment_method, reference_number, status, payment_date,
			   loan_balance_after, interest_balance_after,
			   reversed_at, reversed_by, reversal_reason, notes, cash_session_id,
			   created_by, created_at, updated_at
		FROM payments
		WHERE payment_number = $1
	`

	return r.scanPayment(r.db.QueryRowContext(ctx, query, paymentNumber))
}

// List retrieves payments with pagination and filters
func (r *PaymentRepository) List(ctx context.Context, params repository.PaymentListParams) (*repository.PaginatedResult[domain.Payment], error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 20
	}

	// Base query with JOINs for related entities
	baseQuery := `
		FROM payments p
		LEFT JOIN loans l ON p.loan_id = l.id
		LEFT JOIN customers c ON p.customer_id = c.id
		WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if params.BranchID > 0 {
		argCount++
		baseQuery += fmt.Sprintf(" AND p.branch_id = $%d", argCount)
		args = append(args, params.BranchID)
	}

	if params.CustomerID != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND p.customer_id = $%d", argCount)
		args = append(args, *params.CustomerID)
	}

	if params.LoanID != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND p.loan_id = $%d", argCount)
		args = append(args, *params.LoanID)
	}

	if params.Status != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND p.status = $%d", argCount)
		args = append(args, *params.Status)
	}

	if params.Method != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND p.payment_method = $%d", argCount)
		args = append(args, *params.Method)
	}

	if params.DateFrom != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND p.payment_date >= $%d", argCount)
		args = append(args, *params.DateFrom)
	}

	if params.DateTo != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND p.payment_date <= $%d", argCount)
		args = append(args, *params.DateTo)
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count payments: %w", err)
	}

	// Get data
	orderBy := "p.created_at"
	if params.OrderBy != "" {
		orderBy = "p." + params.OrderBy
	}
	order := "DESC"
	if params.Order == "asc" {
		order = "ASC"
	}

	offset := (params.Page - 1) * params.PerPage
	dataQuery := fmt.Sprintf(`
		SELECT p.id, p.payment_number, p.branch_id, p.loan_id, p.customer_id,
			   p.amount, p.principal_amount, p.interest_amount, p.late_fee_amount,
			   p.payment_method, p.reference_number, p.status, p.payment_date,
			   p.loan_balance_after, p.interest_balance_after,
			   p.reversed_at, p.reversed_by, p.reversal_reason, p.notes, p.cash_session_id,
			   p.created_by, p.created_at, p.updated_at,
			   l.id, l.loan_number, l.status,
			   c.id, c.first_name, c.last_name, c.identity_number
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		baseQuery, orderBy, order, argCount+1, argCount+2,
	)
	args = append(args, params.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	payments := []domain.Payment{}
	for rows.Next() {
		payment, err := r.scanPaymentRowWithRelations(rows)
		if err != nil {
			return nil, err
		}
		payments = append(payments, *payment)
	}

	totalPages := total / params.PerPage
	if total%params.PerPage > 0 {
		totalPages++
	}

	return &repository.PaginatedResult[domain.Payment]{
		Data:       payments,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

// ListByLoan retrieves all payments for a loan
func (r *PaymentRepository) ListByLoan(ctx context.Context, loanID int64) ([]*domain.Payment, error) {
	query := `
		SELECT id, payment_number, branch_id, loan_id, customer_id,
			   amount, principal_amount, interest_amount, late_fee_amount,
			   payment_method, reference_number, status, payment_date,
			   loan_balance_after, interest_balance_after,
			   reversed_at, reversed_by, reversal_reason, notes, cash_session_id,
			   created_by, created_at, updated_at
		FROM payments
		WHERE loan_id = $1
		ORDER BY payment_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	payments := []*domain.Payment{}
	for rows.Next() {
		payment, err := r.scanPaymentRow(rows)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

// Create creates a new payment
func (r *PaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (
			payment_number, branch_id, loan_id, customer_id,
			amount, principal_amount, interest_amount, late_fee_amount,
			payment_method, reference_number, status, payment_date,
			loan_balance_after, interest_balance_after, notes, cash_session_id, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		payment.PaymentNumber, payment.BranchID, payment.LoanID, payment.CustomerID,
		payment.Amount, payment.PrincipalAmount, payment.InterestAmount, payment.LateFeeAmount,
		payment.PaymentMethod, NullString(payment.ReferenceNumber), payment.Status, payment.PaymentDate,
		payment.LoanBalanceAfter, payment.InterestBalanceAfter,
		NullString(payment.Notes), NullInt64(payment.CashSessionID), payment.CreatedBy,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

// Update updates an existing payment
func (r *PaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	query := `
		UPDATE payments SET
			status = $2, reversed_at = $3, reversed_by = $4, reversal_reason = $5, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		payment.ID, payment.Status,
		NullTime(payment.ReversedAt), NullInt64(payment.ReversedBy), NullString(payment.ReversalReason),
	)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// GenerateNumber generates a unique payment number
func (r *PaymentRepository) GenerateNumber(ctx context.Context) (string, error) {
	yearStr := time.Now().Format("2006")
	var seqNum int

	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(payment_number FROM 'PY-\d{4}-(\d+)') AS INTEGER)), 0) + 1
		FROM payments WHERE payment_number LIKE 'PY-' || $1 || '-%'
	`
	r.db.QueryRowContext(ctx, query, yearStr).Scan(&seqNum)
	if seqNum == 0 {
		seqNum = 1
	}

	return fmt.Sprintf("PY-%s-%06d", yearStr, seqNum), nil
}

// Helper functions
func (r *PaymentRepository) scanPayment(row *sql.Row) (*domain.Payment, error) {
	p := &domain.Payment{}
	var referenceNumber, reversalReason, notes sql.NullString
	var reversedAt sql.NullTime
	var reversedBy, cashSessionID sql.NullInt64
	var createdBy sql.NullInt64

	err := row.Scan(
		&p.ID, &p.PaymentNumber, &p.BranchID, &p.LoanID, &p.CustomerID,
		&p.Amount, &p.PrincipalAmount, &p.InterestAmount, &p.LateFeeAmount,
		&p.PaymentMethod, &referenceNumber, &p.Status, &p.PaymentDate,
		&p.LoanBalanceAfter, &p.InterestBalanceAfter,
		&reversedAt, &reversedBy, &reversalReason, &notes, &cashSessionID,
		&createdBy, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	p.ReferenceNumber = StringPtr(referenceNumber)
	p.ReversalReason = StringPtr(reversalReason)
	p.Notes = StringPtr(notes)
	p.ReversedAt = TimePtr(reversedAt)
	p.ReversedBy = Int64Ptr(reversedBy)
	p.CashSessionID = Int64Ptr(cashSessionID)
	if createdBy.Valid {
		p.CreatedBy = createdBy.Int64
	}

	return p, nil
}

func (r *PaymentRepository) scanPaymentRow(rows *sql.Rows) (*domain.Payment, error) {
	p := &domain.Payment{}
	var referenceNumber, reversalReason, notes sql.NullString
	var reversedAt sql.NullTime
	var reversedBy, cashSessionID sql.NullInt64
	var createdBy sql.NullInt64

	err := rows.Scan(
		&p.ID, &p.PaymentNumber, &p.BranchID, &p.LoanID, &p.CustomerID,
		&p.Amount, &p.PrincipalAmount, &p.InterestAmount, &p.LateFeeAmount,
		&p.PaymentMethod, &referenceNumber, &p.Status, &p.PaymentDate,
		&p.LoanBalanceAfter, &p.InterestBalanceAfter,
		&reversedAt, &reversedBy, &reversalReason, &notes, &cashSessionID,
		&createdBy, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan payment: %w", err)
	}

	p.ReferenceNumber = StringPtr(referenceNumber)
	p.ReversalReason = StringPtr(reversalReason)
	p.Notes = StringPtr(notes)
	p.ReversedAt = TimePtr(reversedAt)
	p.ReversedBy = Int64Ptr(reversedBy)
	p.CashSessionID = Int64Ptr(cashSessionID)
	if createdBy.Valid {
		p.CreatedBy = createdBy.Int64
	}

	return p, nil
}

func (r *PaymentRepository) scanPaymentRowWithRelations(rows *sql.Rows) (*domain.Payment, error) {
	p := &domain.Payment{}
	var referenceNumber, reversalReason, notes sql.NullString
	var reversedAt sql.NullTime
	var reversedBy, cashSessionID sql.NullInt64
	var createdBy sql.NullInt64

	// Loan fields
	var loanID sql.NullInt64
	var loanNumber sql.NullString
	var loanStatus sql.NullString

	// Customer fields
	var custID sql.NullInt64
	var custFirstName, custLastName, custIdentityNumber sql.NullString

	err := rows.Scan(
		&p.ID, &p.PaymentNumber, &p.BranchID, &p.LoanID, &p.CustomerID,
		&p.Amount, &p.PrincipalAmount, &p.InterestAmount, &p.LateFeeAmount,
		&p.PaymentMethod, &referenceNumber, &p.Status, &p.PaymentDate,
		&p.LoanBalanceAfter, &p.InterestBalanceAfter,
		&reversedAt, &reversedBy, &reversalReason, &notes, &cashSessionID,
		&createdBy, &p.CreatedAt, &p.UpdatedAt,
		// Loan
		&loanID, &loanNumber, &loanStatus,
		// Customer
		&custID, &custFirstName, &custLastName, &custIdentityNumber,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan payment with relations: %w", err)
	}

	p.ReferenceNumber = StringPtr(referenceNumber)
	p.ReversalReason = StringPtr(reversalReason)
	p.Notes = StringPtr(notes)
	p.ReversedAt = TimePtr(reversedAt)
	p.ReversedBy = Int64Ptr(reversedBy)
	p.CashSessionID = Int64Ptr(cashSessionID)
	if createdBy.Valid {
		p.CreatedBy = createdBy.Int64
	}

	// Populate Loan relation
	if loanID.Valid {
		p.Loan = &domain.Loan{
			ID:         loanID.Int64,
			LoanNumber: loanNumber.String,
			Status:     domain.LoanStatus(loanStatus.String),
		}
	}

	// Populate Customer relation
	if custID.Valid {
		p.Customer = &domain.Customer{
			ID:             custID.Int64,
			FirstName:      custFirstName.String,
			LastName:       custLastName.String,
			IdentityNumber: custIdentityNumber.String,
		}
	}

	return p, nil
}
