package measurementitem

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

type MeasurementItemModel struct {
	ID         measurementitem.MeasurementItemID
	Code       measurementitem.Code
	Name       measurementitem.Name
	Category   measurementitem.Category
	Unit       measurementitem.Unit
	TrialCount measurementitem.TrialCount
	Bilateral  bool
	ValueType  measurementitem.ValueType
	CreatedAt  time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime:false"`
}

func (MeasurementItemModel) TableName() string {
	return "measurement_items"
}
