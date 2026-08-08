package judgment

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/judgment"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/training"
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

type PrescribedMenuOverrideModel struct {
	ID             judgment.PrescribedMenuOverrideID `gorm:"primaryKey"`
	MeasurementID  measurement.MeasurementID
	SortOrder      training.SortOrder
	Element        *measurementitem.Element
	Part           *training.Part
	TrainingMenuID training.TrainingMenuID
	Amount         training.Amount
	Unit           training.Unit
	Sets           training.Sets
	CreatedAt      time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime:false"`
}

func (PrescribedMenuOverrideModel) TableName() string {
	return "prescribed_menu_overrides"
}
