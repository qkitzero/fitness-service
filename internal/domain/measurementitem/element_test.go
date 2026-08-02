package measurementitem

import (
	"math"
	"testing"
)

func TestNewElement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		element string
		want    Element
	}{
		{"success muscle strength", true, "muscle_strength", ElementMuscleStrength},
		{"success muscle endurance", true, "muscle_endurance", ElementMuscleEndurance},
		{"success flexibility", true, "flexibility", ElementFlexibility},
		{"success agility", true, "agility", ElementAgility},
		{"success balance", true, "balance", ElementBalance},
		{"success mobility", true, "mobility", ElementMobility},
		{"failure empty element", false, "", ""},
		{"failure invalid element", false, "unknown", ""},
		{"failure uppercase element", false, "BALANCE", ""},
		{"failure untrimmed element", false, " balance ", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			element, err := NewElement(tt.element)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && element != tt.want {
				t.Errorf("NewElement() = %v, want %v", element, tt.want)
			}

			if tt.success && element.String() != tt.element {
				t.Errorf("String() = %v, want %v", element.String(), tt.element)
			}
		})
	}
}

func TestElementOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		element Element
		want    int
	}{
		{"muscle strength comes first", ElementMuscleStrength, 1},
		{"muscle endurance comes second", ElementMuscleEndurance, 2},
		{"flexibility comes third", ElementFlexibility, 3},
		{"agility comes fourth", ElementAgility, 4},
		{"balance comes fifth", ElementBalance, 5},
		{"mobility comes sixth", ElementMobility, 6},
		{"unknown element comes last", Element("unknown"), math.MaxInt},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.element.Order(); got != tt.want {
				t.Errorf("Order() = %v, want %v", got, tt.want)
			}
		})
	}
}
