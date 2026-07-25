package customer

import (
	"context"
	"time"

	"github.com/qkitzero/fitness-service/internal/application/auth"
	"github.com/qkitzero/fitness-service/internal/application/user"
	"github.com/qkitzero/fitness-service/internal/domain/customer"
)

type CustomerUsecase interface {
	CreateCustomer(ctx context.Context, groupID customer.GroupID, name customer.Name, nameKana customer.NameKana, gender customer.Gender, birthDate customer.BirthDate, phone *customer.Phone, email *customer.Email, postalCode *customer.PostalCode, prefecture *customer.Prefecture, city *customer.City, street *customer.Street, building *customer.Building, emergencyContactName *customer.EmergencyContactName, emergencyContactRelationship *customer.EmergencyContactRelationship, emergencyContactPhone *customer.Phone) (customer.Customer, error)
	GetCustomer(ctx context.Context, customerID customer.CustomerID) (customer.Customer, error)
	ListCustomers(ctx context.Context, groupID customer.GroupID) ([]customer.Customer, error)
	UpdateCustomer(ctx context.Context, customerID customer.CustomerID, name customer.Name, nameKana customer.NameKana, gender customer.Gender, birthDate customer.BirthDate, phone *customer.Phone, email *customer.Email, postalCode *customer.PostalCode, prefecture *customer.Prefecture, city *customer.City, street *customer.Street, building *customer.Building, emergencyContactName *customer.EmergencyContactName, emergencyContactRelationship *customer.EmergencyContactRelationship, emergencyContactPhone *customer.Phone) (customer.Customer, error)
	DeleteCustomer(ctx context.Context, customerID customer.CustomerID) error
}

type customerUsecase struct {
	authService  auth.AuthService
	userService  user.UserService
	customerRepo customer.CustomerRepository
}

func NewCustomerUsecase(authService auth.AuthService, userService user.UserService, customerRepo customer.CustomerRepository) CustomerUsecase {
	return &customerUsecase{authService: authService, userService: userService, customerRepo: customerRepo}
}

func (u *customerUsecase) verifyGroupMembership(ctx context.Context, groupID customer.GroupID) error {
	groupIDs, err := u.userService.ListMyGroups(ctx)
	if err != nil {
		return err
	}

	for _, id := range groupIDs {
		if id == groupID.String() {
			return nil
		}
	}

	return user.ErrNotGroupMember
}

func (u *customerUsecase) CreateCustomer(ctx context.Context, groupID customer.GroupID, name customer.Name, nameKana customer.NameKana, gender customer.Gender, birthDate customer.BirthDate, phone *customer.Phone, email *customer.Email, postalCode *customer.PostalCode, prefecture *customer.Prefecture, city *customer.City, street *customer.Street, building *customer.Building, emergencyContactName *customer.EmergencyContactName, emergencyContactRelationship *customer.EmergencyContactRelationship, emergencyContactPhone *customer.Phone) (customer.Customer, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	if err := u.verifyGroupMembership(ctx, groupID); err != nil {
		return nil, err
	}

	now := time.Now()

	newCustomer := customer.NewCustomer(customer.NewCustomerID(), groupID, name, nameKana, gender, birthDate, phone, email, postalCode, prefecture, city, street, building, emergencyContactName, emergencyContactRelationship, emergencyContactPhone, now, now)

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

	if err := u.verifyGroupMembership(ctx, foundCustomer.GroupID()); err != nil {
		return nil, err
	}

	return foundCustomer, nil
}

func (u *customerUsecase) ListCustomers(ctx context.Context, groupID customer.GroupID) ([]customer.Customer, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	if err := u.verifyGroupMembership(ctx, groupID); err != nil {
		return nil, err
	}

	customers, err := u.customerRepo.ListByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (u *customerUsecase) UpdateCustomer(ctx context.Context, customerID customer.CustomerID, name customer.Name, nameKana customer.NameKana, gender customer.Gender, birthDate customer.BirthDate, phone *customer.Phone, email *customer.Email, postalCode *customer.PostalCode, prefecture *customer.Prefecture, city *customer.City, street *customer.Street, building *customer.Building, emergencyContactName *customer.EmergencyContactName, emergencyContactRelationship *customer.EmergencyContactRelationship, emergencyContactPhone *customer.Phone) (customer.Customer, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	foundCustomer, err := u.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	if err := u.verifyGroupMembership(ctx, foundCustomer.GroupID()); err != nil {
		return nil, err
	}

	foundCustomer.Update(name, nameKana, gender, birthDate, phone, email, postalCode, prefecture, city, street, building, emergencyContactName, emergencyContactRelationship, emergencyContactPhone)

	if err := u.customerRepo.Update(ctx, foundCustomer); err != nil {
		return nil, err
	}

	return foundCustomer, nil
}

func (u *customerUsecase) DeleteCustomer(ctx context.Context, customerID customer.CustomerID) error {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return err
	}

	foundCustomer, err := u.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}

	if err := u.verifyGroupMembership(ctx, foundCustomer.GroupID()); err != nil {
		return err
	}

	if err := u.customerRepo.Delete(ctx, customerID); err != nil {
		return err
	}

	return nil
}
