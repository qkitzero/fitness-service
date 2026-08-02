CREATE TABLE judgments (
  measurement_id VARCHAR(36) PRIMARY KEY,
  advice VARCHAR(2000) CHECK (advice <> ''),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);
ALTER TABLE judgments ADD CONSTRAINT fk_judgments_measurement
  FOREIGN KEY (measurement_id) REFERENCES measurements (id) ON DELETE CASCADE;
