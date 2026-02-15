-- Create payment method enum
CREATE TYPE payment_method AS ENUM (
    'cash',
    'card',
    'transfer',
    'check',
    'other'
);

-- Create payment status enum
CREATE TYPE payment_status AS ENUM (
    'completed',
    'pending',
    'reversed',
    'failed'
);

-- Create payments table
CREATE TABLE payments (
    id              BIGSERIAL PRIMARY KEY,
    payment_number  VARCHAR(50) NOT NULL UNIQUE,
    branch_id       BIGINT NOT NULL REFERENCES branches(id),
    loan_id         BIGINT NOT NULL REFERENCES loans(id),
    customer_id     BIGINT NOT NULL REFERENCES customers(id),

    -- Payment details
    amount              DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    principal_amount    DECIMAL(12,2) NOT NULL DEFAULT 0,
    interest_amount     DECIMAL(12,2) NOT NULL DEFAULT 0,
    late_fee_amount     DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Method
    payment_method      payment_method NOT NULL DEFAULT 'cash',
    reference_number    VARCHAR(100), -- For card/transfer payments

    -- Status
    status              payment_status NOT NULL DEFAULT 'completed',
    payment_date        TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Balances after payment
    loan_balance_after      DECIMAL(12,2) NOT NULL,
    interest_balance_after  DECIMAL(12,2) NOT NULL,

    -- Reversal info
    reversed_at         TIMESTAMPTZ,
    reversed_by         BIGINT REFERENCES users(id),
    reversal_reason     TEXT,

    -- Notes
    notes               TEXT,

    -- Cash session reference
    cash_session_id     BIGINT, -- Will reference cash_sessions

    -- Audit
    created_by          BIGINT REFERENCES users(id),

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_payments_branch_id ON payments(branch_id);
CREATE INDEX idx_payments_loan_id ON payments(loan_id);
CREATE INDEX idx_payments_customer_id ON payments(customer_id);
CREATE INDEX idx_payments_payment_date ON payments(payment_date);
CREATE INDEX idx_payments_payment_number ON payments(payment_number);
CREATE INDEX idx_payments_status ON payments(status);

CREATE TRIGGER payments_updated_at
    BEFORE UPDATE ON payments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to generate payment number
CREATE OR REPLACE FUNCTION generate_payment_number()
RETURNS VARCHAR AS $$
DECLARE
    seq_num INTEGER;
    year_str VARCHAR(4);
BEGIN
    year_str := TO_CHAR(NOW(), 'YYYY');
    SELECT COALESCE(MAX(CAST(SUBSTRING(payment_number FROM 'PY-\d{4}-(\d+)') AS INTEGER)), 0) + 1
    INTO seq_num
    FROM payments
    WHERE payment_number LIKE 'PY-' || year_str || '-%';

    RETURN 'PY-' || year_str || '-' || LPAD(seq_num::TEXT, 6, '0');
END;
$$ LANGUAGE plpgsql;
