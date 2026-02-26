package domain

import (
	"time"
)

// LoanStatus represents the status of a loan
type LoanStatus string

const (
	LoanStatusActive      LoanStatus = "active"
	LoanStatusPaid        LoanStatus = "paid"
	LoanStatusOverdue     LoanStatus = "overdue"
	LoanStatusDefaulted   LoanStatus = "defaulted"
	LoanStatusRenewed     LoanStatus = "renewed"
	LoanStatusConfiscated LoanStatus = "confiscated"
)

// PaymentPlanType represents the type of payment plan
type PaymentPlanType string

const (
	PaymentPlanSingle         PaymentPlanType = "single"
	PaymentPlanMinimumPayment PaymentPlanType = "minimum_payment"
	PaymentPlanInstallments   PaymentPlanType = "installments"
)

// Loan represents a pawn loan
type Loan struct {
	ID         int64 `json:"id"`
	LoanNumber string `json:"loan_number"`
	BranchID   int64 `json:"branch_id"`
	CustomerID int64 `json:"customer_id"`
	ItemID     int64 `json:"item_id"`

	// Amounts
	LoanAmount         float64 `json:"loan_amount"`
	InterestRate       float64 `json:"interest_rate"`
	InterestAmount     float64 `json:"interest_amount"`
	PrincipalRemaining float64 `json:"principal_remaining"`
	InterestRemaining  float64 `json:"interest_remaining"`
	TotalAmount        float64 `json:"total_amount"`
	AmountPaid         float64 `json:"amount_paid"`

	// Late fees
	LateFeeRate   float64 `json:"late_fee_rate"`
	LateFeeAmount float64 `json:"late_fee_amount"`

	// Dates
	StartDate       Date       `json:"start_date"`
	DueDate         Date       `json:"due_date"`
	PaidDate        *time.Time `json:"paid_date,omitempty"`
	ConfiscatedDate *time.Time `json:"confiscated_date,omitempty"`

	// Payment plan
	PaymentPlanType        PaymentPlanType `json:"payment_plan_type"`
	LoanTermDays           int             `json:"loan_term_days"`
	RequiresMinimumPayment bool            `json:"requires_minimum_payment"`
	MinimumPaymentAmount   *float64        `json:"minimum_payment_amount,omitempty"`
	NextPaymentDueDate     *time.Time      `json:"next_payment_due_date,omitempty"`
	GracePeriodDays        int             `json:"grace_period_days"`

	// Installments
	NumberOfInstallments *int     `json:"number_of_installments,omitempty"`
	InstallmentAmount    *float64 `json:"installment_amount,omitempty"`

	// Status
	Status      LoanStatus `json:"status"`
	DaysOverdue int        `json:"days_overdue"`

	// Renewal info
	RenewedFromID *int64 `json:"renewed_from_id,omitempty"`
	RenewalCount  int    `json:"renewal_count"`

	// Notes
	Notes string `json:"notes,omitempty"`

	// Audit
	CreatedBy int64  `json:"created_by,omitempty"`
	UpdatedBy *int64 `json:"updated_by,omitempty"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Branch   *Branch   `json:"branch,omitempty"`
	Customer *Customer `json:"customer,omitempty"`
	Item     *Item     `json:"item,omitempty"`
}

// TableName returns the database table name
func (Loan) TableName() string {
	return "loans"
}

// RemainingBalance returns the total remaining balance
func (l *Loan) RemainingBalance() float64 {
	return l.PrincipalRemaining + l.InterestRemaining + l.LateFeeAmount
}

// IsOverdue checks if the loan is overdue
func (l *Loan) IsOverdue() bool {
	if l.Status != LoanStatusActive {
		return false
	}
	return time.Now().After(l.DueDate.Time)
}

// IsInGracePeriod checks if the loan is in grace period
func (l *Loan) IsInGracePeriod() bool {
	if !l.IsOverdue() {
		return false
	}
	gracePeriodEnd := l.DueDate.AddDate(0, 0, l.GracePeriodDays)
	return time.Now().Before(gracePeriodEnd)
}

// DaysUntilDue returns the number of days until due date
func (l *Loan) DaysUntilDue() int {
	days := int(time.Until(l.DueDate.Time).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

// CalculateDaysOverdue returns the number of days overdue
func (l *Loan) CalculateDaysOverdue() int {
	if !l.IsOverdue() {
		return 0
	}
	return int(time.Since(l.DueDate.Time).Hours() / 24)
}

// LoanInstallment represents an installment for a loan
type LoanInstallment struct {
	ID                int64     `json:"id"`
	LoanID            int64     `json:"loan_id"`
	InstallmentNumber int       `json:"installment_number"`
	DueDate           time.Time `json:"due_date"`
	PrincipalAmount   float64   `json:"principal_amount"`
	InterestAmount    float64   `json:"interest_amount"`
	TotalAmount       float64   `json:"total_amount"`
	AmountPaid        float64   `json:"amount_paid"`
	IsPaid            bool      `json:"is_paid"`
	PaidDate          *time.Time `json:"paid_date,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName returns the database table name
func (LoanInstallment) TableName() string {
	return "loan_installments"
}

// RemainingAmount returns the remaining amount for this installment
func (li *LoanInstallment) RemainingAmount() float64 {
	return li.TotalAmount - li.AmountPaid
}
