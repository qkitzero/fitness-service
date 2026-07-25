package customer

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
)

type CustomerModel struct {
	ID                           customer.CustomerID
	GroupID                      customer.GroupID
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
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
}

func (CustomerModel) TableName() string {
	return "customers"
}
