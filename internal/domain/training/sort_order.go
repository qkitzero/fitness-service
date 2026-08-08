package training

const (
	sortOrderMin = 1
	sortOrderMax = 32767
)

type SortOrder int

func (s SortOrder) Int() int {
	return int(s)
}

func NewSortOrder(n int) (SortOrder, error) {
	if n < sortOrderMin || n > sortOrderMax {
		return SortOrder(0), ErrInvalidSortOrder
	}
	return SortOrder(n), nil
}
