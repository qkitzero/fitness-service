package measurement

import (
	"errors"
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

func TestNewMeasurementEntry(t *testing.T) {
	t.Parallel()
	itemCreatedAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	itemUpdatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	vital, _ := measurementitem.NewCategory("vital")
	motorFunction, _ := measurementitem.NewCategory("motor_function")
	bpm, _ := measurementitem.NewUnit("bpm")
	mmHg, _ := measurementitem.NewUnit("mmHg")
	kg, _ := measurementitem.NewUnit("kg")
	cm, _ := measurementitem.NewUnit("cm")
	singleTrial, _ := measurementitem.NewTrialCount(1)
	twoTrials, _ := measurementitem.NewTrialCount(2)
	pulseRateCode, _ := measurementitem.NewCode("pulse_rate")
	pulseRateName, _ := measurementitem.NewName("脈拍")
	bloodPressureCode, _ := measurementitem.NewCode("blood_pressure")
	bloodPressureName, _ := measurementitem.NewName("血圧")
	gripStrengthCode, _ := measurementitem.NewCode("grip_strength")
	gripStrengthName, _ := measurementitem.NewName("握力")
	standUpTestCode, _ := measurementitem.NewCode("stand_up_test")
	standUpTestName, _ := measurementitem.NewName("立ち上がり")

	pulseRate := measurementitem.NewMeasurementItem(measurementitem.NewMeasurementItemID(), pulseRateCode, pulseRateName, vital, bpm, singleTrial, false, measurementitem.ValueTypeNumeric, nil, measurementitem.SideAggregationMean, nil, itemCreatedAt, itemUpdatedAt)
	bloodPressure := measurementitem.NewMeasurementItem(measurementitem.NewMeasurementItemID(), bloodPressureCode, bloodPressureName, vital, mmHg, singleTrial, false, measurementitem.ValueTypePaired, nil, measurementitem.SideAggregationMean, nil, itemCreatedAt, itemUpdatedAt)
	gripStrength := measurementitem.NewMeasurementItem(measurementitem.NewMeasurementItemID(), gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, true, measurementitem.ValueTypeNumeric, nil, measurementitem.SideAggregationMean, nil, itemCreatedAt, itemUpdatedAt)
	standUpTest := measurementitem.NewMeasurementItem(measurementitem.NewMeasurementItemID(), standUpTestCode, standUpTestName, motorFunction, cm, singleTrial, true, measurementitem.ValueTypeChoice, nil, measurementitem.SideAggregationMean, nil, itemCreatedAt, itemUpdatedAt)
	unknownValueType := measurementitem.NewMeasurementItem(measurementitem.NewMeasurementItemID(), pulseRateCode, pulseRateName, vital, bpm, singleTrial, false, measurementitem.ValueType("range"), nil, measurementitem.SideAggregationMean, nil, itemCreatedAt, itemUpdatedAt)
	note, _ := NewNote("ふらつきあり")

	tests := []struct {
		name                   string
		item                   measurementitem.MeasurementItem
		unmeasurable           bool
		note                   *Note
		values                 func() []MeasurementValue
		wantErr                error
		wantExpectedValueCount int
	}{
		{
			name: "success numeric single trial",
			item: pulseRate,
			note: note,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(72)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideNone, &value, nil, nil)}
			},
			wantExpectedValueCount: 1,
		},
		{
			name: "success paired",
			item: bloodPressure,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(128)
				valueSecondary, _ := NewValue(82)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideNone, &value, &valueSecondary, nil)}
			},
			wantExpectedValueCount: 1,
		},
		{
			name: "success numeric two trials on both sides",
			item: gripStrength,
			values: func() []MeasurementValue {
				firstTrial, _ := NewTrialIndex(1)
				secondTrial, _ := NewTrialIndex(2)
				firstLeft, _ := NewValue(32.4)
				firstRight, _ := NewValue(33.1)
				secondLeft, _ := NewValue(31.8)
				secondRight, _ := NewValue(33.5)
				return []MeasurementValue{
					NewMeasurementValue(firstTrial, SideLeft, &firstLeft, nil, nil),
					NewMeasurementValue(firstTrial, SideRight, &firstRight, nil, nil),
					NewMeasurementValue(secondTrial, SideLeft, &secondLeft, nil, nil),
					NewMeasurementValue(secondTrial, SideRight, &secondRight, nil, nil),
				}
			},
			wantExpectedValueCount: 4,
		},
		{
			name: "success choice on both sides",
			item: standUpTest,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				left, _ := NewChoice("片足20cm")
				right, _ := NewChoice("両足30cm")
				return []MeasurementValue{
					NewMeasurementValue(trialIndex, SideLeft, nil, nil, &left),
					NewMeasurementValue(trialIndex, SideRight, nil, nil, &right),
				}
			},
			wantExpectedValueCount: 2,
		},
		{
			name:                   "success unmeasurable without values",
			item:                   gripStrength,
			unmeasurable:           true,
			values:                 func() []MeasurementValue { return nil },
			wantExpectedValueCount: 4,
		},
		{
			name: "success partial values",
			item: gripStrength,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(32.4)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideLeft, &value, nil, nil)}
			},
			wantExpectedValueCount: 4,
		},
		{
			name:         "failure unmeasurable with values",
			item:         gripStrength,
			unmeasurable: true,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(32.4)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideLeft, &value, nil, nil)}
			},
			wantErr: ErrInvalidValueCount,
		},
		{
			name: "failure side none on bilateral item",
			item: gripStrength,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(32.4)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideNone, &value, nil, nil)}
			},
			wantErr: ErrInvalidSide,
		},
		{
			name: "failure side left on non bilateral item",
			item: pulseRate,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(72)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideLeft, &value, nil, nil)}
			},
			wantErr: ErrInvalidSide,
		},
		{
			name: "failure numeric without value",
			item: pulseRate,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideNone, nil, nil, nil)}
			},
			wantErr: ErrInvalidValueForType,
		},
		{
			name: "failure numeric with secondary value",
			item: pulseRate,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(72)
				valueSecondary, _ := NewValue(60)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideNone, &value, &valueSecondary, nil)}
			},
			wantErr: ErrInvalidValueForType,
		},
		{
			name: "failure paired without secondary value",
			item: bloodPressure,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(128)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideNone, &value, nil, nil)}
			},
			wantErr: ErrInvalidValueForType,
		},
		{
			name: "failure paired with choice value",
			item: bloodPressure,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(128)
				valueSecondary, _ := NewValue(82)
				valueChoice, _ := NewChoice("片足20cm")
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideNone, &value, &valueSecondary, &valueChoice)}
			},
			wantErr: ErrInvalidValueForType,
		},
		{
			name: "failure choice without choice value",
			item: standUpTest,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideLeft, nil, nil, nil)}
			},
			wantErr: ErrInvalidValueForType,
		},
		{
			name: "failure choice with numeric value",
			item: standUpTest,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(20)
				valueChoice, _ := NewChoice("片足20cm")
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideLeft, &value, nil, &valueChoice)}
			},
			wantErr: ErrInvalidValueForType,
		},
		{
			name: "failure unknown value type",
			item: unknownValueType,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(72)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideNone, &value, nil, nil)}
			},
			wantErr: ErrInvalidValueForType,
		},
		{
			name: "failure trial index over trial count",
			item: gripStrength,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(3)
				value, _ := NewValue(32.4)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideLeft, &value, nil, nil)}
			},
			wantErr: ErrInvalidValueCount,
		},
		{
			name: "failure duplicate trial index and side",
			item: gripStrength,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				first, _ := NewValue(32.4)
				second, _ := NewValue(33.1)
				return []MeasurementValue{
					NewMeasurementValue(trialIndex, SideLeft, &first, nil, nil),
					NewMeasurementValue(trialIndex, SideLeft, &second, nil, nil),
				}
			},
			wantErr: ErrDuplicateValue,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			values := tt.values()

			entry, err := NewMeasurementEntry(tt.item, tt.unmeasurable, tt.note, values)
			if tt.wantErr == nil && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if entry != nil {
					t.Errorf("expected no entry on failure, but got %v", entry)
				}
				return
			}

			if entry.MeasurementItemID() != tt.item.ID() {
				t.Errorf("MeasurementItemID() = %v, want %v", entry.MeasurementItemID(), tt.item.ID())
			}
			if entry.Unmeasurable() != tt.unmeasurable {
				t.Errorf("Unmeasurable() = %v, want %v", entry.Unmeasurable(), tt.unmeasurable)
			}
			if tt.note == nil && entry.Note() != nil {
				t.Errorf("Note() = %v, want nil", entry.Note())
			}
			if tt.note != nil && (entry.Note() == nil || *entry.Note() != *tt.note) {
				t.Errorf("Note() = %v, want %v", entry.Note(), *tt.note)
			}
			if len(entry.Values()) != len(values) {
				t.Errorf("len(Values()) = %v, want %v", len(entry.Values()), len(values))
			}
			expectedValueCount, known := entry.ExpectedValueCount()
			if !known {
				t.Errorf("ExpectedValueCount() known = %v, want %v", known, true)
			}
			if expectedValueCount != tt.wantExpectedValueCount {
				t.Errorf("ExpectedValueCount() = %v, want %v", expectedValueCount, tt.wantExpectedValueCount)
			}
		})
	}
}

