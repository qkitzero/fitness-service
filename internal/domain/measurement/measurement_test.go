package measurement

import (
	"errors"
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/staff"
)

func TestNewMeasurement(t *testing.T) {
	t.Parallel()
	itemCreatedAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	itemUpdatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	motorFunction, _ := measurementitem.NewCategory("motor_function")
	kg, _ := measurementitem.NewUnit("kg")
	twoTrials, _ := measurementitem.NewTrialCount(2)
	gripStrengthCode, _ := measurementitem.NewCode("grip_strength")
	gripStrengthName, _ := measurementitem.NewName("握力")
	gripStrength := measurementitem.NewMeasurementItem(measurementitem.NewMeasurementItemID(), gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, measurementitem.SideModeBilateral, measurementitem.ValueTypeNumeric, nil, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, nil, itemCreatedAt, itemUpdatedAt)

	measuredOn, _ := NewMeasuredOn(2026, 8, 1)
	measuredBy, _ := staff.NewStaffID("google-oauth2|000000000000000000000")
	updatedBy, _ := staff.NewStaffID("google-oauth2|111111111111111111111")
	ageAtMeasurement, _ := NewAgeAtMeasurement(65)
	createdAt := time.Date(2026, 8, 1, 1, 2, 3, 0, time.UTC)
	updatedAt := time.Date(2026, 8, 1, 4, 5, 6, 0, time.UTC)

	tests := []struct {
		name       string
		isDraft    bool
		entries    func() []MeasurementEntry
		wantErr    error
		wantValues int
	}{
		{
			name:    "success new draft measurement without entries",
			isDraft: true,
			entries: func() []MeasurementEntry { return nil },
		},
		{
			name:    "failure no entries on a confirmed measurement",
			entries: func() []MeasurementEntry { return nil },
			wantErr: ErrInvalidValueCount,
		},
		{
			name: "success new measurement with complete entries",
			entries: func() []MeasurementEntry {
				firstTrial, _ := NewTrialIndex(1)
				secondTrial, _ := NewTrialIndex(2)
				firstLeft, _ := NewValue(32.4)
				firstRight, _ := NewValue(33.1)
				secondLeft, _ := NewValue(31.8)
				secondRight, _ := NewValue(33.5)
				entry, _ := NewMeasurementEntry(gripStrength, false, nil, []MeasurementValue{
					NewMeasurementValue(firstTrial, SideLeft, &firstLeft, nil, nil),
					NewMeasurementValue(firstTrial, SideRight, &firstRight, nil, nil),
					NewMeasurementValue(secondTrial, SideLeft, &secondLeft, nil, nil),
					NewMeasurementValue(secondTrial, SideRight, &secondRight, nil, nil),
				})
				return []MeasurementEntry{entry}
			},
			wantValues: 4,
		},
		{
			name:    "success draft measurement with partial entries",
			isDraft: true,
			entries: func() []MeasurementEntry {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(32.4)
				entry, _ := NewMeasurementEntry(gripStrength, false, nil, []MeasurementValue{
					NewMeasurementValue(trialIndex, SideLeft, &value, nil, nil),
				})
				return []MeasurementEntry{entry}
			},
			wantValues: 1,
		},
		{
			name: "success new measurement with unmeasurable entry",
			entries: func() []MeasurementEntry {
				note, _ := NewNote("肩の痛みのため測定不可")
				entry, _ := NewMeasurementEntry(gripStrength, true, note, nil)
				return []MeasurementEntry{entry}
			},
		},
		{
			name: "success confirmed measurement with fewer values than the trial capacity",
			entries: func() []MeasurementEntry {
				trialIndex, _ := NewTrialIndex(1)
				left, _ := NewValue(32.4)
				right, _ := NewValue(33.1)
				entry, _ := NewMeasurementEntry(gripStrength, false, nil, []MeasurementValue{
					NewMeasurementValue(trialIndex, SideLeft, &left, nil, nil),
					NewMeasurementValue(trialIndex, SideRight, &right, nil, nil),
				})
				return []MeasurementEntry{entry}
			},
			wantValues: 2,
		},
		{
			name: "failure no values on a confirmed measurable entry",
			entries: func() []MeasurementEntry {
				entry, _ := NewMeasurementEntry(gripStrength, false, nil, nil)
				return []MeasurementEntry{entry}
			},
			wantErr: ErrInvalidValueCount,
		},
		{
			name:    "failure duplicate measurement item",
			isDraft: true,
			entries: func() []MeasurementEntry {
				first, _ := NewMeasurementEntry(gripStrength, true, nil, nil)
				second, _ := NewMeasurementEntry(gripStrength, true, nil, nil)
				return []MeasurementEntry{first, second}
			},
			wantErr: ErrDuplicateEntry,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewMeasurementID()
			customerID := customer.NewCustomerID()
			entries := tt.entries()

			m, err := NewMeasurement(id, customerID, measuredOn, measuredBy, ageAtMeasurement, updatedBy, tt.isDraft, entries, createdAt, updatedAt)
			if tt.wantErr == nil && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if m != nil {
					t.Errorf("expected no measurement on failure, but got %v", m)
				}
				return
			}

			if m.ID() != id {
				t.Errorf("ID() = %v, want %v", m.ID(), id)
			}
			if m.CustomerID() != customerID {
				t.Errorf("CustomerID() = %v, want %v", m.CustomerID(), customerID)
			}
			if !m.MeasuredOn().Equal(measuredOn.Time) {
				t.Errorf("MeasuredOn() = %v, want %v", m.MeasuredOn(), measuredOn)
			}
			if m.MeasuredBy() != measuredBy {
				t.Errorf("MeasuredBy() = %v, want %v", m.MeasuredBy(), measuredBy)
			}
			if m.AgeAtMeasurement() != ageAtMeasurement {
				t.Errorf("AgeAtMeasurement() = %v, want %v", m.AgeAtMeasurement(), ageAtMeasurement)
			}
			if m.UpdatedBy() != updatedBy {
				t.Errorf("UpdatedBy() = %v, want %v", m.UpdatedBy(), updatedBy)
			}
			if m.IsDraft() != tt.isDraft {
				t.Errorf("IsDraft() = %v, want %v", m.IsDraft(), tt.isDraft)
			}
			if len(m.Entries()) != len(entries) {
				t.Errorf("len(Entries()) = %v, want %v", len(m.Entries()), len(entries))
			}
			values := 0
			for _, entry := range m.Entries() {
				values += len(entry.Values())
			}
			if values != tt.wantValues {
				t.Errorf("values = %v, want %v", values, tt.wantValues)
			}
			if !m.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", m.CreatedAt(), createdAt)
			}
			if !m.UpdatedAt().Equal(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", m.UpdatedAt(), updatedAt)
			}
		})
	}
}

