package customer

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customerv1 "github.com/qkitzero/fitness-service/gen/go/customer/v1"
	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/contact"
	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
	mocksappcustomer "github.com/qkitzero/fitness-service/mocks/application/customer"
	mockscustomer "github.com/qkitzero/fitness-service/mocks/domain/customer"
)

const (
	sampleCustomerID     = "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"
	sampleTenantID       = "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"
	sampleOrganizationID = "3f2b6c1d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
)

func customerSample(ctrl *gomock.Controller, active bool) *mockscustomer.MockCustomer {
	m := mockscustomer.NewMockCustomer(ctrl)
	id, _ := customer.NewCustomerIDFromString(sampleCustomerID)
	tenantID, _ := tenant.NewTenantID(sampleTenantID)
	name, _ := customer.NewName("test customer")
	nameKana, _ := customer.NewNameKana("テストカナ")
	birthDate, _ := customer.NewBirthDate(2000, 1, 1)
	phone, _ := contact.NewPhone("03-1234-5678")
	email, _ := contact.NewEmail("test@example.com")
	postalCode, _ := address.NewPostalCode("123-4567")
	prefecture, _ := address.NewPrefecture("東京都")
	city, _ := address.NewCity("千代田区")
	street, _ := address.NewStreet("1-1-1")
	building, _ := address.NewBuilding("テストビル")
	addr := address.NewAddress(postalCode, prefecture, city, street, building)
	emergencyContactName, _ := customer.NewEmergencyContactName("緊急 太郎")
	emergencyContactRelationship, _ := customer.NewEmergencyContactRelationship("父")
	emergencyContactPhone, _ := contact.NewPhone("090-1234-5678")
	organizationID, _ := organization.NewOrganizationIDFromString(sampleOrganizationID)
	m.EXPECT().ID().Return(id).AnyTimes()
	m.EXPECT().Name().Return(name).AnyTimes()
	m.EXPECT().TenantID().Return(tenantID).AnyTimes()
	m.EXPECT().NameKana().Return(nameKana).AnyTimes()
	m.EXPECT().Gender().Return(customer.GenderMale).AnyTimes()
	m.EXPECT().BirthDate().Return(birthDate).AnyTimes()
	m.EXPECT().Phone().Return(phone).AnyTimes()
	m.EXPECT().Email().Return(email).AnyTimes()
	m.EXPECT().Address().Return(addr).AnyTimes()
	m.EXPECT().EmergencyContactName().Return(emergencyContactName).AnyTimes()
	m.EXPECT().EmergencyContactRelationship().Return(emergencyContactRelationship).AnyTimes()
	m.EXPECT().EmergencyContactPhone().Return(emergencyContactPhone).AnyTimes()
	m.EXPECT().OrganizationID().Return(&organizationID).AnyTimes()
	m.EXPECT().IsActive().Return(active).AnyTimes()
	return m
}

