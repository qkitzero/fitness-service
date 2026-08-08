package judgment

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/judgment"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/training"
)

func TestListByMeasurementID(t *testing.T) {
	t.Parallel()
	columns := []string{"id", "measurement_id", "sort_order", "element", "part", "training_menu_id", "amount", "unit", "sets", "created_at", "updated_at"}
	muscleStrength := measurementitem.ElementMuscleStrength
	upperLimb := training.PartUpperLimb

	tests := []struct {
		name     string
		success  bool
		setup    func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID, trainingMenuID training.TrainingMenuID)
		wantLen  int
		wantSort []int
	}{
		{
			name:    "success list prescribed menu overrides by measurement id",
			success: true,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID, trainingMenuID training.TrainingMenuID) {
				rows := sqlmock.NewRows(columns).
					AddRow(judgment.NewPrescribedMenuOverrideID(), measurementID, 1, muscleStrength, upperLimb, trainingMenuID, 10, training.UnitReps, 3, testCreatedAt, testUpdatedAt).
					AddRow(judgment.NewPrescribedMenuOverrideID(), measurementID, 2, nil, nil, trainingMenuID, 20, training.UnitMinutes, 1, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "prescribed_menu_overrides" WHERE measurement_id = $1 ORDER BY sort_order`)).
					WithArgs(measurementID).
					WillReturnRows(rows)
			},
			wantLen:  2,
			wantSort: []int{1, 2},
		},
		{
			name:    "success list no prescribed menu override",
			success: true,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID, trainingMenuID training.TrainingMenuID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "prescribed_menu_overrides" WHERE measurement_id = $1 ORDER BY sort_order`)).
					WithArgs(measurementID).
					WillReturnRows(sqlmock.NewRows(columns))
			},
		},
		{
			name:    "failure list prescribed menu overrides error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID, trainingMenuID training.TrainingMenuID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "prescribed_menu_overrides" WHERE measurement_id = $1 ORDER BY sort_order`)).
					WithArgs(measurementID).
					WillReturnError(errors.New("list prescribed menu overrides error"))
			},
		},
		{
			name:    "failure invalid stored unit",
			success: false,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID, trainingMenuID training.TrainingMenuID) {
				rows := sqlmock.NewRows(columns).
					AddRow(judgment.NewPrescribedMenuOverrideID(), measurementID, 1, nil, nil, trainingMenuID, 10, training.Unit("hours"), 3, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "prescribed_menu_overrides" WHERE measurement_id = $1 ORDER BY sort_order`)).
					WithArgs(measurementID).
					WillReturnRows(rows)
			},
		},
		{
			name:    "failure invalid stored element",
			success: false,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID, trainingMenuID training.TrainingMenuID) {
				rows := sqlmock.NewRows(columns).
					AddRow(judgment.NewPrescribedMenuOverrideID(), measurementID, 1, measurementitem.Element("unknown"), upperLimb, trainingMenuID, 10, training.UnitReps, 3, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "prescribed_menu_overrides" WHERE measurement_id = $1 ORDER BY sort_order`)).
					WithArgs(measurementID).
					WillReturnRows(rows)
			},
		},
		{
			name:    "failure invalid stored part",
			success: false,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID, trainingMenuID training.TrainingMenuID) {
				rows := sqlmock.NewRows(columns).
					AddRow(judgment.NewPrescribedMenuOverrideID(), measurementID, 1, muscleStrength, training.Part("unknown"), trainingMenuID, 10, training.UnitReps, 3, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "prescribed_menu_overrides" WHERE measurement_id = $1 ORDER BY sort_order`)).
					WithArgs(measurementID).
					WillReturnRows(rows)
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

			measurementID := measurement.NewMeasurementID()
			trainingMenuID := training.NewTrainingMenuID()

			tt.setup(mock, measurementID, trainingMenuID)

			repo := NewPrescribedMenuOverrideRepository(gormDB)

			overrides, err := repo.ListByMeasurementID(context.Background(), measurementID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success {
				if len(overrides) != tt.wantLen {
					t.Fatalf("len(overrides) = %v, want %v", len(overrides), tt.wantLen)
				}
				for i, wantSortOrder := range tt.wantSort {
					if overrides[i].SortOrder().Int() != wantSortOrder {
						t.Errorf("overrides[%d].SortOrder() = %v, want %v", i, overrides[i].SortOrder().Int(), wantSortOrder)
					}
					if overrides[i].MeasurementID() != measurementID {
						t.Errorf("overrides[%d].MeasurementID() = %v, want %v", i, overrides[i].MeasurementID(), measurementID)
					}
					if overrides[i].TrainingMenuID() != trainingMenuID {
						t.Errorf("overrides[%d].TrainingMenuID() = %v, want %v", i, overrides[i].TrainingMenuID(), trainingMenuID)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestReplaceByMeasurementID(t *testing.T) {
	t.Parallel()
	muscleStrength := measurementitem.ElementMuscleStrength
	upperLimb := training.PartUpperLimb

	tests := []struct {
		name      string
		success   bool
		sortOrder []int
		setup     func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID)
	}{
		{
			name:      "success replace prescribed menu overrides",
			success:   true,
			sortOrder: []int{1, 2},
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM measurements WHERE id = $1 FOR UPDATE`)).
					WithArgs(measurementID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(measurementID))
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "prescribed_menu_overrides" WHERE measurement_id = $1`)).
					WithArgs(measurementID).
					WillReturnResult(sqlmock.NewResult(0, 3))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "prescribed_menu_overrides" ("id","measurement_id","sort_order","element","part","training_menu_id","amount","unit","sets","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11),($12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)`)).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectCommit()
			},
		},
		{
			name:    "success replace prescribed menu overrides with an empty set",
			success: true,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM measurements WHERE id = $1 FOR UPDATE`)).
					WithArgs(measurementID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(measurementID))
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "prescribed_menu_overrides" WHERE measurement_id = $1`)).
					WithArgs(measurementID).
					WillReturnResult(sqlmock.NewResult(0, 3))
				mock.ExpectCommit()
			},
		},
		{
			name:      "failure delete prescribed menu overrides error",
			success:   false,
			sortOrder: []int{1},
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM measurements WHERE id = $1 FOR UPDATE`)).
					WithArgs(measurementID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(measurementID))
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "prescribed_menu_overrides" WHERE measurement_id = $1`)).
					WithArgs(measurementID).
					WillReturnError(errors.New("delete prescribed menu overrides error"))
				mock.ExpectRollback()
			},
		},
		{
			name:      "failure insert prescribed menu overrides error",
			success:   false,
			sortOrder: []int{1},
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM measurements WHERE id = $1 FOR UPDATE`)).
					WithArgs(measurementID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(measurementID))
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "prescribed_menu_overrides" WHERE measurement_id = $1`)).
					WithArgs(measurementID).
					WillReturnResult(sqlmock.NewResult(0, 3))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "prescribed_menu_overrides" ("id","measurement_id","sort_order","element","part","training_menu_id","amount","unit","sets","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`)).
					WillReturnError(errors.New("insert prescribed menu overrides error"))
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

			measurementID := measurement.NewMeasurementID()
			amount, _ := training.NewAmount(10)
			sets, _ := training.NewSets(3)

			overrides := make([]judgment.PrescribedMenuOverride, 0, len(tt.sortOrder))
			for _, n := range tt.sortOrder {
				sortOrder, _ := training.NewSortOrder(n)
				overrides = append(overrides, judgment.NewPrescribedMenuOverride(judgment.NewPrescribedMenuOverrideID(), measurementID, sortOrder, &muscleStrength, &upperLimb, training.NewTrainingMenuID(), amount, training.UnitReps, sets, testCreatedAt, testUpdatedAt))
			}

			tt.setup(mock, measurementID)

			repo := NewPrescribedMenuOverrideRepository(gormDB)

			err = repo.ReplaceByMeasurementID(context.Background(), measurementID, overrides)
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

func TestDeleteByMeasurementID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID)
	}{
		{
			name:    "success delete prescribed menu overrides",
			success: true,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "prescribed_menu_overrides" WHERE measurement_id = $1`)).
					WithArgs(measurementID).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectCommit()
			},
		},
		{
			name:    "success delete prescribed menu overrides that do not exist",
			success: true,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "prescribed_menu_overrides" WHERE measurement_id = $1`)).
					WithArgs(measurementID).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure delete prescribed menu overrides error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "prescribed_menu_overrides" WHERE measurement_id = $1`)).
					WithArgs(measurementID).
					WillReturnError(errors.New("delete prescribed menu overrides error"))
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

			measurementID := measurement.NewMeasurementID()

			tt.setup(mock, measurementID)

			repo := NewPrescribedMenuOverrideRepository(gormDB)

			err = repo.DeleteByMeasurementID(context.Background(), measurementID)
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
