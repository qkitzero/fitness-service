package measurement

import "errors"

var (
	ErrMeasurementNotFound     = errors.New("measurement not found")
	ErrInvalidValueCount       = errors.New("invalid value count")
	ErrInvalidValueForType     = errors.New("invalid value for value type")
	ErrInvalidSide             = errors.New("invalid side")
	ErrDuplicateValue          = errors.New("duplicate value")
	ErrDuplicateEntry          = errors.New("duplicate entry")
	ErrInvalidAgeAtMeasurement = errors.New("invalid age at measurement")
)
