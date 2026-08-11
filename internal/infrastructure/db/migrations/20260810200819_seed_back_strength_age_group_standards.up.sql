INSERT INTO age_group_standards (id, measurement_item_id, gender, age_from, age_to, mean, standard_deviation, created_at, updated_at)
SELECT gen_random_uuid()::text, measurement_items.id, seed.gender, seed.age_from, seed.age_to, seed.mean, seed.standard_deviation, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'
FROM (
  VALUES
    ('back_strength', 'male', 20, 24, 160.00, 26.00),
    ('back_strength', 'male', 25, 29, 157.00, 25.50),
    ('back_strength', 'male', 30, 34, 155.00, 25.00),
    ('back_strength', 'male', 35, 39, 150.00, 24.50),
    ('back_strength', 'male', 40, 44, 145.00, 24.00),
    ('back_strength', 'male', 45, 49, 141.00, 23.00),
    ('back_strength', 'male', 50, 54, 133.00, 22.00),
    ('back_strength', 'male', 55, 59, 124.00, 21.00),
    ('back_strength', 'male', 60, 64, 116.00, 20.00),
    ('back_strength', 'male', 65, 69, 108.00, 19.00),
    ('back_strength', 'male', 70, 74, 100.00, 18.00),
    ('back_strength', 'male', 75, 79, 93.00, 17.00),
    ('back_strength', 'female', 20, 24, 95.00, 16.00),
    ('back_strength', 'female', 25, 29, 91.00, 15.50),
    ('back_strength', 'female', 30, 34, 85.00, 14.50),
    ('back_strength', 'female', 35, 39, 82.00, 14.00),
    ('back_strength', 'female', 40, 44, 80.00, 13.50),
    ('back_strength', 'female', 45, 49, 74.00, 12.50),
    ('back_strength', 'female', 50, 54, 66.00, 11.50),
    ('back_strength', 'female', 55, 59, 61.00, 11.00),
    ('back_strength', 'female', 60, 64, 57.00, 10.50),
    ('back_strength', 'female', 65, 69, 53.00, 10.00),
    ('back_strength', 'female', 70, 74, 49.00, 9.50),
    ('back_strength', 'female', 75, 79, 45.00, 9.00)
) AS seed (code, gender, age_from, age_to, mean, standard_deviation)
JOIN measurement_items ON measurement_items.code = seed.code;
