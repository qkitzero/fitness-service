package measurement

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const noteMaxLen = 255

type Note string

func (n Note) String() string {
	return string(n)
}

func NewNote(s string) (*Note, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(t) > noteMaxLen {
		return nil, fmt.Errorf("invalid note")
	}
	for _, r := range t {
		if r == '\n' || r == '\r' || r == '\t' {
			continue
		}
		if unicode.IsControl(r) {
			return nil, fmt.Errorf("invalid note")
		}
	}
	n := Note(t)
	return &n, nil
}
