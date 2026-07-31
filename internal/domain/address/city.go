package address

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const cityMaxLen = 255

type City string

func (c City) String() string {
	return string(c)
}

func NewCity(s string) (*City, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(t) > cityMaxLen {
		return nil, fmt.Errorf("invalid city")
	}
	for _, r := range t {
		if unicode.IsControl(r) {
			return nil, fmt.Errorf("invalid city")
		}
	}
	c := City(t)
	return &c, nil
}
