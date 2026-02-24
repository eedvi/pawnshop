-- Add delivered_at field to items table
-- This field tracks when the item was physically delivered to the customer

ALTER TABLE items
ADD COLUMN delivered_at TIMESTAMP WITH TIME ZONE;

-- Create index for querying pending deliveries
CREATE INDEX idx_items_delivered_at ON items(delivered_at) WHERE delivered_at IS NULL;

-- Add comment explaining the field
COMMENT ON COLUMN items.delivered_at IS 'Timestamp when the item was physically delivered to the customer after loan payoff';

-- Items that are available but not yet delivered (pending pickup)
-- Query: SELECT * FROM items WHERE status = 'available' AND delivered_at IS NULL AND acquisition_type = 'pawn';

-- Items delivered in last 30 days
-- Query: SELECT * FROM items WHERE delivered_at >= NOW() - INTERVAL '30 days';
