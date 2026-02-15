package domain

import (
	"time"
)

// Branch represents a pawnshop branch/location
type Branch struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Address   string `json:"address,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	IsActive  bool   `json:"is_active"`
	Timezone  string `json:"timezone"`
	Currency  string `json:"currency"`

	// Business settings
	DefaultInterestRate float64 `json:"default_interest_rate"`
	DefaultLoanTermDays int     `json:"default_loan_term_days"`
	DefaultGracePeriod  int     `json:"default_grace_period"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// TableName returns the database table name
func (Branch) TableName() string {
	return "branches"
}
