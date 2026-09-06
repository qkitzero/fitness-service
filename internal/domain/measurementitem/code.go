package measurementitem

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const codeMaxLen = 64

type Code string

const CodeHeight Code = "height"

func (c Code) String() string {
	return string(c)
}

func NewCode(s string) (Code, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Code(""), fmt.Errorf("invalid code")
	}
	if utf8.RuneCountInString(s) > codeMaxLen {
		return Code(""), fmt.Errorf("invalid code")
	}
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		case r == '_':
		default:
			return Code(""), fmt.Errorf("invalid code")
		}
	}
	return Code(s), nil
}
