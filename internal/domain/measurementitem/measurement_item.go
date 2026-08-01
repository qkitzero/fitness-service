package measurementitem

import (
	"time"
)

type MeasurementItem interface {
	ID() MeasurementItemID
	Code() Code
	Name() Name
	Category() Category
	Unit() Unit
	TrialCount() TrialCount
	Bilateral() bool
	ValueType() ValueType
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type measurementItem struct {
	id         MeasurementItemID
	code       Code
	name       Name
	category   Category
	unit       Unit
	trialCount TrialCount
	bilateral  bool
	valueType  ValueType
	createdAt  time.Time
	updatedAt  time.Time
}

func (m measurementItem) ID() MeasurementItemID {
	return m.id
}

func (m measurementItem) Code() Code {
	return m.code
}

func (m measurementItem) Name() Name {
	return m.name
}

func (m measurementItem) Category() Category {
	return m.category
}

func (m measurementItem) Unit() Unit {
	return m.unit
}

func (m measurementItem) TrialCount() TrialCount {
	return m.trialCount
}

func (m measurementItem) Bilateral() bool {
	return m.bilateral
}

func (m measurementItem) ValueType() ValueType {
	return m.valueType
}

func (m measurementItem) CreatedAt() time.Time {
	return m.createdAt
}

func (m measurementItem) UpdatedAt() time.Time {
	return m.updatedAt
}

func NewMeasurementItem(
	id MeasurementItemID,
	code Code,
	name Name,
	category Category,
	unit Unit,
	trialCount TrialCount,
	bilateral bool,
	valueType ValueType,
	createdAt time.Time,
	updatedAt time.Time,
) MeasurementItem {
	return &measurementItem{
		id:         id,
		code:       code,
		name:       name,
		category:   category,
		unit:       unit,
		trialCount: trialCount,
		bilateral:  bilateral,
		valueType:  valueType,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}
