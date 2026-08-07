package training

import "fmt"

const (
	setsMin = 1
	setsMax = 99
)

type Sets int

func (s Sets) Int() int {
	return int(s)
}

func NewSets(n int) (Sets, error) {
	if n < setsMin || n > setsMax {
		return Sets(0), fmt.Errorf("invalid sets")
	}
	return Sets(n), nil
}
