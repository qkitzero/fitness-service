DELETE FROM rank_standards
WHERE rank IN ('A', 'B', 'C', 'D', 'E');

DELETE FROM measurement_item_elements
USING (
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
) AS seed (code, element), measurement_items
WHERE measurement_items.code = seed.code
  AND measurement_item_elements.measurement_item_id = measurement_items.id
  AND measurement_item_elements.element = seed.element;

UPDATE measurement_items SET score_direction = NULL, updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code IN (
  'grip_strength',
  'cs30',
  'sit_and_reach',
  'eyes_closed_one_leg_stand',
  'eyes_open_one_leg_stand',
  'functional_reach',
  'two_step',
  'side_step',
  'stick_reaction',
  'timed_up_and_go',
  'walk_5m'
);
