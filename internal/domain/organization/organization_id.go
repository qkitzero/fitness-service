package organization

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type OrganizationID struct {
	uuid.UUID
}

func NewOrganizationID() OrganizationID {
	id := uuid.New()
	return OrganizationID{id}
}

func NewOrganizationIDFromString(s string) (OrganizationID, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return OrganizationID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return OrganizationID{id}, nil
}
