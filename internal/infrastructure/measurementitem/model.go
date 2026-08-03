package measurementitem

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

type MeasurementItemModel struct {
	ID              measurementitem.MeasurementItemID
	Code            measurementitem.Code
	Name            measurementitem.Name
	Category        measurementitem.Category
	Unit            measurementitem.Unit
	TrialCount      measurementitem.TrialCount
	Bilateral       bool
	ValueType       measurementitem.ValueType
	ScoreDirection  *measurementitem.ScoreDirection
	SideAggregation measurementitem.SideAggregation
	CreatedAt       time.Time                     `gorm:"autoCreateTime:false"`
	UpdatedAt       time.Time                     `gorm:"autoUpdateTime:false"`
	Elements        []MeasurementItemElementModel `gorm:"foreignKey:MeasurementItemID;references:ID"`
}

func (MeasurementItemModel) TableName() string {
	return "measurement_items"
}

type MeasurementItemElementModel struct {
	MeasurementItemID measurementitem.MeasurementItemID `gorm:"primaryKey"`
	Element           measurementitem.Element           `gorm:"primaryKey"`
	CreatedAt         time.Time                         `gorm:"autoCreateTime:false"`
	UpdatedAt         time.Time                         `gorm:"autoUpdateTime:false"`
}

func (MeasurementItemElementModel) TableName() string {
	return "measurement_item_elements"
}
