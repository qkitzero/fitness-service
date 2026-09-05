ALTER TABLE measurement_items DROP CONSTRAINT IF EXISTS measurement_items_side_aggregation_check;
ALTER TABLE measurement_items ADD CONSTRAINT measurement_items_side_aggregation_check
  CHECK (side_aggregation IN ('mean', 'best'));
