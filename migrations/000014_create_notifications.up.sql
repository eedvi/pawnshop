-- Notifications system

-- Notification types
CREATE TYPE notification_type AS ENUM (
    'loan_due_reminder',      -- Loan due date approaching
    'loan_overdue',           -- Loan is overdue
    'minimum_payment_due',    -- Minimum payment reminder
    'loan_confiscation',      -- Item about to be confiscated
    'payment_received',       -- Payment confirmation
    'loan_paid_off',          -- Loan fully paid
    'system_alert',           -- System notifications
    'custom'                  -- Custom notifications
);

-- Notification channels
CREATE TYPE notification_channel AS ENUM ('email', 'sms', 'whatsapp', 'push', 'internal');

-- Notification status
CREATE TYPE notification_status AS ENUM ('pending', 'sent', 'delivered', 'failed', 'cancelled');

-- Notification templates
CREATE TABLE notification_templates (
    id              BIGSERIAL PRIMARY KEY,
    code            VARCHAR(50) NOT NULL UNIQUE,
    name            VARCHAR(100) NOT NULL,
    notification_type notification_type NOT NULL,
    channel         notification_channel NOT NULL,

    -- Template content
    subject         VARCHAR(200),  -- For email
    body            TEXT NOT NULL,

    -- Variables available (JSON array of variable names)
    variables       JSONB NOT NULL DEFAULT '[]',

    -- Settings
    is_active       BOOLEAN NOT NULL DEFAULT true,
    is_system       BOOLEAN NOT NULL DEFAULT false,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Notification queue
CREATE TABLE notifications (
    id              BIGSERIAL PRIMARY KEY,

    -- Type and channel
    notification_type notification_type NOT NULL,
    channel         notification_channel NOT NULL,
    template_id     BIGINT REFERENCES notification_templates(id),

    -- Recipient
    customer_id     BIGINT REFERENCES customers(id),
    user_id         BIGINT REFERENCES users(id),
    recipient_email VARCHAR(255),
    recipient_phone VARCHAR(50),

    -- Reference
    reference_type  VARCHAR(50),  -- loan, payment, etc.
    reference_id    BIGINT,
    branch_id       BIGINT REFERENCES branches(id),

    -- Content
    subject         VARCHAR(200),
    body            TEXT NOT NULL,
    data            JSONB,  -- Additional data/variables

    -- Scheduling
    scheduled_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Status tracking
    status          notification_status NOT NULL DEFAULT 'pending',
    sent_at         TIMESTAMPTZ,
    delivered_at    TIMESTAMPTZ,
    failed_at       TIMESTAMPTZ,
    error_message   TEXT,
    retry_count     INTEGER NOT NULL DEFAULT 0,
    max_retries     INTEGER NOT NULL DEFAULT 3,

    -- Audit
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Notification preferences per customer
CREATE TABLE customer_notification_preferences (
    id              BIGSERIAL PRIMARY KEY,
    customer_id     BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,

    -- Channel preferences
    email_enabled   BOOLEAN NOT NULL DEFAULT true,
    sms_enabled     BOOLEAN NOT NULL DEFAULT true,
    whatsapp_enabled BOOLEAN NOT NULL DEFAULT false,

    -- Type preferences
    due_reminders   BOOLEAN NOT NULL DEFAULT true,
    payment_confirmations BOOLEAN NOT NULL DEFAULT true,
    promotional     BOOLEAN NOT NULL DEFAULT false,

    -- Reminder settings
    reminder_days_before INTEGER NOT NULL DEFAULT 3,  -- Days before due date to send reminder

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(customer_id)
);

-- Internal notifications for employees
CREATE TABLE internal_notifications (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    branch_id       BIGINT REFERENCES branches(id),

    -- Content
    title           VARCHAR(200) NOT NULL,
    message         TEXT NOT NULL,
    notification_type VARCHAR(50) NOT NULL DEFAULT 'info',  -- info, warning, error, success

    -- Reference
    reference_type  VARCHAR(50),
    reference_id    BIGINT,
    action_url      VARCHAR(500),

    -- Status
    is_read         BOOLEAN NOT NULL DEFAULT false,
    read_at         TIMESTAMPTZ,

    -- Priority
    priority        INTEGER NOT NULL DEFAULT 0,  -- 0=normal, 1=high, 2=urgent

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_notifications_status ON notifications(status) WHERE status = 'pending';
CREATE INDEX idx_notifications_scheduled ON notifications(scheduled_at) WHERE status = 'pending';
CREATE INDEX idx_notifications_customer ON notifications(customer_id);
CREATE INDEX idx_notifications_reference ON notifications(reference_type, reference_id);
CREATE INDEX idx_internal_notifications_user ON internal_notifications(user_id);
CREATE INDEX idx_internal_notifications_unread ON internal_notifications(user_id, is_read) WHERE is_read = false;

-- Triggers
CREATE TRIGGER notification_templates_updated_at
    BEFORE UPDATE ON notification_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER notifications_updated_at
    BEFORE UPDATE ON notifications
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER customer_notification_preferences_updated_at
    BEFORE UPDATE ON customer_notification_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default templates
INSERT INTO notification_templates (code, name, notification_type, channel, subject, body, variables, is_system) VALUES
    ('LOAN_DUE_EMAIL', 'Recordatorio de Vencimiento (Email)', 'loan_due_reminder', 'email',
     'Recordatorio: Su préstamo vence pronto',
     'Estimado(a) {{customer_name}},\n\nLe recordamos que su préstamo #{{loan_number}} vence el {{due_date}}.\n\nMonto pendiente: {{currency}}{{amount_due}}\n\nPara evitar cargos adicionales, le invitamos a realizar su pago antes de la fecha de vencimiento.\n\nGracias por su preferencia.',
     '["customer_name", "loan_number", "due_date", "amount_due", "currency"]', true),

    ('LOAN_DUE_SMS', 'Recordatorio de Vencimiento (SMS)', 'loan_due_reminder', 'sms',
     NULL,
     'Su prestamo #{{loan_number}} vence el {{due_date}}. Monto: {{currency}}{{amount_due}}. Evite cargos adicionales pagando a tiempo.',
     '["loan_number", "due_date", "amount_due", "currency"]', true),

    ('LOAN_OVERDUE_EMAIL', 'Préstamo Vencido (Email)', 'loan_overdue', 'email',
     'IMPORTANTE: Su préstamo está vencido',
     'Estimado(a) {{customer_name}},\n\nSu préstamo #{{loan_number}} venció el {{due_date}} y se encuentra en mora.\n\nMonto vencido: {{currency}}{{amount_due}}\nDías de mora: {{days_overdue}}\nCargo por mora: {{currency}}{{late_fee}}\n\nPor favor, acérquese a nuestra sucursal para regularizar su situación.\n\nAtentamente.',
     '["customer_name", "loan_number", "due_date", "amount_due", "days_overdue", "late_fee", "currency"]', true),

    ('PAYMENT_RECEIVED_EMAIL', 'Confirmación de Pago (Email)', 'payment_received', 'email',
     'Confirmación de pago recibido',
     'Estimado(a) {{customer_name}},\n\nHemos recibido su pago por {{currency}}{{amount}} para el préstamo #{{loan_number}}.\n\nNo. de recibo: {{payment_number}}\nFecha: {{payment_date}}\nSaldo pendiente: {{currency}}{{remaining_balance}}\n\nGracias por su pago.',
     '["customer_name", "loan_number", "amount", "payment_number", "payment_date", "remaining_balance", "currency"]', true),

    ('CONFISCATION_WARNING', 'Aviso de Confiscación', 'loan_confiscation', 'email',
     'AVISO IMPORTANTE: Proceso de confiscación',
     'Estimado(a) {{customer_name}},\n\nDebido a que su préstamo #{{loan_number}} se encuentra vencido por más de {{days_overdue}} días, le informamos que el artículo en prenda será confiscado si no regulariza su situación antes del {{confiscation_date}}.\n\nArtículo: {{item_name}}\nMonto total adeudado: {{currency}}{{total_due}}\n\nPor favor, contáctenos inmediatamente.',
     '["customer_name", "loan_number", "days_overdue", "confiscation_date", "item_name", "total_due", "currency"]', true);
