package training

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewTrainingMenuID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
	}{
		{"success new training menu id", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewTrainingMenuID()
			if tt.success && id.UUID == uuid.Nil {
				t.Errorf("expected valid training menu id, but got a nil UUID")
			}
		})
	}
}

func TestNewTrainingMenuIDFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		id      string
		want    string
	}{
		{"success new training menu id from string", true, "8e35a3d6-3d29-4c22-a5cd-1c8bbbf5f36c", "8e35a3d6-3d29-4c22-a5cd-1c8bbbf5f36c"},
		{"success trim space", true, "  8e35a3d6-3d29-4c22-a5cd-1c8bbbf5f36c  ", "8e35a3d6-3d29-4c22-a5cd-1c8bbbf5f36c"},
		{"success normalize uppercase", true, "8E35A3D6-3D29-4C22-A5CD-1C8BBBF5F36C", "8e35a3d6-3d29-4c22-a5cd-1c8bbbf5f36c"},
		{"failure empty training menu id", false, "", ""},
		{"failure whitespace training menu id", false, "   ", ""},
		{"failure invalid training menu id", false, "0123456789", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := NewTrainingMenuIDFromString(tt.id)
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
