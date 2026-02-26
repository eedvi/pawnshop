package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// CashRegisterRepository implements repository.CashRegisterRepository
type CashRegisterRepository struct {
	db *DB
}

// NewCashRegisterRepository creates a new CashRegisterRepository
func NewCashRegisterRepository(db *DB) *CashRegisterRepository {
	return &CashRegisterRepository{db: db}
}

// GetByID retrieves a cash register by ID
func (r *CashRegisterRepository) GetByID(ctx context.Context, id int64) (*domain.CashRegister, error) {
	query := `
		SELECT id, branch_id, name, description, is_active, created_at, updated_at
		FROM cash_registers
		WHERE id = $1
	`

	register := &domain.CashRegister{}
	var description sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&register.ID, &register.BranchID, &register.Name, &description,
		&register.IsActive, &register.CreatedAt, &register.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("cash register not found")
		}
		return nil, fmt.Errorf("failed to get cash register: %w", err)
	}

	register.Description = StringPtrVal(description)
	return register, nil
}

// List retrieves all cash registers for a branch
func (r *CashRegisterRepository) List(ctx context.Context, branchID int64) ([]*domain.CashRegister, error) {
	query := `
		SELECT id, branch_id, name, description, is_active, created_at, updated_at
		FROM cash_registers
		WHERE branch_id = $1
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to list cash registers: %w", err)
	}
	defer rows.Close()

	registers := []*domain.CashRegister{}
	for rows.Next() {
		register := &domain.CashRegister{}
		var description sql.NullString

		err := rows.Scan(
			&register.ID, &register.BranchID, &register.Name, &description,
			&register.IsActive, &register.CreatedAt, &register.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cash register: %w", err)
		}

		register.Description = StringPtrVal(description)
		registers = append(registers, register)
	}

	return registers, nil
}

// Create creates a new cash register
func (r *CashRegisterRepository) Create(ctx context.Context, register *domain.CashRegister) error {
	query := `
		INSERT INTO cash_registers (branch_id, name, description, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		register.BranchID, register.Name, NullStringPtr(register.Description), register.IsActive,
	).Scan(&register.ID, &register.CreatedAt, &register.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create cash register: %w", err)
	}

	return nil
}