func TestCreateCustomer(t *testing.T) {
	t.Parallel()
	tid := sampleTenantID
	validDate := &date.Date{Year: 2000, Month: 1, Day: 1}
	futureDate := &date.Date{Year: 2500, Month: 1, Day: 1}
	phone := "03-1234-5678"
	email := "test@example.com"
	postalCode := "123-4567"
	prefecture := "東京都"
	city := "千代田区"
	street := "1-1-1"
	building := "テストビル"
	ecName := "緊急 太郎"
	ecRelationship := "父"
	ecPhone := "090-1234-5678"
	organizationID := sampleOrganizationID
	invalid := "invalid"
	invalidPref := "存在しない県"
	tooLong := strings.Repeat("あ", 256)
	withNull := "テスト\x00顧客"

	tests := []struct {
		name        string
		req         *customerv1.CreateCustomerRequest
		callUsecase bool
		createErr   error
		wantCode    codes.Code
	}{
		{"success", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, Phone: &phone, Email: &email, PostalCode: &postalCode, Prefecture: &prefecture, City: &city, Street: &street, Building: &building, EmergencyContactName: &ecName, EmergencyContactRelationship: &ecRelationship, EmergencyContactPhone: &ecPhone, OrganizationId: &organizationID}, true, nil, codes.OK},
		{"failure invalid tenant id", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: "", NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid name", &customerv1.CreateCustomerRequest{Name: "", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure null character name", &customerv1.CreateCustomerRequest{Name: withNull, TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid name kana", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid gender", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_UNSPECIFIED, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid birth date", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: futureDate}, false, nil, codes.InvalidArgument},
		{"failure invalid phone", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, Phone: &invalid}, false, nil, codes.InvalidArgument},
		{"failure invalid email", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, Email: &invalid}, false, nil, codes.InvalidArgument},
		{"failure invalid postal code", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, PostalCode: &invalid}, false, nil, codes.InvalidArgument},
		{"failure invalid prefecture", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, Prefecture: &invalidPref}, false, nil, codes.InvalidArgument},
		{"failure invalid city", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, City: &tooLong}, false, nil, codes.InvalidArgument},
		{"failure null character city", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, City: &withNull}, false, nil, codes.InvalidArgument},
		{"failure invalid street", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, Street: &tooLong}, false, nil, codes.InvalidArgument},
		{"failure invalid building", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, Building: &tooLong}, false, nil, codes.InvalidArgument},
		{"failure invalid emergency contact name", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, EmergencyContactName: &tooLong}, false, nil, codes.InvalidArgument},
		{"failure invalid emergency contact relationship", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, EmergencyContactRelationship: &tooLong}, false, nil, codes.InvalidArgument},
		{"failure invalid emergency contact phone", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, EmergencyContactPhone: &invalid}, false, nil, codes.InvalidArgument},
		{"failure invalid organization id", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, OrganizationId: &invalid}, false, nil, codes.InvalidArgument},
		{"failure organization not in tenant", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, OrganizationId: &organizationID}, true, customer.ErrOrganizationNotInTenant, codes.InvalidArgument},
		{"failure not tenant member", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, true, tenant.ErrNotMember, codes.PermissionDenied},
		{"failure usecase error", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, true, fmt.Errorf("create customer error"), codes.Internal},
		{"failure unauthenticated is preserved", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure downstream code is not forwarded", &customerv1.CreateCustomerRequest{Name: "test customer", TenantId: tid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, true, status.Error(codes.NotFound, "user not found"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockUsecase := mocksappcustomer.NewMockCustomerUsecase(ctrl)
			if tt.callUsecase {
				tenantID, _ := tenant.NewTenantID(tid)
				mockUsecase.EXPECT().CreateCustomer(gomock.Any(), tenantID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(customerSample(ctrl, true), tt.createErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			res, err := handler.CreateCustomer(ctx, tt.req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK && res.GetCustomerId() != sampleCustomerID {
				t.Errorf("CustomerId = %v, want %v", res.GetCustomerId(), sampleCustomerID)
			}
		})
	}
}

func TestGetCustomer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		customerID     string
		callUsecase    bool
		getCustomerErr error
		wantCode       codes.Code
	}{
		{"success get customer", sampleCustomerID, true, nil, codes.OK},
		{"failure invalid customer id", "", false, nil, codes.InvalidArgument},
		{"failure get customer error", sampleCustomerID, true, fmt.Errorf("get customer error"), codes.Internal},
		{"failure not tenant member", sampleCustomerID, true, tenant.ErrNotMember, codes.PermissionDenied},
		{"failure customer not found", sampleCustomerID, true, customer.ErrCustomerNotFound, codes.NotFound},
		{"failure unauthenticated is preserved", sampleCustomerID, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure downstream code is not forwarded", sampleCustomerID, true, status.Error(codes.InvalidArgument, "user service"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockUsecase := mocksappcustomer.NewMockCustomerUsecase(ctrl)
			if tt.callUsecase {
				customerID, _ := customer.NewCustomerIDFromString(sampleCustomerID)
				mockUsecase.EXPECT().GetCustomer(gomock.Any(), customerID).Return(customerSample(ctrl, true), tt.getCustomerErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			res, err := handler.GetCustomer(ctx, &customerv1.GetCustomerRequest{CustomerId: tt.customerID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				if res.GetCustomer().GetCustomerId() != sampleCustomerID {
					t.Errorf("CustomerId = %v, want %v", res.GetCustomer().GetCustomerId(), sampleCustomerID)
				}
				if res.GetCustomer().GetTenantId() != sampleTenantID {
					t.Errorf("TenantId = %v, want %v", res.GetCustomer().GetTenantId(), sampleTenantID)
				}
				if res.GetCustomer().GetName() != "test customer" {
					t.Errorf("Name = %v, want test customer", res.GetCustomer().GetName())
				}
				if res.GetCustomer().GetOrganizationId() != sampleOrganizationID {
					t.Errorf("OrganizationId = %v, want %v", res.GetCustomer().GetOrganizationId(), sampleOrganizationID)
				}
			}
		})
	}
}

func TestListCustomers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		tenantID         string
		includeInactive  bool
		callUsecase      bool
		listCustomersErr error
		wantCode         codes.Code
	}{
		{"success list customers", sampleTenantID, false, true, nil, codes.OK},
		{"success list customers including inactive", sampleTenantID, true, true, nil, codes.OK},
		{"failure invalid tenant id", "", false, false, nil, codes.InvalidArgument},
		{"failure not tenant member", sampleTenantID, false, true, tenant.ErrNotMember, codes.PermissionDenied},
		{"failure usecase error", sampleTenantID, false, true, fmt.Errorf("list customers error"), codes.Internal},
		{"failure unauthenticated is preserved", sampleTenantID, false, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure downstream code is not forwarded", sampleTenantID, false, true, status.Error(codes.NotFound, "user not found"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockUsecase := mocksappcustomer.NewMockCustomerUsecase(ctrl)
			if tt.callUsecase {
				tenantID, _ := tenant.NewTenantID(sampleTenantID)
				customers := []customer.Customer{customerSample(ctrl, true), customerSample(ctrl, false)}
				mockUsecase.EXPECT().ListCustomers(gomock.Any(), tenantID, tt.includeInactive).Return(customers, tt.listCustomersErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			res, err := handler.ListCustomers(ctx, &customerv1.ListCustomersRequest{TenantId: tt.tenantID, IncludeInactive: tt.includeInactive})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				wantActive := []bool{true, false}
				if len(res.GetCustomers()) != len(wantActive) {
					t.Fatalf("len(customers) = %v, want %v", len(res.GetCustomers()), len(wantActive))
				}
				for i, c := range res.GetCustomers() {
					if c.GetCustomerId() != sampleCustomerID {
						t.Errorf("customers[%d].CustomerId = %v, want %v", i, c.GetCustomerId(), sampleCustomerID)
					}
					if c.GetTenantId() != sampleTenantID {
						t.Errorf("customers[%d].TenantId = %v, want %v", i, c.GetTenantId(), sampleTenantID)
					}
					if c.GetIsActive() != wantActive[i] {
						t.Errorf("customers[%d].IsActive = %v, want %v", i, c.GetIsActive(), wantActive[i])
					}
				}
			}
		})
	}
}

