package organization

import (
	"context"

	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type OrganizationRepository interface {
	Create(ctx context.Context, organization Organization) error
	FindByID(ctx context.Context, organizationID OrganizationID) (Organization, error)
	ListByTenantID(ctx context.Context, tenantID tenant.TenantID) ([]Organization, error)
	Update(ctx context.Context, organization Organization) error
	Delete(ctx context.Context, organizationID OrganizationID) error
}
