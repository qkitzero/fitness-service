ALTER TABLE customers ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;
CREATE INDEX idx_customers_group_id_is_active ON customers (group_id, is_active);
