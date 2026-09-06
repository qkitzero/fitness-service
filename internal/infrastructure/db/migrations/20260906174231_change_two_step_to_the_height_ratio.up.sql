ALTER TABLE measurement_items ADD COLUMN normalization VARCHAR(16) NOT NULL DEFAULT 'none'
  CHECK (normalization IN ('none', 'height_ratio'));

UPDATE measurement_items SET normalization = 'height_ratio', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code = 'two_step';

UPDATE age_group_standards
SET mean = seed.mean, standard_deviation = seed.standard_deviation, updated_at = NOW() AT TIME ZONE 'UTC'
FROM measurement_items, (
  VALUES
    ('male', 20, 1.61, 0.13),
    ('male', 25, 1.60, 0.13),
    ('male', 30, 1.56, 0.14),
    ('male', 35, 1.54, 0.14),
    ('male', 40, 1.52, 0.14),
    ('male', 45, 1.49, 0.14),
    ('male', 50, 1.47, 0.15),
    ('male', 55, 1.43, 0.15),
    ('male', 60, 1.41, 0.16),
    ('male', 65, 1.37, 0.16),
    ('male', 70, 1.32, 0.18),
    ('male', 75, 1.28, 0.18),
    ('female', 20, 1.62, 0.13),
    ('female', 25, 1.60, 0.13),
    ('female', 30, 1.57, 0.13),
    ('female', 35, 1.54, 0.13),
    ('female', 40, 1.51, 0.14),
    ('female', 45, 1.48, 0.14),
    ('female', 50, 1.46, 0.15),
    ('female', 55, 1.42, 0.15),
    ('female', 60, 1.40, 0.16),
    ('female', 65, 1.36, 0.16),
    ('female', 70, 1.31, 0.18),
    ('female', 75, 1.27, 0.18)
) AS seed (gender, age_from, mean, standard_deviation)
WHERE age_group_standards.measurement_item_id = measurement_items.id
  AND measurement_items.code = 'two_step'
  AND age_group_standards.gender = seed.gender
  AND age_group_standards.age_from = seed.age_from;
