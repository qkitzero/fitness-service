ALTER TABLE measurement_items ADD COLUMN trial_aggregation VARCHAR(8) NOT NULL DEFAULT 'best'
  CHECK (trial_aggregation IN ('mean', 'best'));
