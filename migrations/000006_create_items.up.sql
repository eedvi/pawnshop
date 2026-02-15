-- Create item status enum
CREATE TYPE item_status AS ENUM (
    'available',    -- Available for loan or sale
    'collateral',   -- Currently used as loan collateral
    'sold',         -- Sold
    'confiscated',  -- Taken due to loan default
    'transferred',  -- Transferred to another branch
    'damaged',      -- Damaged
    'lost'          -- Lost
);

-- Create items table
CREATE TABLE items (
    id              BIGSERIAL PRIMARY KEY,
    branch_id       BIGINT NOT NULL REFERENCES branches(id),
    category_id     BIGINT REFERENCES categories(id),
    customer_id     BIGINT REFERENCES customers(id), -- Original owner

    -- Identification
    sku             VARCHAR(100) NOT NULL UNIQUE,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    brand           VARCHAR(100),
    model           VARCHAR(100),
    serial_number   VARCHAR(255),
    color           VARCHAR(50),
    condition       VARCHAR(50) DEFAULT 'good', -- new, excellent, good, fair, poor

    -- Valuation
    appraised_value DECIMAL(12,2) NOT NULL, -- Market value
    loan_value      DECIMAL(12,2) NOT NULL, -- Maximum loan amount
    sale_price      DECIMAL(12,2), -- Price for direct sale

    -- Status
    status          item_status NOT NULL DEFAULT 'available',

    -- Physical details (for jewelry)
    weight          DECIMAL(10,4), -- In grams
    purity          VARCHAR(20), -- e.g., 14k, 18k, .925

    -- Additional info
    notes           TEXT,
    tags            TEXT[], -- Array of tags

    -- Acquisition info
    acquisition_type    VARCHAR(50) DEFAULT 'pawn', -- pawn, purchase, confiscation
    acquisition_date    DATE NOT NULL DEFAULT CURRENT_DATE,
    acquisition_price   DECIMAL(12,2),

    -- Media
    photos          TEXT[], -- Array of photo URLs

    -- Audit
    created_by      BIGINT REFERENCES users(id),
    updated_by      BIGINT REFERENCES users(id),

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_items_branch_id ON items(branch_id);
CREATE INDEX idx_items_category_id ON items(category_id);
CREATE INDEX idx_items_customer_id ON items(customer_id);
CREATE INDEX idx_items_status ON items(status);
CREATE INDEX idx_items_sku ON items(sku);
CREATE INDEX idx_items_serial_number ON items(serial_number) WHERE serial_number IS NOT NULL;
CREATE INDEX idx_items_deleted_at ON items(deleted_at) WHERE deleted_at IS NULL;

-- Full text search index
CREATE INDEX idx_items_search ON items USING gin(to_tsvector('spanish', name || ' ' || COALESCE(description, '') || ' ' || COALESCE(brand, '') || ' ' || COALESCE(model, '')));

CREATE TRIGGER items_updated_at
    BEFORE UPDATE ON items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Item history table for tracking movements
CREATE TABLE item_history (
    id              BIGSERIAL PRIMARY KEY,
    item_id         BIGINT NOT NULL REFERENCES items(id),
    action          VARCHAR(50) NOT NULL, -- created, status_change, transfer, appraisal, etc.
    old_status      item_status,
    new_status      item_status,
    old_branch_id   BIGINT REFERENCES branches(id),
    new_branch_id   BIGINT REFERENCES branches(id),
    reference_type  VARCHAR(50), -- loan, sale, transfer
    reference_id    BIGINT,
    notes           TEXT,
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_item_history_item_id ON item_history(item_id);
CREATE INDEX idx_item_history_created_at ON item_history(created_at);
