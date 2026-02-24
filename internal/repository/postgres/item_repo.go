package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// ItemRepository implements repository.ItemRepository
type ItemRepository struct {
	db *DB
}

// NewItemRepository creates a new ItemRepository
func NewItemRepository(db *DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// GetByID retrieves an item by ID
func (r *ItemRepository) GetByID(ctx context.Context, id int64) (*domain.Item, error) {
	query := `
		SELECT id, branch_id, category_id, customer_id, sku, name, description,
			   brand, model, serial_number, color, condition,
			   appraised_value, loan_value, sale_price, status,
			   weight, purity, notes, tags, acquisition_type, acquisition_date, acquisition_price,
			   photos, delivered_at, created_by, updated_by, created_at, updated_at, deleted_at
		FROM items
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanItem(r.db.QueryRowContext(ctx, query, id))
}

// GetBySKU retrieves an item by SKU
func (r *ItemRepository) GetBySKU(ctx context.Context, sku string) (*domain.Item, error) {
	query := `
		SELECT id, branch_id, category_id, customer_id, sku, name, description,
			   brand, model, serial_number, color, condition,
			   appraised_value, loan_value, sale_price, status,
			   weight, purity, notes, tags, acquisition_type, acquisition_date, acquisition_price,
			   photos, delivered_at, created_by, updated_by, created_at, updated_at, deleted_at
		FROM items
		WHERE sku = $1 AND deleted_at IS NULL
	`

	return r.scanItem(r.db.QueryRowContext(ctx, query, sku))
}

// List retrieves items with pagination and filters
func (r *ItemRepository) List(ctx context.Context, params repository.ItemListParams) (*repository.PaginatedResult[domain.Item], error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 20
	}

	// Base query with JOINs for related entities
	fromClause := `
		FROM items i
		LEFT JOIN categories c ON i.category_id = c.id
		LEFT JOIN customers cu ON i.customer_id = cu.id
		LEFT JOIN branches b ON i.branch_id = b.id
		WHERE i.deleted_at IS NULL`
	args := []interface{}{}
	argCount := 0

	if params.BranchID > 0 {
		argCount++
		fromClause += fmt.Sprintf(" AND i.branch_id = $%d", argCount)
		args = append(args, params.BranchID)
	}

	if params.CategoryID != nil {
		argCount++
		fromClause += fmt.Sprintf(" AND i.category_id = $%d", argCount)
		args = append(args, *params.CategoryID)
	}

	if params.CustomerID != nil {
		argCount++
		fromClause += fmt.Sprintf(" AND i.customer_id = $%d", argCount)
		args = append(args, *params.CustomerID)
	}

	if params.Status != nil {
		argCount++
		fromClause += fmt.Sprintf(" AND i.status = $%d", argCount)
		args = append(args, *params.Status)
	}

	if params.Search != "" {
		argCount++
		fromClause += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.sku ILIKE $%d OR i.serial_number ILIKE $%d OR i.brand ILIKE $%d)", argCount, argCount, argCount, argCount)
		args = append(args, "%"+params.Search+"%")
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + fromClause
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count items: %w", err)
	}

	// Get data
	orderBy := "i.created_at"
	if params.OrderBy != "" {
		orderBy = "i." + params.OrderBy
	}
	order := "DESC"
	if params.Order == "asc" {
		order = "ASC"
	}

	offset := (params.Page - 1) * params.PerPage
	dataQuery := fmt.Sprintf(`
		SELECT i.id, i.branch_id, i.category_id, i.customer_id, i.sku, i.name, i.description,
			   i.brand, i.model, i.serial_number, i.color, i.condition,
			   i.appraised_value, i.loan_value, i.sale_price, i.status,
			   i.weight, i.purity, i.notes, i.tags, i.acquisition_type, i.acquisition_date, i.acquisition_price,
			   i.photos, i.delivered_at, i.created_by, i.updated_by, i.created_at, i.updated_at, i.deleted_at,
			   c.id, c.name, c.slug,
			   cu.id, cu.first_name, cu.last_name, cu.identity_number, cu.phone,
			   b.id, b.name, b.code
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		fromClause, orderBy, order, argCount+1, argCount+2,
	)
	args = append(args, params.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}
	defer rows.Close()

	items := []domain.Item{}
	for rows.Next() {
		item, err := r.scanItemRowWithRelations(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	totalPages := total / params.PerPage
	if total%params.PerPage > 0 {
		totalPages++
	}

	return &repository.PaginatedResult[domain.Item]{
		Data:       items,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

// Create creates a new item
func (r *ItemRepository) Create(ctx context.Context, item *domain.Item) error {
	query := `
		INSERT INTO items (
			branch_id, category_id, customer_id, sku, name, description,
			brand, model, serial_number, color, condition,
			appraised_value, loan_value, sale_price, status,
			weight, purity, notes, tags, acquisition_type, acquisition_date, acquisition_price,
			photos, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		item.BranchID, NullInt64(item.CategoryID), NullInt64(item.CustomerID),
		item.SKU, item.Name, NullStringPtr(item.Description),
		NullStringPtr(item.Brand), NullStringPtr(item.Model), NullStringPtr(item.SerialNumber),
		NullStringPtr(item.Color), item.Condition,
		item.AppraisedValue, item.LoanValue, NullFloat64(item.SalePrice), item.Status,
		item.Weight, NullStringPtr(item.Purity), NullStringPtr(item.Notes),
		pq.Array(item.Tags), item.AcquisitionType, item.AcquisitionDate, NullFloat64(item.AcquisitionPrice),
		pq.Array(item.Photos), item.CreatedBy,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return nil
}

// Update updates an existing item
func (r *ItemRepository) Update(ctx context.Context, item *domain.Item) error {
	query := `
		UPDATE items SET
			category_id = $2, name = $3, description = $4,
			brand = $5, model = $6, serial_number = $7, color = $8, condition = $9,
			appraised_value = $10, loan_value = $11, sale_price = $12,
			weight = $13, purity = $14, notes = $15, tags = $16, photos = $17,
			delivered_at = $18, updated_by = $19, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		item.ID, NullInt64(item.CategoryID), item.Name, NullStringPtr(item.Description),
		NullStringPtr(item.Brand), NullStringPtr(item.Model), NullStringPtr(item.SerialNumber),
		NullStringPtr(item.Color), item.Condition,
		item.AppraisedValue, item.LoanValue, NullFloat64(item.SalePrice),
		item.Weight, NullStringPtr(item.Purity), NullStringPtr(item.Notes),
		pq.Array(item.Tags), pq.Array(item.Photos), NullTime(item.DeliveredAt), item.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}

// Delete soft deletes an item
func (r *ItemRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE items SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}

// UpdateStatus updates item status
func (r *ItemRepository) UpdateStatus(ctx context.Context, id int64, status domain.ItemStatus) error {
	query := `UPDATE items SET status = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update item status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}

// GenerateSKU generates a unique SKU for an item
func (r *ItemRepository) GenerateSKU(ctx context.Context, branchID int64) (string, error) {
	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(sku FROM 'IT-\d+-(\d+)') AS INTEGER)), 0) + 1
		FROM items
		WHERE branch_id = $1
	`

	var seqNum int
	if err := r.db.QueryRowContext(ctx, query, branchID).Scan(&seqNum); err != nil {
		seqNum = 1
	}

	return fmt.Sprintf("IT-%d-%06d", branchID, seqNum), nil
}

// CreateHistory creates an item history entry
func (r *ItemRepository) CreateHistory(ctx context.Context, history *domain.ItemHistory) error {
	query := `
		INSERT INTO item_history (item_id, action, old_status, new_status, old_branch_id, new_branch_id,
								  reference_type, reference_id, notes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query,
		history.ItemID, history.Action, history.OldStatus, history.NewStatus,
		NullInt64(history.OldBranchID), NullInt64(history.NewBranchID),
		NullStringPtr(history.ReferenceType), NullInt64(history.ReferenceID),
		NullString(history.Notes), history.CreatedBy,
	).Scan(&history.ID, &history.CreatedAt)

	return err
}

// Helper functions
func (r *ItemRepository) scanItem(row *sql.Row) (*domain.Item, error) {
	item := &domain.Item{}
	var categoryID, customerID sql.NullInt64
	var description, brand, model, serialNumber, color, purity, notes sql.NullString
	var salePrice, acquisitionPrice, weight sql.NullFloat64
	var tags, photos pq.StringArray
	var createdBy, updatedBy sql.NullInt64
	var deletedAt, deliveredAt sql.NullTime

	err := row.Scan(
		&item.ID, &item.BranchID, &categoryID, &customerID,
		&item.SKU, &item.Name, &description,
		&brand, &model, &serialNumber, &color, &item.Condition,
		&item.AppraisedValue, &item.LoanValue, &salePrice, &item.Status,
		&weight, &purity, &notes, &tags,
		&item.AcquisitionType, &item.AcquisitionDate, &acquisitionPrice,
		&photos, &deliveredAt, &createdBy, &updatedBy,
		&item.CreatedAt, &item.UpdatedAt, &deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	item.CategoryID = Int64Ptr(categoryID)
	item.CustomerID = Int64Ptr(customerID)
	item.Description = StringPtrVal(description)
	item.Brand = StringPtrVal(brand)
	item.Model = StringPtrVal(model)
	item.SerialNumber = StringPtrVal(serialNumber)
	item.Color = StringPtrVal(color)
	item.Purity = StringPtrVal(purity)
	item.Notes = StringPtrVal(notes)
	item.SalePrice = Float64Ptr(salePrice)
	item.AcquisitionPrice = Float64Ptr(acquisitionPrice)
	item.Weight = weight.Float64
	item.Tags = tags
	item.Photos = photos
	if createdBy.Valid {
		item.CreatedBy = createdBy.Int64
	}
	if updatedBy.Valid {
		item.UpdatedBy = updatedBy.Int64
	}
	item.DeletedAt = TimePtr(deletedAt)
	item.DeliveredAt = TimePtr(deliveredAt)

	return item, nil
}

func (r *ItemRepository) scanItemRow(rows *sql.Rows) (*domain.Item, error) {
	item := &domain.Item{}
	var categoryID, customerID sql.NullInt64
	var description, brand, model, serialNumber, color, purity, notes sql.NullString
	var salePrice, acquisitionPrice, weight sql.NullFloat64
	var tags, photos pq.StringArray
	var createdBy, updatedBy sql.NullInt64
	var deletedAt, deliveredAt sql.NullTime

	err := rows.Scan(
		&item.ID, &item.BranchID, &categoryID, &customerID,
		&item.SKU, &item.Name, &description,
		&brand, &model, &serialNumber, &color, &item.Condition,
		&item.AppraisedValue, &item.LoanValue, &salePrice, &item.Status,
		&weight, &purity, &notes, &tags,
		&item.AcquisitionType, &item.AcquisitionDate, &acquisitionPrice,
		&photos, &deliveredAt, &createdBy, &updatedBy,
		&item.CreatedAt, &item.UpdatedAt, &deletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan item: %w", err)
	}

	item.CategoryID = Int64Ptr(categoryID)
	item.CustomerID = Int64Ptr(customerID)
	item.Description = StringPtrVal(description)
	item.Brand = StringPtrVal(brand)
	item.Model = StringPtrVal(model)
	item.SerialNumber = StringPtrVal(serialNumber)
	item.Color = StringPtrVal(color)
	item.Purity = StringPtrVal(purity)
	item.Notes = StringPtrVal(notes)
	item.SalePrice = Float64Ptr(salePrice)
	item.AcquisitionPrice = Float64Ptr(acquisitionPrice)
	item.Weight = weight.Float64
	item.Tags = tags
	item.Photos = photos
	if createdBy.Valid {
		item.CreatedBy = createdBy.Int64
	}
	if updatedBy.Valid {
		item.UpdatedBy = updatedBy.Int64
	}
	item.DeletedAt = TimePtr(deletedAt)
	item.DeliveredAt = TimePtr(deliveredAt)

	return item, nil
}

func (r *ItemRepository) scanItemRowWithRelations(rows *sql.Rows) (*domain.Item, error) {
	item := &domain.Item{}
	var categoryID, customerID sql.NullInt64
	var description, brand, model, serialNumber, color, purity, notes sql.NullString
	var salePrice, acquisitionPrice, weight sql.NullFloat64
	var tags, photos pq.StringArray
	var createdBy, updatedBy sql.NullInt64
	var deletedAt, deliveredAt sql.NullTime

	// Category fields
	var catID sql.NullInt64
	var catName, catSlug sql.NullString

	// Customer fields
	var custID sql.NullInt64
	var custFirstName, custLastName, custIdentityNumber, custPhone sql.NullString

	// Branch fields
	var branchID sql.NullInt64
	var branchName, branchCode sql.NullString

	err := rows.Scan(
		&item.ID, &item.BranchID, &categoryID, &customerID,
		&item.SKU, &item.Name, &description,
		&brand, &model, &serialNumber, &color, &item.Condition,
		&item.AppraisedValue, &item.LoanValue, &salePrice, &item.Status,
		&weight, &purity, &notes, &tags,
		&item.AcquisitionType, &item.AcquisitionDate, &acquisitionPrice,
		&photos, &deliveredAt, &createdBy, &updatedBy,
		&item.CreatedAt, &item.UpdatedAt, &deletedAt,
		// Category
		&catID, &catName, &catSlug,
		// Customer
		&custID, &custFirstName, &custLastName, &custIdentityNumber, &custPhone,
		// Branch
		&branchID, &branchName, &branchCode,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan item with relations: %w", err)
	}

	item.CategoryID = Int64Ptr(categoryID)
	item.CustomerID = Int64Ptr(customerID)
	item.Description = StringPtrVal(description)
	item.Brand = StringPtrVal(brand)
	item.Model = StringPtrVal(model)
	item.SerialNumber = StringPtrVal(serialNumber)
	item.Color = StringPtrVal(color)
	item.Purity = StringPtrVal(purity)
	item.Notes = StringPtrVal(notes)
	item.SalePrice = Float64Ptr(salePrice)
	item.AcquisitionPrice = Float64Ptr(acquisitionPrice)
	item.Weight = weight.Float64
	item.Tags = tags
	item.Photos = photos
	if createdBy.Valid {
		item.CreatedBy = createdBy.Int64
	}
	if updatedBy.Valid {
		item.UpdatedBy = updatedBy.Int64
	}
	item.DeletedAt = TimePtr(deletedAt)
	item.DeliveredAt = TimePtr(deliveredAt)

	// Populate Category relation
	if catID.Valid {
		item.Category = &domain.Category{
			ID:   catID.Int64,
			Name: catName.String,
			Slug: catSlug.String,
		}
	}

	// Populate Customer relation
	if custID.Valid {
		item.Customer = &domain.Customer{
			ID:             custID.Int64,
			FirstName:      custFirstName.String,
			LastName:       custLastName.String,
			IdentityNumber: custIdentityNumber.String,
			Phone:          custPhone.String,
		}
	}

	// Populate Branch relation
	if branchID.Valid {
		item.Branch = &domain.Branch{
			ID:   branchID.Int64,
			Name: branchName.String,
			Code: branchCode.String,
		}
	}

	return item, nil
}

// NullTime helper for *time.Time
func NullTimeVal(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: true}
}
