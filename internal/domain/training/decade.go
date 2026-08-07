package training

import "fmt"

type Decade int

const (
	Decade10 Decade = 10
	Decade20 Decade = 20
	Decade30 Decade = 30
	Decade40 Decade = 40
	Decade50 Decade = 50
	Decade60 Decade = 60
	Decade70 Decade = 70
	Decade80 Decade = 80
	Decade90 Decade = 90
)

func (d Decade) Int() int {
	return int(d)
}

func NewDecade(n int) (Decade, error) {
	switch Decade(n) {
	case Decade10, Decade20, Decade30, Decade40, Decade50, Decade60, Decade70, Decade80, Decade90:
		return Decade(n), nil
	default:
		return Decade(0), fmt.Errorf("invalid decade")
	}
}
