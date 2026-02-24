package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// ItemService handles item business logic
type ItemService struct {
	itemRepo     repository.ItemRepository
	branchRepo   repository.BranchRepository
	categoryRepo repository.CategoryRepository
	customerRepo repository.CustomerRepository
}

// NewItemService creates a new ItemService
func NewItemService(
	itemRepo repository.ItemRepository,
	branchRepo repository.BranchRepository,
	categoryRepo repository.CategoryRepository,
	customerRepo repository.CustomerRepository,
) *ItemService {
	return &ItemService{
		itemRepo:     itemRepo,
		branchRepo:   branchRepo,
		categoryRepo: categoryRepo,
		customerRepo: customerRepo,
	}
}

// CreateItemInput represents create item request data
type CreateItemInput struct {
	BranchID         int64    `json:"branch_id" validate:"required"`
	CategoryID       *int64   `json:"category_id"`
	CustomerID       *int64   `json:"customer_id"`
	Name             string   `json:"name" validate:"required,min=2"`
	Description      *string  `json:"description"`
	Brand            *string  `json:"brand"`
	Model            *string  `json:"model"`
	SerialNumber     *string  `json:"serial_number"`
	Color            *string  `json:"color"`
	Condition        string   `json:"condition" validate:"required,oneof=new excellent good fair poor"`
	AppraisedValue   float64  `json:"appraised_value" validate:"required,gt=0"`
	LoanValue        float64  `json:"loan_value" validate:"required,gt=0"`
	SalePrice        *float64 `json:"sale_price"`
	Weight           float64  `json:"weight"`
	Purity           *string  `json:"purity"`
	Notes            *string  `json:"notes"`
	Tags             []string `json:"tags"`
	AcquisitionType  string   `json:"acquisition_type" validate:"required,oneof=pawn purchase consignment"`
	AcquisitionPrice *float64 `json:"acquisition_price"`
	Photos           []string `json:"photos"`
	CreatedBy        int64    `json:"-"`
}

// Create creates a new item
func (s *ItemService) Create(ctx context.Context, input CreateItemInput) (*domain.Item, error) {
	// Validate branch exists
	_, err := s.branchRepo.GetByID(ctx, input.BranchID)
	if err != nil {
		return nil, errors.New("invalid branch")
	}

	// Validate category if provided
	if input.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(ctx, *input.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category")
		}
	}

	// Validate customer if provided
	if input.CustomerID != nil {
		customer, err := s.customerRepo.GetByID(ctx, *input.CustomerID)
		if err != nil {
			return nil, errors.New("invalid customer")
		}
		if customer.IsBlocked {
			return nil, errors.New("customer is blocked")
		}
	}

	// Validate loan value doesn't exceed appraised value
	if input.LoanValue > input.AppraisedValue {
		return nil, errors.New("loan value cannot exceed appraised value")
	}

	// Generate SKU
	sku, err := s.itemRepo.GenerateSKU(ctx, input.BranchID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate SKU: %w", err)
	}

	// All items start as available - status changes to collateral when loan is created
	status := domain.ItemStatusAvailable

	// Create item
	item := &domain.Item{
		BranchID:         input.BranchID,
		CategoryID:       input.CategoryID,
		CustomerID:       input.CustomerID,
		SKU:              sku,
		Name:             input.Name,
		Description:      input.Description,
		Brand:            input.Brand,
		Model:            input.Model,
		SerialNumber:     input.SerialNumber,
		Color:            input.Color,
		Condition:        input.Condition,
		AppraisedValue:   input.AppraisedValue,
		LoanValue:        input.LoanValue,
		SalePrice:        input.SalePrice,
		Status:           status,
		Weight:           input.Weight,
		Purity:           input.Purity,
		Notes:            input.Notes,
		Tags:             input.Tags,
		AcquisitionType:  input.AcquisitionType,
		AcquisitionDate:  time.Now(),
		AcquisitionPrice: input.AcquisitionPrice,
		Photos:           input.Photos,
		CreatedBy:        input.CreatedBy,
	}

	if err := s.itemRepo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	// Create history entry
	s.itemRepo.CreateHistory(ctx, &domain.ItemHistory{
		ItemID:    item.ID,
		Action:    "created",
		NewStatus: string(status),
		Notes:     stringVal(input.Notes),
		CreatedBy: input.CreatedBy,
	})

	return item, nil
}

