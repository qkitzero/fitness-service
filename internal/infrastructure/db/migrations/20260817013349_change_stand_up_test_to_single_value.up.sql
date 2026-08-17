UPDATE measurement_values
SET side = 'none', value = folded.value, updated_at = NOW() AT TIME ZONE 'UTC'
FROM (
  SELECT MIN(side_values.id) AS id, MIN(side_values.value) AS value
  FROM measurement_values AS side_values
  JOIN measurement_entries ON measurement_entries.id = side_values.measurement_entry_id
  JOIN measurement_items ON measurement_items.id = measurement_entries.measurement_item_id
  WHERE measurement_items.code = 'stand_up_test'
    AND side_values.side <> 'none'
  GROUP BY side_values.measurement_entry_id, side_values.trial_index
) AS folded
WHERE measurement_values.id = folded.id;

DELETE FROM measurement_values
USING measurement_entries, measurement_items
WHERE measurement_values.measurement_entry_id = measurement_entries.id
  AND measurement_entries.measurement_item_id = measurement_items.id
  AND measurement_items.code = 'stand_up_test'
  AND measurement_values.side <> 'none';

UPDATE measurement_items
SET bilateral = FALSE, updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code = 'stand_up_test';
