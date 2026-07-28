package customer

import (
	"context"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer Customer) error
	FindByID(ctx context.Context, customerID CustomerID) (Customer, error)
	ListByGroupID(ctx context.Context, groupID GroupID, includeInactive bool) ([]Customer, error)
	Update(ctx context.Context, customer Customer) error
	Delete(ctx context.Context, customerID CustomerID) error
}