// Update updates an existing cash register
func (r *CashRegisterRepository) Update(ctx context.Context, register *domain.CashRegister) error {
	query := `
		UPDATE cash_registers SET
			name = $2, description = $3, is_active = $4, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		register.ID, register.Name, NullStringPtr(register.Description), register.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to update cash register: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("cash register not found")
	}

	return nil
}

// CashSessionRepository implements repository.CashSessionRepository
type CashSessionRepository struct {
	db *DB
}

// NewCashSessionRepository creates a new CashSessionRepository
func NewCashSessionRepository(db *DB) *CashSessionRepository {
	return &CashSessionRepository{db: db}
}

// GetByID retrieves a cash session by ID
func (r *CashSessionRepository) GetByID(ctx context.Context, id int64) (*domain.CashSession, error) {
	query := `
		SELECT id, branch_id, cash_register_id, user_id, status,
			   opening_amount, closing_amount, expected_amount, difference,
			   opened_at, closed_at, closed_by, opening_notes, closing_notes,
			   created_at, updated_at
		FROM cash_sessions
		WHERE id = $1
	`

	return r.scanSession(r.db.QueryRowContext(ctx, query, id))
}

// GetOpenSession retrieves an open session for a user
func (r *CashSessionRepository) GetOpenSession(ctx context.Context, userID int64) (*domain.CashSession, error) {
	query := `
		SELECT id, branch_id, cash_register_id, user_id, status,
			   opening_amount, closing_amount, expected_amount, difference,
			   opened_at, closed_at, closed_by, opening_notes, closing_notes,
			   created_at, updated_at
		FROM cash_sessions
		WHERE user_id = $1 AND status = 'open'
		ORDER BY opened_at DESC
		LIMIT 1
	`

	return r.scanSession(r.db.QueryRowContext(ctx, query, userID))
}

// GetOpenSessionByRegister retrieves an open session for a register
func (r *CashSessionRepository) GetOpenSessionByRegister(ctx context.Context, registerID int64) (*domain.CashSession, error) {
	query := `
		SELECT id, branch_id, cash_register_id, user_id, status,
			   opening_amount, closing_amount, expected_amount, difference,
			   opened_at, closed_at, closed_by, opening_notes, closing_notes,
			   created_at, updated_at
		FROM cash_sessions
		WHERE cash_register_id = $1 AND status = 'open'
		ORDER BY opened_at DESC
		LIMIT 1
	`

	return r.scanSession(r.db.QueryRowContext(ctx, query, registerID))
}

// List retrieves cash sessions with pagination and filters
func (r *CashSessionRepository) List(ctx context.Context, params repository.CashSessionListParams) (*repository.PaginatedResult[domain.CashSession], error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 20
	}

	whereClause := `WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if params.BranchID > 0 {
		argCount++
		whereClause += fmt.Sprintf(" AND cs.branch_id = $%d", argCount)
		args = append(args, params.BranchID)
	}

	if params.UserID != nil {
		argCount++
		whereClause += fmt.Sprintf(" AND cs.user_id = $%d", argCount)
		args = append(args, *params.UserID)
	}

	if params.RegisterID != nil {
		argCount++
		whereClause += fmt.Sprintf(" AND cs.cash_register_id = $%d", argCount)
		args = append(args, *params.RegisterID)
	}

	if params.Status != nil {
		argCount++
		whereClause += fmt.Sprintf(" AND cs.status = $%d", argCount)
		args = append(args, *params.Status)
	}

	if params.DateFrom != nil {
		argCount++
		whereClause += fmt.Sprintf(" AND cs.opened_at >= $%d", argCount)
		args = append(args, *params.DateFrom)
	}

	if params.DateTo != nil {
		argCount++
		whereClause += fmt.Sprintf(" AND cs.opened_at <= $%d", argCount)
		args = append(args, *params.DateTo)
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) FROM cash_sessions cs " + whereClause
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count cash sessions: %w", err)
	}

	// Get data
	orderBy := "cs.opened_at"
	if params.OrderBy != "" {
		orderBy = "cs." + params.OrderBy
	}
	order := "DESC"
	if params.Order == "asc" {
		order = "ASC"
	}

	offset := (params.Page - 1) * params.PerPage

	dataQuery := fmt.Sprintf(`
		SELECT cs.id, cs.branch_id, cs.cash_register_id, cs.user_id, cs.status,
			   cs.opening_amount, cs.closing_amount, cs.expected_amount, cs.difference,
			   cs.opened_at, cs.closed_at, cs.closed_by, cs.opening_notes, cs.closing_notes,
			   cs.created_at, cs.updated_at,
			   cr.id, cr.branch_id, cr.name, cr.code, cr.description, cr.is_active, cr.created_at, cr.updated_at,
			   u.id, u.branch_id, u.first_name, u.last_name, u.email, u.phone
		FROM cash_sessions cs
		LEFT JOIN cash_registers cr ON cs.cash_register_id = cr.id
		LEFT JOIN users u ON cs.user_id = u.id AND u.deleted_at IS NULL
		%s
		ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, order, argCount+1, argCount+2,
	)
	args = append(args, params.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list cash sessions: %w", err)
	}
	defer rows.Close()

	sessions := []domain.CashSession{}
	for rows.Next() {
		session, err := r.scanSessionRowWithRelations(rows)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, *session)
	}

	totalPages := total / params.PerPage
	if total%params.PerPage > 0 {
		totalPages++
	}

	return &repository.PaginatedResult[domain.CashSession]{
		Data:       sessions,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

// Create creates a new cash session
func (r *CashSessionRepository) Create(ctx context.Context, session *domain.CashSession) error {
	query := `
		INSERT INTO cash_sessions (branch_id, cash_register_id, user_id, status, opening_amount, opened_at, opening_notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		session.BranchID, session.CashRegisterID, session.UserID, session.Status,
		session.OpeningAmount, session.OpenedAt, NullStringPtr(session.OpeningNotes),
	).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create cash session: %w", err)
	}

	return nil
}

