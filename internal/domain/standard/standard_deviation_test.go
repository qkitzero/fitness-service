package standard

import (
	"math"
	"testing"
)

func TestNewStandardDeviation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		success           bool
		standardDeviation float64
		want              StandardDeviation
	}{
		{"success two decimal standard deviation", true, 6.85, StandardDeviation(6.85)},
		{"success one decimal standard deviation", true, 6.8, StandardDeviation(6.8)},
		{"success integer standard deviation", true, 7, StandardDeviation(7)},
		{"success smallest standard deviation", true, 0.01, StandardDeviation(0.01)},
		{"success max standard deviation", true, 9999.99, StandardDeviation(9999.99)},
		{"failure zero standard deviation", false, 0, StandardDeviation(0)},
		{"failure negative standard deviation", false, -6.8, StandardDeviation(0)},
		{"failure too large standard deviation", false, 10000, StandardDeviation(0)},
		{"failure three decimal standard deviation", false, 6.855, StandardDeviation(0)},
		{"failure NaN standard deviation", false, math.NaN(), StandardDeviation(0)},
		{"failure infinite standard deviation", false, math.Inf(1), StandardDeviation(0)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			standardDeviation, err := NewStandardDeviation(tt.standardDeviation)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && standardDeviation != tt.want {
				t.Errorf("NewStandardDeviation() = %v, want %v", standardDeviation, tt.want)
			}

			if tt.success && standardDeviation.Float64() != tt.standardDeviation {
				t.Errorf("Float64() = %v, want %v", standardDeviation.Float64(), tt.standardDeviation)
			}
		})
	}
}
