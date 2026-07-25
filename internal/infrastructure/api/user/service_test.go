package user

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/qkitzero/fitness-service/mocks/external/group/v1"
	groupv1 "github.com/qkitzero/user-service/gen/go/group/v1"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
)

func TestListMyGroups(t *testing.T) {
	t.Parallel()
	accessToken := "accessToken"
	groupID := "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"
	tests := []struct {
		name            string
		success         bool
		ctx             context.Context
		listMyGroupsErr error
	}{
		{
			name:            "success list my groups",
			success:         true,
			ctx:             metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+accessToken)),
			listMyGroupsErr: nil,
		},
		{
			name:            "failure missing metadata",
			success:         false,
			ctx:             context.Background(),
			listMyGroupsErr: nil,
		},
		{
			name:            "failure list my groups error",
			success:         false,
			ctx:             metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+accessToken)),
			listMyGroupsErr: errors.New("list my groups error"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockGroupServiceClient(ctrl)
			mockListMyGroupsResponse := &groupv1.ListMyGroupsResponse{
				Groups: []*groupv1.GroupMembership{
					{
						Group: &groupv1.Group{GroupId: groupID, Name: "test group"},
						Role:  "admin",
					},
				},
			}
			mockClient.EXPECT().ListMyGroups(gomock.Any(), gomock.Any()).Return(mockListMyGroupsResponse, tt.listMyGroupsErr).AnyTimes()

			userService := NewUserService(mockClient)

			groupIDs, err := userService.ListMyGroups(tt.ctx)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && (len(groupIDs) != 1 || groupIDs[0] != groupID) {
				t.Errorf("ListMyGroups() = %v, want [%v]", groupIDs, groupID)
			}
		})
	}
}
