UPDATE measurement_items
SET unit = 'cm', value_type = 'choice', score_direction = NULL, updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code = 'stand_up_test';

UPDATE measurement_values
SET value = NULL, value_choice = mapping.choice, updated_at = NOW() AT TIME ZONE 'UTC'
FROM measurement_entries, measurement_items, (
  VALUES
    (1, 'both_50'),
    (2, 'both_40'),
    (3, 'both_30'),
    (4, 'both_20'),
    (5, 'both_10'),
    (6, 'single_40'),
    (7, 'single_30'),
    (8, 'single_20'),
    (9, 'single_10'),
    (10, 'single_0')
) AS mapping (level, choice)
WHERE measurement_values.measurement_entry_id = measurement_entries.id
  AND measurement_entries.measurement_item_id = measurement_items.id
  AND measurement_items.code = 'stand_up_test'
  AND measurement_values.value = mapping.level;

ALTER TABLE measurement_items DROP CONSTRAINT IF EXISTS measurement_items_unit_check;
ALTER TABLE measurement_items ADD CONSTRAINT measurement_items_unit_check
  CHECK (unit IN ('kg', 'cm', 'sec', 'count', 'mmHg', 'percent', 'bpm'));
