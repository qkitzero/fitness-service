DELETE FROM measurement_items
WHERE code IN (
  'blood_pressure',
  'pulse_rate',
  'height',
  'weight',
  'body_fat_percentage',
  'muscle_mass',
  'grip_strength',
  'stand_up_test',
  'cs30',
  'sit_and_reach',
  'stick_reaction',
  'eyes_closed_one_leg_stand',
  'eyes_open_one_leg_stand',
  'functional_reach',
  'two_step',
  'timed_up_and_go',
  'walk_5m',
  'side_step'
);
