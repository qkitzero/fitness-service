package customer

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/contact"
	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type Customer interface {
	ID() CustomerID
	TenantID() tenant.TenantID
	Name() Name
	NameKana() NameKana
	Gender() Gender
	BirthDate() BirthDate
	Phone() *contact.Phone
	Email() *contact.Email
	Address() address.Address
	EmergencyContactName() *EmergencyContactName
	EmergencyContactRelationship() *EmergencyContactRelationship
	EmergencyContactPhone() *contact.Phone
	OrganizationID() *organization.OrganizationID
	IsActive() bool
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(name Name, nameKana NameKana, gender Gender, birthDate BirthDate, phone *contact.Phone, email *contact.Email, addr address.Address, emergencyContactName *EmergencyContactName, emergencyContactRelationship *EmergencyContactRelationship, emergencyContactPhone *contact.Phone, organizationID *organization.OrganizationID)
	SetActive(active bool)
}

type customer struct {
	id                           CustomerID
	tenantID                     tenant.TenantID
	name                         Name
	nameKana                     NameKana
	gender                       Gender
	birthDate                    BirthDate
	phone                        *contact.Phone
	email                        *contact.Email
	address                      address.Address
	emergencyContactName         *EmergencyContactName
	emergencyContactRelationship *EmergencyContactRelationship
	emergencyContactPhone        *contact.Phone
	organizationID               *organization.OrganizationID
	active                       bool
	createdAt                    time.Time
	updatedAt                    time.Time
}

func (c customer) ID() CustomerID {
	return c.id
}

func (c customer) TenantID() tenant.TenantID {
	return c.tenantID
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

func (c customer) Phone() *contact.Phone {
	if c.phone == nil {
		return nil
	}
	p := *c.phone
	return &p
}

func (c customer) Email() *contact.Email {
	if c.email == nil {
		return nil
	}
	e := *c.email
	return &e
}

func (c customer) Address() address.Address {
	return c.address
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

func (c customer) EmergencyContactPhone() *contact.Phone {
	if c.emergencyContactPhone == nil {
		return nil
	}
	p := *c.emergencyContactPhone
	return &p
}

func (c customer) OrganizationID() *organization.OrganizationID {
	if c.organizationID == nil {
		return nil
	}
	o := *c.organizationID
	return &o
}

func (c customer) IsActive() bool {
	return c.active
}

func (c customer) CreatedAt() time.Time {
	return c.createdAt
}

func (c customer) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *customer) Update(name Name, nameKana NameKana, gender Gender, birthDate BirthDate, phone *contact.Phone, email *contact.Email, addr address.Address, emergencyContactName *EmergencyContactName, emergencyContactRelationship *EmergencyContactRelationship, emergencyContactPhone *contact.Phone, organizationID *organization.OrganizationID) {
	c.name = name
	c.nameKana = nameKana
	c.gender = gender
	c.birthDate = birthDate
	c.phone = phone
	c.email = email
	c.address = addr
	c.emergencyContactName = emergencyContactName
	c.emergencyContactRelationship = emergencyContactRelationship
	c.emergencyContactPhone = emergencyContactPhone
	c.organizationID = organizationID
	c.updatedAt = time.Now().UTC()
}

func (c *customer) SetActive(active bool) {
	if c.active == active {
		return
	}
	c.active = active
	c.updatedAt = time.Now().UTC()
}

func NewCustomer(
	id CustomerID,
	tenantID tenant.TenantID,
	name Name,
	nameKana NameKana,
	gender Gender,
	birthDate BirthDate,
	phone *contact.Phone,
	email *contact.Email,
	addr address.Address,
	emergencyContactName *EmergencyContactName,
	emergencyContactRelationship *EmergencyContactRelationship,
	emergencyContactPhone *contact.Phone,
	organizationID *organization.OrganizationID,
	active bool,
	createdAt time.Time,
	updatedAt time.Time,
) Customer {
	return &customer{
		id:                           id,
		tenantID:                     tenantID,
		name:                         name,
		nameKana:                     nameKana,
		gender:                       gender,
		birthDate:                    birthDate,
		phone:                        phone,
		email:                        email,
		address:                      addr,
		emergencyContactName:         emergencyContactName,
		emergencyContactRelationship: emergencyContactRelationship,
		emergencyContactPhone:        emergencyContactPhone,
		organizationID:               organizationID,
		active:                       active,
		createdAt:                    createdAt,
		updatedAt:                    updatedAt,
	}
}
