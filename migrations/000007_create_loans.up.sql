-- Create loan status enum
CREATE TYPE loan_status AS ENUM (
    'active',       -- Loan is active and running
    'paid',         -- Fully paid off
    'overdue',      -- Past due date but within grace period
    'defaulted',    -- Past grace period, pending confiscation
    'renewed',      -- Renewed/extended
    'confiscated'   -- Item confiscated due to non-payment
);

-- Create payment plan type enum
CREATE TYPE payment_plan_type AS ENUM (
    'single',           -- Single payment at end
    'minimum_payment',  -- Monthly minimum payment required
    'installments'      -- Fixed installments
);

-- Create loans table
CREATE TABLE loans (
    id                  BIGSERIAL PRIMARY KEY,
    loan_number         VARCHAR(50) NOT NULL UNIQUE,
    branch_id           BIGINT NOT NULL REFERENCES branches(id),
    customer_id         BIGINT NOT NULL REFERENCES customers(id),
    item_id             BIGINT NOT NULL REFERENCES items(id),

    -- Amounts
    loan_amount             DECIMAL(12,2) NOT NULL CHECK (loan_amount > 0),
    interest_rate           DECIMAL(5,2) NOT NULL CHECK (interest_rate >= 0),
    interest_amount         DECIMAL(12,2) NOT NULL DEFAULT 0,
    principal_remaining     DECIMAL(12,2) NOT NULL,
    interest_remaining      DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_amount            DECIMAL(12,2) NOT NULL,
    amount_paid             DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Late fees
    late_fee_rate           DECIMAL(5,2) DEFAULT 0,
    late_fee_amount         DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Dates
    start_date              DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date                DATE NOT NULL,
    paid_date               DATE,
    confiscated_date        DATE,

    -- Payment plan
    payment_plan_type           payment_plan_type NOT NULL DEFAULT 'minimum_payment',
    loan_term_days              INTEGER NOT NULL,
    requires_minimum_payment    BOOLEAN NOT NULL DEFAULT true,
    minimum_payment_amount      DECIMAL(12,2),
    next_payment_due_date       DATE,
    grace_period_days           INTEGER NOT NULL DEFAULT 5,

    -- Installments (if applicable)
    number_of_installments      INTEGER,
    installment_amount          DECIMAL(12,2),

    -- Status
    status                  loan_status NOT NULL DEFAULT 'active',
    days_overdue            INTEGER NOT NULL DEFAULT 0,

    -- Renewal info
    renewed_from_id         BIGINT REFERENCES loans(id),
    renewal_count           INTEGER NOT NULL DEFAULT 0,

    -- Notes
    notes                   TEXT,

    -- Audit
    created_by              BIGINT REFERENCES users(id),
    updated_by              BIGINT REFERENCES users(id),

    -- Timestamps
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_loans_branch_id ON loans(branch_id);
CREATE INDEX idx_loans_customer_id ON loans(customer_id);
CREATE INDEX idx_loans_item_id ON loans(item_id);
CREATE INDEX idx_loans_status ON loans(status);
CREATE INDEX idx_loans_due_date ON loans(due_date) WHERE status = 'active';
CREATE INDEX idx_loans_loan_number ON loans(loan_number);
CREATE INDEX idx_loans_deleted_at ON loans(deleted_at) WHERE deleted_at IS NULL;

-- Prevent same item in multiple active loans
CREATE UNIQUE INDEX idx_loans_active_item
ON loans(item_id)
WHERE status = 'active' AND deleted_at IS NULL;

CREATE TRIGGER loans_updated_at
    BEFORE UPDATE ON loans
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Loan installments table (for installment plans)
CREATE TABLE loan_installments (
    id              BIGSERIAL PRIMARY KEY,
    loan_id         BIGINT NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    installment_number INTEGER NOT NULL,
    due_date        DATE NOT NULL,
    principal_amount DECIMAL(12,2) NOT NULL,
    interest_amount DECIMAL(12,2) NOT NULL,
    total_amount    DECIMAL(12,2) NOT NULL,
    amount_paid     DECIMAL(12,2) NOT NULL DEFAULT 0,
    is_paid         BOOLEAN NOT NULL DEFAULT false,
    paid_date       DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_loan_installments_loan_id ON loan_installments(loan_id);
CREATE UNIQUE INDEX idx_loan_installments_unique ON loan_installments(loan_id, installment_number);

CREATE TRIGGER loan_installments_updated_at
    BEFORE UPDATE ON loan_installments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to generate loan number
CREATE OR REPLACE FUNCTION generate_loan_number()
RETURNS VARCHAR AS $$
DECLARE
    seq_num INTEGER;
    year_str VARCHAR(4);
BEGIN
    year_str := TO_CHAR(NOW(), 'YYYY');
    SELECT COALESCE(MAX(CAST(SUBSTRING(loan_number FROM 'LN-\d{4}-(\d+)') AS INTEGER)), 0) + 1
    INTO seq_num
    FROM loans
    WHERE loan_number LIKE 'LN-' || year_str || '-%';

    RETURN 'LN-' || year_str || '-' || LPAD(seq_num::TEXT, 6, '0');
END;
$$ LANGUAGE plpgsql;