// UpdateItemInput represents update item request data
type UpdateItemInput struct {
	CategoryID     *int64   `json:"category_id"`
	Name           string   `json:"name" validate:"omitempty,min=2"`
	Description    *string  `json:"description"`
	Brand          *string  `json:"brand"`
	Model          *string  `json:"model"`
	SerialNumber   *string  `json:"serial_number"`
	Color          *string  `json:"color"`
	Condition      string   `json:"condition" validate:"omitempty,oneof=new excellent good fair poor"`
	AppraisedValue *float64 `json:"appraised_value"`
	LoanValue      *float64 `json:"loan_value"`
	SalePrice      *float64 `json:"sale_price"`
	Weight         *float64 `json:"weight"`
	Purity         *string  `json:"purity"`
	Notes          *string  `json:"notes"`
	Tags           []string `json:"tags"`
	Photos         []string `json:"photos"`
	UpdatedBy      int64    `json:"-"`
}

// Update updates an existing item
func (s *ItemService) Update(ctx context.Context, id int64, input UpdateItemInput) (*domain.Item, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("item not found")
	}

	// Can only update items that are available or pawned
	if item.Status != domain.ItemStatusAvailable && item.Status != domain.ItemStatusPawned {
		return nil, errors.New("cannot update item in current status")
	}

	// Validate category if changed
	if input.CategoryID != nil && (item.CategoryID == nil || *input.CategoryID != *item.CategoryID) {
		_, err := s.categoryRepo.GetByID(ctx, *input.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category")
		}
		item.CategoryID = input.CategoryID
	}

	// Update fields
	if input.Name != "" {
		item.Name = input.Name
	}
	if input.Description != nil {
		item.Description = input.Description
	}
	if input.Brand != nil {
		item.Brand = input.Brand
	}
	if input.Model != nil {
		item.Model = input.Model
	}
	if input.SerialNumber != nil {
		item.SerialNumber = input.SerialNumber
	}
	if input.Color != nil {
		item.Color = input.Color
	}
	if input.Condition != "" {
		item.Condition = input.Condition
	}
	if input.AppraisedValue != nil {
		item.AppraisedValue = *input.AppraisedValue
	}
	if input.LoanValue != nil {
		item.LoanValue = *input.LoanValue
	}
	if input.SalePrice != nil {
		item.SalePrice = input.SalePrice
	}
	if input.Weight != nil {
		item.Weight = *input.Weight
	}
	if input.Purity != nil {
		item.Purity = input.Purity
	}
	if input.Notes != nil {
		item.Notes = input.Notes
	}
	if input.Tags != nil {
		item.Tags = input.Tags
	}
	if input.Photos != nil {
		item.Photos = input.Photos
	}
	item.UpdatedBy = input.UpdatedBy

	// Validate loan value doesn't exceed appraised value
	if item.LoanValue > item.AppraisedValue {
		return nil, errors.New("loan value cannot exceed appraised value")
	}

	if err := s.itemRepo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return item, nil
}

// GetByID retrieves an item by ID
func (s *ItemService) GetByID(ctx context.Context, id int64) (*domain.Item, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("item not found")
	}

	// Load related data
	if item.CategoryID != nil {
		item.Category, _ = s.categoryRepo.GetByID(ctx, *item.CategoryID)
	}
	if item.CustomerID != nil {
		item.Customer, _ = s.customerRepo.GetByID(ctx, *item.CustomerID)
	}

	return item, nil
}

// GetBySKU retrieves an item by SKU
func (s *ItemService) GetBySKU(ctx context.Context, sku string) (*domain.Item, error) {
	item, err := s.itemRepo.GetBySKU(ctx, sku)
	if err != nil {
		return nil, errors.New("item not found")
	}

	return item, nil
}

// List retrieves items with pagination and filters
func (s *ItemService) List(ctx context.Context, params repository.ItemListParams) (*repository.PaginatedResult[domain.Item], error) {
	return s.itemRepo.List(ctx, params)
}

