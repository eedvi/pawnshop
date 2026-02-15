-- Create customers table
CREATE TABLE customers (
    id              BIGSERIAL PRIMARY KEY,
    branch_id       BIGINT NOT NULL REFERENCES branches(id),

    -- Personal info
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    identity_type   VARCHAR(20) NOT NULL DEFAULT 'dpi', -- dpi, passport, other
    identity_number VARCHAR(50) NOT NULL,
    birth_date      DATE,
    gender          VARCHAR(10), -- male, female, other

    -- Contact info
    phone           VARCHAR(50) NOT NULL,
    phone_secondary VARCHAR(50),
    email           VARCHAR(255),
    address         TEXT,
    city            VARCHAR(100),
    state           VARCHAR(100),
    postal_code     VARCHAR(20),

    -- Emergency contact
    emergency_contact_name  VARCHAR(255),
    emergency_contact_phone VARCHAR(50),
    emergency_contact_relation VARCHAR(100),

    -- Business info
    occupation      VARCHAR(255),
    workplace       VARCHAR(255),
    monthly_income  DECIMAL(12,2),

    -- Credit info
    credit_limit    DECIMAL(12,2) NOT NULL DEFAULT 0,
    credit_score    INTEGER DEFAULT 50, -- 0-100
    total_loans     INTEGER NOT NULL DEFAULT 0,
    total_paid      DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_defaulted DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Status
    is_active       BOOLEAN NOT NULL DEFAULT true,
    is_blocked      BOOLEAN NOT NULL DEFAULT false,
    blocked_reason  TEXT,

    -- Notes
    notes           TEXT,

    -- Photo
    photo_url       VARCHAR(500),

    -- Timestamps
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_customers_branch_id ON customers(branch_id);
CREATE INDEX idx_customers_identity ON customers(identity_type, identity_number);
CREATE INDEX idx_customers_phone ON customers(phone);
CREATE INDEX idx_customers_name ON customers(first_name, last_name);
CREATE INDEX idx_customers_deleted_at ON customers(deleted_at) WHERE deleted_at IS NULL;

-- Unique constraint on identity within same branch
CREATE UNIQUE INDEX idx_customers_unique_identity
ON customers(branch_id, identity_type, identity_number)
WHERE deleted_at IS NULL;

CREATE TRIGGER customers_updated_at
    BEFORE UPDATE ON customers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
