package standard

import (
	"fmt"
	"math"
)

const (
	zScoreMin       = -99.99
	zScoreMax       = 99.99
	zScoreScale     = 100
	zScoreScaleSlop = 1e-9
)

type ZScore float64

func (z ZScore) Float64() float64 {
	return float64(z)
}

func NewZScore(f float64) (ZScore, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return ZScore(0), fmt.Errorf("invalid z score")
	}
	if f < zScoreMin || f > zScoreMax {
		return ZScore(0), fmt.Errorf("invalid z score")
	}
	scaled := f * zScoreScale
	if math.Abs(scaled-math.Round(scaled)) > zScoreScaleSlop {
		return ZScore(0), fmt.Errorf("invalid z score")
	}
	return ZScore(f), nil
}

func MeanZScore(zScores []ZScore) ZScore {
	if len(zScores) == 0 {
		return ZScore(0)
	}

	sum := int64(0)
	for _, zScore := range zScores {
		sum += int64(math.Round(zScore.Float64() * zScoreScale))
	}

	count := int64(len(zScores))
	if sum < 0 {
		return ZScore(float64(-((-2*sum + count) / (2 * count))) / zScoreScale)
	}
	return ZScore(float64((2*sum+count)/(2*count)) / zScoreScale)
}
