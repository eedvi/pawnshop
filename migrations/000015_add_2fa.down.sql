-- Drop loyalty tables
DROP TABLE IF EXISTS loyalty_points_history;

-- Drop 2FA tables
DROP TABLE IF EXISTS two_factor_challenges;
DROP TABLE IF EXISTS two_factor_backup_codes;

-- Remove columns from users
ALTER TABLE users DROP COLUMN IF EXISTS two_factor_enabled;
ALTER TABLE users DROP COLUMN IF EXISTS two_factor_secret;
ALTER TABLE users DROP COLUMN IF EXISTS two_factor_recovery_codes;
ALTER TABLE users DROP COLUMN IF EXISTS two_factor_confirmed_at;

-- Remove columns from customers
ALTER TABLE customers DROP COLUMN IF EXISTS emergency_contact_name;
ALTER TABLE customers DROP COLUMN IF EXISTS emergency_contact_phone;
ALTER TABLE customers DROP COLUMN IF EXISTS emergency_contact_relationship;
ALTER TABLE customers DROP COLUMN IF EXISTS loyalty_points;
ALTER TABLE customers DROP COLUMN IF EXISTS loyalty_tier;
ALTER TABLE customers DROP COLUMN IF EXISTS loyalty_enrolled_at;
