package measurementitem

import (
	"testing"
	"time"
)

func TestNewMeasurementItem(t *testing.T) {
	t.Parallel()
	higherIsBetter := ScoreDirectionHigherIsBetter
	lowerIsBetter := ScoreDirectionLowerIsBetter
	tests := []struct {
		name            string
		code            string
		measurementName string
		category        string
		unit            string
		trialCount      int
		sideMode        SideMode
		valueType       string
		scoreDirection  *ScoreDirection
		sideAggregation SideAggregation
		elements        []Element
	}{
		{"success new measurement item", "grip_strength", "握力", "motor_function", "kg", 2, SideModeBilateral, "numeric", &higherIsBetter, SideAggregationMean, []Element{ElementMuscleStrength}},
		{"success new measurement item without sides", "blood_pressure", "血圧", "vital", "mmHg", 1, SideModeNone, "paired", nil, SideAggregationMean, nil},
		{"success new measurement item of another category", "body_fat_percentage", "体脂肪率", "body_composition", "percent", 1, SideModeNone, "numeric", nil, SideAggregationMean, []Element{}},
		{"success new measurement item scored lower is better", "walk_5m", "5m歩行", "motor_function", "sec", 2, SideModeNone, "numeric", &lowerIsBetter, SideAggregationMean, []Element{ElementMobility}},
		{"success new measurement item aggregated by the best side", "eyes_open_one_leg_stand", "開眼片足立ち", "motor_function", "sec", 2, SideModeBilateral, "numeric", &higherIsBetter, SideAggregationBest, []Element{ElementBalance}},
		{"success new measurement item of multiple elements", "multi_element_item", "複数要素の測定項目", "motor_function", "count", 1, SideModeNone, "numeric", &higherIsBetter, SideAggregationMean, []Element{ElementAgility, ElementMobility}},
		{"success new measurement item of optional sides aggregated by the worst side", "stand_up_test", "立ち上がり", "motor_function", "level", 1, SideModeOptionalBilateral, "numeric", &higherIsBetter, SideAggregationWorst, []Element{ElementMuscleStrength}},
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

			m := NewMeasurementItem(id, code, measurementName, category, unit, trialCount, tt.sideMode, valueType, tt.scoreDirection, tt.sideAggregation, tt.elements, createdAt, updatedAt)

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
			if m.SideMode() != tt.sideMode {
				t.Errorf("SideMode() = %v, want %v", m.SideMode(), tt.sideMode)
			}
			if m.ValueType() != valueType {
				t.Errorf("ValueType() = %v, want %v", m.ValueType(), valueType)
			}
			if m.SideAggregation() != tt.sideAggregation {
				t.Errorf("SideAggregation() = %v, want %v", m.SideAggregation(), tt.sideAggregation)
			}
			switch {
			case tt.scoreDirection == nil:
				if m.ScoreDirection() != nil {
					t.Errorf("ScoreDirection() = %v, want nil", m.ScoreDirection())
				}
			case m.ScoreDirection() == nil:
				t.Errorf("ScoreDirection() = nil, want %v", *tt.scoreDirection)
			case *m.ScoreDirection() != *tt.scoreDirection:
				t.Errorf("ScoreDirection() = %v, want %v", *m.ScoreDirection(), *tt.scoreDirection)
			}
			if len(m.Elements()) != len(tt.elements) {
				t.Errorf("len(Elements()) = %v, want %v", len(m.Elements()), len(tt.elements))
			}
			for i, element := range tt.elements {
				if i >= len(m.Elements()) {
					break
				}
				if m.Elements()[i] != element {
					t.Errorf("Elements()[%d] = %v, want %v", i, m.Elements()[i], element)
				}
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

func TestMeasurementItemIsolatesMutableState(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		mutate func(m MeasurementItem, scoreDirection *ScoreDirection, elements []Element)
	}{
		{
			name: "mutating the constructor arguments does not affect the measurement item",
			mutate: func(_ MeasurementItem, scoreDirection *ScoreDirection, elements []Element) {
				*scoreDirection = ScoreDirectionLowerIsBetter
				elements[0] = ElementMobility
			},
		},
		{
			name: "mutating the getter results does not affect the measurement item",
			mutate: func(m MeasurementItem, _ *ScoreDirection, _ []Element) {
				*m.ScoreDirection() = ScoreDirectionLowerIsBetter
				m.Elements()[0] = ElementMobility
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, _ := NewCode("grip_strength")
			measurementName, _ := NewName("握力")
			category, _ := NewCategory("motor_function")
			unit, _ := NewUnit("kg")
			trialCount, _ := NewTrialCount(2)
			valueType, _ := NewValueType("numeric")
			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
			scoreDirection := ScoreDirectionHigherIsBetter
			elements := []Element{ElementMuscleStrength}

			m := NewMeasurementItem(NewMeasurementItemID(), code, measurementName, category, unit, trialCount, SideModeBilateral, valueType, &scoreDirection, SideAggregationMean, elements, createdAt, updatedAt)

			tt.mutate(m, &scoreDirection, elements)

			if *m.ScoreDirection() != ScoreDirectionHigherIsBetter {
				t.Errorf("ScoreDirection() = %v, want %v", *m.ScoreDirection(), ScoreDirectionHigherIsBetter)
			}
			if m.Elements()[0] != ElementMuscleStrength {
				t.Errorf("Elements()[0] = %v, want %v", m.Elements()[0], ElementMuscleStrength)
			}
		})
	}
}
