package measurement

import "fmt"

type Side string

const (
	SideNone  Side = "none"
	SideLeft  Side = "left"
	SideRight Side = "right"
)

func (s Side) String() string {
	return string(s)
}

func NewSide(s string) (Side, error) {
	switch Side(s) {
	case SideNone, SideLeft, SideRight:
		return Side(s), nil
	default:
		return Side(""), fmt.Errorf("invalid side")
	}
}
