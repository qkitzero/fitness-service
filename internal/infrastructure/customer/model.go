package customer

import (
	"time"

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
	Phone                        *customer.Phone
	Email                        *customer.Email
	PostalCode                   *customer.PostalCode
	Prefecture                   *customer.Prefecture
	City                         *customer.City
	Street                       *customer.Street
	Building                     *customer.Building
	EmergencyContactName         *customer.EmergencyContactName
	EmergencyContactRelationship *customer.EmergencyContactRelationship
	EmergencyContactPhone        *customer.Phone
	OrganizationID               *organization.OrganizationID
	IsActive                     bool
	CreatedAt                    time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt                    time.Time `gorm:"autoUpdateTime:false"`
}

func (CustomerModel) TableName() string {
	return "customers"
}
