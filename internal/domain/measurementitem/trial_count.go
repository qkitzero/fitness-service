package measurementitem

import "fmt"

const trialCountMin = 1

type TrialCount int

func (t TrialCount) Int() int {
	return int(t)
}

func NewTrialCount(n int) (TrialCount, error) {
	if n < trialCountMin {
		return TrialCount(0), fmt.Errorf("invalid trial count")
	}
	return TrialCount(n), nil
}
