package customer

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"
)

const emailMaxLen = 255

type Email string

func (e Email) String() string {
	return string(e)
}

func NewEmail(s string) (*Email, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(t) > emailMaxLen {
		return nil, fmt.Errorf("invalid email")
	}
	addr, err := mail.ParseAddress(t)
	if err != nil || addr.Address != t {
		return nil, fmt.Errorf("invalid email")
	}
	e := Email(addr.Address)
	return &e, nil
}
