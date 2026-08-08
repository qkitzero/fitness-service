package judgment

import (
	"context"

	"github.com/qkitzero/fitness-service/internal/domain/measurement"
)

type PrescribedMenuOverrideRepository interface {
	ListByMeasurementID(ctx context.Context, measurementID measurement.MeasurementID) ([]PrescribedMenuOverride, error)
	ReplaceByMeasurementID(ctx context.Context, measurementID measurement.MeasurementID, overrides []PrescribedMenuOverride) error
	DeleteByMeasurementID(ctx context.Context, measurementID measurement.MeasurementID) error
}
