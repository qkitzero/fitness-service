package training

import (
	"fmt"
	"math"
)

type Part string

const (
	PartUpperLimb Part = "upper_limb"
	PartLowerLimb Part = "lower_limb"
	PartWholeBody Part = "whole_body"
)

var orderedParts = []Part{
	PartUpperLimb,
	PartLowerLimb,
	PartWholeBody,
}

func (p Part) String() string {
	return string(p)
}

func (p Part) Order() int {
	for i, part := range orderedParts {
		if p == part {
			return i + 1
		}
	}
	return math.MaxInt
}

func NewPart(s string) (Part, error) {
	for _, part := range orderedParts {
		if Part(s) == part {
			return part, nil
		}
	}
	return Part(""), fmt.Errorf("invalid part")
}
