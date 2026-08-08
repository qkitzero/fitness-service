package judgment

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/training"
)

type PrescribedMenuOverride interface {
	ID() PrescribedMenuOverrideID
	MeasurementID() measurement.MeasurementID
	SortOrder() training.SortOrder
	Element() *measurementitem.Element
	Part() *training.Part
	TrainingMenuID() training.TrainingMenuID
	Amount() training.Amount
	Unit() training.Unit
	Sets() training.Sets
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type prescribedMenuOverride struct {
	id             PrescribedMenuOverrideID
	measurementID  measurement.MeasurementID
	sortOrder      training.SortOrder
	element        *measurementitem.Element
	part           *training.Part
	trainingMenuID training.TrainingMenuID
	amount         training.Amount
	unit           training.Unit
	sets           training.Sets
	createdAt      time.Time
	updatedAt      time.Time
}

func (p prescribedMenuOverride) ID() PrescribedMenuOverrideID {
	return p.id
}

func (p prescribedMenuOverride) MeasurementID() measurement.MeasurementID {
	return p.measurementID
}

func (p prescribedMenuOverride) SortOrder() training.SortOrder {
	return p.sortOrder
}

func (p prescribedMenuOverride) Element() *measurementitem.Element {
	if p.element == nil {
		return nil
	}
	e := *p.element
	return &e
}

func (p prescribedMenuOverride) Part() *training.Part {
	if p.part == nil {
		return nil
	}
	part := *p.part
	return &part
}

func (p prescribedMenuOverride) TrainingMenuID() training.TrainingMenuID {
	return p.trainingMenuID
}

func (p prescribedMenuOverride) Amount() training.Amount {
	return p.amount
}

func (p prescribedMenuOverride) Unit() training.Unit {
	return p.unit
}

func (p prescribedMenuOverride) Sets() training.Sets {
	return p.sets
}

func (p prescribedMenuOverride) CreatedAt() time.Time {
	return p.createdAt
}

func (p prescribedMenuOverride) UpdatedAt() time.Time {
	return p.updatedAt
}

func NewPrescribedMenuOverride(
	id PrescribedMenuOverrideID,
	measurementID measurement.MeasurementID,
	sortOrder training.SortOrder,
	element *measurementitem.Element,
	part *training.Part,
	trainingMenuID training.TrainingMenuID,
	amount training.Amount,
	unit training.Unit,
	sets training.Sets,
	createdAt time.Time,
	updatedAt time.Time,
) PrescribedMenuOverride {
	p := &prescribedMenuOverride{
		id:             id,
		measurementID:  measurementID,
		sortOrder:      sortOrder,
		trainingMenuID: trainingMenuID,
		amount:         amount,
		unit:           unit,
		sets:           sets,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
	if element != nil {
		e := *element
		p.element = &e
	}
	if part != nil {
		pt := *part
		p.part = &pt
	}
	return p
}
