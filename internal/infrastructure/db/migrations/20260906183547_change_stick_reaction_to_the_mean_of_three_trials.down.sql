UPDATE measurement_items
SET trial_count = 5, trial_aggregation = 'best', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code = 'stick_reaction';
