package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// Account Repository
type accountRepository struct {
	db *DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *DB) repository.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *domain.Account) error {
	query := `
		INSERT INTO accounts (code, name, account_type, parent_id, description, is_active, is_system)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		account.Code,
		account.Name,
		account.AccountType,
		account.ParentID,
		account.Description,
		account.IsActive,
		account.IsSystem,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
}

func (r *accountRepository) GetByID(ctx context.Context, id int64) (*domain.Account, error) {
	query := `
		SELECT id, code, name, account_type, parent_id, description, is_active, is_system, created_at, updated_at
		FROM accounts
		WHERE id = $1`

	account := &domain.Account{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.Code,
		&account.Name,
		&account.AccountType,
		&account.ParentID,
		&account.Description,
		&account.IsActive,
		&account.IsSystem,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *accountRepository) GetByCode(ctx context.Context, code string) (*domain.Account, error) {
	query := `
		SELECT id, code, name, account_type, parent_id, description, is_active, is_system, created_at, updated_at
		FROM accounts
		WHERE code = $1`

	account := &domain.Account{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&account.ID,
		&account.Code,
		&account.Name,
		&account.AccountType,
		&account.ParentID,
		&account.Description,
		&account.IsActive,
		&account.IsSystem,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *accountRepository) Update(ctx context.Context, account *domain.Account) error {
	query := `
		UPDATE accounts SET
			name = $2,
			description = $3,
			is_active = $4,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRowContext(ctx, query,
		account.ID,
		account.Name,
		account.Description,
		account.IsActive,
	).Scan(&account.UpdatedAt)
}

