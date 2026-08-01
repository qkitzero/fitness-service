package measurementitem

import "fmt"

type ValueType string

const (
	ValueTypeNumeric ValueType = "numeric"
	ValueTypePaired  ValueType = "paired"
	ValueTypeChoice  ValueType = "choice"
)

func (v ValueType) String() string {
	return string(v)
}

func NewValueType(s string) (ValueType, error) {
	switch ValueType(s) {
	case ValueTypeNumeric, ValueTypePaired, ValueTypeChoice:
		return ValueType(s), nil
	default:
		return ValueType(""), fmt.Errorf("invalid value type")
	}
}
