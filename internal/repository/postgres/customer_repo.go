package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// CustomerRepository implements repository.CustomerRepository
type CustomerRepository struct {
	db *DB
}

// NewCustomerRepository creates a new CustomerRepository
func NewCustomerRepository(db *DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// GetByID retrieves a customer by ID
func (r *CustomerRepository) GetByID(ctx context.Context, id int64) (*domain.Customer, error) {
	query := `
		SELECT id, branch_id, first_name, last_name, identity_type, identity_number,
			   birth_date, gender, phone, phone_secondary, email, address, city, state, postal_code,
			   emergency_contact_name, emergency_contact_phone, emergency_contact_relation,
			   occupation, workplace, monthly_income,
			   credit_limit, credit_score, total_loans, total_paid, total_defaulted,
			   is_active, is_blocked, blocked_reason, notes, photo_url,
			   created_by, created_at, updated_at, deleted_at
		FROM customers
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanCustomer(r.db.QueryRowContext(ctx, query, id))
}

// GetByIdentity retrieves a customer by identity
func (r *CustomerRepository) GetByIdentity(ctx context.Context, branchID int64, identityType, identityNumber string) (*domain.Customer, error) {
	query := `
		SELECT id, branch_id, first_name, last_name, identity_type, identity_number,
			   birth_date, gender, phone, phone_secondary, email, address, city, state, postal_code,
			   emergency_contact_name, emergency_contact_phone, emergency_contact_relation,
			   occupation, workplace, monthly_income,
			   credit_limit, credit_score, total_loans, total_paid, total_defaulted,
			   is_active, is_blocked, blocked_reason, notes, photo_url,
			   created_by, created_at, updated_at, deleted_at
		FROM customers
		WHERE branch_id = $1 AND identity_type = $2 AND identity_number = $3 AND deleted_at IS NULL
	`

	return r.scanCustomer(r.db.QueryRowContext(ctx, query, branchID, identityType, identityNumber))
}

// List retrieves customers with pagination and filters
func (r *CustomerRepository) List(ctx context.Context, params repository.CustomerListParams) (*repository.PaginatedResult[domain.Customer], error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 20
	}

	baseQuery := `FROM customers WHERE deleted_at IS NULL`
	args := []interface{}{}
	argCount := 0

	if params.BranchID > 0 {
		argCount++
		baseQuery += fmt.Sprintf(" AND branch_id = $%d", argCount)
		args = append(args, params.BranchID)
	}

	if params.IsActive != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *params.IsActive)
	}

	if params.IsBlocked != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND is_blocked = $%d", argCount)
		args = append(args, *params.IsBlocked)
	}

	if params.Search != "" {
		argCount++
		baseQuery += fmt.Sprintf(" AND (first_name ILIKE $%d OR last_name ILIKE $%d OR identity_number ILIKE $%d OR phone ILIKE $%d)", argCount, argCount, argCount, argCount)
		args = append(args, "%"+params.Search+"%")
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count customers: %w", err)
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
		SELECT id, branch_id, first_name, last_name, identity_type, identity_number,
			   birth_date, gender, phone, phone_secondary, email, address, city, state, postal_code,
			   emergency_contact_name, emergency_contact_phone, emergency_contact_relation,
			   occupation, workplace, monthly_income,
			   credit_limit, credit_score, total_loans, total_paid, total_defaulted,
			   is_active, is_blocked, blocked_reason, notes, photo_url,
			   created_by, created_at, updated_at, deleted_at
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		baseQuery, orderBy, order, argCount+1, argCount+2,
	)
	args = append(args, params.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	customers := []domain.Customer{}
	for rows.Next() {
		customer, err := r.scanCustomerRow(rows)
		if err != nil {
			return nil, err
		}
		customers = append(customers, *customer)
	}

	totalPages := total / params.PerPage
	if total%params.PerPage > 0 {
		totalPages++
	}

	return &repository.PaginatedResult[domain.Customer]{
		Data:       customers,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

// Create creates a new customer
func (r *CustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO customers (
			branch_id, first_name, last_name, identity_type, identity_number,
			birth_date, gender, phone, phone_secondary, email, address, city, state, postal_code,
			emergency_contact_name, emergency_contact_phone, emergency_contact_relation,
			occupation, workplace, monthly_income,
			credit_limit, credit_score, is_active, notes, photo_url, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		customer.BranchID, customer.FirstName, customer.LastName,
		customer.IdentityType, customer.IdentityNumber,
		NullTime(customer.BirthDate), NullString(customer.Gender),
		customer.Phone, NullString(customer.PhoneSecondary), NullString(customer.Email),
		NullString(customer.Address), NullString(customer.City), NullString(customer.State), NullString(customer.PostalCode),
		NullString(customer.EmergencyContactName), NullString(customer.EmergencyContactPhone), NullString(customer.EmergencyContactRelation),
		NullString(customer.Occupation), NullString(customer.Workplace), NullFloat64(&customer.MonthlyIncome),
		customer.CreditLimit, customer.CreditScore, customer.IsActive,
		NullString(customer.Notes), NullString(customer.PhotoURL), customer.CreatedBy,
	).Scan(&customer.ID, &customer.CreatedAt, &customer.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	return nil
}

// Update updates an existing customer
func (r *CustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	query := `
		UPDATE customers SET
			first_name = $2, last_name = $3, identity_type = $4, identity_number = $5,
			birth_date = $6, gender = $7, phone = $8, phone_secondary = $9, email = $10,
			address = $11, city = $12, state = $13, postal_code = $14,
			emergency_contact_name = $15, emergency_contact_phone = $16, emergency_contact_relation = $17,
			occupation = $18, workplace = $19, monthly_income = $20,
			credit_limit = $21, is_active = $22, is_blocked = $23, blocked_reason = $24,
			notes = $25, photo_url = $26, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		customer.ID, customer.FirstName, customer.LastName,
		customer.IdentityType, customer.IdentityNumber,
		NullTime(customer.BirthDate), NullString(customer.Gender),
		customer.Phone, NullString(customer.PhoneSecondary), NullString(customer.Email),
		NullString(customer.Address), NullString(customer.City), NullString(customer.State), NullString(customer.PostalCode),
		NullString(customer.EmergencyContactName), NullString(customer.EmergencyContactPhone), NullString(customer.EmergencyContactRelation),
		NullString(customer.Occupation), NullString(customer.Workplace), NullFloat64(&customer.MonthlyIncome),
		customer.CreditLimit, customer.IsActive, customer.IsBlocked, NullString(customer.BlockedReason),
		NullString(customer.Notes), NullString(customer.PhotoURL),
	)
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

// Delete soft deletes a customer
func (r *CustomerRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE customers SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

// UpdateCreditInfo updates customer credit information
func (r *CustomerRepository) UpdateCreditInfo(ctx context.Context, id int64, info repository.CustomerCreditUpdate) error {
	query := `UPDATE customers SET `
	args := []interface{}{id}
	updates := []string{}
	argCount := 1

	if info.CreditLimit != nil {
		argCount++
		updates = append(updates, fmt.Sprintf("credit_limit = $%d", argCount))
		args = append(args, *info.CreditLimit)
	}
	if info.CreditScore != nil {
		argCount++
		updates = append(updates, fmt.Sprintf("credit_score = $%d", argCount))
		args = append(args, *info.CreditScore)
	}
	if info.TotalLoans != nil {
		argCount++
		updates = append(updates, fmt.Sprintf("total_loans = $%d", argCount))
		args = append(args, *info.TotalLoans)
	}
	if info.TotalPaid != nil {
		argCount++
		updates = append(updates, fmt.Sprintf("total_paid = $%d", argCount))
		args = append(args, *info.TotalPaid)
	}
	if info.TotalDefaulted != nil {
		argCount++
		updates = append(updates, fmt.Sprintf("total_defaulted = $%d", argCount))
		args = append(args, *info.TotalDefaulted)
	}

	if len(updates) == 0 {
		return nil
	}

	query += updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += ", updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL"

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// Helper functions
func (r *CustomerRepository) scanCustomer(row *sql.Row) (*domain.Customer, error) {
	c := &domain.Customer{}
	var birthDate, deletedAt sql.NullTime
	var gender, phoneSecondary, email, address, city, state, postalCode sql.NullString
	var emergencyName, emergencyPhone, emergencyRelation sql.NullString
	var occupation, workplace, blockedReason, notes, photoURL sql.NullString
	var monthlyIncome sql.NullFloat64
	var createdBy sql.NullInt64

	err := row.Scan(
		&c.ID, &c.BranchID, &c.FirstName, &c.LastName, &c.IdentityType, &c.IdentityNumber,
		&birthDate, &gender, &c.Phone, &phoneSecondary, &email, &address, &city, &state, &postalCode,
		&emergencyName, &emergencyPhone, &emergencyRelation,
		&occupation, &workplace, &monthlyIncome,
		&c.CreditLimit, &c.CreditScore, &c.TotalLoans, &c.TotalPaid, &c.TotalDefaulted,
		&c.IsActive, &c.IsBlocked, &blockedReason, &notes, &photoURL,
		&createdBy, &c.CreatedAt, &c.UpdatedAt, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	c.BirthDate = TimePtr(birthDate)
	c.Gender = StringPtr(gender)
	c.PhoneSecondary = StringPtr(phoneSecondary)
	c.Email = StringPtr(email)
	c.Address = StringPtr(address)
	c.City = StringPtr(city)
	c.State = StringPtr(state)
	c.PostalCode = StringPtr(postalCode)
	c.EmergencyContactName = StringPtr(emergencyName)
	c.EmergencyContactPhone = StringPtr(emergencyPhone)
	c.EmergencyContactRelation = StringPtr(emergencyRelation)
	c.Occupation = StringPtr(occupation)
	c.Workplace = StringPtr(workplace)
	if monthlyIncome.Valid {
		c.MonthlyIncome = monthlyIncome.Float64
	}
	c.BlockedReason = StringPtr(blockedReason)
	c.Notes = StringPtr(notes)
	c.PhotoURL = StringPtr(photoURL)
	if createdBy.Valid {
		c.CreatedBy = createdBy.Int64
	}
	c.DeletedAt = TimePtr(deletedAt)

	return c, nil
}

func (r *CustomerRepository) scanCustomerRow(rows *sql.Rows) (*domain.Customer, error) {
	c := &domain.Customer{}
	var birthDate, deletedAt sql.NullTime
	var gender, phoneSecondary, email, address, city, state, postalCode sql.NullString
	var emergencyName, emergencyPhone, emergencyRelation sql.NullString
	var occupation, workplace, blockedReason, notes, photoURL sql.NullString
	var monthlyIncome sql.NullFloat64
	var createdBy sql.NullInt64

	err := rows.Scan(
		&c.ID, &c.BranchID, &c.FirstName, &c.LastName, &c.IdentityType, &c.IdentityNumber,
		&birthDate, &gender, &c.Phone, &phoneSecondary, &email, &address, &city, &state, &postalCode,
		&emergencyName, &emergencyPhone, &emergencyRelation,
		&occupation, &workplace, &monthlyIncome,
		&c.CreditLimit, &c.CreditScore, &c.TotalLoans, &c.TotalPaid, &c.TotalDefaulted,
		&c.IsActive, &c.IsBlocked, &blockedReason, &notes, &photoURL,
		&createdBy, &c.CreatedAt, &c.UpdatedAt, &deletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan customer: %w", err)
	}

	c.BirthDate = TimePtr(birthDate)
	c.Gender = StringPtr(gender)
	c.PhoneSecondary = StringPtr(phoneSecondary)
	c.Email = StringPtr(email)
	c.Address = StringPtr(address)
	c.City = StringPtr(city)
	c.State = StringPtr(state)
	c.PostalCode = StringPtr(postalCode)
	c.EmergencyContactName = StringPtr(emergencyName)
	c.EmergencyContactPhone = StringPtr(emergencyPhone)
	c.EmergencyContactRelation = StringPtr(emergencyRelation)
	c.Occupation = StringPtr(occupation)
	c.Workplace = StringPtr(workplace)
	if monthlyIncome.Valid {
		c.MonthlyIncome = monthlyIncome.Float64
	}
	c.BlockedReason = StringPtr(blockedReason)
	c.Notes = StringPtr(notes)
	c.PhotoURL = StringPtr(photoURL)
	if createdBy.Valid {
		c.CreatedBy = createdBy.Int64
	}
	c.DeletedAt = TimePtr(deletedAt)

	return c, nil
}
