-- Drop accounting tables
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS expense_categories;
DROP TABLE IF EXISTS daily_balances;
DROP TABLE IF EXISTS accounting_entry_lines;
DROP TABLE IF EXISTS accounting_entries;
DROP TABLE IF EXISTS accounts;

-- Drop types
DROP TYPE IF EXISTS entry_type;
DROP TYPE IF EXISTS account_type;
