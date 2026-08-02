package measurementitem

import "fmt"

type SideAggregation string

const (
	SideAggregationMean SideAggregation = "mean"
	SideAggregationBest SideAggregation = "best"
)

func (s SideAggregation) String() string {
	return string(s)
}

func NewSideAggregation(s string) (SideAggregation, error) {
	switch SideAggregation(s) {
	case SideAggregationMean, SideAggregationBest:
		return SideAggregation(s), nil
	default:
		return SideAggregation(""), fmt.Errorf("invalid side aggregation")
	}
}
