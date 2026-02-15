package domain

import (
	"time"
)

// SaleStatus represents the status of a sale
type SaleStatus string

const (
	SaleStatusCompleted     SaleStatus = "completed"
	SaleStatusPending       SaleStatus = "pending"
	SaleStatusCancelled     SaleStatus = "cancelled"
	SaleStatusRefunded      SaleStatus = "refunded"
	SaleStatusPartialRefund SaleStatus = "partial_refund"
)

// Sale represents a sale of an item
type Sale struct {
	ID         int64  `json:"id"`
	BranchID   int64  `json:"branch_id"`
	ItemID     int64  `json:"item_id"`
	CustomerID *int64 `json:"customer_id,omitempty"` // Optional customer
	SaleNumber string `json:"sale_number"`
	SaleType   string `json:"sale_type"` // direct, layaway

	// Pricing
	SalePrice       float64  `json:"sale_price"`
	DiscountAmount  float64  `json:"discount_amount"`
	DiscountReason  *string  `json:"discount_reason,omitempty"`
	FinalPrice      float64  `json:"final_price"`

	// Payment
	PaymentMethod   PaymentMethod `json:"payment_method"`
	ReferenceNumber *string       `json:"reference_number,omitempty"`

	// Status
	Status   SaleStatus `json:"status"`
	SaleDate time.Time  `json:"sale_date"`

	// Refund info
	RefundAmount  *float64   `json:"refund_amount,omitempty"`
	RefundReason  *string    `json:"refund_reason,omitempty"`
	RefundedAt    *time.Time `json:"refunded_at,omitempty"`
	RefundedBy    *int64     `json:"refunded_by,omitempty"`

	// Notes
	Notes *string `json:"notes,omitempty"`

	// Cash session reference
	CashSessionID *int64 `json:"cash_session_id,omitempty"`

	// Audit
	CreatedBy int64 `json:"created_by,omitempty"`
	UpdatedBy int64 `json:"updated_by,omitempty"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Branch   *Branch   `json:"branch,omitempty"`
	Item     *Item     `json:"item,omitempty"`
	Customer *Customer `json:"customer,omitempty"`
}

// TableName returns the database table name
func (Sale) TableName() string {
	return "sales"
}

// IsRefunded checks if the sale has been refunded
func (s *Sale) IsRefunded() bool {
	return s.Status == SaleStatusRefunded
}

// CanBeRefunded checks if the sale can be refunded
func (s *Sale) CanBeRefunded() bool {
	return s.Status == SaleStatusCompleted
}
