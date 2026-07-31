package contact

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	phonePattern    = regexp.MustCompile(`^0\d{9,10}$`)
	phoneSeparators = strings.NewReplacer("-", "", "－", "", " ", "", "　", "")
)

type Phone string

func (p Phone) String() string {
	return string(p)
}

func NewPhone(s string) (*Phone, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil, nil
	}
	digits := phoneSeparators.Replace(t)
	if !phonePattern.MatchString(digits) {
		return nil, fmt.Errorf("invalid phone")
	}
	p := Phone(digits)
	return &p, nil
}
