package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// SaleRepository implements repository.SaleRepository
type SaleRepository struct {
	db *DB
}

// NewSaleRepository creates a new SaleRepository
func NewSaleRepository(db *DB) *SaleRepository {
	return &SaleRepository{db: db}
}

// GetByID retrieves a sale by ID
func (r *SaleRepository) GetByID(ctx context.Context, id int64) (*domain.Sale, error) {
	query := `
		SELECT id, branch_id, item_id, customer_id, sale_number, sale_type,
			   sale_price, discount_amount, discount_reason, final_price,
			   payment_method, reference_number, status, sale_date,
			   refund_amount, refund_reason, refunded_at, refunded_by,
			   notes, cash_session_id, created_by, updated_by,
			   created_at, updated_at, deleted_at
		FROM sales
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanSale(r.db.QueryRowContext(ctx, query, id))
}

// GetByNumber retrieves a sale by sale number
func (r *SaleRepository) GetByNumber(ctx context.Context, saleNumber string) (*domain.Sale, error) {
	query := `
		SELECT id, branch_id, item_id, customer_id, sale_number, sale_type,
			   sale_price, discount_amount, discount_reason, final_price,
			   payment_method, reference_number, status, sale_date,
			   refund_amount, refund_reason, refunded_at, refunded_by,
			   notes, cash_session_id, created_by, updated_by,
			   created_at, updated_at, deleted_at
		FROM sales
		WHERE sale_number = $1 AND deleted_at IS NULL
	`

	return r.scanSale(r.db.QueryRowContext(ctx, query, saleNumber))
}

// List retrieves sales with pagination and filters
func (r *SaleRepository) List(ctx context.Context, params repository.SaleListParams) (*repository.PaginatedResult[domain.Sale], error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 20
	}

	// Base query with JOINs for related entities
	baseQuery := `
		FROM sales s
		LEFT JOIN items i ON s.item_id = i.id
		LEFT JOIN customers c ON s.customer_id = c.id
		WHERE s.deleted_at IS NULL`
	args := []interface{}{}
	argCount := 0

	if params.BranchID > 0 {
		argCount++
		baseQuery += fmt.Sprintf(" AND s.branch_id = $%d", argCount)
		args = append(args, params.BranchID)
	}

	if params.CustomerID != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND s.customer_id = $%d", argCount)
		args = append(args, *params.CustomerID)
	}

	if params.ItemID != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND s.item_id = $%d", argCount)
		args = append(args, *params.ItemID)
	}

	if params.Status != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND s.status = $%d", argCount)
		args = append(args, *params.Status)
	}

	if params.DateFrom != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND s.sale_date >= $%d", argCount)
		args = append(args, *params.DateFrom)
	}

	if params.DateTo != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND s.sale_date <= $%d", argCount)
		args = append(args, *params.DateTo)
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count sales: %w", err)
	}

	// Get data
	orderBy := "s.created_at"
	if params.OrderBy != "" {
		orderBy = "s." + params.OrderBy
	}
	order := "DESC"
	if params.Order == "asc" {
		order = "ASC"
	}

	offset := (params.Page - 1) * params.PerPage
	dataQuery := fmt.Sprintf(`
		SELECT s.id, s.branch_id, s.item_id, s.customer_id, s.sale_number, s.sale_type,
			   s.sale_price, s.discount_amount, s.discount_reason, s.final_price,
			   s.payment_method, s.reference_number, s.status, s.sale_date,
			   s.refund_amount, s.refund_reason, s.refunded_at, s.refunded_by,
			   s.notes, s.cash_session_id, s.created_by, s.updated_by,
			   s.created_at, s.updated_at, s.deleted_at,
			   i.id, i.name, i.sku,
			   c.id, c.first_name, c.last_name, c.identity_number
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		baseQuery, orderBy, order, argCount+1, argCount+2,
	)
	args = append(args, params.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list sales: %w", err)
	}
	defer rows.Close()

	sales := []domain.Sale{}
	for rows.Next() {
		sale, err := r.scanSaleRowWithRelations(rows)
		if err != nil {
			return nil, err
		}
		sales = append(sales, *sale)
	}

	totalPages := total / params.PerPage
	if total%params.PerPage > 0 {
		totalPages++
	}

	return &repository.PaginatedResult[domain.Sale]{
		Data:       sales,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

// Create creates a new sale
func (r *SaleRepository) Create(ctx context.Context, sale *domain.Sale) error {
	query := `
		INSERT INTO sales (
			branch_id, item_id, customer_id, sale_number, sale_type,
			sale_price, discount_amount, discount_reason, final_price,
			payment_method, reference_number, status, sale_date,
			notes, cash_session_id, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		sale.BranchID, sale.ItemID, NullInt64(sale.CustomerID), sale.SaleNumber, sale.SaleType,
		sale.SalePrice, sale.DiscountAmount, NullStringPtr(sale.DiscountReason), sale.FinalPrice,
		sale.PaymentMethod, NullStringPtr(sale.ReferenceNumber), sale.Status, sale.SaleDate,
		NullStringPtr(sale.Notes), NullInt64(sale.CashSessionID), sale.CreatedBy,
	).Scan(&sale.ID, &sale.CreatedAt, &sale.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create sale: %w", err)
	}

	return nil
}

// Update updates an existing sale
func (r *SaleRepository) Update(ctx context.Context, sale *domain.Sale) error {
	query := `
		UPDATE sales SET
			status = $2, refund_amount = $3, refund_reason = $4,
			refunded_at = $5, refunded_by = $6, notes = $7,
			updated_by = $8, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		sale.ID, sale.Status, NullFloat64(sale.RefundAmount), NullStringPtr(sale.RefundReason),
		sale.RefundedAt, NullInt64(sale.RefundedBy), NullStringPtr(sale.Notes), sale.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to update sale: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("sale not found")
	}

	return nil
}

