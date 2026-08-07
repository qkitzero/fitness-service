package training

import "testing"

func TestNewSets(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		sets    int
		want    Sets
	}{
		{"success min sets", true, 1, Sets(1)},
		{"success typical sets", true, 3, Sets(3)},
		{"success max sets", true, 99, Sets(99)},
		{"failure zero sets", false, 0, Sets(0)},
		{"failure negative sets", false, -1, Sets(0)},
		{"failure sets above the max", false, 100, Sets(0)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sets, err := NewSets(tt.sets)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && sets != tt.want {
				t.Errorf("NewSets() = %v, want %v", sets, tt.want)
			}

			if tt.success && sets.Int() != tt.sets {
				t.Errorf("Int() = %v, want %v", sets.Int(), tt.sets)
			}
		})
	}
}
