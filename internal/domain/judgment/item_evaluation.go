package judgment

import (
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/standard"
)

type ItemEvaluation interface {
	MeasurementItemID() measurementitem.MeasurementItemID
	Value() measurement.Value
	Mean() standard.Mean
	ZScore() standard.ZScore
	Rank() standard.Rank
	StandardDeviation() standard.StandardDeviation
	AgeRange() standard.AgeRange
}

type itemEvaluation struct {
	measurementItemID measurementitem.MeasurementItemID
	value             measurement.Value
	mean              standard.Mean
	zScore            standard.ZScore
	rank              standard.Rank
	standardDeviation standard.StandardDeviation
	ageRange          standard.AgeRange
}

func (i itemEvaluation) MeasurementItemID() measurementitem.MeasurementItemID {
	return i.measurementItemID
}

func (i itemEvaluation) Value() measurement.Value {
	return i.value
}

func (i itemEvaluation) Mean() standard.Mean {
	return i.mean
}

func (i itemEvaluation) ZScore() standard.ZScore {
	return i.zScore
}

func (i itemEvaluation) Rank() standard.Rank {
	return i.rank
}

func (i itemEvaluation) StandardDeviation() standard.StandardDeviation {
	return i.standardDeviation
}

func (i itemEvaluation) AgeRange() standard.AgeRange {
	return i.ageRange
}

func newItemEvaluation(
	measurementItemID measurementitem.MeasurementItemID,
	value measurement.Value,
	mean standard.Mean,
	zScore standard.ZScore,
	rank standard.Rank,
	standardDeviation standard.StandardDeviation,
	ageRange standard.AgeRange,
) ItemEvaluation {
	return &itemEvaluation{
		measurementItemID: measurementItemID,
		value:             value,
		mean:              mean,
		zScore:            zScore,
		rank:              rank,
		standardDeviation: standardDeviation,
		ageRange:          ageRange,
	}
}
