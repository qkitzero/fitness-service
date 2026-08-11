package measurement

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/staff"
)

const (
	testMeasurementID      = "11111111-1111-1111-1111-111111111111"
	testCustomerID         = "22222222-2222-2222-2222-222222222222"
	testMeasurementItemID  = "33333333-3333-3333-3333-333333333333"
	testMeasurementEntryID = "44444444-4444-4444-4444-444444444444"
	testStaffID            = "google-oauth2|000000000000000000000"

	insertMeasurementSQL = `INSERT INTO "measurements" ("id","customer_id","measured_on","measured_by","age_at_measurement","updated_by","is_draft","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	insertEntrySQL       = `INSERT INTO "measurement_entries" ("id","measurement_id","measurement_item_id","unmeasurable","note","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7)`
	insertValuesSQL      = `INSERT INTO "measurement_values" ("id","measurement_entry_id","trial_index","side","value","value_secondary","value_choice","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9),($10,$11,$12,$13,$14,$15,$16,$17,$18)`
	selectMeasurementSQL = `SELECT * FROM "measurements" WHERE id = $1 ORDER BY "measurements"."id" LIMIT $2`
	selectEntriesSQL     = `SELECT * FROM "measurement_entries" WHERE "measurement_entries"."measurement_id" = $1 ORDER BY measurement_entries.measurement_item_id`
	selectValuesSQL      = `SELECT * FROM "measurement_values" WHERE "measurement_values"."measurement_entry_id" = $1 ORDER BY measurement_values.trial_index, measurement_values.side`
)

var (
	testMeasuredOn  = time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	testCreatedAt   = time.Date(2026, 8, 1, 1, 2, 3, 0, time.UTC)
	testUpdatedAt   = time.Date(2026, 8, 2, 4, 5, 6, 0, time.UTC)
	measurementCols = []string{"id", "customer_id", "measured_on", "measured_by", "age_at_measurement", "updated_by", "is_draft", "created_at", "updated_at"}
	entryCols       = []string{"id", "measurement_id", "measurement_item_id", "unmeasurable", "note", "created_at", "updated_at"}
	valueCols       = []string{"id", "measurement_entry_id", "trial_index", "side", "value", "value_secondary", "value_choice", "created_at", "updated_at"}
)

func TestCreate(t *testing.T) {
	t.Parallel()
	measurementID, _ := measurement.NewMeasurementIDFromString(testMeasurementID)
	customerID, _ := customer.NewCustomerIDFromString(testCustomerID)
	measurementItemID, _ := measurementitem.NewMeasurementItemIDFromString(testMeasurementItemID)
	measuredOn, _ := measurement.NewMeasuredOn(2026, 8, 1)
	measuredBy, _ := staff.NewStaffID(testStaffID)
	ageAtMeasurement, _ := measurement.NewAgeAtMeasurement(65)
	trialIndex, _ := measurement.NewTrialIndex(1)
	left, _ := measurement.NewValue(32.4)
	right, _ := measurement.NewValue(33.1)
	note, _ := measurement.NewNote("ふらつきあり")
	entry := measurement.ReconstructMeasurementEntry(measurementItemID, false, note, []measurement.MeasurementValue{
		measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &left, nil, nil),
		measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &right, nil, nil),
	})
	withEntries := measurement.ReconstructMeasurement(measurementID, customerID, measuredOn, measuredBy, ageAtMeasurement, measuredBy, false, []measurement.MeasurementEntry{entry}, testCreatedAt, testCreatedAt)
	withoutEntries := measurement.ReconstructMeasurement(measurementID, customerID, measuredOn, measuredBy, ageAtMeasurement, measuredBy, true, nil, testCreatedAt, testCreatedAt)

	tests := []struct {
		name        string
		success     bool
		measurement measurement.Measurement
		setup       func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "success create measurement with entries",
			success:     true,
			measurement: withEntries,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(insertMeasurementSQL)).
					WithArgs(testMeasurementID, testCustomerID, testMeasuredOn, testStaffID, 65, testStaffID, false, testCreatedAt, testCreatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(insertEntrySQL)).
					WithArgs(sqlmock.AnyArg(), testMeasurementID, testMeasurementItemID, false, "ふらつきあり", testCreatedAt, testCreatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(insertValuesSQL)).
					WithArgs(
						sqlmock.AnyArg(), sqlmock.AnyArg(), 1, "left", 32.4, nil, nil, testCreatedAt, testCreatedAt,
						sqlmock.AnyArg(), sqlmock.AnyArg(), 1, "right", 33.1, nil, nil, testCreatedAt, testCreatedAt,
					).
					WillReturnResult(sqlmock.NewResult(1, 2))
				mock.ExpectCommit()
			},
		},
		{
			name:        "success create measurement without entries",
			success:     true,
			measurement: withoutEntries,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(insertMeasurementSQL)).
					WithArgs(testMeasurementID, testCustomerID, testMeasuredOn, testStaffID, 65, testStaffID, true, testCreatedAt, testCreatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:        "failure insert measurement error",
			success:     false,
			measurement: withoutEntries,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(insertMeasurementSQL)).
					WillReturnError(errors.New("insert measurement error"))
				mock.ExpectRollback()
			},
		},
		{
			name:        "failure insert entry error",
			success:     false,
			measurement: withEntries,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(insertMeasurementSQL)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(insertEntrySQL)).
					WillReturnError(errors.New("insert entry error"))
				mock.ExpectRollback()
			},
		},
		{
			name:        "failure insert value error",
			success:     false,
			measurement: withEntries,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(insertMeasurementSQL)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(insertEntrySQL)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(insertValuesSQL)).
					WillReturnError(errors.New("insert value error"))
				mock.ExpectRollback()
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

			repo := NewMeasurementRepository(gormDB)

			err = repo.Create(context.Background(), tt.measurement)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		success     bool
		wantErr     error
		wantEntries int
		wantValues  int
		wantNote    string
		setup       func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "success find measurement with entries and values",
			success:     true,
			wantEntries: 1,
			wantValues:  2,
			wantNote:    "ふらつきあり",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectMeasurementSQL)).
					WithArgs(testMeasurementID, 1).
					WillReturnRows(sqlmock.NewRows(measurementCols).
						AddRow(testMeasurementID, testCustomerID, testMeasuredOn, testStaffID, 65, testStaffID, false, testCreatedAt, testUpdatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(selectEntriesSQL)).
					WithArgs(testMeasurementID).
					WillReturnRows(sqlmock.NewRows(entryCols).
						AddRow(testMeasurementEntryID, testMeasurementID, testMeasurementItemID, false, "ふらつきあり", testCreatedAt, testUpdatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(selectValuesSQL)).
					WithArgs(testMeasurementEntryID).
					WillReturnRows(sqlmock.NewRows(valueCols).
						AddRow("55555555-5555-5555-5555-555555555555", testMeasurementEntryID, 1, "left", 32.4, nil, nil, testCreatedAt, testUpdatedAt).
						AddRow("66666666-6666-6666-6666-666666666666", testMeasurementEntryID, 1, "right", 33.1, nil, nil, testCreatedAt, testUpdatedAt))
			},
		},
		{
			name:    "success find measurement without entries",
			success: true,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectMeasurementSQL)).
					WithArgs(testMeasurementID, 1).
					WillReturnRows(sqlmock.NewRows(measurementCols).
						AddRow(testMeasurementID, testCustomerID, testMeasuredOn, testStaffID, 65, testStaffID, true, testCreatedAt, testUpdatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(selectEntriesSQL)).
					WithArgs(testMeasurementID).
					WillReturnRows(sqlmock.NewRows(entryCols))
			},
		},
		{
			name:    "failure measurement not found",
			success: false,
			wantErr: measurement.ErrMeasurementNotFound,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectMeasurementSQL)).
					WithArgs(testMeasurementID, 1).
					WillReturnRows(sqlmock.NewRows(measurementCols))
			},
		},
		{
			name:    "failure find measurement error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectMeasurementSQL)).
					WithArgs(testMeasurementID, 1).
					WillReturnError(errors.New("find measurement error"))
			},
		},
		{
			name:    "failure preload entries error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectMeasurementSQL)).
					WithArgs(testMeasurementID, 1).
					WillReturnRows(sqlmock.NewRows(measurementCols).
						AddRow(testMeasurementID, testCustomerID, testMeasuredOn, testStaffID, 65, testStaffID, false, testCreatedAt, testUpdatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(selectEntriesSQL)).
					WithArgs(testMeasurementID).
					WillReturnError(errors.New("preload entries error"))
			},
		},
		{
			name:    "failure preload values error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectMeasurementSQL)).
					WithArgs(testMeasurementID, 1).
					WillReturnRows(sqlmock.NewRows(measurementCols).
						AddRow(testMeasurementID, testCustomerID, testMeasuredOn, testStaffID, 65, testStaffID, false, testCreatedAt, testUpdatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(selectEntriesSQL)).
					WithArgs(testMeasurementID).
					WillReturnRows(sqlmock.NewRows(entryCols).
						AddRow(testMeasurementEntryID, testMeasurementID, testMeasurementItemID, false, nil, testCreatedAt, testUpdatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(selectValuesSQL)).
					WithArgs(testMeasurementEntryID).
					WillReturnError(errors.New("preload values error"))
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

			repo := NewMeasurementRepository(gormDB)

			measurementID, err := measurement.NewMeasurementIDFromString(testMeasurementID)
			if err != nil {
				t.Fatalf("failed to new measurement id: %v", err)
			}

			foundMeasurement, err := repo.FindByID(context.Background(), measurementID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.success {
				if foundMeasurement.ID().String() != testMeasurementID {
					t.Errorf("ID() = %v, want %v", foundMeasurement.ID().String(), testMeasurementID)
				}
				if foundMeasurement.CustomerID().String() != testCustomerID {
					t.Errorf("CustomerID() = %v, want %v", foundMeasurement.CustomerID().String(), testCustomerID)
				}
				if !foundMeasurement.MeasuredOn().Equal(testMeasuredOn) {
					t.Errorf("MeasuredOn() = %v, want %v", foundMeasurement.MeasuredOn(), testMeasuredOn)
				}
				if foundMeasurement.MeasuredBy().String() != testStaffID {
					t.Errorf("MeasuredBy() = %v, want %v", foundMeasurement.MeasuredBy().String(), testStaffID)
				}
				if foundMeasurement.AgeAtMeasurement().Int() != 65 {
					t.Errorf("AgeAtMeasurement() = %v, want %v", foundMeasurement.AgeAtMeasurement().Int(), 65)
				}
				if !foundMeasurement.CreatedAt().Equal(testCreatedAt) {
					t.Errorf("CreatedAt() = %v, want %v", foundMeasurement.CreatedAt(), testCreatedAt)
				}
				if !foundMeasurement.UpdatedAt().Equal(testUpdatedAt) {
					t.Errorf("UpdatedAt() = %v, want %v", foundMeasurement.UpdatedAt(), testUpdatedAt)
				}
				entries := foundMeasurement.Entries()
				if len(entries) != tt.wantEntries {
					t.Errorf("len(Entries()) = %v, want %v", len(entries), tt.wantEntries)
				}
				values := 0
				for _, entry := range entries {
					values += len(entry.Values())
					if entry.MeasurementItemID().String() != testMeasurementItemID {
						t.Errorf("MeasurementItemID() = %v, want %v", entry.MeasurementItemID().String(), testMeasurementItemID)
					}
					if tt.wantNote == "" && entry.Note() != nil {
						t.Errorf("Note() = %v, want nil", entry.Note())
					}
					if tt.wantNote != "" && (entry.Note() == nil || entry.Note().String() != tt.wantNote) {
						t.Errorf("Note() = %v, want %v", entry.Note(), tt.wantNote)
					}
					for _, value := range entry.Values() {
						if value.Value() == nil {
							t.Errorf("Value() = nil, want a value")
						}
						if value.ValueSecondary() != nil {
							t.Errorf("ValueSecondary() = %v, want nil", value.ValueSecondary())
						}
						if value.ValueChoice() != nil {
							t.Errorf("ValueChoice() = %v, want nil", value.ValueChoice())
						}
					}
				}
				if values != tt.wantValues {
					t.Errorf("values = %v, want %v", values, tt.wantValues)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListByCustomerID(t *testing.T) {
	t.Parallel()
	selectByCustomerIDSQL := `SELECT * FROM "measurements" WHERE customer_id = $1 ORDER BY measured_on DESC, id`
	selectEntriesInSQL := `SELECT * FROM "measurement_entries" WHERE "measurement_entries"."measurement_id" IN ($1,$2) ORDER BY measurement_entries.measurement_item_id`
	tests := []struct {
		name             string
		success          bool
		wantMeasurements int
		setup            func(mock sqlmock.Sqlmock)
	}{
		{
			name:             "success list measurements",
			success:          true,
			wantMeasurements: 2,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectByCustomerIDSQL)).
					WithArgs(testCustomerID).
					WillReturnRows(sqlmock.NewRows(measurementCols).
						AddRow(testMeasurementID, testCustomerID, testMeasuredOn, testStaffID, 65, testStaffID, false, testCreatedAt, testUpdatedAt).
						AddRow("77777777-7777-7777-7777-777777777777", testCustomerID, testMeasuredOn.AddDate(0, -1, 0), testStaffID, 65, testStaffID, false, testCreatedAt, testUpdatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(selectEntriesInSQL)).
					WithArgs(testMeasurementID, "77777777-7777-7777-7777-777777777777").
					WillReturnRows(sqlmock.NewRows(entryCols).
						AddRow(testMeasurementEntryID, testMeasurementID, testMeasurementItemID, false, nil, testCreatedAt, testUpdatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(selectValuesSQL)).
					WithArgs(testMeasurementEntryID).
					WillReturnRows(sqlmock.NewRows(valueCols).
						AddRow("55555555-5555-5555-5555-555555555555", testMeasurementEntryID, 1, "none", 72.0, nil, nil, testCreatedAt, testUpdatedAt))
			},
		},
		{
			name:             "success list no measurements",
			success:          true,
			wantMeasurements: 0,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectByCustomerIDSQL)).
					WithArgs(testCustomerID).
					WillReturnRows(sqlmock.NewRows(measurementCols))
			},
		},
		{
			name:    "failure list measurements error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectByCustomerIDSQL)).
					WithArgs(testCustomerID).
					WillReturnError(errors.New("list measurements error"))
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

			repo := NewMeasurementRepository(gormDB)

			customerID, err := customer.NewCustomerIDFromString(testCustomerID)
			if err != nil {
				t.Fatalf("failed to new customer id: %v", err)
			}

			measurements, err := repo.ListByCustomerID(context.Background(), customerID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && measurements != nil {
				t.Errorf("expected no measurements on failure, but got %v", len(measurements))
			}
			if tt.success && len(measurements) != tt.wantMeasurements {
				t.Errorf("len(measurements) = %v, want %v", len(measurements), tt.wantMeasurements)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListByCustomerIDs(t *testing.T) {
	t.Parallel()
	otherCustomerID := "88888888-8888-8888-8888-888888888888"
	otherMeasurementID := "99999999-9999-9999-9999-999999999999"
	selectByCustomerIDsSQL := `SELECT * FROM "measurements" WHERE customer_id IN ($1,$2) ORDER BY customer_id, measured_on DESC, id`
	selectEntriesInSQL := `SELECT * FROM "measurement_entries" WHERE "measurement_entries"."measurement_id" IN ($1,$2) ORDER BY measurement_entries.measurement_item_id`

	type measurementRow struct {
		measurementID string
		customerID    string
		age           int
		isDraft       bool
		entries       int
	}

	tests := []struct {
		name             string
		success          bool
		customerIDs      []string
		wantMeasurements []measurementRow
		setup            func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "success list measurements of multiple customers",
			success:     true,
			customerIDs: []string{testCustomerID, otherCustomerID},
			wantMeasurements: []measurementRow{
				{measurementID: testMeasurementID, customerID: testCustomerID, age: 65, isDraft: false, entries: 1},
				{measurementID: otherMeasurementID, customerID: otherCustomerID, age: 48, isDraft: true, entries: 0},
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectByCustomerIDsSQL)).
					WithArgs(testCustomerID, otherCustomerID).
					WillReturnRows(sqlmock.NewRows(measurementCols).
						AddRow(testMeasurementID, testCustomerID, testMeasuredOn, testStaffID, 65, testStaffID, false, testCreatedAt, testUpdatedAt).
						AddRow(otherMeasurementID, otherCustomerID, testMeasuredOn, testStaffID, 48, testStaffID, true, testCreatedAt, testUpdatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(selectEntriesInSQL)).
					WithArgs(testMeasurementID, otherMeasurementID).
					WillReturnRows(sqlmock.NewRows(entryCols).
						AddRow(testMeasurementEntryID, testMeasurementID, testMeasurementItemID, false, nil, testCreatedAt, testUpdatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(selectValuesSQL)).
					WithArgs(testMeasurementEntryID).
					WillReturnRows(sqlmock.NewRows(valueCols).
						AddRow("55555555-5555-5555-5555-555555555555", testMeasurementEntryID, 1, "none", 72.0, nil, nil, testCreatedAt, testUpdatedAt))
			},
		},
		{
			name:        "success list no measurements of customers without any measurement",
			success:     true,
			customerIDs: []string{testCustomerID, otherCustomerID},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectByCustomerIDsSQL)).
					WithArgs(testCustomerID, otherCustomerID).
					WillReturnRows(sqlmock.NewRows(measurementCols))
			},
		},
		{
			name:        "success list no measurements without customer ids",
			success:     true,
			customerIDs: []string{},
			setup:       func(mock sqlmock.Sqlmock) {},
		},
		{
			name:        "failure list measurements error",
			success:     false,
			customerIDs: []string{testCustomerID, otherCustomerID},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectByCustomerIDsSQL)).
					WithArgs(testCustomerID, otherCustomerID).
					WillReturnError(errors.New("list measurements error"))
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

			repo := NewMeasurementRepository(gormDB)

			customerIDs := make([]customer.CustomerID, 0, len(tt.customerIDs))
			for _, id := range tt.customerIDs {
				customerID, err := customer.NewCustomerIDFromString(id)
				if err != nil {
					t.Fatalf("failed to new customer id: %v", err)
				}
				customerIDs = append(customerIDs, customerID)
			}

			measurements, err := repo.ListByCustomerIDs(context.Background(), customerIDs)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && measurements != nil {
				t.Errorf("expected no measurements on failure, but got %v", len(measurements))
			}
			if tt.success {
				if len(measurements) != len(tt.wantMeasurements) {
					t.Fatalf("len(measurements) = %v, want %v", len(measurements), len(tt.wantMeasurements))
				}
				for i, want := range tt.wantMeasurements {
					m := measurements[i]
					if m.ID().String() != want.measurementID {
						t.Errorf("measurements[%d].ID() = %v, want %v", i, m.ID(), want.measurementID)
					}
					if m.CustomerID().String() != want.customerID {
						t.Errorf("measurements[%d].CustomerID() = %v, want %v", i, m.CustomerID(), want.customerID)
					}
					if m.AgeAtMeasurement().Int() != want.age {
						t.Errorf("measurements[%d].AgeAtMeasurement() = %v, want %v", i, m.AgeAtMeasurement().Int(), want.age)
					}
					if m.IsDraft() != want.isDraft {
						t.Errorf("measurements[%d].IsDraft() = %v, want %v", i, m.IsDraft(), want.isDraft)
					}
					if len(m.Entries()) != want.entries {
						t.Errorf("len(measurements[%d].Entries()) = %v, want %v", i, len(m.Entries()), want.entries)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	updateMeasurementSQL := `UPDATE "measurements" SET "measured_on"=$1,"measured_by"=$2,"age_at_measurement"=$3,"updated_by"=$4,"is_draft"=$5,"updated_at"=$6 WHERE id = $7`
	deleteEntriesSQL := `DELETE FROM "measurement_entries" WHERE measurement_id = $1`

	measurementID, _ := measurement.NewMeasurementIDFromString(testMeasurementID)
	customerID, _ := customer.NewCustomerIDFromString(testCustomerID)
	measurementItemID, _ := measurementitem.NewMeasurementItemIDFromString(testMeasurementItemID)
	measuredOn, _ := measurement.NewMeasuredOn(2026, 8, 1)
	measuredBy, _ := staff.NewStaffID(testStaffID)
	ageAtMeasurement, _ := measurement.NewAgeAtMeasurement(65)
	trialIndex, _ := measurement.NewTrialIndex(1)
	left, _ := measurement.NewValue(32.4)
	right, _ := measurement.NewValue(33.1)
	note, _ := measurement.NewNote("ふらつきあり")
	entry := measurement.ReconstructMeasurementEntry(measurementItemID, false, note, []measurement.MeasurementValue{
		measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &left, nil, nil),
		measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &right, nil, nil),
	})
	withEntries := measurement.ReconstructMeasurement(measurementID, customerID, measuredOn, measuredBy, ageAtMeasurement, measuredBy, false, []measurement.MeasurementEntry{entry}, testCreatedAt, testCreatedAt)
	withoutEntries := measurement.ReconstructMeasurement(measurementID, customerID, measuredOn, measuredBy, ageAtMeasurement, measuredBy, false, nil, testCreatedAt, testCreatedAt)

	tests := []struct {
		name        string
		success     bool
		wantErr     error
		measurement measurement.Measurement
		setup       func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "success update measurement replaces entries and values",
			success:     true,
			measurement: withEntries,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(updateMeasurementSQL)).
					WithArgs(testMeasuredOn, testStaffID, 65, testStaffID, false, testCreatedAt, testMeasurementID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(deleteEntriesSQL)).
					WithArgs(testMeasurementID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(insertEntrySQL)).
					WithArgs(sqlmock.AnyArg(), testMeasurementID, testMeasurementItemID, false, "ふらつきあり", testCreatedAt, testCreatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(insertValuesSQL)).
					WithArgs(
						sqlmock.AnyArg(), sqlmock.AnyArg(), 1, "left", 32.4, nil, nil, testCreatedAt, testCreatedAt,
						sqlmock.AnyArg(), sqlmock.AnyArg(), 1, "right", 33.1, nil, nil, testCreatedAt, testCreatedAt,
					).
					WillReturnResult(sqlmock.NewResult(1, 2))
				mock.ExpectCommit()
			},
		},
		{
			name:        "success update measurement without entries",
			success:     true,
			measurement: withoutEntries,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(updateMeasurementSQL)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(deleteEntriesSQL)).
					WithArgs(testMeasurementID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:        "failure measurement not found",
			success:     false,
			wantErr:     measurement.ErrMeasurementNotFound,
			measurement: withoutEntries,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(updateMeasurementSQL)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
		},
		{
			name:        "failure update measurement error",
			success:     false,
			measurement: withoutEntries,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(updateMeasurementSQL)).
					WillReturnError(errors.New("update measurement error"))
				mock.ExpectRollback()
			},
		},
		{
			name:        "failure delete entries error",
			success:     false,
			measurement: withoutEntries,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(updateMeasurementSQL)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(deleteEntriesSQL)).
					WillReturnError(errors.New("delete entries error"))
				mock.ExpectRollback()
			},
		},
		{
			name:        "failure insert entry error",
			success:     false,
			measurement: withEntries,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(updateMeasurementSQL)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(deleteEntriesSQL)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(insertEntrySQL)).
					WillReturnError(errors.New("insert entry error"))
				mock.ExpectRollback()
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

			repo := NewMeasurementRepository(gormDB)

			err = repo.Update(context.Background(), tt.measurement)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	deleteMeasurementSQL := `DELETE FROM "measurements" WHERE id = $1`
	tests := []struct {
		name    string
		success bool
		wantErr error
		setup   func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "success delete measurement",
			success: true,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(deleteMeasurementSQL)).
					WithArgs(testMeasurementID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure measurement not found",
			success: false,
			wantErr: measurement.ErrMeasurementNotFound,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(deleteMeasurementSQL)).
					WithArgs(testMeasurementID).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
		},
		{
			name:    "failure delete measurement error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(deleteMeasurementSQL)).
					WithArgs(testMeasurementID).
					WillReturnError(errors.New("delete measurement error"))
				mock.ExpectRollback()
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

			repo := NewMeasurementRepository(gormDB)

			measurementID, err := measurement.NewMeasurementIDFromString(testMeasurementID)
			if err != nil {
				t.Fatalf("failed to new measurement id: %v", err)
			}

			err = repo.Delete(context.Background(), measurementID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
