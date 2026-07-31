package tenant

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) tenant.ProfileRepository {
	return &profileRepository{db: db}
}

func toModel(p tenant.Profile) ProfileModel {
	addr := p.Address()
	return ProfileModel{
		TenantID:    p.TenantID(),
		PostalCode:  addr.PostalCode(),
		Prefecture:  addr.Prefecture(),
		City:        addr.City(),
		Street:      addr.Street(),
		Building:    addr.Building(),
		Phone:       p.Phone(),
		Email:       p.Email(),
		HomepageURL: p.HomepageURL(),
		Note:        p.Note(),
		CreatedAt:   p.CreatedAt(),
		UpdatedAt:   p.UpdatedAt(),
	}
}

func toDomain(m ProfileModel) tenant.Profile {
	return tenant.NewProfile(
		m.TenantID,
		address.NewAddress(m.PostalCode, m.Prefecture, m.City, m.Street, m.Building),
		m.Phone,
		m.Email,
		m.HomepageURL,
		m.Note,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

func (r *profileRepository) FindByTenantID(ctx context.Context, tenantID tenant.TenantID) (tenant.Profile, error) {
	var profileModel ProfileModel
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).First(&profileModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, tenant.ErrProfileNotFound
	}
	if err != nil {
		return nil, err
	}

	return toDomain(profileModel), nil
}

func (r *profileRepository) Upsert(ctx context.Context, p tenant.Profile) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		profileModel := toModel(p)

		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "tenant_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"postal_code",
				"prefecture",
				"city",
				"street",
				"building",
				"phone",
				"email",
				"homepage_url",
				"note",
				"updated_at",
			}),
		}).Create(&profileModel).Error; err != nil {
			return err
		}

		return nil
	})
}
