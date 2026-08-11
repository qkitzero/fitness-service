INSERT INTO measurement_items (id, code, name, category, unit, trial_count, bilateral, value_type, score_direction, side_aggregation, created_at, updated_at)
VALUES
  ('32ea7f00-d3ee-4197-ac0e-9373a033b69e', 'side_step', '反復横跳び', 'motor_function', 'count', 1, FALSE, 'numeric', 'higher_is_better', 'mean', NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC');

INSERT INTO measurement_item_elements (measurement_item_id, element, created_at, updated_at)
VALUES
  ('32ea7f00-d3ee-4197-ac0e-9373a033b69e', 'agility', NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC');

INSERT INTO age_group_standards (id, measurement_item_id, gender, age_from, age_to, mean, standard_deviation, created_at, updated_at)
SELECT gen_random_uuid()::text, measurement_items.id, seed.gender, seed.age_from, seed.age_to, seed.mean, seed.standard_deviation, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'
FROM (
  VALUES
    ('side_step', 'male', 20, 24, 55.00, 7.00),
    ('side_step', 'male', 25, 29, 53.00, 7.00),
    ('side_step', 'male', 30, 34, 51.00, 7.00),
    ('side_step', 'male', 35, 39, 49.00, 6.80),
    ('side_step', 'male', 40, 44, 47.00, 6.80),
    ('side_step', 'male', 45, 49, 45.00, 6.50),
    ('side_step', 'male', 50, 54, 43.00, 6.50),
    ('side_step', 'male', 55, 59, 41.00, 6.30),
    ('side_step', 'male', 60, 64, 39.00, 6.30),
    ('side_step', 'male', 65, 69, 37.00, 6.00),
    ('side_step', 'male', 70, 74, 35.00, 6.00),
    ('side_step', 'male', 75, 79, 33.00, 6.00),
    ('side_step', 'female', 20, 24, 47.00, 6.00),
    ('side_step', 'female', 25, 29, 45.00, 6.00),
    ('side_step', 'female', 30, 34, 44.00, 6.00),
    ('side_step', 'female', 35, 39, 43.00, 5.80),
    ('side_step', 'female', 40, 44, 42.00, 5.80),
    ('side_step', 'female', 45, 49, 41.00, 5.80),
    ('side_step', 'female', 50, 54, 39.00, 5.50),
    ('side_step', 'female', 55, 59, 38.00, 5.50),
    ('side_step', 'female', 60, 64, 36.00, 5.50),
    ('side_step', 'female', 65, 69, 34.00, 5.30),
    ('side_step', 'female', 70, 74, 32.00, 5.30),
    ('side_step', 'female', 75, 79, 30.00, 5.30)
) AS seed (code, gender, age_from, age_to, mean, standard_deviation)
JOIN measurement_items ON measurement_items.code = seed.code;
