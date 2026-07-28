package customer

import (
	"strings"
	"testing"
)

func TestNewEmergencyContactName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		input   string
		want    string
	}{
		{"success", true, "緊急 太郎", "緊急 太郎"},
		{"success blank", true, "  ", ""},
		{"success max length", true, strings.Repeat("あ", 255), strings.Repeat("あ", 255)},
		{"failure too long", false, strings.Repeat("あ", 256), ""},
		{"failure null character", false, "緊急\x00太郎", ""},
		{"failure newline", false, "緊急\n太郎", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewEmergencyContactName(tt.input)
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
