package customer

import (
	"context"
	"errors"
	"time"

	"github.com/qkitzero/fitness-service/internal/application/auth"
	"github.com/qkitzero/fitness-service/internal/application/user"
	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type CustomerUsecase interface {
	CreateCustomer(ctx context.Context, tenantID tenant.TenantID, name customer.Name, nameKana customer.NameKana, gender customer.Gender, birthDate customer.BirthDate, phone *customer.Phone, email *customer.Email, postalCode *customer.PostalCode, prefecture *customer.Prefecture, city *customer.City, street *customer.Street, building *customer.Building, emergencyContactName *customer.EmergencyContactName, emergencyContactRelationship *customer.EmergencyContactRelationship, emergencyContactPhone *customer.Phone, organizationID *organization.OrganizationID) (customer.Customer, error)
	GetCustomer(ctx context.Context, customerID customer.CustomerID) (customer.Customer, error)
	ListCustomers(ctx context.Context, tenantID tenant.TenantID, includeInactive bool) ([]customer.Customer, error)
	UpdateCustomer(ctx context.Context, customerID customer.CustomerID, name customer.Name, nameKana customer.NameKana, gender customer.Gender, birthDate customer.BirthDate, phone *customer.Phone, email *customer.Email, postalCode *customer.PostalCode, prefecture *customer.Prefecture, city *customer.City, street *customer.Street, building *customer.Building, emergencyContactName *customer.EmergencyContactName, emergencyContactRelationship *customer.EmergencyContactRelationship, emergencyContactPhone *customer.Phone, organizationID *organization.OrganizationID) (customer.Customer, error)
	SetCustomerActive(ctx context.Context, customerID customer.CustomerID, active bool) (customer.Customer, error)
	DeleteCustomer(ctx context.Context, customerID customer.CustomerID) error
}

type customerUsecase struct {
	authService      auth.AuthService
	userService      user.UserService
	customerRepo     customer.CustomerRepository
	organizationRepo organization.OrganizationRepository
}

func NewCustomerUsecase(authService auth.AuthService, userService user.UserService, customerRepo customer.CustomerRepository, organizationRepo organization.OrganizationRepository) CustomerUsecase {
	return &customerUsecase{authService: authService, userService: userService, customerRepo: customerRepo, organizationRepo: organizationRepo}
}

func (u *customerUsecase) verifyTenantMembership(ctx context.Context, tenantID tenant.TenantID) error {
	groupIDs, err := u.userService.ListMyGroups(ctx)
	if err != nil {
		return err
	}

	for _, id := range groupIDs {
		if id == tenantID.String() {
			return nil
		}
	}

	return tenant.ErrNotMember
}

func (u *customerUsecase) verifyOrganizationTenant(ctx context.Context, tenantID tenant.TenantID, organizationID *organization.OrganizationID) error {
	if organizationID == nil {
		return nil
	}

	foundOrganization, err := u.organizationRepo.FindByID(ctx, *organizationID)
	if err != nil {
		if errors.Is(err, organization.ErrOrganizationNotFound) {
			return customer.ErrOrganizationNotInTenant
		}
		return err
	}

	if foundOrganization.TenantID() != tenantID {
		return customer.ErrOrganizationNotInTenant
	}

	return nil
}

func (u *customerUsecase) findOwnedCustomer(ctx context.Context, customerID customer.CustomerID) (customer.Customer, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	foundCustomer, err := u.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	if err := u.verifyTenantMembership(ctx, foundCustomer.TenantID()); err != nil {
		if errors.Is(err, tenant.ErrNotMember) {
			return nil, customer.ErrCustomerNotFound
		}
		return nil, err
	}

	return foundCustomer, nil
}

func (u *customerUsecase) CreateCustomer(ctx context.Context, tenantID tenant.TenantID, name customer.Name, nameKana customer.NameKana, gender customer.Gender, birthDate customer.BirthDate, phone *customer.Phone, email *customer.Email, postalCode *customer.PostalCode, prefecture *customer.Prefecture, city *customer.City, street *customer.Street, building *customer.Building, emergencyContactName *customer.EmergencyContactName, emergencyContactRelationship *customer.EmergencyContactRelationship, emergencyContactPhone *customer.Phone, organizationID *organization.OrganizationID) (customer.Customer, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	if err := u.verifyTenantMembership(ctx, tenantID); err != nil {
		return nil, err
	}

	if err := u.verifyOrganizationTenant(ctx, tenantID, organizationID); err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	newCustomer := customer.NewCustomer(customer.NewCustomerID(), tenantID, name, nameKana, gender, birthDate, phone, email, postalCode, prefecture, city, street, building, emergencyContactName, emergencyContactRelationship, emergencyContactPhone, organizationID, true, now, now)

	if err := u.customerRepo.Create(ctx, newCustomer); err != nil {
		return nil, err
	}

	return newCustomer, nil
}

func (u *customerUsecase) GetCustomer(ctx context.Context, customerID customer.CustomerID) (customer.Customer, error) {
	return u.findOwnedCustomer(ctx, customerID)
}

func (u *customerUsecase) ListCustomers(ctx context.Context, tenantID tenant.TenantID, includeInactive bool) ([]customer.Customer, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	if err := u.verifyTenantMembership(ctx, tenantID); err != nil {
		return nil, err
	}

	customers, err := u.customerRepo.ListByTenantID(ctx, tenantID, includeInactive)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (u *customerUsecase) UpdateCustomer(ctx context.Context, customerID customer.CustomerID, name customer.Name, nameKana customer.NameKana, gender customer.Gender, birthDate customer.BirthDate, phone *customer.Phone, email *customer.Email, postalCode *customer.PostalCode, prefecture *customer.Prefecture, city *customer.City, street *customer.Street, building *customer.Building, emergencyContactName *customer.EmergencyContactName, emergencyContactRelationship *customer.EmergencyContactRelationship, emergencyContactPhone *customer.Phone, organizationID *organization.OrganizationID) (customer.Customer, error) {
	foundCustomer, err := u.findOwnedCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}

	if err := u.verifyOrganizationTenant(ctx, foundCustomer.TenantID(), organizationID); err != nil {
		return nil, err
	}

	foundCustomer.Update(name, nameKana, gender, birthDate, phone, email, postalCode, prefecture, city, street, building, emergencyContactName, emergencyContactRelationship, emergencyContactPhone, organizationID)

	if err := u.customerRepo.Update(ctx, foundCustomer); err != nil {
		return nil, err
	}

	return foundCustomer, nil
}

func (u *customerUsecase) SetCustomerActive(ctx context.Context, customerID customer.CustomerID, active bool) (customer.Customer, error) {
	foundCustomer, err := u.findOwnedCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}

	foundCustomer.SetActive(active)

	if err := u.customerRepo.UpdateActive(ctx, foundCustomer); err != nil {
		return nil, err
	}

	return foundCustomer, nil
}

func (u *customerUsecase) DeleteCustomer(ctx context.Context, customerID customer.CustomerID) error {
	if _, err := u.findOwnedCustomer(ctx, customerID); err != nil {
		return err
	}

	if err := u.customerRepo.Delete(ctx, customerID); err != nil {
		return err
	}

	return nil
}
