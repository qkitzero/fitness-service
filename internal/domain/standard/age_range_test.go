package standard

import "testing"

func TestNewAgeRange(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		from    int
		to      int
	}{
		{"success new age range", true, 45, 49},
		{"success single age", true, 45, 45},
		{"success lower bound", true, 0, 4},
		{"success upper bound", true, 145, 150},
		{"success whole range", true, 0, 150},
		{"failure reversed age range", false, 49, 45},
		{"failure negative age from", false, -1, 49},
		{"failure negative age to", false, 0, -1},
		{"failure age from above upper bound", false, 151, 151},
		{"failure age to above upper bound", false, 145, 151},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ageRange, err := NewAgeRange(tt.from, tt.to)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && ageRange.From() != tt.from {
				t.Errorf("From() = %v, want %v", ageRange.From(), tt.from)
			}

			if tt.success && ageRange.To() != tt.to {
				t.Errorf("To() = %v, want %v", ageRange.To(), tt.to)
			}
		})
	}
}

func TestAgeRangeContains(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		from int
		to   int
		age  int
		want bool
	}{
		{"contains the lower bound", 45, 49, 45, true},
		{"contains the upper bound", 45, 49, 49, true},
		{"contains an age between the bounds", 45, 49, 47, true},
		{"does not contain an age below the lower bound", 45, 49, 44, false},
		{"does not contain an age above the upper bound", 45, 49, 50, false},
		{"contains the only age of a single age range", 45, 45, 45, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ageRange, err := NewAgeRange(tt.from, tt.to)
			if err != nil {
				t.Fatalf("failed to new age range: %v", err)
			}

			if got := ageRange.Contains(tt.age); got != tt.want {
				t.Errorf("Contains(%v) = %v, want %v", tt.age, got, tt.want)
			}
		})
	}
}
