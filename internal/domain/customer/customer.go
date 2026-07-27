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
	Phone() *Phone
	Email() *Email
	PostalCode() *PostalCode
	Prefecture() *Prefecture
	City() *City
	Street() *Street
	Building() *Building
	EmergencyContactName() *EmergencyContactName
	EmergencyContactRelationship() *EmergencyContactRelationship
	EmergencyContactPhone() *Phone
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(name Name, nameKana NameKana, gender Gender, birthDate BirthDate, phone *Phone, email *Email, postalCode *PostalCode, prefecture *Prefecture, city *City, street *Street, building *Building, emergencyContactName *EmergencyContactName, emergencyContactRelationship *EmergencyContactRelationship, emergencyContactPhone *Phone)
}

type customer struct {
	id                           CustomerID
	groupID                      GroupID
	name                         Name
	nameKana                     NameKana
	gender                       Gender
	birthDate                    BirthDate
	phone                        *Phone
	email                        *Email
	postalCode                   *PostalCode
	prefecture                   *Prefecture
	city                         *City
	street                       *Street
	building                     *Building
	emergencyContactName         *EmergencyContactName
	emergencyContactRelationship *EmergencyContactRelationship
	emergencyContactPhone        *Phone
	createdAt                    time.Time
	updatedAt                    time.Time
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

func (c customer) Phone() *Phone {
	if c.phone == nil {
		return nil
	}
	p := *c.phone
	return &p
}

func (c customer) Email() *Email {
	if c.email == nil {
		return nil
	}
	e := *c.email
	return &e
}

func (c customer) PostalCode() *PostalCode {
	if c.postalCode == nil {
		return nil
	}
	pc := *c.postalCode
	return &pc
}

func (c customer) Prefecture() *Prefecture {
	if c.prefecture == nil {
		return nil
	}
	p := *c.prefecture
	return &p
}

func (c customer) City() *City {
	if c.city == nil {
		return nil
	}
	city := *c.city
	return &city
}

func (c customer) Street() *Street {
	if c.street == nil {
		return nil
	}
	s := *c.street
	return &s
}

func (c customer) Building() *Building {
	if c.building == nil {
		return nil
	}
	b := *c.building
	return &b
}

func (c customer) EmergencyContactName() *EmergencyContactName {
	if c.emergencyContactName == nil {
		return nil
	}
	n := *c.emergencyContactName
	return &n
}

func (c customer) EmergencyContactRelationship() *EmergencyContactRelationship {
	if c.emergencyContactRelationship == nil {
		return nil
	}
	r := *c.emergencyContactRelationship
	return &r
}

func (c customer) EmergencyContactPhone() *Phone {
	if c.emergencyContactPhone == nil {
		return nil
	}
	p := *c.emergencyContactPhone
	return &p
}

func (c customer) CreatedAt() time.Time {
	return c.createdAt
}

func (c customer) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *customer) Update(name Name, nameKana NameKana, gender Gender, birthDate BirthDate, phone *Phone, email *Email, postalCode *PostalCode, prefecture *Prefecture, city *City, street *Street, building *Building, emergencyContactName *EmergencyContactName, emergencyContactRelationship *EmergencyContactRelationship, emergencyContactPhone *Phone) {
	c.name = name
	c.nameKana = nameKana
	c.gender = gender
	c.birthDate = birthDate
	c.phone = phone
	c.email = email
	c.postalCode = postalCode
	c.prefecture = prefecture
	c.city = city
	c.street = street
	c.building = building
	c.emergencyContactName = emergencyContactName
	c.emergencyContactRelationship = emergencyContactRelationship
	c.emergencyContactPhone = emergencyContactPhone
	c.updatedAt = time.Now().UTC()
}

func NewCustomer(
	id CustomerID,
	groupID GroupID,
	name Name,
	nameKana NameKana,
	gender Gender,
	birthDate BirthDate,
	phone *Phone,
	email *Email,
	postalCode *PostalCode,
	prefecture *Prefecture,
	city *City,
	street *Street,
	building *Building,
	emergencyContactName *EmergencyContactName,
	emergencyContactRelationship *EmergencyContactRelationship,
	emergencyContactPhone *Phone,
	createdAt time.Time,
	updatedAt time.Time,
) Customer {
	return &customer{
		id:                           id,
		groupID:                      groupID,
		name:                         name,
		nameKana:                     nameKana,
		gender:                       gender,
		birthDate:                    birthDate,
		phone:                        phone,
		email:                        email,
		postalCode:                   postalCode,
		prefecture:                   prefecture,
		city:                         city,
		street:                       street,
		building:                     building,
		emergencyContactName:         emergencyContactName,
		emergencyContactRelationship: emergencyContactRelationship,
		emergencyContactPhone:        emergencyContactPhone,
		createdAt:                    createdAt,
		updatedAt:                    updatedAt,
	}
}
