CREATE TABLE fixed_menus (
  sort_order SMALLINT PRIMARY KEY CHECK (sort_order >= 1),
  training_menu_id VARCHAR(36) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);
ALTER TABLE fixed_menus ADD CONSTRAINT fk_fixed_menus_training_menu
  FOREIGN KEY (training_menu_id) REFERENCES training_menus (id) ON DELETE RESTRICT;
