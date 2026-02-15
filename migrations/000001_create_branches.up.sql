-- Create branches table
CREATE TABLE branches (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    code            VARCHAR(50) NOT NULL UNIQUE,
    address         TEXT,
    phone           VARCHAR(50),
    email           VARCHAR(255),
    is_active       BOOLEAN NOT NULL DEFAULT true,
    timezone        VARCHAR(100) NOT NULL DEFAULT 'America/Guatemala',
    currency        VARCHAR(10) NOT NULL DEFAULT 'GTQ',

    -- Business settings
    default_interest_rate   DECIMAL(5,2) NOT NULL DEFAULT 10.00,
    default_loan_term_days  INTEGER NOT NULL DEFAULT 30,
    default_grace_period    INTEGER NOT NULL DEFAULT 5,

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- Create index for soft delete queries
CREATE INDEX idx_branches_deleted_at ON branches(deleted_at) WHERE deleted_at IS NULL;

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER branches_updated_at
    BEFORE UPDATE ON branches
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default branch
INSERT INTO branches (name, code, address, phone, email)
VALUES ('Sucursal Principal', 'MAIN', 'Dirección Principal', '0000-0000', 'main@pawnshop.com');
