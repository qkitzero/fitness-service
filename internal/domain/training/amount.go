package training

import "fmt"

const (
	amountMin = 1
	amountMax = 999
)

type Amount int

func (a Amount) Int() int {
	return int(a)
}

func NewAmount(n int) (Amount, error) {
	if n < amountMin || n > amountMax {
		return Amount(0), fmt.Errorf("invalid amount")
	}
	return Amount(n), nil
}
