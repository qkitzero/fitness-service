package training

import (
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

func TestNewTrainingMenu(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		code         string
		trainingName string
		element      measurementitem.Element
		part         string
		amount       int
		unit         string
		sets         int
		instruction  string
	}{
		{"success new training menu", "wall_push_up", "壁押し", measurementitem.ElementMuscleStrength, "upper_limb", 10, "reps", 3, "壁に手をついて肘を曲げ伸ばしする。"},
		{"success new training menu measured in seconds", "front_plank", "フロントプランク", measurementitem.ElementMuscleStrength, "whole_body", 30, "seconds", 3, "肘とつま先で体を支え、体を一直線に保つ。"},
		{"success new training menu measured in minutes", "indoor_walk", "室内歩行", measurementitem.ElementMobility, "whole_body", 5, "minutes", 1, "手すりや壁の近くを、無理のない速さで歩く。"},
		{"success new training menu of the lower limb", "squat", "スクワット", measurementitem.ElementMuscleStrength, "lower_limb", 10, "reps", 3, "足を肩幅に開き、膝がつま先より前に出ないように腰を落とす。"},
		{"success new training menu of another element", "tandem_stand", "タンデム立位", measurementitem.ElementBalance, "whole_body", 30, "seconds", 3, "片足のかかとにもう一方のつま先をつけて立ち、姿勢を保つ。"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewTrainingMenuID()
			code, _ := NewCode(tt.code)
			trainingName, _ := NewName(tt.trainingName)
			part, _ := NewPart(tt.part)
			amount, _ := NewAmount(tt.amount)
			unit, _ := NewUnit(tt.unit)
			sets, _ := NewSets(tt.sets)
			instruction, _ := NewInstruction(tt.instruction)
			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

			m := NewTrainingMenu(id, code, trainingName, tt.element, part, amount, unit, sets, instruction, createdAt, updatedAt)

			if m.ID() != id {
				t.Errorf("ID() = %v, want %v", m.ID(), id)
			}
			if m.Code() != code {
				t.Errorf("Code() = %v, want %v", m.Code(), code)
			}
			if m.Name() != trainingName {
				t.Errorf("Name() = %v, want %v", m.Name(), trainingName)
			}
			if m.Element() != tt.element {
				t.Errorf("Element() = %v, want %v", m.Element(), tt.element)
			}
			if m.Part() != part {
				t.Errorf("Part() = %v, want %v", m.Part(), part)
			}
			if m.Amount() != amount {
				t.Errorf("Amount() = %v, want %v", m.Amount(), amount)
			}
			if m.Unit() != unit {
				t.Errorf("Unit() = %v, want %v", m.Unit(), unit)
			}
			if m.Sets() != sets {
				t.Errorf("Sets() = %v, want %v", m.Sets(), sets)
			}
			if m.Instruction() != instruction {
				t.Errorf("Instruction() = %v, want %v", m.Instruction(), instruction)
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
