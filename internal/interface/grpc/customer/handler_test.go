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
	"github.com/qkitzero/fitness-service/internal/application/user"
	"github.com/qkitzero/fitness-service/internal/domain/customer"
	mocksappcustomer "github.com/qkitzero/fitness-service/mocks/application/customer"
	mockscustomer "github.com/qkitzero/fitness-service/mocks/domain/customer"
)

func customerSample(ctrl *gomock.Controller) *mockscustomer.MockCustomer {
	m := mockscustomer.NewMockCustomer(ctrl)
	name, _ := customer.NewName("test customer")
	groupID, _ := customer.NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	nameKana, _ := customer.NewNameKana("テストカナ")
	birthDate, _ := customer.NewBirthDate(2000, 1, 1)
	phone, _ := customer.NewPhone("03-1234-5678")
	email, _ := customer.NewEmail("test@example.com")
	postalCode, _ := customer.NewPostalCode("123-4567")
	prefecture, _ := customer.NewPrefecture("東京都")
	city, _ := customer.NewCity("千代田区")
	street, _ := customer.NewStreet("1-1-1")
	building, _ := customer.NewBuilding("テストビル")
	emergencyContactName, _ := customer.NewEmergencyContactName("緊急 太郎")
	emergencyContactRelationship, _ := customer.NewEmergencyContactRelationship("父")
	emergencyContactPhone, _ := customer.NewPhone("090-1234-5678")
	m.EXPECT().ID().Return(customer.NewCustomerID()).AnyTimes()
	m.EXPECT().Name().Return(name).AnyTimes()
	m.EXPECT().GroupID().Return(groupID).AnyTimes()
	m.EXPECT().NameKana().Return(nameKana).AnyTimes()
	m.EXPECT().Gender().Return(customer.GenderMale).AnyTimes()
	m.EXPECT().BirthDate().Return(birthDate).AnyTimes()
	m.EXPECT().Phone().Return(phone).AnyTimes()
	m.EXPECT().Email().Return(email).AnyTimes()
	m.EXPECT().PostalCode().Return(postalCode).AnyTimes()
	m.EXPECT().Prefecture().Return(prefecture).AnyTimes()
	m.EXPECT().City().Return(city).AnyTimes()
	m.EXPECT().Street().Return(street).AnyTimes()
	m.EXPECT().Building().Return(building).AnyTimes()
	m.EXPECT().EmergencyContactName().Return(emergencyContactName).AnyTimes()
	m.EXPECT().EmergencyContactRelationship().Return(emergencyContactRelationship).AnyTimes()
	m.EXPECT().EmergencyContactPhone().Return(emergencyContactPhone).AnyTimes()
	return m
}

