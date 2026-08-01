package measurement

import (
	"database/sql/driver"
	"testing"
	"time"
)

func TestNewMeasuredOn(t *testing.T) {
	t.Parallel()
	tomorrow := time.Now().UTC().AddDate(0, 0, 1)
	dayAfterTomorrow := time.Now().UTC().AddDate(0, 0, 2)
	tests := []struct {
		name    string
		success bool
		year    int32
		month   int32
		day     int32
	}{
		{"success new measured on", true, 2026, 1, 1},
		{"success leap day", true, 2024, 2, 29},
		{"success tomorrow in utc covers today in any time zone", true, int32(tomorrow.Year()), int32(tomorrow.Month()), int32(tomorrow.Day())},
		{"failure the day after tomorrow in utc", false, int32(dayAfterTomorrow.Year()), int32(dayAfterTomorrow.Month()), int32(dayAfterTomorrow.Day())},
		{"failure future measured on", false, 2500, 1, 1},
		{"failure invalid measured on", false, 2026, 2, 30},
		{"failure zero measured on", false, 0, 0, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewMeasuredOn(tt.year, tt.month, tt.day)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestMeasuredOnScan(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		value   any
	}{
		{"success scan", true, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"failure invalid type", false, "invalid type"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var measuredOn MeasuredOn
			err := measuredOn.Scan(tt.value)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestMeasuredOnValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		year  int32
		month int32
		day   int32
		want  driver.Value
	}{
		{"success value", 2026, 1, 1, driver.Value(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))},
		{"success leap day", 2024, 2, 29, driver.Value(time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC))},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			measuredOn, err := NewMeasuredOn(tt.year, tt.month, tt.day)
			if err != nil {
				t.Fatalf("failed to new measured on: %v", err)
			}

			value, err := measuredOn.Value()
			if err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if value != tt.want {
				t.Errorf("Value() = %v, want %v", value, tt.want)
			}
		})
	}
}
