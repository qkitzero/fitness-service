ALTER TABLE measurement_items ALTER COLUMN side_mode SET DEFAULT 'none';

ALTER TABLE measurement_items ADD COLUMN bilateral BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE measurement_items
SET bilateral = TRUE, updated_at = NOW() AT TIME ZONE 'UTC'
WHERE side_mode = 'bilateral';

ALTER TABLE measurement_items ALTER COLUMN bilateral DROP DEFAULT;
