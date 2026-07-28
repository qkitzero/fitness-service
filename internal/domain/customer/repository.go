package customer

import (
	"context"

	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer Customer) error
	FindByID(ctx context.Context, customerID CustomerID) (Customer, error)
	ListByTenantID(ctx context.Context, tenantID tenant.TenantID, includeInactive bool) ([]Customer, error)
	Update(ctx context.Context, customer Customer) error
	UpdateActive(ctx context.Context, customer Customer) error
	Delete(ctx context.Context, customerID CustomerID) error
}
