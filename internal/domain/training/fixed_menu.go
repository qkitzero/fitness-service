package training

import (
	"time"
)

type FixedMenu interface {
	SortOrder() SortOrder
	TrainingMenuID() TrainingMenuID
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type fixedMenu struct {
	sortOrder      SortOrder
	trainingMenuID TrainingMenuID
	createdAt      time.Time
	updatedAt      time.Time
}

func (f fixedMenu) SortOrder() SortOrder {
	return f.sortOrder
}

func (f fixedMenu) TrainingMenuID() TrainingMenuID {
	return f.trainingMenuID
}

func (f fixedMenu) CreatedAt() time.Time {
	return f.createdAt
}

func (f fixedMenu) UpdatedAt() time.Time {
	return f.updatedAt
}

func NewFixedMenu(
	sortOrder SortOrder,
	trainingMenuID TrainingMenuID,
	createdAt time.Time,
	updatedAt time.Time,
) FixedMenu {
	return &fixedMenu{
		sortOrder:      sortOrder,
		trainingMenuID: trainingMenuID,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}
