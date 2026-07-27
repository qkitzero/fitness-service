package organization

import (
	"context"
)

type OrganizationRepository interface {
	Create(ctx context.Context, organization Organization) error
	FindByID(ctx context.Context, organizationID OrganizationID) (Organization, error)
	ListByGroupID(ctx context.Context, groupID GroupID) ([]Organization, error)
	Update(ctx context.Context, organization Organization) error
	Delete(ctx context.Context, organizationID OrganizationID) error
}
