package training

import "testing"

func TestNewUnit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		unit    string
		want    Unit
	}{
		{"success reps", true, "reps", UnitReps},
		{"success seconds", true, "seconds", UnitSeconds},
		{"success minutes", true, "minutes", UnitMinutes},
		{"failure empty unit", false, "", ""},
		{"failure invalid unit", false, "hours", ""},
		{"failure uppercase unit", false, "REPS", ""},
		{"failure untrimmed unit", false, " reps ", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			unit, err := NewUnit(tt.unit)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && unit != tt.want {
				t.Errorf("NewUnit() = %v, want %v", unit, tt.want)
			}

			if tt.success && unit.String() != tt.unit {
				t.Errorf("String() = %v, want %v", unit.String(), tt.unit)
			}
		})
	}
}
