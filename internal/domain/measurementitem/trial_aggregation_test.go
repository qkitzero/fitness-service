package measurementitem

import "testing"

func TestNewTrialAggregation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		success          bool
		trialAggregation string
		want             TrialAggregation
	}{
		{"success mean", true, "mean", TrialAggregationMean},
		{"success best", true, "best", TrialAggregationBest},
		{"failure empty trial aggregation", false, "", ""},
		{"failure invalid trial aggregation", false, "worst", ""},
		{"failure uppercase trial aggregation", false, "MEAN", ""},
		{"failure untrimmed trial aggregation", false, " mean ", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trialAggregation, err := NewTrialAggregation(tt.trialAggregation)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if trialAggregation != tt.want {
				t.Errorf("NewTrialAggregation() = %v, want %v", trialAggregation, tt.want)
			}

			if tt.success && trialAggregation.String() != tt.trialAggregation {
				t.Errorf("String() = %v, want %v", trialAggregation.String(), tt.trialAggregation)
			}
		})
	}
}
