package measurementitem

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

const (
	listSQL                   = `SELECT * FROM "measurement_items"`
	findByIDsSQL              = `SELECT * FROM "measurement_items" WHERE id IN ($1,$2)`
	elementsByItemIDSQL       = `SELECT * FROM "measurement_item_elements" WHERE "measurement_item_elements"."measurement_item_id" = $1`
	elementsByTwoItemIDsSQL   = `SELECT * FROM "measurement_item_elements" WHERE "measurement_item_elements"."measurement_item_id" IN ($1,$2)`
	elementsByEightItemIDsSQL = `SELECT * FROM "measurement_item_elements" WHERE "measurement_item_elements"."measurement_item_id" IN ($1,$2,$3,$4,$5,$6,$7,$8)`

	bloodPressureID     = "c229b597-a5a0-4e02-9a47-10512f366005"
	pulseRateID         = "0506df56-d97f-445a-83fb-92b7843af9a4"
	heightID            = "80389534-2f6d-4e59-82b9-11f9f93d00cf"
	weightID            = "99b20390-501d-4663-a685-cff1d87a215b"
	bodyFatPercentageID = "e9b4d4cc-6edd-49d5-9ccd-ed2754e1c13d"
	muscleMassID        = "1d98b4c8-676e-4c1b-8d01-b092e571d6af"
	cs30ID              = "2625e4d7-7608-45d1-a41b-304f40c6c721"
	gripStrengthID      = "45c2f5cd-ae75-4b2e-8302-69051f0343d5"
	sideStepID          = "32ea7f00-d3ee-4197-ac0e-9373a033b69e"
	unknownID           = "00000000-0000-0000-0000-000000000000"
)

var (
	measurementItemColumns = []string{"id", "code", "name", "category", "unit", "trial_count", "bilateral", "value_type", "score_direction", "created_at", "updated_at"}
	elementColumns         = []string{"measurement_item_id", "element", "created_at", "updated_at"}
)

