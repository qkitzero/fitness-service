UPDATE age_group_standards
SET mean = seed.mean, standard_deviation = seed.standard_deviation, updated_at = NOW() AT TIME ZONE 'UTC'
FROM measurement_items, (
  VALUES
    ('male', 20, 275.00, 25.00),
    ('male', 25, 272.00, 25.00),
    ('male', 30, 268.00, 25.00),
    ('male', 35, 264.00, 25.00),
    ('male', 40, 260.00, 26.00),
    ('male', 45, 255.00, 26.00),
    ('male', 50, 250.00, 27.00),
    ('male', 55, 244.00, 27.00),
    ('male', 60, 238.00, 28.00),
    ('male', 65, 230.00, 28.00),
    ('male', 70, 220.00, 30.00),
    ('male', 75, 210.00, 30.00),
    ('female', 20, 255.00, 23.00),
    ('female', 25, 252.00, 23.00),
    ('female', 30, 248.00, 23.00),
    ('female', 35, 244.00, 23.00),
    ('female', 40, 240.00, 24.00),
    ('female', 45, 235.00, 24.00),
    ('female', 50, 230.00, 25.00),
    ('female', 55, 224.00, 25.00),
    ('female', 60, 218.00, 26.00),
    ('female', 65, 210.00, 26.00),
    ('female', 70, 200.00, 28.00),
    ('female', 75, 190.00, 28.00)
) AS seed (gender, age_from, mean, standard_deviation)
WHERE age_group_standards.measurement_item_id = measurement_items.id
  AND measurement_items.code = 'two_step'
  AND age_group_standards.gender = seed.gender
  AND age_group_standards.age_from = seed.age_from;

ALTER TABLE measurement_items DROP COLUMN normalization;
