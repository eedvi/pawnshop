-- Add description column to audit_logs table
-- This allows storing human-readable descriptions of audit actions

ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS description TEXT;

-- Add index for searching by description
CREATE INDEX IF NOT EXISTS idx_audit_logs_description ON audit_logs USING gin(to_tsvector('spanish', COALESCE(description, '')));

-- Add comment
COMMENT ON COLUMN audit_logs.description IS 'Human-readable description of the audit action';
