package training

import (
	"time"
)

type AgeDecadeMenu interface {
	Decade() Decade
	SortOrder() SortOrder
	TrainingMenuID() TrainingMenuID
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type ageDecadeMenu struct {
	decade         Decade
	sortOrder      SortOrder
	trainingMenuID TrainingMenuID
	createdAt      time.Time
	updatedAt      time.Time
}

func (a ageDecadeMenu) Decade() Decade {
	return a.decade
}

func (a ageDecadeMenu) SortOrder() SortOrder {
	return a.sortOrder
}

func (a ageDecadeMenu) TrainingMenuID() TrainingMenuID {
	return a.trainingMenuID
}

func (a ageDecadeMenu) CreatedAt() time.Time {
	return a.createdAt
}

func (a ageDecadeMenu) UpdatedAt() time.Time {
	return a.updatedAt
}

func NewAgeDecadeMenu(
	decade Decade,
	sortOrder SortOrder,
	trainingMenuID TrainingMenuID,
	createdAt time.Time,
	updatedAt time.Time,
) AgeDecadeMenu {
	return &ageDecadeMenu{
		decade:         decade,
		sortOrder:      sortOrder,
		trainingMenuID: trainingMenuID,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}
