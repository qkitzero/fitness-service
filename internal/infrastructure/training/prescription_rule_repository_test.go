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
)

const (
	listElementMenusSQL   = `SELECT * FROM "element_menus"`
	listFixedMenusSQL     = `SELECT * FROM "fixed_menus" ORDER BY sort_order`
	listAgeDecadeMenusSQL = `SELECT * FROM "age_decade_menus" ORDER BY decade, sort_order`

	ruleMenuID1 = "1a4c9b2e-6d70-4f18-8b35-90c7e21a4f01"
	ruleMenuID2 = "2b5d0c3f-7e81-4029-9c46-01d8f32b5002"
	ruleMenuID3 = "3c6e1d40-8f92-413a-8d57-12e904436103"
	ruleMenuID4 = "4d7f2e51-9013-424b-9e68-23fa15547204"
	ruleMenuID5 = "5e803f62-0124-435c-8f79-340b26658305"
)

var (
	elementMenuColumns   = []string{"element", "part", "level", "training_menu_id", "created_at", "updated_at"}
	fixedMenuColumns     = []string{"sort_order", "training_menu_id", "created_at", "updated_at"}
	ageDecadeMenuColumns = []string{"decade", "sort_order", "training_menu_id", "created_at", "updated_at"}
)

func TestListElementMenus(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	tests := []struct {
		name                string
		success             bool
		wantElements        []string
		wantParts           []string
		wantLevels          []int
		wantTrainingMenuIDs []string
		setup               func(mock sqlmock.Sqlmock)
	}{
		{
			name:                "success list element menus in element order then part order then level order",
			success:             true,
			wantElements:        []string{"muscle_strength", "muscle_strength", "muscle_strength", "flexibility", "mobility"},
			wantParts:           []string{"upper_limb", "upper_limb", "lower_limb", "upper_limb", "whole_body"},
			wantLevels:          []int{1, 2, 3, 2, 3},
			wantTrainingMenuIDs: []string{ruleMenuID1, ruleMenuID2, ruleMenuID3, ruleMenuID4, ruleMenuID5},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listElementMenusSQL)).
					WillReturnRows(sqlmock.NewRows(elementMenuColumns).
						AddRow("mobility", "whole_body", 3, ruleMenuID5, createdAt, updatedAt).
						AddRow("muscle_strength", "upper_limb", 2, ruleMenuID2, createdAt, updatedAt).
						AddRow("flexibility", "upper_limb", 2, ruleMenuID4, createdAt, updatedAt).
						AddRow("muscle_strength", "lower_limb", 3, ruleMenuID3, createdAt, updatedAt).
						AddRow("muscle_strength", "upper_limb", 1, ruleMenuID1, createdAt, updatedAt))
			},
		},
		{
			name:                "success list no element menus",
			success:             true,
			wantElements:        []string{},
			wantParts:           []string{},
			wantLevels:          []int{},
			wantTrainingMenuIDs: []string{},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listElementMenusSQL)).
					WillReturnRows(sqlmock.NewRows(elementMenuColumns))
			},
		},
		{
			name:    "failure list element menus error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listElementMenusSQL)).
					WillReturnError(errors.New("list element menus error"))
			},
		},
		{
			name:    "failure unknown element",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listElementMenusSQL)).
					WillReturnRows(sqlmock.NewRows(elementMenuColumns).
						AddRow("explosive_power", "lower_limb", 3, ruleMenuID3, createdAt, updatedAt))
			},
		},
		{
			name:    "failure unknown part",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listElementMenusSQL)).
					WillReturnRows(sqlmock.NewRows(elementMenuColumns).
						AddRow("muscle_strength", "trunk", 3, ruleMenuID3, createdAt, updatedAt))
			},
		},
		{
			name:    "failure level above the upper bound",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listElementMenusSQL)).
					WillReturnRows(sqlmock.NewRows(elementMenuColumns).
						AddRow("muscle_strength", "lower_limb", 6, ruleMenuID3, createdAt, updatedAt))
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

			repo := NewPrescriptionRuleRepository(gormDB)

			elementMenus, err := repo.ListElementMenus(context.Background())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && elementMenus != nil {
				t.Errorf("expected no element menus on failure, but got %v", len(elementMenus))
			}
			if tt.success {
				if len(elementMenus) != len(tt.wantElements) {
					t.Errorf("len(elementMenus) = %v, want %v", len(elementMenus), len(tt.wantElements))
				}
				for i := range tt.wantElements {
					if i >= len(elementMenus) {
						break
					}
					e := elementMenus[i]
					if e.Element().String() != tt.wantElements[i] {
						t.Errorf("elementMenus[%d].Element() = %v, want %v", i, e.Element().String(), tt.wantElements[i])
					}
					if e.Part().String() != tt.wantParts[i] {
						t.Errorf("elementMenus[%d].Part() = %v, want %v", i, e.Part().String(), tt.wantParts[i])
					}
					if e.Level().Int() != tt.wantLevels[i] {
						t.Errorf("elementMenus[%d].Level() = %v, want %v", i, e.Level().Int(), tt.wantLevels[i])
					}
					if e.TrainingMenuID().String() != tt.wantTrainingMenuIDs[i] {
						t.Errorf("elementMenus[%d].TrainingMenuID() = %v, want %v", i, e.TrainingMenuID().String(), tt.wantTrainingMenuIDs[i])
					}
					if !e.CreatedAt().Equal(createdAt) {
						t.Errorf("elementMenus[%d].CreatedAt() = %v, want %v", i, e.CreatedAt(), createdAt)
					}
					if !e.UpdatedAt().Equal(updatedAt) {
						t.Errorf("elementMenus[%d].UpdatedAt() = %v, want %v", i, e.UpdatedAt(), updatedAt)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListFixedMenus(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	tests := []struct {
		name                string
		success             bool
		wantSortOrders      []int
		wantTrainingMenuIDs []string
		setup               func(mock sqlmock.Sqlmock)
	}{
		{
			name:                "success list fixed menus ordered by the database",
			success:             true,
			wantSortOrders:      []int{1, 2, 3},
			wantTrainingMenuIDs: []string{ruleMenuID1, ruleMenuID2, ruleMenuID3},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listFixedMenusSQL)).
					WillReturnRows(sqlmock.NewRows(fixedMenuColumns).
						AddRow(1, ruleMenuID1, createdAt, updatedAt).
						AddRow(2, ruleMenuID2, createdAt, updatedAt).
						AddRow(3, ruleMenuID3, createdAt, updatedAt))
			},
		},
		{
			name:                "success list no fixed menus",
			success:             true,
			wantSortOrders:      []int{},
			wantTrainingMenuIDs: []string{},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listFixedMenusSQL)).
					WillReturnRows(sqlmock.NewRows(fixedMenuColumns))
			},
		},
		{
			name:    "failure list fixed menus error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listFixedMenusSQL)).
					WillReturnError(errors.New("list fixed menus error"))
			},
		},
		{
			name:    "failure sort order below the lower bound",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listFixedMenusSQL)).
					WillReturnRows(sqlmock.NewRows(fixedMenuColumns).
						AddRow(0, ruleMenuID1, createdAt, updatedAt))
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

			repo := NewPrescriptionRuleRepository(gormDB)

			fixedMenus, err := repo.ListFixedMenus(context.Background())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && fixedMenus != nil {
				t.Errorf("expected no fixed menus on failure, but got %v", len(fixedMenus))
			}
			if tt.success {
				if len(fixedMenus) != len(tt.wantSortOrders) {
					t.Errorf("len(fixedMenus) = %v, want %v", len(fixedMenus), len(tt.wantSortOrders))
				}
				for i := range tt.wantSortOrders {
					if i >= len(fixedMenus) {
						break
					}
					f := fixedMenus[i]
					if f.SortOrder().Int() != tt.wantSortOrders[i] {
						t.Errorf("fixedMenus[%d].SortOrder() = %v, want %v", i, f.SortOrder().Int(), tt.wantSortOrders[i])
					}
					if f.TrainingMenuID().String() != tt.wantTrainingMenuIDs[i] {
						t.Errorf("fixedMenus[%d].TrainingMenuID() = %v, want %v", i, f.TrainingMenuID().String(), tt.wantTrainingMenuIDs[i])
					}
					if !f.CreatedAt().Equal(createdAt) {
						t.Errorf("fixedMenus[%d].CreatedAt() = %v, want %v", i, f.CreatedAt(), createdAt)
					}
					if !f.UpdatedAt().Equal(updatedAt) {
						t.Errorf("fixedMenus[%d].UpdatedAt() = %v, want %v", i, f.UpdatedAt(), updatedAt)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListAgeDecadeMenus(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	tests := []struct {
		name                string
		success             bool
		wantDecades         []int
		wantSortOrders      []int
		wantTrainingMenuIDs []string
		setup               func(mock sqlmock.Sqlmock)
	}{
		{
			name:                "success list age decade menus ordered by the database",
			success:             true,
			wantDecades:         []int{20, 20, 90},
			wantSortOrders:      []int{1, 2, 1},
			wantTrainingMenuIDs: []string{ruleMenuID1, ruleMenuID2, ruleMenuID3},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listAgeDecadeMenusSQL)).
					WillReturnRows(sqlmock.NewRows(ageDecadeMenuColumns).
						AddRow(20, 1, ruleMenuID1, createdAt, updatedAt).
						AddRow(20, 2, ruleMenuID2, createdAt, updatedAt).
						AddRow(90, 1, ruleMenuID3, createdAt, updatedAt))
			},
		},
		{
			name:                "success list no age decade menus",
			success:             true,
			wantDecades:         []int{},
			wantSortOrders:      []int{},
			wantTrainingMenuIDs: []string{},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listAgeDecadeMenusSQL)).
					WillReturnRows(sqlmock.NewRows(ageDecadeMenuColumns))
			},
		},
		{
			name:    "failure list age decade menus error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listAgeDecadeMenusSQL)).
					WillReturnError(errors.New("list age decade menus error"))
			},
		},
		{
			name:    "failure decade not a multiple of ten",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listAgeDecadeMenusSQL)).
					WillReturnRows(sqlmock.NewRows(ageDecadeMenuColumns).
						AddRow(15, 1, ruleMenuID1, createdAt, updatedAt))
			},
		},
		{
			name:    "failure sort order below the lower bound",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listAgeDecadeMenusSQL)).
					WillReturnRows(sqlmock.NewRows(ageDecadeMenuColumns).
						AddRow(20, 0, ruleMenuID1, createdAt, updatedAt))
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

			repo := NewPrescriptionRuleRepository(gormDB)

			ageDecadeMenus, err := repo.ListAgeDecadeMenus(context.Background())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && ageDecadeMenus != nil {
				t.Errorf("expected no age decade menus on failure, but got %v", len(ageDecadeMenus))
			}
			if tt.success {
				if len(ageDecadeMenus) != len(tt.wantDecades) {
					t.Errorf("len(ageDecadeMenus) = %v, want %v", len(ageDecadeMenus), len(tt.wantDecades))
				}
				for i := range tt.wantDecades {
					if i >= len(ageDecadeMenus) {
						break
					}
					a := ageDecadeMenus[i]
					if a.Decade().Int() != tt.wantDecades[i] {
						t.Errorf("ageDecadeMenus[%d].Decade() = %v, want %v", i, a.Decade().Int(), tt.wantDecades[i])
					}
					if a.SortOrder().Int() != tt.wantSortOrders[i] {
						t.Errorf("ageDecadeMenus[%d].SortOrder() = %v, want %v", i, a.SortOrder().Int(), tt.wantSortOrders[i])
					}
					if a.TrainingMenuID().String() != tt.wantTrainingMenuIDs[i] {
						t.Errorf("ageDecadeMenus[%d].TrainingMenuID() = %v, want %v", i, a.TrainingMenuID().String(), tt.wantTrainingMenuIDs[i])
					}
					if !a.CreatedAt().Equal(createdAt) {
						t.Errorf("ageDecadeMenus[%d].CreatedAt() = %v, want %v", i, a.CreatedAt(), createdAt)
					}
					if !a.UpdatedAt().Equal(updatedAt) {
						t.Errorf("ageDecadeMenus[%d].UpdatedAt() = %v, want %v", i, a.UpdatedAt(), updatedAt)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
