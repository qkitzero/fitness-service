package customer

import "errors"

var (
	ErrCustomerNotFound        = errors.New("customer not found")
	ErrOrganizationNotInTenant = errors.New("organization not in tenant")
	ErrCustomerInUse           = errors.New("customer in use")
)
