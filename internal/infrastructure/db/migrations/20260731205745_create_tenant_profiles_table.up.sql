CREATE TABLE tenant_profiles (
  tenant_id VARCHAR(36) PRIMARY KEY,
  postal_code VARCHAR(16),
  prefecture VARCHAR(16),
  city VARCHAR(255),
  street VARCHAR(255),
  building VARCHAR(255),
  phone VARCHAR(16),
  email VARCHAR(255),
  homepage_url VARCHAR(255),
  note VARCHAR(255),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);
