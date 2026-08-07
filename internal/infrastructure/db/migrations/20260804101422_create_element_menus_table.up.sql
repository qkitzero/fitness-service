CREATE TABLE element_menus (
  element VARCHAR(32) NOT NULL CHECK (element IN ('muscle_strength', 'muscle_endurance', 'flexibility', 'agility', 'balance', 'mobility')),
  part VARCHAR(16) NOT NULL CHECK (part IN ('upper_limb', 'lower_limb', 'whole_body')),
  level SMALLINT NOT NULL CHECK (level BETWEEN 1 AND 5),
  training_menu_id VARCHAR(36) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  PRIMARY KEY (element, part, level)
);
ALTER TABLE element_menus ADD CONSTRAINT fk_element_menus_training_menu
  FOREIGN KEY (training_menu_id, element, part) REFERENCES training_menus (id, element, part) ON DELETE RESTRICT;
