-- Fix sales table schema to match domain model and repository

-- Add missing columns
ALTER TABLE sales ADD COLUMN IF NOT EXISTS sale_type VARCHAR(50) NOT NULL DEFAULT 'direct';
ALTER TABLE sales ADD COLUMN IF NOT EXISTS sale_price DECIMAL(12,2);
ALTER TABLE sales ADD COLUMN IF NOT EXISTS discount_reason TEXT;
ALTER TABLE sales ADD COLUMN IF NOT EXISTS reference_number VARCHAR(100);
ALTER TABLE sales ADD COLUMN IF NOT EXISTS refund_amount DECIMAL(12,2);
ALTER TABLE sales ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id);
ALTER TABLE sales ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Rename columns to match expected names
ALTER TABLE sales RENAME COLUMN payment_reference TO payment_reference_old;
ALTER TABLE sales ADD COLUMN payment_reference_temp VARCHAR(100);
UPDATE sales SET payment_reference_temp = payment_reference_old;
ALTER TABLE sales DROP COLUMN payment_reference_old;

-- Copy original_price to sale_price if not set
UPDATE sales SET sale_price = original_price WHERE sale_price IS NULL;

-- Copy payment_reference to reference_number
UPDATE sales SET reference_number = payment_reference_temp WHERE reference_number IS NULL;
ALTER TABLE sales DROP COLUMN IF EXISTS payment_reference_temp;

-- Make sale_price NOT NULL after populating
ALTER TABLE sales ALTER COLUMN sale_price SET NOT NULL;

-- Add index for deleted_at for soft delete queries
CREATE INDEX IF NOT EXISTS idx_sales_deleted_at ON sales(deleted_at);