func TestReconstructMeasurement(t *testing.T) {
	t.Parallel()
	measuredOn, _ := NewMeasuredOn(2026, 8, 1)
	measuredBy, _ := staff.NewStaffID("google-oauth2|000000000000000000000")
	updatedBy, _ := staff.NewStaffID("google-oauth2|111111111111111111111")
	ageAtMeasurement, _ := NewAgeAtMeasurement(65)
	createdAt := time.Date(2026, 8, 1, 1, 2, 3, 0, time.UTC)
	updatedAt := time.Date(2026, 8, 1, 4, 5, 6, 0, time.UTC)

	tests := []struct {
		name    string
		isDraft bool
		entries func() []MeasurementEntry
	}{
		{
			name:    "success reconstruct measurement without entries",
			entries: func() []MeasurementEntry { return nil },
		},
		{
			name: "success reconstruct measurement whose entries are not validated",
			entries: func() []MeasurementEntry {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(32.4)
				return []MeasurementEntry{
					ReconstructMeasurementEntry(measurementitem.NewMeasurementItemID(), false, nil, []MeasurementValue{
						NewMeasurementValue(trialIndex, SideLeft, &value, nil, nil),
					}),
				}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewMeasurementID()
			customerID := customer.NewCustomerID()
			entries := tt.entries()

			m := ReconstructMeasurement(id, customerID, measuredOn, measuredBy, ageAtMeasurement, updatedBy, tt.isDraft, entries, createdAt, updatedAt)

			if m.ID() != id {
				t.Errorf("ID() = %v, want %v", m.ID(), id)
			}
			if m.CustomerID() != customerID {
				t.Errorf("CustomerID() = %v, want %v", m.CustomerID(), customerID)
			}
			if m.IsDraft() != tt.isDraft {
				t.Errorf("IsDraft() = %v, want %v", m.IsDraft(), tt.isDraft)
			}
			if len(m.Entries()) != len(entries) {
				t.Errorf("len(Entries()) = %v, want %v", len(m.Entries()), len(entries))
			}
			if !m.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", m.CreatedAt(), createdAt)
			}
			if !m.UpdatedAt().Equal(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", m.UpdatedAt(), updatedAt)
			}
		})
	}
}

