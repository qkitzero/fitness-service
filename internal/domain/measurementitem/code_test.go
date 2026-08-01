package measurementitem

import (
	"strings"
	"testing"
)

func TestNewCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		code    string
		want    Code
	}{
		{"success new code", true, "grip_strength", "grip_strength"},
		{"success digits", true, "cs30", "cs30"},
		{"success trim space", true, "  walk_5m  ", "walk_5m"},
		{"success max length code", true, strings.Repeat("a", 64), Code(strings.Repeat("a", 64))},
		{"failure empty code", false, "", ""},
		{"failure whitespace code", false, "   ", ""},
		{"failure too long code", false, strings.Repeat("a", 65), ""},
		{"failure uppercase code", false, "GripStrength", ""},
		{"failure hyphen code", false, "grip-strength", ""},
		{"failure inner space code", false, "grip strength", ""},
		{"failure multibyte code", false, "握力", ""},
		{"failure null character code", false, "grip\x00strength", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, err := NewCode(tt.code)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && code != tt.want {
				t.Errorf("NewCode() = %v, want %v", code, tt.want)
			}

			if tt.success && code.String() != tt.want.String() {
				t.Errorf("String() = %v, want %v", code.String(), tt.want.String())
			}
		})
	}
}
