package staff

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const staffIDMaxLen = 255

type StaffID string

func (s StaffID) String() string {
	return string(s)
}

func NewStaffID(s string) (StaffID, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return StaffID(""), fmt.Errorf("invalid staff id")
	}
	if utf8.RuneCountInString(s) > staffIDMaxLen {
		return StaffID(""), fmt.Errorf("invalid staff id")
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return StaffID(""), fmt.Errorf("invalid staff id")
		}
	}
	return StaffID(s), nil
}
