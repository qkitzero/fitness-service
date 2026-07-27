package organization

import (
	"time"
)

type Organization interface {
	ID() OrganizationID
	GroupID() GroupID
	Name() Name
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(name Name)
}

type organization struct {
	id        OrganizationID
	groupID   GroupID
	name      Name
	createdAt time.Time
	updatedAt time.Time
}

func (o organization) ID() OrganizationID {
	return o.id
}

func (o organization) GroupID() GroupID {
	return o.groupID
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
	groupID GroupID,
	name Name,
	createdAt time.Time,
	updatedAt time.Time,
) Organization {
	return &organization{
		id:        id,
		groupID:   groupID,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
