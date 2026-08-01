package measurement

import (
	"fmt"
	"math"
)

const (
	valueMin       = 0
	valueMax       = 9999.99
	valueScale     = 100
	valueScaleSlop = 1e-9
)

type Value float64

func (v Value) Float64() float64 {
	return float64(v)
}

func NewValue(f float64) (Value, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return Value(0), fmt.Errorf("invalid value")
	}
	if f < valueMin || f > valueMax {
		return Value(0), fmt.Errorf("invalid value")
	}
	scaled := f * valueScale
	if math.Abs(scaled-math.Round(scaled)) > valueScaleSlop {
		return Value(0), fmt.Errorf("invalid value")
	}
	return Value(f), nil
}
