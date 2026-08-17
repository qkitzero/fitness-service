package measurementitem

import "testing"

func TestNewUnit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		unit    string
		want    Unit
	}{
		{"success kg", true, "kg", UnitKg},
		{"success cm", true, "cm", UnitCm},
		{"success sec", true, "sec", UnitSec},
		{"success count", true, "count", UnitCount},
		{"success mmHg", true, "mmHg", UnitMmHg},
		{"success percent", true, "percent", UnitPercent},
		{"success bpm", true, "bpm", UnitBpm},
		{"success level", true, "level", UnitLevel},
		{"failure empty unit", false, "", ""},
		{"failure invalid unit", false, "m", ""},
		{"failure lowercase mmhg", false, "mmhg", ""},
		{"failure untrimmed unit", false, " kg ", ""},
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
