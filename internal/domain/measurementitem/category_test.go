package measurementitem

import (
	"math"
	"testing"
)

func TestNewCategory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		success  bool
		category string
		want     Category
	}{
		{"success vital", true, "vital", CategoryVital},
		{"success physique", true, "physique", CategoryPhysique},
		{"success body composition", true, "body_composition", CategoryBodyComposition},
		{"success motor function", true, "motor_function", CategoryMotorFunction},
		{"failure empty category", false, "", ""},
		{"failure invalid category", false, "unknown", ""},
		{"failure uppercase category", false, "VITAL", ""},
		{"failure untrimmed category", false, " vital ", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			category, err := NewCategory(tt.category)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && category != tt.want {
				t.Errorf("NewCategory() = %v, want %v", category, tt.want)
			}

			if tt.success && category.String() != tt.category {
				t.Errorf("String() = %v, want %v", category.String(), tt.category)
			}
		})
	}
}

func TestCategoryOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		category Category
		want     int
	}{
		{"vital comes first", CategoryVital, 1},
		{"physique comes second", CategoryPhysique, 2},
		{"body composition comes third", CategoryBodyComposition, 3},
		{"motor function comes fourth", CategoryMotorFunction, 4},
		{"unknown category comes last", Category("unknown"), math.MaxInt},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.category.Order(); got != tt.want {
				t.Errorf("Order() = %v, want %v", got, tt.want)
			}
		})
	}
}