// GenerateNumber generates a unique sale number
func (r *SaleRepository) GenerateNumber(ctx context.Context) (string, error) {
	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(sale_number FROM 'SL-(\d+)') AS INTEGER)), 0) + 1
		FROM sales
	`

	var seqNum int
	if err := r.db.QueryRowContext(ctx, query).Scan(&seqNum); err != nil {
		seqNum = 1
	}

	return fmt.Sprintf("SL-%08d", seqNum), nil
}

// Helper functions
func (r *SaleRepository) scanSale(row *sql.Row) (*domain.Sale, error) {
	sale := &domain.Sale{}
	var customerID, cashSessionID, refundedBy, createdBy, updatedBy sql.NullInt64
	var discountReason, referenceNumber, refundReason, notes sql.NullString
	var refundAmount sql.NullFloat64
	var refundedAt, deletedAt sql.NullTime

	err := row.Scan(
		&sale.ID, &sale.BranchID, &sale.ItemID, &customerID, &sale.SaleNumber, &sale.SaleType,
		&sale.SalePrice, &sale.DiscountAmount, &discountReason, &sale.FinalPrice,
		&sale.PaymentMethod, &referenceNumber, &sale.Status, &sale.SaleDate,
		&refundAmount, &refundReason, &refundedAt, &refundedBy,
		&notes, &cashSessionID, &createdBy, &updatedBy,
		&sale.CreatedAt, &sale.UpdatedAt, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sale not found")
		}
		return nil, fmt.Errorf("failed to get sale: %w", err)
	}

	sale.CustomerID = Int64Ptr(customerID)
	sale.CashSessionID = Int64Ptr(cashSessionID)
	sale.DiscountReason = StringPtrVal(discountReason)
	sale.ReferenceNumber = StringPtrVal(referenceNumber)
	sale.RefundAmount = Float64Ptr(refundAmount)
	sale.RefundReason = StringPtrVal(refundReason)
	sale.Notes = StringPtrVal(notes)
	if refundedAt.Valid {
		sale.RefundedAt = &refundedAt.Time
	}
	sale.RefundedBy = Int64Ptr(refundedBy)
	if createdBy.Valid {
		sale.CreatedBy = createdBy.Int64
	}
	if updatedBy.Valid {
		sale.UpdatedBy = updatedBy.Int64
	}
	sale.DeletedAt = TimePtr(deletedAt)

	return sale, nil
}

func (r *SaleRepository) scanSaleRow(rows *sql.Rows) (*domain.Sale, error) {
	sale := &domain.Sale{}
	var customerID, cashSessionID, refundedBy, createdBy, updatedBy sql.NullInt64
	var discountReason, referenceNumber, refundReason, notes sql.NullString
	var refundAmount sql.NullFloat64
	var refundedAt, deletedAt sql.NullTime

	err := rows.Scan(
		&sale.ID, &sale.BranchID, &sale.ItemID, &customerID, &sale.SaleNumber, &sale.SaleType,
		&sale.SalePrice, &sale.DiscountAmount, &discountReason, &sale.FinalPrice,
		&sale.PaymentMethod, &referenceNumber, &sale.Status, &sale.SaleDate,
		&refundAmount, &refundReason, &refundedAt, &refundedBy,
		&notes, &cashSessionID, &createdBy, &updatedBy,
		&sale.CreatedAt, &sale.UpdatedAt, &deletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan sale: %w", err)
	}

	sale.CustomerID = Int64Ptr(customerID)
	sale.CashSessionID = Int64Ptr(cashSessionID)
	sale.DiscountReason = StringPtrVal(discountReason)
	sale.ReferenceNumber = StringPtrVal(referenceNumber)
	sale.RefundAmount = Float64Ptr(refundAmount)
	sale.RefundReason = StringPtrVal(refundReason)
	sale.Notes = StringPtrVal(notes)
	if refundedAt.Valid {
		sale.RefundedAt = &refundedAt.Time
	}
	sale.RefundedBy = Int64Ptr(refundedBy)
	if createdBy.Valid {
		sale.CreatedBy = createdBy.Int64
	}
	if updatedBy.Valid {
		sale.UpdatedBy = updatedBy.Int64
	}
	sale.DeletedAt = TimePtr(deletedAt)

	return sale, nil
}

func (r *SaleRepository) scanSaleRowWithRelations(rows *sql.Rows) (*domain.Sale, error) {
	sale := &domain.Sale{}
	var customerID sql.NullInt64
	var discountAmount, salePrice, refundAmount sql.NullFloat64
	var discountReason, referenceNumber, refundReason, notes sql.NullString
	var refundedAt sql.NullTime
	var refundedBy, cashSessionID, createdBy, updatedBy sql.NullInt64
	var deletedAt sql.NullTime

	// Item fields
	var itemID sql.NullInt64
	var itemName, itemSKU sql.NullString

	// Customer fields
	var custID sql.NullInt64
	var custFirstName, custLastName, custIdentityNumber sql.NullString

	err := rows.Scan(
		&sale.ID, &sale.BranchID, &sale.ItemID, &customerID,
		&sale.SaleNumber, &sale.SaleType,
		&salePrice, &discountAmount, &discountReason, &sale.FinalPrice,
		&sale.PaymentMethod, &referenceNumber, &sale.Status, &sale.SaleDate,
		&refundAmount, &refundReason, &refundedAt, &refundedBy,
		&notes, &cashSessionID, &createdBy, &updatedBy,
		&sale.CreatedAt, &sale.UpdatedAt, &deletedAt,
		// Item
		&itemID, &itemName, &itemSKU,
		// Customer
		&custID, &custFirstName, &custLastName, &custIdentityNumber,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan sale with relations: %w", err)
	}

	sale.CustomerID = Int64Ptr(customerID)
	sale.SalePrice = salePrice.Float64
	sale.DiscountAmount = discountAmount.Float64
	sale.DiscountReason = StringPtrVal(discountReason)
	sale.ReferenceNumber = StringPtrVal(referenceNumber)
	sale.RefundAmount = Float64Ptr(refundAmount)
	sale.RefundReason = StringPtrVal(refundReason)
	sale.RefundedAt = TimePtr(refundedAt)
	sale.RefundedBy = Int64Ptr(refundedBy)
	sale.Notes = StringPtrVal(notes)
	sale.CashSessionID = Int64Ptr(cashSessionID)
	if createdBy.Valid {
		sale.CreatedBy = createdBy.Int64
	}
	if updatedBy.Valid {
		sale.UpdatedBy = updatedBy.Int64
	}
	sale.DeletedAt = TimePtr(deletedAt)

	// Populate Item relation
	if itemID.Valid {
		sale.Item = &domain.Item{
			ID:   itemID.Int64,
			Name: itemName.String,
			SKU:  itemSKU.String,
		}
	}

	// Populate Customer relation (optional)
	if custID.Valid {
		sale.Customer = &domain.Customer{
			ID:             custID.Int64,
			FirstName:      custFirstName.String,
			LastName:       custLastName.String,
			IdentityNumber: custIdentityNumber.String,
		}
	}

	return sale, nil
}
