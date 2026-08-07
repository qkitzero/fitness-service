package training

import "testing"

func TestNewDecade(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		decade  int
		want    Decade
	}{
		{"success teens", true, 10, Decade10},
		{"success twenties", true, 20, Decade20},
		{"success thirties", true, 30, Decade30},
		{"success forties", true, 40, Decade40},
		{"success fifties", true, 50, Decade50},
		{"success sixties", true, 60, Decade60},
		{"success seventies", true, 70, Decade70},
		{"success eighties", true, 80, Decade80},
		{"success nineties", true, 90, Decade90},
		{"failure zero decade", false, 0, Decade(0)},
		{"failure negative decade", false, -10, Decade(0)},
		{"failure decade below the min", false, 5, Decade(0)},
		{"failure decade not a multiple of ten", false, 15, Decade(0)},
		{"failure decade above the max", false, 100, Decade(0)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decade, err := NewDecade(tt.decade)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && decade != tt.want {
				t.Errorf("NewDecade() = %v, want %v", decade, tt.want)
			}

			if tt.success && decade.Int() != tt.decade {
				t.Errorf("Int() = %v, want %v", decade.Int(), tt.decade)
			}
		})
	}
}
