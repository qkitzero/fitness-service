UPDATE measurement_items
SET side_mode = 'optional_bilateral', side_aggregation = 'worst', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE code = 'stand_up_test';
