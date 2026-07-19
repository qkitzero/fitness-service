package customer

import (
	"time"
)

type Customer interface {
	ID() CustomerID
	GroupID() GroupID
	Name() Name
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(name Name)
}

type customer struct {
	id        CustomerID
	groupID   GroupID
	name      Name
	createdAt time.Time
	updatedAt time.Time
}

func (c customer) ID() CustomerID {
	return c.id
}

func (c customer) GroupID() GroupID {
	return c.groupID
}

func (c customer) Name() Name {
	return c.name
}

func (c customer) CreatedAt() time.Time {
	return c.createdAt
}

func (c customer) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *customer) Update(name Name) {
	c.name = name
	c.updatedAt = time.Now()
}

func NewCustomer(
	id CustomerID,
	groupID GroupID,
	name Name,
	createdAt time.Time,
	updatedAt time.Time,
) Customer {
	return &customer{
		id:        id,
		groupID:   groupID,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
