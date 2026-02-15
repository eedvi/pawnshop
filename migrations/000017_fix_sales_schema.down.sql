-- Revert sales table schema changes
DROP INDEX IF EXISTS idx_sales_deleted_at;

ALTER TABLE sales DROP COLUMN IF EXISTS sale_type;
ALTER TABLE sales DROP COLUMN IF EXISTS sale_price;
ALTER TABLE sales DROP COLUMN IF EXISTS discount_reason;
ALTER TABLE sales DROP COLUMN IF EXISTS reference_number;
ALTER TABLE sales DROP COLUMN IF EXISTS refund_amount;
ALTER TABLE sales DROP COLUMN IF EXISTS updated_by;
ALTER TABLE sales DROP COLUMN IF EXISTS deleted_at;

-- Note: payment_reference column change is not reversible without data loss