func TestUpdateCustomer(t *testing.T) {
	t.Parallel()
	cid := sampleCustomerID
	validDate := &date.Date{Year: 1999, Month: 12, Day: 31}
	futureDate := &date.Date{Year: 2500, Month: 1, Day: 1}
	organizationID := sampleOrganizationID
	invalid := "invalid"
	withNull := "コウシン\x00顧客"

	tests := []struct {
		name        string
		req         *customerv1.UpdateCustomerRequest
		callUsecase bool
		updateErr   error
		wantCode    codes.Code
	}{
		{"success update customer", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate, OrganizationId: &organizationID}, true, nil, codes.OK},
		{"failure invalid customer id", &customerv1.UpdateCustomerRequest{CustomerId: "", Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid name", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure null character name", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: withNull, NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid name kana", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid gender", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_UNSPECIFIED, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid birth date", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: futureDate}, false, nil, codes.InvalidArgument},
		{"failure invalid phone", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate, Phone: &invalid}, false, nil, codes.InvalidArgument},
		{"failure invalid organization id", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate, OrganizationId: &invalid}, false, nil, codes.InvalidArgument},
		{"failure organization not in tenant", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate, OrganizationId: &organizationID}, true, customer.ErrOrganizationNotInTenant, codes.InvalidArgument},
		{"failure usecase error", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, true, fmt.Errorf("update customer error"), codes.Internal},
		{"failure not tenant member", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, true, tenant.ErrNotMember, codes.PermissionDenied},
		{"failure customer not found", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, true, customer.ErrCustomerNotFound, codes.NotFound},
		{"failure unauthenticated is preserved", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure downstream code is not forwarded", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, true, status.Error(codes.InvalidArgument, "user service"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockUsecase := mocksappcustomer.NewMockCustomerUsecase(ctrl)
			if tt.callUsecase {
				customerID, _ := customer.NewCustomerIDFromString(cid)
				mockUsecase.EXPECT().UpdateCustomer(gomock.Any(), customerID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(customerSample(ctrl, true), tt.updateErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			res, err := handler.UpdateCustomer(ctx, tt.req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				if res.GetCustomer().GetCustomerId() != sampleCustomerID {
					t.Errorf("CustomerId = %v, want %v", res.GetCustomer().GetCustomerId(), sampleCustomerID)
				}
				if res.GetCustomer().GetName() != "test customer" {
					t.Errorf("Name = %v, want test customer", res.GetCustomer().GetName())
				}
				if res.GetCustomer().GetOrganizationId() != sampleOrganizationID {
					t.Errorf("OrganizationId = %v, want %v", res.GetCustomer().GetOrganizationId(), sampleOrganizationID)
				}
			}
		})
	}
}

