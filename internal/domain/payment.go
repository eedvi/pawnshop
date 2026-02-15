package domain

import (
	"time"
)

// PaymentMethod represents the method of payment
type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodCard     PaymentMethod = "card"
	PaymentMethodTransfer PaymentMethod = "transfer"
	PaymentMethodCheck    PaymentMethod = "check"
	PaymentMethodOther    PaymentMethod = "other"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusReversed  PaymentStatus = "reversed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

// Payment represents a loan payment
type Payment struct {
	ID            int64  `json:"id"`
	PaymentNumber string `json:"payment_number"`
	BranchID      int64  `json:"branch_id"`
	LoanID        int64  `json:"loan_id"`
	CustomerID    int64  `json:"customer_id"`

	// Payment details
	Amount          float64 `json:"amount"`
	PrincipalAmount float64 `json:"principal_amount"`
	InterestAmount  float64 `json:"interest_amount"`
	LateFeeAmount   float64 `json:"late_fee_amount"`

	// Method
	PaymentMethod   PaymentMethod `json:"payment_method"`
	ReferenceNumber string        `json:"reference_number,omitempty"`

	// Status
	Status      PaymentStatus `json:"status"`
	PaymentDate time.Time     `json:"payment_date"`

	// Balances after payment
	LoanBalanceAfter     float64 `json:"loan_balance_after"`
	InterestBalanceAfter float64 `json:"interest_balance_after"`

	// Reversal info
	ReversedAt      *time.Time `json:"reversed_at,omitempty"`
	ReversedBy      *int64     `json:"reversed_by,omitempty"`
	ReversalReason  string     `json:"reversal_reason,omitempty"`

	// Notes
	Notes string `json:"notes,omitempty"`

	// Cash session reference
	CashSessionID *int64 `json:"cash_session_id,omitempty"`

	// Audit
	CreatedBy int64 `json:"created_by,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Branch   *Branch   `json:"branch,omitempty"`
	Loan     *Loan     `json:"loan,omitempty"`
	Customer *Customer `json:"customer,omitempty"`
}

// TableName returns the database table name
func (Payment) TableName() string {
	return "payments"
}

// IsReversed checks if the payment has been reversed
func (p *Payment) IsReversed() bool {
	return p.Status == PaymentStatusReversed
}

// CanBeReversed checks if the payment can be reversed
func (p *Payment) CanBeReversed() bool {
	return p.Status == PaymentStatusCompleted
}
