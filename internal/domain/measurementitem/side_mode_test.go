package measurementitem

import "testing"

func TestNewSideMode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		success  bool
		sideMode string
		want     SideMode
	}{
		{"success none", true, "none", SideModeNone},
		{"success bilateral", true, "bilateral", SideModeBilateral},
		{"success optional bilateral", true, "optional_bilateral", SideModeOptionalBilateral},
		{"failure empty side mode", false, "", ""},
		{"failure invalid side mode", false, "unilateral", ""},
		{"failure uppercase side mode", false, "NONE", ""},
		{"failure untrimmed side mode", false, " none ", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sideMode, err := NewSideMode(tt.sideMode)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if sideMode != tt.want {
				t.Errorf("NewSideMode() = %v, want %v", sideMode, tt.want)
			}

			if tt.success && sideMode.String() != tt.sideMode {
				t.Errorf("String() = %v, want %v", sideMode.String(), tt.sideMode)
			}
		})
	}
}
