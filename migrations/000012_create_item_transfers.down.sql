-- Remove transfer column from items
ALTER TABLE items DROP COLUMN IF EXISTS in_transfer;

-- Drop item_transfers table
DROP TABLE IF EXISTS item_transfers;

-- Drop transfer_status type
DROP TYPE IF EXISTS transfer_status;
