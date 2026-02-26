package domain

import (
	"time"
)

// CashSessionStatus represents the status of a cash session
type CashSessionStatus string

const (
	CashSessionStatusOpen   CashSessionStatus = "open"
	CashSessionStatusClosed CashSessionStatus = "closed"
)

// CashMovementType represents the type of cash movement
type CashMovementType string

const (
	CashMovementTypeIncome  CashMovementType = "income"
	CashMovementTypeExpense CashMovementType = "expense"
)

// CashRegister represents a physical cash register
type CashRegister struct {
	ID          int64   `json:"id"`
	BranchID    int64   `json:"branch_id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"is_active"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Branch *Branch `json:"branch,omitempty"`
}

// TableName returns the database table name
func (CashRegister) TableName() string {
	return "cash_registers"
}

// CashSession represents a cash register session
type CashSession struct {
	ID               int64 `json:"id"`
	BranchID         int64 `json:"branch_id"`
	CashRegisterID   int64 `json:"cash_register_id"`
	UserID           int64 `json:"user_id"`

	// Amounts
	OpeningAmount  float64  `json:"opening_amount"`
	ClosingAmount  *float64 `json:"closing_amount,omitempty"`
	ExpectedAmount *float64 `json:"expected_amount,omitempty"`
	Difference     *float64 `json:"difference,omitempty"`

	// Status
	Status CashSessionStatus `json:"status"`

	// Timestamps
	OpenedAt time.Time  `json:"opened_at"`
	ClosedAt *time.Time `json:"closed_at,omitempty"`

	// Notes
	OpeningNotes *string `json:"opening_notes,omitempty"`
	ClosingNotes *string `json:"closing_notes,omitempty"`

	// Audit
	ClosedBy *int64 `json:"closed_by,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	CashRegister *CashRegister   `json:"register,omitempty"`
	User         *User           `json:"user,omitempty"`
	Branch       *Branch         `json:"branch,omitempty"`
	Movements    []*CashMovement `json:"movements,omitempty"`
}

// TableName returns the database table name
func (CashSession) TableName() string {
	return "cash_sessions"
}

// IsOpen checks if the session is open
func (cs *CashSession) IsOpen() bool {
	return cs.Status == CashSessionStatusOpen
}

// CashMovement represents a cash movement in a session
type CashMovement struct {
	ID        int64 `json:"id"`
	BranchID  int64 `json:"branch_id"`
	SessionID int64 `json:"session_id"`

	// Movement details
	MovementType  CashMovementType `json:"movement_type"`
	Amount        float64          `json:"amount"`
	PaymentMethod PaymentMethod    `json:"payment_method"`

	// Reference to related entity
	ReferenceType *string `json:"reference_type,omitempty"` // loan, payment, sale
	ReferenceID   *int64  `json:"reference_id,omitempty"`

	// Description
	Description string `json:"description"`

	// Balance
	BalanceAfter float64 `json:"balance_after"`

	// Audit
	CreatedBy int64     `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	CashSession *CashSession `json:"cash_session,omitempty"`
	Branch      *Branch      `json:"branch,omitempty"`
}

// TableName returns the database table name
func (CashMovement) TableName() string {
	return "cash_movements"
}

// IsIncome checks if the movement is an income
func (cm *CashMovement) IsIncome() bool {
	return cm.MovementType == CashMovementTypeIncome
}

// IsExpense checks if the movement is an expense
func (cm *CashMovement) IsExpense() bool {
	return cm.MovementType == CashMovementTypeExpense
}
