-- Add late_fee_remaining column to track remaining late fees (separate from historical total)
ALTER TABLE loans ADD COLUMN IF NOT EXISTS late_fee_remaining DECIMAL(12, 2) DEFAULT 0;

-- Initialize late_fee_remaining with current late_fee_amount for existing loans
-- (for unpaid loans, remaining = amount; for paid loans, remaining = 0)
UPDATE loans
SET late_fee_remaining = CASE
    WHEN status IN ('paid', 'confiscated') THEN 0
    ELSE late_fee_amount
END;
