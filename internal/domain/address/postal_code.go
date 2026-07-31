package address

import (
	"fmt"
	"regexp"
	"strings"
)

var postalCodePattern = regexp.MustCompile(`^\d{3}-?\d{4}$`)

type PostalCode string

func (pc PostalCode) String() string {
	return string(pc)
}

func NewPostalCode(s string) (*PostalCode, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil, nil
	}
	if !postalCodePattern.MatchString(t) {
		return nil, fmt.Errorf("invalid postal code")
	}
	pc := PostalCode(strings.ReplaceAll(t, "-", ""))
	return &pc, nil
}
