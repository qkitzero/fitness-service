package standard

import (
	"context"
	"fmt"
	"sort"

	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/standard"
)

type rankStandardRepository struct {
	db *gorm.DB
}

func NewRankStandardRepository(db *gorm.DB) standard.RankStandardRepository {
	return &rankStandardRepository{db: db}
}

func toRankStandard(m RankStandardModel) (standard.RankStandard, error) {
	rank, err := standard.NewRank(m.Rank.String())
	if err != nil {
		return nil, fmt.Errorf("rank standard %q: %w", m.Rank, err)
	}
	var zScoreMin *standard.ZScore
	if m.ZScoreMin != nil {
		z, err := standard.NewZScore(m.ZScoreMin.Float64())
		if err != nil {
			return nil, fmt.Errorf("rank standard %q: %w", m.Rank, err)
		}
		zScoreMin = &z
	}
	var zScoreMax *standard.ZScore
	if m.ZScoreMax != nil {
		z, err := standard.NewZScore(m.ZScoreMax.Float64())
		if err != nil {
			return nil, fmt.Errorf("rank standard %q: %w", m.Rank, err)
		}
		zScoreMax = &z
	}

	rankStandard, err := standard.NewRankStandard(rank, zScoreMin, zScoreMax, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("rank standard %q: %w", m.Rank, err)
	}

	return rankStandard, nil
}

func (r *rankStandardRepository) List(ctx context.Context) ([]standard.RankStandard, error) {
	var rankStandardModels []RankStandardModel
	if err := r.db.WithContext(ctx).Find(&rankStandardModels).Error; err != nil {
		return nil, err
	}

	rankStandards := make([]standard.RankStandard, 0, len(rankStandardModels))
	for _, rankStandardModel := range rankStandardModels {
		rankStandard, err := toRankStandard(rankStandardModel)
		if err != nil {
			return nil, err
		}
		rankStandards = append(rankStandards, rankStandard)
	}

	sort.Slice(rankStandards, func(i, j int) bool {
		return rankStandards[i].Rank().Order() < rankStandards[j].Rank().Order()
	})

	return rankStandards, nil
}
