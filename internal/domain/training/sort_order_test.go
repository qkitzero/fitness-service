package training

import "testing"

func TestNewSortOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		success   bool
		sortOrder int
		want      SortOrder
	}{
		{"success min sort order", true, 1, SortOrder(1)},
		{"success large sort order", true, 100, SortOrder(100)},
		{"failure zero sort order", false, 0, SortOrder(0)},
		{"failure negative sort order", false, -1, SortOrder(0)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sortOrder, err := NewSortOrder(tt.sortOrder)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && sortOrder != tt.want {
				t.Errorf("NewSortOrder() = %v, want %v", sortOrder, tt.want)
			}

			if tt.success && sortOrder.Int() != tt.sortOrder {
				t.Errorf("Int() = %v, want %v", sortOrder.Int(), tt.sortOrder)
			}
		})
	}
}
