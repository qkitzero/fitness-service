package judgment

import "fmt"

const (
	motorAgeMin = 0
	motorAgeMax = 150
)

type MotorAge int

func (m MotorAge) Int() int {
	return int(m)
}

func NewMotorAge(n int) (MotorAge, error) {
	if n < motorAgeMin || n > motorAgeMax {
		return MotorAge(0), fmt.Errorf("invalid motor age")
	}
	return MotorAge(n), nil
}
