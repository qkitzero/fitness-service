package measurementitem

import "fmt"

type SideMode string

const (
	SideModeNone              SideMode = "none"
	SideModeBilateral         SideMode = "bilateral"
	SideModeOptionalBilateral SideMode = "optional_bilateral"
)

func (s SideMode) String() string {
	return string(s)
}

func NewSideMode(s string) (SideMode, error) {
	switch SideMode(s) {
	case SideModeNone, SideModeBilateral, SideModeOptionalBilateral:
		return SideMode(s), nil
	default:
		return SideMode(""), fmt.Errorf("invalid side mode")
	}
}
