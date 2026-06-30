package customer

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customerv1 "github.com/qkitzero/fitness-service/gen/go/customer/v1"
	"github.com/qkitzero/fitness-service/internal/domain/customer"
	mocksappcustomer "github.com/qkitzero/fitness-service/mocks/application/customer"
	mockscustomer "github.com/qkitzero/fitness-service/mocks/domain/customer"
)

func TestCreateCustomer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		ctx               context.Context
		customerName      string
		callUsecase       bool
		createCustomerErr error
		wantCode          codes.Code
	}{
		{"success create customer", context.Background(), "test customer", true, nil, codes.OK},
		{"failure invalid name", context.Background(), "", false, nil, codes.InvalidArgument},
		{"failure usecase error", context.Background(), "test customer", true, fmt.Errorf("create customer error"), codes.Internal},
		{"failure status preserved", context.Background(), "test customer", true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappcustomer.NewMockCustomerUsecase(ctrl)
			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().CreateCustomer(tt.ctx, gomock.Any()).Return(mockCustomer, tt.createCustomerErr).Times(1)
				mockCustomer.EXPECT().ID().Return(customer.NewCustomerID()).AnyTimes()
			}

			handler := NewCustomerHandler(mockUsecase)

			req := &customerv1.CreateCustomerRequest{
				Name: tt.customerName,
			}

			_, err := handler.CreateCustomer(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestGetCustomer(t *testing.T) {
	t.Parallel()

	mockCustomerSample := func(ctrl *gomock.Controller) *mockscustomer.MockCustomer {
		m := mockscustomer.NewMockCustomer(ctrl)
		name, _ := customer.NewName("test customer")
		m.EXPECT().ID().Return(customer.NewCustomerID()).AnyTimes()
		m.EXPECT().Name().Return(name).AnyTimes()
		return m
	}

	tests := []struct {
		name           string
		ctx            context.Context
		customerID     string
		callUsecase    bool
		getCustomerErr error
		wantCode       codes.Code
	}{
		{"success get customer", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, nil, codes.OK},
		{"failure invalid customer id", context.Background(), "", false, nil, codes.InvalidArgument},
		{"failure get customer error", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, fmt.Errorf("get customer error"), codes.Internal},
		{"failure customer not found", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, customer.ErrCustomerNotFound, codes.NotFound},
		{"failure status preserved", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappcustomer.NewMockCustomerUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().GetCustomer(tt.ctx, gomock.Any()).Return(mockCustomerSample(ctrl), tt.getCustomerErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			req := &customerv1.GetCustomerRequest{
				CustomerId: tt.customerID,
			}

			_, err := handler.GetCustomer(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestUpdateCustomer(t *testing.T) {
	t.Parallel()

	mockCustomerSample := func(ctrl *gomock.Controller) *mockscustomer.MockCustomer {
		m := mockscustomer.NewMockCustomer(ctrl)
		name, _ := customer.NewName("updated test customer")
		m.EXPECT().ID().Return(customer.NewCustomerID()).AnyTimes()
		m.EXPECT().Name().Return(name).AnyTimes()
		return m
	}

	tests := []struct {
		name              string
		ctx               context.Context
		customerID        string
		customerName      string
		callUsecase       bool
		updateCustomerErr error
		wantCode          codes.Code
	}{
		{"success update customer", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "updated test customer", true, nil, codes.OK},
		{"failure invalid customer id", context.Background(), "", "updated test customer", false, nil, codes.InvalidArgument},
		{"failure invalid name", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "", false, nil, codes.InvalidArgument},
		{"failure usecase error", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "updated test customer", true, fmt.Errorf("update customer error"), codes.Internal},
		{"failure customer not found", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "updated test customer", true, customer.ErrCustomerNotFound, codes.NotFound},
		{"failure status preserved", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "updated test customer", true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappcustomer.NewMockCustomerUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().UpdateCustomer(tt.ctx, gomock.Any(), gomock.Any()).Return(mockCustomerSample(ctrl), tt.updateCustomerErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			req := &customerv1.UpdateCustomerRequest{
				CustomerId: tt.customerID,
				Name:       tt.customerName,
			}

			_, err := handler.UpdateCustomer(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestDeleteCustomer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		ctx               context.Context
		customerID        string
		callUsecase       bool
		deleteCustomerErr error
		wantCode          codes.Code
	}{
		{"success delete customer", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, nil, codes.OK},
		{"failure invalid customer id", context.Background(), "", false, nil, codes.InvalidArgument},
		{"failure usecase error", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, fmt.Errorf("delete customer error"), codes.Internal},
		{"failure customer not found", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, customer.ErrCustomerNotFound, codes.NotFound},
		{"failure status preserved", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappcustomer.NewMockCustomerUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().DeleteCustomer(tt.ctx, gomock.Any()).Return(tt.deleteCustomerErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			req := &customerv1.DeleteCustomerRequest{
				CustomerId: tt.customerID,
			}

			_, err := handler.DeleteCustomer(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}
