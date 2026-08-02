package standard

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewAgeGroupStandardID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
	}{
		{"success new age group standard id", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewAgeGroupStandardID()
			if tt.success && id.UUID == uuid.Nil {
				t.Errorf("expected valid age group standard id, but got a nil UUID")
			}
		})
	}
}

func TestNewAgeGroupStandardIDFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		id      string
		want    string
	}{
		{"success new age group standard id from string", true, "a2f5f8bb-8a9d-4c3b-9f0e-0a1cb1f1b1d0", "a2f5f8bb-8a9d-4c3b-9f0e-0a1cb1f1b1d0"},
		{"success trim space", true, "  a2f5f8bb-8a9d-4c3b-9f0e-0a1cb1f1b1d0  ", "a2f5f8bb-8a9d-4c3b-9f0e-0a1cb1f1b1d0"},
		{"success normalize uppercase", true, "A2F5F8BB-8A9D-4C3B-9F0E-0A1CB1F1B1D0", "a2f5f8bb-8a9d-4c3b-9f0e-0a1cb1f1b1d0"},
		{"failure empty age group standard id", false, "", ""},
		{"failure whitespace age group standard id", false, "   ", ""},
		{"failure invalid age group standard id", false, "0123456789", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := NewAgeGroupStandardIDFromString(tt.id)
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
