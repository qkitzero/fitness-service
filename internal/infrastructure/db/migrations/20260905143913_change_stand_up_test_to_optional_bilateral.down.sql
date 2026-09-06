UPDATE measurement_values
SET side = 'none', value = folded.value, updated_at = NOW() AT TIME ZONE 'UTC'
FROM (
  SELECT COALESCE(MIN(measurement_values.id) FILTER (WHERE measurement_values.side = 'none'), MIN(measurement_values.id)) AS id,
    MIN(measurement_values.value) AS value
  FROM measurement_values
  JOIN measurement_entries ON measurement_entries.id = measurement_values.measurement_entry_id
  JOIN measurement_items ON measurement_items.id = measurement_entries.measurement_item_id
  WHERE measurement_items.code = 'stand_up_test'
  GROUP BY measurement_values.measurement_entry_id, measurement_values.trial_index
  HAVING COUNT(*) FILTER (WHERE measurement_values.side <> 'none') > 0
) AS folded
WHERE measurement_values.id = folded.id;

DELETE FROM measurement_values
USING measurement_entries, measurement_items
WHERE measurement_values.measurement_entry_id = measurement_entries.id
  AND measurement_entries.measurement_item_id = measurement_items.id
  AND measurement_items.code = 'stand_up_test'
  AND measurement_values.side <> 'none';

UPDATE measurement_items
SET side_mode = 'none', side_aggregation = 'mean', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code = 'stand_up_test';
