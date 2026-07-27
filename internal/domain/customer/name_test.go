package customer

import (
	"strings"
	"testing"
)

func TestNewName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		success      bool
		customerName string
		want         string
	}{
		{"success new name", true, "test customer", "test customer"},
		{"success trim space", true, "  test customer  ", "test customer"},
		{"success max length name", true, strings.Repeat("あ", 255), strings.Repeat("あ", 255)},
		{"failure empty name", false, "", ""},
		{"failure whitespace name", false, "   ", ""},
		{"failure too long", false, strings.Repeat("あ", 256), ""},
		{"failure null character", false, "test\x00customer", ""},
		{"failure newline", false, "test\ncustomer", ""},
		{"failure tab", false, "test\tcustomer", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			name, err := NewName(tt.customerName)
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
