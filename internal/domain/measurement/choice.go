package measurement

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const choiceMaxLen = 32

type Choice string

func (c Choice) String() string {
	return string(c)
}

func NewChoice(s string) (Choice, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Choice(""), fmt.Errorf("invalid choice")
	}
	if utf8.RuneCountInString(s) > choiceMaxLen {
		return Choice(""), fmt.Errorf("invalid choice")
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return Choice(""), fmt.Errorf("invalid choice")
		}
	}
	return Choice(s), nil
}
