package measurementitem

import "testing"

func TestNewTrialCount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		success    bool
		trialCount int
		want       TrialCount
	}{
		{"success single trial", true, 1, TrialCount(1)},
		{"success multiple trials", true, 5, TrialCount(5)},
		{"failure zero trial count", false, 0, TrialCount(0)},
		{"failure negative trial count", false, -1, TrialCount(0)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trialCount, err := NewTrialCount(tt.trialCount)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && trialCount != tt.want {
				t.Errorf("NewTrialCount() = %v, want %v", trialCount, tt.want)
			}

			if tt.success && trialCount.Int() != tt.trialCount {
				t.Errorf("Int() = %v, want %v", trialCount.Int(), tt.trialCount)
			}
		})
	}
}
