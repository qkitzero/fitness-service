package tenant

import "errors"

var (
	ErrNotMember       = errors.New("not a tenant member")
	ErrProfileNotFound = errors.New("tenant profile not found")
)