// Delete soft deletes an item
func (s *ItemService) Delete(ctx context.Context, id int64, deletedBy int64) error {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("item not found")
	}

	// Can only delete items that are available
	if item.Status != domain.ItemStatusAvailable {
		return errors.New("can only delete available items")
	}

	// Create history before deletion
	s.itemRepo.CreateHistory(ctx, &domain.ItemHistory{
		ItemID:    item.ID,
		Action:    "deleted",
		OldStatus: string(item.Status),
		CreatedBy: deletedBy,
	})

	return s.itemRepo.Delete(ctx, id)
}

// UpdateStatusInput represents update status request data
type UpdateStatusInput struct {
	Status    domain.ItemStatus `json:"status" validate:"required"`
	Notes     string            `json:"notes"`
	UpdatedBy int64             `json:"-"`
}

// UpdateStatus updates item status
func (s *ItemService) UpdateStatus(ctx context.Context, id int64, input UpdateStatusInput) error {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("item not found")
	}

	oldStatus := item.Status

	// Validate status transition
	if !isValidStatusTransition(oldStatus, input.Status) {
		return errors.New("invalid status transition")
	}

	if err := s.itemRepo.UpdateStatus(ctx, id, input.Status); err != nil {
		return fmt.Errorf("failed to update item status: %w", err)
	}

	// Create history entry
	s.itemRepo.CreateHistory(ctx, &domain.ItemHistory{
		ItemID:    id,
		Action:    "status_changed",
		OldStatus: string(oldStatus),
		NewStatus: string(input.Status),
		Notes:     input.Notes,
		CreatedBy: input.UpdatedBy,
	})

	return nil
}

// MarkForSale marks an item as available for sale
func (s *ItemService) MarkForSale(ctx context.Context, id int64, salePrice float64, updatedBy int64) error {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("item not found")
	}

	// Can only mark confiscated or available items for sale
	if item.Status != domain.ItemStatusConfiscated && item.Status != domain.ItemStatusAvailable {
		return errors.New("can only mark confiscated or available items for sale")
	}

	// Cannot mark for sale items that were delivered to customer
	if item.AcquisitionType == domain.AcquisitionTypePawn && item.DeliveredAt != nil {
		return errors.New("cannot mark for sale items that have been delivered to customer")
	}

	item.SalePrice = &salePrice
	item.Status = domain.ItemStatusForSale
	item.UpdatedBy = updatedBy

	if err := s.itemRepo.Update(ctx, item); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	s.itemRepo.CreateHistory(ctx, &domain.ItemHistory{
		ItemID:    id,
		Action:    "marked_for_sale",
		OldStatus: string(item.Status),
		NewStatus: string(domain.ItemStatusForSale),
		Notes:     "Marked for sale at price: " + formatCurrency(salePrice),
		CreatedBy: updatedBy,
	})

	return nil
}

// GetAvailableForSale retrieves items available for sale
func (s *ItemService) GetAvailableForSale(ctx context.Context, branchID int64) ([]*domain.Item, error) {
	status := domain.ItemStatusForSale
	result, err := s.itemRepo.List(ctx, repository.ItemListParams{
		BranchID: branchID,
		Status:   &status,
		PaginationParams: repository.PaginationParams{
			Page:    1,
			PerPage: 1000, // Get all items for sale
			OrderBy: "created_at",
			Order:   "desc",
		},
	})
	if err != nil {
		return nil, err
	}

	items := make([]*domain.Item, len(result.Data))
	for i := range result.Data {
		items[i] = &result.Data[i]
	}
	return items, nil
}

// isValidStatusTransition validates if a status transition is allowed
func isValidStatusTransition(from, to domain.ItemStatus) bool {
	transitions := map[domain.ItemStatus][]domain.ItemStatus{
		domain.ItemStatusAvailable:   {domain.ItemStatusPawned, domain.ItemStatusForSale, domain.ItemStatusSold, domain.ItemStatusTransferred},
		domain.ItemStatusPawned:      {domain.ItemStatusAvailable, domain.ItemStatusConfiscated},
		domain.ItemStatusForSale:     {domain.ItemStatusSold, domain.ItemStatusAvailable},
		domain.ItemStatusSold:        {}, // Cannot transition from sold
		domain.ItemStatusConfiscated: {domain.ItemStatusForSale, domain.ItemStatusAvailable},
		domain.ItemStatusTransferred: {domain.ItemStatusAvailable},
	}

	allowedTransitions, ok := transitions[from]
	if !ok {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == to {
			return true
		}
	}
	return false
}