// Update updates an existing cash session
func (r *CashSessionRepository) Update(ctx context.Context, session *domain.CashSession) error {
	query := `
		UPDATE cash_sessions SET
			status = $2, closing_amount = $3, expected_amount = $4, difference = $5,
			closed_at = $6, closed_by = $7, closing_notes = $8, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		session.ID, session.Status, NullFloat64(session.ClosingAmount),
		NullFloat64(session.ExpectedAmount), NullFloat64(session.Difference),
		session.ClosedAt, NullInt64(session.ClosedBy), NullStringPtr(session.ClosingNotes),
	)
	if err != nil {
		return fmt.Errorf("failed to update cash session: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("cash session not found")
	}

	return nil
}

// Close closes a cash session
func (r *CashSessionRepository) Close(ctx context.Context, id int64, data repository.CashSessionCloseData) error {
	query := `
		UPDATE cash_sessions SET
			status = 'closed', closing_amount = $2, expected_amount = $3, difference = $4,
			closed_at = NOW(), closed_by = $5, closing_notes = $6, updated_at = NOW()
		WHERE id = $1 AND status = 'open'
	`

	result, err := r.db.ExecContext(ctx, query,
		id, data.ClosingAmount, data.ExpectedAmount, data.Difference, data.ClosedBy, data.ClosingNotes,
	)
	if err != nil {
		return fmt.Errorf("failed to close cash session: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("cash session not found or already closed")
	}

	return nil
}

// Helper functions for CashSessionRepository
func (r *CashSessionRepository) scanSession(row *sql.Row) (*domain.CashSession, error) {
	session := &domain.CashSession{}
	var closingAmount, expectedAmount, difference sql.NullFloat64
	var closedAt sql.NullTime
	var closedBy sql.NullInt64
	var openingNotes, closingNotes sql.NullString

	err := row.Scan(
		&session.ID, &session.BranchID, &session.CashRegisterID, &session.UserID, &session.Status,
		&session.OpeningAmount, &closingAmount, &expectedAmount, &difference,
		&session.OpenedAt, &closedAt, &closedBy, &openingNotes, &closingNotes,
		&session.CreatedAt, &session.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("cash session not found")
		}
		return nil, fmt.Errorf("failed to get cash session: %w", err)
	}

	session.ClosingAmount = Float64Ptr(closingAmount)
	session.ExpectedAmount = Float64Ptr(expectedAmount)
	session.Difference = Float64Ptr(difference)
	if closedAt.Valid {
		session.ClosedAt = &closedAt.Time
	}
	session.ClosedBy = Int64Ptr(closedBy)
	session.OpeningNotes = StringPtrVal(openingNotes)
	session.ClosingNotes = StringPtrVal(closingNotes)

	return session, nil
}

func (r *CashSessionRepository) scanSessionRow(rows *sql.Rows) (*domain.CashSession, error) {
	session := &domain.CashSession{}
	var closingAmount, expectedAmount, difference sql.NullFloat64
	var closedAt sql.NullTime
	var closedBy sql.NullInt64
	var openingNotes, closingNotes sql.NullString

	err := rows.Scan(
		&session.ID, &session.BranchID, &session.CashRegisterID, &session.UserID, &session.Status,
		&session.OpeningAmount, &closingAmount, &expectedAmount, &difference,
		&session.OpenedAt, &closedAt, &closedBy, &openingNotes, &closingNotes,
		&session.CreatedAt, &session.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan cash session: %w", err)
	}

	session.ClosingAmount = Float64Ptr(closingAmount)
	session.ExpectedAmount = Float64Ptr(expectedAmount)
	session.Difference = Float64Ptr(difference)
	if closedAt.Valid {
		session.ClosedAt = &closedAt.Time
	}
	session.ClosedBy = Int64Ptr(closedBy)
	session.OpeningNotes = StringPtrVal(openingNotes)
	session.ClosingNotes = StringPtrVal(closingNotes)

	return session, nil
}

func (r *CashSessionRepository) scanSessionRowWithRelations(rows *sql.Rows) (*domain.CashSession, error) {
	session := &domain.CashSession{}
	var closingAmount, expectedAmount, difference sql.NullFloat64
	var closedAt sql.NullTime
	var closedBy sql.NullInt64
	var openingNotes, closingNotes sql.NullString

	// Register fields
	var registerID, registerBranchID sql.NullInt64
	var registerName, registerCode sql.NullString
	var registerDescription sql.NullString
	var registerIsActive sql.NullBool
	var registerCreatedAt, registerUpdatedAt sql.NullTime

	// User fields
	var userID, userBranchID sql.NullInt64
	var userFirstName, userLastName, userEmail, userPhone sql.NullString

	err := rows.Scan(
		&session.ID, &session.BranchID, &session.CashRegisterID, &session.UserID, &session.Status,
		&session.OpeningAmount, &closingAmount, &expectedAmount, &difference,
		&session.OpenedAt, &closedAt, &closedBy, &openingNotes, &closingNotes,
		&session.CreatedAt, &session.UpdatedAt,
		// Register
		&registerID, &registerBranchID, &registerName, &registerCode, &registerDescription, &registerIsActive, &registerCreatedAt, &registerUpdatedAt,
		// User
		&userID, &userBranchID, &userFirstName, &userLastName, &userEmail, &userPhone,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan cash session with relations: %w", err)
	}

	session.ClosingAmount = Float64Ptr(closingAmount)
	session.ExpectedAmount = Float64Ptr(expectedAmount)
	session.Difference = Float64Ptr(difference)
	if closedAt.Valid {
		session.ClosedAt = &closedAt.Time
	}
	session.ClosedBy = Int64Ptr(closedBy)
	session.OpeningNotes = StringPtrVal(openingNotes)
	session.ClosingNotes = StringPtrVal(closingNotes)

	// Populate register
	if registerID.Valid {
		session.CashRegister = &domain.CashRegister{
			ID:          registerID.Int64,
			BranchID:    registerBranchID.Int64,
			Name:        registerName.String,
			Code:        registerCode.String,
			Description: StringPtrVal(registerDescription),
			IsActive:    registerIsActive.Bool,
			CreatedAt:   registerCreatedAt.Time,
			UpdatedAt:   registerUpdatedAt.Time,
		}
	}

	// Populate user
	if userID.Valid {
		session.User = &domain.User{
			ID:        userID.Int64,
			BranchID:  Int64Ptr(userBranchID),
			FirstName: userFirstName.String,
			LastName:  userLastName.String,
			Email:     userEmail.String,
			Phone:     userPhone.String,
		}
	}

	return session, nil
}

// CashMovementRepository implements repository.CashMovementRepository
type CashMovementRepository struct {
	db *DB
}

// NewCashMovementRepository creates a new CashMovementRepository
func NewCashMovementRepository(db *DB) *CashMovementRepository {
	return &CashMovementRepository{db: db}
}

// GetByID retrieves a cash movement by ID
func (r *CashMovementRepository) GetByID(ctx context.Context, id int64) (*domain.CashMovement, error) {
	query := `
		SELECT id, branch_id, session_id, movement_type, amount,
			   payment_method, reference_type, reference_id, description,
			   balance_after, created_by, created_at
		FROM cash_movements
		WHERE id = $1
	`

	movement := &domain.CashMovement{}
	var referenceType sql.NullString
	var referenceID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&movement.ID, &movement.BranchID, &movement.SessionID, &movement.MovementType, &movement.Amount,
		&movement.PaymentMethod, &referenceType, &referenceID, &movement.Description,
		&movement.BalanceAfter, &movement.CreatedBy, &movement.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("cash movement not found")
		}
		return nil, fmt.Errorf("failed to get cash movement: %w", err)
	}

	movement.ReferenceType = StringPtrVal(referenceType)
	movement.ReferenceID = Int64Ptr(referenceID)

	return movement, nil
}

// List retrieves cash movements with pagination and filters
func (r *CashMovementRepository) List(ctx context.Context, params repository.CashMovementListParams) (*repository.PaginatedResult[domain.CashMovement], error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 20
	}

	baseQuery := `FROM cash_movements WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if params.BranchID > 0 {
		argCount++
		baseQuery += fmt.Sprintf(" AND branch_id = $%d", argCount)
		args = append(args, params.BranchID)
	}

	if params.SessionID != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND session_id = $%d", argCount)
		args = append(args, *params.SessionID)
	}

	if params.MovementType != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND movement_type = $%d", argCount)
		args = append(args, *params.MovementType)
	}

	if params.PaymentMethod != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND payment_method = $%d", argCount)
		args = append(args, *params.PaymentMethod)
	}

	if params.DateFrom != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, *params.DateFrom)
	}

	if params.DateTo != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, *params.DateTo)
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count cash movements: %w", err)
	}

	// Get data
	orderBy := "created_at"
	if params.OrderBy != "" {
		orderBy = params.OrderBy
	}
	order := "DESC"
	if params.Order == "asc" {
		order = "ASC"
	}

	offset := (params.Page - 1) * params.PerPage
	dataQuery := fmt.Sprintf(`
		SELECT id, branch_id, session_id, movement_type, amount,
			   payment_method, reference_type, reference_id, description,
			   balance_after, created_by, created_at
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		baseQuery, orderBy, order, argCount+1, argCount+2,
	)
	args = append(args, params.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list cash movements: %w", err)
	}
	defer rows.Close()

	movements := []domain.CashMovement{}
	for rows.Next() {
		movement := &domain.CashMovement{}
		var referenceType sql.NullString
		var referenceID sql.NullInt64

		err := rows.Scan(
			&movement.ID, &movement.BranchID, &movement.SessionID, &movement.MovementType, &movement.Amount,
			&movement.PaymentMethod, &referenceType, &referenceID, &movement.Description,
			&movement.BalanceAfter, &movement.CreatedBy, &movement.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cash movement: %w", err)
		}

		movement.ReferenceType = StringPtrVal(referenceType)
		movement.ReferenceID = Int64Ptr(referenceID)
		movements = append(movements, *movement)
	}

	totalPages := total / params.PerPage
	if total%params.PerPage > 0 {
		totalPages++
	}

	return &repository.PaginatedResult[domain.CashMovement]{
		Data:       movements,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

// ListBySession retrieves all cash movements for a session
func (r *CashMovementRepository) ListBySession(ctx context.Context, sessionID int64) ([]*domain.CashMovement, error) {
	query := `
		SELECT id, branch_id, session_id, movement_type, amount,
			   payment_method, reference_type, reference_id, description,
			   balance_after, created_by, created_at
		FROM cash_movements
		WHERE session_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list cash movements: %w", err)
	}
	defer rows.Close()

	movements := []*domain.CashMovement{}
	for rows.Next() {
		movement := &domain.CashMovement{}
		var referenceType sql.NullString
		var referenceID sql.NullInt64

		err := rows.Scan(
			&movement.ID, &movement.BranchID, &movement.SessionID, &movement.MovementType, &movement.Amount,
			&movement.PaymentMethod, &referenceType, &referenceID, &movement.Description,
			&movement.BalanceAfter, &movement.CreatedBy, &movement.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cash movement: %w", err)
		}

		movement.ReferenceType = StringPtrVal(referenceType)
		movement.ReferenceID = Int64Ptr(referenceID)
		movements = append(movements, movement)
	}

	return movements, nil
}

