package customer

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
	mockscustomer "github.com/qkitzero/fitness-service/mocks/domain/customer"
)

func TestCreateCustomer(t *testing.T) {
	t.Parallel()
	name, _ := customer.NewName("test customer")

	tests := []struct {
		name      string
		success   bool
		ctx       context.Context
		createErr error
	}{
		{"success create customer", true, context.Background(), nil},
		{"failure create error", false, context.Background(), errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockCustomerRepository.EXPECT().Create(tt.ctx, gomock.Any()).Return(tt.createErr).AnyTimes()

			u := NewCustomerUsecase(mockCustomerRepository)

			_, err := u.CreateCustomer(tt.ctx, name)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestGetCustomer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		findByIDErr error
	}{
		{"success get customer", true, context.Background(), nil},
		{"failure find by id error", false, context.Background(), errors.New("find by id error")},
		{"failure customer not found", false, context.Background(), customer.ErrCustomerNotFound},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockCustomerRepository.EXPECT().FindByID(tt.ctx, gomock.Any()).Return(mockCustomer, tt.findByIDErr).AnyTimes()

			u := NewCustomerUsecase(mockCustomerRepository)

			_, err := u.GetCustomer(tt.ctx, customer.NewCustomerID())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestUpdateCustomer(t *testing.T) {
	t.Parallel()
	name, _ := customer.NewName("updated test customer")

	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		findByIDErr error
		updateErr   error
	}{
		{"success update customer", true, context.Background(), nil, nil},
		{"failure find by id error", false, context.Background(), errors.New("find by id error"), nil},
		{"failure customer not found", false, context.Background(), customer.ErrCustomerNotFound, nil},
		{"failure update error", false, context.Background(), nil, errors.New("update error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().Update(gomock.Any()).AnyTimes()
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockCustomerRepository.EXPECT().FindByID(tt.ctx, gomock.Any()).Return(mockCustomer, tt.findByIDErr).AnyTimes()
			mockCustomerRepository.EXPECT().Update(tt.ctx, gomock.Any()).Return(tt.updateErr).AnyTimes()

			u := NewCustomerUsecase(mockCustomerRepository)

			_, err := u.UpdateCustomer(tt.ctx, customer.NewCustomerID(), name)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestDeleteCustomer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		findByIDErr error
		deleteErr   error
	}{
		{"success delete customer", true, context.Background(), nil, nil},
		{"failure find by id error", false, context.Background(), errors.New("find by id error"), nil},
		{"failure customer not found", false, context.Background(), customer.ErrCustomerNotFound, nil},
		{"failure delete error", false, context.Background(), nil, errors.New("delete error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockCustomerRepository.EXPECT().FindByID(tt.ctx, gomock.Any()).Return(mockCustomer, tt.findByIDErr).AnyTimes()
			mockCustomerRepository.EXPECT().Delete(tt.ctx, gomock.Any()).Return(tt.deleteErr).AnyTimes()

			u := NewCustomerUsecase(mockCustomerRepository)

			err := u.DeleteCustomer(tt.ctx, customer.NewCustomerID())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
