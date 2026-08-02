package standard

import (
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

func TestNewAgeGroupStandard(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		gender            string
		ageFrom           int
		ageTo             int
		mean              float64
		standardDeviation float64
	}{
		{"success new age group standard", "male", 45, 49, 46.5, 6.8},
		{"success new age group standard of female", "female", 45, 49, 27.7, 5.1},
		{"success new age group standard of a single age", "male", 20, 20, 47.31, 7.02},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewAgeGroupStandardID()
			measurementItemID := measurementitem.NewMeasurementItemID()
			gender, _ := NewGender(tt.gender)
			ageRange, _ := NewAgeRange(tt.ageFrom, tt.ageTo)
			mean, _ := NewMean(tt.mean)
			standardDeviation, _ := NewStandardDeviation(tt.standardDeviation)
			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

			a := NewAgeGroupStandard(id, measurementItemID, gender, ageRange, mean, standardDeviation, createdAt, updatedAt)

			if a.ID() != id {
				t.Errorf("ID() = %v, want %v", a.ID(), id)
			}
			if a.MeasurementItemID() != measurementItemID {
				t.Errorf("MeasurementItemID() = %v, want %v", a.MeasurementItemID(), measurementItemID)
			}
			if a.Gender() != gender {
				t.Errorf("Gender() = %v, want %v", a.Gender(), gender)
			}
			if a.AgeRange() != ageRange {
				t.Errorf("AgeRange() = %v, want %v", a.AgeRange(), ageRange)
			}
			if a.Mean() != mean {
				t.Errorf("Mean() = %v, want %v", a.Mean(), mean)
			}
			if a.StandardDeviation() != standardDeviation {
				t.Errorf("StandardDeviation() = %v, want %v", a.StandardDeviation(), standardDeviation)
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
