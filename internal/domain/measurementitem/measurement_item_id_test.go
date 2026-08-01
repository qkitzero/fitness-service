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
