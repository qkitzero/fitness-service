package organization

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	organizationv1 "github.com/qkitzero/fitness-service/gen/go/organization/v1"
	"github.com/qkitzero/fitness-service/internal/application/user"
	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
	mocksapporganization "github.com/qkitzero/fitness-service/mocks/application/organization"
	mocksorganization "github.com/qkitzero/fitness-service/mocks/domain/organization"
)

func TestCreateOrganization(t *testing.T) {
	t.Parallel()
	gid := "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"
	oid := "3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f"
	tooLong := strings.Repeat("あ", 256)

	tests := []struct {
		name        string
		req         *organizationv1.CreateOrganizationRequest
		callUsecase bool
		createErr   error
		wantCode    codes.Code
	}{
		{"success create organization", &organizationv1.CreateOrganizationRequest{GroupId: gid, Name: "テスト株式会社"}, true, nil, codes.OK},
		{"failure invalid group id", &organizationv1.CreateOrganizationRequest{GroupId: "", Name: "テスト株式会社"}, false, nil, codes.InvalidArgument},
		{"failure invalid name", &organizationv1.CreateOrganizationRequest{GroupId: gid, Name: ""}, false, nil, codes.InvalidArgument},
		{"failure too long name", &organizationv1.CreateOrganizationRequest{GroupId: gid, Name: tooLong}, false, nil, codes.InvalidArgument},
		{"failure null character name", &organizationv1.CreateOrganizationRequest{GroupId: gid, Name: "テスト\x00株式会社"}, false, nil, codes.InvalidArgument},
		{"failure not tenant member", &organizationv1.CreateOrganizationRequest{GroupId: gid, Name: "テスト株式会社"}, true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure usecase error", &organizationv1.CreateOrganizationRequest{GroupId: gid, Name: "テスト株式会社"}, true, fmt.Errorf("create organization error"), codes.Internal},
		{"failure unauthenticated is preserved", &organizationv1.CreateOrganizationRequest{GroupId: gid, Name: "テスト株式会社"}, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure downstream code is not forwarded", &organizationv1.CreateOrganizationRequest{GroupId: gid, Name: "テスト株式会社"}, true, status.Error(codes.NotFound, "user not found"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockUsecase := mocksapporganization.NewMockOrganizationUsecase(ctrl)
			if tt.callUsecase {
				organizationID, _ := organization.NewOrganizationIDFromString(oid)
				tenantID, _ := tenant.NewTenantID(gid)
				name, _ := organization.NewName("テスト株式会社")
				mockOrganization := mocksorganization.NewMockOrganization(ctrl)
				mockOrganization.EXPECT().ID().Return(organizationID).AnyTimes()
				mockOrganization.EXPECT().TenantID().Return(tenantID).AnyTimes()
				mockOrganization.EXPECT().Name().Return(name).AnyTimes()
				mockUsecase.EXPECT().CreateOrganization(gomock.Any(), tenantID, name).Return(mockOrganization, tt.createErr).Times(1)
			}

			handler := NewOrganizationHandler(mockUsecase)

			res, err := handler.CreateOrganization(ctx, tt.req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK && res.GetOrganizationId() != oid {
				t.Errorf("OrganizationId = %v, want %v", res.GetOrganizationId(), oid)
			}
		})
	}
}

func TestGetOrganization(t *testing.T) {
	t.Parallel()
	gid := "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"
	oid := "3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f"

	tests := []struct {
		name               string
		organizationID     string
		callUsecase        bool
		getOrganizationErr error
		wantCode           codes.Code
	}{
		{"success get organization", oid, true, nil, codes.OK},
		{"failure invalid organization id", "", false, nil, codes.InvalidArgument},
		{"failure usecase error", oid, true, fmt.Errorf("get organization error"), codes.Internal},
		{"failure not tenant member", oid, true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure organization not found", oid, true, organization.ErrOrganizationNotFound, codes.NotFound},
		{"failure unauthenticated is preserved", oid, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure downstream code is not forwarded", oid, true, status.Error(codes.InvalidArgument, "user service"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockUsecase := mocksapporganization.NewMockOrganizationUsecase(ctrl)
			if tt.callUsecase {
				organizationID, _ := organization.NewOrganizationIDFromString(oid)
				tenantID, _ := tenant.NewTenantID(gid)
				name, _ := organization.NewName("テスト株式会社")
				mockOrganization := mocksorganization.NewMockOrganization(ctrl)
				mockOrganization.EXPECT().ID().Return(organizationID).AnyTimes()
				mockOrganization.EXPECT().TenantID().Return(tenantID).AnyTimes()
				mockOrganization.EXPECT().Name().Return(name).AnyTimes()
				mockUsecase.EXPECT().GetOrganization(gomock.Any(), organizationID).Return(mockOrganization, tt.getOrganizationErr).Times(1)
			}

			handler := NewOrganizationHandler(mockUsecase)

			res, err := handler.GetOrganization(ctx, &organizationv1.GetOrganizationRequest{OrganizationId: tt.organizationID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				if res.GetOrganization().GetOrganizationId() != oid {
					t.Errorf("OrganizationId = %v, want %v", res.GetOrganization().GetOrganizationId(), oid)
				}
				if res.GetOrganization().GetGroupId() != gid {
					t.Errorf("GroupId = %v, want %v", res.GetOrganization().GetGroupId(), gid)
				}
				if res.GetOrganization().GetName() != "テスト株式会社" {
					t.Errorf("Name = %v, want テスト株式会社", res.GetOrganization().GetName())
				}
			}
		})
	}
}

func TestListOrganizations(t *testing.T) {
	t.Parallel()
	gid := "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"
	wantNames := []string{"テスト株式会社", "テスト工業株式会社"}

	tests := []struct {
		name                 string
		tenantID             string
		callUsecase          bool
		listOrganizationsErr error
		wantCode             codes.Code
	}{
		{"success list organizations", gid, true, nil, codes.OK},
		{"failure invalid group id", "", false, nil, codes.InvalidArgument},
		{"failure not tenant member", gid, true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure usecase error", gid, true, fmt.Errorf("list organizations error"), codes.Internal},
		{"failure unauthenticated is preserved", gid, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure downstream code is not forwarded", gid, true, status.Error(codes.NotFound, "user not found"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockUsecase := mocksapporganization.NewMockOrganizationUsecase(ctrl)
			if tt.callUsecase {
				tenantID, _ := tenant.NewTenantID(gid)
				organizations := make([]organization.Organization, 0, len(wantNames))
				for _, wantName := range wantNames {
					name, _ := organization.NewName(wantName)
					mockOrganization := mocksorganization.NewMockOrganization(ctrl)
					mockOrganization.EXPECT().ID().Return(organization.NewOrganizationID()).AnyTimes()
					mockOrganization.EXPECT().TenantID().Return(tenantID).AnyTimes()
					mockOrganization.EXPECT().Name().Return(name).AnyTimes()
					organizations = append(organizations, mockOrganization)
				}
				mockUsecase.EXPECT().ListOrganizations(gomock.Any(), tenantID).Return(organizations, tt.listOrganizationsErr).Times(1)
			}

			handler := NewOrganizationHandler(mockUsecase)

			res, err := handler.ListOrganizations(ctx, &organizationv1.ListOrganizationsRequest{GroupId: tt.tenantID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				if len(res.GetOrganizations()) != len(wantNames) {
					t.Errorf("len(organizations) = %v, want %v", len(res.GetOrganizations()), len(wantNames))
				}
				for i, wantName := range wantNames {
					if i >= len(res.GetOrganizations()) {
						break
					}
					if res.GetOrganizations()[i].GetName() != wantName {
						t.Errorf("organizations[%d].Name = %v, want %v", i, res.GetOrganizations()[i].GetName(), wantName)
					}
					if res.GetOrganizations()[i].GetGroupId() != gid {
						t.Errorf("organizations[%d].GroupId = %v, want %v", i, res.GetOrganizations()[i].GetGroupId(), gid)
					}
				}
			}
		})
	}
}

func TestUpdateOrganization(t *testing.T) {
	t.Parallel()
	gid := "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"
	oid := "3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f"
	tooLong := strings.Repeat("あ", 256)

	tests := []struct {
		name        string
		req         *organizationv1.UpdateOrganizationRequest
		callUsecase bool
		updateErr   error
		wantCode    codes.Code
	}{
		{"success update organization", &organizationv1.UpdateOrganizationRequest{OrganizationId: oid, Name: "更新株式会社"}, true, nil, codes.OK},
		{"failure invalid organization id", &organizationv1.UpdateOrganizationRequest{OrganizationId: "", Name: "更新株式会社"}, false, nil, codes.InvalidArgument},
		{"failure invalid name", &organizationv1.UpdateOrganizationRequest{OrganizationId: oid, Name: ""}, false, nil, codes.InvalidArgument},
		{"failure too long name", &organizationv1.UpdateOrganizationRequest{OrganizationId: oid, Name: tooLong}, false, nil, codes.InvalidArgument},
		{"failure null character name", &organizationv1.UpdateOrganizationRequest{OrganizationId: oid, Name: "更新\x00株式会社"}, false, nil, codes.InvalidArgument},
		{"failure usecase error", &organizationv1.UpdateOrganizationRequest{OrganizationId: oid, Name: "更新株式会社"}, true, fmt.Errorf("update organization error"), codes.Internal},
		{"failure not tenant member", &organizationv1.UpdateOrganizationRequest{OrganizationId: oid, Name: "更新株式会社"}, true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure organization not found", &organizationv1.UpdateOrganizationRequest{OrganizationId: oid, Name: "更新株式会社"}, true, organization.ErrOrganizationNotFound, codes.NotFound},
		{"failure unauthenticated is preserved", &organizationv1.UpdateOrganizationRequest{OrganizationId: oid, Name: "更新株式会社"}, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure downstream code is not forwarded", &organizationv1.UpdateOrganizationRequest{OrganizationId: oid, Name: "更新株式会社"}, true, status.Error(codes.InvalidArgument, "user service"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockUsecase := mocksapporganization.NewMockOrganizationUsecase(ctrl)
			if tt.callUsecase {
				organizationID, _ := organization.NewOrganizationIDFromString(oid)
				tenantID, _ := tenant.NewTenantID(gid)
				name, _ := organization.NewName("更新株式会社")
				mockOrganization := mocksorganization.NewMockOrganization(ctrl)
				mockOrganization.EXPECT().ID().Return(organizationID).AnyTimes()
				mockOrganization.EXPECT().TenantID().Return(tenantID).AnyTimes()
				mockOrganization.EXPECT().Name().Return(name).AnyTimes()
				mockUsecase.EXPECT().UpdateOrganization(gomock.Any(), organizationID, name).Return(mockOrganization, tt.updateErr).Times(1)
			}

			handler := NewOrganizationHandler(mockUsecase)

			res, err := handler.UpdateOrganization(ctx, tt.req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				if res.GetOrganization().GetOrganizationId() != oid {
					t.Errorf("OrganizationId = %v, want %v", res.GetOrganization().GetOrganizationId(), oid)
				}
				if res.GetOrganization().GetName() != "更新株式会社" {
					t.Errorf("Name = %v, want 更新株式会社", res.GetOrganization().GetName())
				}
			}
		})
	}
}

func TestDeleteOrganization(t *testing.T) {
	t.Parallel()
	oid := "3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f"

	tests := []struct {
		name                  string
		organizationID        string
		callUsecase           bool
		deleteOrganizationErr error
		wantCode              codes.Code
	}{
		{"success delete organization", oid, true, nil, codes.OK},
		{"failure invalid organization id", "", false, nil, codes.InvalidArgument},
		{"failure usecase error", oid, true, fmt.Errorf("delete organization error"), codes.Internal},
		{"failure not tenant member", oid, true, user.ErrNotGroupMember, codes.PermissionDenied},
		{"failure organization not found", oid, true, organization.ErrOrganizationNotFound, codes.NotFound},
		{"failure unauthenticated is preserved", oid, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure downstream code is not forwarded", oid, true, status.Error(codes.NotFound, "user not found"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockUsecase := mocksapporganization.NewMockOrganizationUsecase(ctrl)
			if tt.callUsecase {
				organizationID, _ := organization.NewOrganizationIDFromString(oid)
				mockUsecase.EXPECT().DeleteOrganization(gomock.Any(), organizationID).Return(tt.deleteOrganizationErr).Times(1)
			}

			handler := NewOrganizationHandler(mockUsecase)

			_, err := handler.DeleteOrganization(ctx, &organizationv1.DeleteOrganizationRequest{OrganizationId: tt.organizationID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestToProtoOrganization(t *testing.T) {
	t.Parallel()
	id, _ := organization.NewOrganizationIDFromString("3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f")
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := organization.NewName("テスト株式会社")
	now := time.Now().UTC()

	got := toProtoOrganization(organization.NewOrganization(id, tenantID, name, now, now))
	if got.GetOrganizationId() != id.String() {
		t.Errorf("OrganizationId = %v, want %v", got.GetOrganizationId(), id.String())
	}
	if got.GetGroupId() != tenantID.String() {
		t.Errorf("GroupId = %v, want %v", got.GetGroupId(), tenantID.String())
	}
	if got.GetName() != "テスト株式会社" {
		t.Errorf("Name = %v, want テスト株式会社", got.GetName())
	}
}
