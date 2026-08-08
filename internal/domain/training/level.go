package training

import "fmt"

const (
	levelMin = 1
	levelMax = 5
)

type Level int

func (l Level) Int() int {
	return int(l)
}

func NewLevel(n int) (Level, error) {
	if n < levelMin || n > levelMax {
		return Level(0), fmt.Errorf("invalid level")
	}
	return Level(n), nil
}
