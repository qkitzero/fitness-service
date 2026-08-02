package standard

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/standard"
)

type AgeGroupStandardModel struct {
	ID                standard.AgeGroupStandardID
	MeasurementItemID measurementitem.MeasurementItemID
	Gender            standard.Gender
	AgeFrom           int
	AgeTo             int
	Mean              standard.Mean
	StandardDeviation standard.StandardDeviation
	CreatedAt         time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime:false"`
}

func (AgeGroupStandardModel) TableName() string {
	return "age_group_standards"
}

type RankStandardModel struct {
	Rank      standard.Rank `gorm:"primaryKey"`
	ZScoreMin *standard.ZScore
	ZScoreMax *standard.ZScore
	CreatedAt time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}

func (RankStandardModel) TableName() string {
	return "rank_standards"
}
