package standard

import (
	"math"
	"testing"
)

func TestNewMean(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		mean    float64
		want    Mean
	}{
		{"success two decimal mean", true, 46.53, Mean(46.53)},
		{"success one decimal mean", true, 46.5, Mean(46.5)},
		{"success integer mean", true, 46, Mean(46)},
		{"success zero mean", true, 0, Mean(0)},
		{"success max mean", true, 9999.99, Mean(9999.99)},
		{"failure negative mean", false, -0.01, Mean(0)},
		{"failure too large mean", false, 10000, Mean(0)},
		{"failure three decimal mean", false, 46.535, Mean(0)},
		{"failure NaN mean", false, math.NaN(), Mean(0)},
		{"failure infinite mean", false, math.Inf(1), Mean(0)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mean, err := NewMean(tt.mean)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && mean != tt.want {
				t.Errorf("NewMean() = %v, want %v", mean, tt.want)
			}

			if tt.success && mean.Float64() != tt.mean {
				t.Errorf("Float64() = %v, want %v", mean.Float64(), tt.mean)
			}
		})
	}
}
