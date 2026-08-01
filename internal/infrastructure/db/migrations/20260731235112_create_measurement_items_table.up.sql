CREATE TABLE measurement_items (
  id VARCHAR(36) PRIMARY KEY,
  code VARCHAR(64) NOT NULL UNIQUE CHECK (code ~ '^[a-z0-9_]+$'),
  name VARCHAR(255) NOT NULL CHECK (name <> ''),
  category VARCHAR(32) NOT NULL CHECK (category IN ('vital', 'physique', 'body_composition', 'motor_function')),
  unit VARCHAR(16) NOT NULL CHECK (unit IN ('kg', 'cm', 'sec', 'count', 'mmHg', 'percent', 'bpm')),
  trial_count SMALLINT NOT NULL CHECK (trial_count >= 1),
  bilateral BOOLEAN NOT NULL,
  value_type VARCHAR(16) NOT NULL CHECK (value_type IN ('numeric', 'paired', 'choice')),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);
