package tenant

import (
	"strings"
	"testing"
)

func TestNewTenantID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		success  bool
		tenantID string
		want     string
	}{
		{"success new tenant id", true, "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"},
		{"success trim whitespace", true, "  0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b  ", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"},
		{"success upstream format is not validated", true, "group_00000000", "group_00000000"},
		{"failure empty tenant id", false, "", ""},
		{"failure whitespace tenant id", false, "   ", ""},
		{"failure too long tenant id", false, strings.Repeat("a", 37), ""},
		{"failure control character", false, "0f4a1a2b\x00", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tenantID, err := NewTenantID(tt.tenantID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && tenantID.String() != tt.want {
				t.Errorf("String() = %v, want %v", tenantID.String(), tt.want)
			}
		})
	}
}
