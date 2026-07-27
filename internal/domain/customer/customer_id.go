package customer

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type CustomerID struct {
	uuid.UUID
}

func NewCustomerID() CustomerID {
	id := uuid.New()
	return CustomerID{id}
}

func NewCustomerIDFromString(s string) (CustomerID, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return CustomerID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return CustomerID{id}, nil
}
