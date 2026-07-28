CREATE INDEX idx_organizations_tenant_id ON organizations (tenant_id);
CREATE INDEX idx_customers_tenant_id_is_active ON customers (tenant_id, is_active);
