package organization

import (
	"context"
	"errors"
	"time"

	"github.com/qkitzero/fitness-service/internal/application/auth"
	"github.com/qkitzero/fitness-service/internal/application/user"
	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type OrganizationUsecase interface {
	CreateOrganization(ctx context.Context, tenantID tenant.TenantID, name organization.Name) (organization.Organization, error)
	GetOrganization(ctx context.Context, organizationID organization.OrganizationID) (organization.Organization, error)
	ListOrganizations(ctx context.Context, tenantID tenant.TenantID) ([]organization.Organization, error)
	UpdateOrganization(ctx context.Context, organizationID organization.OrganizationID, name organization.Name) (organization.Organization, error)
	DeleteOrganization(ctx context.Context, organizationID organization.OrganizationID) error
}

type organizationUsecase struct {
	authService      auth.AuthService
	userService      user.UserService
	organizationRepo organization.OrganizationRepository
}

func NewOrganizationUsecase(authService auth.AuthService, userService user.UserService, organizationRepo organization.OrganizationRepository) OrganizationUsecase {
	return &organizationUsecase{authService: authService, userService: userService, organizationRepo: organizationRepo}
}

func (u *organizationUsecase) verifyTenantMembership(ctx context.Context, tenantID tenant.TenantID) error {
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

func (u *organizationUsecase) CreateOrganization(ctx context.Context, tenantID tenant.TenantID, name organization.Name) (organization.Organization, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	if err := u.verifyTenantMembership(ctx, tenantID); err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	newOrganization := organization.NewOrganization(organization.NewOrganizationID(), tenantID, name, now, now)

	if err := u.organizationRepo.Create(ctx, newOrganization); err != nil {
		return nil, err
	}

	return newOrganization, nil
}

func (u *organizationUsecase) GetOrganization(ctx context.Context, organizationID organization.OrganizationID) (organization.Organization, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	foundOrganization, err := u.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if err := u.verifyTenantMembership(ctx, foundOrganization.TenantID()); err != nil {
		if errors.Is(err, tenant.ErrNotMember) {
			return nil, organization.ErrOrganizationNotFound
		}
		return nil, err
	}

	return foundOrganization, nil
}

func (u *organizationUsecase) ListOrganizations(ctx context.Context, tenantID tenant.TenantID) ([]organization.Organization, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	if err := u.verifyTenantMembership(ctx, tenantID); err != nil {
		return nil, err
	}

	organizations, err := u.organizationRepo.ListByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func (u *organizationUsecase) UpdateOrganization(ctx context.Context, organizationID organization.OrganizationID, name organization.Name) (organization.Organization, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	foundOrganization, err := u.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if err := u.verifyTenantMembership(ctx, foundOrganization.TenantID()); err != nil {
		if errors.Is(err, tenant.ErrNotMember) {
			return nil, organization.ErrOrganizationNotFound
		}
		return nil, err
	}

	foundOrganization.Update(name)

	if err := u.organizationRepo.Update(ctx, foundOrganization); err != nil {
		return nil, err
	}

	return foundOrganization, nil
}

func (u *organizationUsecase) DeleteOrganization(ctx context.Context, organizationID organization.OrganizationID) error {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return err
	}

	foundOrganization, err := u.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return err
	}

	if err := u.verifyTenantMembership(ctx, foundOrganization.TenantID()); err != nil {
		if errors.Is(err, tenant.ErrNotMember) {
			return organization.ErrOrganizationNotFound
		}
		return err
	}

	if err := u.organizationRepo.Delete(ctx, organizationID); err != nil {
		return err
	}

	return nil
}
