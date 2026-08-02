package standard

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type AgeGroupStandardID struct {
	uuid.UUID
}

func NewAgeGroupStandardID() AgeGroupStandardID {
	id := uuid.New()
	return AgeGroupStandardID{id}
}

func NewAgeGroupStandardIDFromString(s string) (AgeGroupStandardID, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return AgeGroupStandardID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return AgeGroupStandardID{id}, nil
}
