package customer

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
)

type CustomerModel struct {
	ID        customer.CustomerID
	Name      customer.Name
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (CustomerModel) TableName() string {
	return "customers"
}
