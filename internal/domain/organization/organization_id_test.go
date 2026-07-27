package organization

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewOrganizationID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
	}{
		{"success new organization id", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewOrganizationID()
			if tt.success && id.UUID == uuid.Nil {
				t.Errorf("expected valid organization id, but got a nil UUID")
			}
		})
	}
}

func TestNewOrganizationIDFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		id      string
		want    string
	}{
		{"success new organization id from string", true, "3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f", "3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f"},
		{"success trim space", true, "  3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f  ", "3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f"},
		{"success normalize uppercase", true, "3D1E6A5C-7B8F-4C2D-9A0E-1F2B3C4D5E6F", "3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f"},
		{"failure empty organization id", false, "", ""},
		{"failure whitespace organization id", false, "   ", ""},
		{"failure invalid organization id", false, "0123456789", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := NewOrganizationIDFromString(tt.id)
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
