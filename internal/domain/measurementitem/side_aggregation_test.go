package measurementitem

import "testing"

func TestNewSideAggregation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		success         bool
		sideAggregation string
		want            SideAggregation
	}{
		{"success mean", true, "mean", SideAggregationMean},
		{"success best", true, "best", SideAggregationBest},
		{"failure empty side aggregation", false, "", ""},
		{"failure invalid side aggregation", false, "median", ""},
		{"failure uppercase side aggregation", false, "MEAN", ""},
		{"failure untrimmed side aggregation", false, " mean ", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sideAggregation, err := NewSideAggregation(tt.sideAggregation)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if sideAggregation != tt.want {
				t.Errorf("NewSideAggregation() = %v, want %v", sideAggregation, tt.want)
			}

			if tt.success && sideAggregation.String() != tt.sideAggregation {
				t.Errorf("String() = %v, want %v", sideAggregation.String(), tt.sideAggregation)
			}
		})
	}
}
