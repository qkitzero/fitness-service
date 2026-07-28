package organization

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	organizationv1 "github.com/qkitzero/fitness-service/gen/go/organization/v1"
	apporganization "github.com/qkitzero/fitness-service/internal/application/organization"
	appuser "github.com/qkitzero/fitness-service/internal/application/user"
	domainorganization "github.com/qkitzero/fitness-service/internal/domain/organization"
	domaintenant "github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type OrganizationHandler struct {
	organizationv1.UnimplementedOrganizationServiceServer
	organizationUsecase apporganization.OrganizationUsecase
}

func NewOrganizationHandler(
	organizationUsecase apporganization.OrganizationUsecase,
) *OrganizationHandler {
	return &OrganizationHandler{
		organizationUsecase: organizationUsecase,
	}
}

func toProtoOrganization(o domainorganization.Organization) *organizationv1.Organization {
	return &organizationv1.Organization{
		OrganizationId: o.ID().String(),
		TenantId:       o.TenantID().String(),
		Name:           o.Name().String(),
	}
}

func mapOrganizationError(err error, op string) error {
	if errors.Is(err, appuser.ErrNotGroupMember) {
		return status.Error(codes.PermissionDenied, err.Error())
	}
	if errors.Is(err, domainorganization.ErrOrganizationNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.Unauthenticated, codes.PermissionDenied:
			return err
		}
	}
	log.Printf("%s: internal error: %v", op, err)
	return status.Error(codes.Internal, "internal error")
}

func (h *OrganizationHandler) CreateOrganization(ctx context.Context, req *organizationv1.CreateOrganizationRequest) (*organizationv1.CreateOrganizationResponse, error) {
	tenantID, err := domaintenant.NewTenantID(req.GetTenantId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	name, err := domainorganization.NewName(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	organization, err := h.organizationUsecase.CreateOrganization(ctx, tenantID, name)
	if err != nil {
		return nil, mapOrganizationError(err, "CreateOrganization")
	}

	return &organizationv1.CreateOrganizationResponse{
		OrganizationId: organization.ID().String(),
	}, nil
}

func (h *OrganizationHandler) GetOrganization(ctx context.Context, req *organizationv1.GetOrganizationRequest) (*organizationv1.GetOrganizationResponse, error) {
	organizationID, err := domainorganization.NewOrganizationIDFromString(req.GetOrganizationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	organization, err := h.organizationUsecase.GetOrganization(ctx, organizationID)
	if err != nil {
		return nil, mapOrganizationError(err, "GetOrganization")
	}

	return &organizationv1.GetOrganizationResponse{
		Organization: toProtoOrganization(organization),
	}, nil
}

func (h *OrganizationHandler) ListOrganizations(ctx context.Context, req *organizationv1.ListOrganizationsRequest) (*organizationv1.ListOrganizationsResponse, error) {
	tenantID, err := domaintenant.NewTenantID(req.GetTenantId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	organizations, err := h.organizationUsecase.ListOrganizations(ctx, tenantID)
	if err != nil {
		return nil, mapOrganizationError(err, "ListOrganizations")
	}

	organizationMessages := make([]*organizationv1.Organization, 0, len(organizations))
	for _, o := range organizations {
		organizationMessages = append(organizationMessages, toProtoOrganization(o))
	}

	return &organizationv1.ListOrganizationsResponse{
		Organizations: organizationMessages,
	}, nil
}

func (h *OrganizationHandler) UpdateOrganization(ctx context.Context, req *organizationv1.UpdateOrganizationRequest) (*organizationv1.UpdateOrganizationResponse, error) {
	organizationID, err := domainorganization.NewOrganizationIDFromString(req.GetOrganizationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	name, err := domainorganization.NewName(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	organization, err := h.organizationUsecase.UpdateOrganization(ctx, organizationID, name)
	if err != nil {
		return nil, mapOrganizationError(err, "UpdateOrganization")
	}

	return &organizationv1.UpdateOrganizationResponse{
		Organization: toProtoOrganization(organization),
	}, nil
}

func (h *OrganizationHandler) DeleteOrganization(ctx context.Context, req *organizationv1.DeleteOrganizationRequest) (*organizationv1.DeleteOrganizationResponse, error) {
	organizationID, err := domainorganization.NewOrganizationIDFromString(req.GetOrganizationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.organizationUsecase.DeleteOrganization(ctx, organizationID); err != nil {
		return nil, mapOrganizationError(err, "DeleteOrganization")
	}

	return &organizationv1.DeleteOrganizationResponse{}, nil
}
