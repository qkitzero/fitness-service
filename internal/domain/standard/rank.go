package standard

import (
	"fmt"
	"math"
)

type Rank string

const (
	RankA Rank = "A"
	RankB Rank = "B"
	RankC Rank = "C"
	RankD Rank = "D"
	RankE Rank = "E"
)

var orderedRanks = []Rank{
	RankA,
	RankB,
	RankC,
	RankD,
	RankE,
}

func (r Rank) String() string {
	return string(r)
}

func (r Rank) Order() int {
	for i, rank := range orderedRanks {
		if r == rank {
			return i + 1
		}
	}
	return math.MaxInt
}

func NewRank(s string) (Rank, error) {
	for _, rank := range orderedRanks {
		if Rank(s) == rank {
			return rank, nil
		}
	}
	return Rank(""), fmt.Errorf("invalid rank")
}
