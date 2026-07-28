DROP INDEX IF EXISTS idx_customers_organization_id;
ALTER TABLE customers DROP CONSTRAINT IF EXISTS fk_customers_organization;
ALTER TABLE customers DROP COLUMN IF EXISTS organization_id;
