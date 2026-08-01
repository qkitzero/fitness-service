package measurement

const (
	ageAtMeasurementMin = 0
	ageAtMeasurementMax = 150
)

type AgeAtMeasurement int

func (a AgeAtMeasurement) Int() int {
	return int(a)
}

func NewAgeAtMeasurement(n int) (AgeAtMeasurement, error) {
	if n < ageAtMeasurementMin || n > ageAtMeasurementMax {
		return AgeAtMeasurement(0), ErrInvalidAgeAtMeasurement
	}
	return AgeAtMeasurement(n), nil
}
