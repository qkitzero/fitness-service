package measurement

import (
	"time"

	"github.com/google/uuid"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/staff"
)

type MeasurementModel struct {
	ID               measurement.MeasurementID
	CustomerID       customer.CustomerID
	MeasuredOn       measurement.MeasuredOn
	MeasuredBy       staff.StaffID
	AgeAtMeasurement measurement.AgeAtMeasurement
	UpdatedBy        staff.StaffID
	IsDraft          bool
	CreatedAt        time.Time               `gorm:"autoCreateTime:false"`
	UpdatedAt        time.Time               `gorm:"autoUpdateTime:false"`
	Entries          []MeasurementEntryModel `gorm:"foreignKey:MeasurementID;references:ID"`
}

func (MeasurementModel) TableName() string {
	return "measurements"
}

type MeasurementEntryModel struct {
	ID                uuid.UUID
	MeasurementID     measurement.MeasurementID
	MeasurementItemID measurementitem.MeasurementItemID
	Unmeasurable      bool
	Note              *measurement.Note
	CreatedAt         time.Time               `gorm:"autoCreateTime:false"`
	UpdatedAt         time.Time               `gorm:"autoUpdateTime:false"`
	Values            []MeasurementValueModel `gorm:"foreignKey:MeasurementEntryID;references:ID"`
}

func (MeasurementEntryModel) TableName() string {
	return "measurement_entries"
}

type MeasurementValueModel struct {
	ID                 uuid.UUID
	MeasurementEntryID uuid.UUID
	TrialIndex         measurement.TrialIndex
	Side               measurement.Side
	Value              *measurement.Value
	ValueSecondary     *measurement.Value
	ValueChoice        *measurement.Choice
	CreatedAt          time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime:false"`
}

func (MeasurementValueModel) TableName() string {
	return "measurement_values"
}
