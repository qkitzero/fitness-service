ALTER TABLE measurement_items DROP COLUMN IF EXISTS bilateral;

ALTER TABLE measurement_items ALTER COLUMN side_mode DROP DEFAULT;
