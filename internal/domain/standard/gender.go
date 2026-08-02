package standard

import "fmt"

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

func (g Gender) String() string {
	return string(g)
}

func NewGender(s string) (Gender, error) {
	switch Gender(s) {
	case GenderMale, GenderFemale:
		return Gender(s), nil
	default:
		return Gender(""), fmt.Errorf("invalid gender")
	}
}
