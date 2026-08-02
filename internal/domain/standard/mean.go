package standard

import (
	"fmt"
	"math"
)

const (
	meanMin       = 0
	meanMax       = 9999.99
	meanScale     = 100
	meanScaleSlop = 1e-9
)

type Mean float64

func (m Mean) Float64() float64 {
	return float64(m)
}

func NewMean(f float64) (Mean, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return Mean(0), fmt.Errorf("invalid mean")
	}
	if f < meanMin || f > meanMax {
		return Mean(0), fmt.Errorf("invalid mean")
	}
	scaled := f * meanScale
	if math.Abs(scaled-math.Round(scaled)) > meanScaleSlop {
		return Mean(0), fmt.Errorf("invalid mean")
	}
	return Mean(f), nil
}
