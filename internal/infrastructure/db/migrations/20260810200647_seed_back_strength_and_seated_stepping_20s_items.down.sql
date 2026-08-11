DELETE FROM measurement_entries
USING measurement_items
WHERE measurement_items.id = measurement_entries.measurement_item_id
  AND measurement_items.code IN (
    'back_strength',
    'seated_stepping_20s'
  );

DELETE FROM age_group_standards
USING measurement_items
WHERE measurement_items.id = age_group_standards.measurement_item_id
  AND measurement_items.code IN (
    'back_strength',
    'seated_stepping_20s'
  );

DELETE FROM measurement_items
WHERE code IN (
  'back_strength',
  'seated_stepping_20s'
);
