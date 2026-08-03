ALTER TABLE measurement_items ADD COLUMN score_direction VARCHAR(16)
  CHECK (score_direction IN ('higher_is_better', 'lower_is_better'));
