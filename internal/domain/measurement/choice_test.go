package measurement

import (
	"strings"
	"testing"
)

func TestNewChoice(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		choice  string
		want    string
	}{
		{"success new choice", true, "片足20cm", "片足20cm"},
		{"success trim whitespace", true, "  両足50cm  ", "両足50cm"},
		{"success max length choice", true, strings.Repeat("あ", 32), strings.Repeat("あ", 32)},
		{"failure empty choice", false, "", ""},
		{"failure whitespace choice", false, "   ", ""},
		{"failure too long choice", false, strings.Repeat("あ", 33), ""},
		{"failure control character", false, "片足20cm\x00", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			choice, err := NewChoice(tt.choice)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if choice.String() != tt.want {
				t.Errorf("String() = %v, want %v", choice.String(), tt.want)
			}
		})
	}
}
