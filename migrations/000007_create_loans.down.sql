DROP FUNCTION IF EXISTS generate_loan_number();
DROP TRIGGER IF EXISTS loan_installments_updated_at ON loan_installments;
DROP TABLE IF EXISTS loan_installments;
DROP TRIGGER IF EXISTS loans_updated_at ON loans;
DROP TABLE IF EXISTS loans;
DROP TYPE IF EXISTS payment_plan_type;
DROP TYPE IF EXISTS loan_status;
