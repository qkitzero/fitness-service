CREATE TABLE measurement_item_elements (
  measurement_item_id VARCHAR(36) NOT NULL,
  element VARCHAR(32) NOT NULL CHECK (element IN ('muscle_strength', 'muscle_endurance', 'flexibility', 'agility', 'balance', 'mobility')),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  PRIMARY KEY (measurement_item_id, element)
);
ALTER TABLE measurement_item_elements ADD CONSTRAINT fk_measurement_item_elements_measurement_item
  FOREIGN KEY (measurement_item_id) REFERENCES measurement_items (id) ON DELETE CASCADE;
