package training

import (
	"testing"
	"time"
)

func TestNewFixedMenu(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		sortOrder int
	}{
		{"success new fixed menu", 1},
		{"success new fixed menu of a later sort order", 5},
		{"success new fixed menu of a middle sort order", 3},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sortOrder, _ := NewSortOrder(tt.sortOrder)
			trainingMenuID := NewTrainingMenuID()
			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

			f := NewFixedMenu(sortOrder, trainingMenuID, createdAt, updatedAt)

			if f.SortOrder() != sortOrder {
				t.Errorf("SortOrder() = %v, want %v", f.SortOrder(), sortOrder)
			}
			if f.TrainingMenuID() != trainingMenuID {
				t.Errorf("TrainingMenuID() = %v, want %v", f.TrainingMenuID(), trainingMenuID)
			}
			if !f.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", f.CreatedAt(), createdAt)
			}
			if !f.UpdatedAt().Equal(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", f.UpdatedAt(), updatedAt)
			}
		})
	}
}
