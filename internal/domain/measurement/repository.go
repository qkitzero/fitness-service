package measurement

import (
	"context"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
)

type ListCursor struct {
	CustomerID    customer.CustomerID
	MeasuredOn    MeasuredOn
	MeasurementID MeasurementID
}

type MeasurementRepository interface {
	Create(ctx context.Context, measurement Measurement) error
	FindByID(ctx context.Context, measurementID MeasurementID) (Measurement, error)
	ListByCustomerID(ctx context.Context, customerID customer.CustomerID) ([]Measurement, error)
	ListPageByCustomerIDs(ctx context.Context, customerIDs []customer.CustomerID, after *ListCursor, limit int) ([]Measurement, error)
	Update(ctx context.Context, measurement Measurement) error
	Delete(ctx context.Context, measurementID MeasurementID) error
}
