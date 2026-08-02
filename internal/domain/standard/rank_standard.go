package standard

import (
	"time"
)

type RankStandard interface {
	Rank() Rank
	ZScoreMin() *ZScore
	ZScoreMax() *ZScore
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type rankStandard struct {
	rank      Rank
	zScoreMin *ZScore
	zScoreMax *ZScore
	createdAt time.Time
	updatedAt time.Time
}

func (r rankStandard) Rank() Rank {
	return r.rank
}

func (r rankStandard) ZScoreMin() *ZScore {
	if r.zScoreMin == nil {
		return nil
	}
	z := *r.zScoreMin
	return &z
}

func (r rankStandard) ZScoreMax() *ZScore {
	if r.zScoreMax == nil {
		return nil
	}
	z := *r.zScoreMax
	return &z
}

func (r rankStandard) CreatedAt() time.Time {
	return r.createdAt
}

func (r rankStandard) UpdatedAt() time.Time {
	return r.updatedAt
}

func NewRankStandard(
	rank Rank,
	zScoreMin *ZScore,
	zScoreMax *ZScore,
	createdAt time.Time,
	updatedAt time.Time,
) (RankStandard, error) {
	if zScoreMin == nil && zScoreMax == nil {
		return nil, ErrInvalidZScoreRange
	}
	if zScoreMin != nil && zScoreMax != nil && *zScoreMin >= *zScoreMax {
		return nil, ErrInvalidZScoreRange
	}

	r := &rankStandard{
		rank:      rank,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	if zScoreMin != nil {
		z := *zScoreMin
		r.zScoreMin = &z
	}
	if zScoreMax != nil {
		z := *zScoreMax
		r.zScoreMax = &z
	}
	return r, nil
}
