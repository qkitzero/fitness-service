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
)

const listRankStandardsSQL = `SELECT * FROM "rank_standards"`

var rankStandardColumns = []string{"rank", "z_score_min", "z_score_max", "created_at", "updated_at"}

func TestList(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	zScore15 := 1.5
	zScore05 := 0.5
	zScoreMinus05 := -0.5
	zScoreMinus15 := -1.5
	tests := []struct {
		name          string
		success       bool
		wantRanks     []string
		wantZScoreMin []*float64
		wantZScoreMax []*float64
		setup         func(mock sqlmock.Sqlmock)
	}{
		{
			name:          "success list rank standards in rank order",
			success:       true,
			wantRanks:     []string{"A", "B", "C", "D", "E"},
			wantZScoreMin: []*float64{&zScore15, &zScore05, &zScoreMinus05, &zScoreMinus15, nil},
			wantZScoreMax: []*float64{nil, &zScore15, &zScore05, &zScoreMinus05, &zScoreMinus15},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listRankStandardsSQL)).
					WillReturnRows(sqlmock.NewRows(rankStandardColumns).
						AddRow("C", -0.5, 0.5, createdAt, updatedAt).
						AddRow("E", nil, -1.5, createdAt, updatedAt).
						AddRow("A", 1.5, nil, createdAt, updatedAt).
						AddRow("D", -1.5, -0.5, createdAt, updatedAt).
						AddRow("B", 0.5, 1.5, createdAt, updatedAt))
			},
		},
		{
			name:          "success list no rank standards",
			success:       true,
			wantRanks:     []string{},
			wantZScoreMin: []*float64{},
			wantZScoreMax: []*float64{},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listRankStandardsSQL)).
					WillReturnRows(sqlmock.NewRows(rankStandardColumns))
			},
		},
		{
			name:    "failure list rank standards error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listRankStandardsSQL)).
					WillReturnError(errors.New("list rank standards error"))
			},
		},
		{
			name:    "failure unknown rank",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listRankStandardsSQL)).
					WillReturnRows(sqlmock.NewRows(rankStandardColumns).
						AddRow("F", -0.5, 0.5, createdAt, updatedAt))
			},
		},
		{
			name:    "failure z score below the lower bound",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listRankStandardsSQL)).
					WillReturnRows(sqlmock.NewRows(rankStandardColumns).
						AddRow("C", -100.0, 0.5, createdAt, updatedAt))
			},
		},
		{
			name:    "failure z score above the upper bound",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listRankStandardsSQL)).
					WillReturnRows(sqlmock.NewRows(rankStandardColumns).
						AddRow("C", -0.5, 100.0, createdAt, updatedAt))
			},
		},
		{
			name:    "failure reversed z score bounds",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listRankStandardsSQL)).
					WillReturnRows(sqlmock.NewRows(rankStandardColumns).
						AddRow("C", 0.5, -0.5, createdAt, updatedAt))
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

			repo := NewRankStandardRepository(gormDB)

			rankStandards, err := repo.List(context.Background())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && rankStandards != nil {
				t.Errorf("expected no rank standards on failure, but got %v", len(rankStandards))
			}
			if tt.success {
				if len(rankStandards) != len(tt.wantRanks) {
					t.Errorf("len(rankStandards) = %v, want %v", len(rankStandards), len(tt.wantRanks))
				}
				for i := range tt.wantRanks {
					if i >= len(rankStandards) {
						break
					}
					r := rankStandards[i]
					if r.Rank().String() != tt.wantRanks[i] {
						t.Errorf("rankStandards[%d].Rank() = %v, want %v", i, r.Rank().String(), tt.wantRanks[i])
					}
					switch {
					case tt.wantZScoreMin[i] == nil:
						if r.ZScoreMin() != nil {
							t.Errorf("rankStandards[%d].ZScoreMin() = %v, want nil", i, r.ZScoreMin())
						}
					case r.ZScoreMin() == nil:
						t.Errorf("rankStandards[%d].ZScoreMin() = nil, want %v", i, *tt.wantZScoreMin[i])
					case r.ZScoreMin().Float64() != *tt.wantZScoreMin[i]:
						t.Errorf("rankStandards[%d].ZScoreMin() = %v, want %v", i, r.ZScoreMin().Float64(), *tt.wantZScoreMin[i])
					}
					switch {
					case tt.wantZScoreMax[i] == nil:
						if r.ZScoreMax() != nil {
							t.Errorf("rankStandards[%d].ZScoreMax() = %v, want nil", i, r.ZScoreMax())
						}
					case r.ZScoreMax() == nil:
						t.Errorf("rankStandards[%d].ZScoreMax() = nil, want %v", i, *tt.wantZScoreMax[i])
					case r.ZScoreMax().Float64() != *tt.wantZScoreMax[i]:
						t.Errorf("rankStandards[%d].ZScoreMax() = %v, want %v", i, r.ZScoreMax().Float64(), *tt.wantZScoreMax[i])
					}
					if !r.CreatedAt().Equal(createdAt) {
						t.Errorf("rankStandards[%d].CreatedAt() = %v, want %v", i, r.CreatedAt(), createdAt)
					}
					if !r.UpdatedAt().Equal(updatedAt) {
						t.Errorf("rankStandards[%d].UpdatedAt() = %v, want %v", i, r.UpdatedAt(), updatedAt)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