func (r *accountRepository) List(ctx context.Context) ([]*domain.Account, error) {
	query := `
		SELECT id, code, name, account_type, parent_id, description, is_active, is_system, created_at, updated_at
		FROM accounts
		ORDER BY code`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		account := &domain.Account{}
		if err := rows.Scan(
			&account.ID,
			&account.Code,
			&account.Name,
			&account.AccountType,
			&account.ParentID,
			&account.Description,
			&account.IsActive,
			&account.IsSystem,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

func (r *accountRepository) ListByType(ctx context.Context, accountType string) ([]*domain.Account, error) {
	query := `
		SELECT id, code, name, account_type, parent_id, description, is_active, is_system, created_at, updated_at
		FROM accounts
		WHERE account_type = $1
		ORDER BY code`

	rows, err := r.db.QueryContext(ctx, query, accountType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		account := &domain.Account{}
		if err := rows.Scan(
			&account.ID,
			&account.Code,
			&account.Name,
			&account.AccountType,
			&account.ParentID,
			&account.Description,
			&account.IsActive,
			&account.IsSystem,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

func (r *accountRepository) ListChildren(ctx context.Context, parentID int64) ([]*domain.Account, error) {
	query := `
		SELECT id, code, name, account_type, parent_id, description, is_active, is_system, created_at, updated_at
		FROM accounts
		WHERE parent_id = $1
		ORDER BY code`

	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		account := &domain.Account{}
		if err := rows.Scan(
			&account.ID,
			&account.Code,
			&account.Name,
			&account.AccountType,
			&account.ParentID,
			&account.Description,
			&account.IsActive,
			&account.IsSystem,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

func (r *accountRepository) GetTree(ctx context.Context) ([]*domain.Account, error) {
	// First get all accounts
	accounts, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	// Build map by ID
	accountMap := make(map[int64]*domain.Account)
	for _, acc := range accounts {
		accountMap[acc.ID] = acc
	}

	// Build tree
	var rootAccounts []*domain.Account
	for _, acc := range accounts {
		if acc.ParentID == nil {
			rootAccounts = append(rootAccounts, acc)
		} else {
			parent := accountMap[*acc.ParentID]
			if parent != nil {
				parent.Children = append(parent.Children, acc)
			}
		}
	}

	return rootAccounts, nil
}

// Accounting Entry Repository
type accountingEntryRepository struct {
	db *DB
}

// NewAccountingEntryRepository creates a new accounting entry repository
func NewAccountingEntryRepository(db *DB) repository.AccountingEntryRepository {
	return &accountingEntryRepository{db: db}
}

func (r *accountingEntryRepository) Create(ctx context.Context, entry *domain.AccountingEntry) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert entry
	query := `
		INSERT INTO accounting_entries (
			entry_number, branch_id, entry_date, description,
			reference_type, reference_id, total_debit, total_credit,
			is_posted, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query,
		entry.EntryNumber,
		entry.BranchID,
		entry.EntryDate,
		entry.Description,
		entry.ReferenceType,
		entry.ReferenceID,
		entry.TotalDebit,
		entry.TotalCredit,
		entry.IsPosted,
		entry.CreatedBy,
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert lines
	for _, line := range entry.Lines {
		line.EntryID = entry.ID
		lineQuery := `
			INSERT INTO accounting_entry_lines (entry_id, account_id, entry_type, amount, description)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at`

		err = tx.QueryRowContext(ctx, lineQuery,
			line.EntryID,
			line.AccountID,
			line.EntryType,
			line.Amount,
			line.Description,
		).Scan(&line.ID, &line.CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *accountingEntryRepository) GetByID(ctx context.Context, id int64) (*domain.AccountingEntry, error) {
	query := `
		SELECT id, entry_number, branch_id, entry_date, description,
			   reference_type, reference_id, total_debit, total_credit,
			   is_posted, posted_at, posted_by, created_by, created_at, updated_at
		FROM accounting_entries
		WHERE id = $1`

	entry := &domain.AccountingEntry{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entry.ID,
		&entry.EntryNumber,
		&entry.BranchID,
		&entry.EntryDate,
		&entry.Description,
		&entry.ReferenceType,
		&entry.ReferenceID,
		&entry.TotalDebit,
		&entry.TotalCredit,
		&entry.IsPosted,
		&entry.PostedAt,
		&entry.PostedBy,
		&entry.CreatedBy,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Load lines
	lines, err := r.getEntryLines(ctx, id)
	if err != nil {
		return nil, err
	}
	entry.Lines = lines

	return entry, nil
}

func (r *accountingEntryRepository) getEntryLines(ctx context.Context, entryID int64) ([]*domain.AccountingEntryLine, error) {
	query := `
		SELECT id, entry_id, account_id, entry_type, amount, description, created_at
		FROM accounting_entry_lines
		WHERE entry_id = $1
		ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, entryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []*domain.AccountingEntryLine
	for rows.Next() {
		line := &domain.AccountingEntryLine{}
		if err := rows.Scan(
			&line.ID,
			&line.EntryID,
			&line.AccountID,
			&line.EntryType,
			&line.Amount,
			&line.Description,
			&line.CreatedAt,
		); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, rows.Err()
}

func (r *accountingEntryRepository) GetByNumber(ctx context.Context, number string) (*domain.AccountingEntry, error) {
	query := `SELECT id FROM accounting_entries WHERE entry_number = $1`
	var id int64
	err := r.db.QueryRowContext(ctx, query, number).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *accountingEntryRepository) Update(ctx context.Context, entry *domain.AccountingEntry) error {
	query := `
		UPDATE accounting_entries SET
			description = $2,
			reference_type = $3,
			reference_id = $4,
			updated_at = NOW()
		WHERE id = $1 AND is_posted = false
		RETURNING updated_at`

	return r.db.QueryRowContext(ctx, query,
		entry.ID,
		entry.Description,
		entry.ReferenceType,
		entry.ReferenceID,
	).Scan(&entry.UpdatedAt)
}

func (r *accountingEntryRepository) List(ctx context.Context, filter repository.AccountingEntryFilter) ([]*domain.AccountingEntry, int64, error) {
	var conditions []string
	var args []interface{}
	argPos := 1

	if filter.BranchID != nil {
		conditions = append(conditions, fmt.Sprintf("ae.branch_id = $%d", argPos))
		args = append(args, *filter.BranchID)
		argPos++
	}
	if filter.ReferenceType != nil {
		conditions = append(conditions, fmt.Sprintf("ae.reference_type = $%d", argPos))
		args = append(args, *filter.ReferenceType)
		argPos++
	}
	if filter.ReferenceID != nil {
		conditions = append(conditions, fmt.Sprintf("ae.reference_id = $%d", argPos))
		args = append(args, *filter.ReferenceID)
		argPos++
	}
	if filter.IsPosted != nil {
		conditions = append(conditions, fmt.Sprintf("ae.is_posted = $%d", argPos))
		args = append(args, *filter.IsPosted)
		argPos++
	}
	if filter.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("ae.entry_date >= $%d", argPos))
		args = append(args, *filter.DateFrom)
		argPos++
	}
	if filter.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("ae.entry_date <= $%d", argPos))
		args = append(args, *filter.DateTo)
		argPos++
	}
	if filter.AccountID != nil {
		conditions = append(conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM accounting_entry_lines l WHERE l.entry_id = ae.id AND l.account_id = $%d)", argPos))
		args = append(args, *filter.AccountID)
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM accounting_entries ae %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT ae.id, ae.entry_number, ae.branch_id, ae.entry_date, ae.description,
			   ae.reference_type, ae.reference_id, ae.total_debit, ae.total_credit,
			   ae.is_posted, ae.posted_at, ae.posted_by, ae.created_by, ae.created_at, ae.updated_at
		FROM accounting_entries ae
		%s
		ORDER BY ae.entry_date DESC, ae.id DESC
		LIMIT $%d OFFSET $%d`, whereClause, argPos, argPos+1)

	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []*domain.AccountingEntry
	for rows.Next() {
		entry := &domain.AccountingEntry{}
		if err := rows.Scan(
			&entry.ID,
			&entry.EntryNumber,
			&entry.BranchID,
			&entry.EntryDate,
			&entry.Description,
			&entry.ReferenceType,
			&entry.ReferenceID,
			&entry.TotalDebit,
			&entry.TotalCredit,
			&entry.IsPosted,
			&entry.PostedAt,
			&entry.PostedBy,
			&entry.CreatedBy,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		entries = append(entries, entry)
	}

	return entries, total, rows.Err()
}

func (r *accountingEntryRepository) ListByBranch(ctx context.Context, branchID int64, filter repository.AccountingEntryFilter) ([]*domain.AccountingEntry, int64, error) {
	filter.BranchID = &branchID
	return r.List(ctx, filter)
}

func (r *accountingEntryRepository) ListByReference(ctx context.Context, refType string, refID int64) ([]*domain.AccountingEntry, error) {
	query := `
		SELECT id, entry_number, branch_id, entry_date, description,
			   reference_type, reference_id, total_debit, total_credit,
			   is_posted, posted_at, posted_by, created_by, created_at, updated_at
		FROM accounting_entries
		WHERE reference_type = $1 AND reference_id = $2
		ORDER BY entry_date DESC`

	rows, err := r.db.QueryContext(ctx, query, refType, refID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*domain.AccountingEntry
	for rows.Next() {
		entry := &domain.AccountingEntry{}
		if err := rows.Scan(
			&entry.ID,
			&entry.EntryNumber,
			&entry.BranchID,
			&entry.EntryDate,
			&entry.Description,
			&entry.ReferenceType,
			&entry.ReferenceID,
			&entry.TotalDebit,
			&entry.TotalCredit,
			&entry.IsPosted,
			&entry.PostedAt,
			&entry.PostedBy,
			&entry.CreatedBy,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

func (r *accountingEntryRepository) Post(ctx context.Context, id int64, postedBy int64) error {
	query := `
		UPDATE accounting_entries SET
			is_posted = true,
			posted_at = NOW(),
			posted_by = $2,
			updated_at = NOW()
		WHERE id = $1 AND is_posted = false`

	result, err := r.db.ExecContext(ctx, query, id, postedBy)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("entry already posted or not found")
	}
	return nil
}

func (r *accountingEntryRepository) GenerateEntryNumber(ctx context.Context) (string, error) {
	now := time.Now()
	prefix := fmt.Sprintf("JE-%s-", now.Format("20060102"))

	query := `
		SELECT COUNT(*) + 1 FROM accounting_entries
		WHERE entry_number LIKE $1`

	var seq int
	if err := r.db.QueryRowContext(ctx, query, prefix+"%").Scan(&seq); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%04d", prefix, seq), nil
}

func (r *accountingEntryRepository) GetAccountBalance(ctx context.Context, accountID int64, asOfDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(
			SUM(CASE WHEN l.entry_type = 'debit' THEN l.amount ELSE -l.amount END),
			0
		)
		FROM accounting_entry_lines l
		JOIN accounting_entries e ON l.entry_id = e.id
		WHERE l.account_id = $1 AND e.entry_date <= $2 AND e.is_posted = true`

	var balance float64
	if err := r.db.QueryRowContext(ctx, query, accountID, asOfDate).Scan(&balance); err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *accountingEntryRepository) GetAccountBalanceByBranch(ctx context.Context, accountID int64, branchID int64, asOfDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(
			SUM(CASE WHEN l.entry_type = 'debit' THEN l.amount ELSE -l.amount END),
			0
		)
		FROM accounting_entry_lines l
		JOIN accounting_entries e ON l.entry_id = e.id
		WHERE l.account_id = $1 AND e.branch_id = $2 AND e.entry_date <= $3 AND e.is_posted = true`

	var balance float64
	if err := r.db.QueryRowContext(ctx, query, accountID, branchID, asOfDate).Scan(&balance); err != nil {
		return 0, err
	}
	return balance, nil
}

// Daily Balance Repository
type dailyBalanceRepository struct {
	db *DB
}

// NewDailyBalanceRepository creates a new daily balance repository
func NewDailyBalanceRepository(db *DB) repository.DailyBalanceRepository {
	return &dailyBalanceRepository{db: db}
}

func (r *dailyBalanceRepository) Create(ctx context.Context, balance *domain.DailyBalance) error {
	query := `
		INSERT INTO daily_balances (
			branch_id, balance_date, loan_disbursements, interest_income,
			late_fee_income, sales_income, other_income, operational_expenses,
			refunds, other_expenses, cash_opening, cash_closing,
			total_loans_active, total_loans_count, net_income
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		balance.BranchID,
		balance.BalanceDate,
		balance.LoanDisbursements,
		balance.InterestIncome,
		balance.LateFeeIncome,
		balance.SalesIncome,
		balance.OtherIncome,
		balance.OperationalExpenses,
		balance.Refunds,
		balance.OtherExpenses,
		balance.CashOpening,
		balance.CashClosing,
		balance.TotalLoansActive,
		balance.TotalLoansCount,
		balance.NetIncome,
	).Scan(&balance.ID, &balance.CreatedAt, &balance.UpdatedAt)
}

func (r *dailyBalanceRepository) GetByID(ctx context.Context, id int64) (*domain.DailyBalance, error) {
	query := `
		SELECT id, branch_id, balance_date, loan_disbursements, interest_income,
			   late_fee_income, sales_income, other_income, operational_expenses,
			   refunds, other_expenses, cash_opening, cash_closing,
			   total_loans_active, total_loans_count, net_income, created_at, updated_at
		FROM daily_balances
		WHERE id = $1`

	balance := &domain.DailyBalance{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&balance.ID,
		&balance.BranchID,
		&balance.BalanceDate,
		&balance.LoanDisbursements,
		&balance.InterestIncome,
		&balance.LateFeeIncome,
		&balance.SalesIncome,
		&balance.OtherIncome,
		&balance.OperationalExpenses,
		&balance.Refunds,
		&balance.OtherExpenses,
		&balance.CashOpening,
		&balance.CashClosing,
		&balance.TotalLoansActive,
		&balance.TotalLoansCount,
		&balance.NetIncome,
		&balance.CreatedAt,
		&balance.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (r *dailyBalanceRepository) GetByBranchAndDate(ctx context.Context, branchID int64, date time.Time) (*domain.DailyBalance, error) {
	query := `
		SELECT id, branch_id, balance_date, loan_disbursements, interest_income,
			   late_fee_income, sales_income, other_income, operational_expenses,
			   refunds, other_expenses, cash_opening, cash_closing,
			   total_loans_active, total_loans_count, net_income, created_at, updated_at
		FROM daily_balances
		WHERE branch_id = $1 AND DATE(balance_date) = DATE($2)`

	balance := &domain.DailyBalance{}
	err := r.db.QueryRowContext(ctx, query, branchID, date).Scan(
		&balance.ID,
		&balance.BranchID,
		&balance.BalanceDate,
		&balance.LoanDisbursements,
		&balance.InterestIncome,
		&balance.LateFeeIncome,
		&balance.SalesIncome,
		&balance.OtherIncome,
		&balance.OperationalExpenses,
		&balance.Refunds,
		&balance.OtherExpenses,
		&balance.CashOpening,
		&balance.CashClosing,
		&balance.TotalLoansActive,
		&balance.TotalLoansCount,
		&balance.NetIncome,
		&balance.CreatedAt,
		&balance.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (r *dailyBalanceRepository) Update(ctx context.Context, balance *domain.DailyBalance) error {
	query := `
		UPDATE daily_balances SET
			loan_disbursements = $2,
			interest_income = $3,
			late_fee_income = $4,
			sales_income = $5,
			other_income = $6,
			operational_expenses = $7,
			refunds = $8,
			other_expenses = $9,
			cash_opening = $10,
			cash_closing = $11,
			total_loans_active = $12,
			total_loans_count = $13,
			net_income = $14,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRowContext(ctx, query,
		balance.ID,
		balance.LoanDisbursements,
		balance.InterestIncome,
		balance.LateFeeIncome,
		balance.SalesIncome,
		balance.OtherIncome,
		balance.OperationalExpenses,
		balance.Refunds,
		balance.OtherExpenses,
		balance.CashOpening,
		balance.CashClosing,
		balance.TotalLoansActive,
		balance.TotalLoansCount,
		balance.NetIncome,
	).Scan(&balance.UpdatedAt)
}

func (r *dailyBalanceRepository) Upsert(ctx context.Context, balance *domain.DailyBalance) error {
	query := `
		INSERT INTO daily_balances (
			branch_id, balance_date, loan_disbursements, interest_income,
			late_fee_income, sales_income, other_income, operational_expenses,
			refunds, other_expenses, cash_opening, cash_closing,
			total_loans_active, total_loans_count, net_income
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (branch_id, balance_date) DO UPDATE SET
			loan_disbursements = EXCLUDED.loan_disbursements,
			interest_income = EXCLUDED.interest_income,
			late_fee_income = EXCLUDED.late_fee_income,
			sales_income = EXCLUDED.sales_income,
			other_income = EXCLUDED.other_income,
			operational_expenses = EXCLUDED.operational_expenses,
			refunds = EXCLUDED.refunds,
			other_expenses = EXCLUDED.other_expenses,
			cash_opening = EXCLUDED.cash_opening,
			cash_closing = EXCLUDED.cash_closing,
			total_loans_active = EXCLUDED.total_loans_active,
			total_loans_count = EXCLUDED.total_loans_count,
			net_income = EXCLUDED.net_income,
			updated_at = NOW()
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		balance.BranchID,
		balance.BalanceDate,
		balance.LoanDisbursements,
		balance.InterestIncome,
		balance.LateFeeIncome,
		balance.SalesIncome,
		balance.OtherIncome,
		balance.OperationalExpenses,
		balance.Refunds,
		balance.OtherExpenses,
		balance.CashOpening,
		balance.CashClosing,
		balance.TotalLoansActive,
		balance.TotalLoansCount,
		balance.NetIncome,
	).Scan(&balance.ID, &balance.CreatedAt, &balance.UpdatedAt)
}

func (r *dailyBalanceRepository) ListByBranch(ctx context.Context, branchID int64, dateFrom, dateTo time.Time) ([]*domain.DailyBalance, error) {
	query := `
		SELECT id, branch_id, balance_date, loan_disbursements, interest_income,
			   late_fee_income, sales_income, other_income, operational_expenses,
			   refunds, other_expenses, cash_opening, cash_closing,
			   total_loans_active, total_loans_count, net_income, created_at, updated_at
		FROM daily_balances
		WHERE branch_id = $1 AND balance_date >= $2 AND balance_date <= $3
		ORDER BY balance_date DESC`

	rows, err := r.db.QueryContext(ctx, query, branchID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []*domain.DailyBalance
	for rows.Next() {
		balance := &domain.DailyBalance{}
		if err := rows.Scan(
			&balance.ID,
			&balance.BranchID,
			&balance.BalanceDate,
			&balance.LoanDisbursements,
			&balance.InterestIncome,
			&balance.LateFeeIncome,
			&balance.SalesIncome,
			&balance.OtherIncome,
			&balance.OperationalExpenses,
			&balance.Refunds,
			&balance.OtherExpenses,
			&balance.CashOpening,
			&balance.CashClosing,
			&balance.TotalLoansActive,
			&balance.TotalLoansCount,
			&balance.NetIncome,
			&balance.CreatedAt,
			&balance.UpdatedAt,
		); err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}

	return balances, rows.Err()
}

func (r *dailyBalanceRepository) GetSummary(ctx context.Context, branchID *int64, dateFrom, dateTo time.Time) (*repository.DailyBalanceSummary, error) {
	var query string
	var args []interface{}

	if branchID != nil {
		query = `
			SELECT
				COALESCE(SUM(loan_disbursements), 0),
				COALESCE(SUM(interest_income), 0),
				COALESCE(SUM(late_fee_income), 0),
				COALESCE(SUM(sales_income), 0),
				COALESCE(SUM(other_income), 0),
				COALESCE(SUM(operational_expenses), 0),
				COALESCE(SUM(refunds), 0),
				COALESCE(SUM(other_expenses), 0),
				COALESCE(SUM(net_income), 0)
			FROM daily_balances
			WHERE branch_id = $1 AND balance_date >= $2 AND balance_date <= $3`
		args = []interface{}{*branchID, dateFrom, dateTo}
	} else {
		query = `
			SELECT
				COALESCE(SUM(loan_disbursements), 0),
				COALESCE(SUM(interest_income), 0),
				COALESCE(SUM(late_fee_income), 0),
				COALESCE(SUM(sales_income), 0),
				COALESCE(SUM(other_income), 0),
				COALESCE(SUM(operational_expenses), 0),
				COALESCE(SUM(refunds), 0),
				COALESCE(SUM(other_expenses), 0),
				COALESCE(SUM(net_income), 0)
			FROM daily_balances
			WHERE balance_date >= $1 AND balance_date <= $2`
		args = []interface{}{dateFrom, dateTo}
	}

	summary := &repository.DailyBalanceSummary{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&summary.TotalLoanDisbursements,
		&summary.TotalInterestIncome,
		&summary.TotalLateFeeIncome,
		&summary.TotalSalesIncome,
		&summary.TotalOtherIncome,
		&summary.TotalOperationalExpenses,
		&summary.TotalRefunds,
		&summary.TotalOtherExpenses,
		&summary.TotalNetIncome,
	)
	if err != nil {
		return nil, err
	}
	return summary, nil
}
