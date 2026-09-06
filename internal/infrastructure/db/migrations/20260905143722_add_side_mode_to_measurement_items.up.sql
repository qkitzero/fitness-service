ALTER TABLE measurement_items ADD COLUMN side_mode VARCHAR(32) NOT NULL DEFAULT 'none'
  CHECK (side_mode IN ('none', 'bilateral', 'optional_bilateral'));

UPDATE measurement_items
SET side_mode = 'bilateral', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE bilateral;
