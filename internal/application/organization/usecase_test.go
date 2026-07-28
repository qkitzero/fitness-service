package organization

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/qkitzero/fitness-service/internal/application/user"
	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
	mocksappauth "github.com/qkitzero/fitness-service/mocks/application/auth"
	mocksappuser "github.com/qkitzero/fitness-service/mocks/application/user"
	mocksorganization "github.com/qkitzero/fitness-service/mocks/domain/organization"
)

func TestCreateOrganization(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := organization.NewName("テスト株式会社")

	tests := []struct {
		name            string
		success         bool
		wantErr         error
		callCreate      bool
		ctx             context.Context
		userID          string
		verifyTokenErr  error
		myTenantIDs     []string
		listMyGroupsErr error
		createErr       error
	}{
		{"success create organization", true, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{tenantID.String()}, nil, nil},
		{"failure verify token error", false, nil, false, context.Background(), "", errors.New("verify token error"), []string{tenantID.String()}, nil, nil},
		{"failure not tenant member", false, user.ErrNotGroupMember, false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, nil, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, errors.New("list my groups error"), nil},
		{"failure create error", false, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{tenantID.String()}, nil, errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockOrganizationRepository := mocksorganization.NewMockOrganizationRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callCreate {
				mockOrganizationRepository.EXPECT().Create(tt.ctx, gomock.Any()).Return(tt.createErr).Times(1)
			}

			u := NewOrganizationUsecase(mockAuthService, mockUserService, mockOrganizationRepository)

			createdOrganization, err := u.CreateOrganization(tt.ctx, tenantID, name)
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
				if createdOrganization.TenantID() != tenantID {
					t.Errorf("TenantID() = %v, want %v", createdOrganization.TenantID(), tenantID)
				}
				if createdOrganization.Name() != name {
					t.Errorf("Name() = %v, want %v", createdOrganization.Name(), name)
				}
				if createdOrganization.ID().UUID == uuid.Nil {
					t.Errorf("expected generated organization id, but got a nil UUID")
				}
				if !createdOrganization.CreatedAt().Equal(createdOrganization.UpdatedAt()) {
					t.Errorf("CreatedAt() = %v, UpdatedAt() = %v, want equal", createdOrganization.CreatedAt(), createdOrganization.UpdatedAt())
				}
			}
		})
	}
}

func TestGetOrganization(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")

	tests := []struct {
		name            string
		success         bool
		wantErr         error
		callFindByID    bool
		ctx             context.Context
		userID          string
		verifyTokenErr  error
		findByIDErr     error
		myTenantIDs     []string
		listMyGroupsErr error
	}{
		{"success get organization", true, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{tenantID.String()}, nil},
		{"failure verify token error", false, nil, false, context.Background(), "", errors.New("verify token error"), nil, []string{tenantID.String()}, nil},
		{"failure find by id error", false, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, errors.New("find by id error"), []string{tenantID.String()}, nil},
		{"failure organization not found", false, organization.ErrOrganizationNotFound, true, context.Background(), "google-oauth2|000000000000000000000", nil, organization.ErrOrganizationNotFound, []string{tenantID.String()}, nil},
		{"failure other tenant is hidden as not found", false, organization.ErrOrganizationNotFound, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil},
		{"failure list my groups error", false, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, nil, errors.New("list my groups error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			organizationID := organization.NewOrganizationID()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockOrganization := mocksorganization.NewMockOrganization(ctrl)
			mockOrganization.EXPECT().TenantID().Return(tenantID).AnyTimes()
			mockOrganizationRepository := mocksorganization.NewMockOrganizationRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindByID {
				mockOrganizationRepository.EXPECT().FindByID(tt.ctx, organizationID).Return(mockOrganization, tt.findByIDErr).Times(1)
			}

			u := NewOrganizationUsecase(mockAuthService, mockUserService, mockOrganizationRepository)

			foundOrganization, err := u.GetOrganization(tt.ctx, organizationID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.success && foundOrganization != mockOrganization {
				t.Errorf("expected the organization returned by the repository")
			}
		})
	}
}

func TestListOrganizations(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")

	tests := []struct {
		name               string
		success            bool
		wantErr            error
		callListByTenantID bool
		ctx                context.Context
		userID             string
		verifyTokenErr     error
		myTenantIDs        []string
		listMyGroupsErr    error
		listByTenantIDErr  error
	}{
		{"success list organizations", true, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{tenantID.String()}, nil, nil},
		{"failure verify token error", false, nil, false, context.Background(), "", errors.New("verify token error"), []string{tenantID.String()}, nil, nil},
		{"failure not tenant member", false, user.ErrNotGroupMember, false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, nil, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, errors.New("list my groups error"), nil},
		{"failure list by tenant id error", false, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{tenantID.String()}, nil, errors.New("list by group id error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockOrganization := mocksorganization.NewMockOrganization(ctrl)
			mockOrganizationRepository := mocksorganization.NewMockOrganizationRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callListByTenantID {
				mockOrganizationRepository.EXPECT().ListByTenantID(tt.ctx, tenantID).Return([]organization.Organization{mockOrganization}, tt.listByTenantIDErr).Times(1)
			}

			u := NewOrganizationUsecase(mockAuthService, mockUserService, mockOrganizationRepository)

			organizations, err := u.ListOrganizations(tt.ctx, tenantID)
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
				if len(organizations) != 1 {
					t.Errorf("len(organizations) = %v, want %v", len(organizations), 1)
				}
				if len(organizations) == 1 && organizations[0] != mockOrganization {
					t.Errorf("expected the organizations returned by the repository")
				}
			}
		})
	}
}

