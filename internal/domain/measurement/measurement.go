package measurement

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/staff"
)

type Measurement interface {
	ID() MeasurementID
	CustomerID() customer.CustomerID
	MeasuredOn() MeasuredOn
	MeasuredBy() staff.StaffID
	AgeAtMeasurement() AgeAtMeasurement
	UpdatedBy() staff.StaffID
	IsDraft() bool
	Entries() []MeasurementEntry
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(measuredOn MeasuredOn, measuredBy staff.StaffID, ageAtMeasurement AgeAtMeasurement, updatedBy staff.StaffID, isDraft bool, entries []MeasurementEntry) error
}

type measurement struct {
	id               MeasurementID
	customerID       customer.CustomerID
	measuredOn       MeasuredOn
	measuredBy       staff.StaffID
	ageAtMeasurement AgeAtMeasurement
	updatedBy        staff.StaffID
	isDraft          bool
	entries          []MeasurementEntry
	createdAt        time.Time
	updatedAt        time.Time
}

func (m measurement) ID() MeasurementID {
	return m.id
}

func (m measurement) CustomerID() customer.CustomerID {
	return m.customerID
}

func (m measurement) MeasuredOn() MeasuredOn {
	return m.measuredOn
}

func (m measurement) MeasuredBy() staff.StaffID {
	return m.measuredBy
}

func (m measurement) AgeAtMeasurement() AgeAtMeasurement {
	return m.ageAtMeasurement
}

func (m measurement) UpdatedBy() staff.StaffID {
	return m.updatedBy
}

func (m measurement) IsDraft() bool {
	return m.isDraft
}

func (m measurement) Entries() []MeasurementEntry {
	entries := make([]MeasurementEntry, len(m.entries))
	copy(entries, m.entries)
	return entries
}

func (m measurement) CreatedAt() time.Time {
	return m.createdAt
}

func (m measurement) UpdatedAt() time.Time {
	return m.updatedAt
}

func (m *measurement) Update(measuredOn MeasuredOn, measuredBy staff.StaffID, ageAtMeasurement AgeAtMeasurement, updatedBy staff.StaffID, isDraft bool, entries []MeasurementEntry) error {
	if err := verifyEntries(isDraft, entries); err != nil {
		return err
	}

	m.measuredOn = measuredOn
	m.measuredBy = measuredBy
	m.ageAtMeasurement = ageAtMeasurement
	m.updatedBy = updatedBy
	m.isDraft = isDraft
	m.entries = make([]MeasurementEntry, len(entries))
	copy(m.entries, entries)
	m.updatedAt = time.Now().UTC()

	return nil
}

func verifyEntries(isDraft bool, entries []MeasurementEntry) error {
	if !isDraft && len(entries) == 0 {
		return ErrInvalidValueCount
	}

	seen := make(map[measurementitem.MeasurementItemID]struct{}, len(entries))
	for _, entry := range entries {
		if _, ok := seen[entry.MeasurementItemID()]; ok {
			return ErrDuplicateEntry
		}
		seen[entry.MeasurementItemID()] = struct{}{}

		if isDraft || entry.Unmeasurable() {
			continue
		}
		expectedValueCount, ok := entry.ExpectedValueCount()
		valueCount := len(entry.Values())
		if !ok || valueCount == 0 || valueCount > expectedValueCount {
			return ErrInvalidValueCount
		}
	}
	return nil
}

func NewMeasurement(
	id MeasurementID,
	customerID customer.CustomerID,
	measuredOn MeasuredOn,
	measuredBy staff.StaffID,
	ageAtMeasurement AgeAtMeasurement,
	updatedBy staff.StaffID,
	isDraft bool,
	entries []MeasurementEntry,
	createdAt time.Time,
	updatedAt time.Time,
) (Measurement, error) {
	if err := verifyEntries(isDraft, entries); err != nil {
		return nil, err
	}

	return newMeasurement(id, customerID, measuredOn, measuredBy, ageAtMeasurement, updatedBy, isDraft, entries, createdAt, updatedAt), nil
}

func ReconstructMeasurement(
	id MeasurementID,
	customerID customer.CustomerID,
	measuredOn MeasuredOn,
	measuredBy staff.StaffID,
	ageAtMeasurement AgeAtMeasurement,
	updatedBy staff.StaffID,
	isDraft bool,
	entries []MeasurementEntry,
	createdAt time.Time,
	updatedAt time.Time,
) Measurement {
	return newMeasurement(id, customerID, measuredOn, measuredBy, ageAtMeasurement, updatedBy, isDraft, entries, createdAt, updatedAt)
}

func newMeasurement(
	id MeasurementID,
	customerID customer.CustomerID,
	measuredOn MeasuredOn,
	measuredBy staff.StaffID,
	ageAtMeasurement AgeAtMeasurement,
	updatedBy staff.StaffID,
	isDraft bool,
	entries []MeasurementEntry,
	createdAt time.Time,
	updatedAt time.Time,
) Measurement {
	m := &measurement{
		id:               id,
		customerID:       customerID,
		measuredOn:       measuredOn,
		measuredBy:       measuredBy,
		ageAtMeasurement: ageAtMeasurement,
		updatedBy:        updatedBy,
		isDraft:          isDraft,
		entries:          make([]MeasurementEntry, len(entries)),
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
	copy(m.entries, entries)
	return m
}
