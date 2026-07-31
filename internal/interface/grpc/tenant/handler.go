package tenant

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tenantv1 "github.com/qkitzero/fitness-service/gen/go/tenant/v1"
	apptenant "github.com/qkitzero/fitness-service/internal/application/tenant"
	domainaddress "github.com/qkitzero/fitness-service/internal/domain/address"
	domaincontact "github.com/qkitzero/fitness-service/internal/domain/contact"
	domaintenant "github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type ProfileHandler struct {
	tenantv1.UnimplementedProfileServiceServer
	profileUsecase apptenant.ProfileUsecase
}

func NewProfileHandler(
	profileUsecase apptenant.ProfileUsecase,
) *ProfileHandler {
	return &ProfileHandler{
		profileUsecase: profileUsecase,
	}
}

func toProtoProfile(p domaintenant.Profile) *tenantv1.Profile {
	msg := &tenantv1.Profile{
		TenantId: p.TenantID().String(),
	}
	addr := p.Address()
	if v := addr.PostalCode(); v != nil {
		s := v.String()
		msg.PostalCode = &s
	}
	if v := addr.Prefecture(); v != nil {
		s := v.String()
		msg.Prefecture = &s
	}
	if v := addr.City(); v != nil {
		s := v.String()
		msg.City = &s
	}
	if v := addr.Street(); v != nil {
		s := v.String()
		msg.Street = &s
	}
	if v := addr.Building(); v != nil {
		s := v.String()
		msg.Building = &s
	}
	if v := p.Phone(); v != nil {
		s := v.String()
		msg.Phone = &s
	}
	if v := p.Email(); v != nil {
		s := v.String()
		msg.Email = &s
	}
	if v := p.HomepageURL(); v != nil {
		s := v.String()
		msg.HomepageUrl = &s
	}
	if v := p.Note(); v != nil {
		s := v.String()
		msg.Note = &s
	}
	return msg
}

func parseProfilePatch(req *tenantv1.UpsertProfileRequest) (apptenant.ProfilePatch, error) {
	var patch apptenant.ProfilePatch
	var err error
	if req.PostalCode != nil {
		patch.HasPostalCode = true
		if patch.PostalCode, err = domainaddress.NewPostalCode(*req.PostalCode); err != nil {
			return patch, err
		}
	}
	if req.Prefecture != nil {
		patch.HasPrefecture = true
		if patch.Prefecture, err = domainaddress.NewPrefecture(*req.Prefecture); err != nil {
			return patch, err
		}
	}
	if req.City != nil {
		patch.HasCity = true
		if patch.City, err = domainaddress.NewCity(*req.City); err != nil {
			return patch, err
		}
	}
	if req.Street != nil {
		patch.HasStreet = true
		if patch.Street, err = domainaddress.NewStreet(*req.Street); err != nil {
			return patch, err
		}
	}
	if req.Building != nil {
		patch.HasBuilding = true
		if patch.Building, err = domainaddress.NewBuilding(*req.Building); err != nil {
			return patch, err
		}
	}
	if req.Phone != nil {
		patch.HasPhone = true
		if patch.Phone, err = domaincontact.NewPhone(*req.Phone); err != nil {
			return patch, err
		}
	}
	if req.Email != nil {
		patch.HasEmail = true
		if patch.Email, err = domaincontact.NewEmail(*req.Email); err != nil {
			return patch, err
		}
	}
	if req.HomepageUrl != nil {
		patch.HasHomepageURL = true
		if patch.HomepageURL, err = domaintenant.NewHomepageURL(*req.HomepageUrl); err != nil {
			return patch, err
		}
	}
	if req.Note != nil {
		patch.HasNote = true
		if patch.Note, err = domaintenant.NewNote(*req.Note); err != nil {
			return patch, err
		}
	}
	return patch, nil
}

func mapProfileError(err error, op string) error {
	if errors.Is(err, domaintenant.ErrNotMember) {
		return status.Error(codes.PermissionDenied, err.Error())
	}
	if errors.Is(err, domaintenant.ErrProfileNotFound) {
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

func (h *ProfileHandler) GetProfile(ctx context.Context, req *tenantv1.GetProfileRequest) (*tenantv1.GetProfileResponse, error) {
	tenantID, err := domaintenant.NewTenantID(req.GetTenantId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	profile, err := h.profileUsecase.GetProfile(ctx, tenantID)
	if err != nil {
		return nil, mapProfileError(err, "GetProfile")
	}

	return &tenantv1.GetProfileResponse{
		Profile: toProtoProfile(profile),
	}, nil
}

func (h *ProfileHandler) UpsertProfile(ctx context.Context, req *tenantv1.UpsertProfileRequest) (*tenantv1.UpsertProfileResponse, error) {
	tenantID, err := domaintenant.NewTenantID(req.GetTenantId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	patch, err := parseProfilePatch(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	profile, err := h.profileUsecase.UpsertProfile(ctx, tenantID, patch)
	if err != nil {
		return nil, mapProfileError(err, "UpsertProfile")
	}

	return &tenantv1.UpsertProfileResponse{
		Profile: toProtoProfile(profile),
	}, nil
}
