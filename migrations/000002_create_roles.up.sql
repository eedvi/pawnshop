-- Create roles table
CREATE TABLE roles (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,
    display_name    VARCHAR(255) NOT NULL,
    description     TEXT,
    permissions     JSONB NOT NULL DEFAULT '[]',
    is_system       BOOLEAN NOT NULL DEFAULT false,

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default roles with permissions
INSERT INTO roles (name, display_name, description, permissions, is_system) VALUES
('super_admin', 'Super Administrador', 'Acceso total al sistema',
 '["*"]', true),

('admin', 'Administrador', 'Administrador de sucursal',
 '["users.read", "users.create", "users.update", "customers.*", "items.*", "loans.*", "payments.*", "sales.*", "cash.*", "reports.*", "settings.read"]', true),

('manager', 'Gerente', 'Gerente de sucursal',
 '["customers.*", "items.*", "loans.*", "payments.*", "sales.*", "cash.*", "reports.read"]', true),

('cashier', 'Cajero', 'Operador de caja',
 '["customers.read", "customers.create", "items.read", "items.create", "loans.read", "loans.create", "payments.*", "sales.*", "cash.own"]', true),

('seller', 'Vendedor', 'Vendedor de artículos',
 '["customers.read", "items.read", "sales.read", "sales.create"]', true);
