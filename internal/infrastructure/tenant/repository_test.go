package tenant

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

	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/contact"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
	mockstenant "github.com/qkitzero/fitness-service/mocks/domain/tenant"
)

var (
	testCreatedAt = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	testUpdatedAt = time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
)

func TestFindByTenantID(t *testing.T) {
	t.Parallel()
	columns := []string{"tenant_id", "postal_code", "prefecture", "city", "street", "building", "phone", "email", "homepage_url", "note", "created_at", "updated_at"}
	tests := []struct {
		name            string
		success         bool
		wantErr         error
		tenantID        tenant.TenantID
		setup           func(mock sqlmock.Sqlmock, tenantID tenant.TenantID)
		wantPostalCode  string
		wantPrefecture  string
		wantCity        string
		wantStreet      string
		wantBuilding    string
		wantPhone       string
		wantEmail       string
		wantHomepageURL string
		wantNote        string
	}{
		{
			name:     "success find profile by tenant id",
			success:  true,
			wantErr:  nil,
			tenantID: tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, tenantID tenant.TenantID) {
				profileRows := sqlmock.NewRows(columns).
					AddRow(tenantID, "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "0312345678", "test@example.com", "https://example.com", "テスト備考", testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tenant_profiles" WHERE tenant_id = $1 ORDER BY "tenant_profiles"."tenant_id" LIMIT $2`)).
					WithArgs(tenantID, 1).
					WillReturnRows(profileRows)
			},
			wantPostalCode:  "1234567",
			wantPrefecture:  "東京都",
			wantCity:        "千代田区",
			wantStreet:      "1-1-1",
			wantBuilding:    "テストビル",
			wantPhone:       "0312345678",
			wantEmail:       "test@example.com",
			wantHomepageURL: "https://example.com",
			wantNote:        "テスト備考",
		},
		{
			name:     "success find profile by tenant id without optional fields",
			success:  true,
			wantErr:  nil,
			tenantID: tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, tenantID tenant.TenantID) {
				profileRows := sqlmock.NewRows(columns).
					AddRow(tenantID, nil, nil, nil, nil, nil, nil, nil, nil, nil, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tenant_profiles" WHERE tenant_id = $1 ORDER BY "tenant_profiles"."tenant_id" LIMIT $2`)).
					WithArgs(tenantID, 1).
					WillReturnRows(profileRows)
			},
		},
		{
			name:     "failure profile not found",
			success:  false,
			wantErr:  tenant.ErrProfileNotFound,
			tenantID: tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, tenantID tenant.TenantID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tenant_profiles" WHERE tenant_id = $1 ORDER BY "tenant_profiles"."tenant_id" LIMIT $2`)).
					WithArgs(tenantID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:     "failure find profile error",
			success:  false,
			wantErr:  nil,
			tenantID: tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, tenantID tenant.TenantID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tenant_profiles" WHERE tenant_id = $1 ORDER BY "tenant_profiles"."tenant_id" LIMIT $2`)).
					WithArgs(tenantID, 1).
					WillReturnError(errors.New("find profile error"))
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

			tt.setup(mock, tt.tenantID)

			repo := NewProfileRepository(gormDB)

			p, err := repo.FindByTenantID(context.Background(), tt.tenantID)
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
				if p.TenantID() != tt.tenantID {
					t.Errorf("TenantID() = %v, want %v", p.TenantID(), tt.tenantID)
				}
				addr := p.Address()
				if tt.wantPostalCode == "" && addr.PostalCode() != nil {
					t.Errorf("Address().PostalCode() = %v, want nil", addr.PostalCode())
				}
				if tt.wantPostalCode != "" && (addr.PostalCode() == nil || addr.PostalCode().String() != tt.wantPostalCode) {
					t.Errorf("Address().PostalCode() = %v, want %v", addr.PostalCode(), tt.wantPostalCode)
				}
				if tt.wantPrefecture == "" && addr.Prefecture() != nil {
					t.Errorf("Address().Prefecture() = %v, want nil", addr.Prefecture())
				}
				if tt.wantPrefecture != "" && (addr.Prefecture() == nil || addr.Prefecture().String() != tt.wantPrefecture) {
					t.Errorf("Address().Prefecture() = %v, want %v", addr.Prefecture(), tt.wantPrefecture)
				}
				if tt.wantCity == "" && addr.City() != nil {
					t.Errorf("Address().City() = %v, want nil", addr.City())
				}
				if tt.wantCity != "" && (addr.City() == nil || addr.City().String() != tt.wantCity) {
					t.Errorf("Address().City() = %v, want %v", addr.City(), tt.wantCity)
				}
				if tt.wantStreet == "" && addr.Street() != nil {
					t.Errorf("Address().Street() = %v, want nil", addr.Street())
				}
				if tt.wantStreet != "" && (addr.Street() == nil || addr.Street().String() != tt.wantStreet) {
					t.Errorf("Address().Street() = %v, want %v", addr.Street(), tt.wantStreet)
				}
				if tt.wantBuilding == "" && addr.Building() != nil {
					t.Errorf("Address().Building() = %v, want nil", addr.Building())
				}
				if tt.wantBuilding != "" && (addr.Building() == nil || addr.Building().String() != tt.wantBuilding) {
					t.Errorf("Address().Building() = %v, want %v", addr.Building(), tt.wantBuilding)
				}
				if tt.wantPhone == "" && p.Phone() != nil {
					t.Errorf("Phone() = %v, want nil", p.Phone())
				}
				if tt.wantPhone != "" && (p.Phone() == nil || p.Phone().String() != tt.wantPhone) {
					t.Errorf("Phone() = %v, want %v", p.Phone(), tt.wantPhone)
				}
				if tt.wantEmail == "" && p.Email() != nil {
					t.Errorf("Email() = %v, want nil", p.Email())
				}
				if tt.wantEmail != "" && (p.Email() == nil || p.Email().String() != tt.wantEmail) {
					t.Errorf("Email() = %v, want %v", p.Email(), tt.wantEmail)
				}
				if tt.wantHomepageURL == "" && p.HomepageURL() != nil {
					t.Errorf("HomepageURL() = %v, want nil", p.HomepageURL())
				}
				if tt.wantHomepageURL != "" && (p.HomepageURL() == nil || p.HomepageURL().String() != tt.wantHomepageURL) {
					t.Errorf("HomepageURL() = %v, want %v", p.HomepageURL(), tt.wantHomepageURL)
				}
				if tt.wantNote == "" && p.Note() != nil {
					t.Errorf("Note() = %v, want nil", p.Note())
				}
				if tt.wantNote != "" && (p.Note() == nil || p.Note().String() != tt.wantNote) {
					t.Errorf("Note() = %v, want %v", p.Note(), tt.wantNote)
				}
				if !p.CreatedAt().Equal(testCreatedAt) || !p.UpdatedAt().Equal(testUpdatedAt) {
					t.Errorf("timestamps = %v/%v, want %v/%v", p.CreatedAt(), p.UpdatedAt(), testCreatedAt, testUpdatedAt)
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
	postalCode, _ := address.NewPostalCode("123-4567")
	prefecture, _ := address.NewPrefecture("東京都")
	city, _ := address.NewCity("千代田区")
	street, _ := address.NewStreet("1-1-1")
	building, _ := address.NewBuilding("テストビル")
	addr := address.NewAddress(postalCode, prefecture, city, street, building)
	phone, _ := contact.NewPhone("03-1234-5678")
	email, _ := contact.NewEmail("test@example.com")
	homepageURL, _ := tenant.NewHomepageURL("https://example.com")
	note, _ := tenant.NewNote("テスト備考")

	tests := []struct {
		name        string
		success     bool
		address     address.Address
		phone       *contact.Phone
		email       *contact.Email
		homepageURL *tenant.HomepageURL
		note        *tenant.Note
		setup       func(mock sqlmock.Sqlmock, profile tenant.Profile)
	}{
		{
			name:        "success upsert profile",
			success:     true,
			address:     addr,
			phone:       phone,
			email:       email,
			homepageURL: homepageURL,
			note:        note,
			setup: func(mock sqlmock.Sqlmock, profile tenant.Profile) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "tenant_profiles" ("tenant_id","postal_code","prefecture","city","street","building","phone","email","homepage_url","note","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) ON CONFLICT ("tenant_id") DO UPDATE SET "postal_code"="excluded"."postal_code","prefecture"="excluded"."prefecture","city"="excluded"."city","street"="excluded"."street","building"="excluded"."building","phone"="excluded"."phone","email"="excluded"."email","homepage_url"="excluded"."homepage_url","note"="excluded"."note","updated_at"="excluded"."updated_at"`)).
					WithArgs(profile.TenantID(), "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "0312345678", "test@example.com", "https://example.com", "テスト備考", testCreatedAt, testUpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:        "success upsert profile without optional fields",
			success:     true,
			address:     address.Address{},
			phone:       nil,
			email:       nil,
			homepageURL: nil,
			note:        nil,
			setup: func(mock sqlmock.Sqlmock, profile tenant.Profile) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "tenant_profiles" ("tenant_id","postal_code","prefecture","city","street","building","phone","email","homepage_url","note","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) ON CONFLICT ("tenant_id") DO UPDATE SET "postal_code"="excluded"."postal_code","prefecture"="excluded"."prefecture","city"="excluded"."city","street"="excluded"."street","building"="excluded"."building","phone"="excluded"."phone","email"="excluded"."email","homepage_url"="excluded"."homepage_url","note"="excluded"."note","updated_at"="excluded"."updated_at"`)).
					WithArgs(profile.TenantID(), nil, nil, nil, nil, nil, nil, nil, nil, nil, testCreatedAt, testUpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:        "failure upsert profile error",
			success:     false,
			address:     addr,
			phone:       phone,
			email:       email,
			homepageURL: homepageURL,
			note:        note,
			setup: func(mock sqlmock.Sqlmock, profile tenant.Profile) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "tenant_profiles" ("tenant_id","postal_code","prefecture","city","street","building","phone","email","homepage_url","note","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) ON CONFLICT ("tenant_id") DO UPDATE SET "postal_code"="excluded"."postal_code","prefecture"="excluded"."prefecture","city"="excluded"."city","street"="excluded"."street","building"="excluded"."building","phone"="excluded"."phone","email"="excluded"."email","homepage_url"="excluded"."homepage_url","note"="excluded"."note","updated_at"="excluded"."updated_at"`)).
					WithArgs(profile.TenantID(), "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "0312345678", "test@example.com", "https://example.com", "テスト備考", testCreatedAt, testUpdatedAt).
					WillReturnError(errors.New("upsert profile error"))

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

			mockProfile := mockstenant.NewMockProfile(ctrl)
			mockProfile.EXPECT().TenantID().Return(tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
			mockProfile.EXPECT().Address().Return(tt.address).AnyTimes()
			mockProfile.EXPECT().Phone().Return(tt.phone).AnyTimes()
			mockProfile.EXPECT().Email().Return(tt.email).AnyTimes()
			mockProfile.EXPECT().HomepageURL().Return(tt.homepageURL).AnyTimes()
			mockProfile.EXPECT().Note().Return(tt.note).AnyTimes()
			mockProfile.EXPECT().CreatedAt().Return(testCreatedAt).AnyTimes()
			mockProfile.EXPECT().UpdatedAt().Return(testUpdatedAt).AnyTimes()

			tt.setup(mock, mockProfile)

			repo := NewProfileRepository(gormDB)

			err = repo.Upsert(context.Background(), mockProfile)
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
