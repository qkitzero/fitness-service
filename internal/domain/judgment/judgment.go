package judgment

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurement"
)

type Judgment interface {
	MeasurementID() measurement.MeasurementID
	Advice() *Advice
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(advice *Advice)
}

type judgment struct {
	measurementID measurement.MeasurementID
	advice        *Advice
	createdAt     time.Time
	updatedAt     time.Time
}

func (j judgment) MeasurementID() measurement.MeasurementID {
	return j.measurementID
}

func (j judgment) Advice() *Advice {
	if j.advice == nil {
		return nil
	}
	a := *j.advice
	return &a
}

func (j judgment) CreatedAt() time.Time {
	return j.createdAt
}

func (j judgment) UpdatedAt() time.Time {
	return j.updatedAt
}

func (j *judgment) Update(advice *Advice) {
	j.advice = nil
	if advice != nil {
		a := *advice
		j.advice = &a
	}
	j.updatedAt = time.Now().UTC()
}

func NewJudgment(
	measurementID measurement.MeasurementID,
	advice *Advice,
	createdAt time.Time,
	updatedAt time.Time,
) Judgment {
	return newJudgment(measurementID, advice, createdAt, updatedAt)
}

func ReconstructJudgment(
	measurementID measurement.MeasurementID,
	advice *Advice,
	createdAt time.Time,
	updatedAt time.Time,
) Judgment {
	return newJudgment(measurementID, advice, createdAt, updatedAt)
}

func newJudgment(
	measurementID measurement.MeasurementID,
	advice *Advice,
	createdAt time.Time,
	updatedAt time.Time,
) Judgment {
	j := &judgment{
		measurementID: measurementID,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
	if advice != nil {
		a := *advice
		j.advice = &a
	}
	return j
}
