package training

import "testing"

func TestNewLevel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		level   int
		want    Level
	}{
		{"success min level", true, 1, Level(1)},
		{"success middle level", true, 3, Level(3)},
		{"success max level", true, 5, Level(5)},
		{"failure zero level", false, 0, Level(0)},
		{"failure negative level", false, -1, Level(0)},
		{"failure level above the max", false, 6, Level(0)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			level, err := NewLevel(tt.level)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && level != tt.want {
				t.Errorf("NewLevel() = %v, want %v", level, tt.want)
			}

			if tt.success && level.Int() != tt.level {
				t.Errorf("Int() = %v, want %v", level.Int(), tt.level)
			}
		})
	}
}
