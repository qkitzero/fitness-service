package customer

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type GroupID string

func (g GroupID) String() string {
	return string(g)
}

func NewGroupID(s string) (GroupID, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return GroupID(""), fmt.Errorf("invalid group id: %w", err)
	}
	return GroupID(id.String()), nil
}
