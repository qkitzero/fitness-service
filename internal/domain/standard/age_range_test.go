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

func TestAgeRangeMedian(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		from int
		to   int
		want int
	}{
		{"median of an odd width range", 45, 49, 47},
		{"median of an even width range", 40, 49, 45},
		{"median of a single age range", 45, 45, 45},
		{"median of the lowest range", 0, 4, 2},
		{"lower bound of an open ended range", 65, 150, 65},
		{"lower bound of the whole range", 0, 150, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ageRange, err := NewAgeRange(tt.from, tt.to)
			if err != nil {
				t.Fatalf("failed to new age range: %v", err)
			}

			if got := ageRange.Median(); got != tt.want {
				t.Errorf("Median() = %v, want %v", got, tt.want)
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

func TestAgeRangeDistance(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		from int
		to   int
		age  int
		want int
	}{
		{"an age inside the range is zero", 45, 49, 47, 0},
		{"the lower bound is zero", 45, 49, 45, 0},
		{"the upper bound is zero", 45, 49, 49, 0},
		{"an age below the lower bound", 45, 49, 40, 5},
		{"an age just below the lower bound", 45, 49, 44, 1},
		{"an age above the upper bound", 45, 49, 60, 11},
		{"an age just above the upper bound", 45, 49, 50, 1},
		{"an age below an open ended range", 65, 150, 40, 25},
		{"every age of the whole range is zero", 0, 150, 150, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ageRange, err := NewAgeRange(tt.from, tt.to)
			if err != nil {
				t.Fatalf("failed to new age range: %v", err)
			}

			if got := ageRange.Distance(tt.age); got != tt.want {
				t.Errorf("Distance(%v) = %v, want %v", tt.age, got, tt.want)
			}
		})
	}
}
