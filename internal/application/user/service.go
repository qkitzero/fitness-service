package user

import (
	"context"
)

type UserService interface {
	ListMyGroups(ctx context.Context) ([]string, error)
}