// Create creates a new cash movement
func (r *CashMovementRepository) Create(ctx context.Context, movement *domain.CashMovement) error {
	query := `
		INSERT INTO cash_movements (
			branch_id, session_id, movement_type, amount,
			payment_method, reference_type, reference_id, description,
			balance_after, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query,
		movement.BranchID, movement.SessionID, movement.MovementType, movement.Amount,
		movement.PaymentMethod, NullStringPtr(movement.ReferenceType), NullInt64(movement.ReferenceID),
		movement.Description, movement.BalanceAfter, movement.CreatedBy,
	).Scan(&movement.ID, &movement.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create cash movement: %w", err)
	}

	return nil
}

// GetSessionBalance retrieves the current balance for a session
func (r *CashMovementRepository) GetSessionBalance(ctx context.Context, sessionID int64) (float64, error) {
	query := `
		SELECT COALESCE(
			(SELECT balance_after FROM cash_movements WHERE session_id = $1 ORDER BY created_at DESC LIMIT 1),
			(SELECT opening_amount FROM cash_sessions WHERE id = $1)
		)
	`

	var balance float64
	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to get session balance: %w", err)
	}

	return balance, nil
}

// GetSessionSummary retrieves a summary of movements for a session
func (r *CashMovementRepository) GetSessionSummary(ctx context.Context, sessionID int64) (*CashSessionSummary, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN movement_type = 'income' THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN movement_type = 'expense' THEN amount ELSE 0 END), 0) as total_expense,
			COALESCE(SUM(CASE WHEN movement_type = 'income' AND payment_method = 'cash' THEN amount ELSE 0 END), 0) as cash_income,
			COALESCE(SUM(CASE WHEN movement_type = 'expense' AND payment_method = 'cash' THEN amount ELSE 0 END), 0) as cash_expense,
			COUNT(*) as total_movements
		FROM cash_movements
		WHERE session_id = $1
	`

	summary := &CashSessionSummary{}
	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&summary.TotalIncome, &summary.TotalExpense,
		&summary.CashIncome, &summary.CashExpense, &summary.TotalMovements,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get session summary: %w", err)
	}

	return summary, nil
}

// CashSessionSummary contains summary data for a cash session
type CashSessionSummary struct {
	TotalIncome    float64 `json:"total_income"`
	TotalExpense   float64 `json:"total_expense"`
	CashIncome     float64 `json:"cash_income"`
	CashExpense    float64 `json:"cash_expense"`
	TotalMovements int     `json:"total_movements"`
}