func TestUpdateOrganization(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := organization.NewName("更新株式会社")

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
		myTenantIDs     []string
		listMyGroupsErr error
		updateErr       error
	}{
		{"success update organization", true, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{tenantID.String()}, nil, nil},
		{"failure verify token error", false, nil, false, false, context.Background(), "", errors.New("verify token error"), nil, []string{tenantID.String()}, nil, nil},
		{"failure find by id error", false, nil, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, errors.New("find by id error"), []string{tenantID.String()}, nil, nil},
		{"failure organization not found", false, organization.ErrOrganizationNotFound, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, organization.ErrOrganizationNotFound, []string{tenantID.String()}, nil, nil},
		{"failure other tenant is hidden as not found", false, organization.ErrOrganizationNotFound, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, nil, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, nil, errors.New("list my groups error"), nil},
		{"failure update error", false, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{tenantID.String()}, nil, errors.New("update error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			organizationID := organization.NewOrganizationID()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockOrganization := mocksorganization.NewMockOrganization(ctrl)
			mockOrganization.EXPECT().TenantID().Return(tenantID).AnyTimes()
			mockOrganizationRepository := mocksorganization.NewMockOrganizationRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindByID {
				mockOrganizationRepository.EXPECT().FindByID(tt.ctx, organizationID).Return(mockOrganization, tt.findByIDErr).Times(1)
			}
			if tt.callUpdate {
				mockOrganization.EXPECT().Update(name).Times(1)
				mockOrganizationRepository.EXPECT().Update(tt.ctx, mockOrganization).Return(tt.updateErr).Times(1)
			}

			u := NewOrganizationUsecase(mockAuthService, mockUserService, mockOrganizationRepository)

			updatedOrganization, err := u.UpdateOrganization(tt.ctx, organizationID, name)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.success && updatedOrganization != mockOrganization {
				t.Errorf("expected the organization returned by the repository")
			}
		})
	}
}

func TestDeleteOrganization(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")

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
		myTenantIDs     []string
		listMyGroupsErr error
		deleteErr       error
	}{
		{"success delete organization", true, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{tenantID.String()}, nil, nil},
		{"failure verify token error", false, nil, false, false, context.Background(), "", errors.New("verify token error"), nil, []string{tenantID.String()}, nil, nil},
		{"failure find by id error", false, nil, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, errors.New("find by id error"), []string{tenantID.String()}, nil, nil},
		{"failure organization not found", false, organization.ErrOrganizationNotFound, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, organization.ErrOrganizationNotFound, []string{tenantID.String()}, nil, nil},
		{"failure other tenant is hidden as not found", false, organization.ErrOrganizationNotFound, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, nil, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, nil, errors.New("list my groups error"), nil},
		{"failure delete error", false, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, []string{tenantID.String()}, nil, errors.New("delete error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			organizationID := organization.NewOrganizationID()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockOrganization := mocksorganization.NewMockOrganization(ctrl)
			mockOrganization.EXPECT().TenantID().Return(tenantID).AnyTimes()
			mockOrganizationRepository := mocksorganization.NewMockOrganizationRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindByID {
				mockOrganizationRepository.EXPECT().FindByID(tt.ctx, organizationID).Return(mockOrganization, tt.findByIDErr).Times(1)
			}
			if tt.callDelete {
				mockOrganizationRepository.EXPECT().Delete(tt.ctx, organizationID).Return(tt.deleteErr).Times(1)
			}

			u := NewOrganizationUsecase(mockAuthService, mockUserService, mockOrganizationRepository)

			err := u.DeleteOrganization(tt.ctx, organizationID)
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
