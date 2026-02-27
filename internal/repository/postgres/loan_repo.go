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

// LoanRepository implements repository.LoanRepository
type LoanRepository struct {
	db *DB
}

// NewLoanRepository creates a new LoanRepository
func NewLoanRepository(db *DB) *LoanRepository {
	return &LoanRepository{db: db}
}

// GetByID retrieves a loan by ID
func (r *LoanRepository) GetByID(ctx context.Context, id int64) (*domain.Loan, error) {
	query := `
		SELECT id, loan_number, branch_id, customer_id, item_id,
			   loan_amount, interest_rate, interest_amount, principal_remaining, interest_remaining,
			   total_amount, amount_paid, late_fee_rate, late_fee_amount, late_fee_remaining,
			   start_date, due_date, paid_date, confiscated_date,
			   payment_plan_type, loan_term_days, requires_minimum_payment,
			   minimum_payment_amount, next_payment_due_date, grace_period_days,
			   number_of_installments, installment_amount,
			   status, days_overdue, renewed_from_id, renewal_count, notes,
			   created_by, updated_by, created_at, updated_at, deleted_at
		FROM loans
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanLoan(r.db.QueryRowContext(ctx, query, id))
}

// GetByNumber retrieves a loan by loan number
func (r *LoanRepository) GetByNumber(ctx context.Context, loanNumber string) (*domain.Loan, error) {
	query := `
		SELECT id, loan_number, branch_id, customer_id, item_id,
			   loan_amount, interest_rate, interest_amount, principal_remaining, interest_remaining,
			   total_amount, amount_paid, late_fee_rate, late_fee_amount, late_fee_remaining,
			   start_date, due_date, paid_date, confiscated_date,
			   payment_plan_type, loan_term_days, requires_minimum_payment,
			   minimum_payment_amount, next_payment_due_date, grace_period_days,
			   number_of_installments, installment_amount,
			   status, days_overdue, renewed_from_id, renewal_count, notes,
			   created_by, updated_by, created_at, updated_at, deleted_at
		FROM loans
		WHERE loan_number = $1 AND deleted_at IS NULL
	`

	return r.scanLoan(r.db.QueryRowContext(ctx, query, loanNumber))
}

// List retrieves loans with pagination and filters
func (r *LoanRepository) List(ctx context.Context, params repository.LoanListParams) (*repository.PaginatedResult[domain.Loan], error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 20
	}

	// Base query with JOINs for related entities
	baseQuery := `
		FROM loans l
		LEFT JOIN customers c ON l.customer_id = c.id
		LEFT JOIN items i ON l.item_id = i.id
		WHERE l.deleted_at IS NULL`
	args := []interface{}{}
	argCount := 0

	if params.BranchID > 0 {
		argCount++
		baseQuery += fmt.Sprintf(" AND l.branch_id = $%d", argCount)
		args = append(args, params.BranchID)
	}

	if params.CustomerID != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND l.customer_id = $%d", argCount)
		args = append(args, *params.CustomerID)
	}

	if params.ItemID != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND l.item_id = $%d", argCount)
		args = append(args, *params.ItemID)
	}

	if params.Status != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND l.status = $%d", argCount)
		args = append(args, *params.Status)
	}

	if params.DueBefore != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND l.due_date <= $%d", argCount)
		args = append(args, *params.DueBefore)
	}

	if params.DueAfter != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND l.due_date >= $%d", argCount)
		args = append(args, *params.DueAfter)
	}

	if params.Search != "" {
		argCount++
		baseQuery += fmt.Sprintf(" AND (l.loan_number ILIKE $%d OR c.first_name ILIKE $%d OR c.last_name ILIKE $%d OR c.identity_number ILIKE $%d OR i.name ILIKE $%d)", argCount, argCount, argCount, argCount, argCount)
		args = append(args, "%"+params.Search+"%")
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count loans: %w", err)
	}

	// Get data
	orderBy := "l.created_at"
	if params.OrderBy != "" {
		orderBy = "l." + params.OrderBy
	}
	order := "DESC"
	if params.Order == "asc" {
		order = "ASC"
	}

	offset := (params.Page - 1) * params.PerPage
	dataQuery := fmt.Sprintf(`
		SELECT l.id, l.loan_number, l.branch_id, l.customer_id, l.item_id,
			   l.loan_amount, l.interest_rate, l.interest_amount, l.principal_remaining, l.interest_remaining,
			   l.total_amount, l.amount_paid, l.late_fee_rate, l.late_fee_amount, l.late_fee_remaining,
			   l.start_date, l.due_date, l.paid_date, l.confiscated_date,
			   l.payment_plan_type, l.loan_term_days, l.requires_minimum_payment,
			   l.minimum_payment_amount, l.next_payment_due_date, l.grace_period_days,
			   l.number_of_installments, l.installment_amount,
			   l.status, l.days_overdue, l.renewed_from_id, l.renewal_count, l.notes,
			   l.created_by, l.updated_by, l.created_at, l.updated_at, l.deleted_at,
			   c.id, c.first_name, c.last_name, c.identity_number,
			   i.id, i.name, i.sku
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		baseQuery, orderBy, order, argCount+1, argCount+2,
	)
	args = append(args, params.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list loans: %w", err)
	}
	defer rows.Close()

	loans := []domain.Loan{}
	for rows.Next() {
		loan, err := r.scanLoanRowWithRelations(rows)
		if err != nil {
			return nil, err
		}
		loans = append(loans, *loan)
	}

	totalPages := total / params.PerPage
	if total%params.PerPage > 0 {
		totalPages++
	}

	return &repository.PaginatedResult[domain.Loan]{
		Data:       loans,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

// Create creates a new loan
func (r *LoanRepository) Create(ctx context.Context, loan *domain.Loan) error {
	query := `
		INSERT INTO loans (
			loan_number, branch_id, customer_id, item_id,
			loan_amount, interest_rate, interest_amount, principal_remaining, interest_remaining,
			total_amount, late_fee_rate,
			start_date, due_date,
			payment_plan_type, loan_term_days, requires_minimum_payment,
			minimum_payment_amount, next_payment_due_date, grace_period_days,
			number_of_installments, installment_amount,
			status, notes, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		loan.LoanNumber, loan.BranchID, loan.CustomerID, loan.ItemID,
		loan.LoanAmount, loan.InterestRate, loan.InterestAmount,
		loan.PrincipalRemaining, loan.InterestRemaining, loan.TotalAmount, loan.LateFeeRate,
		loan.StartDate, loan.DueDate,
		loan.PaymentPlanType, loan.LoanTermDays, loan.RequiresMinimumPayment,
		NullFloat64(loan.MinimumPaymentAmount), NullTime(loan.NextPaymentDueDate), loan.GracePeriodDays,
		loan.NumberOfInstallments, NullFloat64(loan.InstallmentAmount),
		loan.Status, NullString(loan.Notes), loan.CreatedBy,
	).Scan(&loan.ID, &loan.CreatedAt, &loan.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create loan: %w", err)
	}

	return nil
}

// Update updates an existing loan
func (r *LoanRepository) Update(ctx context.Context, loan *domain.Loan) error {
	query := `
		UPDATE loans SET
			interest_amount = $2, principal_remaining = $3, interest_remaining = $4,
			total_amount = $5, amount_paid = $6, late_fee_amount = $7, late_fee_remaining = $8,
			due_date = $9, paid_date = $10, confiscated_date = $11,
			minimum_payment_amount = $12, next_payment_due_date = $13,
			status = $14, days_overdue = $15, renewal_count = $16, notes = $17,
			updated_by = $18, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		loan.ID, loan.InterestAmount, loan.PrincipalRemaining, loan.InterestRemaining,
		loan.TotalAmount, loan.AmountPaid, loan.LateFeeAmount, loan.LateFeeRemaining,
		loan.DueDate, NullTime(loan.PaidDate), NullTime(loan.ConfiscatedDate),
		NullFloat64(loan.MinimumPaymentAmount), NullTime(loan.NextPaymentDueDate),
		loan.Status, loan.DaysOverdue, loan.RenewalCount, NullString(loan.Notes),
		loan.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("loan not found")
	}

	return nil
}

