package judgment

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewPrescribedMenuOverrideID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
	}{
		{"success new prescribed menu override id", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewPrescribedMenuOverrideID()
			if tt.success && id.UUID == uuid.Nil {
				t.Errorf("expected valid prescribed menu override id, but got a nil UUID")
			}
		})
	}
}
