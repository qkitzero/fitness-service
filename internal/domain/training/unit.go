package training

import "fmt"

type Unit string

const (
	UnitReps    Unit = "reps"
	UnitSeconds Unit = "seconds"
	UnitMinutes Unit = "minutes"
)

func (u Unit) String() string {
	return string(u)
}

func NewUnit(s string) (Unit, error) {
	switch Unit(s) {
	case UnitReps, UnitSeconds, UnitMinutes:
		return Unit(s), nil
	default:
		return Unit(""), fmt.Errorf("invalid unit")
	}
}
