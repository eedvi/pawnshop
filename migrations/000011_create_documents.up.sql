-- Create document type enum
CREATE TYPE document_type AS ENUM (
    'loan_contract',
    'loan_receipt',
    'payment_receipt',
    'sale_receipt',
    'confiscation_notice',
    'other'
);

-- Create documents table for tracking generated documents
CREATE TABLE documents (
    id              BIGSERIAL PRIMARY KEY,
    branch_id       BIGINT NOT NULL REFERENCES branches(id),
    document_type   document_type NOT NULL,
    document_number VARCHAR(100) NOT NULL,

    -- Reference to related entity
    reference_type  VARCHAR(50) NOT NULL, -- loan, payment, sale
    reference_id    BIGINT NOT NULL,

    -- File info
    file_path       VARCHAR(500),
    file_url        VARCHAR(500),
    file_size       INTEGER,
    mime_type       VARCHAR(100) DEFAULT 'application/pdf',

    -- Content hash for integrity
    content_hash    VARCHAR(64),

    -- Audit
    created_by      BIGINT REFERENCES users(id),

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_documents_branch_id ON documents(branch_id);
CREATE INDEX idx_documents_type ON documents(document_type);
CREATE INDEX idx_documents_reference ON documents(reference_type, reference_id);
CREATE INDEX idx_documents_created_at ON documents(created_at);

-- Create audit log table
CREATE TABLE audit_logs (
    id              BIGSERIAL PRIMARY KEY,
    branch_id       BIGINT REFERENCES branches(id),
    user_id         BIGINT REFERENCES users(id),

    -- Action details
    action          VARCHAR(100) NOT NULL,
    entity_type     VARCHAR(100) NOT NULL,
    entity_id       BIGINT,

    -- Data
    old_values      JSONB,
    new_values      JSONB,
    ip_address      VARCHAR(45),
    user_agent      TEXT,

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_branch_id ON audit_logs(branch_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Create settings table for global and branch-specific settings
CREATE TABLE settings (
    id              BIGSERIAL PRIMARY KEY,
    branch_id       BIGINT REFERENCES branches(id), -- NULL for global settings
    key             VARCHAR(255) NOT NULL,
    value           JSONB NOT NULL,
    description     TEXT,

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_settings_key ON settings(COALESCE(branch_id, 0), key);

CREATE TRIGGER settings_updated_at
    BEFORE UPDATE ON settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default global settings
INSERT INTO settings (key, value, description) VALUES
('company_name', '"Casa de Empeño"', 'Nombre de la empresa'),
('company_phone', '""', 'Teléfono de la empresa'),
('company_email', '""', 'Email de la empresa'),
('company_address', '""', 'Dirección de la empresa'),
('company_logo', '""', 'URL del logo'),
('tax_id', '""', 'NIT de la empresa'),
('default_interest_rate', '10', 'Tasa de interés por defecto (%)'),
('default_loan_term_days', '30', 'Plazo de préstamo por defecto (días)'),
('default_grace_period_days', '5', 'Período de gracia por defecto (días)'),
('require_minimum_payment', 'true', 'Requerir pago mínimo mensual'),
('minimum_payment_percent', '10', 'Porcentaje de pago mínimo (%)'),
('late_fee_rate', '5', 'Tasa de mora (%)'),
('currency', '"GTQ"', 'Moneda por defecto'),
('currency_symbol', '"Q"', 'Símbolo de moneda'),
('date_format', '"DD/MM/YYYY"', 'Formato de fecha'),
('time_format', '"HH:mm"', 'Formato de hora'),
('receipt_footer', '"Gracias por su preferencia"', 'Pie de página de recibos');

-- Create refresh tokens table
CREATE TABLE refresh_tokens (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash      VARCHAR(255) NOT NULL UNIQUE,
    device_info     TEXT,
    ip_address      VARCHAR(45),
    expires_at      TIMESTAMPTZ NOT NULL,
    revoked_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
