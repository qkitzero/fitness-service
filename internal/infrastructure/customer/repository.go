package customer

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) customer.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, c customer.Customer) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		customerModel := CustomerModel{
			ID:        c.ID(),
			Name:      c.Name(),
			CreatedAt: c.CreatedAt(),
			UpdatedAt: c.UpdatedAt(),
		}

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

	return customer.NewCustomer(
		customerModel.ID,
		customerModel.Name,
		customerModel.CreatedAt,
		customerModel.UpdatedAt,
	), nil
}

func (r *customerRepository) Update(ctx context.Context, c customer.Customer) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		customerModel := CustomerModel{
			ID:        c.ID(),
			Name:      c.Name(),
			CreatedAt: c.CreatedAt(),
			UpdatedAt: c.UpdatedAt(),
		}

		if err := tx.Save(&customerModel).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *customerRepository) Delete(ctx context.Context, id customer.CustomerID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&CustomerModel{}).Error; err != nil {
			return err
		}

		return nil
	})
}
