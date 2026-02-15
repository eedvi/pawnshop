-- Accounting module for basic financial tracking

-- Account types
CREATE TYPE account_type AS ENUM ('asset', 'liability', 'equity', 'income', 'expense');

-- Chart of accounts
CREATE TABLE accounts (
    id              BIGSERIAL PRIMARY KEY,
    code            VARCHAR(20) NOT NULL UNIQUE,
    name            VARCHAR(100) NOT NULL,
    account_type    account_type NOT NULL,
    parent_id       BIGINT REFERENCES accounts(id),
    description     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    is_system       BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Entry types
CREATE TYPE entry_type AS ENUM ('debit', 'credit');

-- Accounting entries (journal entries)
CREATE TABLE accounting_entries (
    id              BIGSERIAL PRIMARY KEY,
    entry_number    VARCHAR(50) NOT NULL UNIQUE,
    branch_id       BIGINT NOT NULL REFERENCES branches(id),

    -- Entry details
    entry_date      DATE NOT NULL DEFAULT CURRENT_DATE,
    description     TEXT NOT NULL,

    -- Reference to source document
    reference_type  VARCHAR(50), -- loan, payment, sale, cash_movement, etc.
    reference_id    BIGINT,

    -- Totals (for quick access)
    total_debit     DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_credit    DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Status
    is_posted       BOOLEAN NOT NULL DEFAULT false,
    posted_at       TIMESTAMPTZ,
    posted_by       BIGINT REFERENCES users(id),

    -- Audit
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Accounting entry lines (individual debits/credits)
CREATE TABLE accounting_entry_lines (
    id              BIGSERIAL PRIMARY KEY,
    entry_id        BIGINT NOT NULL REFERENCES accounting_entries(id) ON DELETE CASCADE,
    account_id      BIGINT NOT NULL REFERENCES accounts(id),
    entry_type      entry_type NOT NULL,
    amount          DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Daily balance summaries per branch
CREATE TABLE daily_balances (
    id              BIGSERIAL PRIMARY KEY,
    branch_id       BIGINT NOT NULL REFERENCES branches(id),
    balance_date    DATE NOT NULL,

    -- Income
    loan_disbursements      DECIMAL(12,2) NOT NULL DEFAULT 0,
    interest_income         DECIMAL(12,2) NOT NULL DEFAULT 0,
    late_fee_income         DECIMAL(12,2) NOT NULL DEFAULT 0,
    sales_income            DECIMAL(12,2) NOT NULL DEFAULT 0,
    other_income            DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Expenses
    operational_expenses    DECIMAL(12,2) NOT NULL DEFAULT 0,
    refunds                 DECIMAL(12,2) NOT NULL DEFAULT 0,
    other_expenses          DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Cash position
    cash_opening            DECIMAL(12,2) NOT NULL DEFAULT 0,
    cash_closing            DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Loan portfolio
    total_loans_active      DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_loans_count       INTEGER NOT NULL DEFAULT 0,

    -- Calculated
    net_income              DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(branch_id, balance_date)
);

-- Expense categories
CREATE TABLE expense_categories (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(20) NOT NULL UNIQUE,
    description     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Operational expenses
CREATE TABLE expenses (
    id              BIGSERIAL PRIMARY KEY,
    expense_number  VARCHAR(50) NOT NULL UNIQUE,
    branch_id       BIGINT NOT NULL REFERENCES branches(id),
    category_id     BIGINT REFERENCES expense_categories(id),

    -- Details
    description     TEXT NOT NULL,
    amount          DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    expense_date    DATE NOT NULL DEFAULT CURRENT_DATE,

    -- Payment info
    payment_method  VARCHAR(20) NOT NULL DEFAULT 'cash',
    receipt_number  VARCHAR(100),
    vendor          VARCHAR(200),

    -- Approval
    approved_by     BIGINT REFERENCES users(id),
    approved_at     TIMESTAMPTZ,

    -- Audit
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_accounting_entries_branch ON accounting_entries(branch_id);
CREATE INDEX idx_accounting_entries_date ON accounting_entries(entry_date);
CREATE INDEX idx_accounting_entries_reference ON accounting_entries(reference_type, reference_id);
CREATE INDEX idx_accounting_entry_lines_entry ON accounting_entry_lines(entry_id);
CREATE INDEX idx_accounting_entry_lines_account ON accounting_entry_lines(account_id);
CREATE INDEX idx_daily_balances_branch_date ON daily_balances(branch_id, balance_date);
CREATE INDEX idx_expenses_branch ON expenses(branch_id);
CREATE INDEX idx_expenses_date ON expenses(expense_date);

-- Triggers
CREATE TRIGGER accounting_entries_updated_at
    BEFORE UPDATE ON accounting_entries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER daily_balances_updated_at
    BEFORE UPDATE ON daily_balances
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER expenses_updated_at
    BEFORE UPDATE ON expenses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default accounts (Chart of Accounts)
INSERT INTO accounts (code, name, account_type, is_system) VALUES
    ('1000', 'Activos', 'asset', true),
    ('1100', 'Caja y Bancos', 'asset', true),
    ('1110', 'Caja General', 'asset', true),
    ('1120', 'Bancos', 'asset', true),
    ('1200', 'Cuentas por Cobrar', 'asset', true),
    ('1210', 'Préstamos por Cobrar', 'asset', true),
    ('1220', 'Intereses por Cobrar', 'asset', true),
    ('1300', 'Inventario', 'asset', true),
    ('1310', 'Artículos en Prenda', 'asset', true),
    ('1320', 'Artículos para Venta', 'asset', true),
    ('2000', 'Pasivos', 'liability', true),
    ('2100', 'Cuentas por Pagar', 'liability', true),
    ('3000', 'Capital', 'equity', true),
    ('3100', 'Capital Social', 'equity', true),
    ('3200', 'Utilidades Retenidas', 'equity', true),
    ('4000', 'Ingresos', 'income', true),
    ('4100', 'Ingresos por Intereses', 'income', true),
    ('4200', 'Ingresos por Mora', 'income', true),
    ('4300', 'Ingresos por Ventas', 'income', true),
    ('4400', 'Otros Ingresos', 'income', true),
    ('5000', 'Gastos', 'expense', true),
    ('5100', 'Gastos Operativos', 'expense', true),
    ('5110', 'Salarios', 'expense', true),
    ('5120', 'Alquiler', 'expense', true),
    ('5130', 'Servicios', 'expense', true),
    ('5200', 'Gastos Administrativos', 'expense', true),
    ('5300', 'Otros Gastos', 'expense', true);

-- Insert default expense categories
INSERT INTO expense_categories (code, name) VALUES
    ('SAL', 'Salarios y Sueldos'),
    ('RNT', 'Alquiler'),
    ('UTL', 'Servicios (Agua, Luz, Tel)'),
    ('SUP', 'Suministros de Oficina'),
    ('MNT', 'Mantenimiento'),
    ('TRN', 'Transporte'),
    ('MKT', 'Publicidad y Marketing'),
    ('OTH', 'Otros Gastos');
