package tenant

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"
)

const homepageURLMaxLen = 255

type HomepageURL string

func (h HomepageURL) String() string {
	return string(h)
}

func NewHomepageURL(s string) (*HomepageURL, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(t) > homepageURLMaxLen {
		return nil, fmt.Errorf("invalid homepage url")
	}
	for _, r := range t {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return nil, fmt.Errorf("invalid homepage url")
		}
	}
	parsed, err := url.Parse(t)
	if err != nil {
		return nil, fmt.Errorf("invalid homepage url")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("invalid homepage url")
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("invalid homepage url")
	}
	if parsed.User != nil {
		return nil, fmt.Errorf("invalid homepage url")
	}
	parsed.Host = strings.ToLower(parsed.Host)
	h := HomepageURL(parsed.String())
	return &h, nil
}