// GenerateNumber generates a unique loan number
func (r *LoanRepository) GenerateNumber(ctx context.Context) (string, error) {
	query := `SELECT generate_loan_number()`
	var loanNumber string
	err := r.db.QueryRowContext(ctx, query).Scan(&loanNumber)
	if err != nil {
		// Fallback if function doesn't exist
		yearStr := time.Now().Format("2006")
		var seqNum int
		fallbackQuery := `
			SELECT COALESCE(MAX(CAST(SUBSTRING(loan_number FROM 'LN-\d{4}-(\d+)') AS INTEGER)), 0) + 1
			FROM loans WHERE loan_number LIKE 'LN-' || $1 || '-%'
		`
		r.db.QueryRowContext(ctx, fallbackQuery, yearStr).Scan(&seqNum)
		if seqNum == 0 {
			seqNum = 1
		}
		loanNumber = fmt.Sprintf("LN-%s-%06d", yearStr, seqNum)
	}
	return loanNumber, nil
}

// GetOverdueLoans retrieves overdue loans for a branch
func (r *LoanRepository) GetOverdueLoans(ctx context.Context, branchID int64) ([]*domain.Loan, error) {
	query := `
		SELECT id, loan_number, branch_id, customer_id, item_id,
			   loan_amount, interest_rate, interest_amount, principal_remaining, interest_remaining,
			   total_amount, amount_paid, late_fee_rate, late_fee_amount, late_fee_remaining,
			   start_date, due_date, paid_date, confiscated_date,
			   payment_plan_type, loan_term_days, requires_minimum_payment,
			   minimum_payment_amount, next_payment_due_date, grace_period_days,
			   number_of_installments, installment_amount,
			   status, days_overdue, renewed_from_id, renewal_count, notes,
			   created_by, updated_by, created_at, updated_at, deleted_at
		FROM loans
		WHERE (branch_id = $1 OR $1 = 0)
		  AND status IN ('active', 'overdue')
		  AND due_date < NOW()
		  AND deleted_at IS NULL
		ORDER BY due_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue loans: %w", err)
	}
	defer rows.Close()

	loans := []*domain.Loan{}
	for rows.Next() {
		loan, err := r.scanLoanRow(rows)
		if err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	return loans, nil
}

// UpdateStatus updates loan status
func (r *LoanRepository) UpdateStatus(ctx context.Context, id int64, status domain.LoanStatus) error {
	query := `UPDATE loans SET status = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update loan status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("loan not found")
	}

	return nil
}

