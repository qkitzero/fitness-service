package organization

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type Organization interface {
	ID() OrganizationID
	TenantID() tenant.TenantID
	Name() Name
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(name Name)
}

type organization struct {
	id        OrganizationID
	tenantID  tenant.TenantID
	name      Name
	createdAt time.Time
	updatedAt time.Time
}

func (o organization) ID() OrganizationID {
	return o.id
}

func (o organization) TenantID() tenant.TenantID {
	return o.tenantID
}

func (o organization) Name() Name {
	return o.name
}

func (o organization) CreatedAt() time.Time {
	return o.createdAt
}

func (o organization) UpdatedAt() time.Time {
	return o.updatedAt
}

func (o *organization) Update(name Name) {
	o.name = name
	o.updatedAt = time.Now().UTC()
}

func NewOrganization(
	id OrganizationID,
	tenantID tenant.TenantID,
	name Name,
	createdAt time.Time,
	updatedAt time.Time,
) Organization {
	return &organization{
		id:        id,
		tenantID:  tenantID,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
