-- Remove late_fee_remaining column
ALTER TABLE loans DROP COLUMN IF EXISTS late_fee_remaining;
