package standard

import (
	"math"
	"testing"
)

func TestNewZScore(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		zScore  float64
		want    ZScore
	}{
		{"success positive z score", true, 1.5, ZScore(1.5)},
		{"success negative z score", true, -1.5, ZScore(-1.5)},
		{"success zero z score", true, 0, ZScore(0)},
		{"success max z score", true, 99.99, ZScore(99.99)},
		{"success min z score", true, -99.99, ZScore(-99.99)},
		{"failure too large z score", false, 100, ZScore(0)},
		{"failure too small z score", false, -100, ZScore(0)},
		{"failure three decimal z score", false, 1.505, ZScore(0)},
		{"failure NaN z score", false, math.NaN(), ZScore(0)},
		{"failure infinite z score", false, math.Inf(-1), ZScore(0)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			zScore, err := NewZScore(tt.zScore)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && zScore != tt.want {
				t.Errorf("NewZScore() = %v, want %v", zScore, tt.want)
			}

			if tt.success && zScore.Float64() != tt.zScore {
				t.Errorf("Float64() = %v, want %v", zScore.Float64(), tt.zScore)
			}
		})
	}
}
