ALTER TABLE customers ADD COLUMN organization_id VARCHAR(36);
ALTER TABLE customers ADD CONSTRAINT fk_customers_organization FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE RESTRICT;
CREATE INDEX idx_customers_organization_id ON customers (organization_id);
