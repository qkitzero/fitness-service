package measurement

import "testing"

func TestNewTrialIndex(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		success    bool
		trialIndex int
		want       int
	}{
		{"success new trial index", true, 1, 1},
		{"success second trial index", true, 2, 2},
		{"failure zero trial index", false, 0, 0},
		{"failure negative trial index", false, -1, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trialIndex, err := NewTrialIndex(tt.trialIndex)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if trialIndex.Int() != tt.want {
				t.Errorf("Int() = %v, want %v", trialIndex.Int(), tt.want)
			}
		})
	}
}
