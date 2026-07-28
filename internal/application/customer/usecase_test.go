package customer

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/qkitzero/fitness-service/internal/application/user"
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
		wantErr         error
		callCreate      bool
		ctx             context.Context
		userID          string
		verifyTokenErr  error
		myGroupIDs      []string
		listMyGroupsErr error
		createErr       error
	}{
		{"success create customer", true, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{groupID.String()}, nil, nil},
		{"failure verify token error", false, nil, false, context.Background(), "", errors.New("verify token error"), []string{groupID.String()}, nil, nil},
		{"failure not group member", false, user.ErrNotGroupMember, false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, nil, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, errors.New("list my groups error"), nil},
		{"failure create error", false, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{groupID.String()}, nil, errors.New("create error")},
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
			if tt.callCreate {
				mockCustomerRepository.EXPECT().Create(tt.ctx, gomock.Any()).Return(tt.createErr).Times(1)
			}

			u := NewCustomerUsecase(mockAuthService, mockUserService, mockCustomerRepository)

			createdCustomer, err := u.CreateCustomer(tt.ctx, groupID, name, nameKana, gender, birthDate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.success {
				if createdCustomer.GroupID() != groupID {
					t.Errorf("GroupID() = %v, want %v", createdCustomer.GroupID(), groupID)
				}
				if createdCustomer.Name() != name {
					t.Errorf("Name() = %v, want %v", createdCustomer.Name(), name)
				}
				if createdCustomer.ID().UUID == uuid.Nil {
					t.Errorf("expected generated customer id, but got a nil UUID")
				}
				if !createdCustomer.CreatedAt().Equal(createdCustomer.UpdatedAt()) {
					t.Errorf("CreatedAt() = %v, UpdatedAt() = %v, want equal", createdCustomer.CreatedAt(), createdCustomer.UpdatedAt())
				}
				if !createdCustomer.IsActive() {
					t.Errorf("IsActive() = %v, want %v", createdCustomer.IsActive(), true)
				}
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
		wantErr         error
		callFindByID    bool
		ctx             context.Context
		userID          string
		verifyTokenErr  error
		findByIDErr     error
		myGroupIDs      []string
		listMyGroupsErr error
	}{
		{"success get customer", true, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{groupID.String()}, nil},
		{"failure verify token error", false, nil, false, context.Background(), "", errors.New("verify token error"), nil, []string{groupID.String()}, nil},
		{"failure find by id error", false, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, errors.New("find by id error"), []string{groupID.String()}, nil},
		{"failure customer not found", false, customer.ErrCustomerNotFound, true, context.Background(), "google-oauth2|000000000000000000000", nil, customer.ErrCustomerNotFound, []string{groupID.String()}, nil},
		{"failure other group is hidden as not found", false, customer.ErrCustomerNotFound, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil},
		{"failure list my groups error", false, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, nil, errors.New("list my groups error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			customerID := customer.NewCustomerID()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().GroupID().Return(groupID).AnyTimes()
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myGroupIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindByID {
				mockCustomerRepository.EXPECT().FindByID(tt.ctx, customerID).Return(mockCustomer, tt.findByIDErr).Times(1)
			}

			u := NewCustomerUsecase(mockAuthService, mockUserService, mockCustomerRepository)

			foundCustomer, err := u.GetCustomer(tt.ctx, customerID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.success && foundCustomer != mockCustomer {
				t.Errorf("expected the customer returned by the repository")
			}
		})
	}
}

func TestListCustomers(t *testing.T) {
	t.Parallel()
	groupID, _ := customer.NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")

	tests := []struct {
		name              string
		success           bool
		wantErr           error
		callListByGroupID bool
		includeInactive   bool
		ctx               context.Context
		userID            string
		verifyTokenErr    error
		myGroupIDs        []string
		listMyGroupsErr   error
		listByGroupIDErr  error
	}{
		{"success list customers", true, nil, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{groupID.String()}, nil, nil},
		{"success list customers including inactive", true, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{groupID.String()}, nil, nil},
		{"failure verify token error", false, nil, false, false, context.Background(), "", errors.New("verify token error"), []string{groupID.String()}, nil, nil},
		{"failure not group member", false, user.ErrNotGroupMember, false, false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, nil, false, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, errors.New("list my groups error"), nil},
		{"failure list by group id error", false, nil, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{groupID.String()}, nil, errors.New("list by group id error")},
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
			if tt.callListByGroupID {
				mockCustomerRepository.EXPECT().ListByGroupID(tt.ctx, groupID, tt.includeInactive).Return([]customer.Customer{mockCustomer}, tt.listByGroupIDErr).Times(1)
			}

			u := NewCustomerUsecase(mockAuthService, mockUserService, mockCustomerRepository)

			customers, err := u.ListCustomers(tt.ctx, groupID, tt.includeInactive)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.success {
				if len(customers) != 1 {
					t.Errorf("len(customers) = %v, want %v", len(customers), 1)
				}
				if len(customers) == 1 && customers[0] != mockCustomer {
					t.Errorf("expected the customers returned by the repository")
				}
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
		wantErr         error
		callFindByID    bool
		callUpdate      bool
		ctx             context.Context
		userID          string
		verifyTokenErr  error
		findByIDErr     error
		myGroupIDs      []string
		listMyGroupsErr error
		updateErr       error
	}{
		{"success update customer", true, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{groupID.String()}, nil, nil},
		{"failure verify token error", false, nil, false, false, context.Background(), "", errors.New("verify token error"), nil, []string{groupID.String()}, nil, nil},
		{"failure find by id error", false, nil, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, errors.New("find by id error"), []string{groupID.String()}, nil, nil},
		{"failure customer not found", false, customer.ErrCustomerNotFound, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, customer.ErrCustomerNotFound, []string{groupID.String()}, nil, nil},
		{"failure other group is hidden as not found", false, customer.ErrCustomerNotFound, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, nil, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, nil, errors.New("list my groups error"), nil},
		{"failure update error", false, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{groupID.String()}, nil, errors.New("update error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			customerID := customer.NewCustomerID()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().GroupID().Return(groupID).AnyTimes()
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myGroupIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindByID {
				mockCustomerRepository.EXPECT().FindByID(tt.ctx, customerID).Return(mockCustomer, tt.findByIDErr).Times(1)
			}
			if tt.callUpdate {
				mockCustomer.EXPECT().Update(name, nameKana, gender, birthDate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil).Times(1)
				mockCustomerRepository.EXPECT().Update(tt.ctx, mockCustomer).Return(tt.updateErr).Times(1)
			}

			u := NewCustomerUsecase(mockAuthService, mockUserService, mockCustomerRepository)

			updatedCustomer, err := u.UpdateCustomer(tt.ctx, customerID, name, nameKana, gender, birthDate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.success && updatedCustomer != mockCustomer {
				t.Errorf("expected the customer returned by the repository")
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
		wantErr         error
		callFindByID    bool
		callDelete      bool
		ctx             context.Context
		userID          string
		verifyTokenErr  error
		findByIDErr     error
		myGroupIDs      []string
		listMyGroupsErr error
		deleteErr       error
	}{
		{"success delete customer", true, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{groupID.String()}, nil, nil},
		{"failure verify token error", false, nil, false, false, context.Background(), "", errors.New("verify token error"), nil, []string{groupID.String()}, nil, nil},
		{"failure find by id error", false, nil, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, errors.New("find by id error"), []string{groupID.String()}, nil, nil},
		{"failure customer not found", false, customer.ErrCustomerNotFound, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, customer.ErrCustomerNotFound, []string{groupID.String()}, nil, nil},
		{"failure other group is hidden as not found", false, customer.ErrCustomerNotFound, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, nil, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, nil, errors.New("list my groups error"), nil},
		{"failure delete error", false, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{groupID.String()}, nil, errors.New("delete error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			customerID := customer.NewCustomerID()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().GroupID().Return(groupID).AnyTimes()
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myGroupIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindByID {
				mockCustomerRepository.EXPECT().FindByID(tt.ctx, customerID).Return(mockCustomer, tt.findByIDErr).Times(1)
			}
			if tt.callDelete {
				mockCustomerRepository.EXPECT().Delete(tt.ctx, customerID).Return(tt.deleteErr).Times(1)
			}

			u := NewCustomerUsecase(mockAuthService, mockUserService, mockCustomerRepository)

			err := u.DeleteCustomer(tt.ctx, customerID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
