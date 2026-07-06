package customer

import (
	"context"
	"time"

	"github.com/qkitzero/fitness-service/internal/application/auth"
	"github.com/qkitzero/fitness-service/internal/domain/customer"
)

type CustomerUsecase interface {
	CreateCustomer(ctx context.Context, name customer.Name) (customer.Customer, error)
	GetCustomer(ctx context.Context, customerID customer.CustomerID) (customer.Customer, error)
	UpdateCustomer(ctx context.Context, customerID customer.CustomerID, name customer.Name) (customer.Customer, error)
	DeleteCustomer(ctx context.Context, customerID customer.CustomerID) error
}

type customerUsecase struct {
	authService  auth.AuthService
	customerRepo customer.CustomerRepository
}

func NewCustomerUsecase(authService auth.AuthService, customerRepo customer.CustomerRepository) CustomerUsecase {
	return &customerUsecase{authService: authService, customerRepo: customerRepo}
}

func (u *customerUsecase) CreateCustomer(ctx context.Context, name customer.Name) (customer.Customer, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	now := time.Now()

	newCustomer := customer.NewCustomer(customer.NewCustomerID(), name, now, now)

	if err := u.customerRepo.Create(ctx, newCustomer); err != nil {
		return nil, err
	}

	return newCustomer, nil
}

func (u *customerUsecase) GetCustomer(ctx context.Context, customerID customer.CustomerID) (customer.Customer, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	foundCustomer, err := u.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	return foundCustomer, nil
}

func (u *customerUsecase) UpdateCustomer(ctx context.Context, customerID customer.CustomerID, name customer.Name) (customer.Customer, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

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
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return err
	}

	if _, err := u.customerRepo.FindByID(ctx, customerID); err != nil {
		return err
	}

	if err := u.customerRepo.Delete(ctx, customerID); err != nil {
		return err
	}

	return nil
}
