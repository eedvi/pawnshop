-- Add Two-Factor Authentication support to users

-- Add 2FA columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS two_factor_enabled BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS two_factor_secret VARCHAR(64);
ALTER TABLE users ADD COLUMN IF NOT EXISTS two_factor_recovery_codes TEXT[];  -- Array of encrypted recovery codes
ALTER TABLE users ADD COLUMN IF NOT EXISTS two_factor_confirmed_at TIMESTAMPTZ;

-- Backup codes tracking
CREATE TABLE two_factor_backup_codes (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash       VARCHAR(255) NOT NULL,  -- Hashed backup code
    used_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2FA login attempts/challenges
CREATE TABLE two_factor_challenges (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    challenge_token VARCHAR(255) NOT NULL UNIQUE,
    ip_address      VARCHAR(45),
    user_agent      TEXT,
    expires_at      TIMESTAMPTZ NOT NULL,
    verified_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add emergency contact fields to customers
ALTER TABLE customers ADD COLUMN IF NOT EXISTS emergency_contact_name VARCHAR(100);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS emergency_contact_phone VARCHAR(50);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS emergency_contact_relationship VARCHAR(50);

-- Add loyalty/points fields to customers
ALTER TABLE customers ADD COLUMN IF NOT EXISTS loyalty_points INTEGER NOT NULL DEFAULT 0;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS loyalty_tier VARCHAR(20) DEFAULT 'standard';  -- standard, silver, gold, platinum
ALTER TABLE customers ADD COLUMN IF NOT EXISTS loyalty_enrolled_at TIMESTAMPTZ;

-- Loyalty points history
CREATE TABLE loyalty_points_history (
    id              BIGSERIAL PRIMARY KEY,
    customer_id     BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    branch_id       BIGINT REFERENCES branches(id),

    -- Points change
    points_change   INTEGER NOT NULL,  -- Positive for earned, negative for redeemed
    points_balance  INTEGER NOT NULL,  -- Balance after this transaction

    -- Reference
    reference_type  VARCHAR(50),  -- loan, payment, redemption, bonus, adjustment
    reference_id    BIGINT,
    description     TEXT,

    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_two_factor_backup_codes_user ON two_factor_backup_codes(user_id);
CREATE INDEX idx_two_factor_challenges_user ON two_factor_challenges(user_id);
CREATE INDEX idx_two_factor_challenges_token ON two_factor_challenges(challenge_token);
CREATE INDEX idx_two_factor_challenges_expires ON two_factor_challenges(expires_at);
CREATE INDEX idx_loyalty_points_history_customer ON loyalty_points_history(customer_id);
