DELETE FROM measurement_values
USING measurement_entries, measurement_items
WHERE measurement_values.measurement_entry_id = measurement_entries.id
  AND measurement_entries.measurement_item_id = measurement_items.id
  AND measurement_items.code = 'stick_reaction'
  AND measurement_values.trial_index > 3;

DELETE FROM measurement_entries
USING measurements, measurement_items
WHERE measurement_entries.measurement_id = measurements.id
  AND measurement_entries.measurement_item_id = measurement_items.id
  AND measurement_items.code = 'stick_reaction'
  AND measurements.is_draft = FALSE
  AND measurement_entries.unmeasurable = FALSE
  AND NOT EXISTS (
    SELECT 1 FROM measurement_values
    WHERE measurement_values.measurement_entry_id = measurement_entries.id
  );

UPDATE measurement_items
SET trial_count = 3, trial_aggregation = 'mean', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code = 'stick_reaction';
