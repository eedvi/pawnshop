-- Add missing item_status enum values
ALTER TYPE item_status ADD VALUE IF NOT EXISTS 'pawned' AFTER 'available';
ALTER TYPE item_status ADD VALUE IF NOT EXISTS 'for_sale' AFTER 'collateral';
ALTER TYPE item_status ADD VALUE IF NOT EXISTS 'in_transfer' AFTER 'transferred';
