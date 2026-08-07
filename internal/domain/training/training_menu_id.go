package training

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type TrainingMenuID struct {
	uuid.UUID
}

func NewTrainingMenuID() TrainingMenuID {
	id := uuid.New()
	return TrainingMenuID{id}
}

func NewTrainingMenuIDFromString(s string) (TrainingMenuID, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return TrainingMenuID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return TrainingMenuID{id}, nil
}
