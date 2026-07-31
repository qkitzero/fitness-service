package customer

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/contact"
	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type CustomerModel struct {
	ID                           customer.CustomerID
	TenantID                     tenant.TenantID
	Name                         customer.Name
	NameKana                     customer.NameKana
	Gender                       customer.Gender
	BirthDate                    customer.BirthDate
	Phone                        *contact.Phone
	Email                        *contact.Email
	PostalCode                   *address.PostalCode
	Prefecture                   *address.Prefecture
	City                         *address.City
	Street                       *address.Street
	Building                     *address.Building
	EmergencyContactName         *customer.EmergencyContactName
	EmergencyContactRelationship *customer.EmergencyContactRelationship
	EmergencyContactPhone        *contact.Phone
	OrganizationID               *organization.OrganizationID
	IsActive                     bool
	CreatedAt                    time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt                    time.Time `gorm:"autoUpdateTime:false"`
}

func (CustomerModel) TableName() string {
	return "customers"
}
