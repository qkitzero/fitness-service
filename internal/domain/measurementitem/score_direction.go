package measurementitem

import "fmt"

type ScoreDirection string

const (
	ScoreDirectionHigherIsBetter ScoreDirection = "higher_is_better"
	ScoreDirectionLowerIsBetter  ScoreDirection = "lower_is_better"
)

func (s ScoreDirection) String() string {
	return string(s)
}

func NewScoreDirection(s string) (ScoreDirection, error) {
	switch ScoreDirection(s) {
	case ScoreDirectionHigherIsBetter, ScoreDirectionLowerIsBetter:
		return ScoreDirection(s), nil
	default:
		return ScoreDirection(""), fmt.Errorf("invalid score direction")
	}
}
