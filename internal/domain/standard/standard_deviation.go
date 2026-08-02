package standard

import (
	"fmt"
	"math"
)

const (
	standardDeviationMax       = 9999.99
	standardDeviationScale     = 100
	standardDeviationScaleSlop = 1e-9
)

type StandardDeviation float64

func (s StandardDeviation) Float64() float64 {
	return float64(s)
}

func NewStandardDeviation(f float64) (StandardDeviation, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return StandardDeviation(0), fmt.Errorf("invalid standard deviation")
	}
	if f <= 0 || f > standardDeviationMax {
		return StandardDeviation(0), fmt.Errorf("invalid standard deviation")
	}
	scaled := f * standardDeviationScale
	if math.Abs(scaled-math.Round(scaled)) > standardDeviationScaleSlop {
		return StandardDeviation(0), fmt.Errorf("invalid standard deviation")
	}
	return StandardDeviation(f), nil
}
