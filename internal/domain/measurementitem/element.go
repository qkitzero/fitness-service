package measurementitem

import (
	"fmt"
	"math"
)

type Element string

const (
	ElementMuscleStrength  Element = "muscle_strength"
	ElementMuscleEndurance Element = "muscle_endurance"
	ElementFlexibility     Element = "flexibility"
	ElementAgility         Element = "agility"
	ElementBalance         Element = "balance"
	ElementMobility        Element = "mobility"
)

var orderedElements = []Element{
	ElementMuscleStrength,
	ElementMuscleEndurance,
	ElementFlexibility,
	ElementAgility,
	ElementBalance,
	ElementMobility,
}

func (e Element) String() string {
	return string(e)
}

func (e Element) Order() int {
	for i, element := range orderedElements {
		if e == element {
			return i + 1
		}
	}
	return math.MaxInt
}

func NewElement(s string) (Element, error) {
	for _, element := range orderedElements {
		if Element(s) == element {
			return element, nil
		}
	}
	return Element(""), fmt.Errorf("invalid element")
}
