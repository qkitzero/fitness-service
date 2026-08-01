package measurementitem

import (
	"fmt"
	"math"
)

type Category string

const (
	CategoryVital           Category = "vital"
	CategoryPhysique        Category = "physique"
	CategoryBodyComposition Category = "body_composition"
	CategoryMotorFunction   Category = "motor_function"
)

var orderedCategories = []Category{
	CategoryVital,
	CategoryPhysique,
	CategoryBodyComposition,
	CategoryMotorFunction,
}

func (c Category) String() string {
	return string(c)
}

func (c Category) Order() int {
	for i, category := range orderedCategories {
		if c == category {
			return i + 1
		}
	}
	return math.MaxInt
}

func NewCategory(s string) (Category, error) {
	for _, category := range orderedCategories {
		if Category(s) == category {
			return category, nil
		}
	}
	return Category(""), fmt.Errorf("invalid category")
}
