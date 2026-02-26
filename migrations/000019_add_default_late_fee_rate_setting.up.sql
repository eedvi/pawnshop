-- Add default late fee rate setting (1% per day as default)
INSERT INTO settings (key, value, description, branch_id)
VALUES (
    'default_late_fee_rate',
    '1.0',
    'Default late fee rate percentage per day for overdue loans',
    NULL
)
ON CONFLICT (key, COALESCE(branch_id, 0)) DO NOTHING;
