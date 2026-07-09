package customer

import (
	"time"
)

type Customer interface {
	ID() CustomerID
	Name() Name
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(name Name)
}

type customer struct {
	id        CustomerID
	name      Name
	createdAt time.Time
	updatedAt time.Time
}

func (c customer) ID() CustomerID {
	return c.id
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
	name Name,
	createdAt time.Time,
	updatedAt time.Time,
) Customer {
	return &customer{
		id:        id,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
