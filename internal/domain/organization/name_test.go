package organization

import (
	"strings"
	"testing"
)

func TestNewName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		success          bool
		organizationName string
		want             string
	}{
		{"success new name", true, "テスト株式会社", "テスト株式会社"},
		{"success trim space", true, "  テスト株式会社  ", "テスト株式会社"},
		{"success max length name", true, strings.Repeat("あ", 255), strings.Repeat("あ", 255)},
		{"failure empty name", false, "", ""},
		{"failure whitespace name", false, "   ", ""},
		{"failure too long", false, strings.Repeat("あ", 256), ""},
		{"failure null character", false, "テスト\x00株式会社", ""},
		{"failure newline", false, "テスト\n株式会社", ""},
		{"failure tab", false, "テスト\t株式会社", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			name, err := NewName(tt.organizationName)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && name.String() != tt.want {
				t.Errorf("String() = %v, want %v", name.String(), tt.want)
			}
		})
	}
}
