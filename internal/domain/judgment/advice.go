package judgment

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const adviceMaxLen = 2000

type Advice string

func (a Advice) String() string {
	return string(a)
}

func NewAdvice(s string) (*Advice, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(t) > adviceMaxLen {
		return nil, fmt.Errorf("invalid advice")
	}
	for _, r := range t {
		if r == '\n' || r == '\r' || r == '\t' {
			continue
		}
		if unicode.IsControl(r) {
			return nil, fmt.Errorf("invalid advice")
		}
	}
	a := Advice(t)
	return &a, nil
}