func TestUpdateMeasurement(t *testing.T) {
	t.Parallel()
	itemCreatedAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	itemUpdatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	vital, _ := measurementitem.NewCategory("vital")
	bpm, _ := measurementitem.NewUnit("bpm")
	singleTrial, _ := measurementitem.NewTrialCount(1)
	pulseRateCode, _ := measurementitem.NewCode("pulse_rate")
	pulseRateName, _ := measurementitem.NewName("脈拍")
	pulseRate := measurementitem.NewMeasurementItem(measurementitem.NewMeasurementItemID(), pulseRateCode, pulseRateName, vital, bpm, singleTrial, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, nil, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, nil, itemCreatedAt, itemUpdatedAt)

	measuredOn, _ := NewMeasuredOn(2026, 8, 1)
	updatedMeasuredOn, _ := NewMeasuredOn(2026, 7, 31)
	measuredBy, _ := staff.NewStaffID("google-oauth2|000000000000000000000")
	updatedMeasuredBy, _ := staff.NewStaffID("google-oauth2|222222222222222222222")
	updatedBy, _ := staff.NewStaffID("google-oauth2|111111111111111111111")
	newUpdatedBy, _ := staff.NewStaffID("google-oauth2|333333333333333333333")
	ageAtMeasurement, _ := NewAgeAtMeasurement(65)
	updatedAgeAtMeasurement, _ := NewAgeAtMeasurement(64)
	createdAt := time.Date(2026, 8, 1, 1, 2, 3, 0, time.UTC)
	updatedAt := time.Date(2026, 8, 1, 4, 5, 6, 0, time.UTC)

	tests := []struct {
		name       string
		isDraft    bool
		entries    func() []MeasurementEntry
		wantErr    error
		wantValues int
	}{
		{
			name: "success update measurement with complete entries",
			entries: func() []MeasurementEntry {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(72)
				entry, _ := NewMeasurementEntry(pulseRate, false, nil, []MeasurementValue{
					NewMeasurementValue(trialIndex, SideNone, &value, nil, nil),
				})
				return []MeasurementEntry{entry}
			},
			wantValues: 1,
		},
		{
			name:    "success update measurement to a draft without values",
			isDraft: true,
			entries: func() []MeasurementEntry {
				entry, _ := NewMeasurementEntry(pulseRate, false, nil, nil)
				return []MeasurementEntry{entry}
			},
		},
		{
			name: "failure partial entries on a confirmed measurement",
			entries: func() []MeasurementEntry {
				entry, _ := NewMeasurementEntry(pulseRate, false, nil, nil)
				return []MeasurementEntry{entry}
			},
			wantErr: ErrInvalidValueCount,
		},
		{
			name:    "failure duplicate measurement item",
			isDraft: true,
			entries: func() []MeasurementEntry {
				first, _ := NewMeasurementEntry(pulseRate, true, nil, nil)
				second, _ := NewMeasurementEntry(pulseRate, true, nil, nil)
				return []MeasurementEntry{first, second}
			},
			wantErr: ErrDuplicateEntry,
		},
		{
			name:    "failure no entries on a confirmed measurement",
			entries: func() []MeasurementEntry { return nil },
			wantErr: ErrInvalidValueCount,
		},
		{
			name: "failure reconstructed entry cannot confirm a measurement",
			entries: func() []MeasurementEntry {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(72)
				return []MeasurementEntry{
					ReconstructMeasurementEntry(pulseRate.ID(), false, nil, []MeasurementValue{
						NewMeasurementValue(trialIndex, SideNone, &value, nil, nil),
					}),
				}
			},
			wantErr: ErrInvalidValueCount,
		},
		{
			name:    "success reconstructed entry on a draft measurement",
			isDraft: true,
			entries: func() []MeasurementEntry {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(72)
				return []MeasurementEntry{
					ReconstructMeasurementEntry(pulseRate.ID(), false, nil, []MeasurementValue{
						NewMeasurementValue(trialIndex, SideNone, &value, nil, nil),
					}),
				}
			},
			wantValues: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := ReconstructMeasurement(NewMeasurementID(), customer.NewCustomerID(), measuredOn, measuredBy, ageAtMeasurement, updatedBy, false, nil, createdAt, updatedAt)
			entries := tt.entries()

			err := m.Update(updatedMeasuredOn, updatedMeasuredBy, updatedAgeAtMeasurement, newUpdatedBy, tt.isDraft, entries)
			if tt.wantErr == nil && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if !m.MeasuredOn().Equal(measuredOn.Time) {
					t.Errorf("MeasuredOn() = %v, want %v", m.MeasuredOn(), measuredOn)
				}
				if len(m.Entries()) != 0 {
					t.Errorf("len(Entries()) = %v, want %v", len(m.Entries()), 0)
				}
				if !m.UpdatedAt().Equal(updatedAt) {
					t.Errorf("UpdatedAt() = %v, want %v", m.UpdatedAt(), updatedAt)
				}
				return
			}

			if !m.MeasuredOn().Equal(updatedMeasuredOn.Time) {
				t.Errorf("MeasuredOn() = %v, want %v", m.MeasuredOn(), updatedMeasuredOn)
			}
			if m.MeasuredBy() != updatedMeasuredBy {
				t.Errorf("MeasuredBy() = %v, want %v", m.MeasuredBy(), updatedMeasuredBy)
			}
			if m.AgeAtMeasurement() != updatedAgeAtMeasurement {
				t.Errorf("AgeAtMeasurement() = %v, want %v", m.AgeAtMeasurement(), updatedAgeAtMeasurement)
			}
			if m.UpdatedBy() != newUpdatedBy {
				t.Errorf("UpdatedBy() = %v, want %v", m.UpdatedBy(), newUpdatedBy)
			}
			if m.IsDraft() != tt.isDraft {
				t.Errorf("IsDraft() = %v, want %v", m.IsDraft(), tt.isDraft)
			}
			if len(m.Entries()) != len(entries) {
				t.Errorf("len(Entries()) = %v, want %v", len(m.Entries()), len(entries))
			}
			values := 0
			for _, entry := range m.Entries() {
				values += len(entry.Values())
			}
			if values != tt.wantValues {
				t.Errorf("values = %v, want %v", values, tt.wantValues)
			}
			if !m.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", m.CreatedAt(), createdAt)
			}
			if !m.UpdatedAt().After(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want after %v", m.UpdatedAt(), updatedAt)
			}
		})
	}
}