// BeginTx starts a new transaction
func (r *LoanRepository) BeginTx(ctx context.Context) (repository.Transaction, error) {
	return r.db.BeginTx(ctx)
}

// CreateTx creates a loan within a transaction
func (r *LoanRepository) CreateTx(ctx context.Context, tx repository.Transaction, loan *domain.Loan) error {
	pgTx := tx.(*Tx)

	query := `
		INSERT INTO loans (
			loan_number, branch_id, customer_id, item_id,
			loan_amount, interest_rate, interest_amount, principal_remaining, interest_remaining,
			total_amount, late_fee_rate,
			start_date, due_date,
			payment_plan_type, loan_term_days, requires_minimum_payment,
			minimum_payment_amount, next_payment_due_date, grace_period_days,
			number_of_installments, installment_amount,
			status, notes, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
		RETURNING id, created_at, updated_at
	`

	err := pgTx.QueryRowContext(ctx, query,
		loan.LoanNumber, loan.BranchID, loan.CustomerID, loan.ItemID,
		loan.LoanAmount, loan.InterestRate, loan.InterestAmount,
		loan.PrincipalRemaining, loan.InterestRemaining, loan.TotalAmount, loan.LateFeeRate,
		loan.StartDate, loan.DueDate,
		loan.PaymentPlanType, loan.LoanTermDays, loan.RequiresMinimumPayment,
		NullFloat64(loan.MinimumPaymentAmount), NullTime(loan.NextPaymentDueDate), loan.GracePeriodDays,
		loan.NumberOfInstallments, NullFloat64(loan.InstallmentAmount),
		loan.Status, NullString(loan.Notes), loan.CreatedBy,
	).Scan(&loan.ID, &loan.CreatedAt, &loan.UpdatedAt)

	return err
}

