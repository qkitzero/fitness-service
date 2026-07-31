package tenant

import (
	"strings"
	"testing"
)

func TestNewHomepageURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		input   string
		want    string
	}{
		{"success https", true, "https://example.com", "https://example.com"},
		{"success http", true, "http://example.com", "http://example.com"},
		{"success with path", true, "https://example.com/clinic", "https://example.com/clinic"},
		{"success normalizes the scheme and host case", true, "HTTPS://Example.COM/Clinic", "https://example.com/Clinic"},
		{"success blank", true, "  ", ""},
		{"failure too long", false, "https://example.com/" + strings.Repeat("a", 256), ""},
		{"failure invalid scheme", false, "ftp://example.com", ""},
		{"failure missing scheme", false, "example.com", ""},
		{"failure missing host", false, "https://", ""},
		{"failure unparsable", false, "https://[::1", ""},
		{"failure userinfo", false, "https://admin:s3cret@example.com", ""},
		{"failure user only", false, "https://admin@example.com", ""},
		{"failure null character", false, "https://exa\x00mple.com", ""},
		{"failure space", false, "https://example.com/x y", ""},
		{"failure no-break space", false, "https://example.com/x\u00a0y", ""},
		{"failure line separator", false, "https://example.com/x\u2028y", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewHomepageURL(tt.input)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && tt.want == "" && got != nil {
				t.Errorf("expected nil, but got %v", *got)
			}
			if tt.success && tt.want != "" && (got == nil || got.String() != tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
