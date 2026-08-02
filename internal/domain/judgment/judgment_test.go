package judgment

import (
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurement"
)

func TestNewJudgment(t *testing.T) {
	t.Parallel()
	advice, _ := NewAdvice("週2回のスクワットを継続してください")
	tests := []struct {
		name        string
		advice      *Advice
		reconstruct bool
		want        string
	}{
		{"success new judgment with advice", advice, false, "週2回のスクワットを継続してください"},
		{"success new judgment without advice", nil, false, ""},
		{"success reconstruct judgment with advice", advice, true, "週2回のスクワットを継続してください"},
		{"success reconstruct judgment without advice", nil, true, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			measurementID := measurement.NewMeasurementID()
			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

			j := NewJudgment(measurementID, tt.advice, createdAt, updatedAt)
			if tt.reconstruct {
				j = ReconstructJudgment(measurementID, tt.advice, createdAt, updatedAt)
			}

			if j.MeasurementID() != measurementID {
				t.Errorf("MeasurementID() = %v, want %v", j.MeasurementID(), measurementID)
			}
			if tt.advice == nil {
				if j.Advice() != nil {
					t.Errorf("Advice() = %v, want nil", *j.Advice())
				}
			} else if j.Advice() == nil || j.Advice().String() != tt.want {
				t.Errorf("Advice() = %v, want %v", j.Advice(), tt.want)
			}
			if !j.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", j.CreatedAt(), createdAt)
			}
			if !j.UpdatedAt().Equal(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", j.UpdatedAt(), updatedAt)
			}
		})
	}
}

func TestJudgmentUpdate(t *testing.T) {
	t.Parallel()
	advice, _ := NewAdvice("週2回のスクワットを継続してください")
	updatedAdvice, _ := NewAdvice("週3回のスクワットに増やしてください")
	tests := []struct {
		name    string
		advice  *Advice
		update  *Advice
		wantNil bool
		want    string
	}{
		{"success update advice", advice, updatedAdvice, false, "週3回のスクワットに増やしてください"},
		{"success set advice on a judgment without advice", nil, advice, false, "週2回のスクワットを継続してください"},
		{"success clear advice", advice, nil, true, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

			j := NewJudgment(measurement.NewMeasurementID(), tt.advice, createdAt, updatedAt)

			j.Update(tt.update)

			if tt.wantNil {
				if j.Advice() != nil {
					t.Errorf("Advice() = %v, want nil", *j.Advice())
				}
			} else if j.Advice() == nil || j.Advice().String() != tt.want {
				t.Errorf("Advice() = %v, want %v", j.Advice(), tt.want)
			}
			if !j.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", j.CreatedAt(), createdAt)
			}
			if !j.UpdatedAt().After(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want after %v", j.UpdatedAt(), updatedAt)
			}
		})
	}
}

func TestJudgmentIsolatesMutableState(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		mutate func(j Judgment, advice *Advice)
	}{
		{
			name: "mutating the constructor argument does not affect the judgment",
			mutate: func(_ Judgment, advice *Advice) {
				*advice = Advice("改ざんされたアドバイス")
			},
		},
		{
			name: "mutating the getter result does not affect the judgment",
			mutate: func(j Judgment, _ *Advice) {
				*j.Advice() = Advice("改ざんされたアドバイス")
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
			advice := Advice("週2回のスクワットを継続してください")

			j := NewJudgment(measurement.NewMeasurementID(), &advice, createdAt, updatedAt)

			tt.mutate(j, &advice)

			if j.Advice() == nil || j.Advice().String() != "週2回のスクワットを継続してください" {
				t.Errorf("Advice() = %v, want %v", j.Advice(), "週2回のスクワットを継続してください")
			}
		})
	}
}
