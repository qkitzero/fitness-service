package measurementitem

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewMeasurementItemID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
	}{
		{"success new measurement item id", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewMeasurementItemID()
			if tt.success && id.UUID == uuid.Nil {
				t.Errorf("expected valid measurement item id, but got a nil UUID")
			}
		})
	}
}

func TestNewMeasurementItemIDFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		id      string
		want    string
	}{
		{"success new measurement item id from string", true, "45c2f5cd-ae75-4b2e-8302-69051f0343d5", "45c2f5cd-ae75-4b2e-8302-69051f0343d5"},
		{"success trim space", true, "  45c2f5cd-ae75-4b2e-8302-69051f0343d5  ", "45c2f5cd-ae75-4b2e-8302-69051f0343d5"},
		{"success normalize uppercase", true, "45C2F5CD-AE75-4B2E-8302-69051F0343D5", "45c2f5cd-ae75-4b2e-8302-69051f0343d5"},
		{"failure empty measurement item id", false, "", ""},
		{"failure whitespace measurement item id", false, "   ", ""},
		{"failure invalid measurement item id", false, "0123456789", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := NewMeasurementItemIDFromString(tt.id)
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
