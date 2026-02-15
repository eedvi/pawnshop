-- Create sale status enum
CREATE TYPE sale_status AS ENUM (
    'completed',
    'pending',
    'cancelled',
    'refunded'
);

-- Create sales table
CREATE TABLE sales (
    id              BIGSERIAL PRIMARY KEY,
    sale_number     VARCHAR(50) NOT NULL UNIQUE,
    branch_id       BIGINT NOT NULL REFERENCES branches(id),
    item_id         BIGINT NOT NULL REFERENCES items(id),
    customer_id     BIGINT REFERENCES customers(id), -- Optional customer

    -- Pricing
    original_price  DECIMAL(12,2) NOT NULL,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    final_price     DECIMAL(12,2) NOT NULL,
    tax_amount      DECIMAL(12,2) DEFAULT 0,
    total_amount    DECIMAL(12,2) NOT NULL,

    -- Payment
    payment_method  payment_method NOT NULL DEFAULT 'cash',
    payment_reference VARCHAR(100),

    -- Status
    status          sale_status NOT NULL DEFAULT 'completed',
    sale_date       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Refund info
    refunded_at     TIMESTAMPTZ,
    refunded_by     BIGINT REFERENCES users(id),
    refund_reason   TEXT,

    -- Notes
    notes           TEXT,

    -- Cash session reference
    cash_session_id BIGINT, -- Will reference cash_sessions

    -- Warranty (if applicable)
    warranty_days   INTEGER DEFAULT 0,
    warranty_expiry DATE,

    -- Audit
    created_by      BIGINT REFERENCES users(id),

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_sales_branch_id ON sales(branch_id);
CREATE INDEX idx_sales_item_id ON sales(item_id);
CREATE INDEX idx_sales_customer_id ON sales(customer_id);
CREATE INDEX idx_sales_sale_date ON sales(sale_date);
CREATE INDEX idx_sales_sale_number ON sales(sale_number);
CREATE INDEX idx_sales_status ON sales(status);

CREATE TRIGGER sales_updated_at
    BEFORE UPDATE ON sales
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to generate sale number
CREATE OR REPLACE FUNCTION generate_sale_number()
RETURNS VARCHAR AS $$
DECLARE
    seq_num INTEGER;
    year_str VARCHAR(4);
BEGIN
    year_str := TO_CHAR(NOW(), 'YYYY');
    SELECT COALESCE(MAX(CAST(SUBSTRING(sale_number FROM 'SL-\d{4}-(\d+)') AS INTEGER)), 0) + 1
    INTO seq_num
    FROM sales
    WHERE sale_number LIKE 'SL-' || year_str || '-%';

    RETURN 'SL-' || year_str || '-' || LPAD(seq_num::TEXT, 6, '0');
END;
$$ LANGUAGE plpgsql;
