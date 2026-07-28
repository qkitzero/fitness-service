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

	"github.com/qkitzero/fitness-service/internal/domain/customer"
	mockscustomer "github.com/qkitzero/fitness-service/mocks/domain/customer"
)

var (
	testBirthDate = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	testCreatedAt = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	testUpdatedAt = time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
)

func TestCreate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, customer customer.Customer)
	}{
		{
			name:    "success create customer",
			success: true,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "customers" ("id","group_id","name","name_kana","gender","birth_date","phone","email","postal_code","prefecture","city","street","building","emergency_contact_name","emergency_contact_relationship","emergency_contact_phone","is_active","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`)).
					WithArgs(customer.ID(), customer.GroupID(), customer.Name(), customer.NameKana(), customer.Gender(), testBirthDate, "0312345678", "test@example.com", "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "緊急 太郎", "父", "09012345678", true, testCreatedAt, testUpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "failure create customer error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "customers" ("id","group_id","name","name_kana","gender","birth_date","phone","email","postal_code","prefecture","city","street","building","emergency_contact_name","emergency_contact_relationship","emergency_contact_phone","is_active","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`)).
					WithArgs(customer.ID(), customer.GroupID(), customer.Name(), customer.NameKana(), customer.Gender(), testBirthDate, "0312345678", "test@example.com", "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "緊急 太郎", "父", "09012345678", true, testCreatedAt, testUpdatedAt).
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

			phone, _ := customer.NewPhone("03-1234-5678")
			email, _ := customer.NewEmail("test@example.com")
			postalCode, _ := customer.NewPostalCode("123-4567")
			prefecture, _ := customer.NewPrefecture("東京都")
			city, _ := customer.NewCity("千代田区")
			street, _ := customer.NewStreet("1-1-1")
			building, _ := customer.NewBuilding("テストビル")
			emergencyContactName, _ := customer.NewEmergencyContactName("緊急 太郎")
			emergencyContactRelationship, _ := customer.NewEmergencyContactRelationship("父")
			emergencyContactPhone, _ := customer.NewPhone("090-1234-5678")

			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().ID().Return(customer.CustomerID{UUID: uuid.New()}).AnyTimes()
			mockCustomer.EXPECT().GroupID().Return(customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
			mockCustomer.EXPECT().Name().Return(customer.Name("test customer")).AnyTimes()
			mockCustomer.EXPECT().NameKana().Return(customer.NameKana("テストカナ")).AnyTimes()
			mockCustomer.EXPECT().Gender().Return(customer.GenderMale).AnyTimes()
			mockCustomer.EXPECT().BirthDate().Return(customer.BirthDate{Time: testBirthDate}).AnyTimes()
			mockCustomer.EXPECT().IsActive().Return(true).AnyTimes()
			mockCustomer.EXPECT().Phone().Return(phone).AnyTimes()
			mockCustomer.EXPECT().Email().Return(email).AnyTimes()
			mockCustomer.EXPECT().PostalCode().Return(postalCode).AnyTimes()
			mockCustomer.EXPECT().Prefecture().Return(prefecture).AnyTimes()
			mockCustomer.EXPECT().City().Return(city).AnyTimes()
			mockCustomer.EXPECT().Street().Return(street).AnyTimes()
			mockCustomer.EXPECT().Building().Return(building).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactName().Return(emergencyContactName).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactRelationship().Return(emergencyContactRelationship).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactPhone().Return(emergencyContactPhone).AnyTimes()
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

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCreateNilOptionals(t *testing.T) {
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
	mockCustomer.EXPECT().GroupID().Return(customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
	mockCustomer.EXPECT().Name().Return(customer.Name("test customer")).AnyTimes()
	mockCustomer.EXPECT().NameKana().Return(customer.NameKana("テストカナ")).AnyTimes()
	mockCustomer.EXPECT().Gender().Return(customer.GenderMale).AnyTimes()
	mockCustomer.EXPECT().BirthDate().Return(customer.BirthDate{Time: testBirthDate}).AnyTimes()
	mockCustomer.EXPECT().IsActive().Return(true).AnyTimes()
	mockCustomer.EXPECT().Phone().Return(nil).AnyTimes()
	mockCustomer.EXPECT().Email().Return(nil).AnyTimes()
	mockCustomer.EXPECT().PostalCode().Return(nil).AnyTimes()
	mockCustomer.EXPECT().Prefecture().Return(nil).AnyTimes()
	mockCustomer.EXPECT().City().Return(nil).AnyTimes()
	mockCustomer.EXPECT().Street().Return(nil).AnyTimes()
	mockCustomer.EXPECT().Building().Return(nil).AnyTimes()
	mockCustomer.EXPECT().EmergencyContactName().Return(nil).AnyTimes()
	mockCustomer.EXPECT().EmergencyContactRelationship().Return(nil).AnyTimes()
	mockCustomer.EXPECT().EmergencyContactPhone().Return(nil).AnyTimes()
	mockCustomer.EXPECT().CreatedAt().Return(testCreatedAt).AnyTimes()
	mockCustomer.EXPECT().UpdatedAt().Return(testUpdatedAt).AnyTimes()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "customers" ("id","group_id","name","name_kana","gender","birth_date","phone","email","postal_code","prefecture","city","street","building","emergency_contact_name","emergency_contact_relationship","emergency_contact_phone","is_active","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`)).
		WithArgs(mockCustomer.ID(), mockCustomer.GroupID(), mockCustomer.Name(), mockCustomer.NameKana(), mockCustomer.Gender(), testBirthDate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, true, testCreatedAt, testUpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewCustomerRepository(gormDB)

	if err := repo.Create(context.Background(), mockCustomer); err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFindByIDNilOptionals(t *testing.T) {
	t.Parallel()
	columns := []string{"id", "group_id", "name", "name_kana", "gender", "birth_date", "phone", "email", "postal_code", "prefecture", "city", "street", "building", "emergency_contact_name", "emergency_contact_relationship", "emergency_contact_phone", "is_active", "created_at", "updated_at"}
	customerID := customer.CustomerID{UUID: uuid.New()}

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to new sqlmock: %s", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %s", err)
	}

	customerRows := sqlmock.NewRows(columns).
		AddRow(customerID, "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "test customer", "テストカナ", "male", testBirthDate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, true, testCreatedAt, testUpdatedAt)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE id = $1 ORDER BY "customers"."id" LIMIT $2`)).
		WithArgs(customerID, 1).
		WillReturnRows(customerRows)

	repo := NewCustomerRepository(gormDB)

	c, err := repo.FindByID(context.Background(), customerID)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	if c.Phone() != nil || c.Email() != nil || c.PostalCode() != nil || c.Prefecture() != nil || c.City() != nil || c.Street() != nil || c.Building() != nil || c.EmergencyContactName() != nil || c.EmergencyContactRelationship() != nil || c.EmergencyContactPhone() != nil {
		t.Errorf("expected nil optional fields")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFindByID(t *testing.T) {
	t.Parallel()
	columns := []string{"id", "group_id", "name", "name_kana", "gender", "birth_date", "phone", "email", "postal_code", "prefecture", "city", "street", "building", "emergency_contact_name", "emergency_contact_relationship", "emergency_contact_phone", "is_active", "created_at", "updated_at"}
	tests := []struct {
		name       string
		success    bool
		wantErr    error
		customerID customer.CustomerID
		setup      func(mock sqlmock.Sqlmock, customerID customer.CustomerID)
	}{
		{
			name:       "success find customer by id",
			success:    true,
			wantErr:    nil,
			customerID: customer.CustomerID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, customerID customer.CustomerID) {
				customerRows := sqlmock.NewRows(columns).
					AddRow(customerID, "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "test customer", "テストカナ", "male", testBirthDate, "0312345678", "test@example.com", "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "緊急 太郎", "父", "09012345678", false, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE id = $1 ORDER BY "customers"."id" LIMIT $2`)).
					WithArgs(customerID, 1).
					WillReturnRows(customerRows)
			},
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
				if c.IsActive() {
					t.Errorf("IsActive() = %v, want %v", c.IsActive(), false)
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

func TestListByGroupID(t *testing.T) {
	t.Parallel()
	columns := []string{"id", "group_id", "name", "name_kana", "gender", "birth_date", "phone", "email", "postal_code", "prefecture", "city", "street", "building", "emergency_contact_name", "emergency_contact_relationship", "emergency_contact_phone", "is_active", "created_at", "updated_at"}
	tests := []struct {
		name            string
		success         bool
		wantNames       []string
		wantActive      []bool
		includeInactive bool
		groupID         customer.GroupID
		setup           func(mock sqlmock.Sqlmock, groupID customer.GroupID)
	}{
		{
			name:            "success list active customers by group id",
			success:         true,
			wantNames:       []string{"test customer 1", "test customer 2"},
			wantActive:      []bool{true, true},
			includeInactive: false,
			groupID:         customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, groupID customer.GroupID) {
				customerRows := sqlmock.NewRows(columns).
					AddRow(uuid.New().String(), groupID, "test customer 1", "テストカナ", "male", testBirthDate, "0312345678", "test1@example.com", "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "緊急 太郎", "父", "09012345678", true, testCreatedAt, testUpdatedAt).
					AddRow(uuid.New().String(), groupID, "test customer 2", "テストカナ", "female", testBirthDate, "0312345679", "test2@example.com", "1234568", "大阪府", "大阪市", "2-2-2", "更新ビル", "緊急 花子", "母", "08012345678", true, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE group_id = $1 AND is_active = $2 ORDER BY created_at, id`)).
					WithArgs(groupID, true).
					WillReturnRows(customerRows)
			},
		},
		{
			name:            "success list customers including inactive",
			success:         true,
			wantNames:       []string{"test customer 1", "test customer 2"},
			wantActive:      []bool{true, false},
			includeInactive: true,
			groupID:         customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, groupID customer.GroupID) {
				customerRows := sqlmock.NewRows(columns).
					AddRow(uuid.New().String(), groupID, "test customer 1", "テストカナ", "male", testBirthDate, "0312345678", "test1@example.com", "1234567", "東京都", "千代田区", "1-1-1", "テストビル", "緊急 太郎", "父", "09012345678", true, testCreatedAt, testUpdatedAt).
					AddRow(uuid.New().String(), groupID, "test customer 2", "テストカナ", "female", testBirthDate, "0312345679", "test2@example.com", "1234568", "大阪府", "大阪市", "2-2-2", "更新ビル", "緊急 花子", "母", "08012345678", false, testCreatedAt, testUpdatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE group_id = $1 ORDER BY created_at, id`)).
					WithArgs(groupID).
					WillReturnRows(customerRows)
			},
		},
		{
			name:            "success list no customers",
			success:         true,
			wantNames:       []string{},
			wantActive:      []bool{},
			includeInactive: false,
			groupID:         customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, groupID customer.GroupID) {
				customerRows := sqlmock.NewRows(columns)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE group_id = $1 AND is_active = $2 ORDER BY created_at, id`)).
					WithArgs(groupID, true).
					WillReturnRows(customerRows)
			},
		},
		{
			name:            "failure list customers error",
			success:         false,
			wantNames:       nil,
			wantActive:      nil,
			includeInactive: false,
			groupID:         customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, groupID customer.GroupID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE group_id = $1 AND is_active = $2 ORDER BY created_at, id`)).
					WithArgs(groupID, true).
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

			tt.setup(mock, tt.groupID)

			repo := NewCustomerRepository(gormDB)

			customers, err := repo.ListByGroupID(context.Background(), tt.groupID, tt.includeInactive)
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
					if customers[i].GroupID() != tt.groupID {
						t.Errorf("customers[%d].GroupID() = %v, want %v", i, customers[i].GroupID(), tt.groupID)
					}
					if customers[i].IsActive() != tt.wantActive[i] {
						t.Errorf("customers[%d].IsActive() = %v, want %v", i, customers[i].IsActive(), tt.wantActive[i])
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
	tests := []struct {
		name    string
		success bool
		wantErr error
		setup   func(mock sqlmock.Sqlmock, customer customer.Customer)
	}{
		{
			name:    "success update customer",
			success: true,
			wantErr: nil,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "group_id"=$1,"name"=$2,"name_kana"=$3,"gender"=$4,"birth_date"=$5,"phone"=$6,"email"=$7,"postal_code"=$8,"prefecture"=$9,"city"=$10,"street"=$11,"building"=$12,"emergency_contact_name"=$13,"emergency_contact_relationship"=$14,"emergency_contact_phone"=$15,"created_at"=$16,"updated_at"=$17 WHERE id = $18`)).
					WithArgs(customer.GroupID(), customer.Name(), customer.NameKana(), customer.Gender(), updatedBirthDate, "08012345678", "updated@example.com", "5432100", "大阪府", "大阪市", "2-2-2", "更新ビル", "更新 花子", "母", "07012345678", testCreatedAt, testUpdatedAt, customer.ID()).
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

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "group_id"=$1,"name"=$2,"name_kana"=$3,"gender"=$4,"birth_date"=$5,"phone"=$6,"email"=$7,"postal_code"=$8,"prefecture"=$9,"city"=$10,"street"=$11,"building"=$12,"emergency_contact_name"=$13,"emergency_contact_relationship"=$14,"emergency_contact_phone"=$15,"created_at"=$16,"updated_at"=$17 WHERE id = $18`)).
					WithArgs(customer.GroupID(), customer.Name(), customer.NameKana(), customer.Gender(), updatedBirthDate, "08012345678", "updated@example.com", "5432100", "大阪府", "大阪市", "2-2-2", "更新ビル", "更新 花子", "母", "07012345678", testCreatedAt, testUpdatedAt, customer.ID()).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectRollback()
			},
		},
		{
			name:    "failure update customer error",
			success: false,
			wantErr: nil,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "group_id"=$1,"name"=$2,"name_kana"=$3,"gender"=$4,"birth_date"=$5,"phone"=$6,"email"=$7,"postal_code"=$8,"prefecture"=$9,"city"=$10,"street"=$11,"building"=$12,"emergency_contact_name"=$13,"emergency_contact_relationship"=$14,"emergency_contact_phone"=$15,"created_at"=$16,"updated_at"=$17 WHERE id = $18`)).
					WithArgs(customer.GroupID(), customer.Name(), customer.NameKana(), customer.Gender(), updatedBirthDate, "08012345678", "updated@example.com", "5432100", "大阪府", "大阪市", "2-2-2", "更新ビル", "更新 花子", "母", "07012345678", testCreatedAt, testUpdatedAt, customer.ID()).
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

			phone, _ := customer.NewPhone("080-1234-5678")
			email, _ := customer.NewEmail("updated@example.com")
			postalCode, _ := customer.NewPostalCode("543-2100")
			prefecture, _ := customer.NewPrefecture("大阪府")
			city, _ := customer.NewCity("大阪市")
			street, _ := customer.NewStreet("2-2-2")
			building, _ := customer.NewBuilding("更新ビル")
			emergencyContactName, _ := customer.NewEmergencyContactName("更新 花子")
			emergencyContactRelationship, _ := customer.NewEmergencyContactRelationship("母")
			emergencyContactPhone, _ := customer.NewPhone("070-1234-5678")

			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().ID().Return(customer.CustomerID{UUID: uuid.New()}).AnyTimes()
			mockCustomer.EXPECT().GroupID().Return(customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
			mockCustomer.EXPECT().Name().Return(customer.Name("updated customer")).AnyTimes()
			mockCustomer.EXPECT().NameKana().Return(customer.NameKana("コウシンカナ")).AnyTimes()
			mockCustomer.EXPECT().Gender().Return(customer.GenderFemale).AnyTimes()
			mockCustomer.EXPECT().BirthDate().Return(customer.BirthDate{Time: updatedBirthDate}).AnyTimes()
			mockCustomer.EXPECT().IsActive().Return(false).AnyTimes()
			mockCustomer.EXPECT().Phone().Return(phone).AnyTimes()
			mockCustomer.EXPECT().Email().Return(email).AnyTimes()
			mockCustomer.EXPECT().PostalCode().Return(postalCode).AnyTimes()
			mockCustomer.EXPECT().Prefecture().Return(prefecture).AnyTimes()
			mockCustomer.EXPECT().City().Return(city).AnyTimes()
			mockCustomer.EXPECT().Street().Return(street).AnyTimes()
			mockCustomer.EXPECT().Building().Return(building).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactName().Return(emergencyContactName).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactRelationship().Return(emergencyContactRelationship).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactPhone().Return(emergencyContactPhone).AnyTimes()
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

func TestUpdateNilOptionals(t *testing.T) {
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
	mockCustomer.EXPECT().GroupID().Return(customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
	mockCustomer.EXPECT().Name().Return(customer.Name("updated customer")).AnyTimes()
	mockCustomer.EXPECT().NameKana().Return(customer.NameKana("コウシンカナ")).AnyTimes()
	mockCustomer.EXPECT().Gender().Return(customer.GenderFemale).AnyTimes()
	mockCustomer.EXPECT().BirthDate().Return(customer.BirthDate{Time: testBirthDate}).AnyTimes()
	mockCustomer.EXPECT().IsActive().Return(true).AnyTimes()
	mockCustomer.EXPECT().Phone().Return(nil).AnyTimes()
	mockCustomer.EXPECT().Email().Return(nil).AnyTimes()
	mockCustomer.EXPECT().PostalCode().Return(nil).AnyTimes()
	mockCustomer.EXPECT().Prefecture().Return(nil).AnyTimes()
	mockCustomer.EXPECT().City().Return(nil).AnyTimes()
	mockCustomer.EXPECT().Street().Return(nil).AnyTimes()
	mockCustomer.EXPECT().Building().Return(nil).AnyTimes()
	mockCustomer.EXPECT().EmergencyContactName().Return(nil).AnyTimes()
	mockCustomer.EXPECT().EmergencyContactRelationship().Return(nil).AnyTimes()
	mockCustomer.EXPECT().EmergencyContactPhone().Return(nil).AnyTimes()
	mockCustomer.EXPECT().CreatedAt().Return(testCreatedAt).AnyTimes()
	mockCustomer.EXPECT().UpdatedAt().Return(testUpdatedAt).AnyTimes()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "group_id"=$1,"name"=$2,"name_kana"=$3,"gender"=$4,"birth_date"=$5,"phone"=$6,"email"=$7,"postal_code"=$8,"prefecture"=$9,"city"=$10,"street"=$11,"building"=$12,"emergency_contact_name"=$13,"emergency_contact_relationship"=$14,"emergency_contact_phone"=$15,"created_at"=$16,"updated_at"=$17 WHERE id = $18`)).
		WithArgs(mockCustomer.GroupID(), mockCustomer.Name(), mockCustomer.NameKana(), mockCustomer.Gender(), testBirthDate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, testCreatedAt, testUpdatedAt, mockCustomer.ID()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewCustomerRepository(gormDB)

	if err := repo.Update(context.Background(), mockCustomer); err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
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
			mockCustomer.EXPECT().GroupID().Return(customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
			mockCustomer.EXPECT().Name().Return(customer.Name("test customer")).AnyTimes()
			mockCustomer.EXPECT().NameKana().Return(customer.NameKana("テストカナ")).AnyTimes()
			mockCustomer.EXPECT().Gender().Return(customer.GenderMale).AnyTimes()
			mockCustomer.EXPECT().BirthDate().Return(customer.BirthDate{Time: testBirthDate}).AnyTimes()
			mockCustomer.EXPECT().IsActive().Return(false).AnyTimes()
			mockCustomer.EXPECT().Phone().Return(nil).AnyTimes()
			mockCustomer.EXPECT().Email().Return(nil).AnyTimes()
			mockCustomer.EXPECT().PostalCode().Return(nil).AnyTimes()
			mockCustomer.EXPECT().Prefecture().Return(nil).AnyTimes()
			mockCustomer.EXPECT().City().Return(nil).AnyTimes()
			mockCustomer.EXPECT().Street().Return(nil).AnyTimes()
			mockCustomer.EXPECT().Building().Return(nil).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactName().Return(nil).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactRelationship().Return(nil).AnyTimes()
			mockCustomer.EXPECT().EmergencyContactPhone().Return(nil).AnyTimes()
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
