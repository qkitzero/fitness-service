package measurement

import "testing"

func TestNewSide(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		side    string
		want    string
	}{
		{"success none", true, "none", "none"},
		{"success left", true, "left", "left"},
		{"success right", true, "right", "right"},
		{"failure empty side", false, "", ""},
		{"failure unknown side", false, "both", ""},
		{"failure uppercase side", false, "LEFT", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			side, err := NewSide(tt.side)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if side.String() != tt.want {
				t.Errorf("String() = %v, want %v", side.String(), tt.want)
			}
		})
	}
}
