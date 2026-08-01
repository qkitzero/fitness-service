package measurementitem

import (
	"github.com/google/uuid"
)

type MeasurementItemID struct {
	uuid.UUID
}

func NewMeasurementItemID() MeasurementItemID {
	id := uuid.New()
	return MeasurementItemID{id}
}
