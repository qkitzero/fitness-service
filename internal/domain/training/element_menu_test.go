package training

import (
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

func TestNewElementMenu(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		element measurementitem.Element
		part    string
		level   int
	}{
		{"success new element menu", measurementitem.ElementMuscleStrength, "upper_limb", 2},
		{"success new element menu of the lowest level", measurementitem.ElementBalance, "lower_limb", 1},
		{"success new element menu of the highest level", measurementitem.ElementMobility, "whole_body", 5},
		{"success new element menu of another element", measurementitem.ElementFlexibility, "lower_limb", 3},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			part, _ := NewPart(tt.part)
			level, _ := NewLevel(tt.level)
			trainingMenuID := NewTrainingMenuID()
			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

			e := NewElementMenu(tt.element, part, level, trainingMenuID, createdAt, updatedAt)

			if e.Element() != tt.element {
				t.Errorf("Element() = %v, want %v", e.Element(), tt.element)
			}
			if e.Part() != part {
				t.Errorf("Part() = %v, want %v", e.Part(), part)
			}
			if e.Level() != level {
				t.Errorf("Level() = %v, want %v", e.Level(), level)
			}
			if e.TrainingMenuID() != trainingMenuID {
				t.Errorf("TrainingMenuID() = %v, want %v", e.TrainingMenuID(), trainingMenuID)
			}
			if !e.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", e.CreatedAt(), createdAt)
			}
			if !e.UpdatedAt().Equal(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", e.UpdatedAt(), updatedAt)
			}
		})
	}
}
