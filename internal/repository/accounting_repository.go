package repository

import (
	"context"
	"time"

	"pawnshop/internal/domain"
)

// AccountRepository defines the interface for chart of accounts operations
type AccountRepository interface {
	// Create creates a new account
	Create(ctx context.Context, account *domain.Account) error

	// GetByID retrieves an account by ID
	GetByID(ctx context.Context, id int64) (*domain.Account, error)

	// GetByCode retrieves an account by code
	GetByCode(ctx context.Context, code string) (*domain.Account, error)

	// Update updates an existing account
	Update(ctx context.Context, account *domain.Account) error

	// List retrieves all accounts
	List(ctx context.Context) ([]*domain.Account, error)

	// ListByType retrieves accounts by type
	ListByType(ctx context.Context, accountType string) ([]*domain.Account, error)

	// ListChildren retrieves child accounts for a parent
	ListChildren(ctx context.Context, parentID int64) ([]*domain.Account, error)

	// GetTree retrieves the full account tree
	GetTree(ctx context.Context) ([]*domain.Account, error)
}

// AccountingEntryRepository defines the interface for journal entry operations
type AccountingEntryRepository interface {
	// Create creates a new accounting entry with its lines
	Create(ctx context.Context, entry *domain.AccountingEntry) error

	// GetByID retrieves an entry by ID with its lines
	GetByID(ctx context.Context, id int64) (*domain.AccountingEntry, error)

	// GetByNumber retrieves an entry by entry number
	GetByNumber(ctx context.Context, number string) (*domain.AccountingEntry, error)

	// Update updates an existing entry
	Update(ctx context.Context, entry *domain.AccountingEntry) error

	// List retrieves entries with filtering
	List(ctx context.Context, filter AccountingEntryFilter) ([]*domain.AccountingEntry, int64, error)

	// ListByBranch retrieves entries for a branch
	ListByBranch(ctx context.Context, branchID int64, filter AccountingEntryFilter) ([]*domain.AccountingEntry, int64, error)

	// ListByReference retrieves entries by reference
	ListByReference(ctx context.Context, refType string, refID int64) ([]*domain.AccountingEntry, error)

	// Post posts an entry (marks it as finalized)
	Post(ctx context.Context, id int64, postedBy int64) error

	// GenerateEntryNumber generates a unique entry number
	GenerateEntryNumber(ctx context.Context) (string, error)

	// GetAccountBalance retrieves the balance for an account
	GetAccountBalance(ctx context.Context, accountID int64, asOfDate time.Time) (float64, error)

	// GetAccountBalanceByBranch retrieves the balance for an account in a specific branch
	GetAccountBalanceByBranch(ctx context.Context, accountID int64, branchID int64, asOfDate time.Time) (float64, error)
}

// AccountingEntryFilter contains filters for listing accounting entries
type AccountingEntryFilter struct {
	BranchID      *int64
	AccountID     *int64
	ReferenceType *string
	ReferenceID   *int64
	IsPosted      *bool
	DateFrom      *string
	DateTo        *string
	Page          int
	PageSize      int
}

// DailyBalanceRepository defines the interface for daily balance operations
type DailyBalanceRepository interface {
	// Create creates or updates a daily balance
	Create(ctx context.Context, balance *domain.DailyBalance) error

	// GetByID retrieves a daily balance by ID
	GetByID(ctx context.Context, id int64) (*domain.DailyBalance, error)

	// GetByBranchAndDate retrieves a daily balance for a branch on a specific date
	GetByBranchAndDate(ctx context.Context, branchID int64, date time.Time) (*domain.DailyBalance, error)

	// Update updates an existing daily balance
	Update(ctx context.Context, balance *domain.DailyBalance) error

	// Upsert creates or updates a daily balance
	Upsert(ctx context.Context, balance *domain.DailyBalance) error

	// ListByBranch retrieves daily balances for a branch
	ListByBranch(ctx context.Context, branchID int64, dateFrom, dateTo time.Time) ([]*domain.DailyBalance, error)

	// GetSummary retrieves aggregated balances for a period
	GetSummary(ctx context.Context, branchID *int64, dateFrom, dateTo time.Time) (*DailyBalanceSummary, error)
}

// DailyBalanceSummary represents aggregated daily balance data
type DailyBalanceSummary struct {
	TotalLoanDisbursements  float64 `json:"total_loan_disbursements"`
	TotalInterestIncome     float64 `json:"total_interest_income"`
	TotalLateFeeIncome      float64 `json:"total_late_fee_income"`
	TotalSalesIncome        float64 `json:"total_sales_income"`
	TotalOtherIncome        float64 `json:"total_other_income"`
	TotalOperationalExpenses float64 `json:"total_operational_expenses"`
	TotalRefunds            float64 `json:"total_refunds"`
	TotalOtherExpenses      float64 `json:"total_other_expenses"`
	TotalNetIncome          float64 `json:"total_net_income"`
}

// ExpenseCategoryRepository defines the interface for expense category operations
type ExpenseCategoryRepository interface {
	// Create creates a new expense category
	Create(ctx context.Context, category *domain.ExpenseCategory) error

	// GetByID retrieves an expense category by ID
	GetByID(ctx context.Context, id int64) (*domain.ExpenseCategory, error)

	// GetByCode retrieves an expense category by code
	GetByCode(ctx context.Context, code string) (*domain.ExpenseCategory, error)

	// Update updates an existing expense category
	Update(ctx context.Context, category *domain.ExpenseCategory) error

	// List retrieves all expense categories
	List(ctx context.Context, includeInactive bool) ([]*domain.ExpenseCategory, error)
}

// ExpenseRepository defines the interface for expense operations
type ExpenseRepository interface {
	// Create creates a new expense
	Create(ctx context.Context, expense *domain.Expense) error

	// GetByID retrieves an expense by ID
	GetByID(ctx context.Context, id int64) (*domain.Expense, error)

	// GetByNumber retrieves an expense by number
	GetByNumber(ctx context.Context, number string) (*domain.Expense, error)

	// Update updates an existing expense
	Update(ctx context.Context, expense *domain.Expense) error

	// Delete deletes an expense
	Delete(ctx context.Context, id int64) error

	// List retrieves expenses with filtering
	List(ctx context.Context, filter ExpenseFilter) ([]*domain.Expense, int64, error)

	// ListByBranch retrieves expenses for a branch
	ListByBranch(ctx context.Context, branchID int64, filter ExpenseFilter) ([]*domain.Expense, int64, error)

	// Approve approves an expense
	Approve(ctx context.Context, id int64, approvedBy int64) error

	// GenerateExpenseNumber generates a unique expense number
	GenerateExpenseNumber(ctx context.Context) (string, error)

	// GetTotalByBranchAndDate retrieves total expenses for a branch on a date
	GetTotalByBranchAndDate(ctx context.Context, branchID int64, date time.Time) (float64, error)
}

// ExpenseFilter contains filters for listing expenses
type ExpenseFilter struct {
	BranchID    *int64
	CategoryID  *int64
	IsApproved  *bool
	DateFrom    *string
	DateTo      *string
	MinAmount   *float64
	MaxAmount   *float64
	Page        int
	PageSize    int
}
