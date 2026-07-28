package organization

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
	mocksorganization "github.com/qkitzero/fitness-service/mocks/domain/organization"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, organization organization.Organization)
	}{
		{
			name:    "success create organization",
			success: true,
			setup: func(mock sqlmock.Sqlmock, organization organization.Organization) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "organizations" ("id","group_id","name","created_at","updated_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(organization.ID(), organization.TenantID(), organization.Name(), createdAt, updatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "failure create organization error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, organization organization.Organization) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "organizations" ("id","group_id","name","created_at","updated_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(organization.ID(), organization.TenantID(), organization.Name(), createdAt, updatedAt).
					WillReturnError(errors.New("create organization error"))

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

			mockOrganization := mocksorganization.NewMockOrganization(ctrl)
			mockOrganization.EXPECT().ID().Return(organization.OrganizationID{UUID: uuid.New()}).AnyTimes()
			mockOrganization.EXPECT().TenantID().Return(tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
			mockOrganization.EXPECT().Name().Return(organization.Name("テスト株式会社")).AnyTimes()
			mockOrganization.EXPECT().CreatedAt().Return(createdAt).AnyTimes()
			mockOrganization.EXPECT().UpdatedAt().Return(updatedAt).AnyTimes()

			tt.setup(mock, mockOrganization)

			repo := NewOrganizationRepository(gormDB)

			err = repo.Create(context.Background(), mockOrganization)
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
	columns := []string{"id", "group_id", "name", "created_at", "updated_at"}
	tests := []struct {
		name           string
		success        bool
		wantErr        error
		organizationID organization.OrganizationID
		setup          func(mock sqlmock.Sqlmock, organizationID organization.OrganizationID)
	}{
		{
			name:           "success find organization by id",
			success:        true,
			wantErr:        nil,
			organizationID: organization.OrganizationID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, organizationID organization.OrganizationID) {
				organizationRows := sqlmock.NewRows(columns).
					AddRow(organizationID, "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "テスト株式会社", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE id = $1 ORDER BY "organizations"."id" LIMIT $2`)).
					WithArgs(organizationID, 1).
					WillReturnRows(organizationRows)
			},
		},
		{
			name:           "failure organization not found",
			success:        false,
			wantErr:        organization.ErrOrganizationNotFound,
			organizationID: organization.OrganizationID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, organizationID organization.OrganizationID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE id = $1 ORDER BY "organizations"."id" LIMIT $2`)).
					WithArgs(organizationID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:           "failure find organization error",
			success:        false,
			wantErr:        nil,
			organizationID: organization.OrganizationID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, organizationID organization.OrganizationID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE id = $1 ORDER BY "organizations"."id" LIMIT $2`)).
					WithArgs(organizationID, 1).
					WillReturnError(errors.New("find organization error"))
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

			tt.setup(mock, tt.organizationID)

			repo := NewOrganizationRepository(gormDB)

			o, err := repo.FindByID(context.Background(), tt.organizationID)
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
				if o.ID() != tt.organizationID {
					t.Errorf("ID() = %v, want %v", o.ID(), tt.organizationID)
				}
				if o.TenantID().String() != "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b" {
					t.Errorf("TenantID() = %v, want 0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", o.TenantID().String())
				}
				if o.Name().String() != "テスト株式会社" {
					t.Errorf("Name() = %v, want テスト株式会社", o.Name().String())
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListByTenantID(t *testing.T) {
	t.Parallel()
	columns := []string{"id", "group_id", "name", "created_at", "updated_at"}
	tests := []struct {
		name      string
		success   bool
		wantNames []string
		tenantID  tenant.TenantID
		setup     func(mock sqlmock.Sqlmock, tenantID tenant.TenantID)
	}{
		{
			name:      "success list organizations by tenant id",
			success:   true,
			wantNames: []string{"テスト株式会社", "テスト工業株式会社"},
			tenantID:  tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, tenantID tenant.TenantID) {
				organizationRows := sqlmock.NewRows(columns).
					AddRow(uuid.New().String(), tenantID, "テスト株式会社", time.Now(), time.Now()).
					AddRow(uuid.New().String(), tenantID, "テスト工業株式会社", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE group_id = $1 ORDER BY created_at, id`)).
					WithArgs(tenantID).
					WillReturnRows(organizationRows)
			},
		},
		{
			name:      "success list no organizations",
			success:   true,
			wantNames: []string{},
			tenantID:  tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, tenantID tenant.TenantID) {
				organizationRows := sqlmock.NewRows(columns)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE group_id = $1 ORDER BY created_at, id`)).
					WithArgs(tenantID).
					WillReturnRows(organizationRows)
			},
		},
		{
			name:      "failure list organizations error",
			success:   false,
			wantNames: nil,
			tenantID:  tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, tenantID tenant.TenantID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE group_id = $1 ORDER BY created_at, id`)).
					WithArgs(tenantID).
					WillReturnError(errors.New("list organizations error"))
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

			repo := NewOrganizationRepository(gormDB)

			organizations, err := repo.ListByTenantID(context.Background(), tt.tenantID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success {
				if len(organizations) != len(tt.wantNames) {
					t.Errorf("len(organizations) = %v, want %v", len(organizations), len(tt.wantNames))
				}
				for i, wantName := range tt.wantNames {
					if i >= len(organizations) {
						break
					}
					if organizations[i].Name().String() != wantName {
						t.Errorf("organizations[%d].Name() = %v, want %v", i, organizations[i].Name().String(), wantName)
					}
					if organizations[i].TenantID() != tt.tenantID {
						t.Errorf("organizations[%d].TenantID() = %v, want %v", i, organizations[i].TenantID(), tt.tenantID)
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
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	tests := []struct {
		name    string
		success bool
		wantErr error
		setup   func(mock sqlmock.Sqlmock, organization organization.Organization)
	}{
		{
			name:    "success update organization",
			success: true,
			wantErr: nil,
			setup: func(mock sqlmock.Sqlmock, organization organization.Organization) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "organizations" SET "group_id"=$1,"name"=$2,"created_at"=$3,"updated_at"=$4 WHERE id = $5`)).
					WithArgs(organization.TenantID(), organization.Name(), createdAt, updatedAt, organization.ID()).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "failure organization not found",
			success: false,
			wantErr: organization.ErrOrganizationNotFound,
			setup: func(mock sqlmock.Sqlmock, organization organization.Organization) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "organizations" SET "group_id"=$1,"name"=$2,"created_at"=$3,"updated_at"=$4 WHERE id = $5`)).
					WithArgs(organization.TenantID(), organization.Name(), createdAt, updatedAt, organization.ID()).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectRollback()
			},
		},
		{
			name:    "failure update organization error",
			success: false,
			wantErr: nil,
			setup: func(mock sqlmock.Sqlmock, organization organization.Organization) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "organizations" SET "group_id"=$1,"name"=$2,"created_at"=$3,"updated_at"=$4 WHERE id = $5`)).
					WithArgs(organization.TenantID(), organization.Name(), createdAt, updatedAt, organization.ID()).
					WillReturnError(errors.New("update organization error"))

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

			mockOrganization := mocksorganization.NewMockOrganization(ctrl)
			mockOrganization.EXPECT().ID().Return(organization.OrganizationID{UUID: uuid.New()}).AnyTimes()
			mockOrganization.EXPECT().TenantID().Return(tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
			mockOrganization.EXPECT().Name().Return(organization.Name("更新株式会社")).AnyTimes()
			mockOrganization.EXPECT().CreatedAt().Return(createdAt).AnyTimes()
			mockOrganization.EXPECT().UpdatedAt().Return(updatedAt).AnyTimes()

			tt.setup(mock, mockOrganization)

			repo := NewOrganizationRepository(gormDB)

			err = repo.Update(context.Background(), mockOrganization)
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
	tests := []struct {
		name           string
		success        bool
		wantErr        error
		organizationID organization.OrganizationID
		setup          func(mock sqlmock.Sqlmock, organizationID organization.OrganizationID)
	}{
		{
			name:           "success delete organization",
			success:        true,
			wantErr:        nil,
			organizationID: organization.OrganizationID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, organizationID organization.OrganizationID) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "organizations" WHERE id = $1`)).
					WithArgs(organizationID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:           "failure organization not found",
			success:        false,
			wantErr:        organization.ErrOrganizationNotFound,
			organizationID: organization.OrganizationID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, organizationID organization.OrganizationID) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "organizations" WHERE id = $1`)).
					WithArgs(organizationID).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectRollback()
			},
		},
		{
			name:           "failure delete organization error",
			success:        false,
			wantErr:        nil,
			organizationID: organization.OrganizationID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, organizationID organization.OrganizationID) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "organizations" WHERE id = $1`)).
					WithArgs(organizationID).
					WillReturnError(errors.New("delete organization error"))

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

			tt.setup(mock, tt.organizationID)

			repo := NewOrganizationRepository(gormDB)

			err = repo.Delete(context.Background(), tt.organizationID)
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
