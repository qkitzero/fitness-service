package measurementitem

import "testing"

func TestNewScoreDirection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		scoreDirection string
		want           ScoreDirection
	}{
		{"success higher is better", true, "higher_is_better", ScoreDirectionHigherIsBetter},
		{"success lower is better", true, "lower_is_better", ScoreDirectionLowerIsBetter},
		{"failure empty score direction", false, "", ""},
		{"failure invalid score direction", false, "unknown", ""},
		{"failure uppercase score direction", false, "HIGHER_IS_BETTER", ""},
		{"failure untrimmed score direction", false, " higher_is_better ", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			scoreDirection, err := NewScoreDirection(tt.scoreDirection)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && scoreDirection != tt.want {
				t.Errorf("NewScoreDirection() = %v, want %v", scoreDirection, tt.want)
			}

			if tt.success && scoreDirection.String() != tt.scoreDirection {
				t.Errorf("String() = %v, want %v", scoreDirection.String(), tt.scoreDirection)
			}
		})
	}
}
