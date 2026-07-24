package customer

import "fmt"

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

func (g Gender) String() string {
	return string(g)
}

func NewGender(s string) (Gender, error) {
	switch Gender(s) {
	case GenderMale, GenderFemale, GenderOther:
		return Gender(s), nil
	default:
		return Gender(""), fmt.Errorf("invalid gender")
	}
}