func TestReconstructMeasurementEntry(t *testing.T) {
	t.Parallel()
	note, _ := NewNote("ふらつきあり")

	tests := []struct {
		name         string
		unmeasurable bool
		note         *Note
		values       func() []MeasurementValue
	}{
		{
			name: "success reconstruct measurement entry",
			note: note,
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(1)
				value, _ := NewValue(72)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideNone, &value, nil, nil)}
			},
		},
		{
			name:         "success reconstruct unmeasurable entry without values",
			unmeasurable: true,
			values:       func() []MeasurementValue { return nil },
		},
		{
			name: "success reconstruct entry has no expected value count",
			values: func() []MeasurementValue {
				trialIndex, _ := NewTrialIndex(9)
				value, _ := NewValue(72)
				return []MeasurementValue{NewMeasurementValue(trialIndex, SideLeft, &value, nil, nil)}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			measurementItemID := measurementitem.NewMeasurementItemID()
			values := tt.values()

			entry := ReconstructMeasurementEntry(measurementItemID, tt.unmeasurable, tt.note, values)

			if entry.MeasurementItemID() != measurementItemID {
				t.Errorf("MeasurementItemID() = %v, want %v", entry.MeasurementItemID(), measurementItemID)
			}
			if entry.Unmeasurable() != tt.unmeasurable {
				t.Errorf("Unmeasurable() = %v, want %v", entry.Unmeasurable(), tt.unmeasurable)
			}
			if tt.note == nil && entry.Note() != nil {
				t.Errorf("Note() = %v, want nil", entry.Note())
			}
			if tt.note != nil && (entry.Note() == nil || *entry.Note() != *tt.note) {
				t.Errorf("Note() = %v, want %v", entry.Note(), *tt.note)
			}
			if len(entry.Values()) != len(values) {
				t.Errorf("len(Values()) = %v, want %v", len(entry.Values()), len(values))
			}
			if _, known := entry.ExpectedValueCount(); known {
				t.Errorf("ExpectedValueCount() known = %v, want %v", known, false)
			}
		})
	}
}
