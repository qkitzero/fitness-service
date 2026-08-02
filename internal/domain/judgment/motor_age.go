package judgment

import (
	"github.com/qkitzero/fitness-service/internal/domain/standard"
)

type MotorAge int

func (m MotorAge) Int() int {
	return int(m)
}

func NewMotorAge(ageRange standard.AgeRange) MotorAge {
	return MotorAge(ageRange.Median())
}
