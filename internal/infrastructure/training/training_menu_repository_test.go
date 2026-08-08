package training

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/training"
)

const (
	listSQL               = `SELECT * FROM "training_menus"`
	findByIDsSQL          = `SELECT * FROM "training_menus" WHERE id IN ($1,$2)`
	findByThreeIDsSQL     = `SELECT * FROM "training_menus" WHERE id IN ($1,$2,$3)`
	armRaiseID            = "3d1d7e6a-3f6a-4c0e-8a6d-2a4f26b5c101"
	wallPushUpID          = "9e0f4c22-19a1-4f7d-8f1b-6de1c9a53202"
	squatID               = "60da87f0-3f50-4c42-8801-0313a3be3e38"
	frontPlankID          = "b1d0a1c7-5c2a-4e91-9d13-7f4b52e6c404"
	shoulderRotationID    = "c7a2f8b3-2e64-4d15-8c72-1b9d3f6a7505"
	indoorWalkID          = "e4f6b9d1-7c38-4a26-9b54-8d2e5f1a3606"
	unknownTrainingMenuID = "00000000-0000-0000-0000-000000000000"
	wallPushUpInstruction = "壁から一歩離れて立ち、肩幅で壁に手をつく。"
	armRaiseInstruction   = "椅子に座り、両腕を前から頭の高さまでゆっくり上げて下ろす。"
)

var trainingMenuColumns = []string{"id", "code", "name", "element", "part", "amount", "unit", "sets", "instruction", "created_at", "updated_at"}

