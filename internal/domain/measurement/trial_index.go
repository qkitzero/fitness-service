package measurement

import "fmt"

const trialIndexMin = 1

type TrialIndex int

func (t TrialIndex) Int() int {
	return int(t)
}

func NewTrialIndex(n int) (TrialIndex, error) {
	if n < trialIndexMin {
		return TrialIndex(0), fmt.Errorf("invalid trial index")
	}
	return TrialIndex(n), nil
}
