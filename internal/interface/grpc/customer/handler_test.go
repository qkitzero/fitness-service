package customer

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customerv1 "github.com/qkitzero/fitness-service/gen/go/customer/v1"
	"github.com/qkitzero/fitness-service/internal/application/user"
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
		groupID           string
		nameKana          string
		gender            customerv1.Gender
		birthDate         *date.Date
		callUsecase       bool
		createCustomerErr error
		wantCode          codes.Code
	}{
		{"success create customer", context.Background(), "test customer", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "テストカナ", customerv1.Gender_GENDER_MALE, &date.Date{Year: 2000, Month: 1, Day: 1}, true, nil, codes.OK},
		{"failure invalid name", context.Background(), "", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "テストカナ", customerv1.Gender_GENDER_MALE, &date.Date{Year: 2000, Month: 1, Day: 1}, false, nil, codes.InvalidArgument},
		{"failure invalid group id", context.Background(), "test customer", "", "テストカナ", customerv1.Gender_GENDER_MALE, &date.Date{Year: 2000, Month: 1, Day: 1}, false, nil, codes.InvalidArgument},
		{"failure invalid name kana", context.Background(), "test customer", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "", customerv1.Gender_GENDER_MALE, &date.Date{Year: 2000, Month: 1, Day: 1}, false, nil, codes.InvalidArgument},
		{"failure invalid gender", context.Background(), "test customer", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "テストカナ", customerv1.Gender_GENDER_UNSPECIFIED, &date.Date{Year: 2000, Month: 1, Day: 1}, false, nil, codes.InvalidArgument},
		{"failure invalid birth date", context.Background(), "test customer", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "テストカナ", customerv1.Gender_GENDER_MALE, &date.Date{Year: 2500, Month: 1, Day: 1}, false, nil, codes.InvalidArgument},
		{"failure not group member", context.Background(), "test customer", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "テストカナ", customerv1.Gender_GENDER_MALE, &date.Date{Year: 2000, Month: 1, Day: 1}, true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure usecase error", context.Background(), "test customer", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "テストカナ", customerv1.Gender_GENDER_MALE, &date.Date{Year: 2000, Month: 1, Day: 1}, true, fmt.Errorf("create customer error"), codes.Internal},
		{"failure status preserved", context.Background(), "test customer", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "テストカナ", customerv1.Gender_GENDER_MALE, &date.Date{Year: 2000, Month: 1, Day: 1}, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
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
				mockUsecase.EXPECT().CreateCustomer(tt.ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCustomer, tt.createCustomerErr).Times(1)
				mockCustomer.EXPECT().ID().Return(customer.NewCustomerID()).AnyTimes()
			}

			handler := NewCustomerHandler(mockUsecase)

			req := &customerv1.CreateCustomerRequest{
				Name:      tt.customerName,
				GroupId:   tt.groupID,
				NameKana:  tt.nameKana,
				Gender:    tt.gender,
				BirthDate: tt.birthDate,
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
		groupID, _ := customer.NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
		nameKana, _ := customer.NewNameKana("テストカナ")
		birthDate, _ := customer.NewBirthDate(2000, 1, 1)
		m.EXPECT().ID().Return(customer.NewCustomerID()).AnyTimes()
		m.EXPECT().Name().Return(name).AnyTimes()
		m.EXPECT().GroupID().Return(groupID).AnyTimes()
		m.EXPECT().NameKana().Return(nameKana).AnyTimes()
		m.EXPECT().Gender().Return(customer.GenderMale).AnyTimes()
		m.EXPECT().BirthDate().Return(birthDate).AnyTimes()
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
		{"failure not group member", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, user.ErrNotGroupMember, codes.PermissionDenied},
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

func TestListCustomers(t *testing.T) {
	t.Parallel()

	mockCustomerSample := func(ctrl *gomock.Controller) *mockscustomer.MockCustomer {
		m := mockscustomer.NewMockCustomer(ctrl)
		name, _ := customer.NewName("test customer")
		groupID, _ := customer.NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
		nameKana, _ := customer.NewNameKana("テストカナ")
		birthDate, _ := customer.NewBirthDate(2000, 1, 1)
		m.EXPECT().ID().Return(customer.NewCustomerID()).AnyTimes()
		m.EXPECT().Name().Return(name).AnyTimes()
		m.EXPECT().GroupID().Return(groupID).AnyTimes()
		m.EXPECT().NameKana().Return(nameKana).AnyTimes()
		m.EXPECT().Gender().Return(customer.GenderMale).AnyTimes()
		m.EXPECT().BirthDate().Return(birthDate).AnyTimes()
		return m
	}

	tests := []struct {
		name             string
		ctx              context.Context
		groupID          string
		callUsecase      bool
		listCustomersErr error
		wantCode         codes.Code
	}{
		{"success list customers", context.Background(), "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", true, nil, codes.OK},
		{"failure invalid group id", context.Background(), "", false, nil, codes.InvalidArgument},
		{"failure not group member", context.Background(), "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure usecase error", context.Background(), "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", true, fmt.Errorf("list customers error"), codes.Internal},
		{"failure status preserved", context.Background(), "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappcustomer.NewMockCustomerUsecase(ctrl)
			if tt.callUsecase {
				customers := []customer.Customer{mockCustomerSample(ctrl), mockCustomerSample(ctrl)}
				mockUsecase.EXPECT().ListCustomers(tt.ctx, gomock.Any()).Return(customers, tt.listCustomersErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			req := &customerv1.ListCustomersRequest{
				GroupId: tt.groupID,
			}

			res, err := handler.ListCustomers(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK && len(res.GetCustomers()) != 2 {
				t.Errorf("len(customers) = %v, want %v", len(res.GetCustomers()), 2)
			}
		})
	}
}

func TestUpdateCustomer(t *testing.T) {
	t.Parallel()

	mockCustomerSample := func(ctrl *gomock.Controller) *mockscustomer.MockCustomer {
		m := mockscustomer.NewMockCustomer(ctrl)
		name, _ := customer.NewName("updated test customer")
		groupID, _ := customer.NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
		nameKana, _ := customer.NewNameKana("コウシンカナ")
		birthDate, _ := customer.NewBirthDate(1999, 12, 31)
		m.EXPECT().ID().Return(customer.NewCustomerID()).AnyTimes()
		m.EXPECT().Name().Return(name).AnyTimes()
		m.EXPECT().GroupID().Return(groupID).AnyTimes()
		m.EXPECT().NameKana().Return(nameKana).AnyTimes()
		m.EXPECT().Gender().Return(customer.GenderFemale).AnyTimes()
		m.EXPECT().BirthDate().Return(birthDate).AnyTimes()
		return m
	}

	tests := []struct {
		name              string
		ctx               context.Context
		customerID        string
		customerName      string
		nameKana          string
		gender            customerv1.Gender
		birthDate         *date.Date
		callUsecase       bool
		updateCustomerErr error
		wantCode          codes.Code
	}{
		{"success update customer", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "updated test customer", "コウシンカナ", customerv1.Gender_GENDER_FEMALE, &date.Date{Year: 1999, Month: 12, Day: 31}, true, nil, codes.OK},
		{"failure invalid customer id", context.Background(), "", "updated test customer", "コウシンカナ", customerv1.Gender_GENDER_FEMALE, &date.Date{Year: 1999, Month: 12, Day: 31}, false, nil, codes.InvalidArgument},
		{"failure invalid name", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "", "コウシンカナ", customerv1.Gender_GENDER_FEMALE, &date.Date{Year: 1999, Month: 12, Day: 31}, false, nil, codes.InvalidArgument},
		{"failure invalid name kana", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "updated test customer", "", customerv1.Gender_GENDER_FEMALE, &date.Date{Year: 1999, Month: 12, Day: 31}, false, nil, codes.InvalidArgument},
		{"failure invalid gender", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "updated test customer", "コウシンカナ", customerv1.Gender_GENDER_UNSPECIFIED, &date.Date{Year: 1999, Month: 12, Day: 31}, false, nil, codes.InvalidArgument},
		{"failure invalid birth date", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "updated test customer", "コウシンカナ", customerv1.Gender_GENDER_FEMALE, &date.Date{Year: 2500, Month: 1, Day: 1}, false, nil, codes.InvalidArgument},
		{"failure usecase error", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "updated test customer", "コウシンカナ", customerv1.Gender_GENDER_FEMALE, &date.Date{Year: 1999, Month: 12, Day: 31}, true, fmt.Errorf("update customer error"), codes.Internal},
		{"failure not group member", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "updated test customer", "コウシンカナ", customerv1.Gender_GENDER_FEMALE, &date.Date{Year: 1999, Month: 12, Day: 31}, true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure customer not found", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "updated test customer", "コウシンカナ", customerv1.Gender_GENDER_FEMALE, &date.Date{Year: 1999, Month: 12, Day: 31}, true, customer.ErrCustomerNotFound, codes.NotFound},
		{"failure status preserved", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "updated test customer", "コウシンカナ", customerv1.Gender_GENDER_FEMALE, &date.Date{Year: 1999, Month: 12, Day: 31}, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappcustomer.NewMockCustomerUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().UpdateCustomer(tt.ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCustomerSample(ctrl), tt.updateCustomerErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			req := &customerv1.UpdateCustomerRequest{
				CustomerId: tt.customerID,
				Name:       tt.customerName,
				NameKana:   tt.nameKana,
				Gender:     tt.gender,
				BirthDate:  tt.birthDate,
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
		{"failure not group member", context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, user.ErrNotGroupMember, codes.PermissionDenied},
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

func TestToDomainGender(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		gender  customerv1.Gender
		want    customer.Gender
	}{
		{"male", true, customerv1.Gender_GENDER_MALE, customer.GenderMale},
		{"female", true, customerv1.Gender_GENDER_FEMALE, customer.GenderFemale},
		{"other", true, customerv1.Gender_GENDER_OTHER, customer.GenderOther},
		{"failure unspecified", false, customerv1.Gender_GENDER_UNSPECIFIED, customer.Gender("")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := toDomainGender(tt.gender)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if got != tt.want {
				t.Errorf("toDomainGender() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToProtoGender(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		gender customer.Gender
		want   customerv1.Gender
	}{
		{"male", customer.GenderMale, customerv1.Gender_GENDER_MALE},
		{"female", customer.GenderFemale, customerv1.Gender_GENDER_FEMALE},
		{"other", customer.GenderOther, customerv1.Gender_GENDER_OTHER},
		{"unknown", customer.Gender("unknown"), customerv1.Gender_GENDER_UNSPECIFIED},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := toProtoGender(tt.gender); got != tt.want {
				t.Errorf("toProtoGender() = %v, want %v", got, tt.want)
			}
		})
	}
}
