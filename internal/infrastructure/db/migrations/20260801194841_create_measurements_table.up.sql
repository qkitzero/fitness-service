CREATE TABLE measurements (
  id VARCHAR(36) PRIMARY KEY,
  customer_id VARCHAR(36) NOT NULL,
  measured_on DATE NOT NULL,
  measured_by VARCHAR(255) NOT NULL CHECK (measured_by <> ''),
  age_at_measurement SMALLINT NOT NULL CHECK (age_at_measurement BETWEEN 0 AND 150),
  updated_by VARCHAR(255) NOT NULL CHECK (updated_by <> ''),
  is_draft BOOLEAN NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);
ALTER TABLE measurements ADD CONSTRAINT fk_measurements_customer
  FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE RESTRICT;
CREATE INDEX idx_measurements_customer_id ON measurements (customer_id);

CREATE TABLE measurement_entries (
  id VARCHAR(36) PRIMARY KEY,
  measurement_id VARCHAR(36) NOT NULL,
  measurement_item_id VARCHAR(36) NOT NULL,
  unmeasurable BOOLEAN NOT NULL,
  note VARCHAR(255) CHECK (note <> ''),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  UNIQUE (measurement_id, measurement_item_id)
);
ALTER TABLE measurement_entries ADD CONSTRAINT fk_measurement_entries_measurement
  FOREIGN KEY (measurement_id) REFERENCES measurements (id) ON DELETE CASCADE;
ALTER TABLE measurement_entries ADD CONSTRAINT fk_measurement_entries_measurement_item
  FOREIGN KEY (measurement_item_id) REFERENCES measurement_items (id) ON DELETE RESTRICT;
CREATE INDEX idx_measurement_entries_measurement_item_id ON measurement_entries (measurement_item_id);

CREATE TABLE measurement_values (
  id VARCHAR(36) PRIMARY KEY,
  measurement_entry_id VARCHAR(36) NOT NULL,
  trial_index SMALLINT NOT NULL CHECK (trial_index >= 1),
  side VARCHAR(8) NOT NULL CHECK (side IN ('none', 'left', 'right')),
  value NUMERIC(6, 2) CHECK (value >= 0),
  value_secondary NUMERIC(6, 2) CHECK (value_secondary >= 0),
  value_choice VARCHAR(32) CHECK (value_choice <> ''),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  UNIQUE (measurement_entry_id, trial_index, side)
);
ALTER TABLE measurement_values ADD CONSTRAINT fk_measurement_values_measurement_entry
  FOREIGN KEY (measurement_entry_id) REFERENCES measurement_entries (id) ON DELETE CASCADE;
