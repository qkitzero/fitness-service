package measurement

import (
	"context"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
)

type MeasurementRepository interface {
	Create(ctx context.Context, measurement Measurement) error
	FindByID(ctx context.Context, measurementID MeasurementID) (Measurement, error)
	ListByCustomerID(ctx context.Context, customerID customer.CustomerID) ([]Measurement, error)
	Update(ctx context.Context, measurement Measurement) error
	Delete(ctx context.Context, measurementID MeasurementID) error
}
