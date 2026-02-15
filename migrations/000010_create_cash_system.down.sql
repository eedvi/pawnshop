ALTER TABLE sales DROP CONSTRAINT IF EXISTS fk_sales_cash_session;
ALTER TABLE payments DROP CONSTRAINT IF EXISTS fk_payments_cash_session;
DROP TABLE IF EXISTS cash_movements;
DROP TYPE IF EXISTS cash_movement_type;
DROP TRIGGER IF EXISTS cash_sessions_updated_at ON cash_sessions;
DROP TABLE IF EXISTS cash_sessions;
DROP TYPE IF EXISTS cash_session_status;
DROP TRIGGER IF EXISTS cash_registers_updated_at ON cash_registers;
DROP TABLE IF EXISTS cash_registers;
