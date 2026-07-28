ALTER TABLE customers RENAME COLUMN group_id TO tenant_id;
ALTER INDEX idx_customers_group_id_is_active RENAME TO idx_customers_tenant_id_is_active;
ALTER TABLE organizations RENAME COLUMN group_id TO tenant_id;
ALTER INDEX idx_organizations_group_id RENAME TO idx_organizations_tenant_id;
