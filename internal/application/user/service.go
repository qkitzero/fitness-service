package user

import (
	"context"
	"errors"
)

var ErrNotGroupMember = errors.New("not a group member")

type UserService interface {
	ListMyGroups(ctx context.Context) ([]string, error)
}
