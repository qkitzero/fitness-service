package customer

import "testing"

func TestNewGender(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		gender  string
		want    Gender
	}{
		{"success male", true, "male", GenderMale},
		{"success female", true, "female", GenderFemale},
		{"success other", true, "other", GenderOther},
		{"failure empty gender", false, "", ""},
		{"failure invalid gender", false, "unknown", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gender, err := NewGender(tt.gender)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && gender != tt.want {
				t.Errorf("NewGender() = %v, want %v", gender, tt.want)
			}

			if tt.success && gender.String() != tt.gender {
				t.Errorf("String() = %v, want %v", gender.String(), tt.gender)
			}
		})
	}
}