func TestSetCustomerActive(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                 string
		customerID           string
		active               bool
		callUsecase          bool
		setCustomerActiveErr error
		wantCode             codes.Code
	}{
		{"success deactivate customer", sampleCustomerID, false, true, nil, codes.OK},
		{"success activate customer", sampleCustomerID, true, true, nil, codes.OK},
		{"failure invalid customer id", "", false, false, nil, codes.InvalidArgument},
		{"failure usecase error", sampleCustomerID, false, true, fmt.Errorf("set customer active error"), codes.Internal},
		{"failure not tenant member", sampleCustomerID, false, true, tenant.ErrNotMember, codes.PermissionDenied},
		{"failure customer not found", sampleCustomerID, false, true, customer.ErrCustomerNotFound, codes.NotFound},
		{"failure unauthenticated is preserved", sampleCustomerID, false, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure downstream code is not forwarded", sampleCustomerID, false, true, status.Error(codes.NotFound, "user not found"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockUsecase := mocksappcustomer.NewMockCustomerUsecase(ctrl)
			if tt.callUsecase {
				customerID, _ := customer.NewCustomerIDFromString(sampleCustomerID)
				mockUsecase.EXPECT().SetCustomerActive(gomock.Any(), customerID, tt.active).Return(customerSample(ctrl, tt.active), tt.setCustomerActiveErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			res, err := handler.SetCustomerActive(ctx, &customerv1.SetCustomerActiveRequest{CustomerId: tt.customerID, IsActive: tt.active})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				if res.GetCustomer().GetCustomerId() != sampleCustomerID {
					t.Errorf("CustomerId = %v, want %v", res.GetCustomer().GetCustomerId(), sampleCustomerID)
				}
				if res.GetCustomer().GetIsActive() != tt.active {
					t.Errorf("IsActive = %v, want %v", res.GetCustomer().GetIsActive(), tt.active)
				}
			}
		})
	}
}

