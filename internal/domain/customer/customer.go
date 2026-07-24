package customer

import (
	"time"
)

type Customer interface {
	ID() CustomerID
	GroupID() GroupID
	Name() Name
	NameKana() NameKana
	Gender() Gender
	BirthDate() BirthDate
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(name Name, nameKana NameKana, gender Gender, birthDate BirthDate)
}

type customer struct {
	id        CustomerID
	groupID   GroupID
	name      Name
	nameKana  NameKana
	gender    Gender
	birthDate BirthDate
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

func (c customer) NameKana() NameKana {
	return c.nameKana
}

func (c customer) Gender() Gender {
	return c.gender
}

func (c customer) BirthDate() BirthDate {
	return c.birthDate
}

func (c customer) CreatedAt() time.Time {
	return c.createdAt
}

func (c customer) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *customer) Update(name Name, nameKana NameKana, gender Gender, birthDate BirthDate) {
	c.name = name
	c.nameKana = nameKana
	c.gender = gender
	c.birthDate = birthDate
	c.updatedAt = time.Now()
}

func NewCustomer(
	id CustomerID,
	groupID GroupID,
	name Name,
	nameKana NameKana,
	gender Gender,
	birthDate BirthDate,
	createdAt time.Time,
	updatedAt time.Time,
) Customer {
	return &customer{
		id:        id,
		groupID:   groupID,
		name:      name,
		nameKana:  nameKana,
		gender:    gender,
		birthDate: birthDate,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
