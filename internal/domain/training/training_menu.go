package training

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

type TrainingMenu interface {
	ID() TrainingMenuID
	Code() Code
	Name() Name
	Element() measurementitem.Element
	Part() Part
	Amount() Amount
	Unit() Unit
	Sets() Sets
	Instruction() Instruction
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type trainingMenu struct {
	id          TrainingMenuID
	code        Code
	name        Name
	element     measurementitem.Element
	part        Part
	amount      Amount
	unit        Unit
	sets        Sets
	instruction Instruction
	createdAt   time.Time
	updatedAt   time.Time
}

func (t trainingMenu) ID() TrainingMenuID {
	return t.id
}

func (t trainingMenu) Code() Code {
	return t.code
}

func (t trainingMenu) Name() Name {
	return t.name
}

func (t trainingMenu) Element() measurementitem.Element {
	return t.element
}

func (t trainingMenu) Part() Part {
	return t.part
}

func (t trainingMenu) Amount() Amount {
	return t.amount
}

func (t trainingMenu) Unit() Unit {
	return t.unit
}

func (t trainingMenu) Sets() Sets {
	return t.sets
}

func (t trainingMenu) Instruction() Instruction {
	return t.instruction
}

func (t trainingMenu) CreatedAt() time.Time {
	return t.createdAt
}

func (t trainingMenu) UpdatedAt() time.Time {
	return t.updatedAt
}

func NewTrainingMenu(
	id TrainingMenuID,
	code Code,
	name Name,
	element measurementitem.Element,
	part Part,
	amount Amount,
	unit Unit,
	sets Sets,
	instruction Instruction,
	createdAt time.Time,
	updatedAt time.Time,
) TrainingMenu {
	return &trainingMenu{
		id:          id,
		code:        code,
		name:        name,
		element:     element,
		part:        part,
		amount:      amount,
		unit:        unit,
		sets:        sets,
		instruction: instruction,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
