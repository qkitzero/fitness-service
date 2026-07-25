package customer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const emergencyContactRelationshipMaxLen = 255

type EmergencyContactRelationship string

func (r EmergencyContactRelationship) String() string {
	return string(r)
}

func NewEmergencyContactRelationship(s string) (*EmergencyContactRelationship, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(t) > emergencyContactRelationshipMaxLen {
		return nil, fmt.Errorf("invalid emergency contact relationship")
	}
	r := EmergencyContactRelationship(t)
	return &r, nil
}
