-- Create cash register table
CREATE TABLE cash_registers (
    id              BIGSERIAL PRIMARY KEY,
    branch_id       BIGINT NOT NULL REFERENCES branches(id),
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(50) NOT NULL,
    description     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT true,

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cash_registers_branch_id ON cash_registers(branch_id);
CREATE UNIQUE INDEX idx_cash_registers_code ON cash_registers(branch_id, code);

CREATE TRIGGER cash_registers_updated_at
    BEFORE UPDATE ON cash_registers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create cash session status enum
CREATE TYPE cash_session_status AS ENUM (
    'open',
    'closed'
);

-- Create cash sessions table
CREATE TABLE cash_sessions (
    id                  BIGSERIAL PRIMARY KEY,
    cash_register_id    BIGINT NOT NULL REFERENCES cash_registers(id),
    user_id             BIGINT NOT NULL REFERENCES users(id),
    branch_id           BIGINT NOT NULL REFERENCES branches(id),

    -- Amounts
    opening_amount      DECIMAL(12,2) NOT NULL,
    closing_amount      DECIMAL(12,2),
    expected_amount     DECIMAL(12,2),
    difference          DECIMAL(12,2),

    -- Totals by type
    total_income        DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_expense       DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_loans_disbursed DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_payments_received DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_sales         DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Breakdown by payment method
    total_cash          DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_card          DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_transfer      DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_other         DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Status
    status              cash_session_status NOT NULL DEFAULT 'open',

    -- Timestamps
    opened_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at           TIMESTAMPTZ,

    -- Notes
    opening_notes       TEXT,
    closing_notes       TEXT,

    -- Audit
    closed_by           BIGINT REFERENCES users(id),

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cash_sessions_cash_register_id ON cash_sessions(cash_register_id);
CREATE INDEX idx_cash_sessions_user_id ON cash_sessions(user_id);
CREATE INDEX idx_cash_sessions_branch_id ON cash_sessions(branch_id);
CREATE INDEX idx_cash_sessions_status ON cash_sessions(status);
CREATE INDEX idx_cash_sessions_opened_at ON cash_sessions(opened_at);

-- Only one open session per cash register
CREATE UNIQUE INDEX idx_cash_sessions_open_register
ON cash_sessions(cash_register_id)
WHERE status = 'open';

-- Only one open session per user
CREATE UNIQUE INDEX idx_cash_sessions_open_user
ON cash_sessions(user_id)
WHERE status = 'open';

CREATE TRIGGER cash_sessions_updated_at
    BEFORE UPDATE ON cash_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create cash movement type enum
CREATE TYPE cash_movement_type AS ENUM (
    'income_loan_disbursement',  -- Egreso: Desembolso de préstamo
    'income_payment',            -- Ingreso: Cobro de pago
    'income_sale',               -- Ingreso: Venta
    'income_other',              -- Ingreso: Otros
    'expense_loan_disbursement', -- Egreso: Desembolso de préstamo
    'expense_return',            -- Egreso: Devolución
    'expense_supplier',          -- Egreso: Pago a proveedor
    'expense_other',             -- Egreso: Otros
    'adjustment_positive',       -- Ajuste positivo
    'adjustment_negative'        -- Ajuste negativo
);

-- Create cash movements table
CREATE TABLE cash_movements (
    id                  BIGSERIAL PRIMARY KEY,
    cash_session_id     BIGINT NOT NULL REFERENCES cash_sessions(id),
    branch_id           BIGINT NOT NULL REFERENCES branches(id),

    -- Movement details
    movement_type       cash_movement_type NOT NULL,
    amount              DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    payment_method      payment_method NOT NULL DEFAULT 'cash',

    -- Reference to related entity
    reference_type      VARCHAR(50), -- loan, payment, sale
    reference_id        BIGINT,

    -- Description
    description         TEXT NOT NULL,
    notes               TEXT,

    -- Balance
    balance_after       DECIMAL(12,2) NOT NULL,

    -- Audit
    created_by          BIGINT REFERENCES users(id),

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cash_movements_session_id ON cash_movements(cash_session_id);
CREATE INDEX idx_cash_movements_branch_id ON cash_movements(branch_id);
CREATE INDEX idx_cash_movements_type ON cash_movements(movement_type);
CREATE INDEX idx_cash_movements_reference ON cash_movements(reference_type, reference_id);
CREATE INDEX idx_cash_movements_created_at ON cash_movements(created_at);

-- Add foreign key from payments to cash_sessions
ALTER TABLE payments
ADD CONSTRAINT fk_payments_cash_session
FOREIGN KEY (cash_session_id) REFERENCES cash_sessions(id);

-- Add foreign key from sales to cash_sessions
ALTER TABLE sales
ADD CONSTRAINT fk_sales_cash_session
FOREIGN KEY (cash_session_id) REFERENCES cash_sessions(id);

-- Insert default cash register for main branch
INSERT INTO cash_registers (branch_id, name, code, description)
VALUES (1, 'Caja Principal', 'CAJA-01', 'Caja principal de la sucursal');
