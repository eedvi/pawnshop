package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// SaleService handles sale business logic
type SaleService struct {
	saleRepo     repository.SaleRepository
	itemRepo     repository.ItemRepository
	customerRepo repository.CustomerRepository
	branchRepo   repository.BranchRepository
}

// NewSaleService creates a new SaleService
func NewSaleService(
	saleRepo repository.SaleRepository,
	itemRepo repository.ItemRepository,
	customerRepo repository.CustomerRepository,
	branchRepo repository.BranchRepository,
) *SaleService {
	return &SaleService{
		saleRepo:     saleRepo,
		itemRepo:     itemRepo,
		customerRepo: customerRepo,
		branchRepo:   branchRepo,
	}
}

// CreateSaleInput represents create sale request data
type CreateSaleInput struct {
	BranchID        int64   `json:"branch_id" validate:"required"`
	ItemID          int64   `json:"item_id" validate:"required"`
	CustomerID      *int64  `json:"customer_id"`
	SaleType        string  `json:"sale_type" validate:"required,oneof=direct layaway"`
	DiscountAmount  float64 `json:"discount_amount" validate:"gte=0"`
	DiscountReason  *string `json:"discount_reason"`
	PaymentMethod   string  `json:"payment_method" validate:"required,oneof=cash card transfer check other"`
	ReferenceNumber *string `json:"reference_number"`
	Notes           *string `json:"notes"`
	CashSessionID   *int64  `json:"cash_session_id"`
	CreatedBy       int64   `json:"-"`
}

// SaleResult contains the result of a sale
type SaleResult struct {
	Sale *domain.Sale `json:"sale"`
	Item *domain.Item `json:"item"`
}

// Create creates a new sale
func (s *SaleService) Create(ctx context.Context, input CreateSaleInput) (*SaleResult, error) {
	// Validate branch
	_, err := s.branchRepo.GetByID(ctx, input.BranchID)
	if err != nil {
		return nil, errors.New("invalid branch")
	}

	// Get and validate item
	item, err := s.itemRepo.GetByID(ctx, input.ItemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	// Check item is for sale
	if item.Status != domain.ItemStatusForSale {
		return nil, errors.New("item is not available for sale")
	}

	// Check item belongs to branch
	if item.BranchID != input.BranchID {
		return nil, errors.New("item does not belong to this branch")
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

	// Get sale price from item
	if item.SalePrice == nil || *item.SalePrice <= 0 {
		return nil, errors.New("item does not have a valid sale price")
	}
	salePrice := *item.SalePrice

	// Validate discount
	if input.DiscountAmount > salePrice {
		return nil, errors.New("discount cannot exceed sale price")
	}

	// Calculate final price
	finalPrice := salePrice - input.DiscountAmount

	// Generate sale number
	saleNumber, err := s.saleRepo.GenerateNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate sale number: %w", err)
	}

	// Create sale record
	sale := &domain.Sale{
		BranchID:        input.BranchID,
		ItemID:          input.ItemID,
		CustomerID:      input.CustomerID,
		SaleNumber:      saleNumber,
		SaleType:        input.SaleType,
		SalePrice:       salePrice,
		DiscountAmount:  input.DiscountAmount,
		DiscountReason:  input.DiscountReason,
		FinalPrice:      finalPrice,
		PaymentMethod:   domain.PaymentMethod(input.PaymentMethod),
		ReferenceNumber: input.ReferenceNumber,
		Status:          domain.SaleStatusCompleted,
		SaleDate:        time.Now(),
		Notes:           input.Notes,
		CashSessionID:   input.CashSessionID,
		CreatedBy:       input.CreatedBy,
	}

	// Save sale
	if err := s.saleRepo.Create(ctx, sale); err != nil {
		return nil, fmt.Errorf("failed to create sale: %w", err)
	}

	// Update item status to sold
	if err := s.itemRepo.UpdateStatus(ctx, item.ID, domain.ItemStatusSold); err != nil {
		return nil, fmt.Errorf("failed to update item status: %w", err)
	}

	// Create item history
	s.itemRepo.CreateHistory(ctx, &domain.ItemHistory{
		ItemID:        item.ID,
		Action:        "sold",
		OldStatus:     string(domain.ItemStatusForSale),
		NewStatus:     string(domain.ItemStatusSold),
		ReferenceType: strPtr("sale"),
		ReferenceID:   &sale.ID,
		Notes:         "Sold via sale: " + saleNumber,
		CreatedBy:     input.CreatedBy,
	})

	// Reload item with updated status
	item, _ = s.itemRepo.GetByID(ctx, item.ID)

	return &SaleResult{
		Sale: sale,
		Item: item,
	}, nil
}

// GetByID retrieves a sale by ID
func (s *SaleService) GetByID(ctx context.Context, id int64) (*domain.Sale, error) {
	sale, err := s.saleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("sale not found")
	}

	// Load related data
	sale.Item, _ = s.itemRepo.GetByID(ctx, sale.ItemID)
	if sale.CustomerID != nil {
		sale.Customer, _ = s.customerRepo.GetByID(ctx, *sale.CustomerID)
	}

	return sale, nil
}

