package measurementitem

import (
	"context"
)

type MeasurementItemRepository interface {
	List(ctx context.Context) ([]MeasurementItem, error)
}
