INSERT INTO measurement_items (id, code, name, category, unit, trial_count, bilateral, value_type, score_direction, side_aggregation, created_at, updated_at)
VALUES
  ('c8e6df71-e55a-40bf-adf3-25fe89038eaf', 'back_strength', '背筋力', 'motor_function', 'kg', 2, FALSE, 'numeric', 'higher_is_better', 'mean', NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'),
  ('8133ca64-6a30-4a31-8114-a4b0dfdf06b3', 'seated_stepping_20s', '座位ステップ（20秒）', 'motor_function', 'count', 1, FALSE, 'numeric', 'higher_is_better', 'mean', NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC');

INSERT INTO measurement_item_elements (measurement_item_id, element, created_at, updated_at)
SELECT measurement_items.id, seed.element, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'
FROM (
  VALUES
    ('back_strength', 'muscle_strength'),
    ('seated_stepping_20s', 'agility')
) AS seed (code, element)
JOIN measurement_items ON measurement_items.code = seed.code;
