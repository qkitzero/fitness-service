package judgment

import (
	"testing"
)

func TestNewMotorAge(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		age  int
		want MotorAge
	}{
		{"success new motor age", 47, MotorAge(47)},
		{"success new motor age from the lowest age", 0, MotorAge(0)},
		{"success new motor age from an age below the registered age groups", 18, MotorAge(18)},
		{"success new motor age from an age above the registered age groups", 99, MotorAge(99)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			motorAge := NewMotorAge(tt.age)

			if motorAge != tt.want {
				t.Errorf("NewMotorAge() = %v, want %v", motorAge, tt.want)
			}

			if motorAge.Int() != int(tt.want) {
				t.Errorf("Int() = %v, want %v", motorAge.Int(), int(tt.want))
			}
		})
	}
}
