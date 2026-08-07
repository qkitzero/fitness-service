package training

import "fmt"

const sortOrderMin = 1

type SortOrder int

func (s SortOrder) Int() int {
	return int(s)
}

func NewSortOrder(n int) (SortOrder, error) {
	if n < sortOrderMin {
		return SortOrder(0), fmt.Errorf("invalid sort order")
	}
	return SortOrder(n), nil
}
