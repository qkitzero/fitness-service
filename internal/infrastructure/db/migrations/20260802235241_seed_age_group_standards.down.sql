DELETE FROM age_group_standards
USING measurement_items
WHERE measurement_items.id = age_group_standards.measurement_item_id
  AND measurement_items.code IN (
    'grip_strength',
    'cs30',
    'sit_and_reach',
    'stick_reaction',
    'side_step',
    'eyes_closed_one_leg_stand',
    'eyes_open_one_leg_stand',
    'functional_reach',
    'two_step',
    'timed_up_and_go',
    'walk_5m'
  )
  AND age_group_standards.age_from IN (20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75);
