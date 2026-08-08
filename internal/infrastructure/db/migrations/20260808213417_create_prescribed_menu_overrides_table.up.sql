CREATE TABLE prescribed_menu_overrides (
  id VARCHAR(36) PRIMARY KEY,
  measurement_id VARCHAR(36) NOT NULL,
  sort_order SMALLINT NOT NULL CHECK (sort_order >= 1),
  element VARCHAR(32) CHECK (element IN ('muscle_strength', 'muscle_endurance', 'flexibility', 'agility', 'balance', 'mobility')),
  part VARCHAR(16) CHECK (part IN ('upper_limb', 'lower_limb', 'whole_body')),
  training_menu_id VARCHAR(36) NOT NULL,
  amount SMALLINT NOT NULL CHECK (amount BETWEEN 1 AND 999),
  unit VARCHAR(16) NOT NULL CHECK (unit IN ('reps', 'seconds', 'minutes')),
  sets SMALLINT NOT NULL CHECK (sets BETWEEN 1 AND 99),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  CHECK ((element IS NULL) = (part IS NULL)),
  UNIQUE (measurement_id, sort_order)
);
ALTER TABLE prescribed_menu_overrides ADD CONSTRAINT fk_prescribed_menu_overrides_measurement
  FOREIGN KEY (measurement_id) REFERENCES measurements (id) ON DELETE CASCADE;
ALTER TABLE prescribed_menu_overrides ADD CONSTRAINT fk_prescribed_menu_overrides_training_menu
  FOREIGN KEY (training_menu_id) REFERENCES training_menus (id) ON DELETE RESTRICT;
ALTER TABLE prescribed_menu_overrides ADD CONSTRAINT fk_prescribed_menu_overrides_training_menu_label
  FOREIGN KEY (training_menu_id, element, part) REFERENCES training_menus (id, element, part) ON DELETE RESTRICT;
CREATE INDEX idx_prescribed_menu_overrides_training_menu_id ON prescribed_menu_overrides (training_menu_id);
