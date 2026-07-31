package customer

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

	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/contact"
	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
	mockscustomer "github.com/qkitzero/fitness-service/mocks/domain/customer"
)

var (
	testBirthDate      = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	testOrganizationID = "3f2b6c1d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	testCreatedAt      = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	testUpdatedAt      = time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
)

func TestCreate(t *testing.T) {
	t.Parallel()
	phone, _ := contact.NewPhone("03-1234-5678")
	email, _ := contact.NewEmail("test@example.com")
	postalCode, _ := address.NewPostalCode("123-4567")
	prefecture, _ := address.NewPrefecture("東京都")
	city, _ := address.NewCity("千代田区")
	street, _ := address.NewStreet("1-1-1")
	building, _ := address.NewBuilding("テストビル")
	addr := address.NewAddress(postalCode, prefecture, city, street, building)
	emergencyContactName, _ := customer.NewEmergencyContactName("緊急 太郎")
	emergencyContactRelationship, _ := customer.NewEmergencyContactRelationship("父")
	emergencyContactPhone, _ := contact.NewPhone("090-1234-5678")
	organizationID, _ := organization.NewOrganizationIDFromString(testOrganizationID)
	foreignKeyViolation := errors.New(`pq: insert or update on table "customers" violates foreign key constraint "fk_customers_organization"`)

	tests := []struct {
		name                         string
		success                      bool
		wantErr                      error
		phone                        *contact.Phone
		email                        *contact.Email
		address                      address.Address
		emergencyContactName         *customer.EmergencyContactName
		emergencyContactRelationship *customer.EmergencyContactRelationship
		emergencyContactPhone        *contact.Phone
		organizationID               *organization.OrganizationID
		setup                        func(mock sqlmock.Sqlmock, customer customer.Customer)
	}{
		{
			name:                         "success create customer",
			success:                      true,
			wantErr:                      nil,
			phone:                        phone,
			email:                        email,
			address:                      addr,
			emergencyContactName:         emergencyContactName,
			emergencyContactRelationship: emergencyContactRelationship,
			emergencyContactPhone:        emergencyContactPhone,
			organizationID:               &organizationID,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "customers" ("id","tenant_id","name","name_kana","gender","birth_date","phone","email","postal_code","prefecture","city","street","building","emergency_contact_name","emergency_contact_relationship","emergency_contact_phone","organization_id","is_active","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`)).
					WithArgs(customer.ID(), customer.TenantID(), customer.Name(), customer.NameKana(), customer.Gender(), testBirthDate, "0312345678", "test@example.com", "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "緊急 太郎", "父", "09012345678", testOrganizationID, true, testCreatedAt, testUpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:                         "success create customer without optional fields",
			success:                      true,
			wantErr:                      nil,
			phone:                        nil,
			email:                        nil,
			address:                      address.Address{},
			emergencyContactName:         nil,
			emergencyContactRelationship: nil,
			emergencyContactPhone:        nil,
			organizationID:               nil,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "customers" ("id","tenant_id","name","name_kana","gender","birth_date","phone","email","postal_code","prefecture","city","street","building","emergency_contact_name","emergency_contact_relationship","emergency_contact_phone","organization_id","is_active","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`)).
					WithArgs(customer.ID(), customer.TenantID(), customer.Name(), customer.NameKana(), customer.Gender(), testBirthDate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, true, testCreatedAt, testUpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:                         "failure unknown organization",
			success:                      false,
			wantErr:                      foreignKeyViolation,
			phone:                        nil,
			email:                        nil,
			address:                      address.Address{},
			emergencyContactName:         nil,
			emergencyContactRelationship: nil,
			emergencyContactPhone:        nil,
			organizationID:               &organizationID,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "customers" ("id","tenant_id","name","name_kana","gender","birth_date","phone","email","postal_code","prefecture","city","street","building","emergency_contact_name","emergency_contact_relationship","emergency_contact_phone","organization_id","is_active","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`)).
					WithArgs(customer.ID(), customer.TenantID(), customer.Name(), customer.NameKana(), customer.Gender(), testBirthDate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, testOrganizationID, true, testCreatedAt, testUpdatedAt).
					WillReturnError(foreignKeyViolation)

				mock.ExpectRollback()
			},
		},
		{
			name:                         "failure create customer error",
			success:                      false,
			wantErr:                      nil,
			phone:                        phone,
			email:                        email,
			address:                      addr,
			emergencyContactName:         emergencyContactName,
			emergencyContactRelationship: emergencyContactRelationship,
			emergencyContactPhone:        emergencyContactPhone,
			organizationID:               &organizationID,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "customers" ("id","tenant_id","name","name_kana","gender","birth_date","phone","email","postal_code","prefecture","city","street","building","emergency_contact_name","emergency_contact_relationship","emergency_contact_phone","organization_id","is_active","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`)).
					WithArgs(customer.ID(), customer.TenantID(), customer.Name(), customer.NameKana(), customer.Gender(), testBirthDate, "0312345678", "test@example.com", "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "緊急 太郎", "父", "09012345678", testOrganizationID, true, testCreatedAt, testUpdatedAt).
					WillReturnError(errors.New("create customer error"))

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

			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().ID().Return(customer.CustomerID{UUID: uuid.New()}).AnyTimes()
			mockCustomer.EXPECT().TenantID().Return(tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
			mockCustomer.EXPECT().Name().Return(customer.Name("test customer")).AnyTimes()
			mockCustomer.EXPECT().NameKana().Return(customer.NameKana("テストカナ")).AnyTimes()
			mockCustomer.EXPECT().Gender().Return(customer.GenderMale).AnyTimes()
			mockCustomer.EXPECT().BirthDate().Return(customer.BirthDate{Time: testBirthDate}).AnyTimes()
			mockCustomer.EXPECT().IsActive().Return(true).AnyTimes()
			mockCustomer.EXPECT().Phone().Return(tt.phone).AnyTimes()
			mockCustomer.EXPECT().Email().Return(tt.email).AnyTimes()
			mockCustomer.EXPECT().Address().Return(tt.address).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactName().Return(tt.emergencyContactName).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactRelationship().Return(tt.emergencyContactRelationship).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactPhone().Return(tt.emergencyContactPhone).AnyTimes()
			mockCustomer.EXPECT().OrganizationID().Return(tt.organizationID).AnyTimes()
			mockCustomer.EXPECT().CreatedAt().Return(testCreatedAt).AnyTimes()
			mockCustomer.EXPECT().UpdatedAt().Return(testUpdatedAt).AnyTimes()

			tt.setup(mock, mockCustomer)

			repo := NewCustomerRepository(gormDB)

			err = repo.Create(context.Background(), mockCustomer)
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

func TestFindByID(t *testing.T) {
	t.Parallel()
	columns := []string{"id", "tenant_id", "name", "name_kana", "gender", "birth_date", "phone", "email", "postal_code", "prefecture", "city", "street", "building", "emergency_contact_name", "emergency_contact_relationship", "emergency_contact_phone", "organization_id", "is_active", "created_at", "updated_at"}
	tests := []struct {
		name               string
		success            bool
		wantErr            error
		customerID         customer.CustomerID
		setup              func(mock sqlmock.Sqlmock, customerID customer.CustomerID)
		wantPhone          string
		wantEmail          string
		wantPostalCode     string
		wantPrefecture     string
		wantCity           string
		wantStreet         string
		wantBuilding       string
		wantOrganizationID string
		wantIsActive       bool
	}{
		{
			name:       "success find customer by id",
			success:    true,
			wantErr:    nil,
			customerID: customer.CustomerID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, customerID customer.CustomerID) {
				customerRows := sqlmock.NewRows(columns).
					AddRow(customerID, "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "test customer", "テストカナ", "male", testBirthDate, "0312345678", "test@example.com", "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "緊急 太郎", "父", "09012345678", testOrganizationID, false, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE id = $1 ORDER BY "customers"."id" LIMIT $2`)).
					WithArgs(customerID, 1).
					WillReturnRows(customerRows)
			},
			wantPhone:          "0312345678",
			wantEmail:          "test@example.com",
			wantPostalCode:     "1234567",
			wantPrefecture:     "東京都",
			wantCity:           "千代田区",
			wantStreet:         "1-1-1",
			wantBuilding:       "テストビル",
			wantOrganizationID: testOrganizationID,
			wantIsActive:       false,
		},
		{
			name:       "success find customer by id without optional fields",
			success:    true,
			wantErr:    nil,
			customerID: customer.CustomerID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, customerID customer.CustomerID) {
				customerRows := sqlmock.NewRows(columns).
					AddRow(customerID, "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "test customer", "テストカナ", "male", testBirthDate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, true, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE id = $1 ORDER BY "customers"."id" LIMIT $2`)).
					WithArgs(customerID, 1).
					WillReturnRows(customerRows)
			},
			wantPhone:          "",
			wantEmail:          "",
			wantPostalCode:     "",
			wantPrefecture:     "",
			wantCity:           "",
			wantStreet:         "",
			wantBuilding:       "",
			wantOrganizationID: "",
			wantIsActive:       true,
		},
		{
			name:       "failure customer not found",
			success:    false,
			wantErr:    customer.ErrCustomerNotFound,
			customerID: customer.CustomerID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, customerID customer.CustomerID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE id = $1 ORDER BY "customers"."id" LIMIT $2`)).
					WithArgs(customerID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:       "failure find customer error",
			success:    false,
			wantErr:    nil,
			customerID: customer.CustomerID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, customerID customer.CustomerID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE id = $1 ORDER BY "customers"."id" LIMIT $2`)).
					WithArgs(customerID, 1).
					WillReturnError(errors.New("find customer error"))
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

			tt.setup(mock, tt.customerID)

			repo := NewCustomerRepository(gormDB)

			c, err := repo.FindByID(context.Background(), tt.customerID)
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
				if c.ID() != tt.customerID {
					t.Errorf("ID() = %v, want %v", c.ID(), tt.customerID)
				}
				if c.Name().String() != "test customer" {
					t.Errorf("Name() = %v, want test customer", c.Name().String())
				}
				if c.IsActive() != tt.wantIsActive {
					t.Errorf("IsActive() = %v, want %v", c.IsActive(), tt.wantIsActive)
				}
				if tt.wantOrganizationID == "" && c.OrganizationID() != nil {
					t.Errorf("OrganizationID() = %v, want nil", c.OrganizationID())
				}
				if tt.wantOrganizationID != "" && (c.OrganizationID() == nil || c.OrganizationID().String() != tt.wantOrganizationID) {
					t.Errorf("OrganizationID() = %v, want %v", c.OrganizationID(), tt.wantOrganizationID)
				}
				if tt.wantPhone == "" && c.Phone() != nil {
					t.Errorf("Phone() = %v, want nil", c.Phone())
				}
				if tt.wantPhone != "" && (c.Phone() == nil || c.Phone().String() != tt.wantPhone) {
					t.Errorf("Phone() = %v, want %v", c.Phone(), tt.wantPhone)
				}
				if tt.wantEmail == "" && c.Email() != nil {
					t.Errorf("Email() = %v, want nil", c.Email())
				}
				if tt.wantEmail != "" && (c.Email() == nil || c.Email().String() != tt.wantEmail) {
					t.Errorf("Email() = %v, want %v", c.Email(), tt.wantEmail)
				}
				addr := c.Address()
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
				if !c.CreatedAt().Equal(testCreatedAt) || !c.UpdatedAt().Equal(testUpdatedAt) {
					t.Errorf("timestamps = %v/%v, want %v/%v", c.CreatedAt(), c.UpdatedAt(), testCreatedAt, testUpdatedAt)
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
	columns := []string{"id", "tenant_id", "name", "name_kana", "gender", "birth_date", "phone", "email", "postal_code", "prefecture", "city", "street", "building", "emergency_contact_name", "emergency_contact_relationship", "emergency_contact_phone", "organization_id", "is_active", "created_at", "updated_at"}
	tests := []struct {
		name                string
		success             bool
		wantNames           []string
		wantActive          []bool
		wantOrganizationIDs []string
		includeInactive     bool
		tenantID            tenant.TenantID
		setup               func(mock sqlmock.Sqlmock, tenantID tenant.TenantID)
	}{
		{
			name:                "success list active customers by tenant id",
			success:             true,
			wantNames:           []string{"test customer 1", "test customer 2"},
			wantActive:          []bool{true, true},
			wantOrganizationIDs: []string{testOrganizationID, ""},
			includeInactive:     false,
			tenantID:            tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, tenantID tenant.TenantID) {
				customerRows := sqlmock.NewRows(columns).
					AddRow(uuid.New().String(), tenantID, "test customer 1", "テストカナ", "male", testBirthDate, "0312345678", "test1@example.com", "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "緊急 太郎", "父", "09012345678", testOrganizationID, true, testCreatedAt, testUpdatedAt).
					AddRow(uuid.New().String(), tenantID, "test customer 2", "テストカナ", "female", testBirthDate, "0312345679", "test2@example.com", "1234568", "大阪府", "大阪市", "2-2-2", "更新ビル", "緊急 花子", "母", "08012345678", nil, true, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE tenant_id = $1 AND is_active = $2 ORDER BY created_at, id`)).
					WithArgs(tenantID, true).
					WillReturnRows(customerRows)
			},
		},
		{
			name:                "success list customers including inactive",
			success:             true,
			wantNames:           []string{"test customer 1", "test customer 2"},
			wantActive:          []bool{true, false},
			wantOrganizationIDs: []string{testOrganizationID, ""},
			includeInactive:     true,
			tenantID:            tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, tenantID tenant.TenantID) {
				customerRows := sqlmock.NewRows(columns).
					AddRow(uuid.New().String(), tenantID, "test customer 1", "テストカナ", "male", testBirthDate, "0312345678", "test1@example.com", "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "緊急 太郎", "父", "09012345678", testOrganizationID, true, testCreatedAt, testUpdatedAt).
					AddRow(uuid.New().String(), tenantID, "test customer 2", "テストカナ", "female", testBirthDate, "0312345679", "test2@example.com", "1234568", "大阪府", "大阪市", "2-2-2", "更新ビル", "緊急 花子", "母", "08012345678", nil, false, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE tenant_id = $1 ORDER BY created_at, id`)).
					WithArgs(tenantID).
					WillReturnRows(customerRows)
			},
		},
		{
			name:                "success list no customers",
			success:             true,
			wantNames:           []string{},
			wantActive:          []bool{},
			wantOrganizationIDs: []string{},
			includeInactive:     false,
			tenantID:            tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, tenantID tenant.TenantID) {
				customerRows := sqlmock.NewRows(columns)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE tenant_id = $1 AND is_active = $2 ORDER BY created_at, id`)).
					WithArgs(tenantID, true).
					WillReturnRows(customerRows)
			},
		},
		{
			name:                "failure list customers error",
			success:             false,
			wantNames:           nil,
			wantActive:          nil,
			wantOrganizationIDs: nil,
			includeInactive:     false,
			tenantID:            tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, tenantID tenant.TenantID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE tenant_id = $1 AND is_active = $2 ORDER BY created_at, id`)).
					WithArgs(tenantID, true).
					WillReturnError(errors.New("list customers error"))
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

			repo := NewCustomerRepository(gormDB)

			customers, err := repo.ListByTenantID(context.Background(), tt.tenantID, tt.includeInactive)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success {
				if len(customers) != len(tt.wantNames) {
					t.Errorf("len(customers) = %v, want %v", len(customers), len(tt.wantNames))
				}
				for i, wantName := range tt.wantNames {
					if i >= len(customers) {
						break
					}
					if customers[i].Name().String() != wantName {
						t.Errorf("customers[%d].Name() = %v, want %v", i, customers[i].Name().String(), wantName)
					}
					if customers[i].TenantID() != tt.tenantID {
						t.Errorf("customers[%d].TenantID() = %v, want %v", i, customers[i].TenantID(), tt.tenantID)
					}
					if customers[i].IsActive() != tt.wantActive[i] {
						t.Errorf("customers[%d].IsActive() = %v, want %v", i, customers[i].IsActive(), tt.wantActive[i])
					}
					if tt.wantOrganizationIDs[i] == "" && customers[i].OrganizationID() != nil {
						t.Errorf("customers[%d].OrganizationID() = %v, want nil", i, customers[i].OrganizationID())
					}
					if tt.wantOrganizationIDs[i] != "" && (customers[i].OrganizationID() == nil || customers[i].OrganizationID().String() != tt.wantOrganizationIDs[i]) {
						t.Errorf("customers[%d].OrganizationID() = %v, want %v", i, customers[i].OrganizationID(), tt.wantOrganizationIDs[i])
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
	updatedBirthDate := time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)
	phone, _ := contact.NewPhone("080-1234-5678")
	email, _ := contact.NewEmail("updated@example.com")
	postalCode, _ := address.NewPostalCode("543-2100")
	prefecture, _ := address.NewPrefecture("大阪府")
	city, _ := address.NewCity("大阪市")
	street, _ := address.NewStreet("2-2-2")
	building, _ := address.NewBuilding("更新ビル")
	addr := address.NewAddress(postalCode, prefecture, city, street, building)
	emergencyContactName, _ := customer.NewEmergencyContactName("更新 花子")
	emergencyContactRelationship, _ := customer.NewEmergencyContactRelationship("母")
	emergencyContactPhone, _ := contact.NewPhone("070-1234-5678")
	organizationID, _ := organization.NewOrganizationIDFromString(testOrganizationID)

	tests := []struct {
		name                         string
		success                      bool
		wantErr                      error
		birthDate                    time.Time
		phone                        *contact.Phone
		email                        *contact.Email
		address                      address.Address
		emergencyContactName         *customer.EmergencyContactName
		emergencyContactRelationship *customer.EmergencyContactRelationship
		emergencyContactPhone        *contact.Phone
		organizationID               *organization.OrganizationID
		setup                        func(mock sqlmock.Sqlmock, customer customer.Customer)
	}{
		{
			name:                         "success update customer",
			success:                      true,
			wantErr:                      nil,
			birthDate:                    updatedBirthDate,
			phone:                        phone,
			email:                        email,
			address:                      addr,
			emergencyContactName:         emergencyContactName,
			emergencyContactRelationship: emergencyContactRelationship,
			emergencyContactPhone:        emergencyContactPhone,
			organizationID:               &organizationID,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "tenant_id"=$1,"name"=$2,"name_kana"=$3,"gender"=$4,"birth_date"=$5,"phone"=$6,"email"=$7,"postal_code"=$8,"prefecture"=$9,"city"=$10,"street"=$11,"building"=$12,"emergency_contact_name"=$13,"emergency_contact_relationship"=$14,"emergency_contact_phone"=$15,"organization_id"=$16,"created_at"=$17,"updated_at"=$18 WHERE id = $19`)).
					WithArgs(customer.TenantID(), customer.Name(), customer.NameKana(), customer.Gender(), updatedBirthDate, "08012345678", "updated@example.com", "5432100", "大阪府", "大阪市", "2-2-2", "更新ビル", "更新 花子", "母", "07012345678", testOrganizationID, testCreatedAt, testUpdatedAt, customer.ID()).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:                         "success update customer clearing optional fields",
			success:                      true,
			wantErr:                      nil,
			birthDate:                    testBirthDate,
			phone:                        nil,
			email:                        nil,
			address:                      address.Address{},
			emergencyContactName:         nil,
			emergencyContactRelationship: nil,
			emergencyContactPhone:        nil,
			organizationID:               nil,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "tenant_id"=$1,"name"=$2,"name_kana"=$3,"gender"=$4,"birth_date"=$5,"phone"=$6,"email"=$7,"postal_code"=$8,"prefecture"=$9,"city"=$10,"street"=$11,"building"=$12,"emergency_contact_name"=$13,"emergency_contact_relationship"=$14,"emergency_contact_phone"=$15,"organization_id"=$16,"created_at"=$17,"updated_at"=$18 WHERE id = $19`)).
					WithArgs(customer.TenantID(), customer.Name(), customer.NameKana(), customer.Gender(), testBirthDate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, testCreatedAt, testUpdatedAt, customer.ID()).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:                         "failure customer not found",
			success:                      false,
			wantErr:                      customer.ErrCustomerNotFound,
			birthDate:                    updatedBirthDate,
			phone:                        phone,
			email:                        email,
			address:                      addr,
			emergencyContactName:         emergencyContactName,
			emergencyContactRelationship: emergencyContactRelationship,
			emergencyContactPhone:        emergencyContactPhone,
			organizationID:               &organizationID,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "tenant_id"=$1,"name"=$2,"name_kana"=$3,"gender"=$4,"birth_date"=$5,"phone"=$6,"email"=$7,"postal_code"=$8,"prefecture"=$9,"city"=$10,"street"=$11,"building"=$12,"emergency_contact_name"=$13,"emergency_contact_relationship"=$14,"emergency_contact_phone"=$15,"organization_id"=$16,"created_at"=$17,"updated_at"=$18 WHERE id = $19`)).
					WithArgs(customer.TenantID(), customer.Name(), customer.NameKana(), customer.Gender(), updatedBirthDate, "08012345678", "updated@example.com", "5432100", "大阪府", "大阪市", "2-2-2", "更新ビル", "更新 花子", "母", "07012345678", testOrganizationID, testCreatedAt, testUpdatedAt, customer.ID()).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectRollback()
			},
		},
		{
			name:                         "failure update customer error",
			success:                      false,
			wantErr:                      nil,
			birthDate:                    updatedBirthDate,
			phone:                        phone,
			email:                        email,
			address:                      addr,
			emergencyContactName:         emergencyContactName,
			emergencyContactRelationship: emergencyContactRelationship,
			emergencyContactPhone:        emergencyContactPhone,
			organizationID:               &organizationID,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "tenant_id"=$1,"name"=$2,"name_kana"=$3,"gender"=$4,"birth_date"=$5,"phone"=$6,"email"=$7,"postal_code"=$8,"prefecture"=$9,"city"=$10,"street"=$11,"building"=$12,"emergency_contact_name"=$13,"emergency_contact_relationship"=$14,"emergency_contact_phone"=$15,"organization_id"=$16,"created_at"=$17,"updated_at"=$18 WHERE id = $19`)).
					WithArgs(customer.TenantID(), customer.Name(), customer.NameKana(), customer.Gender(), updatedBirthDate, "08012345678", "updated@example.com", "5432100", "大阪府", "大阪市", "2-2-2", "更新ビル", "更新 花子", "母", "07012345678", testOrganizationID, testCreatedAt, testUpdatedAt, customer.ID()).
					WillReturnError(errors.New("update customer error"))

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

			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().ID().Return(customer.CustomerID{UUID: uuid.New()}).AnyTimes()
			mockCustomer.EXPECT().TenantID().Return(tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
			mockCustomer.EXPECT().Name().Return(customer.Name("updated customer")).AnyTimes()
			mockCustomer.EXPECT().NameKana().Return(customer.NameKana("コウシンカナ")).AnyTimes()
			mockCustomer.EXPECT().Gender().Return(customer.GenderFemale).AnyTimes()
			mockCustomer.EXPECT().BirthDate().Return(customer.BirthDate{Time: tt.birthDate}).AnyTimes()
			mockCustomer.EXPECT().IsActive().Return(false).AnyTimes()
			mockCustomer.EXPECT().Phone().Return(tt.phone).AnyTimes()
			mockCustomer.EXPECT().Email().Return(tt.email).AnyTimes()
			mockCustomer.EXPECT().Address().Return(tt.address).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactName().Return(tt.emergencyContactName).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactRelationship().Return(tt.emergencyContactRelationship).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactPhone().Return(tt.emergencyContactPhone).AnyTimes()
			mockCustomer.EXPECT().OrganizationID().Return(tt.organizationID).AnyTimes()
			mockCustomer.EXPECT().CreatedAt().Return(testCreatedAt).AnyTimes()
			mockCustomer.EXPECT().UpdatedAt().Return(testUpdatedAt).AnyTimes()

			tt.setup(mock, mockCustomer)

			repo := NewCustomerRepository(gormDB)

			err = repo.Update(context.Background(), mockCustomer)
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

func TestUpdateActive(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		wantErr error
		setup   func(mock sqlmock.Sqlmock, customer customer.Customer)
	}{
		{
			name:    "success update active",
			success: true,
			wantErr: nil,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "is_active"=$1,"updated_at"=$2 WHERE id = $3`)).
					WithArgs(false, testUpdatedAt, customer.ID()).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "failure customer not found",
			success: false,
			wantErr: customer.ErrCustomerNotFound,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "is_active"=$1,"updated_at"=$2 WHERE id = $3`)).
					WithArgs(false, testUpdatedAt, customer.ID()).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectRollback()
			},
		},
		{
			name:    "failure update active error",
			success: false,
			wantErr: nil,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "is_active"=$1,"updated_at"=$2 WHERE id = $3`)).
					WithArgs(false, testUpdatedAt, customer.ID()).
					WillReturnError(errors.New("update active error"))

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

			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().ID().Return(customer.CustomerID{UUID: uuid.New()}).AnyTimes()
			mockCustomer.EXPECT().TenantID().Return(tenant.TenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
			mockCustomer.EXPECT().Name().Return(customer.Name("test customer")).AnyTimes()
			mockCustomer.EXPECT().NameKana().Return(customer.NameKana("テストカナ")).AnyTimes()
			mockCustomer.EXPECT().Gender().Return(customer.GenderMale).AnyTimes()
			mockCustomer.EXPECT().BirthDate().Return(customer.BirthDate{Time: testBirthDate}).AnyTimes()
			mockCustomer.EXPECT().IsActive().Return(false).AnyTimes()
			mockCustomer.EXPECT().Phone().Return(nil).AnyTimes()
			mockCustomer.EXPECT().Email().Return(nil).AnyTimes()
			mockCustomer.EXPECT().Address().Return(address.Address{}).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactName().Return(nil).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactRelationship().Return(nil).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactPhone().Return(nil).AnyTimes()
			mockCustomer.EXPECT().OrganizationID().Return(nil).AnyTimes()
			mockCustomer.EXPECT().CreatedAt().Return(testCreatedAt).AnyTimes()
			mockCustomer.EXPECT().UpdatedAt().Return(testUpdatedAt).AnyTimes()

			tt.setup(mock, mockCustomer)

			repo := NewCustomerRepository(gormDB)

			err = repo.UpdateActive(context.Background(), mockCustomer)
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
		name       string
		success    bool
		wantErr    error
		customerID customer.CustomerID
		setup      func(mock sqlmock.Sqlmock, customerID customer.CustomerID)
	}{
		{
			name:       "success delete customer",
			success:    true,
			wantErr:    nil,
			customerID: customer.CustomerID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, customerID customer.CustomerID) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "customers" WHERE id = $1`)).
					WithArgs(customerID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:       "failure customer not found",
			success:    false,
			wantErr:    customer.ErrCustomerNotFound,
			customerID: customer.CustomerID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, customerID customer.CustomerID) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "customers" WHERE id = $1`)).
					WithArgs(customerID).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectRollback()
			},
		},
		{
			name:       "failure delete customer error",
			success:    false,
			wantErr:    nil,
			customerID: customer.CustomerID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, customerID customer.CustomerID) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "customers" WHERE id = $1`)).
					WithArgs(customerID).
					WillReturnError(errors.New("delete customer error"))

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

			tt.setup(mock, tt.customerID)

			repo := NewCustomerRepository(gormDB)

			err = repo.Delete(context.Background(), tt.customerID)
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
