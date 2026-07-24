package customer

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
	mocksappauth "github.com/qkitzero/fitness-service/mocks/application/auth"
	mocksappuser "github.com/qkitzero/fitness-service/mocks/application/user"
	mockscustomer "github.com/qkitzero/fitness-service/mocks/domain/customer"
)

func TestCreateCustomer(t *testing.T) {
	t.Parallel()
	groupID, _ := customer.NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := customer.NewName("test customer")
	nameKana, _ := customer.NewNameKana("テストカナ")
	gender, _ := customer.NewGender("male")
	birthDate, _ := customer.NewBirthDate(2000, 1, 1)

	tests := []struct {
		name            string
		success         bool
		ctx             context.Context
		userID          string
		verifyTokenErr  error
		myGroupIDs      []string
		listMyGroupsErr error
		createErr       error
	}{
		{"success create customer", true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{groupID.String()}, nil, nil},
		{"failure verify token error", false, context.Background(), "", errors.New("verify token error"), []string{groupID.String()}, nil, nil},
		{"failure not group member", false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, errors.New("list my groups error"), nil},
		{"failure create error", false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{groupID.String()}, nil, errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myGroupIDs, tt.listMyGroupsErr).AnyTimes()
			mockCustomerRepository.EXPECT().Create(tt.ctx, gomock.Any()).Return(tt.createErr).AnyTimes()

			u := NewCustomerUsecase(mockAuthService, mockUserService, mockCustomerRepository)

			_, err := u.CreateCustomer(tt.ctx, groupID, name, nameKana, gender, birthDate)
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
	groupID, _ := customer.NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")

	tests := []struct {
		name            string
		success         bool
		ctx             context.Context
		userID          string
		verifyTokenErr  error
		findByIDErr     error
		myGroupIDs      []string
		listMyGroupsErr error
	}{
		{"success get customer", true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{groupID.String()}, nil},
		{"failure verify token error", false, context.Background(), "", errors.New("verify token error"), nil, []string{groupID.String()}, nil},
		{"failure find by id error", false, context.Background(), "google-oauth2|000000000000000000000", nil, errors.New("find by id error"), []string{groupID.String()}, nil},
		{"failure customer not found", false, context.Background(), "google-oauth2|000000000000000000000", nil, customer.ErrCustomerNotFound, []string{groupID.String()}, nil},
		{"failure not group member", false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil},
		{"failure list my groups error", false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, nil, errors.New("list my groups error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().GroupID().Return(groupID).AnyTimes()
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockCustomerRepository.EXPECT().FindByID(tt.ctx, gomock.Any()).Return(mockCustomer, tt.findByIDErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myGroupIDs, tt.listMyGroupsErr).AnyTimes()

			u := NewCustomerUsecase(mockAuthService, mockUserService, mockCustomerRepository)

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

func TestListCustomers(t *testing.T) {
	t.Parallel()
	groupID, _ := customer.NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")

	tests := []struct {
		name             string
		success          bool
		ctx              context.Context
		userID           string
		verifyTokenErr   error
		myGroupIDs       []string
		listMyGroupsErr  error
		listByGroupIDErr error
	}{
		{"success list customers", true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{groupID.String()}, nil, nil},
		{"failure verify token error", false, context.Background(), "", errors.New("verify token error"), []string{groupID.String()}, nil, nil},
		{"failure not group member", false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, errors.New("list my groups error"), nil},
		{"failure list by group id error", false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{groupID.String()}, nil, errors.New("list by group id error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myGroupIDs, tt.listMyGroupsErr).AnyTimes()
			mockCustomerRepository.EXPECT().ListByGroupID(tt.ctx, groupID).Return([]customer.Customer{mockCustomer}, tt.listByGroupIDErr).AnyTimes()

			u := NewCustomerUsecase(mockAuthService, mockUserService, mockCustomerRepository)

			_, err := u.ListCustomers(tt.ctx, groupID)
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
	groupID, _ := customer.NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := customer.NewName("updated test customer")
	nameKana, _ := customer.NewNameKana("コウシンカナ")
	gender, _ := customer.NewGender("female")
	birthDate, _ := customer.NewBirthDate(1999, 12, 31)

	tests := []struct {
		name            string
		success         bool
		ctx             context.Context
		userID          string
		verifyTokenErr  error
		findByIDErr     error
		myGroupIDs      []string
		listMyGroupsErr error
		updateErr       error
	}{
		{"success update customer", true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{groupID.String()}, nil, nil},
		{"failure verify token error", false, context.Background(), "", errors.New("verify token error"), nil, []string{groupID.String()}, nil, nil},
		{"failure find by id error", false, context.Background(), "google-oauth2|000000000000000000000", nil, errors.New("find by id error"), []string{groupID.String()}, nil, nil},
		{"failure customer not found", false, context.Background(), "google-oauth2|000000000000000000000", nil, customer.ErrCustomerNotFound, []string{groupID.String()}, nil, nil},
		{"failure not group member", false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, nil, errors.New("list my groups error"), nil},
		{"failure update error", false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{groupID.String()}, nil, errors.New("update error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().GroupID().Return(groupID).AnyTimes()
			mockCustomer.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockCustomerRepository.EXPECT().FindByID(tt.ctx, gomock.Any()).Return(mockCustomer, tt.findByIDErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myGroupIDs, tt.listMyGroupsErr).AnyTimes()
			mockCustomerRepository.EXPECT().Update(tt.ctx, gomock.Any()).Return(tt.updateErr).AnyTimes()

			u := NewCustomerUsecase(mockAuthService, mockUserService, mockCustomerRepository)

			_, err := u.UpdateCustomer(tt.ctx, customer.NewCustomerID(), name, nameKana, gender, birthDate)
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
	groupID, _ := customer.NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")

	tests := []struct {
		name            string
		success         bool
		ctx             context.Context
		userID          string
		verifyTokenErr  error
		findByIDErr     error
		myGroupIDs      []string
		listMyGroupsErr error
		deleteErr       error
	}{
		{"success delete customer", true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{groupID.String()}, nil, nil},
		{"failure verify token error", false, context.Background(), "", errors.New("verify token error"), nil, []string{groupID.String()}, nil, nil},
		{"failure find by id error", false, context.Background(), "google-oauth2|000000000000000000000", nil, errors.New("find by id error"), []string{groupID.String()}, nil, nil},
		{"failure customer not found", false, context.Background(), "google-oauth2|000000000000000000000", nil, customer.ErrCustomerNotFound, []string{groupID.String()}, nil, nil},
		{"failure not group member", false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, nil, errors.New("list my groups error"), nil},
		{"failure delete error", false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{groupID.String()}, nil, errors.New("delete error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().GroupID().Return(groupID).AnyTimes()
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockCustomerRepository.EXPECT().FindByID(tt.ctx, gomock.Any()).Return(mockCustomer, tt.findByIDErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myGroupIDs, tt.listMyGroupsErr).AnyTimes()
			mockCustomerRepository.EXPECT().Delete(tt.ctx, gomock.Any()).Return(tt.deleteErr).AnyTimes()

			u := NewCustomerUsecase(mockAuthService, mockUserService, mockCustomerRepository)

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
