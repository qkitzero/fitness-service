package training

import "testing"

func TestNewAmount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		amount  int
		want    Amount
	}{
		{"success min amount", true, 1, Amount(1)},
		{"success typical amount", true, 10, Amount(10)},
		{"success max amount", true, 999, Amount(999)},
		{"failure zero amount", false, 0, Amount(0)},
		{"failure negative amount", false, -1, Amount(0)},
		{"failure amount above the max", false, 1000, Amount(0)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			amount, err := NewAmount(tt.amount)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && amount != tt.want {
				t.Errorf("NewAmount() = %v, want %v", amount, tt.want)
			}

			if tt.success && amount.Int() != tt.amount {
				t.Errorf("Int() = %v, want %v", amount.Int(), tt.amount)
			}
		})
	}
}
