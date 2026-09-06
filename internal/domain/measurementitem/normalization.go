package measurementitem

import "fmt"

type Normalization string

const (
	NormalizationNone        Normalization = "none"
	NormalizationHeightRatio Normalization = "height_ratio"
)

func (n Normalization) String() string {
	return string(n)
}

func NewNormalization(s string) (Normalization, error) {
	switch Normalization(s) {
	case NormalizationNone, NormalizationHeightRatio:
		return Normalization(s), nil
	default:
		return Normalization(""), fmt.Errorf("invalid normalization")
	}
}
