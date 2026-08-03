UPDATE measurement_items SET score_direction = 'higher_is_better', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code IN (
  'grip_strength',
  'cs30',
  'sit_and_reach',
  'eyes_closed_one_leg_stand',
  'eyes_open_one_leg_stand',
  'functional_reach',
  'two_step',
  'side_step'
);

UPDATE measurement_items SET score_direction = 'lower_is_better', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code IN (
  'stick_reaction',
  'timed_up_and_go',
  'walk_5m'
);

INSERT INTO measurement_item_elements (measurement_item_id, element, created_at, updated_at)
SELECT measurement_items.id, seed.element, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'
FROM (
  VALUES
    ('grip_strength', 'muscle_strength'),
    ('stand_up_test', 'muscle_strength'),
    ('cs30', 'muscle_endurance'),
    ('sit_and_reach', 'flexibility'),
    ('stick_reaction', 'agility'),
    ('side_step', 'agility'),
    ('eyes_closed_one_leg_stand', 'balance'),
    ('eyes_open_one_leg_stand', 'balance'),
    ('functional_reach', 'balance'),
    ('two_step', 'mobility'),
    ('timed_up_and_go', 'mobility'),
    ('walk_5m', 'mobility')
) AS seed (code, element)
JOIN measurement_items ON measurement_items.code = seed.code;

INSERT INTO rank_standards (rank, z_score_min, z_score_max, created_at, updated_at)
VALUES
  ('A', 1.50, NULL, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'),
  ('B', 0.50, 1.50, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'),
  ('C', -0.50, 0.50, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'),
  ('D', -1.50, -0.50, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC'),
  ('E', NULL, -1.50, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC');
