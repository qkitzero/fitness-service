package customer

import (
	"strings"
	"testing"
)

func TestNewEmail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		input   string
		want    string
	}{
		{"success", true, "test@example.com", "test@example.com"},
		{"success blank", true, "  ", ""},
		{"failure too long", false, strings.Repeat("a", 256) + "@example.com", ""},
		{"failure invalid", false, "invalid", ""},
		{"failure display name", false, "Tanaka Taro <taro@example.com>", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewEmail(tt.input)
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