func TestDeleteCustomer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		customerID        string
		callUsecase       bool
		deleteCustomerErr error
		wantCode          codes.Code
	}{
		{"success delete customer", sampleCustomerID, true, nil, codes.OK},
		{"failure invalid customer id", "", false, nil, codes.InvalidArgument},
		{"failure usecase error", sampleCustomerID, true, fmt.Errorf("delete customer error"), codes.Internal},
		{"failure not tenant member", sampleCustomerID, true, tenant.ErrNotMember, codes.PermissionDenied},
		{"failure customer not found", sampleCustomerID, true, customer.ErrCustomerNotFound, codes.NotFound},
		{"failure unauthenticated is preserved", sampleCustomerID, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure downstream code is not forwarded", sampleCustomerID, true, status.Error(codes.NotFound, "user not found"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockUsecase := mocksappcustomer.NewMockCustomerUsecase(ctrl)
			if tt.callUsecase {
				customerID, _ := customer.NewCustomerIDFromString(sampleCustomerID)
				mockUsecase.EXPECT().DeleteCustomer(gomock.Any(), customerID).Return(tt.deleteCustomerErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			_, err := handler.DeleteCustomer(ctx, &customerv1.DeleteCustomerRequest{CustomerId: tt.customerID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestParseCustomerFields(t *testing.T) {
	t.Parallel()
	phone := "03-1234-5678"
	email := "test@example.com"
	postalCode := "123-4567"
	prefecture := "東京都"
	city := "千代田区"
	street := "1-1-1"
	building := "テストビル"
	ecName := "緊急 太郎"
	ecRelationship := "父"
	ecPhone := "090-1234-5678"
	organizationID := sampleOrganizationID

	tests := []struct {
		name                             string
		req                              *customerv1.CreateCustomerRequest
		wantPhone                        string
		wantEmail                        string
		wantPostalCode                   string
		wantPrefecture                   string
		wantCity                         string
		wantStreet                       string
		wantBuilding                     string
		wantEmergencyContactName         string
		wantEmergencyContactRelationship string
		wantEmergencyContactPhone        string
		wantOrganizationID               string
	}{
		{
			"success parse customer fields",
			&customerv1.CreateCustomerRequest{
				Name:                         "test customer",
				NameKana:                     "テストカナ",
				Gender:                       customerv1.Gender_GENDER_MALE,
				BirthDate:                    &date.Date{Year: 2000, Month: 1, Day: 1},
				Phone:                        &phone,
				Email:                        &email,
				PostalCode:                   &postalCode,
				Prefecture:                   &prefecture,
				City:                         &city,
				Street:                       &street,
				Building:                     &building,
				EmergencyContactName:         &ecName,
				EmergencyContactRelationship: &ecRelationship,
				EmergencyContactPhone:        &ecPhone,
				OrganizationId:               &organizationID,
			},
			"0312345678", "test@example.com", "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "緊急 太郎", "父", "09012345678", sampleOrganizationID,
		},
		{
			"success parse customer fields without optional fields",
			&customerv1.CreateCustomerRequest{
				Name:      "test customer",
				NameKana:  "テストカナ",
				Gender:    customerv1.Gender_GENDER_MALE,
				BirthDate: &date.Date{Year: 2000, Month: 1, Day: 1},
			},
			"", "", "", "", "", "", "", "", "", "", "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseCustomerFields(tt.req)
			if err != nil {
				t.Fatalf("expected no error, but got %v", err)
			}
			if got.name.String() != "test customer" {
				t.Errorf("name = %v, want test customer", got.name.String())
			}
			if got.nameKana.String() != "テストカナ" {
				t.Errorf("nameKana = %v, want テストカナ", got.nameKana.String())
			}
			if got.gender != customer.GenderMale {
				t.Errorf("gender = %v, want %v", got.gender, customer.GenderMale)
			}
			if got.birthDate.Year() != 2000 || got.birthDate.Month() != time.January || got.birthDate.Day() != 1 {
				t.Errorf("birthDate = %v, want 2000-01-01", got.birthDate)
			}
			if tt.wantPhone == "" && got.phone != nil {
				t.Errorf("phone = %v, want nil", got.phone)
			}
			if tt.wantPhone != "" && (got.phone == nil || got.phone.String() != tt.wantPhone) {
				t.Errorf("phone = %v, want %v", got.phone, tt.wantPhone)
			}
			if tt.wantEmail == "" && got.email != nil {
				t.Errorf("email = %v, want nil", got.email)
			}
			if tt.wantEmail != "" && (got.email == nil || got.email.String() != tt.wantEmail) {
				t.Errorf("email = %v, want %v", got.email, tt.wantEmail)
			}
			if tt.wantPostalCode == "" && got.address.PostalCode() != nil {
				t.Errorf("address.PostalCode() = %v, want nil", got.address.PostalCode())
			}
			if tt.wantPostalCode != "" && (got.address.PostalCode() == nil || got.address.PostalCode().String() != tt.wantPostalCode) {
				t.Errorf("address.PostalCode() = %v, want %v", got.address.PostalCode(), tt.wantPostalCode)
			}
			if tt.wantPrefecture == "" && got.address.Prefecture() != nil {
				t.Errorf("address.Prefecture() = %v, want nil", got.address.Prefecture())
			}
			if tt.wantPrefecture != "" && (got.address.Prefecture() == nil || got.address.Prefecture().String() != tt.wantPrefecture) {
				t.Errorf("address.Prefecture() = %v, want %v", got.address.Prefecture(), tt.wantPrefecture)
			}
			if tt.wantCity == "" && got.address.City() != nil {
				t.Errorf("address.City() = %v, want nil", got.address.City())
			}
			if tt.wantCity != "" && (got.address.City() == nil || got.address.City().String() != tt.wantCity) {
				t.Errorf("address.City() = %v, want %v", got.address.City(), tt.wantCity)
			}
			if tt.wantStreet == "" && got.address.Street() != nil {
				t.Errorf("address.Street() = %v, want nil", got.address.Street())
			}
			if tt.wantStreet != "" && (got.address.Street() == nil || got.address.Street().String() != tt.wantStreet) {
				t.Errorf("address.Street() = %v, want %v", got.address.Street(), tt.wantStreet)
			}
			if tt.wantBuilding == "" && got.address.Building() != nil {
				t.Errorf("address.Building() = %v, want nil", got.address.Building())
			}
			if tt.wantBuilding != "" && (got.address.Building() == nil || got.address.Building().String() != tt.wantBuilding) {
				t.Errorf("address.Building() = %v, want %v", got.address.Building(), tt.wantBuilding)
			}
			if tt.wantEmergencyContactName == "" && got.emergencyContactName != nil {
				t.Errorf("emergencyContactName = %v, want nil", got.emergencyContactName)
			}
			if tt.wantEmergencyContactName != "" && (got.emergencyContactName == nil || got.emergencyContactName.String() != tt.wantEmergencyContactName) {
				t.Errorf("emergencyContactName = %v, want %v", got.emergencyContactName, tt.wantEmergencyContactName)
			}
			if tt.wantEmergencyContactRelationship == "" && got.emergencyContactRelationship != nil {
				t.Errorf("emergencyContactRelationship = %v, want nil", got.emergencyContactRelationship)
			}
			if tt.wantEmergencyContactRelationship != "" && (got.emergencyContactRelationship == nil || got.emergencyContactRelationship.String() != tt.wantEmergencyContactRelationship) {
				t.Errorf("emergencyContactRelationship = %v, want %v", got.emergencyContactRelationship, tt.wantEmergencyContactRelationship)
			}
			if tt.wantEmergencyContactPhone == "" && got.emergencyContactPhone != nil {
				t.Errorf("emergencyContactPhone = %v, want nil", got.emergencyContactPhone)
			}
			if tt.wantEmergencyContactPhone != "" && (got.emergencyContactPhone == nil || got.emergencyContactPhone.String() != tt.wantEmergencyContactPhone) {
				t.Errorf("emergencyContactPhone = %v, want %v", got.emergencyContactPhone, tt.wantEmergencyContactPhone)
			}
			if tt.wantOrganizationID == "" && got.organizationID != nil {
				t.Errorf("organizationID = %v, want nil", got.organizationID)
			}
			if tt.wantOrganizationID != "" && (got.organizationID == nil || got.organizationID.String() != tt.wantOrganizationID) {
				t.Errorf("organizationID = %v, want %v", got.organizationID, tt.wantOrganizationID)
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

func TestToProtoCustomer(t *testing.T) {
	t.Parallel()
	id, _ := customer.NewCustomerIDFromString(sampleCustomerID)
	tenantID, _ := tenant.NewTenantID(sampleTenantID)
	name, _ := customer.NewName("test customer")
	nameKana, _ := customer.NewNameKana("テストカナ")
	gender, _ := customer.NewGender("male")
	birthDate, _ := customer.NewBirthDate(2000, 1, 1)
	phone, _ := contact.NewPhone("03-1234-5678")
	email, _ := contact.NewEmail("test@example.com")
	postalCode, _ := address.NewPostalCode("123-4567")
	prefecture, _ := address.NewPrefecture("東京都")
	city, _ := address.NewCity("千代田区")
	street, _ := address.NewStreet("1-1-1")
	building, _ := address.NewBuilding("テストビル")
	addr := address.NewAddress(postalCode, prefecture, city, street, building)
	emergencyContactName, _ := customer.NewEmergencyContactName("緊急 太郎")
	emergencyContactRelationship, _ := customer.NewEmergencyContactRelationship("父")
	emergencyContactPhone, _ := contact.NewPhone("090-1234-5678")
	organizationID, _ := organization.NewOrganizationIDFromString(sampleOrganizationID)
	now := time.Now().UTC()

	tests := []struct {
		name                             string
		c                                customer.Customer
		wantPhone                        string
		wantEmail                        string
		wantPostalCode                   string
		wantPrefecture                   string
		wantCity                         string
		wantStreet                       string
		wantBuilding                     string
		wantEmergencyContactName         string
		wantEmergencyContactRelationship string
		wantEmergencyContactPhone        string
		wantOrganizationID               string
		wantIsActive                     bool
	}{
		{
			"success to proto customer",
			customer.NewCustomer(id, tenantID, name, nameKana, gender, birthDate, phone, email, addr, emergencyContactName, emergencyContactRelationship, emergencyContactPhone, &organizationID, true, now, now),
			"0312345678", "test@example.com", "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "緊急 太郎", "父", "09012345678", sampleOrganizationID, true,
		},
		{
			"success to proto customer without optional fields",
			customer.NewCustomer(id, tenantID, name, nameKana, gender, birthDate, nil, nil, address.Address{}, nil, nil, nil, nil, false, now, now),
			"", "", "", "", "", "", "", "", "", "", "", false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := toProtoCustomer(tt.c)

			if got.GetCustomerId() != sampleCustomerID {
				t.Errorf("CustomerId = %v, want %v", got.GetCustomerId(), sampleCustomerID)
			}
			if got.GetName() != "test customer" {
				t.Errorf("Name = %v, want test customer", got.GetName())
			}
			if got.GetTenantId() != sampleTenantID {
				t.Errorf("TenantId = %v, want %v", got.GetTenantId(), sampleTenantID)
			}
			if got.GetNameKana() != "テストカナ" {
				t.Errorf("NameKana = %v, want テストカナ", got.GetNameKana())
			}
			if got.GetGender() != customerv1.Gender_GENDER_MALE {
				t.Errorf("Gender = %v, want MALE", got.GetGender())
			}
			if got.GetBirthDate().GetYear() != 2000 || got.GetBirthDate().GetMonth() != 1 || got.GetBirthDate().GetDay() != 1 {
				t.Errorf("BirthDate = %v, want 2000-1-1", got.GetBirthDate())
			}
			if tt.wantPhone == "" && got.Phone != nil {
				t.Errorf("Phone = %v, want nil", got.Phone)
			}
			if tt.wantPhone != "" && got.GetPhone() != tt.wantPhone {
				t.Errorf("Phone = %v, want %v", got.GetPhone(), tt.wantPhone)
			}
			if tt.wantEmail == "" && got.Email != nil {
				t.Errorf("Email = %v, want nil", got.Email)
			}
			if tt.wantEmail != "" && got.GetEmail() != tt.wantEmail {
				t.Errorf("Email = %v, want %v", got.GetEmail(), tt.wantEmail)
			}
			if tt.wantPostalCode == "" && got.PostalCode != nil {
				t.Errorf("PostalCode = %v, want nil", got.PostalCode)
			}
			if tt.wantPostalCode != "" && got.GetPostalCode() != tt.wantPostalCode {
				t.Errorf("PostalCode = %v, want %v", got.GetPostalCode(), tt.wantPostalCode)
			}
			if tt.wantPrefecture == "" && got.Prefecture != nil {
				t.Errorf("Prefecture = %v, want nil", got.Prefecture)
			}
			if tt.wantPrefecture != "" && got.GetPrefecture() != tt.wantPrefecture {
				t.Errorf("Prefecture = %v, want %v", got.GetPrefecture(), tt.wantPrefecture)
			}
			if tt.wantCity == "" && got.City != nil {
				t.Errorf("City = %v, want nil", got.City)
			}
			if tt.wantCity != "" && got.GetCity() != tt.wantCity {
				t.Errorf("City = %v, want %v", got.GetCity(), tt.wantCity)
			}
			if tt.wantStreet == "" && got.Street != nil {
				t.Errorf("Street = %v, want nil", got.Street)
			}
			if tt.wantStreet != "" && got.GetStreet() != tt.wantStreet {
				t.Errorf("Street = %v, want %v", got.GetStreet(), tt.wantStreet)
			}
			if tt.wantBuilding == "" && got.Building != nil {
				t.Errorf("Building = %v, want nil", got.Building)
			}
			if tt.wantBuilding != "" && got.GetBuilding() != tt.wantBuilding {
				t.Errorf("Building = %v, want %v", got.GetBuilding(), tt.wantBuilding)
			}
			if tt.wantEmergencyContactName == "" && got.EmergencyContactName != nil {
				t.Errorf("EmergencyContactName = %v, want nil", got.EmergencyContactName)
			}
			if tt.wantEmergencyContactName != "" && got.GetEmergencyContactName() != tt.wantEmergencyContactName {
				t.Errorf("EmergencyContactName = %v, want %v", got.GetEmergencyContactName(), tt.wantEmergencyContactName)
			}
			if tt.wantEmergencyContactRelationship == "" && got.EmergencyContactRelationship != nil {
				t.Errorf("EmergencyContactRelationship = %v, want nil", got.EmergencyContactRelationship)
			}
			if tt.wantEmergencyContactRelationship != "" && got.GetEmergencyContactRelationship() != tt.wantEmergencyContactRelationship {
				t.Errorf("EmergencyContactRelationship = %v, want %v", got.GetEmergencyContactRelationship(), tt.wantEmergencyContactRelationship)
			}
			if tt.wantEmergencyContactPhone == "" && got.EmergencyContactPhone != nil {
				t.Errorf("EmergencyContactPhone = %v, want nil", got.EmergencyContactPhone)
			}
			if tt.wantEmergencyContactPhone != "" && got.GetEmergencyContactPhone() != tt.wantEmergencyContactPhone {
				t.Errorf("EmergencyContactPhone = %v, want %v", got.GetEmergencyContactPhone(), tt.wantEmergencyContactPhone)
			}
			if tt.wantOrganizationID == "" && got.OrganizationId != nil {
				t.Errorf("OrganizationId = %v, want nil", got.OrganizationId)
			}
			if tt.wantOrganizationID != "" && got.GetOrganizationId() != tt.wantOrganizationID {
				t.Errorf("OrganizationId = %v, want %v", got.GetOrganizationId(), tt.wantOrganizationID)
			}
			if got.GetIsActive() != tt.wantIsActive {
				t.Errorf("IsActive = %v, want %v", got.GetIsActive(), tt.wantIsActive)
			}
		})
	}
}
