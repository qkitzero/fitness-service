package judgment

import "testing"

func TestPrescriptionSourceString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		source PrescriptionSource
		want   string
	}{
		{"success element", PrescriptionSourceElement, "element"},
		{"success fixed", PrescriptionSourceFixed, "fixed"},
		{"success age decade", PrescriptionSourceAgeDecade, "age_decade"},
		{"success manual", PrescriptionSourceManual, "manual"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.source.String() != tt.want {
				t.Errorf("String() = %v, want %v", tt.source.String(), tt.want)
			}
		})
	}
}
