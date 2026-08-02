package standard

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
	"github.com/qkitzero/fitness-service/internal/domain/standard"
)

const (
	testGripStrengthID    = "45c2f5cd-ae75-4b2e-8302-69051f0343d5"
	testSitAndReachID     = "700ef351-eed2-478e-9e54-2574ea070ac2"
	testAgeGroupID40s     = "a2f5f8bb-8a9d-4c3b-9f0e-0a1cb1f1b1d0"
	testAgeGroupID50s     = "b3a6a9cc-9bae-4d4c-af1f-1b2dc2f2c2e1"
	testAgeGroupIDAnother = "c4b7badd-acbf-4e5d-bf20-2c3ed3f3d3f2"

	listByItemIDSQL           = `SELECT * FROM "age_group_standards" WHERE measurement_item_id = $1 ORDER BY gender, age_from`
	listByItemIDsAndGenderSQL = `SELECT * FROM "age_group_standards" WHERE measurement_item_id IN ($1,$2) AND gender = $3 ORDER BY measurement_item_id, age_from`
)

var ageGroupStandardColumns = []string{"id", "measurement_item_id", "gender", "age_from", "age_to", "mean", "standard_deviation", "created_at", "updated_at"}

func TestListByItemID(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	tests := []struct {
		name         string
		success      bool
		wantIDs      []string
		wantAgeFroms []int
		wantAgeTos   []int
		wantMeans    []float64
		setup        func(mock sqlmock.Sqlmock)
	}{
		{
			name:         "success list age group standards of a measurement item",
			success:      true,
			wantIDs:      []string{testAgeGroupID40s, testAgeGroupID50s},
			wantAgeFroms: []int{45, 50},
			wantAgeTos:   []int{49, 54},
			wantMeans:    []float64{46.5, 45.1},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listByItemIDSQL)).
					WithArgs(testGripStrengthID).
					WillReturnRows(sqlmock.NewRows(ageGroupStandardColumns).
						AddRow(testAgeGroupID40s, testGripStrengthID, "male", 45, 49, 46.5, 6.8, createdAt, updatedAt).
						AddRow(testAgeGroupID50s, testGripStrengthID, "male", 50, 54, 45.1, 6.9, createdAt, updatedAt))
			},
		},
		{
			name:         "success list no age group standards",
			success:      true,
			wantIDs:      []string{},
			wantAgeFroms: []int{},
			wantAgeTos:   []int{},
			wantMeans:    []float64{},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listByItemIDSQL)).
					WithArgs(testGripStrengthID).
					WillReturnRows(sqlmock.NewRows(ageGroupStandardColumns))
			},
		},
		{
			name:    "failure list age group standards error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listByItemIDSQL)).
					WithArgs(testGripStrengthID).
					WillReturnError(errors.New("list age group standards error"))
			},
		},
		{
			name:    "failure unknown gender",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listByItemIDSQL)).
					WithArgs(testGripStrengthID).
					WillReturnRows(sqlmock.NewRows(ageGroupStandardColumns).
						AddRow(testAgeGroupID40s, testGripStrengthID, "other", 45, 49, 46.5, 6.8, createdAt, updatedAt))
			},
		},
		{
			name:    "failure reversed age range",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listByItemIDSQL)).
					WithArgs(testGripStrengthID).
					WillReturnRows(sqlmock.NewRows(ageGroupStandardColumns).
						AddRow(testAgeGroupID40s, testGripStrengthID, "male", 49, 45, 46.5, 6.8, createdAt, updatedAt))
			},
		},
		{
			name:    "failure negative mean",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listByItemIDSQL)).
					WithArgs(testGripStrengthID).
					WillReturnRows(sqlmock.NewRows(ageGroupStandardColumns).
						AddRow(testAgeGroupID40s, testGripStrengthID, "male", 45, 49, -46.5, 6.8, createdAt, updatedAt))
			},
		},
		{
			name:    "failure zero standard deviation",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listByItemIDSQL)).
					WithArgs(testGripStrengthID).
					WillReturnRows(sqlmock.NewRows(ageGroupStandardColumns).
						AddRow(testAgeGroupID40s, testGripStrengthID, "male", 45, 49, 46.5, 0, createdAt, updatedAt))
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

			repo := NewAgeGroupStandardRepository(gormDB)

			measurementItemID, err := measurementitem.NewMeasurementItemIDFromString(testGripStrengthID)
			if err != nil {
				t.Fatalf("failed to new measurement item id: %v", err)
			}

			ageGroupStandards, err := repo.ListByItemID(context.Background(), measurementItemID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && ageGroupStandards != nil {
				t.Errorf("expected no age group standards on failure, but got %v", len(ageGroupStandards))
			}
			if tt.success {
				if len(ageGroupStandards) != len(tt.wantIDs) {
					t.Errorf("len(ageGroupStandards) = %v, want %v", len(ageGroupStandards), len(tt.wantIDs))
				}
				for i := range tt.wantIDs {
					if i >= len(ageGroupStandards) {
						break
					}
					a := ageGroupStandards[i]
					if a.ID().String() != tt.wantIDs[i] {
						t.Errorf("ageGroupStandards[%d].ID() = %v, want %v", i, a.ID().String(), tt.wantIDs[i])
					}
					if a.MeasurementItemID() != measurementItemID {
						t.Errorf("ageGroupStandards[%d].MeasurementItemID() = %v, want %v", i, a.MeasurementItemID(), measurementItemID)
					}
					if a.AgeRange().From() != tt.wantAgeFroms[i] {
						t.Errorf("ageGroupStandards[%d].AgeRange().From() = %v, want %v", i, a.AgeRange().From(), tt.wantAgeFroms[i])
					}
					if a.AgeRange().To() != tt.wantAgeTos[i] {
						t.Errorf("ageGroupStandards[%d].AgeRange().To() = %v, want %v", i, a.AgeRange().To(), tt.wantAgeTos[i])
					}
					if a.Mean().Float64() != tt.wantMeans[i] {
						t.Errorf("ageGroupStandards[%d].Mean() = %v, want %v", i, a.Mean().Float64(), tt.wantMeans[i])
					}
					if !a.CreatedAt().Equal(createdAt) {
						t.Errorf("ageGroupStandards[%d].CreatedAt() = %v, want %v", i, a.CreatedAt(), createdAt)
					}
					if !a.UpdatedAt().Equal(updatedAt) {
						t.Errorf("ageGroupStandards[%d].UpdatedAt() = %v, want %v", i, a.UpdatedAt(), updatedAt)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListByItemIDsAndGender(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	tests := []struct {
		name         string
		success      bool
		ids          []string
		gender       standard.Gender
		wantIDs      []string
		wantGenders  []string
		wantAgeFroms []int
		setup        func(mock sqlmock.Sqlmock)
	}{
		{
			name:         "success list age group standards of every age group",
			success:      true,
			ids:          []string{testGripStrengthID, testSitAndReachID},
			gender:       standard.GenderMale,
			wantIDs:      []string{testAgeGroupID40s, testAgeGroupID50s, testAgeGroupIDAnother},
			wantGenders:  []string{"male", "male", "male"},
			wantAgeFroms: []int{45, 50, 45},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listByItemIDsAndGenderSQL)).
					WithArgs(testGripStrengthID, testSitAndReachID, "male").
					WillReturnRows(sqlmock.NewRows(ageGroupStandardColumns).
						AddRow(testAgeGroupID40s, testGripStrengthID, "male", 45, 49, 46.5, 6.8, createdAt, updatedAt).
						AddRow(testAgeGroupID50s, testGripStrengthID, "male", 50, 54, 45.1, 6.9, createdAt, updatedAt).
						AddRow(testAgeGroupIDAnother, testSitAndReachID, "male", 45, 49, 38.2, 9.7, createdAt, updatedAt))
			},
		},
		{
			name:         "success no ids does not query",
			success:      true,
			ids:          nil,
			gender:       standard.GenderFemale,
			wantIDs:      []string{},
			wantGenders:  []string{},
			wantAgeFroms: []int{},
			setup:        func(mock sqlmock.Sqlmock) {},
		},
		{
			name:    "failure list age group standards error",
			success: false,
			ids:     []string{testGripStrengthID, testSitAndReachID},
			gender:  standard.GenderMale,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listByItemIDsAndGenderSQL)).
					WithArgs(testGripStrengthID, testSitAndReachID, "male").
					WillReturnError(errors.New("list age group standards error"))
			},
		},
		{
			name:    "failure age above the upper bound",
			success: false,
			ids:     []string{testGripStrengthID, testSitAndReachID},
			gender:  standard.GenderMale,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listByItemIDsAndGenderSQL)).
					WithArgs(testGripStrengthID, testSitAndReachID, "male").
					WillReturnRows(sqlmock.NewRows(ageGroupStandardColumns).
						AddRow(testAgeGroupID40s, testGripStrengthID, "male", 45, 151, 46.5, 6.8, createdAt, updatedAt))
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

			repo := NewAgeGroupStandardRepository(gormDB)

			measurementItemIDs := make([]measurementitem.MeasurementItemID, 0, len(tt.ids))
			for _, id := range tt.ids {
				measurementItemID, err := measurementitem.NewMeasurementItemIDFromString(id)
				if err != nil {
					t.Fatalf("failed to new measurement item id: %v", err)
				}
				measurementItemIDs = append(measurementItemIDs, measurementItemID)
			}

			ageGroupStandards, err := repo.ListByItemIDsAndGender(context.Background(), measurementItemIDs, tt.gender)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && ageGroupStandards != nil {
				t.Errorf("expected no age group standards on failure, but got %v", len(ageGroupStandards))
			}
			if tt.success {
				if len(ageGroupStandards) != len(tt.wantIDs) {
					t.Errorf("len(ageGroupStandards) = %v, want %v", len(ageGroupStandards), len(tt.wantIDs))
				}
				for i := range tt.wantIDs {
					if i >= len(ageGroupStandards) {
						break
					}
					a := ageGroupStandards[i]
					if a.ID().String() != tt.wantIDs[i] {
						t.Errorf("ageGroupStandards[%d].ID() = %v, want %v", i, a.ID().String(), tt.wantIDs[i])
					}
					if a.Gender().String() != tt.wantGenders[i] {
						t.Errorf("ageGroupStandards[%d].Gender() = %v, want %v", i, a.Gender().String(), tt.wantGenders[i])
					}
					if a.AgeRange().From() != tt.wantAgeFroms[i] {
						t.Errorf("ageGroupStandards[%d].AgeRange().From() = %v, want %v", i, a.AgeRange().From(), tt.wantAgeFroms[i])
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
