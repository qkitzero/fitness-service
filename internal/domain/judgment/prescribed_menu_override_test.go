package judgment

import (
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/training"
)

func TestNewPrescribedMenuOverride(t *testing.T) {
	t.Parallel()
	muscleStrength := measurementitem.ElementMuscleStrength
	upperLimb := training.PartUpperLimb
	tests := []struct {
		name      string
		sortOrder int
		element   *measurementitem.Element
		part      *training.Part
		amount    int
		unit      training.Unit
		sets      int
	}{
		{"success new prescribed menu override", 1, &muscleStrength, &upperLimb, 10, training.UnitReps, 3},
		{"success new prescribed menu override without labels", 2, nil, nil, 30, training.UnitSeconds, 1},
		{"success new prescribed menu override of the lowest bounds", 3, &muscleStrength, &upperLimb, 1, training.UnitMinutes, 1},
		{"success new prescribed menu override of the highest bounds", 4, nil, nil, 999, training.UnitReps, 99},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewPrescribedMenuOverrideID()
			measurementID := measurement.NewMeasurementID()
			sortOrder, _ := training.NewSortOrder(tt.sortOrder)
			trainingMenuID := training.NewTrainingMenuID()
			amount, _ := training.NewAmount(tt.amount)
			sets, _ := training.NewSets(tt.sets)
			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

			o := NewPrescribedMenuOverride(id, measurementID, sortOrder, tt.element, tt.part, trainingMenuID, amount, tt.unit, sets, createdAt, updatedAt)

			if o.ID() != id {
				t.Errorf("ID() = %v, want %v", o.ID(), id)
			}
			if o.MeasurementID() != measurementID {
				t.Errorf("MeasurementID() = %v, want %v", o.MeasurementID(), measurementID)
			}
			if o.SortOrder() != sortOrder {
				t.Errorf("SortOrder() = %v, want %v", o.SortOrder(), sortOrder)
			}
			switch {
			case tt.element == nil:
				if o.Element() != nil {
					t.Errorf("Element() = %v, want nil", *o.Element())
				}
			case o.Element() == nil:
				t.Errorf("Element() = nil, want %v", *tt.element)
			case *o.Element() != *tt.element:
				t.Errorf("Element() = %v, want %v", *o.Element(), *tt.element)
			}
			switch {
			case tt.part == nil:
				if o.Part() != nil {
					t.Errorf("Part() = %v, want nil", *o.Part())
				}
			case o.Part() == nil:
				t.Errorf("Part() = nil, want %v", *tt.part)
			case *o.Part() != *tt.part:
				t.Errorf("Part() = %v, want %v", *o.Part(), *tt.part)
			}
			if o.TrainingMenuID() != trainingMenuID {
				t.Errorf("TrainingMenuID() = %v, want %v", o.TrainingMenuID(), trainingMenuID)
			}
			if o.Amount() != amount {
				t.Errorf("Amount() = %v, want %v", o.Amount(), amount)
			}
			if o.Unit() != tt.unit {
				t.Errorf("Unit() = %v, want %v", o.Unit(), tt.unit)
			}
			if o.Sets() != sets {
				t.Errorf("Sets() = %v, want %v", o.Sets(), sets)
			}
			if !o.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", o.CreatedAt(), createdAt)
			}
			if !o.UpdatedAt().Equal(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", o.UpdatedAt(), updatedAt)
			}
		})
	}
}

func TestPrescribedMenuOverrideLabelsAreCopied(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		element measurementitem.Element
		part    training.Part
	}{
		{"success mutating the caller labels does not affect the override", measurementitem.ElementMuscleStrength, training.PartUpperLimb},
		{"success mutating another pair of caller labels does not affect the override", measurementitem.ElementBalance, training.PartWholeBody},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sortOrder, _ := training.NewSortOrder(1)
			amount, _ := training.NewAmount(10)
			sets, _ := training.NewSets(3)
			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			element, part := tt.element, tt.part

			o := NewPrescribedMenuOverride(NewPrescribedMenuOverrideID(), measurement.NewMeasurementID(), sortOrder, &element, &part, training.NewTrainingMenuID(), amount, training.UnitReps, sets, createdAt, createdAt)

			element = measurementitem.ElementMobility
			part = training.PartLowerLimb

			if *o.Element() != tt.element {
				t.Errorf("Element() = %v, want %v", *o.Element(), tt.element)
			}
			if *o.Part() != tt.part {
				t.Errorf("Part() = %v, want %v", *o.Part(), tt.part)
			}

			*o.Element() = measurementitem.ElementAgility
			*o.Part() = training.PartLowerLimb

			if *o.Element() != tt.element {
				t.Errorf("Element() = %v, want %v", *o.Element(), tt.element)
			}
			if *o.Part() != tt.part {
				t.Errorf("Part() = %v, want %v", *o.Part(), tt.part)
			}
		})
	}
}