func TestList(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	tests := []struct {
		name             string
		success          bool
		wantIDs          []string
		wantCodes        []string
		wantNames        []string
		wantElements     []string
		wantParts        []string
		wantAmounts      []int
		wantUnits        []string
		wantSets         []int
		wantInstructions []string
		setup            func(mock sqlmock.Sqlmock)
	}{
		{
			name:             "success list training menus in element order then part order then code order",
			success:          true,
			wantIDs:          []string{armRaiseID, wallPushUpID, squatID, frontPlankID, shoulderRotationID, indoorWalkID},
			wantCodes:        []string{"arm_raise", "wall_push_up", "squat", "front_plank", "shoulder_rotation", "indoor_walk"},
			wantNames:        []string{"腕上げ", "壁押し", "スクワット", "フロントプランク", "肩回し", "室内歩行"},
			wantElements:     []string{"muscle_strength", "muscle_strength", "muscle_strength", "muscle_strength", "flexibility", "mobility"},
			wantParts:        []string{"upper_limb", "upper_limb", "lower_limb", "whole_body", "upper_limb", "whole_body"},
			wantAmounts:      []int{10, 10, 10, 30, 10, 5},
			wantUnits:        []string{"reps", "reps", "reps", "seconds", "reps", "minutes"},
			wantSets:         []int{2, 3, 3, 3, 2, 1},
			wantInstructions: []string{armRaiseInstruction, wallPushUpInstruction, "足を肩幅に開き、腰を落とす。", "肘とつま先で体を支える。", "肘で大きな円を描く。", "無理のない速さで歩く。"},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns).
						AddRow(indoorWalkID, "indoor_walk", "室内歩行", "mobility", "whole_body", 5, "minutes", 1, "無理のない速さで歩く。", createdAt, updatedAt).
						AddRow(shoulderRotationID, "shoulder_rotation", "肩回し", "flexibility", "upper_limb", 10, "reps", 2, "肘で大きな円を描く。", createdAt, updatedAt).
						AddRow(frontPlankID, "front_plank", "フロントプランク", "muscle_strength", "whole_body", 30, "seconds", 3, "肘とつま先で体を支える。", createdAt, updatedAt).
						AddRow(squatID, "squat", "スクワット", "muscle_strength", "lower_limb", 10, "reps", 3, "足を肩幅に開き、腰を落とす。", createdAt, updatedAt).
						AddRow(wallPushUpID, "wall_push_up", "壁押し", "muscle_strength", "upper_limb", 10, "reps", 3, wallPushUpInstruction, createdAt, updatedAt).
						AddRow(armRaiseID, "arm_raise", "腕上げ", "muscle_strength", "upper_limb", 10, "reps", 2, armRaiseInstruction, createdAt, updatedAt))
			},
		},
		{
			name:             "success list no training menus",
			success:          true,
			wantIDs:          []string{},
			wantCodes:        []string{},
			wantNames:        []string{},
			wantElements:     []string{},
			wantParts:        []string{},
			wantAmounts:      []int{},
			wantUnits:        []string{},
			wantSets:         []int{},
			wantInstructions: []string{},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns))
			},
		},
		{
			name:    "failure list training menus error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnError(errors.New("list training menus error"))
			},
		},
		{
			name:    "failure invalid code",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns).
						AddRow(squatID, "Squat", "スクワット", "muscle_strength", "lower_limb", 10, "reps", 3, "足を肩幅に開き、腰を落とす。", createdAt, updatedAt))
			},
		},
		{
			name:    "failure empty name",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns).
						AddRow(squatID, "squat", "", "muscle_strength", "lower_limb", 10, "reps", 3, "足を肩幅に開き、腰を落とす。", createdAt, updatedAt))
			},
		},
		{
			name:    "failure unknown element",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns).
						AddRow(squatID, "squat", "スクワット", "explosive_power", "lower_limb", 10, "reps", 3, "足を肩幅に開き、腰を落とす。", createdAt, updatedAt))
			},
		},
		{
			name:    "failure unknown part",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns).
						AddRow(squatID, "squat", "スクワット", "muscle_strength", "trunk", 10, "reps", 3, "足を肩幅に開き、腰を落とす。", createdAt, updatedAt))
			},
		},
		{
			name:    "failure amount above the upper bound",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns).
						AddRow(squatID, "squat", "スクワット", "muscle_strength", "lower_limb", 1000, "reps", 3, "足を肩幅に開き、腰を落とす。", createdAt, updatedAt))
			},
		},
		{
			name:    "failure unknown unit",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns).
						AddRow(squatID, "squat", "スクワット", "muscle_strength", "lower_limb", 10, "hours", 3, "足を肩幅に開き、腰を落とす。", createdAt, updatedAt))
			},
		},
		{
			name:    "failure sets below the lower bound",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns).
						AddRow(squatID, "squat", "スクワット", "muscle_strength", "lower_limb", 10, "reps", 0, "足を肩幅に開き、腰を落とす。", createdAt, updatedAt))
			},
		},
		{
			name:    "failure empty instruction",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listSQL)).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns).
						AddRow(squatID, "squat", "スクワット", "muscle_strength", "lower_limb", 10, "reps", 3, "", createdAt, updatedAt))
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

			repo := NewTrainingMenuRepository(gormDB)

			trainingMenus, err := repo.List(context.Background())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && trainingMenus != nil {
				t.Errorf("expected no training menus on failure, but got %v", len(trainingMenus))
			}
			if tt.success {
				if len(trainingMenus) != len(tt.wantCodes) {
					t.Errorf("len(trainingMenus) = %v, want %v", len(trainingMenus), len(tt.wantCodes))
				}
				for i := range tt.wantCodes {
					if i >= len(trainingMenus) {
						break
					}
					m := trainingMenus[i]
					if m.ID().String() != tt.wantIDs[i] {
						t.Errorf("trainingMenus[%d].ID() = %v, want %v", i, m.ID().String(), tt.wantIDs[i])
					}
					if m.Code().String() != tt.wantCodes[i] {
						t.Errorf("trainingMenus[%d].Code() = %v, want %v", i, m.Code().String(), tt.wantCodes[i])
					}
					if m.Name().String() != tt.wantNames[i] {
						t.Errorf("trainingMenus[%d].Name() = %v, want %v", i, m.Name().String(), tt.wantNames[i])
					}
					if m.Element().String() != tt.wantElements[i] {
						t.Errorf("trainingMenus[%d].Element() = %v, want %v", i, m.Element().String(), tt.wantElements[i])
					}
					if m.Part().String() != tt.wantParts[i] {
						t.Errorf("trainingMenus[%d].Part() = %v, want %v", i, m.Part().String(), tt.wantParts[i])
					}
					if m.Amount().Int() != tt.wantAmounts[i] {
						t.Errorf("trainingMenus[%d].Amount() = %v, want %v", i, m.Amount().Int(), tt.wantAmounts[i])
					}
					if m.Unit().String() != tt.wantUnits[i] {
						t.Errorf("trainingMenus[%d].Unit() = %v, want %v", i, m.Unit().String(), tt.wantUnits[i])
					}
					if m.Sets().Int() != tt.wantSets[i] {
						t.Errorf("trainingMenus[%d].Sets() = %v, want %v", i, m.Sets().Int(), tt.wantSets[i])
					}
					if m.Instruction().String() != tt.wantInstructions[i] {
						t.Errorf("trainingMenus[%d].Instruction() = %v, want %v", i, m.Instruction().String(), tt.wantInstructions[i])
					}
					if !m.CreatedAt().Equal(createdAt) {
						t.Errorf("trainingMenus[%d].CreatedAt() = %v, want %v", i, m.CreatedAt(), createdAt)
					}
					if !m.UpdatedAt().Equal(updatedAt) {
						t.Errorf("trainingMenus[%d].UpdatedAt() = %v, want %v", i, m.UpdatedAt(), updatedAt)
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
		name      string
		success   bool
		ids       []string
		wantCodes []string
		setup     func(mock sqlmock.Sqlmock)
	}{
		{
			name:      "success find training menus by ids",
			success:   true,
			ids:       []string{squatID, wallPushUpID},
			wantCodes: []string{"squat", "wall_push_up"},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(findByIDsSQL)).
					WithArgs(squatID, wallPushUpID).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns).
						AddRow(squatID, "squat", "スクワット", "muscle_strength", "lower_limb", 10, "reps", 3, "足を肩幅に開き、腰を落とす。", createdAt, updatedAt).
						AddRow(wallPushUpID, "wall_push_up", "壁押し", "muscle_strength", "upper_limb", 10, "reps", 3, wallPushUpInstruction, createdAt, updatedAt))
			},
		},
		{
			name:      "success find no training menus for an empty id slice",
			success:   true,
			ids:       []string{},
			wantCodes: []string{},
			setup:     func(_ sqlmock.Sqlmock) {},
		},
		{
			name:      "success skip ids that do not exist",
			success:   true,
			ids:       []string{squatID, unknownTrainingMenuID, wallPushUpID},
			wantCodes: []string{"squat", "wall_push_up"},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(findByThreeIDsSQL)).
					WithArgs(squatID, unknownTrainingMenuID, wallPushUpID).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns).
						AddRow(squatID, "squat", "スクワット", "muscle_strength", "lower_limb", 10, "reps", 3, "足を肩幅に開き、腰を落とす。", createdAt, updatedAt).
						AddRow(wallPushUpID, "wall_push_up", "壁押し", "muscle_strength", "upper_limb", 10, "reps", 3, wallPushUpInstruction, createdAt, updatedAt))
			},
		},
		{
			name:    "failure find training menus error",
			success: false,
			ids:     []string{squatID, wallPushUpID},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(findByIDsSQL)).
					WithArgs(squatID, wallPushUpID).
					WillReturnError(errors.New("find training menus error"))
			},
		},
		{
			name:    "failure unknown element",
			success: false,
			ids:     []string{squatID, wallPushUpID},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(findByIDsSQL)).
					WithArgs(squatID, wallPushUpID).
					WillReturnRows(sqlmock.NewRows(trainingMenuColumns).
						AddRow(squatID, "squat", "スクワット", "explosive_power", "lower_limb", 10, "reps", 3, "足を肩幅に開き、腰を落とす。", createdAt, updatedAt))
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

			trainingMenuIDs := make([]training.TrainingMenuID, 0, len(tt.ids))
			for _, id := range tt.ids {
				trainingMenuID, err := training.NewTrainingMenuIDFromString(id)
				if err != nil {
					t.Fatalf("failed to new training menu id: %s", err)
				}
				trainingMenuIDs = append(trainingMenuIDs, trainingMenuID)
			}

			repo := NewTrainingMenuRepository(gormDB)

			trainingMenus, err := repo.FindByIDs(context.Background(), trainingMenuIDs)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && trainingMenus != nil {
				t.Errorf("expected no training menus on failure, but got %v", len(trainingMenus))
			}
			if tt.success {
				if len(trainingMenus) != len(tt.wantCodes) {
					t.Errorf("len(trainingMenus) = %v, want %v", len(trainingMenus), len(tt.wantCodes))
				}
				for i := range tt.wantCodes {
					if i >= len(trainingMenus) {
						break
					}
					if trainingMenus[i].Code().String() != tt.wantCodes[i] {
						t.Errorf("trainingMenus[%d].Code() = %v, want %v", i, trainingMenus[i].Code().String(), tt.wantCodes[i])
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
