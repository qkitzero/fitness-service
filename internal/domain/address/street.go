package address

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const streetMaxLen = 255

type Street string

func (st Street) String() string {
	return string(st)
}

func NewStreet(s string) (*Street, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(t) > streetMaxLen {
		return nil, fmt.Errorf("invalid street")
	}
	for _, r := range t {
		if unicode.IsControl(r) {
			return nil, fmt.Errorf("invalid street")
		}
	}
	st := Street(t)
	return &st, nil
}