func TestCreateCustomer(t *testing.T) {
	t.Parallel()
	gid := "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"
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
	invalid := "invalid"
	invalidPref := "存在しない県"
	tooLong := strings.Repeat("あ", 256)

	tests := []struct {
		name        string
		req         *customerv1.CreateCustomerRequest
		callUsecase bool
		createErr   error
		wantCode    codes.Code
	}{
		{"success", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, Phone: &phone, Email: &email, PostalCode: &postalCode, Prefecture: &prefecture, City: &city, Street: &street, Building: &building, EmergencyContactName: &ecName, EmergencyContactRelationship: &ecRelationship, EmergencyContactPhone: &ecPhone}, true, nil, codes.OK},
		{"failure invalid group id", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: "", NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid name", &customerv1.CreateCustomerRequest{Name: "", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid name kana", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid gender", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_UNSPECIFIED, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid birth date", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: futureDate}, false, nil, codes.InvalidArgument},
		{"failure invalid phone", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, Phone: &invalid}, false, nil, codes.InvalidArgument},
		{"failure invalid email", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, Email: &invalid}, false, nil, codes.InvalidArgument},
		{"failure invalid postal code", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, PostalCode: &invalid}, false, nil, codes.InvalidArgument},
		{"failure invalid prefecture", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, Prefecture: &invalidPref}, false, nil, codes.InvalidArgument},
		{"failure invalid city", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, City: &tooLong}, false, nil, codes.InvalidArgument},
		{"failure invalid street", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, Street: &tooLong}, false, nil, codes.InvalidArgument},
		{"failure invalid building", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, Building: &tooLong}, false, nil, codes.InvalidArgument},
		{"failure invalid emergency contact name", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, EmergencyContactName: &tooLong}, false, nil, codes.InvalidArgument},
		{"failure invalid emergency contact relationship", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, EmergencyContactRelationship: &tooLong}, false, nil, codes.InvalidArgument},
		{"failure invalid emergency contact phone", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate, EmergencyContactPhone: &invalid}, false, nil, codes.InvalidArgument},
		{"failure not group member", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure usecase error", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, true, fmt.Errorf("create customer error"), codes.Internal},
		{"failure status preserved", &customerv1.CreateCustomerRequest{Name: "test customer", GroupId: gid, NameKana: "テストカナ", Gender: customerv1.Gender_GENDER_MALE, BirthDate: validDate}, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
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
				mockUsecase.EXPECT().CreateCustomer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(customerSample(ctrl), tt.createErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			_, err := handler.CreateCustomer(ctx, tt.req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
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
		{"success get customer", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, nil, codes.OK},
		{"failure invalid customer id", "", false, nil, codes.InvalidArgument},
		{"failure get customer error", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, fmt.Errorf("get customer error"), codes.Internal},
		{"failure not group member", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure customer not found", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, customer.ErrCustomerNotFound, codes.NotFound},
		{"failure status preserved", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
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
				mockUsecase.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).Return(customerSample(ctrl), tt.getCustomerErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			_, err := handler.GetCustomer(ctx, &customerv1.GetCustomerRequest{CustomerId: tt.customerID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestListCustomers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		groupID          string
		callUsecase      bool
		listCustomersErr error
		wantCode         codes.Code
	}{
		{"success list customers", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", true, nil, codes.OK},
		{"failure invalid group id", "", false, nil, codes.InvalidArgument},
		{"failure not group member", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure usecase error", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", true, fmt.Errorf("list customers error"), codes.Internal},
		{"failure status preserved", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
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
				customers := []customer.Customer{customerSample(ctrl), customerSample(ctrl)}
				mockUsecase.EXPECT().ListCustomers(gomock.Any(), gomock.Any()).Return(customers, tt.listCustomersErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			res, err := handler.ListCustomers(ctx, &customerv1.ListCustomersRequest{GroupId: tt.groupID})
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
	cid := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"
	validDate := &date.Date{Year: 1999, Month: 12, Day: 31}
	futureDate := &date.Date{Year: 2500, Month: 1, Day: 1}
	invalid := "invalid"

	tests := []struct {
		name        string
		req         *customerv1.UpdateCustomerRequest
		callUsecase bool
		updateErr   error
		wantCode    codes.Code
	}{
		{"success update customer", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, true, nil, codes.OK},
		{"failure invalid customer id", &customerv1.UpdateCustomerRequest{CustomerId: "", Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid name", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid name kana", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid gender", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_UNSPECIFIED, BirthDate: validDate}, false, nil, codes.InvalidArgument},
		{"failure invalid birth date", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: futureDate}, false, nil, codes.InvalidArgument},
		{"failure invalid phone", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate, Phone: &invalid}, false, nil, codes.InvalidArgument},
		{"failure usecase error", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, true, fmt.Errorf("update customer error"), codes.Internal},
		{"failure not group member", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure customer not found", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, true, customer.ErrCustomerNotFound, codes.NotFound},
		{"failure status preserved", &customerv1.UpdateCustomerRequest{CustomerId: cid, Name: "updated test customer", NameKana: "コウシンカナ", Gender: customerv1.Gender_GENDER_FEMALE, BirthDate: validDate}, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
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
				mockUsecase.EXPECT().UpdateCustomer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(customerSample(ctrl), tt.updateErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			_, err := handler.UpdateCustomer(ctx, tt.req)
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
		customerID        string
		callUsecase       bool
		deleteCustomerErr error
		wantCode          codes.Code
	}{
		{"success delete customer", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, nil, codes.OK},
		{"failure invalid customer id", "", false, nil, codes.InvalidArgument},
		{"failure usecase error", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, fmt.Errorf("delete customer error"), codes.Internal},
		{"failure not group member", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure customer not found", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, customer.ErrCustomerNotFound, codes.NotFound},
		{"failure status preserved", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
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
				mockUsecase.EXPECT().DeleteCustomer(gomock.Any(), gomock.Any()).Return(tt.deleteCustomerErr).Times(1)
			}

			handler := NewCustomerHandler(mockUsecase)

			_, err := handler.DeleteCustomer(ctx, &customerv1.DeleteCustomerRequest{CustomerId: tt.customerID})
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

func TestToProtoCustomer(t *testing.T) {
	t.Parallel()
	id, _ := customer.NewCustomerIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	groupID, _ := customer.NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := customer.NewName("test customer")
	nameKana, _ := customer.NewNameKana("テストカナ")
	gender, _ := customer.NewGender("male")
	birthDate, _ := customer.NewBirthDate(2000, 1, 1)
	phone, _ := customer.NewPhone("03-1234-5678")
	email, _ := customer.NewEmail("test@example.com")
	postalCode, _ := customer.NewPostalCode("123-4567")
	prefecture, _ := customer.NewPrefecture("東京都")
	city, _ := customer.NewCity("千代田区")
	street, _ := customer.NewStreet("1-1-1")
	building, _ := customer.NewBuilding("テストビル")
	emergencyContactName, _ := customer.NewEmergencyContactName("緊急 太郎")
	emergencyContactRelationship, _ := customer.NewEmergencyContactRelationship("父")
	emergencyContactPhone, _ := customer.NewPhone("090-1234-5678")
	now := time.Now()

	full := customer.NewCustomer(id, groupID, name, nameKana, gender, birthDate, phone, email, postalCode, prefecture, city, street, building, emergencyContactName, emergencyContactRelationship, emergencyContactPhone, now, now)
	got := toProtoCustomer(full)
	if got.GetCustomerId() != id.String() {
		t.Errorf("CustomerId = %v, want %v", got.GetCustomerId(), id.String())
	}
	if got.GetName() != "test customer" {
		t.Errorf("Name = %v, want test customer", got.GetName())
	}
	if got.GetGroupId() != groupID.String() {
		t.Errorf("GroupId = %v, want %v", got.GetGroupId(), groupID.String())
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
	if got.GetPhone() != "0312345678" {
		t.Errorf("Phone = %v, want 0312345678", got.GetPhone())
	}
	if got.GetEmail() != "test@example.com" {
		t.Errorf("Email = %v, want test@example.com", got.GetEmail())
	}
	if got.GetPostalCode() != "1234567" {
		t.Errorf("PostalCode = %v, want 1234567", got.GetPostalCode())
	}
	if got.GetPrefecture() != "東京都" {
		t.Errorf("Prefecture = %v, want 東京都", got.GetPrefecture())
	}
	if got.GetCity() != "千代田区" {
		t.Errorf("City = %v, want 千代田区", got.GetCity())
	}
	if got.GetStreet() != "1-1-1" {
		t.Errorf("Street = %v, want 1-1-1", got.GetStreet())
	}
	if got.GetBuilding() != "テストビル" {
		t.Errorf("Building = %v, want テストビル", got.GetBuilding())
	}
	if got.GetEmergencyContactName() != "緊急 太郎" {
		t.Errorf("EmergencyContactName = %v, want 緊急 太郎", got.GetEmergencyContactName())
	}
	if got.GetEmergencyContactRelationship() != "父" {
		t.Errorf("EmergencyContactRelationship = %v, want 父", got.GetEmergencyContactRelationship())
	}
	if got.GetEmergencyContactPhone() != "09012345678" {
		t.Errorf("EmergencyContactPhone = %v, want 09012345678", got.GetEmergencyContactPhone())
	}

	empty := customer.NewCustomer(id, groupID, name, nameKana, gender, birthDate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, now, now)
	gotEmpty := toProtoCustomer(empty)
	if gotEmpty.Phone != nil || gotEmpty.Email != nil || gotEmpty.PostalCode != nil || gotEmpty.Prefecture != nil || gotEmpty.City != nil || gotEmpty.Street != nil || gotEmpty.Building != nil || gotEmpty.EmergencyContactName != nil || gotEmpty.EmergencyContactRelationship != nil || gotEmpty.EmergencyContactPhone != nil {
		t.Errorf("expected nil optional proto fields for empty customer")
	}
}
