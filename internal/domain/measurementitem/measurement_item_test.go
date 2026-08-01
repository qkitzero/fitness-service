package measurementitem

import (
	"testing"
	"time"
)

func TestNewMeasurementItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		code            string
		measurementName string
		category        string
		unit            string
		trialCount      int
		bilateral       bool
		valueType       string
	}{
		{"success new measurement item", "grip_strength", "握力", "motor_function", "kg", 2, true, "numeric"},
		{"success new measurement item without bilateral", "blood_pressure", "血圧", "vital", "mmHg", 1, false, "paired"},
		{"success new measurement item of another category", "body_fat_percentage", "体脂肪率", "body_composition", "percent", 1, false, "numeric"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewMeasurementItemID()
			code, _ := NewCode(tt.code)
			measurementName, _ := NewName(tt.measurementName)
			category, _ := NewCategory(tt.category)
			unit, _ := NewUnit(tt.unit)
			trialCount, _ := NewTrialCount(tt.trialCount)
			valueType, _ := NewValueType(tt.valueType)
			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

			m := NewMeasurementItem(id, code, measurementName, category, unit, trialCount, tt.bilateral, valueType, createdAt, updatedAt)

			if m.ID() != id {
				t.Errorf("ID() = %v, want %v", m.ID(), id)
			}
			if m.Code() != code {
				t.Errorf("Code() = %v, want %v", m.Code(), code)
			}
			if m.Name() != measurementName {
				t.Errorf("Name() = %v, want %v", m.Name(), measurementName)
			}
			if m.Category() != category {
				t.Errorf("Category() = %v, want %v", m.Category(), category)
			}
			if m.Unit() != unit {
				t.Errorf("Unit() = %v, want %v", m.Unit(), unit)
			}
			if m.TrialCount() != trialCount {
				t.Errorf("TrialCount() = %v, want %v", m.TrialCount(), trialCount)
			}
			if m.Bilateral() != tt.bilateral {
				t.Errorf("Bilateral() = %v, want %v", m.Bilateral(), tt.bilateral)
			}
			if m.ValueType() != valueType {
				t.Errorf("ValueType() = %v, want %v", m.ValueType(), valueType)
			}
			if !m.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", m.CreatedAt(), createdAt)
			}
			if !m.UpdatedAt().Equal(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", m.UpdatedAt(), updatedAt)
			}
		})
	}
}
