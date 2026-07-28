package tenant

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const tenantIDMaxLen = 36

type TenantID string

func (t TenantID) String() string {
	return string(t)
}

func NewTenantID(s string) (TenantID, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return TenantID(""), fmt.Errorf("invalid tenant id")
	}
	if utf8.RuneCountInString(s) > tenantIDMaxLen {
		return TenantID(""), fmt.Errorf("invalid tenant id")
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return TenantID(""), fmt.Errorf("invalid tenant id")
		}
	}
	return TenantID(s), nil
}
