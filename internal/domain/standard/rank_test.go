package standard

import (
	"math"
	"testing"
)

func TestNewRank(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		rank    string
		want    Rank
	}{
		{"success rank a", true, "A", RankA},
		{"success rank b", true, "B", RankB},
		{"success rank c", true, "C", RankC},
		{"success rank d", true, "D", RankD},
		{"success rank e", true, "E", RankE},
		{"failure empty rank", false, "", ""},
		{"failure invalid rank", false, "F", ""},
		{"failure lowercase rank", false, "a", ""},
		{"failure untrimmed rank", false, " A ", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rank, err := NewRank(tt.rank)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && rank != tt.want {
				t.Errorf("NewRank() = %v, want %v", rank, tt.want)
			}

			if tt.success && rank.String() != tt.rank {
				t.Errorf("String() = %v, want %v", rank.String(), tt.rank)
			}
		})
	}
}

func TestRankOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		rank Rank
		want int
	}{
		{"rank a comes first", RankA, 1},
		{"rank b comes second", RankB, 2},
		{"rank c comes third", RankC, 3},
		{"rank d comes fourth", RankD, 4},
		{"rank e comes fifth", RankE, 5},
		{"unknown rank comes last", Rank("F"), math.MaxInt},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.rank.Order(); got != tt.want {
				t.Errorf("Order() = %v, want %v", got, tt.want)
			}
		})
	}
}
