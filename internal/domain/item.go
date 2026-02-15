package domain

import (
	"time"
)

// ItemStatus represents the status of an item
type ItemStatus string

const (
	ItemStatusAvailable   ItemStatus = "available"
	ItemStatusPawned      ItemStatus = "pawned"
	ItemStatusCollateral  ItemStatus = "collateral"
	ItemStatusForSale     ItemStatus = "for_sale"
	ItemStatusSold        ItemStatus = "sold"
	ItemStatusConfiscated ItemStatus = "confiscated"
	ItemStatusTransferred ItemStatus = "transferred"
	ItemStatusInTransfer  ItemStatus = "in_transfer"
	ItemStatusDamaged     ItemStatus = "damaged"
	ItemStatusLost        ItemStatus = "lost"
)

// ItemCondition represents the condition of an item
type ItemCondition string

const (
	ItemConditionNew       ItemCondition = "new"
	ItemConditionExcellent ItemCondition = "excellent"
	ItemConditionGood      ItemCondition = "good"
	ItemConditionFair      ItemCondition = "fair"
	ItemConditionPoor      ItemCondition = "poor"
)

// Item represents a pawn item
type Item struct {
	ID         int64  `json:"id"`
	BranchID   int64  `json:"branch_id"`
	CategoryID *int64 `json:"category_id,omitempty"`
	CustomerID *int64 `json:"customer_id,omitempty"` // Original owner

	// Identification
	SKU          string  `json:"sku"`
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	Brand        *string `json:"brand,omitempty"`
	Model        *string `json:"model,omitempty"`
	SerialNumber *string `json:"serial_number,omitempty"`
	Color        *string `json:"color,omitempty"`
	Condition    string  `json:"condition"`

	// Valuation
	AppraisedValue float64  `json:"appraised_value"` // Market value
	LoanValue      float64  `json:"loan_value"`      // Maximum loan amount
	SalePrice      *float64 `json:"sale_price,omitempty"`

	// Status
	Status ItemStatus `json:"status"`

	// Physical details (for jewelry)
	Weight float64  `json:"weight,omitempty"` // In grams
	Purity *string  `json:"purity,omitempty"` // e.g., 14k, 18k, .925

	// Additional info
	Notes *string  `json:"notes,omitempty"`
	Tags  []string `json:"tags,omitempty"`

	// Acquisition info
	AcquisitionType  string     `json:"acquisition_type"` // pawn, purchase, confiscation
	AcquisitionDate  time.Time  `json:"acquisition_date"`
	AcquisitionPrice *float64   `json:"acquisition_price,omitempty"`

	// Media
	Photos []string `json:"photos,omitempty"`

	// Audit
	CreatedBy int64 `json:"created_by,omitempty"`
	UpdatedBy int64 `json:"updated_by,omitempty"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Branch   *Branch   `json:"branch,omitempty"`
	Category *Category `json:"category,omitempty"`
	Customer *Customer `json:"customer,omitempty"`
}

// TableName returns the database table name
func (Item) TableName() string {
	return "items"
}

// IsAvailable checks if the item is available for loan or sale
func (i *Item) IsAvailable() bool {
	return i.Status == ItemStatusAvailable
}

// CanBeSold checks if the item can be sold
func (i *Item) CanBeSold() bool {
	return i.Status == ItemStatusAvailable || i.Status == ItemStatusConfiscated
}

// Acquisition type constants
const (
	AcquisitionTypePawn         = "pawn"
	AcquisitionTypePurchase     = "purchase"
	AcquisitionTypeConfiscation = "confiscation"
)

// ItemHistory represents a history entry for an item
type ItemHistory struct {
	ID            int64     `json:"id"`
	ItemID        int64     `json:"item_id"`
	Action        string    `json:"action"`
	OldStatus     string    `json:"old_status,omitempty"`
	NewStatus     string    `json:"new_status,omitempty"`
	OldBranchID   *int64    `json:"old_branch_id,omitempty"`
	NewBranchID   *int64    `json:"new_branch_id,omitempty"`
	ReferenceType *string   `json:"reference_type,omitempty"`
	ReferenceID   *int64    `json:"reference_id,omitempty"`
	Notes         string    `json:"notes,omitempty"`
	CreatedBy     int64     `json:"created_by,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// TableName returns the database table name
func (ItemHistory) TableName() string {
	return "item_history"
}
