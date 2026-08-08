package judgment

import (
	"github.com/google/uuid"
)

type PrescribedMenuOverrideID struct {
	uuid.UUID
}

func NewPrescribedMenuOverrideID() PrescribedMenuOverrideID {
	id := uuid.New()
	return PrescribedMenuOverrideID{id}
}