// CreateInstallments creates installments for a loan
func (r *LoanRepository) CreateInstallments(ctx context.Context, installments []*domain.LoanInstallment) error {
	query := `
		INSERT INTO loan_installments (loan_id, installment_number, due_date, principal_amount, interest_amount, total_amount)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	for _, inst := range installments {
		err := r.db.QueryRowContext(ctx, query,
			inst.LoanID, inst.InstallmentNumber, inst.DueDate,
			inst.PrincipalAmount, inst.InterestAmount, inst.TotalAmount,
		).Scan(&inst.ID, &inst.CreatedAt, &inst.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create installment: %w", err)
		}
	}

	return nil
}

// CreateInstallmentsTx creates installments for a loan within a transaction
func (r *LoanRepository) CreateInstallmentsTx(ctx context.Context, tx repository.Transaction, installments []*domain.LoanInstallment) error {
	pgTx := tx.(*Tx)

	query := `
		INSERT INTO loan_installments (loan_id, installment_number, due_date, principal_amount, interest_amount, total_amount)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	for _, inst := range installments {
		err := pgTx.Tx.QueryRowContext(ctx, query,
			inst.LoanID, inst.InstallmentNumber, inst.DueDate,
			inst.PrincipalAmount, inst.InterestAmount, inst.TotalAmount,
		).Scan(&inst.ID, &inst.CreatedAt, &inst.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to create installment: %w", err)
		}
	}

	return nil
}

// GetInstallments retrieves installments for a loan
func (r *LoanRepository) GetInstallments(ctx context.Context, loanID int64) ([]*domain.LoanInstallment, error) {
	query := `
		SELECT id, loan_id, installment_number, due_date, principal_amount, interest_amount,
			   total_amount, amount_paid, is_paid, paid_date, created_at, updated_at
		FROM loan_installments
		WHERE loan_id = $1
		ORDER BY installment_number
	`

	rows, err := r.db.QueryContext(ctx, query, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get installments: %w", err)
	}
	defer rows.Close()

	installments := []*domain.LoanInstallment{}
	for rows.Next() {
		inst := &domain.LoanInstallment{}
		var paidDate sql.NullTime

		err := rows.Scan(
			&inst.ID, &inst.LoanID, &inst.InstallmentNumber, &inst.DueDate,
			&inst.PrincipalAmount, &inst.InterestAmount, &inst.TotalAmount,
			&inst.AmountPaid, &inst.IsPaid, &paidDate, &inst.CreatedAt, &inst.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan installment: %w", err)
		}

		inst.PaidDate = TimePtr(paidDate)
		installments = append(installments, inst)
	}

	return installments, nil
}

