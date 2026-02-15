package domain

import "time"

// Account types
const (
	AccountTypeAsset     = "asset"
	AccountTypeLiability = "liability"
	AccountTypeEquity    = "equity"
	AccountTypeIncome    = "income"
	AccountTypeExpense   = "expense"
)

// Entry types
const (
	EntryTypeDebit  = "debit"
	EntryTypeCredit = "credit"
)

// Account represents a chart of accounts entry
type Account struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	AccountType string `json:"account_type"`
	ParentID    *int64 `json:"parent_id,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
	IsSystem    bool   `json:"is_system"`

	// Relations
	Parent   *Account   `json:"parent,omitempty"`
	Children []*Account `json:"children,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AccountingEntry represents a journal entry
type AccountingEntry struct {
	ID          int64  `json:"id"`
	EntryNumber string `json:"entry_number"`
	BranchID    int64  `json:"branch_id"`
	Branch      *Branch `json:"branch,omitempty"`

	// Entry details
	EntryDate   time.Time `json:"entry_date"`
	Description string    `json:"description"`

	// Reference
	ReferenceType string `json:"reference_type,omitempty"`
	ReferenceID   *int64 `json:"reference_id,omitempty"`

	// Totals
	TotalDebit  float64 `json:"total_debit"`
	TotalCredit float64 `json:"total_credit"`

	// Status
	IsPosted bool       `json:"is_posted"`
	PostedAt *time.Time `json:"posted_at,omitempty"`
	PostedBy *int64     `json:"posted_by,omitempty"`

	// Lines
	Lines []*AccountingEntryLine `json:"lines,omitempty"`

	// Audit
	CreatedBy *int64    `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AccountingEntryLine represents a single debit or credit line
type AccountingEntryLine struct {
	ID          int64    `json:"id"`
	EntryID     int64    `json:"entry_id"`
	AccountID   int64    `json:"account_id"`
	Account     *Account `json:"account,omitempty"`
	EntryType   string   `json:"entry_type"` // debit or credit
	Amount      float64  `json:"amount"`
	Description string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// DailyBalance represents daily financial summary for a branch
type DailyBalance struct {
	ID          int64     `json:"id"`
	BranchID    int64     `json:"branch_id"`
	Branch      *Branch   `json:"branch,omitempty"`
	BalanceDate time.Time `json:"balance_date"`

	// Income
	LoanDisbursements float64 `json:"loan_disbursements"`
	InterestIncome    float64 `json:"interest_income"`
	LateFeeIncome     float64 `json:"late_fee_income"`
	SalesIncome       float64 `json:"sales_income"`
	OtherIncome       float64 `json:"other_income"`

	// Expenses
	OperationalExpenses float64 `json:"operational_expenses"`
	Refunds             float64 `json:"refunds"`
	OtherExpenses       float64 `json:"other_expenses"`

	// Cash position
	CashOpening float64 `json:"cash_opening"`
	CashClosing float64 `json:"cash_closing"`

	// Loan portfolio
	TotalLoansActive float64 `json:"total_loans_active"`
	TotalLoansCount  int     `json:"total_loans_count"`

	// Calculated
	NetIncome float64 `json:"net_income"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TotalIncome calculates total income
func (d *DailyBalance) TotalIncome() float64 {
	return d.InterestIncome + d.LateFeeIncome + d.SalesIncome + d.OtherIncome
}

// TotalExpenses calculates total expenses
func (d *DailyBalance) TotalExpenses() float64 {
	return d.OperationalExpenses + d.Refunds + d.OtherExpenses + d.LoanDisbursements
}

// ExpenseCategory represents a category for expenses
type ExpenseCategory struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Expense represents an operational expense
type Expense struct {
	ID            int64  `json:"id"`
	ExpenseNumber string `json:"expense_number"`
	BranchID      int64  `json:"branch_id"`
	Branch        *Branch `json:"branch,omitempty"`
	CategoryID    *int64 `json:"category_id,omitempty"`
	Category      *ExpenseCategory `json:"category,omitempty"`

	// Details
	Description   string    `json:"description"`
	Amount        float64   `json:"amount"`
	ExpenseDate   time.Time `json:"expense_date"`

	// Payment info
	PaymentMethod string `json:"payment_method"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
	Vendor        string `json:"vendor,omitempty"`

	// Approval
	ApprovedBy *int64     `json:"approved_by,omitempty"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`

	// Audit
	CreatedBy *int64    `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsApproved checks if expense is approved
func (e *Expense) IsApproved() bool {
	return e.ApprovedBy != nil
}
