package judgment

import (
	"context"

	"github.com/qkitzero/fitness-service/internal/domain/measurement"
)

type JudgmentRepository interface {
	FindByMeasurementID(ctx context.Context, measurementID measurement.MeasurementID) (Judgment, error)
	Upsert(ctx context.Context, judgment Judgment) error
}
