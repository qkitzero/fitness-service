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
	ScoreDirection() *ScoreDirection
	Elements() []Element
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type measurementItem struct {
	id             MeasurementItemID
	code           Code
	name           Name
	category       Category
	unit           Unit
	trialCount     TrialCount
	bilateral      bool
	valueType      ValueType
	scoreDirection *ScoreDirection
	elements       []Element
	createdAt      time.Time
	updatedAt      time.Time
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

func (m measurementItem) ScoreDirection() *ScoreDirection {
	if m.scoreDirection == nil {
		return nil
	}
	s := *m.scoreDirection
	return &s
}

func (m measurementItem) Elements() []Element {
	elements := make([]Element, len(m.elements))
	copy(elements, m.elements)
	return elements
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
	scoreDirection *ScoreDirection,
	elements []Element,
	createdAt time.Time,
	updatedAt time.Time,
) MeasurementItem {
	m := &measurementItem{
		id:         id,
		code:       code,
		name:       name,
		category:   category,
		unit:       unit,
		trialCount: trialCount,
		bilateral:  bilateral,
		valueType:  valueType,
		elements:   make([]Element, len(elements)),
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
	copy(m.elements, elements)
	if scoreDirection != nil {
		s := *scoreDirection
		m.scoreDirection = &s
	}
	return m
}
