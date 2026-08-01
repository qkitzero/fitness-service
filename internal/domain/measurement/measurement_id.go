package measurement

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type MeasurementID struct {
	uuid.UUID
}

func NewMeasurementID() MeasurementID {
	id := uuid.New()
	return MeasurementID{id}
}

func NewMeasurementIDFromString(s string) (MeasurementID, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return MeasurementID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return MeasurementID{id}, nil
}
