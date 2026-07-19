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
	"github.com/qkitzero/fitness-service/testutil"
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

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "customers" ("id","group_id","name","created_at","updated_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(customer.ID(), customer.GroupID(), customer.Name(), testutil.AnyTime{}, testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "failure create customer error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "customers" ("id","group_id","name","created_at","updated_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(customer.ID(), customer.GroupID(), customer.Name(), testutil.AnyTime{}, testutil.AnyTime{}).
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
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().ID().Return(customer.CustomerID{UUID: uuid.New()}).AnyTimes()
			mockCustomer.EXPECT().GroupID().Return(customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
			mockCustomer.EXPECT().Name().Return(customer.Name("test customer")).AnyTimes()
			mockCustomer.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
			mockCustomer.EXPECT().UpdatedAt().Return(time.Now()).AnyTimes()

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

func TestFindByID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		success    bool
		customerID customer.CustomerID
		setup      func(mock sqlmock.Sqlmock, customerID customer.CustomerID)
	}{
		{
			name:       "success find customer by id",
			success:    true,
			customerID: customer.CustomerID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, customerID customer.CustomerID) {
				customerRows := sqlmock.NewRows([]string{"id", "group_id", "name", "created_at", "updated_at"}).
					AddRow(customerID, "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "test customer", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE id = $1 ORDER BY "customers"."id" LIMIT $2`)).
					WithArgs(customerID, 1).
					WillReturnRows(customerRows)
			},
		},
		{
			name:       "failure customer not found",
			success:    false,
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
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			tt.setup(mock, tt.customerID)

			repo := NewCustomerRepository(gormDB)

			_, err = repo.FindByID(context.Background(), tt.customerID)
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

func TestListByGroupID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		success   bool
		wantCount int
		groupID   customer.GroupID
		setup     func(mock sqlmock.Sqlmock, groupID customer.GroupID)
	}{
		{
			name:      "success list customers by group id",
			success:   true,
			wantCount: 2,
			groupID:   customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, groupID customer.GroupID) {
				customerRows := sqlmock.NewRows([]string{"id", "group_id", "name", "created_at", "updated_at"}).
					AddRow(uuid.New().String(), groupID, "test customer 1", time.Now(), time.Now()).
					AddRow(uuid.New().String(), groupID, "test customer 2", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE group_id = $1 ORDER BY created_at, id`)).
					WithArgs(groupID).
					WillReturnRows(customerRows)
			},
		},
		{
			name:      "success list no customers",
			success:   true,
			wantCount: 0,
			groupID:   customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, groupID customer.GroupID) {
				customerRows := sqlmock.NewRows([]string{"id", "group_id", "name", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE group_id = $1 ORDER BY created_at, id`)).
					WithArgs(groupID).
					WillReturnRows(customerRows)
			},
		},
		{
			name:      "failure list customers error",
			success:   false,
			wantCount: 0,
			groupID:   customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"),
			setup: func(mock sqlmock.Sqlmock, groupID customer.GroupID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE group_id = $1 ORDER BY created_at, id`)).
					WithArgs(groupID).
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
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			tt.setup(mock, tt.groupID)

			repo := NewCustomerRepository(gormDB)

			customers, err := repo.ListByGroupID(context.Background(), tt.groupID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && len(customers) != tt.wantCount {
				t.Errorf("len(customers) = %v, want %v", len(customers), tt.wantCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, customer customer.Customer)
	}{
		{
			name:    "success update customer",
			success: true,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "group_id"=$1,"name"=$2,"created_at"=$3,"updated_at"=$4 WHERE "id" = $5`)).
					WithArgs(customer.GroupID(), customer.Name(), testutil.AnyTime{}, testutil.AnyTime{}, customer.ID()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "failure update customer error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, customer customer.Customer) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "group_id"=$1,"name"=$2,"created_at"=$3,"updated_at"=$4 WHERE "id" = $5`)).
					WithArgs(customer.GroupID(), customer.Name(), testutil.AnyTime{}, testutil.AnyTime{}, customer.ID()).
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
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCustomer := mockscustomer.NewMockCustomer(ctrl)
			mockCustomer.EXPECT().ID().Return(customer.CustomerID{UUID: uuid.New()}).AnyTimes()
			mockCustomer.EXPECT().GroupID().Return(customer.GroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")).AnyTimes()
			mockCustomer.EXPECT().Name().Return(customer.Name("updated customer")).AnyTimes()
			mockCustomer.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
			mockCustomer.EXPECT().UpdatedAt().Return(time.Now()).AnyTimes()

			tt.setup(mock, mockCustomer)

			repo := NewCustomerRepository(gormDB)

			err = repo.Update(context.Background(), mockCustomer)
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

func TestDelete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		success    bool
		customerID customer.CustomerID
		setup      func(mock sqlmock.Sqlmock, customerID customer.CustomerID)
	}{
		{
			name:       "success delete customer",
			success:    true,
			customerID: customer.CustomerID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, customerID customer.CustomerID) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "customers" WHERE id = $1`)).
					WithArgs(customerID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:       "failure delete customer error",
			success:    false,
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
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
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

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
