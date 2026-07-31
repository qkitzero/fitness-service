package address

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const buildingMaxLen = 255

type Building string

func (b Building) String() string {
	return string(b)
}

func NewBuilding(s string) (*Building, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(t) > buildingMaxLen {
		return nil, fmt.Errorf("invalid building")
	}
	for _, r := range t {
		if unicode.IsControl(r) {
			return nil, fmt.Errorf("invalid building")
		}
	}
	b := Building(t)
	return &b, nil
}
