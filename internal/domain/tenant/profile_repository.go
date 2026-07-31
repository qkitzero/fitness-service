package tenant

import (
	"context"
)

type ProfileRepository interface {
	FindByTenantID(ctx context.Context, tenantID TenantID) (Profile, error)
	Upsert(ctx context.Context, profile Profile) error
}