func TestList(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	tests := []struct {
		name                string
		success             bool
		wantIDs             []string
		wantCodes           []string
		wantNames           []string
		wantCategories      []string
		wantUnits           []string
		wantTrialCounts     []int
		wantBilaterals      []bool
		wantValueTypes      []string
		wantScoreDirections []string
		wantElements        [][]string
		setup               func(mock sqlmock.Sqlmock)
	}{
		{
			name:                "success list measurement items in category order then code order",
			success:             true,
			wantIDs:             []string{bloodPressureID, pulseRateID, heightID, weightID, bodyFatPercentageID, muscleMassID, cs30ID, gripStrengthID},
			wantCodes:           []string{"blood_pressure", "pulse_rate", "height", "weight", "body_fat_percentage", "muscle_mass", "cs30", "grip_strength"},
			wantNames:           []string{"血圧", "脈拍", "身長", "体重", "体脂肪率", "筋肉量", "CS-30（30秒立ち座り）", "握力"},
			wantCategories:      []string{"vital", "vital", "physique", "physique", "body_composition", "body_composition", "motor_function", "motor_function"},
			wantUnits:           []string{"mmHg", "bpm", "cm", "kg", "percent", "kg", "count", "kg"},
			wantTrialCounts:     []int{1, 1, 1, 1, 1, 1, 1, 2},
			wantBilaterals:      []bool{false, false, false, false, false, false, false, true},
			wantValueTypes:      []string{"paired", "numeric", "numeric", "numeric", "numeric", "numeric", "numeric", "numeric"},
			wantScoreDirections: []string{"", "", "", "", "", "", "higher_is_better", "higher_is_better"},
			wantElements:        [][]string{{}, {}, {}, {}, {}, {}, {"muscle_endurance"}, {"muscle_strength"}},
			setup: func(mock sqlmock.Sqlmock) {
				measurementItemRows := sqlmock.NewRows(measurementItemColumns).
					AddRow(muscleMassID, "muscle_mass", "筋肉量", "body_composition", "kg", 1, false, "numeric", nil, createdAt, updatedAt).
					AddRow(bodyFatPercentageID, "body_fat_percentage", "体脂肪率", "body_composition", "percent", 1, false, "numeric", nil, createdAt, updatedAt).
					AddRow(gripStrengthID, "grip_strength", "握力", "motor_function", "kg", 2, true, "numeric", "higher_is_better", createdAt, updatedAt).
					AddRow(cs30ID, "cs30", "CS-30（30秒立ち座り）", "motor_function", "count", 1, false, "numeric", "higher_is_better", createdAt, updatedAt).
					AddRow(weightID, "weight", "体重", "physique", "kg", 1, false, "numeric", nil, createdAt, updatedAt).
					AddRow(heightID, "height", "身長", "physique", "cm", 1, false, "numeric", nil, createdAt, updatedAt).
					AddRow(pulseRateID, "pulse_rate", "脈拍", "vital", "bpm", 1, false, "numeric", nil, createdAt, updatedAt).
					AddRow(bloodPressureID, "blood_pressure", "血圧", "vital", "mmHg", 1, false, "paired", nil, createdAt, updatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(measurementItemRows)
				mock.ExpectQuery(regexp.QuoteMeta(elementsByEightItemIDsSQL)).
					WithArgs(muscleMassID, bodyFatPercentageID, gripStrengthID, cs30ID, weightID, heightID, pulseRateID, bloodPressureID).
					WillReturnRows(sqlmock.NewRows(elementColumns).
						AddRow(gripStrengthID, "muscle_strength", createdAt, updatedAt).
						AddRow(cs30ID, "muscle_endurance", createdAt, updatedAt))
			},
		},
		{
			name:                "success list elements in business order",
			success:             true,
			wantIDs:             []string{sideStepID},
			wantCodes:           []string{"side_step"},
			wantNames:           []string{"反復横跳び"},
			wantCategories:      []string{"motor_function"},
			wantUnits:           []string{"count"},
			wantTrialCounts:     []int{1},
			wantBilaterals:      []bool{false},
			wantValueTypes:      []string{"numeric"},
			wantScoreDirections: []string{"higher_is_better"},
			wantElements:        [][]string{{"agility", "mobility"}},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(sideStepID, "side_step", "反復横跳び", "motor_function", "count", 1, false, "numeric", "higher_is_better", createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByItemIDSQL)).
					WithArgs(sideStepID).
					WillReturnRows(sqlmock.NewRows(elementColumns).
						AddRow(sideStepID, "mobility", createdAt, updatedAt).
						AddRow(sideStepID, "agility", createdAt, updatedAt))
			},
		},
		{
			name:                "success list no measurement items",
			success:             true,
			wantIDs:             []string{},
			wantCodes:           []string{},
			wantNames:           []string{},
			wantCategories:      []string{},
			wantUnits:           []string{},
			wantTrialCounts:     []int{},
			wantBilaterals:      []bool{},
			wantValueTypes:      []string{},
			wantScoreDirections: []string{},
			wantElements:        [][]string{},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns))
			},
		},
		{
			name:    "failure list measurement items error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnError(errors.New("list measurement items error"))
			},
		},
		{
			name:    "failure list measurement item elements error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(gripStrengthID, "grip_strength", "握力", "motor_function", "kg", 2, true, "numeric", "higher_is_better", createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByItemIDSQL)).
					WithArgs(gripStrengthID).
					WillReturnError(errors.New("list measurement item elements error"))
			},
		},
		{
			name:    "failure invalid id",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow("not-a-uuid", "grip_strength", "握力", "motor_function", "kg", 2, true, "numeric", "higher_is_better", createdAt, updatedAt))
			},
		},
		{
			name:    "failure invalid code",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(gripStrengthID, "GRIP STRENGTH!!", "握力", "motor_function", "kg", 2, true, "numeric", "higher_is_better", createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByItemIDSQL)).
					WithArgs(gripStrengthID).
					WillReturnRows(sqlmock.NewRows(elementColumns))
			},
		},
		{
			name:    "failure empty name",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(gripStrengthID, "grip_strength", "", "motor_function", "kg", 2, true, "numeric", "higher_is_better", createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByItemIDSQL)).
					WithArgs(gripStrengthID).
					WillReturnRows(sqlmock.NewRows(elementColumns))
			},
		},
		{
			name:    "failure unknown category",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(gripStrengthID, "grip_strength", "握力", "vitals", "kg", 2, true, "numeric", "higher_is_better", createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByItemIDSQL)).
					WithArgs(gripStrengthID).
					WillReturnRows(sqlmock.NewRows(elementColumns))
			},
		},
		{
			name:    "failure unknown unit",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(gripStrengthID, "grip_strength", "握力", "motor_function", "kilogram", 2, true, "numeric", "higher_is_better", createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByItemIDSQL)).
					WithArgs(gripStrengthID).
					WillReturnRows(sqlmock.NewRows(elementColumns))
			},
		},
		{
			name:    "failure trial count below one",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(gripStrengthID, "grip_strength", "握力", "motor_function", "kg", -1, true, "numeric", "higher_is_better", createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByItemIDSQL)).
					WithArgs(gripStrengthID).
					WillReturnRows(sqlmock.NewRows(elementColumns))
			},
		},
		{
			name:    "failure unknown value type",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(gripStrengthID, "grip_strength", "握力", "motor_function", "kg", 2, true, "bogus", "higher_is_better", createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByItemIDSQL)).
					WithArgs(gripStrengthID).
					WillReturnRows(sqlmock.NewRows(elementColumns))
			},
		},
		{
			name:    "failure unknown score direction",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(gripStrengthID, "grip_strength", "握力", "motor_function", "kg", 2, true, "numeric", "bigger_is_better", createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByItemIDSQL)).
					WithArgs(gripStrengthID).
					WillReturnRows(sqlmock.NewRows(elementColumns))
			},
		},
		{
			name:    "failure unknown element",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(gripStrengthID, "grip_strength", "握力", "motor_function", "kg", 2, true, "numeric", "higher_is_better", createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByItemIDSQL)).
					WithArgs(gripStrengthID).
					WillReturnRows(sqlmock.NewRows(elementColumns).
						AddRow(gripStrengthID, "power", createdAt, updatedAt))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.setup(mock)

			repo := NewMeasurementItemRepository(gormDB)

			measurementItems, err := repo.List(context.Background())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && measurementItems != nil {
				t.Errorf("expected no measurement items on failure, but got %v", len(measurementItems))
			}
			if tt.success {
				if len(measurementItems) != len(tt.wantCodes) {
					t.Errorf("len(measurementItems) = %v, want %v", len(measurementItems), len(tt.wantCodes))
				}
				for i := range tt.wantCodes {
					if i >= len(measurementItems) {
						break
					}
					m := measurementItems[i]
					if m.ID().String() != tt.wantIDs[i] {
						t.Errorf("measurementItems[%d].ID() = %v, want %v", i, m.ID().String(), tt.wantIDs[i])
					}
					if m.Code().String() != tt.wantCodes[i] {
						t.Errorf("measurementItems[%d].Code() = %v, want %v", i, m.Code().String(), tt.wantCodes[i])
					}
					if m.Name().String() != tt.wantNames[i] {
						t.Errorf("measurementItems[%d].Name() = %v, want %v", i, m.Name().String(), tt.wantNames[i])
					}
					if m.Category().String() != tt.wantCategories[i] {
						t.Errorf("measurementItems[%d].Category() = %v, want %v", i, m.Category().String(), tt.wantCategories[i])
					}
					if m.Unit().String() != tt.wantUnits[i] {
						t.Errorf("measurementItems[%d].Unit() = %v, want %v", i, m.Unit().String(), tt.wantUnits[i])
					}
					if m.TrialCount().Int() != tt.wantTrialCounts[i] {
						t.Errorf("measurementItems[%d].TrialCount() = %v, want %v", i, m.TrialCount().Int(), tt.wantTrialCounts[i])
					}
					if m.Bilateral() != tt.wantBilaterals[i] {
						t.Errorf("measurementItems[%d].Bilateral() = %v, want %v", i, m.Bilateral(), tt.wantBilaterals[i])
					}
					if m.ValueType().String() != tt.wantValueTypes[i] {
						t.Errorf("measurementItems[%d].ValueType() = %v, want %v", i, m.ValueType().String(), tt.wantValueTypes[i])
					}
					switch {
					case tt.wantScoreDirections[i] == "":
						if m.ScoreDirection() != nil {
							t.Errorf("measurementItems[%d].ScoreDirection() = %v, want nil", i, m.ScoreDirection())
						}
					case m.ScoreDirection() == nil:
						t.Errorf("measurementItems[%d].ScoreDirection() = nil, want %v", i, tt.wantScoreDirections[i])
					case m.ScoreDirection().String() != tt.wantScoreDirections[i]:
						t.Errorf("measurementItems[%d].ScoreDirection() = %v, want %v", i, m.ScoreDirection().String(), tt.wantScoreDirections[i])
					}
					if len(m.Elements()) != len(tt.wantElements[i]) {
						t.Errorf("len(measurementItems[%d].Elements()) = %v, want %v", i, len(m.Elements()), len(tt.wantElements[i]))
					}
					for j := range tt.wantElements[i] {
						if j >= len(m.Elements()) {
							break
						}
						if m.Elements()[j].String() != tt.wantElements[i][j] {
							t.Errorf("measurementItems[%d].Elements()[%d] = %v, want %v", i, j, m.Elements()[j].String(), tt.wantElements[i][j])
						}
					}
					if !m.CreatedAt().Equal(createdAt) {
						t.Errorf("measurementItems[%d].CreatedAt() = %v, want %v", i, m.CreatedAt(), createdAt)
					}
					if !m.UpdatedAt().Equal(updatedAt) {
						t.Errorf("measurementItems[%d].UpdatedAt() = %v, want %v", i, m.UpdatedAt(), updatedAt)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFindByIDs(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	tests := []struct {
		name         string
		success      bool
		ids          []string
		wantCodes    []string
		wantElements [][]string
		setup        func(mock sqlmock.Sqlmock)
	}{
		{
			name:         "success find measurement items by ids",
			success:      true,
			ids:          []string{gripStrengthID, pulseRateID},
			wantCodes:    []string{"grip_strength", "pulse_rate"},
			wantElements: [][]string{{"muscle_strength"}, {}},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(findByIDsSQL)).
					WithArgs(gripStrengthID, pulseRateID).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(gripStrengthID, "grip_strength", "握力", "motor_function", "kg", 2, true, "numeric", "higher_is_better", createdAt, updatedAt).
						AddRow(pulseRateID, "pulse_rate", "脈拍", "vital", "bpm", 1, false, "numeric", nil, createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByTwoItemIDsSQL)).
					WithArgs(gripStrengthID, pulseRateID).
					WillReturnRows(sqlmock.NewRows(elementColumns).
						AddRow(gripStrengthID, "muscle_strength", createdAt, updatedAt))
			},
		},
		{
			name:         "success unknown id is not returned",
			success:      true,
			ids:          []string{gripStrengthID, unknownID},
			wantCodes:    []string{"grip_strength"},
			wantElements: [][]string{{"muscle_strength"}},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(findByIDsSQL)).
					WithArgs(gripStrengthID, unknownID).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(gripStrengthID, "grip_strength", "握力", "motor_function", "kg", 2, true, "numeric", "higher_is_better", createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByItemIDSQL)).
					WithArgs(gripStrengthID).
					WillReturnRows(sqlmock.NewRows(elementColumns).
						AddRow(gripStrengthID, "muscle_strength", createdAt, updatedAt))
			},
		},
		{
			name:         "success no ids does not query",
			success:      true,
			ids:          nil,
			wantCodes:    []string{},
			wantElements: [][]string{},
			setup:        func(mock sqlmock.Sqlmock) {},
		},
		{
			name:    "failure find measurement items error",
			success: false,
			ids:     []string{gripStrengthID, pulseRateID},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(findByIDsSQL)).
					WithArgs(gripStrengthID, pulseRateID).
					WillReturnError(errors.New("find measurement items error"))
			},
		},
		{
			name:    "failure unknown category",
			success: false,
			ids:     []string{gripStrengthID, pulseRateID},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(findByIDsSQL)).
					WithArgs(gripStrengthID, pulseRateID).
					WillReturnRows(sqlmock.NewRows(measurementItemColumns).
						AddRow(gripStrengthID, "grip_strength", "握力", "vitals", "kg", 2, true, "numeric", "higher_is_better", createdAt, updatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(elementsByItemIDSQL)).
					WithArgs(gripStrengthID).
					WillReturnRows(sqlmock.NewRows(elementColumns))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.setup(mock)

			repo := NewMeasurementItemRepository(gormDB)

			measurementItemIDs := make([]measurementitem.MeasurementItemID, 0, len(tt.ids))
			for _, id := range tt.ids {
				measurementItemID, err := measurementitem.NewMeasurementItemIDFromString(id)
				if err != nil {
					t.Fatalf("failed to new measurement item id: %v", err)
				}
				measurementItemIDs = append(measurementItemIDs, measurementItemID)
			}

			measurementItems, err := repo.FindByIDs(context.Background(), measurementItemIDs)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && measurementItems != nil {
				t.Errorf("expected no measurement items on failure, but got %v", len(measurementItems))
			}
			if tt.success {
				if len(measurementItems) != len(tt.wantCodes) {
					t.Errorf("len(measurementItems) = %v, want %v", len(measurementItems), len(tt.wantCodes))
				}
				for i := range tt.wantCodes {
					if i >= len(measurementItems) {
						break
					}
					if measurementItems[i].Code().String() != tt.wantCodes[i] {
						t.Errorf("measurementItems[%d].Code() = %v, want %v", i, measurementItems[i].Code().String(), tt.wantCodes[i])
					}
					if len(measurementItems[i].Elements()) != len(tt.wantElements[i]) {
						t.Errorf("len(measurementItems[%d].Elements()) = %v, want %v", i, len(measurementItems[i].Elements()), len(tt.wantElements[i]))
					}
					for j := range tt.wantElements[i] {
						if j >= len(measurementItems[i].Elements()) {
							break
						}
						if measurementItems[i].Elements()[j].String() != tt.wantElements[i][j] {
							t.Errorf("measurementItems[%d].Elements()[%d] = %v, want %v", i, j, measurementItems[i].Elements()[j].String(), tt.wantElements[i][j])
						}
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
