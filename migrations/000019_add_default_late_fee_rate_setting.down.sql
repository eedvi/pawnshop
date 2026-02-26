-- Remove default late fee rate setting
DELETE FROM settings
WHERE key = 'default_late_fee_rate' AND branch_id IS NULL;