// AddPhoto adds a photo URL to an item
func (s *ItemService) AddPhoto(ctx context.Context, itemID int64, photoURL string, updatedBy int64) error {
	item, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return errors.New("item not found")
	}

	// Check if photo already exists
	for _, p := range item.Photos {
		if p == photoURL {
			return nil // Already exists, no-op
		}
	}

	item.Photos = append(item.Photos, photoURL)
	item.UpdatedBy = updatedBy

	if err := s.itemRepo.Update(ctx, item); err != nil {
		return fmt.Errorf("failed to add photo: %w", err)
	}

	return nil
}

// RemovePhoto removes a photo URL from an item
func (s *ItemService) RemovePhoto(ctx context.Context, itemID int64, photoURL string, updatedBy int64) error {
	item, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return errors.New("item not found")
	}

	// Find and remove the photo
	found := false
	newPhotos := make([]string, 0, len(item.Photos))
	for _, p := range item.Photos {
		if p == photoURL {
			found = true
		} else {
			newPhotos = append(newPhotos, p)
		}
	}

	if !found {
		return errors.New("photo not found")
	}

	item.Photos = newPhotos
	item.UpdatedBy = updatedBy

	if err := s.itemRepo.Update(ctx, item); err != nil {
		return fmt.Errorf("failed to remove photo: %w", err)
	}

	return nil
}

// MarkAsDeliveredInput represents mark as delivered request data
type MarkAsDeliveredInput struct {
	ItemID    int64  `json:"item_id" validate:"required"`
	Notes     string `json:"notes"`
	UpdatedBy int64  `json:"-"`
}

// MarkAsDelivered marks an item as physically delivered to the customer
func (s *ItemService) MarkAsDelivered(ctx context.Context, input MarkAsDeliveredInput) (*domain.Item, error) {
	item, err := s.itemRepo.GetByID(ctx, input.ItemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	// Can only mark available items as delivered
	if item.Status != domain.ItemStatusAvailable {
		return nil, errors.New("only available items can be marked as delivered")
	}

	// Item should belong to a customer (pawn acquisition)
	if item.AcquisitionType != domain.AcquisitionTypePawn || item.CustomerID == nil {
		return nil, errors.New("only pawned items with customer can be delivered")
	}

	// Check if already delivered
	if item.DeliveredAt != nil {
		return nil, errors.New("item already marked as delivered")
	}

	// Mark as delivered
	now := time.Now()
	item.DeliveredAt = &now
	item.UpdatedBy = input.UpdatedBy

	if err := s.itemRepo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to mark item as delivered: %w", err)
	}

	// Create history entry
	notes := input.Notes
	if notes == "" {
		notes = "Item physically delivered to customer"
	}
	s.itemRepo.CreateHistory(ctx, &domain.ItemHistory{
		ItemID:    input.ItemID,
		Action:    "delivered",
		OldStatus: string(item.Status),
		NewStatus: string(item.Status),
		Notes:     notes,
		CreatedBy: input.UpdatedBy,
	})

	return item, nil
}

// GetPendingDeliveries retrieves items that are paid but not yet delivered
func (s *ItemService) GetPendingDeliveries(ctx context.Context, branchID int64) ([]*domain.Item, error) {
	status := domain.ItemStatusAvailable
	result, err := s.itemRepo.List(ctx, repository.ItemListParams{
		BranchID: branchID,
		Status:   &status,
		PaginationParams: repository.PaginationParams{
			Page:    1,
			PerPage: 1000,
			OrderBy: "updated_at",
			Order:   "desc",
		},
	})
	if err != nil {
		return nil, err
	}

	// Filter for pawn items not yet delivered
	items := make([]*domain.Item, 0)
	for i := range result.Data {
		item := &result.Data[i]
		if item.IsPendingDelivery() {
			items = append(items, item)
		}
	}

	return items, nil
}

// Helper functions
func stringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func formatCurrency(amount float64) string {
	return "Q" + formatFloat(amount)
}

func formatFloat(f float64) string {
	return string(rune(int(f*100)/100)) // Simple format
}
