package measurement

import (
	"errors"
	"testing"
)

func TestNewAgeAtMeasurement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		age     int
		want    int
	}{
		{"success new age at measurement", true, 65, 65},
		{"success min age at measurement", true, 0, 0},
		{"success max age at measurement", true, 150, 150},
		{"failure negative age at measurement", false, -1, 0},
		{"failure too large age at measurement", false, 151, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			age, err := NewAgeAtMeasurement(tt.age)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && !errors.Is(err, ErrInvalidAgeAtMeasurement) {
				t.Errorf("err = %v, want %v", err, ErrInvalidAgeAtMeasurement)
			}
			if age.Int() != tt.want {
				t.Errorf("Int() = %v, want %v", age.Int(), tt.want)
			}
		})
	}
}
