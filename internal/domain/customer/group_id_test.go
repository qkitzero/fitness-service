package customer

import "testing"

func TestNewGroupID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		groupID string
		want    string
	}{
		{"success new group id", true, "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"},
		{"success normalize uppercase", true, "0F4A1A2B-3C4D-5E6F-7A8B-9C0D1E2F3A4B", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"},
		{"failure empty group id", false, "", ""},
		{"failure whitespace group id", false, "   ", ""},
		{"failure not uuid", false, "not-a-uuid", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			groupID, err := NewGroupID(tt.groupID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && groupID.String() != tt.want {
				t.Errorf("String() = %v, want %v", groupID.String(), tt.want)
			}
		})
	}
}
