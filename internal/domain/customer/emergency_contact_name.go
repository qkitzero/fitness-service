package customer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const emergencyContactNameMaxLen = 255

type EmergencyContactName string

func (n EmergencyContactName) String() string {
	return string(n)
}

func NewEmergencyContactName(s string) (*EmergencyContactName, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(t) > emergencyContactNameMaxLen {
		return nil, fmt.Errorf("invalid emergency contact name")
	}
	for _, r := range t {
		if unicode.IsControl(r) {
			return nil, fmt.Errorf("invalid emergency contact name")
		}
	}
	n := EmergencyContactName(t)
	return &n, nil
}
