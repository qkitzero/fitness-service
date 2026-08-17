UPDATE measurement_items
SET bilateral = TRUE, updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code = 'stand_up_test';

INSERT INTO measurement_values (id, measurement_entry_id, trial_index, side, value, created_at, updated_at)
SELECT gen_random_uuid()::text, folded_values.measurement_entry_id, folded_values.trial_index, 'right', folded_values.value, folded_values.created_at, NOW() AT TIME ZONE 'UTC'
FROM measurement_values AS folded_values
JOIN measurement_entries ON measurement_entries.id = folded_values.measurement_entry_id
JOIN measurement_items ON measurement_items.id = measurement_entries.measurement_item_id
WHERE measurement_items.code = 'stand_up_test'
  AND folded_values.side = 'none';

UPDATE measurement_values
SET side = 'left', updated_at = NOW() AT TIME ZONE 'UTC'
FROM measurement_entries, measurement_items
WHERE measurement_values.measurement_entry_id = measurement_entries.id
  AND measurement_entries.measurement_item_id = measurement_items.id
  AND measurement_items.code = 'stand_up_test'
  AND measurement_values.side = 'none';
