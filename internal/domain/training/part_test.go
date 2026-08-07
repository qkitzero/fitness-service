package training

import (
	"math"
	"testing"
)

func TestNewPart(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		part    string
		want    Part
	}{
		{"success upper limb", true, "upper_limb", PartUpperLimb},
		{"success lower limb", true, "lower_limb", PartLowerLimb},
		{"success whole body", true, "whole_body", PartWholeBody},
		{"failure empty part", false, "", ""},
		{"failure invalid part", false, "unknown", ""},
		{"failure uppercase part", false, "UPPER_LIMB", ""},
		{"failure untrimmed part", false, " upper_limb ", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			part, err := NewPart(tt.part)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && part != tt.want {
				t.Errorf("NewPart() = %v, want %v", part, tt.want)
			}

			if tt.success && part.String() != tt.part {
				t.Errorf("String() = %v, want %v", part.String(), tt.part)
			}
		})
	}
}

func TestPartOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		part Part
		want int
	}{
		{"upper limb comes first", PartUpperLimb, 1},
		{"lower limb comes second", PartLowerLimb, 2},
		{"whole body comes third", PartWholeBody, 3},
		{"unknown part comes last", Part("unknown"), math.MaxInt},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.part.Order(); got != tt.want {
				t.Errorf("Order() = %v, want %v", got, tt.want)
			}
		})
	}
}
