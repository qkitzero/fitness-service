package judgment

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/judgment"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	mocksjudgment "github.com/qkitzero/fitness-service/mocks/domain/judgment"
)

var (
	testCreatedAt = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	testUpdatedAt = time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
)

func TestFindByMeasurementID(t *testing.T) {
	t.Parallel()
	columns := []string{"measurement_id", "advice", "created_at", "updated_at"}
	tests := []struct {
		name       string
		success    bool
		wantErr    error
		setup      func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID)
		wantAdvice string
	}{
		{
			name:    "success find judgment by measurement id",
			success: true,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				judgmentRows := sqlmock.NewRows(columns).
					AddRow(measurementID, "週2回のスクワットを継続してください", testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "judgments" WHERE measurement_id = $1 ORDER BY "judgments"."measurement_id" LIMIT $2`)).
					WithArgs(measurementID, 1).
					WillReturnRows(judgmentRows)
			},
			wantAdvice: "週2回のスクワットを継続してください",
		},
		{
			name:    "success find judgment without advice",
			success: true,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				judgmentRows := sqlmock.NewRows(columns).
					AddRow(measurementID, nil, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "judgments" WHERE measurement_id = $1 ORDER BY "judgments"."measurement_id" LIMIT $2`)).
					WithArgs(measurementID, 1).
					WillReturnRows(judgmentRows)
			},
		},
		{
			name:    "failure judgment not found",
			success: false,
			wantErr: judgment.ErrJudgmentNotFound,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "judgments" WHERE measurement_id = $1 ORDER BY "judgments"."measurement_id" LIMIT $2`)).
					WithArgs(measurementID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "failure find judgment error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "judgments" WHERE measurement_id = $1 ORDER BY "judgments"."measurement_id" LIMIT $2`)).
					WithArgs(measurementID, 1).
					WillReturnError(errors.New("find judgment error"))
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

			repo := NewJudgmentRepository(gormDB)

			j, err := repo.FindByMeasurementID(context.Background(), measurementID)
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
				if j.MeasurementID() != measurementID {
					t.Errorf("MeasurementID() = %v, want %v", j.MeasurementID(), measurementID)
				}
				if tt.wantAdvice == "" && j.Advice() != nil {
					t.Errorf("Advice() = %v, want nil", j.Advice())
				}
				if tt.wantAdvice != "" && (j.Advice() == nil || j.Advice().String() != tt.wantAdvice) {
					t.Errorf("Advice() = %v, want %v", j.Advice(), tt.wantAdvice)
				}
				if !j.CreatedAt().Equal(testCreatedAt) || !j.UpdatedAt().Equal(testUpdatedAt) {
					t.Errorf("timestamps = %v/%v, want %v/%v", j.CreatedAt(), j.UpdatedAt(), testCreatedAt, testUpdatedAt)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpsert(t *testing.T) {
	t.Parallel()
	advice, _ := judgment.NewAdvice("週2回のスクワットを継続してください")

	tests := []struct {
		name    string
		success bool
		advice  *judgment.Advice
		setup   func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID)
	}{
		{
			name:    "success insert judgment",
			success: true,
			advice:  advice,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "judgments" ("measurement_id","advice","created_at","updated_at") VALUES ($1,$2,$3,$4) ON CONFLICT ("measurement_id") DO UPDATE SET "advice"="excluded"."advice","updated_at"="excluded"."updated_at"`)).
					WithArgs(measurementID, "週2回のスクワットを継続してください", testCreatedAt, testUpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "success update judgment",
			success: true,
			advice:  advice,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "judgments" ("measurement_id","advice","created_at","updated_at") VALUES ($1,$2,$3,$4) ON CONFLICT ("measurement_id") DO UPDATE SET "advice"="excluded"."advice","updated_at"="excluded"."updated_at"`)).
					WithArgs(measurementID, "週2回のスクワットを継続してください", testCreatedAt, testUpdatedAt).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "success upsert judgment without advice",
			success: true,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "judgments" ("measurement_id","advice","created_at","updated_at") VALUES ($1,$2,$3,$4) ON CONFLICT ("measurement_id") DO UPDATE SET "advice"="excluded"."advice","updated_at"="excluded"."updated_at"`)).
					WithArgs(measurementID, nil, testCreatedAt, testUpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "failure upsert judgment error",
			success: false,
			advice:  advice,
			setup: func(mock sqlmock.Sqlmock, measurementID measurement.MeasurementID) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "judgments" ("measurement_id","advice","created_at","updated_at") VALUES ($1,$2,$3,$4) ON CONFLICT ("measurement_id") DO UPDATE SET "advice"="excluded"."advice","updated_at"="excluded"."updated_at"`)).
					WithArgs(measurementID, "週2回のスクワットを継続してください", testCreatedAt, testUpdatedAt).
					WillReturnError(errors.New("upsert judgment error"))

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

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			measurementID := measurement.NewMeasurementID()

			mockJudgment := mocksjudgment.NewMockJudgment(ctrl)
			mockJudgment.EXPECT().MeasurementID().Return(measurementID).AnyTimes()
			mockJudgment.EXPECT().Advice().Return(tt.advice).AnyTimes()
			mockJudgment.EXPECT().CreatedAt().Return(testCreatedAt).AnyTimes()
			mockJudgment.EXPECT().UpdatedAt().Return(testUpdatedAt).AnyTimes()

			tt.setup(mock, measurementID)

			repo := NewJudgmentRepository(gormDB)

			err = repo.Upsert(context.Background(), mockJudgment)
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
