package standard

import (
	"context"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

type AgeGroupStandardRepository interface {
	List(ctx context.Context) ([]AgeGroupStandard, error)
	ListByItemIDsAndGender(ctx context.Context, measurementItemIDs []measurementitem.MeasurementItemID, gender Gender) ([]AgeGroupStandard, error)
}
