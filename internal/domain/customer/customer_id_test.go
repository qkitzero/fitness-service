package customer

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewCustomerID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
	}{
		{"success new customer id", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewCustomerID()
			if tt.success && id.UUID == uuid.Nil {
				t.Errorf("expected valid customer id, but got a nil UUID")
			}
		})
	}
}

func TestNewCustomerIDFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		id      string
		want    string
	}{
		{"success new customer id from string", true, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"},
		{"success trim space", true, "  fe8c2263-bbac-4bb9-a41d-b04f5afc4425  ", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"},
		{"success normalize uppercase", true, "FE8C2263-BBAC-4BB9-A41D-B04F5AFC4425", "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"},
		{"failure empty customer id", false, "", ""},
		{"failure whitespace customer id", false, "   ", ""},
		{"failure invalid customer id", false, "0123456789", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := NewCustomerIDFromString(tt.id)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && id.String() != tt.want {
				t.Errorf("String() = %v, want %v", id.String(), tt.want)
			}
		})
	}
}
