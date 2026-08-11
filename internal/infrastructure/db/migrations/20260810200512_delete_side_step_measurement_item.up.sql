DELETE FROM measurement_entries
USING measurement_items
WHERE measurement_items.id = measurement_entries.measurement_item_id
  AND measurement_items.code = 'side_step';

DELETE FROM age_group_standards
USING measurement_items
WHERE measurement_items.id = age_group_standards.measurement_item_id
  AND measurement_items.code = 'side_step';

DELETE FROM measurement_items
WHERE code = 'side_step';
