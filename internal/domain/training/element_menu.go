package training

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

type ElementMenu interface {
	Element() measurementitem.Element
	Part() Part
	Level() Level
	TrainingMenuID() TrainingMenuID
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type elementMenu struct {
	element        measurementitem.Element
	part           Part
	level          Level
	trainingMenuID TrainingMenuID
	createdAt      time.Time
	updatedAt      time.Time
}

func (e elementMenu) Element() measurementitem.Element {
	return e.element
}

func (e elementMenu) Part() Part {
	return e.part
}

func (e elementMenu) Level() Level {
	return e.level
}

func (e elementMenu) TrainingMenuID() TrainingMenuID {
	return e.trainingMenuID
}

func (e elementMenu) CreatedAt() time.Time {
	return e.createdAt
}

func (e elementMenu) UpdatedAt() time.Time {
	return e.updatedAt
}

func NewElementMenu(
	element measurementitem.Element,
	part Part,
	level Level,
	trainingMenuID TrainingMenuID,
	createdAt time.Time,
	updatedAt time.Time,
) ElementMenu {
	return &elementMenu{
		element:        element,
		part:           part,
		level:          level,
		trainingMenuID: trainingMenuID,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}
