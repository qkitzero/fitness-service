package customer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const nameMaxLen = 255

type Name string

func (n Name) String() string {
	return string(n)
}

func NewName(s string) (Name, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Name(""), fmt.Errorf("invalid name")
	}
	if utf8.RuneCountInString(s) > nameMaxLen {
		return Name(""), fmt.Errorf("invalid name")
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return Name(""), fmt.Errorf("invalid name")
		}
	}
	return Name(s), nil
}
