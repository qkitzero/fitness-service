package customer

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type BirthDate struct {
	time.Time
}

func (b *BirthDate) Scan(value any) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan birth date: %v", value)
	}
	b.Time = t
	return nil
}

func (b BirthDate) Value() (driver.Value, error) {
	return b.Time, nil
}

func (b BirthDate) AgeAt(t time.Time) int {
	age := t.Year() - b.Year()
	if t.Month() < b.Month() || (t.Month() == b.Month() && t.Day() < b.Day()) {
		age--
	}
	return age
}

func NewBirthDate(year, month, day int32) (BirthDate, error) {
	birthDate := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)

	if birthDate.Year() != int(year) || birthDate.Month() != time.Month(month) || birthDate.Day() != int(day) {
		return BirthDate{}, fmt.Errorf("invalid birth date: %d-%02d-%02d", year, month, day)
	}

	if birthDate.After(time.Now()) {
		return BirthDate{}, fmt.Errorf("birth date cannot be in the future")
	}

	return BirthDate{birthDate}, nil
}
