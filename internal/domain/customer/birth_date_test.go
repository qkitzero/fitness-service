package customer

import (
	"database/sql/driver"
	"testing"
	"time"
)

func TestNewBirthDate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		year    int32
		month   int32
		day     int32
	}{
		{"success new birth date", true, 2000, 1, 1},
		{"failure future birth date", false, 2500, 1, 1},
		{"failure invalid birth date", false, 2000, 2, 30},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewBirthDate(tt.year, tt.month, tt.day)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestBirthDateAgeAt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		year  int32
		month int32
		day   int32
		at    time.Time
		want  int
	}{
		{"success on the birthday", 2000, 8, 1, time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC), 26},
		{"success the day before the birthday", 2000, 8, 1, time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC), 25},
		{"success the day after the birthday", 2000, 8, 1, time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC), 26},
		{"success earlier month", 2000, 8, 1, time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), 25},
		{"success later month", 2000, 2, 1, time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC), 26},
		{"success leap day birth on a leap year", 2000, 2, 29, time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC), 24},
		{"success leap day birth on the day before in a common year", 2000, 2, 29, time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC), 25},
		{"success leap day birth on the first of march in a common year", 2000, 2, 29, time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), 26},
		{"success same year", 2026, 1, 1, time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC), 0},
		{"success before birth", 2026, 8, 1, time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC), -1},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			birthDate, err := NewBirthDate(tt.year, tt.month, tt.day)
			if err != nil {
				t.Fatalf("failed to new birth date: %v", err)
			}

			if got := birthDate.AgeAt(tt.at); got != tt.want {
				t.Errorf("AgeAt(%v) = %v, want %v", tt.at, got, tt.want)
			}
		})
	}
}

func TestBirthDateScan(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		value   any
	}{
		{"success scan", true, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"failure invalid type", false, "invalid type"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var birthDate BirthDate
			err := birthDate.Scan(tt.value)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestBirthDateValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		year  int32
		month int32
		day   int32
		want  driver.Value
	}{
		{"success value", 2000, 1, 1, driver.Value(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))},
		{"success leap day", 2000, 2, 29, driver.Value(time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC))},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			birthDate, err := NewBirthDate(tt.year, tt.month, tt.day)
			if err != nil {
				t.Fatalf("failed to new birth date: %v", err)
			}

			value, err := birthDate.Value()
			if err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if value != tt.want {
				t.Errorf("Value() = %v, want %v", value, tt.want)
			}
		})
	}
}
