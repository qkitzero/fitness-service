package measurement

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type MeasuredOn struct {
	time.Time
}

func (m *MeasuredOn) Scan(value any) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan measured on: %v", value)
	}
	m.Time = t
	return nil
}

func (m MeasuredOn) Value() (driver.Value, error) {
	return m.Time, nil
}

func NewMeasuredOn(year, month, day int32) (MeasuredOn, error) {
	measuredOn := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)

	if measuredOn.Year() != int(year) || measuredOn.Month() != time.Month(month) || measuredOn.Day() != int(day) {
		return MeasuredOn{}, fmt.Errorf("invalid measured on: %d-%02d-%02d", year, month, day)
	}

	now := time.Now().UTC()
	latest := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)

	if measuredOn.After(latest) {
		return MeasuredOn{}, fmt.Errorf("measured on cannot be in the future")
	}

	return MeasuredOn{measuredOn}, nil
}
