package customer

import "testing"

func TestNewName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		success      bool
		customerName string
	}{
		{"success new name", true, "test customer"},
		{"failure empty name", false, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			name, err := NewName(tt.customerName)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && name.String() != tt.customerName {
				t.Errorf("String() = %v, want %v", name.String(), tt.customerName)
			}
		})
	}
}
