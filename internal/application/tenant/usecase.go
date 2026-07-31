package tenant

import (
	"context"
	"errors"
	"time"

	"github.com/qkitzero/fitness-service/internal/application/auth"
	"github.com/qkitzero/fitness-service/internal/application/user"
	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/contact"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type ProfilePatch struct {
	PostalCode     *address.PostalCode
	HasPostalCode  bool
	Prefecture     *address.Prefecture
	HasPrefecture  bool
	City           *address.City
	HasCity        bool
	Street         *address.Street
	HasStreet      bool
	Building       *address.Building
	HasBuilding    bool
	Phone          *contact.Phone
	HasPhone       bool
	Email          *contact.Email
	HasEmail       bool
	HomepageURL    *tenant.HomepageURL
	HasHomepageURL bool
	Note           *tenant.Note
	HasNote        bool
}

type ProfileUsecase interface {
	GetProfile(ctx context.Context, tenantID tenant.TenantID) (tenant.Profile, error)
	UpsertProfile(ctx context.Context, tenantID tenant.TenantID, patch ProfilePatch) (tenant.Profile, error)
}

type profileUsecase struct {
	authService auth.AuthService
	userService user.UserService
	profileRepo tenant.ProfileRepository
}

func NewProfileUsecase(authService auth.AuthService, userService user.UserService, profileRepo tenant.ProfileRepository) ProfileUsecase {
	return &profileUsecase{authService: authService, userService: userService, profileRepo: profileRepo}
}

func (u *profileUsecase) verifyTenantMembership(ctx context.Context, tenantID tenant.TenantID) error {
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

func (u *profileUsecase) GetProfile(ctx context.Context, tenantID tenant.TenantID) (tenant.Profile, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	if err := u.verifyTenantMembership(ctx, tenantID); err != nil {
		if errors.Is(err, tenant.ErrNotMember) {
			return nil, tenant.ErrProfileNotFound
		}
		return nil, err
	}

	foundProfile, err := u.profileRepo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	return foundProfile, nil
}

func (u *profileUsecase) UpsertProfile(ctx context.Context, tenantID tenant.TenantID, patch ProfilePatch) (tenant.Profile, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	if err := u.verifyTenantMembership(ctx, tenantID); err != nil {
		if errors.Is(err, tenant.ErrNotMember) {
			return nil, tenant.ErrProfileNotFound
		}
		return nil, err
	}

	foundProfile, err := u.profileRepo.FindByTenantID(ctx, tenantID)
	if err != nil && !errors.Is(err, tenant.ErrProfileNotFound) {
		return nil, err
	}

	targetProfile := foundProfile
	if targetProfile == nil {
		now := time.Now().UTC()
		targetProfile = tenant.NewProfile(tenantID, address.Address{}, nil, nil, nil, nil, now, now)
	}

	addr := targetProfile.Address()
	postalCode := addr.PostalCode()
	if patch.HasPostalCode {
		postalCode = patch.PostalCode
	}
	prefecture := addr.Prefecture()
	if patch.HasPrefecture {
		prefecture = patch.Prefecture
	}
	city := addr.City()
	if patch.HasCity {
		city = patch.City
	}
	street := addr.Street()
	if patch.HasStreet {
		street = patch.Street
	}
	building := addr.Building()
	if patch.HasBuilding {
		building = patch.Building
	}
	phone := targetProfile.Phone()
	if patch.HasPhone {
		phone = patch.Phone
	}
	email := targetProfile.Email()
	if patch.HasEmail {
		email = patch.Email
	}
	homepageURL := targetProfile.HomepageURL()
	if patch.HasHomepageURL {
		homepageURL = patch.HomepageURL
	}
	note := targetProfile.Note()
	if patch.HasNote {
		note = patch.Note
	}

	targetProfile.Update(address.NewAddress(postalCode, prefecture, city, street, building), phone, email, homepageURL, note)

	if err := u.profileRepo.Upsert(ctx, targetProfile); err != nil {
		return nil, err
	}

	return targetProfile, nil
}
