package measurementitem

import "testing"

func TestNewNormalization(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		success       bool
		normalization string
		want          Normalization
	}{
		{"success none", true, "none", NormalizationNone},
		{"success height ratio", true, "height_ratio", NormalizationHeightRatio},
		{"failure empty normalization", false, "", ""},
		{"failure invalid normalization", false, "weight_ratio", ""},
		{"failure uppercase normalization", false, "HEIGHT_RATIO", ""},
		{"failure untrimmed normalization", false, " height_ratio ", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			normalization, err := NewNormalization(tt.normalization)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if normalization != tt.want {
				t.Errorf("NewNormalization() = %v, want %v", normalization, tt.want)
			}

			if tt.success && normalization.String() != tt.normalization {
				t.Errorf("String() = %v, want %v", normalization.String(), tt.normalization)
			}
		})
	}
}
