package measurement

import (
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

type MeasurementEntry interface {
	MeasurementItemID() measurementitem.MeasurementItemID
	Unmeasurable() bool
	Note() *Note
	Values() []MeasurementValue
	ExpectedValueCount() (int, bool)
}

type measurementEntry struct {
	measurementItemID       measurementitem.MeasurementItemID
	unmeasurable            bool
	note                    *Note
	values                  []MeasurementValue
	expectedValueCount      int
	expectedValueCountKnown bool
}

func (m measurementEntry) MeasurementItemID() measurementitem.MeasurementItemID {
	return m.measurementItemID
}

func (m measurementEntry) Unmeasurable() bool {
	return m.unmeasurable
}

func (m measurementEntry) Note() *Note {
	if m.note == nil {
		return nil
	}
	n := *m.note
	return &n
}

func (m measurementEntry) Values() []MeasurementValue {
	values := make([]MeasurementValue, len(m.values))
	copy(values, m.values)
	return values
}

func (m measurementEntry) ExpectedValueCount() (int, bool) {
	return m.expectedValueCount, m.expectedValueCountKnown
}

type valueKey struct {
	trialIndex TrialIndex
	side       Side
}

func verifySide(item measurementitem.MeasurementItem, side Side) error {
	switch item.SideMode() {
	case measurementitem.SideModeNone:
		if side != SideNone {
			return ErrInvalidSide
		}
	case measurementitem.SideModeBilateral:
		if side != SideLeft && side != SideRight {
			return ErrInvalidSide
		}
	case measurementitem.SideModeOptionalBilateral:
		if side != SideNone && side != SideLeft && side != SideRight {
			return ErrInvalidSide
		}
	default:
		return ErrInvalidSide
	}
	return nil
}

func verifyValueForType(item measurementitem.MeasurementItem, value MeasurementValue) error {
	switch item.ValueType() {
	case measurementitem.ValueTypeNumeric:
		if value.Value() == nil || value.ValueSecondary() != nil || value.ValueChoice() != nil {
			return ErrInvalidValueForType
		}
	case measurementitem.ValueTypePaired:
		if value.Value() == nil || value.ValueSecondary() == nil || value.ValueChoice() != nil {
			return ErrInvalidValueForType
		}
	case measurementitem.ValueTypeChoice:
		if value.ValueChoice() == nil || value.Value() != nil || value.ValueSecondary() != nil {
			return ErrInvalidValueForType
		}
	default:
		return ErrInvalidValueForType
	}
	return nil
}

func NewMeasurementEntry(
	item measurementitem.MeasurementItem,
	unmeasurable bool,
	note *Note,
	values []MeasurementValue,
) (MeasurementEntry, error) {
	if unmeasurable && len(values) > 0 {
		return nil, ErrInvalidValueCount
	}

	seen := make(map[valueKey]struct{}, len(values))
	for _, value := range values {
		if err := verifySide(item, value.Side()); err != nil {
			return nil, err
		}
		if err := verifyValueForType(item, value); err != nil {
			return nil, err
		}
		if value.TrialIndex().Int() > item.TrialCount().Int() {
			return nil, ErrInvalidValueCount
		}
		key := valueKey{trialIndex: value.TrialIndex(), side: value.Side()}
		if _, ok := seen[key]; ok {
			return nil, ErrDuplicateValue
		}
		seen[key] = struct{}{}
	}

	expectedValueCount := item.TrialCount().Int()
	if item.SideMode() != measurementitem.SideModeNone {
		expectedValueCount *= 2
	}

	return newMeasurementEntry(item.ID(), unmeasurable, note, values, expectedValueCount, true), nil
}

func ReconstructMeasurementEntry(
	measurementItemID measurementitem.MeasurementItemID,
	unmeasurable bool,
	note *Note,
	values []MeasurementValue,
) MeasurementEntry {
	return newMeasurementEntry(measurementItemID, unmeasurable, note, values, 0, false)
}

func newMeasurementEntry(
	measurementItemID measurementitem.MeasurementItemID,
	unmeasurable bool,
	note *Note,
	values []MeasurementValue,
	expectedValueCount int,
	expectedValueCountKnown bool,
) MeasurementEntry {
	m := &measurementEntry{
		measurementItemID:       measurementItemID,
		unmeasurable:            unmeasurable,
		values:                  make([]MeasurementValue, len(values)),
		expectedValueCount:      expectedValueCount,
		expectedValueCountKnown: expectedValueCountKnown,
	}
	copy(m.values, values)
	if note != nil {
		n := *note
		m.note = &n
	}
	return m
}
