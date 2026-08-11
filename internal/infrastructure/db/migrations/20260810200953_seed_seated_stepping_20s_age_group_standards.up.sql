INSERT INTO age_group_standards (id, measurement_item_id, gender, age_from, age_to, mean, standard_deviation, created_at, updated_at)
SELECT gen_random_uuid()::text, measurement_items.id, seed.gender, seed.age_from, seed.age_to, seed.mean, seed.standard_deviation, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'
FROM (
  VALUES
    ('seated_stepping_20s', 'male', 20, 24, 43.00, 6.50),
    ('seated_stepping_20s', 'male', 25, 29, 42.00, 6.50),
    ('seated_stepping_20s', 'male', 30, 34, 41.00, 6.50),
    ('seated_stepping_20s', 'male', 35, 39, 39.00, 6.30),
    ('seated_stepping_20s', 'male', 40, 44, 38.00, 6.30),
    ('seated_stepping_20s', 'male', 45, 49, 37.00, 6.30),
    ('seated_stepping_20s', 'male', 50, 54, 35.00, 6.00),
    ('seated_stepping_20s', 'male', 55, 59, 34.00, 6.00),
    ('seated_stepping_20s', 'male', 60, 64, 32.00, 6.00),
    ('seated_stepping_20s', 'male', 65, 69, 30.00, 5.80),
    ('seated_stepping_20s', 'male', 70, 74, 28.00, 5.80),
    ('seated_stepping_20s', 'male', 75, 79, 26.00, 5.80),
    ('seated_stepping_20s', 'female', 20, 24, 41.00, 6.30),
    ('seated_stepping_20s', 'female', 25, 29, 40.00, 6.30),
    ('seated_stepping_20s', 'female', 30, 34, 39.00, 6.30),
    ('seated_stepping_20s', 'female', 35, 39, 37.00, 6.00),
    ('seated_stepping_20s', 'female', 40, 44, 36.00, 6.00),
    ('seated_stepping_20s', 'female', 45, 49, 35.00, 6.00),
    ('seated_stepping_20s', 'female', 50, 54, 33.00, 5.80),
    ('seated_stepping_20s', 'female', 55, 59, 32.00, 5.80),
    ('seated_stepping_20s', 'female', 60, 64, 30.00, 5.80),
    ('seated_stepping_20s', 'female', 65, 69, 28.00, 5.50),
    ('seated_stepping_20s', 'female', 70, 74, 26.00, 5.50),
    ('seated_stepping_20s', 'female', 75, 79, 24.00, 5.50)
) AS seed (code, gender, age_from, age_to, mean, standard_deviation)
JOIN measurement_items ON measurement_items.code = seed.code;
