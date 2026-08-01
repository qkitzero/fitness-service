package measurement

type MeasurementValue interface {
	TrialIndex() TrialIndex
	Side() Side
	Value() *Value
	ValueSecondary() *Value
	ValueChoice() *Choice
}

type measurementValue struct {
	trialIndex     TrialIndex
	side           Side
	value          *Value
	valueSecondary *Value
	valueChoice    *Choice
}

func (m measurementValue) TrialIndex() TrialIndex {
	return m.trialIndex
}

func (m measurementValue) Side() Side {
	return m.side
}

func (m measurementValue) Value() *Value {
	if m.value == nil {
		return nil
	}
	v := *m.value
	return &v
}

func (m measurementValue) ValueSecondary() *Value {
	if m.valueSecondary == nil {
		return nil
	}
	v := *m.valueSecondary
	return &v
}

func (m measurementValue) ValueChoice() *Choice {
	if m.valueChoice == nil {
		return nil
	}
	c := *m.valueChoice
	return &c
}

func NewMeasurementValue(
	trialIndex TrialIndex,
	side Side,
	value *Value,
	valueSecondary *Value,
	valueChoice *Choice,
) MeasurementValue {
	m := &measurementValue{
		trialIndex: trialIndex,
		side:       side,
	}
	if value != nil {
		v := *value
		m.value = &v
	}
	if valueSecondary != nil {
		v := *valueSecondary
		m.valueSecondary = &v
	}
	if valueChoice != nil {
		c := *valueChoice
		m.valueChoice = &c
	}
	return m
}
