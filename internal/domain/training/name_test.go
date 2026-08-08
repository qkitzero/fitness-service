package training

import (
	"strings"
	"testing"
)

func TestNewName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		success      bool
		trainingName string
		want         Name
	}{
		{"success new name", true, "壁押し", "壁押し"},
		{"success trim space", true, "  スクワット  ", "スクワット"},
		{"success max length name", true, strings.Repeat("a", 255), Name(strings.Repeat("a", 255))},
		{"failure empty name", false, "", ""},
		{"failure whitespace name", false, "   ", ""},
		{"failure too long name", false, strings.Repeat("a", 256), ""},
		{"failure newline name", false, "壁\n押し", ""},
		{"failure tab name", false, "壁\t押し", ""},
		{"failure null character name", false, "壁\x00押し", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trainingName, err := NewName(tt.trainingName)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && trainingName != tt.want {
				t.Errorf("NewName() = %v, want %v", trainingName, tt.want)
			}

			if tt.success && trainingName.String() != tt.want.String() {
				t.Errorf("String() = %v, want %v", trainingName.String(), tt.want.String())
			}
		})
	}
}
