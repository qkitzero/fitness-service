package judgment

import "testing"

func TestNewMotorAge(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		success  bool
		motorAge int
		want     MotorAge
	}{
		{"success new motor age", true, 47, MotorAge(47)},
		{"success lower bound", true, 0, MotorAge(0)},
		{"success upper bound", true, 150, MotorAge(150)},
		{"failure negative motor age", false, -1, MotorAge(0)},
		{"failure motor age above upper bound", false, 151, MotorAge(0)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			motorAge, err := NewMotorAge(tt.motorAge)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if motorAge != tt.want {
				t.Errorf("NewMotorAge() = %v, want %v", motorAge, tt.want)
			}

			if tt.success && motorAge.Int() != tt.motorAge {
				t.Errorf("Int() = %v, want %v", motorAge.Int(), tt.motorAge)
			}
		})
	}
}
