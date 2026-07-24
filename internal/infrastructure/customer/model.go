package customer

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
)

type CustomerModel struct {
	ID        customer.CustomerID
	GroupID   customer.GroupID
	Name      customer.Name
	NameKana  customer.NameKana
	Gender    customer.Gender
	BirthDate customer.BirthDate
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (CustomerModel) TableName() string {
	return "customers"
}
