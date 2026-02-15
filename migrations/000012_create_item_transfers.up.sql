-- Item transfers between branches
CREATE TYPE transfer_status AS ENUM ('pending', 'in_transit', 'completed', 'cancelled');

CREATE TABLE item_transfers (
    id              BIGSERIAL PRIMARY KEY,
    transfer_number VARCHAR(50) NOT NULL UNIQUE,

    -- Item being transferred
    item_id         BIGINT NOT NULL REFERENCES items(id),

    -- Source and destination branches
    from_branch_id  BIGINT NOT NULL REFERENCES branches(id),
    to_branch_id    BIGINT NOT NULL REFERENCES branches(id),

    -- Status tracking
    status          transfer_status NOT NULL DEFAULT 'pending',

    -- Users involved
    requested_by    BIGINT NOT NULL REFERENCES users(id),
    approved_by     BIGINT REFERENCES users(id),
    received_by     BIGINT REFERENCES users(id),

    -- Dates
    requested_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_at     TIMESTAMPTZ,
    shipped_at      TIMESTAMPTZ,
    received_at     TIMESTAMPTZ,
    cancelled_at    TIMESTAMPTZ,

    -- Notes
    request_notes   TEXT,
    approval_notes  TEXT,
    receipt_notes   TEXT,
    cancellation_reason TEXT,

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT different_branches CHECK (from_branch_id != to_branch_id)
);

-- Indexes
CREATE INDEX idx_item_transfers_item ON item_transfers(item_id);
CREATE INDEX idx_item_transfers_from_branch ON item_transfers(from_branch_id);
CREATE INDEX idx_item_transfers_to_branch ON item_transfers(to_branch_id);
CREATE INDEX idx_item_transfers_status ON item_transfers(status);
CREATE INDEX idx_item_transfers_number ON item_transfers(transfer_number);

-- Trigger for updated_at
CREATE TRIGGER item_transfers_updated_at
    BEFORE UPDATE ON item_transfers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add transfer status to items
ALTER TABLE items ADD COLUMN IF NOT EXISTS in_transfer BOOLEAN NOT NULL DEFAULT false;
