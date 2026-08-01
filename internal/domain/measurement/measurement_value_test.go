package measurement

import "testing"

func TestNewMeasurementValue(t *testing.T) {
	t.Parallel()
	value, _ := NewValue(128)
	valueSecondary, _ := NewValue(82)
	valueChoice, _ := NewChoice("片足20cm")

	tests := []struct {
		name           string
		trialIndex     int
		side           Side
		value          *Value
		valueSecondary *Value
		valueChoice    *Choice
	}{
		{"success new numeric value", 1, SideNone, &value, nil, nil},
		{"success new paired value", 1, SideNone, &value, &valueSecondary, nil},
		{"success new choice value", 1, SideLeft, nil, nil, &valueChoice},
		{"success new value without any value", 2, SideRight, nil, nil, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trialIndex, err := NewTrialIndex(tt.trialIndex)
			if err != nil {
				t.Fatalf("failed to new trial index: %v", err)
			}

			m := NewMeasurementValue(trialIndex, tt.side, tt.value, tt.valueSecondary, tt.valueChoice)

			if m.TrialIndex() != trialIndex {
				t.Errorf("TrialIndex() = %v, want %v", m.TrialIndex(), trialIndex)
			}
			if m.Side() != tt.side {
				t.Errorf("Side() = %v, want %v", m.Side(), tt.side)
			}
			if tt.value == nil && m.Value() != nil {
				t.Errorf("Value() = %v, want nil", m.Value())
			}
			if tt.value != nil && (m.Value() == nil || *m.Value() != *tt.value) {
				t.Errorf("Value() = %v, want %v", m.Value(), *tt.value)
			}
			if tt.valueSecondary == nil && m.ValueSecondary() != nil {
				t.Errorf("ValueSecondary() = %v, want nil", m.ValueSecondary())
			}
			if tt.valueSecondary != nil && (m.ValueSecondary() == nil || *m.ValueSecondary() != *tt.valueSecondary) {
				t.Errorf("ValueSecondary() = %v, want %v", m.ValueSecondary(), *tt.valueSecondary)
			}
			if tt.valueChoice == nil && m.ValueChoice() != nil {
				t.Errorf("ValueChoice() = %v, want nil", m.ValueChoice())
			}
			if tt.valueChoice != nil && (m.ValueChoice() == nil || *m.ValueChoice() != *tt.valueChoice) {
				t.Errorf("ValueChoice() = %v, want %v", m.ValueChoice(), *tt.valueChoice)
			}

			if got := m.Value(); got != nil {
				*got = Value(0)
				if *m.Value() != *tt.value {
					t.Errorf("Value() = %v, want %v", m.Value(), *tt.value)
				}
			}
			if got := m.ValueSecondary(); got != nil {
				*got = Value(0)
				if *m.ValueSecondary() != *tt.valueSecondary {
					t.Errorf("ValueSecondary() = %v, want %v", m.ValueSecondary(), *tt.valueSecondary)
				}
			}
			if got := m.ValueChoice(); got != nil {
				*got = Choice("")
				if *m.ValueChoice() != *tt.valueChoice {
					t.Errorf("ValueChoice() = %v, want %v", m.ValueChoice(), *tt.valueChoice)
				}
			}
		})
	}
}
