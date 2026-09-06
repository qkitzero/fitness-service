package measurementitem

import "fmt"

type TrialAggregation string

const (
	TrialAggregationMean TrialAggregation = "mean"
	TrialAggregationBest TrialAggregation = "best"
)

func (t TrialAggregation) String() string {
	return string(t)
}

func NewTrialAggregation(s string) (TrialAggregation, error) {
	switch TrialAggregation(s) {
	case TrialAggregationMean, TrialAggregationBest:
		return TrialAggregation(s), nil
	default:
		return TrialAggregation(""), fmt.Errorf("invalid trial aggregation")
	}
}
