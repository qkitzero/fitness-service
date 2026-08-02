package judgment

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/judgment"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
)

type JudgmentModel struct {
	MeasurementID measurement.MeasurementID `gorm:"primaryKey"`
	Advice        *judgment.Advice
	CreatedAt     time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime:false"`
}

func (JudgmentModel) TableName() string {
	return "judgments"
}
