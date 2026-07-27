package organization

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/organization"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) organization.OrganizationRepository {
	return &organizationRepository{db: db}
}

func toModel(o organization.Organization) OrganizationModel {
	return OrganizationModel{
		ID:        o.ID(),
		GroupID:   o.GroupID(),
		Name:      o.Name(),
		CreatedAt: o.CreatedAt(),
		UpdatedAt: o.UpdatedAt(),
	}
}

func toDomain(m OrganizationModel) organization.Organization {
	return organization.NewOrganization(
		m.ID,
		m.GroupID,
		m.Name,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

func (r *organizationRepository) Create(ctx context.Context, o organization.Organization) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		organizationModel := toModel(o)

		if err := tx.Create(&organizationModel).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *organizationRepository) FindByID(ctx context.Context, id organization.OrganizationID) (organization.Organization, error) {
	var organizationModel OrganizationModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&organizationModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, organization.ErrOrganizationNotFound
	}
	if err != nil {
		return nil, err
	}

	return toDomain(organizationModel), nil
}

func (r *organizationRepository) ListByGroupID(ctx context.Context, groupID organization.GroupID) ([]organization.Organization, error) {
	var organizationModels []OrganizationModel
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Order("created_at, id").Find(&organizationModels).Error; err != nil {
		return nil, err
	}

	organizations := make([]organization.Organization, 0, len(organizationModels))
	for _, organizationModel := range organizationModels {
		organizations = append(organizations, toDomain(organizationModel))
	}

	return organizations, nil
}

func (r *organizationRepository) Update(ctx context.Context, o organization.Organization) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		organizationModel := toModel(o)

		result := tx.Model(&OrganizationModel{}).
			Where("id = ?", organizationModel.ID).
			Select("group_id", "name", "created_at", "updated_at").
			Updates(organizationModel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return organization.ErrOrganizationNotFound
		}

		return nil
	})
}

func (r *organizationRepository) Delete(ctx context.Context, id organization.OrganizationID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ?", id).Delete(&OrganizationModel{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return organization.ErrOrganizationNotFound
		}

		return nil
	})
}
