ALTER TABLE measurement_items DROP CONSTRAINT IF EXISTS measurement_items_unit_check;
ALTER TABLE measurement_items ADD CONSTRAINT measurement_items_unit_check
  CHECK (unit IN ('kg', 'cm', 'sec', 'count', 'mmHg', 'percent', 'bpm', 'level'));

UPDATE measurement_values
SET value = mapping.level, value_choice = NULL, updated_at = NOW() AT TIME ZONE 'UTC'
FROM measurement_entries, measurement_items, (
  VALUES
    ('both_50', 1),
    ('both_40', 2),
    ('both_30', 3),
    ('both_20', 4),
    ('both_10', 5),
    ('single_40', 6),
    ('single_30', 7),
    ('single_20', 8),
    ('single_10', 9),
    ('single_0', 10),
    ('40', 6),
    ('30', 7),
    ('20', 8),
    ('10', 9)
) AS mapping (choice, level)
WHERE measurement_values.measurement_entry_id = measurement_entries.id
  AND measurement_entries.measurement_item_id = measurement_items.id
  AND measurement_items.code = 'stand_up_test'
  AND measurement_values.value_choice = mapping.choice;

UPDATE measurement_items
SET unit = 'level', value_type = 'numeric', score_direction = 'higher_is_better', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code = 'stand_up_test';