// UpdateInstallment updates an installment
func (r *LoanRepository) UpdateInstallment(ctx context.Context, installment *domain.LoanInstallment) error {
	query := `
		UPDATE loan_installments SET
			amount_paid = $2, is_paid = $3, paid_date = $4, updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		installment.ID, installment.AmountPaid, installment.IsPaid, NullTime(installment.PaidDate),
	)
	return err
}

// Helper functions
func (r *LoanRepository) scanLoan(row *sql.Row) (*domain.Loan, error) {
	loan := &domain.Loan{}
	var paidDate, confiscatedDate, nextPaymentDueDate, deletedAt sql.NullTime
	var minimumPaymentAmount, installmentAmount sql.NullFloat64
	var numberOfInstallments, renewedFromID sql.NullInt64
	var notes sql.NullString
	var createdBy, updatedBy sql.NullInt64

	err := row.Scan(
		&loan.ID, &loan.LoanNumber, &loan.BranchID, &loan.CustomerID, &loan.ItemID,
		&loan.LoanAmount, &loan.InterestRate, &loan.InterestAmount,
		&loan.PrincipalRemaining, &loan.InterestRemaining,
		&loan.TotalAmount, &loan.AmountPaid, &loan.LateFeeRate, &loan.LateFeeAmount, &loan.LateFeeRemaining,
		&loan.StartDate, &loan.DueDate, &paidDate, &confiscatedDate,
		&loan.PaymentPlanType, &loan.LoanTermDays, &loan.RequiresMinimumPayment,
		&minimumPaymentAmount, &nextPaymentDueDate, &loan.GracePeriodDays,
		&numberOfInstallments, &installmentAmount,
		&loan.Status, &loan.DaysOverdue, &renewedFromID, &loan.RenewalCount, &notes,
		&createdBy, &updatedBy, &loan.CreatedAt, &loan.UpdatedAt, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("loan not found")
		}
		return nil, fmt.Errorf("failed to get loan: %w", err)
	}

	loan.PaidDate = TimePtr(paidDate)
	loan.ConfiscatedDate = TimePtr(confiscatedDate)
	loan.NextPaymentDueDate = TimePtr(nextPaymentDueDate)
	loan.MinimumPaymentAmount = Float64Ptr(minimumPaymentAmount)
	loan.InstallmentAmount = Float64Ptr(installmentAmount)
	loan.NumberOfInstallments = IntPtr(numberOfInstallments)
	loan.RenewedFromID = Int64Ptr(renewedFromID)
	loan.Notes = StringPtr(notes)
	if createdBy.Valid {
		loan.CreatedBy = createdBy.Int64
	}
	loan.UpdatedBy = Int64Ptr(updatedBy)
	loan.DeletedAt = TimePtr(deletedAt)

	return loan, nil
}

func (r *LoanRepository) scanLoanRow(rows *sql.Rows) (*domain.Loan, error) {
	loan := &domain.Loan{}
	var paidDate, confiscatedDate, nextPaymentDueDate, deletedAt sql.NullTime
	var minimumPaymentAmount, installmentAmount sql.NullFloat64
	var numberOfInstallments, renewedFromID sql.NullInt64
	var notes sql.NullString
	var createdBy, updatedBy sql.NullInt64

	err := rows.Scan(
		&loan.ID, &loan.LoanNumber, &loan.BranchID, &loan.CustomerID, &loan.ItemID,
		&loan.LoanAmount, &loan.InterestRate, &loan.InterestAmount,
		&loan.PrincipalRemaining, &loan.InterestRemaining,
		&loan.TotalAmount, &loan.AmountPaid, &loan.LateFeeRate, &loan.LateFeeAmount, &loan.LateFeeRemaining,
		&loan.StartDate, &loan.DueDate, &paidDate, &confiscatedDate,
		&loan.PaymentPlanType, &loan.LoanTermDays, &loan.RequiresMinimumPayment,
		&minimumPaymentAmount, &nextPaymentDueDate, &loan.GracePeriodDays,
		&numberOfInstallments, &installmentAmount,
		&loan.Status, &loan.DaysOverdue, &renewedFromID, &loan.RenewalCount, &notes,
		&createdBy, &updatedBy, &loan.CreatedAt, &loan.UpdatedAt, &deletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan loan: %w", err)
	}

	loan.PaidDate = TimePtr(paidDate)
	loan.ConfiscatedDate = TimePtr(confiscatedDate)
	loan.NextPaymentDueDate = TimePtr(nextPaymentDueDate)
	loan.MinimumPaymentAmount = Float64Ptr(minimumPaymentAmount)
	loan.InstallmentAmount = Float64Ptr(installmentAmount)
	loan.NumberOfInstallments = IntPtr(numberOfInstallments)
	loan.RenewedFromID = Int64Ptr(renewedFromID)
	loan.Notes = StringPtr(notes)
	if createdBy.Valid {
		loan.CreatedBy = createdBy.Int64
	}
	loan.UpdatedBy = Int64Ptr(updatedBy)
	loan.DeletedAt = TimePtr(deletedAt)

	return loan, nil
}

func (r *LoanRepository) scanLoanRowWithRelations(rows *sql.Rows) (*domain.Loan, error) {
loan := &domain.Loan{}
var paidDate, confiscatedDate, nextPaymentDueDate, deletedAt sql.NullTime
var minimumPaymentAmount, installmentAmount sql.NullFloat64
var numberOfInstallments, renewedFromID sql.NullInt64
var notes sql.NullString
var createdBy, updatedBy sql.NullInt64

// Customer fields
var custID sql.NullInt64
var custFirstName, custLastName, custIdentityNumber sql.NullString

// Item fields
var itemID sql.NullInt64
var itemName, itemSKU sql.NullString

err := rows.Scan(
&loan.ID, &loan.LoanNumber, &loan.BranchID, &loan.CustomerID, &loan.ItemID,
&loan.LoanAmount, &loan.InterestRate, &loan.InterestAmount,
&loan.PrincipalRemaining, &loan.InterestRemaining,
&loan.TotalAmount, &loan.AmountPaid, &loan.LateFeeRate, &loan.LateFeeAmount, &loan.LateFeeRemaining,
&loan.StartDate, &loan.DueDate, &paidDate, &confiscatedDate,
&loan.PaymentPlanType, &loan.LoanTermDays, &loan.RequiresMinimumPayment,
&minimumPaymentAmount, &nextPaymentDueDate, &loan.GracePeriodDays,
&numberOfInstallments, &installmentAmount,
&loan.Status, &loan.DaysOverdue, &renewedFromID, &loan.RenewalCount, &notes,
&createdBy, &updatedBy, &loan.CreatedAt, &loan.UpdatedAt, &deletedAt,
// Customer
&custID, &custFirstName, &custLastName, &custIdentityNumber,
// Item
&itemID, &itemName, &itemSKU,
)

if err != nil {
return nil, fmt.Errorf("failed to scan loan with relations: %w", err)
}

loan.PaidDate = TimePtr(paidDate)
loan.ConfiscatedDate = TimePtr(confiscatedDate)
loan.NextPaymentDueDate = TimePtr(nextPaymentDueDate)
loan.MinimumPaymentAmount = Float64Ptr(minimumPaymentAmount)
loan.InstallmentAmount = Float64Ptr(installmentAmount)
loan.NumberOfInstallments = IntPtr(numberOfInstallments)
loan.RenewedFromID = Int64Ptr(renewedFromID)
loan.Notes = StringPtr(notes)
if createdBy.Valid {
loan.CreatedBy = createdBy.Int64
}
loan.UpdatedBy = Int64Ptr(updatedBy)
loan.DeletedAt = TimePtr(deletedAt)

// Populate Customer relation
if custID.Valid {
loan.Customer = &domain.Customer{
ID:             custID.Int64,
FirstName:      custFirstName.String,
LastName:       custLastName.String,
IdentityNumber: custIdentityNumber.String,
}
}

// Populate Item relation
if itemID.Valid {
loan.Item = &domain.Item{
ID:   itemID.Int64,
Name: itemName.String,
SKU:  itemSKU.String,
}
}

return loan, nil
}
