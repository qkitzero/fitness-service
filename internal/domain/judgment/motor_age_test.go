package judgment

import (
	"testing"

	"github.com/qkitzero/fitness-service/internal/domain/standard"
)

func TestNewMotorAge(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		from int
		to   int
		want MotorAge
	}{
		{"success new motor age from an odd width age range", 45, 49, MotorAge(47)},
		{"success new motor age from an even width age range", 40, 49, MotorAge(45)},
		{"success new motor age from a single age range", 45, 45, MotorAge(45)},
		{"success new motor age from the lowest age range", 0, 4, MotorAge(2)},
		{"success new motor age from an open ended age range", 65, 150, MotorAge(65)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ageRange, err := standard.NewAgeRange(tt.from, tt.to)
			if err != nil {
				t.Fatalf("failed to new age range: %v", err)
			}

			motorAge := NewMotorAge(ageRange)

			if motorAge != tt.want {
				t.Errorf("NewMotorAge() = %v, want %v", motorAge, tt.want)
			}

			if motorAge.Int() != int(tt.want) {
				t.Errorf("Int() = %v, want %v", motorAge.Int(), int(tt.want))
			}
		})
	}
}
