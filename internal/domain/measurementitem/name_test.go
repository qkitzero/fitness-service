package measurementitem

import (
	"strings"
	"testing"
)

func TestNewName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		success         bool
		measurementName string
		want            Name
	}{
		{"success new name", true, "握力", "握力"},
		{"success trim space", true, "  長座体前屈  ", "長座体前屈"},
		{"success max length name", true, strings.Repeat("あ", 255), Name(strings.Repeat("あ", 255))},
		{"failure empty name", false, "", ""},
		{"failure whitespace name", false, "   ", ""},
		{"failure too long name", false, strings.Repeat("あ", 256), ""},
		{"failure null character name", false, "握\x00力", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			measurementName, err := NewName(tt.measurementName)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && measurementName != tt.want {
				t.Errorf("NewName() = %v, want %v", measurementName, tt.want)
			}

			if tt.success && measurementName.String() != tt.want.String() {
				t.Errorf("String() = %v, want %v", measurementName.String(), tt.want.String())
			}
		})
	}
}