// GetByNumber retrieves a sale by sale number
func (s *SaleService) GetByNumber(ctx context.Context, saleNumber string) (*domain.Sale, error) {
	sale, err := s.saleRepo.GetByNumber(ctx, saleNumber)
	if err != nil {
		return nil, errors.New("sale not found")
	}

	return sale, nil
}

// List retrieves sales with pagination and filters
func (s *SaleService) List(ctx context.Context, params repository.SaleListParams) (*repository.PaginatedResult[domain.Sale], error) {
	return s.saleRepo.List(ctx, params)
}

// RefundSaleInput represents refund sale request data
type RefundSaleInput struct {
	SaleID       int64   `json:"sale_id" validate:"required"`
	RefundAmount float64 `json:"refund_amount" validate:"required,gt=0"`
	Reason       string  `json:"reason" validate:"required"`
	RefundedBy   int64   `json:"-"`
}

// Refund processes a sale refund
func (s *SaleService) Refund(ctx context.Context, input RefundSaleInput) (*domain.Sale, error) {
	// Get sale
	sale, err := s.saleRepo.GetByID(ctx, input.SaleID)
	if err != nil {
		return nil, errors.New("sale not found")
	}

	// Validate sale can be refunded
	if sale.Status != domain.SaleStatusCompleted {
		return nil, errors.New("only completed sales can be refunded")
	}

	// Validate refund amount
	if input.RefundAmount > sale.FinalPrice {
		return nil, errors.New("refund amount cannot exceed sale price")
	}

	// Get item
	item, err := s.itemRepo.GetByID(ctx, sale.ItemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	// Check if item is still sold (not transferred to another customer)
	if item.Status != domain.ItemStatusSold {
		return nil, errors.New("item status has changed, cannot process refund")
	}

	// Determine refund type
	isFullRefund := input.RefundAmount >= sale.FinalPrice

	// Update sale
	now := time.Now()
	if isFullRefund {
		sale.Status = domain.SaleStatusRefunded
	} else {
		sale.Status = domain.SaleStatusPartialRefund
	}
	sale.RefundAmount = &input.RefundAmount
	sale.RefundReason = &input.Reason
	sale.RefundedAt = &now
	sale.RefundedBy = &input.RefundedBy
	sale.UpdatedBy = input.RefundedBy

	if err := s.saleRepo.Update(ctx, sale); err != nil {
		return nil, fmt.Errorf("failed to update sale: %w", err)
	}

	// If full refund, return item to available status
	if isFullRefund {
		if err := s.itemRepo.UpdateStatus(ctx, item.ID, domain.ItemStatusAvailable); err != nil {
			return nil, fmt.Errorf("failed to update item status: %w", err)
		}

		// Create item history
		s.itemRepo.CreateHistory(ctx, &domain.ItemHistory{
			ItemID:        item.ID,
			Action:        "returned",
			OldStatus:     string(domain.ItemStatusSold),
			NewStatus:     string(domain.ItemStatusAvailable),
			ReferenceType: strPtr("sale_refund"),
			ReferenceID:   &sale.ID,
			Notes:         "Full refund for sale: " + sale.SaleNumber,
			CreatedBy:     input.RefundedBy,
		})
	}

	return sale, nil
}

// GetSalesSummary retrieves sales summary for a branch
func (s *SaleService) GetSalesSummary(ctx context.Context, branchID int64, dateFrom, dateTo string) (*SalesSummary, error) {
	completedStatus := domain.SaleStatusCompleted
	result, err := s.saleRepo.List(ctx, repository.SaleListParams{
		BranchID: branchID,
		Status:   &completedStatus,
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
		PaginationParams: repository.PaginationParams{
			Page:    1,
			PerPage: 10000, // Get all sales for summary
		},
	})
	if err != nil {
		return nil, err
	}

	summary := &SalesSummary{
		TotalSales:   len(result.Data),
		TotalRevenue: 0,
		ByMethod:     make(map[string]float64),
	}

	for _, sale := range result.Data {
		summary.TotalRevenue += sale.FinalPrice
		summary.TotalDiscount += sale.DiscountAmount
		summary.ByMethod[string(sale.PaymentMethod)] += sale.FinalPrice
	}

	return summary, nil
}

// SalesSummary contains sales summary data
type SalesSummary struct {
	TotalSales    int                `json:"total_sales"`
	TotalRevenue  float64            `json:"total_revenue"`
	TotalDiscount float64            `json:"total_discount"`
	ByMethod      map[string]float64 `json:"by_method"`
}

// Helper function
func strPtr(s string) *string {
	return &s
}
