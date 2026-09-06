package judgment

type MotorAge int

func (m MotorAge) Int() int {
	return int(m)
}

func NewMotorAge(age int) MotorAge {
	return MotorAge(age)
}
