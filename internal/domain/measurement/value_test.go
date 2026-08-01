package measurement

import (
	"math"
	"testing"
)

func TestNewValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		value   float64
		want    float64
	}{
		{"success new value", true, 32.4, 32.4},
		{"success zero value", true, 0, 0},
		{"success max value", true, 9999.99, 9999.99},
		{"success two decimal places", true, 32.45, 32.45},
		{"failure three decimal places", false, 12.345, 0},
		{"failure more precision than the column keeps", false, 70.125, 0},
		{"failure negative value", false, -0.1, 0},
		{"failure too large value", false, 10000, 0},
		{"failure not a number", false, math.NaN(), 0},
		{"failure infinity", false, math.Inf(1), 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, err := NewValue(tt.value)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if value.Float64() != tt.want {
				t.Errorf("Float64() = %v, want %v", value.Float64(), tt.want)
			}
		})
	}
}
