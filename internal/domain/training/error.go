package training

import "errors"

var (
	ErrTrainingMenuNotFound = errors.New("training menu not found")
	ErrInvalidSortOrder     = errors.New("invalid sort order")
)
