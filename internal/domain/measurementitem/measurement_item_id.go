package measurementitem

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type MeasurementItemID struct {
	uuid.UUID
}

func NewMeasurementItemID() MeasurementItemID {
	id := uuid.New()
	return MeasurementItemID{id}
}

func NewMeasurementItemIDFromString(s string) (MeasurementItemID, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return MeasurementItemID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return MeasurementItemID{id}, nil
}
