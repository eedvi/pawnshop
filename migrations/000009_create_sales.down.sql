DROP FUNCTION IF EXISTS generate_sale_number();
DROP TRIGGER IF EXISTS sales_updated_at ON sales;
DROP TABLE IF EXISTS sales;
DROP TYPE IF EXISTS sale_status;
