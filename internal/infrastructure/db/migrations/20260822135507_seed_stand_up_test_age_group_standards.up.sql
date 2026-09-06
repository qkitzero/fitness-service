INSERT INTO age_group_standards (id, measurement_item_id, gender, age_from, age_to, mean, standard_deviation, created_at, updated_at)
SELECT gen_random_uuid()::text, measurement_items.id, seed.gender, seed.age_from, seed.age_to, seed.mean, seed.standard_deviation, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'
FROM (
  VALUES
    ('stand_up_test', 'male', 20, 24, 9.00, 1.50),
    ('stand_up_test', 'male', 25, 29, 8.80, 1.50),
    ('stand_up_test', 'male', 30, 34, 8.50, 1.50),
    ('stand_up_test', 'male', 35, 39, 8.20, 1.50),
    ('stand_up_test', 'male', 40, 44, 8.00, 1.50),
    ('stand_up_test', 'male', 45, 49, 7.50, 1.50),
    ('stand_up_test', 'male', 50, 54, 7.00, 1.50),
    ('stand_up_test', 'male', 55, 59, 6.50, 1.50),
    ('stand_up_test', 'male', 60, 64, 6.00, 1.50),
    ('stand_up_test', 'male', 65, 69, 5.50, 1.50),
    ('stand_up_test', 'male', 70, 74, 5.00, 1.50),
    ('stand_up_test', 'male', 75, 79, 4.50, 1.50),
    ('stand_up_test', 'female', 20, 24, 8.70, 1.50),
    ('stand_up_test', 'female', 25, 29, 8.50, 1.50),
    ('stand_up_test', 'female', 30, 34, 8.20, 1.50),
    ('stand_up_test', 'female', 35, 39, 7.90, 1.50),
    ('stand_up_test', 'female', 40, 44, 7.60, 1.50),
    ('stand_up_test', 'female', 45, 49, 7.10, 1.50),
    ('stand_up_test', 'female', 50, 54, 6.60, 1.50),
    ('stand_up_test', 'female', 55, 59, 6.10, 1.50),
    ('stand_up_test', 'female', 60, 64, 5.60, 1.50),
    ('stand_up_test', 'female', 65, 69, 5.10, 1.50),
    ('stand_up_test', 'female', 70, 74, 4.60, 1.50),
    ('stand_up_test', 'female', 75, 79, 4.10, 1.50)
) AS seed (code, gender, age_from, age_to, mean, standard_deviation)
JOIN measurement_items ON measurement_items.code = seed.code;
