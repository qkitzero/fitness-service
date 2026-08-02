package standard

import (
	"testing"
	"time"
)

func TestNewRankStandard(t *testing.T) {
	t.Parallel()
	zScoreMinA := ZScore(1.5)
	zScoreMinB := ZScore(0.5)
	zScoreMaxB := ZScore(1.5)
	zScoreMaxE := ZScore(-1.5)
	tests := []struct {
		name      string
		success   bool
		rank      Rank
		zScoreMin *ZScore
		zScoreMax *ZScore
	}{
		{"success rank without an upper bound", true, RankA, &zScoreMinA, nil},
		{"success rank with both bounds", true, RankB, &zScoreMinB, &zScoreMaxB},
		{"success rank without a lower bound", true, RankE, nil, &zScoreMaxE},
		{"failure rank without bounds", false, RankC, nil, nil},
		{"failure reversed bounds", false, RankB, &zScoreMaxB, &zScoreMinB},
		{"failure equal bounds", false, RankB, &zScoreMinB, &zScoreMinB},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

			r, err := NewRankStandard(tt.rank, tt.zScoreMin, tt.zScoreMax, createdAt, updatedAt)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success {
				return
			}

			if r.Rank() != tt.rank {
				t.Errorf("Rank() = %v, want %v", r.Rank(), tt.rank)
			}
			switch {
			case tt.zScoreMin == nil:
				if r.ZScoreMin() != nil {
					t.Errorf("ZScoreMin() = %v, want nil", r.ZScoreMin())
				}
			case r.ZScoreMin() == nil:
				t.Errorf("ZScoreMin() = nil, want %v", *tt.zScoreMin)
			case *r.ZScoreMin() != *tt.zScoreMin:
				t.Errorf("ZScoreMin() = %v, want %v", *r.ZScoreMin(), *tt.zScoreMin)
			}
			switch {
			case tt.zScoreMax == nil:
				if r.ZScoreMax() != nil {
					t.Errorf("ZScoreMax() = %v, want nil", r.ZScoreMax())
				}
			case r.ZScoreMax() == nil:
				t.Errorf("ZScoreMax() = nil, want %v", *tt.zScoreMax)
			case *r.ZScoreMax() != *tt.zScoreMax:
				t.Errorf("ZScoreMax() = %v, want %v", *r.ZScoreMax(), *tt.zScoreMax)
			}
			if !r.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", r.CreatedAt(), createdAt)
			}
			if !r.UpdatedAt().Equal(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", r.UpdatedAt(), updatedAt)
			}
		})
	}
}

func TestRankStandardIsolatesMutableState(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		mutate func(r RankStandard, zScoreMin *ZScore, zScoreMax *ZScore)
	}{
		{
			name: "mutating the constructor arguments does not affect the rank standard",
			mutate: func(_ RankStandard, zScoreMin *ZScore, zScoreMax *ZScore) {
				*zScoreMin = ZScore(-99.99)
				*zScoreMax = ZScore(99.99)
			},
		},
		{
			name: "mutating the getter results does not affect the rank standard",
			mutate: func(r RankStandard, _ *ZScore, _ *ZScore) {
				*r.ZScoreMin() = ZScore(-99.99)
				*r.ZScoreMax() = ZScore(99.99)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
			zScoreMin := ZScore(0.5)
			zScoreMax := ZScore(1.5)

			r, err := NewRankStandard(RankB, &zScoreMin, &zScoreMax, createdAt, updatedAt)
			if err != nil {
				t.Fatalf("failed to new rank standard: %v", err)
			}

			tt.mutate(r, &zScoreMin, &zScoreMax)

			if *r.ZScoreMin() != ZScore(0.5) {
				t.Errorf("ZScoreMin() = %v, want %v", *r.ZScoreMin(), ZScore(0.5))
			}
			if *r.ZScoreMax() != ZScore(1.5) {
				t.Errorf("ZScoreMax() = %v, want %v", *r.ZScoreMax(), ZScore(1.5))
			}
		})
	}
}
