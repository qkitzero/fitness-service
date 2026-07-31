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
