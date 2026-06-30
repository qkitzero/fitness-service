package customer

import (
	"context"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
)

type CustomerUsecase interface {
	CreateCustomer(ctx context.Context, name customer.Name) (customer.Customer, error)
	GetCustomer(ctx context.Context, customerID customer.CustomerID) (customer.Customer, error)
	UpdateCustomer(ctx context.Context, customerID customer.CustomerID, name customer.Name) (customer.Customer, error)
	DeleteCustomer(ctx context.Context, customerID customer.CustomerID) error
}

type customerUsecase struct {
	customerRepo customer.CustomerRepository
}

func NewCustomerUsecase(customerRepo customer.CustomerRepository) CustomerUsecase {
	return &customerUsecase{customerRepo: customerRepo}
}

func (u *customerUsecase) CreateCustomer(ctx context.Context, name customer.Name) (customer.Customer, error) {
	now := time.Now()

	newCustomer := customer.NewCustomer(customer.NewCustomerID(), name, now, now)

	if err := u.customerRepo.Create(ctx, newCustomer); err != nil {
		return nil, err
	}

	return newCustomer, nil
}

func (u *customerUsecase) GetCustomer(ctx context.Context, customerID customer.CustomerID) (customer.Customer, error) {
	foundCustomer, err := u.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	return foundCustomer, nil
}

func (u *customerUsecase) UpdateCustomer(ctx context.Context, customerID customer.CustomerID, name customer.Name) (customer.Customer, error) {
	foundCustomer, err := u.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	foundCustomer.Update(name)

	if err := u.customerRepo.Update(ctx, foundCustomer); err != nil {
		return nil, err
	}

	return foundCustomer, nil
}

func (u *customerUsecase) DeleteCustomer(ctx context.Context, customerID customer.CustomerID) error {
	if _, err := u.customerRepo.FindByID(ctx, customerID); err != nil {
		return err
	}

	if err := u.customerRepo.Delete(ctx, customerID); err != nil {
		return err
	}

	return nil
}
