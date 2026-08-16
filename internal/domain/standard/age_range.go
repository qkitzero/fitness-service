package standard

import "fmt"

const (
	ageRangeMin = 0
	ageRangeMax = 150
)

type AgeRange struct {
	from int
	to   int
}

func (a AgeRange) From() int {
	return a.from
}

func (a AgeRange) To() int {
	return a.to
}

func (a AgeRange) Contains(age int) bool {
	return a.Distance(age) == 0
}

func (a AgeRange) Distance(age int) int {
	if age < a.from {
		return a.from - age
	}
	if age > a.to {
		return age - a.to
	}
	return 0
}

func (a AgeRange) Median() int {
	if a.to == ageRangeMax {
		return a.from
	}
	return a.from + (a.to-a.from+1)/2
}

func NewAgeRange(from, to int) (AgeRange, error) {
	if from < ageRangeMin || from > ageRangeMax {
		return AgeRange{}, fmt.Errorf("invalid age range")
	}
	if to < ageRangeMin || to > ageRangeMax {
		return AgeRange{}, fmt.Errorf("invalid age range")
	}
	if from > to {
		return AgeRange{}, fmt.Errorf("invalid age range")
	}
	return AgeRange{from: from, to: to}, nil
}
