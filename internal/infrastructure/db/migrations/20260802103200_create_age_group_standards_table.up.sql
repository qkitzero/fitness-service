CREATE TABLE age_group_standards (
  id VARCHAR(36) PRIMARY KEY,
  measurement_item_id VARCHAR(36) NOT NULL,
  gender VARCHAR(8) NOT NULL CHECK (gender IN ('male', 'female')),
  age_from SMALLINT NOT NULL CHECK (age_from BETWEEN 0 AND 150),
  age_to SMALLINT NOT NULL CHECK (age_to BETWEEN 0 AND 150),
  mean NUMERIC(6, 2) NOT NULL CHECK (mean >= 0),
  standard_deviation NUMERIC(6, 2) NOT NULL CHECK (standard_deviation > 0),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  CHECK (age_from <= age_to),
  UNIQUE (measurement_item_id, gender, age_from)
);
ALTER TABLE age_group_standards ADD CONSTRAINT fk_age_group_standards_measurement_item
  FOREIGN KEY (measurement_item_id) REFERENCES measurement_items (id) ON DELETE RESTRICT;
