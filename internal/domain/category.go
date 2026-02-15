package domain

import (
	"time"
)

// Category represents an item category
type Category struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
	ParentID    *int64  `json:"parent_id,omitempty"`
	Icon        *string `json:"icon,omitempty"`

	// Loan settings
	DefaultInterestRate float64  `json:"default_interest_rate"`
	MinLoanAmount       *float64 `json:"min_loan_amount,omitempty"`
	MaxLoanAmount       *float64 `json:"max_loan_amount,omitempty"`
	LoanToValueRatio    float64  `json:"loan_to_value_ratio"`

	// Display
	SortOrder int  `json:"sort_order"`
	IsActive  bool `json:"is_active"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Parent   *Category   `json:"parent,omitempty"`
	Children []*Category `json:"children,omitempty"`
}

// TableName returns the database table name
func (Category) TableName() string {
	return "categories"
}

// IsRoot checks if the category is a root category (no parent)
func (c *Category) IsRoot() bool {
	return c.ParentID == nil
}
