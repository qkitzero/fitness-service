package user

import (
	"context"
	"errors"

	groupv1 "github.com/qkitzero/user-service/gen/go/group/v1"
	"google.golang.org/grpc/metadata"

	"github.com/qkitzero/fitness-service/internal/application/user"
)

type userService struct {
	client groupv1.GroupServiceClient
}

func NewUserService(client groupv1.GroupServiceClient) user.UserService {
	return &userService{client: client}
}

func (s *userService) ListMyGroups(ctx context.Context) ([]string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is missing")
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	listMyGroupsRequest := &groupv1.ListMyGroupsRequest{}

	listMyGroupsResponse, err := s.client.ListMyGroups(ctx, listMyGroupsRequest)
	if err != nil {
		return nil, err
	}

	groupIDs := make([]string, 0, len(listMyGroupsResponse.GetGroups()))
	for _, membership := range listMyGroupsResponse.GetGroups() {
		groupIDs = append(groupIDs, membership.GetGroup().GetGroupId())
	}

	return groupIDs, nil
}
