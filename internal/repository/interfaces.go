package repository

import (
	"context"
	"pawnshop/internal/domain"
)

// Common pagination parameters
type PaginationParams struct {
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
	OrderBy string `query:"order_by"`
	Order   string `query:"order"` // asc or desc
}

// PaginatedResult contains paginated results with metadata
type PaginatedResult[T any] struct {
	Data       []T `json:"data"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
}

// Transaction interface for database transactions
type Transaction interface {
	Commit() error
	Rollback() error
}

// BranchRepository defines methods for branch operations
type BranchRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Branch, error)
	GetByCode(ctx context.Context, code string) (*domain.Branch, error)
	List(ctx context.Context, params PaginationParams) (*PaginatedResult[domain.Branch], error)
	Create(ctx context.Context, branch *domain.Branch) error
	Update(ctx context.Context, branch *domain.Branch) error
	Delete(ctx context.Context, id int64) error
}

// RoleRepository defines methods for role operations
type RoleRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Role, error)
	GetByName(ctx context.Context, name string) (*domain.Role, error)
	List(ctx context.Context) ([]*domain.Role, error)
	Create(ctx context.Context, role *domain.Role) error
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, id int64) error
}

// UserRepository defines methods for user operations
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, params UserListParams) (*PaginatedResult[domain.User], error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int64) error
	UpdatePassword(ctx context.Context, id int64, passwordHash string) error
	UpdateLastLogin(ctx context.Context, id int64, ip string) error
	IncrementFailedLogins(ctx context.Context, id int64) error
	ResetFailedLogins(ctx context.Context, id int64) error
	LockUser(ctx context.Context, id int64, until *int64) error
}

// UserListParams for filtering user list
type UserListParams struct {
	PaginationParams
	BranchID *int64 `query:"branch_id"`
	RoleID   *int64 `query:"role_id"`
	IsActive *bool  `query:"is_active"`
	Search   string `query:"search"`
}

// CustomerRepository defines methods for customer operations
type CustomerRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Customer, error)
	GetByIdentity(ctx context.Context, branchID int64, identityType, identityNumber string) (*domain.Customer, error)
	List(ctx context.Context, params CustomerListParams) (*PaginatedResult[domain.Customer], error)
	Create(ctx context.Context, customer *domain.Customer) error
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, id int64) error
	UpdateCreditInfo(ctx context.Context, id int64, info CustomerCreditUpdate) error
}

// CustomerListParams for filtering customer list
type CustomerListParams struct {
	PaginationParams
	BranchID  int64  `query:"branch_id"`
	IsActive  *bool  `query:"is_active"`
	IsBlocked *bool  `query:"is_blocked"`
	Search    string `query:"search"`
}

// CustomerCreditUpdate for updating credit info
type CustomerCreditUpdate struct {
	CreditLimit    *float64
	CreditScore    *int
	TotalLoans     *int
	TotalPaid      *float64
	TotalDefaulted *float64
}

// CategoryRepository defines methods for category operations
type CategoryRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Category, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Category, error)
	List(ctx context.Context, params CategoryListParams) ([]*domain.Category, error)
	ListWithChildren(ctx context.Context) ([]*domain.Category, error)
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int64) error
}

// CategoryListParams for filtering category list
type CategoryListParams struct {
	ParentID *int64 `query:"parent_id"`
	IsActive *bool  `query:"is_active"`
}

// ItemRepository defines methods for item operations
type ItemRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Item, error)
	GetBySKU(ctx context.Context, sku string) (*domain.Item, error)
	List(ctx context.Context, params ItemListParams) (*PaginatedResult[domain.Item], error)
	Create(ctx context.Context, item *domain.Item) error
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, id int64) error
	UpdateStatus(ctx context.Context, id int64, status domain.ItemStatus) error
	GenerateSKU(ctx context.Context, branchID int64) (string, error)
	CreateHistory(ctx context.Context, history *domain.ItemHistory) error
}

// ItemListParams for filtering item list
type ItemListParams struct {
	PaginationParams
	BranchID   int64               `query:"branch_id"`
	CategoryID *int64              `query:"category_id"`
	CustomerID *int64              `query:"customer_id"`
	Status     *domain.ItemStatus  `query:"status"`
	Search     string              `query:"search"`
}

// LoanRepository defines methods for loan operations
type LoanRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Loan, error)
	GetByNumber(ctx context.Context, loanNumber string) (*domain.Loan, error)
	List(ctx context.Context, params LoanListParams) (*PaginatedResult[domain.Loan], error)
	Create(ctx context.Context, loan *domain.Loan) error
	Update(ctx context.Context, loan *domain.Loan) error
	GenerateNumber(ctx context.Context) (string, error)
	GetOverdueLoans(ctx context.Context, branchID int64) ([]*domain.Loan, error)
	UpdateStatus(ctx context.Context, id int64, status domain.LoanStatus) error
	BeginTx(ctx context.Context) (Transaction, error)
	CreateTx(ctx context.Context, tx Transaction, loan *domain.Loan) error

	// Installments
	CreateInstallments(ctx context.Context, installments []*domain.LoanInstallment) error
	CreateInstallmentsTx(ctx context.Context, tx Transaction, installments []*domain.LoanInstallment) error
	GetInstallments(ctx context.Context, loanID int64) ([]*domain.LoanInstallment, error)
	UpdateInstallment(ctx context.Context, installment *domain.LoanInstallment) error
}

// LoanListParams for filtering loan list
type LoanListParams struct {
	PaginationParams
	BranchID   int64              `query:"branch_id"`
	CustomerID *int64             `query:"customer_id"`
	ItemID     *int64             `query:"item_id"`
	Status     *domain.LoanStatus `query:"status"`
	DueBefore  *string            `query:"due_before"`
	DueAfter   *string            `query:"due_after"`
	Search     string             `query:"search"`
}

// PaymentRepository defines methods for payment operations
type PaymentRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Payment, error)
	GetByNumber(ctx context.Context, paymentNumber string) (*domain.Payment, error)
	List(ctx context.Context, params PaymentListParams) (*PaginatedResult[domain.Payment], error)
	ListByLoan(ctx context.Context, loanID int64) ([]*domain.Payment, error)
	Create(ctx context.Context, payment *domain.Payment) error
	Update(ctx context.Context, payment *domain.Payment) error
	GenerateNumber(ctx context.Context) (string, error)
}

// PaymentListParams for filtering payment list
type PaymentListParams struct {
	PaginationParams
	BranchID   int64                  `query:"branch_id"`
	CustomerID *int64                 `query:"customer_id"`
	LoanID     *int64                 `query:"loan_id"`
	Status     *domain.PaymentStatus  `query:"status"`
	Method     *domain.PaymentMethod  `query:"method"`
	DateFrom   *string                `query:"date_from"`
	DateTo     *string                `query:"date_to"`
}

// SaleRepository defines methods for sale operations
type SaleRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Sale, error)
	GetByNumber(ctx context.Context, saleNumber string) (*domain.Sale, error)
	List(ctx context.Context, params SaleListParams) (*PaginatedResult[domain.Sale], error)
	Create(ctx context.Context, sale *domain.Sale) error
	Update(ctx context.Context, sale *domain.Sale) error
	GenerateNumber(ctx context.Context) (string, error)
}

// SaleListParams for filtering sale list
type SaleListParams struct {
	PaginationParams
	BranchID   int64              `query:"branch_id"`
	CustomerID *int64             `query:"customer_id"`
	ItemID     *int64             `query:"item_id"`
	Status     *domain.SaleStatus `query:"status"`
	DateFrom   *string            `query:"date_from"`
	DateTo     *string            `query:"date_to"`
}

// CashRegisterRepository defines methods for cash register operations
type CashRegisterRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.CashRegister, error)
	List(ctx context.Context, branchID int64) ([]*domain.CashRegister, error)
	Create(ctx context.Context, register *domain.CashRegister) error
	Update(ctx context.Context, register *domain.CashRegister) error
}

// CashSessionRepository defines methods for cash session operations
type CashSessionRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.CashSession, error)
	GetOpenSession(ctx context.Context, userID int64) (*domain.CashSession, error)
	GetOpenSessionByRegister(ctx context.Context, registerID int64) (*domain.CashSession, error)
	List(ctx context.Context, params CashSessionListParams) (*PaginatedResult[domain.CashSession], error)
	Create(ctx context.Context, session *domain.CashSession) error
	Update(ctx context.Context, session *domain.CashSession) error
	Close(ctx context.Context, id int64, closingData CashSessionCloseData) error
}

// CashSessionListParams for filtering cash session list
type CashSessionListParams struct {
	PaginationParams
	BranchID       int64                      `query:"branch_id"`
	UserID         *int64                     `query:"user_id"`
	RegisterID     *int64                     `query:"register_id"`
	Status         *domain.CashSessionStatus  `query:"status"`
	DateFrom       *string                    `query:"date_from"`
	DateTo         *string                    `query:"date_to"`
}

// CashSessionCloseData contains closing data for a session
type CashSessionCloseData struct {
	ClosingAmount  float64
	ExpectedAmount float64
	Difference     float64
	ClosedBy       int64
	ClosingNotes   string
}

// CashMovementRepository defines methods for cash movement operations
type CashMovementRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.CashMovement, error)
	List(ctx context.Context, params CashMovementListParams) (*PaginatedResult[domain.CashMovement], error)
	ListBySession(ctx context.Context, sessionID int64) ([]*domain.CashMovement, error)
	Create(ctx context.Context, movement *domain.CashMovement) error
}

// CashMovementListParams for filtering cash movement list
type CashMovementListParams struct {
	PaginationParams
	BranchID      int64                     `query:"branch_id"`
	SessionID     *int64                    `query:"session_id"`
	MovementType  *domain.CashMovementType  `query:"type"`
	PaymentMethod *domain.PaymentMethod     `query:"method"`
	DateFrom      *string                   `query:"date_from"`
	DateTo        *string                   `query:"date_to"`
}

// RefreshTokenRepository defines methods for refresh token operations
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *domain.RefreshToken) error
	GetByHash(ctx context.Context, hash string) (*domain.RefreshToken, error)
	Revoke(ctx context.Context, id int64) error
	RevokeAllForUser(ctx context.Context, userID int64) error
	DeleteExpired(ctx context.Context) error
}

// SettingRepository defines methods for settings operations
type SettingRepository interface {
	Get(ctx context.Context, key string, branchID *int64) (*domain.Setting, error)
	GetAll(ctx context.Context, branchID *int64) ([]*domain.Setting, error)
	Set(ctx context.Context, setting *domain.Setting) error
	Delete(ctx context.Context, key string, branchID *int64) error
}

// AuditLogRepository defines methods for audit log operations
type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
	List(ctx context.Context, params AuditLogListParams) (*PaginatedResult[domain.AuditLog], error)
}

// AuditLogListParams for filtering audit log list
type AuditLogListParams struct {
	PaginationParams
	BranchID   *int64  `query:"branch_id"`
	UserID     *int64  `query:"user_id"`
	Action     string  `query:"action"`
	EntityType string  `query:"entity_type"`
	EntityID   *int64  `query:"entity_id"`
	DateFrom   *string `query:"date_from"`
	DateTo     *string `query:"date_to"`
}

// DocumentRepository defines methods for document operations
type DocumentRepository interface {
	Create(ctx context.Context, doc *domain.Document) error
	GetByID(ctx context.Context, id int64) (*domain.Document, error)
	ListByReference(ctx context.Context, refType string, refID int64) ([]*domain.Document, error)
}
