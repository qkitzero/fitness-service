CREATE TABLE training_menus (
  id VARCHAR(36) PRIMARY KEY,
  code VARCHAR(64) NOT NULL UNIQUE CHECK (code ~ '^[a-z0-9_]+$'),
  name VARCHAR(255) NOT NULL CHECK (name <> ''),
  element VARCHAR(32) NOT NULL CHECK (element IN ('muscle_strength', 'muscle_endurance', 'flexibility', 'agility', 'balance', 'mobility')),
  part VARCHAR(16) NOT NULL CHECK (part IN ('upper_limb', 'lower_limb', 'whole_body')),
  amount SMALLINT NOT NULL CHECK (amount BETWEEN 1 AND 999),
  unit VARCHAR(16) NOT NULL CHECK (unit IN ('reps', 'seconds', 'minutes')),
  sets SMALLINT NOT NULL CHECK (sets BETWEEN 1 AND 99),
  instruction VARCHAR(2000) NOT NULL CHECK (instruction <> ''),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  UNIQUE (id, element, part)
);
