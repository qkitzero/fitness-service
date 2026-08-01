package measurementitem

import "testing"

func TestNewValueType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		success   bool
		valueType string
		want      ValueType
	}{
		{"success numeric", true, "numeric", ValueTypeNumeric},
		{"success paired", true, "paired", ValueTypePaired},
		{"success choice", true, "choice", ValueTypeChoice},
		{"failure empty value type", false, "", ""},
		{"failure invalid value type", false, "unknown", ""},
		{"failure uppercase value type", false, "NUMERIC", ""},
		{"failure untrimmed value type", false, " numeric ", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			valueType, err := NewValueType(tt.valueType)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && valueType != tt.want {
				t.Errorf("NewValueType() = %v, want %v", valueType, tt.want)
			}

			if tt.success && valueType.String() != tt.valueType {
				t.Errorf("String() = %v, want %v", valueType.String(), tt.valueType)
			}
		})
	}
}
