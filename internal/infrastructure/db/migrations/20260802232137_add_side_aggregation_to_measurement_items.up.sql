ALTER TABLE measurement_items ADD COLUMN side_aggregation VARCHAR(8) NOT NULL DEFAULT 'mean'
  CHECK (side_aggregation IN ('mean', 'best'));

UPDATE measurement_items SET side_aggregation = 'best', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code IN (
  'eyes_closed_one_leg_stand',
  'eyes_open_one_leg_stand'
);
