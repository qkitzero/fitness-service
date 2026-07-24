package customer

import (
	"fmt"
	"regexp"
	"strings"
)

var nameKanaPattern = regexp.MustCompile(`^[\p{Katakana}ー・ 　]+$`)

type NameKana string

func (nk NameKana) String() string {
	return string(nk)
}

func NewNameKana(s string) (NameKana, error) {
	s = strings.TrimSpace(s)
	if !nameKanaPattern.MatchString(s) {
		return NameKana(""), fmt.Errorf("invalid name kana")
	}
	return NameKana(s), nil
}
