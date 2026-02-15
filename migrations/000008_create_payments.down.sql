DROP FUNCTION IF EXISTS generate_payment_number();
DROP TRIGGER IF EXISTS payments_updated_at ON payments;
DROP TABLE IF EXISTS payments;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_method;
