package customer

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) customer.CustomerRepository {
	return &customerRepository{db: db}
}

func toModel(c customer.Customer) CustomerModel {
	addr := c.Address()
	return CustomerModel{
		ID:                           c.ID(),
		TenantID:                     c.TenantID(),
		Name:                         c.Name(),
		NameKana:                     c.NameKana(),
		Gender:                       c.Gender(),
		BirthDate:                    c.BirthDate(),
		Phone:                        c.Phone(),
		Email:                        c.Email(),
		PostalCode:                   addr.PostalCode(),
		Prefecture:                   addr.Prefecture(),
		City:                         addr.City(),
		Street:                       addr.Street(),
		Building:                     addr.Building(),
		EmergencyContactName:         c.EmergencyContactName(),
		EmergencyContactRelationship: c.EmergencyContactRelationship(),
		EmergencyContactPhone:        c.EmergencyContactPhone(),
		OrganizationID:               c.OrganizationID(),
		IsActive:                     c.IsActive(),
		CreatedAt:                    c.CreatedAt(),
		UpdatedAt:                    c.UpdatedAt(),
	}
}

func toDomain(m CustomerModel) customer.Customer {
	return customer.NewCustomer(
		m.ID,
		m.TenantID,
		m.Name,
		m.NameKana,
		m.Gender,
		m.BirthDate,
		m.Phone,
		m.Email,
		address.NewAddress(m.PostalCode, m.Prefecture, m.City, m.Street, m.Building),
		m.EmergencyContactName,
		m.EmergencyContactRelationship,
		m.EmergencyContactPhone,
		m.OrganizationID,
		m.IsActive,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

func (r *customerRepository) Create(ctx context.Context, c customer.Customer) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		customerModel := toModel(c)

		if err := tx.Create(&customerModel).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *customerRepository) FindByID(ctx context.Context, id customer.CustomerID) (customer.Customer, error) {
	var customerModel CustomerModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&customerModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, customer.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}

	return toDomain(customerModel), nil
}

func (r *customerRepository) ListByTenantID(ctx context.Context, tenantID tenant.TenantID, includeInactive bool) ([]customer.Customer, error) {
	var customerModels []CustomerModel
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}
	if err := query.Order("created_at, id").Find(&customerModels).Error; err != nil {
		return nil, err
	}

	customers := make([]customer.Customer, 0, len(customerModels))
	for _, customerModel := range customerModels {
		customers = append(customers, toDomain(customerModel))
	}

	return customers, nil
}

func (r *customerRepository) Update(ctx context.Context, c customer.Customer) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		customerModel := toModel(c)

		result := tx.Model(&CustomerModel{}).
			Where("id = ?", customerModel.ID).
			Select(
				"tenant_id",
				"name",
				"name_kana",
				"gender",
				"birth_date",
				"phone",
				"email",
				"postal_code",
				"prefecture",
				"city",
				"street",
				"building",
				"emergency_contact_name",
				"emergency_contact_relationship",
				"emergency_contact_phone",
				"organization_id",
				"created_at",
				"updated_at",
			).
			Updates(customerModel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return customer.ErrCustomerNotFound
		}

		return nil
	})
}

func (r *customerRepository) UpdateActive(ctx context.Context, c customer.Customer) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		customerModel := toModel(c)

		result := tx.Model(&CustomerModel{}).
			Where("id = ?", customerModel.ID).
			Select(
				"is_active",
				"updated_at",
			).
			Updates(customerModel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return customer.ErrCustomerNotFound
		}

		return nil
	})
}

func (r *customerRepository) Delete(ctx context.Context, id customer.CustomerID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ?", id).Delete(&CustomerModel{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return customer.ErrCustomerNotFound
		}

		return nil
	})
}
