package training

import (
	"testing"
	"time"
)

func TestNewAgeDecadeMenu(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		decade    int
		sortOrder int
	}{
		{"success new age decade menu", 60, 1},
		{"success new age decade menu of the youngest decade", 10, 5},
		{"success new age decade menu of the oldest decade", 90, 3},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decade, _ := NewDecade(tt.decade)
			sortOrder, _ := NewSortOrder(tt.sortOrder)
			trainingMenuID := NewTrainingMenuID()
			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

			a := NewAgeDecadeMenu(decade, sortOrder, trainingMenuID, createdAt, updatedAt)

			if a.Decade() != decade {
				t.Errorf("Decade() = %v, want %v", a.Decade(), decade)
			}
			if a.SortOrder() != sortOrder {
				t.Errorf("SortOrder() = %v, want %v", a.SortOrder(), sortOrder)
			}
			if a.TrainingMenuID() != trainingMenuID {
				t.Errorf("TrainingMenuID() = %v, want %v", a.TrainingMenuID(), trainingMenuID)
			}
			if !a.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", a.CreatedAt(), createdAt)
			}
			if !a.UpdatedAt().Equal(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", a.UpdatedAt(), updatedAt)
			}
		})
	}
}
