package standard

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

type AgeGroupStandard interface {
	ID() AgeGroupStandardID
	MeasurementItemID() measurementitem.MeasurementItemID
	Gender() Gender
	AgeRange() AgeRange
	Mean() Mean
	StandardDeviation() StandardDeviation
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type ageGroupStandard struct {
	id                AgeGroupStandardID
	measurementItemID measurementitem.MeasurementItemID
	gender            Gender
	ageRange          AgeRange
	mean              Mean
	standardDeviation StandardDeviation
	createdAt         time.Time
	updatedAt         time.Time
}

func (a ageGroupStandard) ID() AgeGroupStandardID {
	return a.id
}

func (a ageGroupStandard) MeasurementItemID() measurementitem.MeasurementItemID {
	return a.measurementItemID
}

func (a ageGroupStandard) Gender() Gender {
	return a.gender
}

func (a ageGroupStandard) AgeRange() AgeRange {
	return a.ageRange
}

func (a ageGroupStandard) Mean() Mean {
	return a.mean
}

func (a ageGroupStandard) StandardDeviation() StandardDeviation {
	return a.standardDeviation
}

func (a ageGroupStandard) CreatedAt() time.Time {
	return a.createdAt
}

func (a ageGroupStandard) UpdatedAt() time.Time {
	return a.updatedAt
}

func NewAgeGroupStandard(
	id AgeGroupStandardID,
	measurementItemID measurementitem.MeasurementItemID,
	gender Gender,
	ageRange AgeRange,
	mean Mean,
	standardDeviation StandardDeviation,
	createdAt time.Time,
	updatedAt time.Time,
) AgeGroupStandard {
	return &ageGroupStandard{
		id:                id,
		measurementItemID: measurementItemID,
		gender:            gender,
		ageRange:          ageRange,
		mean:              mean,
		standardDeviation: standardDeviation,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}
