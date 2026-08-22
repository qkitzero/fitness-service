DELETE FROM age_group_standards
USING measurement_items
WHERE measurement_items.id = age_group_standards.measurement_item_id
  AND measurement_items.code = 'stand_up_test';
