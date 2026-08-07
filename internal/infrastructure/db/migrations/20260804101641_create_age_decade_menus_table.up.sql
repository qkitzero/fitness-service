CREATE TABLE age_decade_menus (
  decade SMALLINT NOT NULL CHECK (decade BETWEEN 10 AND 90 AND decade % 10 = 0),
  sort_order SMALLINT NOT NULL CHECK (sort_order >= 1),
  training_menu_id VARCHAR(36) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  PRIMARY KEY (decade, sort_order)
);
ALTER TABLE age_decade_menus ADD CONSTRAINT fk_age_decade_menus_training_menu
  FOREIGN KEY (training_menu_id) REFERENCES training_menus (id) ON DELETE RESTRICT;
