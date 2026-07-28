DROP INDEX IF EXISTS idx_customers_group_id_is_active;
ALTER TABLE customers DROP COLUMN IF EXISTS is_active;
