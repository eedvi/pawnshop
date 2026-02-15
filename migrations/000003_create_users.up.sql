-- Create users table
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    branch_id       BIGINT REFERENCES branches(id) ON DELETE SET NULL,
    role_id         BIGINT NOT NULL REFERENCES roles(id),

    -- Credentials
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,

    -- Personal info
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    phone           VARCHAR(50),
    avatar_url      VARCHAR(500),

    -- Status
    is_active       BOOLEAN NOT NULL DEFAULT true,
    email_verified  BOOLEAN NOT NULL DEFAULT false,

    -- Security
    failed_login_attempts   INTEGER NOT NULL DEFAULT 0,
    locked_until            TIMESTAMPTZ,
    password_changed_at     TIMESTAMPTZ,
    last_login_at           TIMESTAMPTZ,
    last_login_ip           VARCHAR(45),

    -- 2FA
    two_factor_enabled      BOOLEAN NOT NULL DEFAULT false,
    two_factor_secret       VARCHAR(255),

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_users_branch_id ON users(branch_id);
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;

-- Trigger for updated_at
CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default super admin user (password: admin123 - CHANGE IN PRODUCTION!)
-- Argon2id hash for 'admin123'
INSERT INTO users (branch_id, role_id, email, password_hash, first_name, last_name, is_active, email_verified)
SELECT
    1,
    (SELECT id FROM roles WHERE name = 'super_admin'),
    'admin@pawnshop.com',
    '$argon2id$v=19$m=65536,t=3,p=2$Rx8KhjlnxPmeOK0LY6z0ow$yg/SKFaf0DwZ0Q/xN8Ho33UE2UJGOt2wL2rRGmSeljA',
    'Admin',
    'Sistema',
    true,
    true;
