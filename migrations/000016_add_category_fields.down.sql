-- Remove added fields from categories table
ALTER TABLE categories DROP COLUMN IF EXISTS icon;
ALTER TABLE categories DROP COLUMN IF EXISTS default_interest_rate;
ALTER TABLE categories DROP COLUMN IF EXISTS min_loan_amount;
ALTER TABLE categories DROP COLUMN IF EXISTS max_loan_amount;
ALTER TABLE categories DROP COLUMN IF EXISTS loan_to_value_ratio;
ALTER TABLE categories DROP COLUMN IF EXISTS sort_order;
