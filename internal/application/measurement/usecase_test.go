package measurement

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/staff"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
	mocksappauth "github.com/qkitzero/fitness-service/mocks/application/auth"
	mocksappuser "github.com/qkitzero/fitness-service/mocks/application/user"
	mockscustomer "github.com/qkitzero/fitness-service/mocks/domain/customer"
	mocksmeasurement "github.com/qkitzero/fitness-service/mocks/domain/measurement"
	mocksmeasurementitem "github.com/qkitzero/fitness-service/mocks/domain/measurementitem"
)

func TestCreateMeasurement(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	otherTenantID := "9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	userID := "google-oauth2|000000000000000000000"
	measuredBy, _ := staff.NewStaffID("google-oauth2|999999999999999999999")
	measuredOn, _ := measurement.NewMeasuredOn(2026, 8, 1)
	beforeBirthMeasuredOn, _ := measurement.NewMeasuredOn(2025, 6, 1)
	birthDate, _ := customer.NewBirthDate(1960, 5, 20)
	laterBirthDate, _ := customer.NewBirthDate(2026, 1, 1)

	itemCreatedAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	itemUpdatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	motorFunction, _ := measurementitem.NewCategory("motor_function")
	kg, _ := measurementitem.NewUnit("kg")
	twoTrials, _ := measurementitem.NewTrialCount(2)
	gripStrengthCode, _ := measurementitem.NewCode("grip_strength")
	gripStrengthName, _ := measurementitem.NewName("握力")
	gripStrength := measurementitem.NewMeasurementItem(measurementitem.NewMeasurementItemID(), gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, measurementitem.SideModeBilateral, measurementitem.ValueTypeNumeric, nil, measurementitem.SideAggregationMean, nil, itemCreatedAt, itemUpdatedAt)

	completeEntryInputs := func() []MeasurementEntryInput {
		firstTrial, _ := measurement.NewTrialIndex(1)
		secondTrial, _ := measurement.NewTrialIndex(2)
		firstLeft, _ := measurement.NewValue(32.4)
		firstRight, _ := measurement.NewValue(33.1)
		secondLeft, _ := measurement.NewValue(31.8)
		secondRight, _ := measurement.NewValue(33.5)
		return []MeasurementEntryInput{{
			MeasurementItemID: gripStrength.ID(),
			Values: []measurement.MeasurementValue{
				measurement.NewMeasurementValue(firstTrial, measurement.SideLeft, &firstLeft, nil, nil),
				measurement.NewMeasurementValue(firstTrial, measurement.SideRight, &firstRight, nil, nil),
				measurement.NewMeasurementValue(secondTrial, measurement.SideLeft, &secondLeft, nil, nil),
				measurement.NewMeasurementValue(secondTrial, measurement.SideRight, &secondRight, nil, nil),
			},
		}}
	}

	tests := []struct {
		name             string
		success          bool
		wantErr          error
		ctx              context.Context
		userID           string
		verifyTokenErr   error
		callFindCustomer bool
		birthDate        customer.BirthDate
		findCustomerErr  error
		myTenantIDs      []string
		listMyGroupsErr  error
		measuredOn       measurement.MeasuredOn
		isDraft          bool
		entryInputs      func() []MeasurementEntryInput
		callFindByIDs    bool
		measurementItems []measurementitem.MeasurementItem
		findByIDsErr     error
		callCreate       bool
		createErr        error
	}{
		{
			name:             "success create draft measurement without entries",
			success:          true,
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			isDraft:          true,
			entryInputs:      func() []MeasurementEntryInput { return nil },
			callCreate:       true,
		},
		{
			name:             "success create measurement with entries",
			success:          true,
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs:      completeEntryInputs,
			callFindByIDs:    true,
			measurementItems: []measurementitem.MeasurementItem{gripStrength},
			callCreate:       true,
		},
		{
			name:             "success create draft measurement with partial entries",
			success:          true,
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			isDraft:          true,
			entryInputs: func() []MeasurementEntryInput {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(32.4)
				return []MeasurementEntryInput{{
					MeasurementItemID: gripStrength.ID(),
					Values: []measurement.MeasurementValue{
						measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
					},
				}}
			},
			callFindByIDs:    true,
			measurementItems: []measurementitem.MeasurementItem{gripStrength},
			callCreate:       true,
		},
		{
			name:           "failure verify token error",
			ctx:            context.Background(),
			userID:         "",
			verifyTokenErr: errors.New("verify token error"),
			measuredOn:     measuredOn,
			entryInputs:    func() []MeasurementEntryInput { return nil },
		},
		{
			name:        "failure invalid staff id from token",
			ctx:         context.Background(),
			userID:      "   ",
			measuredOn:  measuredOn,
			entryInputs: func() []MeasurementEntryInput { return nil },
		},
		{
			name:             "failure customer not found",
			wantErr:          customer.ErrCustomerNotFound,
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			findCustomerErr:  customer.ErrCustomerNotFound,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs:      func() []MeasurementEntryInput { return nil },
		},
		{
			name:             "failure other tenant is hidden as not found",
			wantErr:          customer.ErrCustomerNotFound,
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{otherTenantID},
			measuredOn:       measuredOn,
			entryInputs:      func() []MeasurementEntryInput { return nil },
		},
		{
			name:             "failure list my groups error",
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			listMyGroupsErr:  errors.New("list my groups error"),
			measuredOn:       measuredOn,
			entryInputs:      func() []MeasurementEntryInput { return nil },
		},
		{
			name:             "failure measured on before the birth date",
			wantErr:          measurement.ErrInvalidAgeAtMeasurement,
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        laterBirthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       beforeBirthMeasuredOn,
			entryInputs:      func() []MeasurementEntryInput { return nil },
		},
		{
			name:             "failure find measurement items error",
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs:      completeEntryInputs,
			callFindByIDs:    true,
			findByIDsErr:     errors.New("find by ids error"),
		},
		{
			name:             "failure unknown measurement item",
			wantErr:          measurementitem.ErrMeasurementItemNotFound,
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs:      completeEntryInputs,
			callFindByIDs:    true,
			measurementItems: nil,
		},
		{
			name:             "failure invalid side for the measurement item",
			wantErr:          measurement.ErrInvalidSide,
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs: func() []MeasurementEntryInput {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(32.4)
				return []MeasurementEntryInput{{
					MeasurementItemID: gripStrength.ID(),
					Values: []measurement.MeasurementValue{
						measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil),
					},
				}}
			},
			callFindByIDs:    true,
			measurementItems: []measurementitem.MeasurementItem{gripStrength},
		},
		{
			name:             "success create confirmed measurement with fewer values than the trial capacity",
			success:          true,
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs: func() []MeasurementEntryInput {
				trialIndex, _ := measurement.NewTrialIndex(1)
				left, _ := measurement.NewValue(32.4)
				right, _ := measurement.NewValue(33.1)
				return []MeasurementEntryInput{{
					MeasurementItemID: gripStrength.ID(),
					Values: []measurement.MeasurementValue{
						measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &left, nil, nil),
						measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &right, nil, nil),
					},
				}}
			},
			callFindByIDs:    true,
			measurementItems: []measurementitem.MeasurementItem{gripStrength},
			callCreate:       true,
		},
		{
			name:             "failure no values on a confirmed measurable entry",
			wantErr:          measurement.ErrInvalidValueCount,
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs: func() []MeasurementEntryInput {
				return []MeasurementEntryInput{{MeasurementItemID: gripStrength.ID()}}
			},
			callFindByIDs:    true,
			measurementItems: []measurementitem.MeasurementItem{gripStrength},
		},
		{
			name:             "failure no entries on a confirmed measurement",
			wantErr:          measurement.ErrInvalidValueCount,
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs:      func() []MeasurementEntryInput { return nil },
		},
		{
			name:             "failure create error",
			ctx:              context.Background(),
			userID:           userID,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			isDraft:          true,
			entryInputs:      func() []MeasurementEntryInput { return nil },
			callCreate:       true,
			createErr:        errors.New("create error"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			customerID := customer.NewCustomerID()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockMeasurementRepository := mocksmeasurement.NewMockMeasurementRepository(ctrl)
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockMeasurementItemRepository := mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindCustomer {
				var foundCustomer customer.Customer
				if tt.findCustomerErr == nil {
					mockCustomer := mockscustomer.NewMockCustomer(ctrl)
					mockCustomer.EXPECT().TenantID().Return(tenantID).AnyTimes()
					mockCustomer.EXPECT().BirthDate().Return(tt.birthDate).AnyTimes()
					foundCustomer = mockCustomer
				}
				mockCustomerRepository.EXPECT().FindByID(tt.ctx, customerID).Return(foundCustomer, tt.findCustomerErr).Times(1)
			}
			if tt.callFindByIDs {
				mockMeasurementItemRepository.EXPECT().FindByIDs(tt.ctx, gomock.Any()).Return(tt.measurementItems, tt.findByIDsErr).Times(1)
			}
			if tt.callCreate {
				mockMeasurementRepository.EXPECT().Create(tt.ctx, gomock.Any()).Return(tt.createErr).Times(1)
			}

			u := NewMeasurementUsecase(mockAuthService, mockUserService, mockMeasurementRepository, mockCustomerRepository, mockMeasurementItemRepository)

			entryInputs := tt.entryInputs()

			createdMeasurement, err := u.CreateMeasurement(tt.ctx, customerID, tt.measuredOn, measuredBy, tt.isDraft, entryInputs)
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
				if createdMeasurement.ID().UUID == uuid.Nil {
					t.Errorf("expected generated measurement id, but got a nil UUID")
				}
				if createdMeasurement.CustomerID() != customerID {
					t.Errorf("CustomerID() = %v, want %v", createdMeasurement.CustomerID(), customerID)
				}
				if createdMeasurement.MeasuredBy() != measuredBy {
					t.Errorf("MeasuredBy() = %v, want %v", createdMeasurement.MeasuredBy(), measuredBy)
				}
				if createdMeasurement.UpdatedBy().String() != tt.userID {
					t.Errorf("UpdatedBy() = %v, want %v", createdMeasurement.UpdatedBy(), tt.userID)
				}
				if createdMeasurement.AgeAtMeasurement().Int() != tt.birthDate.AgeAt(tt.measuredOn.Time) {
					t.Errorf("AgeAtMeasurement() = %v, want %v", createdMeasurement.AgeAtMeasurement().Int(), tt.birthDate.AgeAt(tt.measuredOn.Time))
				}
				if createdMeasurement.IsDraft() != tt.isDraft {
					t.Errorf("IsDraft() = %v, want %v", createdMeasurement.IsDraft(), tt.isDraft)
				}
				if len(createdMeasurement.Entries()) != len(entryInputs) {
					t.Errorf("len(Entries()) = %v, want %v", len(createdMeasurement.Entries()), len(entryInputs))
				}
				if !createdMeasurement.CreatedAt().Equal(createdMeasurement.UpdatedAt()) {
					t.Errorf("CreatedAt() = %v, UpdatedAt() = %v, want equal", createdMeasurement.CreatedAt(), createdMeasurement.UpdatedAt())
				}
			}
		})
	}
}

func TestGetMeasurement(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	otherTenantID := "9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	userID := "google-oauth2|000000000000000000000"

	tests := []struct {
		name             string
		success          bool
		wantErr          error
		ctx              context.Context
		userID           string
		verifyTokenErr   error
		callFindByID     bool
		findByIDErr      error
		callFindCustomer bool
		findCustomerErr  error
		myTenantIDs      []string
		listMyGroupsErr  error
	}{
		{"success get measurement", true, nil, context.Background(), userID, nil, true, nil, true, nil, []string{tenantID.String()}, nil},
		{"failure verify token error", false, nil, context.Background(), "", errors.New("verify token error"), false, nil, false, nil, []string{tenantID.String()}, nil},
		{"failure find by id error", false, nil, context.Background(), userID, nil, true, errors.New("find by id error"), false, nil, []string{tenantID.String()}, nil},
		{"failure measurement not found", false, measurement.ErrMeasurementNotFound, context.Background(), userID, nil, true, measurement.ErrMeasurementNotFound, false, nil, []string{tenantID.String()}, nil},
		{"failure customer not found is hidden as measurement not found", false, measurement.ErrMeasurementNotFound, context.Background(), userID, nil, true, nil, true, customer.ErrCustomerNotFound, []string{tenantID.String()}, nil},
		{"failure other tenant is hidden as measurement not found", false, measurement.ErrMeasurementNotFound, context.Background(), userID, nil, true, nil, true, nil, []string{otherTenantID}, nil},
		{"failure list my groups error", false, nil, context.Background(), userID, nil, true, nil, true, nil, nil, errors.New("list my groups error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			measurementID := measurement.NewMeasurementID()
			customerID := customer.NewCustomerID()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockMeasurementRepository := mocksmeasurement.NewMockMeasurementRepository(ctrl)
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockMeasurement := mocksmeasurement.NewMockMeasurement(ctrl)
			mockMeasurement.EXPECT().CustomerID().Return(customerID).AnyTimes()
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindByID {
				var foundMeasurement measurement.Measurement
				if tt.findByIDErr == nil {
					foundMeasurement = mockMeasurement
				}
				mockMeasurementRepository.EXPECT().FindByID(tt.ctx, measurementID).Return(foundMeasurement, tt.findByIDErr).Times(1)
			}
			if tt.callFindCustomer {
				var foundCustomer customer.Customer
				if tt.findCustomerErr == nil {
					mockCustomer := mockscustomer.NewMockCustomer(ctrl)
					mockCustomer.EXPECT().TenantID().Return(tenantID).AnyTimes()
					foundCustomer = mockCustomer
				}
				mockCustomerRepository.EXPECT().FindByID(tt.ctx, customerID).Return(foundCustomer, tt.findCustomerErr).Times(1)
			}

			u := NewMeasurementUsecase(mockAuthService, mockUserService, mockMeasurementRepository, mockCustomerRepository, mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl))

			foundMeasurement, err := u.GetMeasurement(tt.ctx, measurementID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.success && foundMeasurement != mockMeasurement {
				t.Errorf("expected the measurement returned by the repository")
			}
		})
	}
}

func TestListMeasurements(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	otherTenantID := "9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	userID := "google-oauth2|000000000000000000000"

	tests := []struct {
		name                 string
		success              bool
		wantErr              error
		ctx                  context.Context
		userID               string
		verifyTokenErr       error
		callFindCustomer     bool
		findCustomerErr      error
		myTenantIDs          []string
		listMyGroupsErr      error
		callListByCustomerID bool
		measurementCount     int
		listByCustomerIDErr  error
	}{
		{"success list measurements", true, nil, context.Background(), userID, nil, true, nil, []string{tenantID.String()}, nil, true, 2, nil},
		{"success list no measurements", true, nil, context.Background(), userID, nil, true, nil, []string{tenantID.String()}, nil, true, 0, nil},
		{"failure verify token error", false, nil, context.Background(), "", errors.New("verify token error"), false, nil, []string{tenantID.String()}, nil, false, 0, nil},
		{"failure customer not found", false, customer.ErrCustomerNotFound, context.Background(), userID, nil, true, customer.ErrCustomerNotFound, []string{tenantID.String()}, nil, false, 0, nil},
		{"failure other tenant is hidden as not found", false, customer.ErrCustomerNotFound, context.Background(), userID, nil, true, nil, []string{otherTenantID}, nil, false, 0, nil},
		{"failure list my groups error", false, nil, context.Background(), userID, nil, true, nil, nil, errors.New("list my groups error"), false, 0, nil},
		{"failure list by customer id error", false, nil, context.Background(), userID, nil, true, nil, []string{tenantID.String()}, nil, true, 0, errors.New("list by customer id error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			customerID := customer.NewCustomerID()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockMeasurementRepository := mocksmeasurement.NewMockMeasurementRepository(ctrl)
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindCustomer {
				var foundCustomer customer.Customer
				if tt.findCustomerErr == nil {
					mockCustomer := mockscustomer.NewMockCustomer(ctrl)
					mockCustomer.EXPECT().TenantID().Return(tenantID).AnyTimes()
					foundCustomer = mockCustomer
				}
				mockCustomerRepository.EXPECT().FindByID(tt.ctx, customerID).Return(foundCustomer, tt.findCustomerErr).Times(1)
			}
			measurements := make([]measurement.Measurement, 0, tt.measurementCount)
			for range tt.measurementCount {
				measurements = append(measurements, mocksmeasurement.NewMockMeasurement(ctrl))
			}
			if tt.callListByCustomerID {
				if tt.listByCustomerIDErr != nil {
					measurements = nil
				}
				mockMeasurementRepository.EXPECT().ListByCustomerID(tt.ctx, customerID).Return(measurements, tt.listByCustomerIDErr).Times(1)
			}

			u := NewMeasurementUsecase(mockAuthService, mockUserService, mockMeasurementRepository, mockCustomerRepository, mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl))

			foundMeasurements, err := u.ListMeasurements(tt.ctx, customerID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.success && len(foundMeasurements) != tt.measurementCount {
				t.Errorf("len(measurements) = %v, want %v", len(foundMeasurements), tt.measurementCount)
			}
		})
	}
}

func TestUpdateMeasurement(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	otherTenantID := "9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	userID := "google-oauth2|000000000000000000000"
	measuredBy, _ := staff.NewStaffID("google-oauth2|999999999999999999999")
	measuredOn, _ := measurement.NewMeasuredOn(2026, 8, 1)
	storedMeasuredOn, _ := measurement.NewMeasuredOn(2026, 7, 1)
	beforeBirthMeasuredOn, _ := measurement.NewMeasuredOn(2025, 6, 1)
	birthDate, _ := customer.NewBirthDate(1960, 5, 20)
	laterBirthDate, _ := customer.NewBirthDate(2026, 1, 1)

	itemCreatedAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	itemUpdatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	vital, _ := measurementitem.NewCategory("vital")
	bpm, _ := measurementitem.NewUnit("bpm")
	singleTrial, _ := measurementitem.NewTrialCount(1)
	pulseRateCode, _ := measurementitem.NewCode("pulse_rate")
	pulseRateName, _ := measurementitem.NewName("脈拍")
	pulseRate := measurementitem.NewMeasurementItem(measurementitem.NewMeasurementItemID(), pulseRateCode, pulseRateName, vital, bpm, singleTrial, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, nil, measurementitem.SideAggregationMean, nil, itemCreatedAt, itemUpdatedAt)

	completeEntryInputs := func() []MeasurementEntryInput {
		trialIndex, _ := measurement.NewTrialIndex(1)
		value, _ := measurement.NewValue(72)
		return []MeasurementEntryInput{{
			MeasurementItemID: pulseRate.ID(),
			Values: []measurement.MeasurementValue{
				measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil),
			},
		}}
	}

	tests := []struct {
		name             string
		success          bool
		wantErr          error
		ctx              context.Context
		userID           string
		verifyTokenErr   error
		callFindByID     bool
		findByIDErr      error
		callFindCustomer bool
		findCustomerErr  error
		birthDate        customer.BirthDate
		myTenantIDs      []string
		listMyGroupsErr  error
		measuredOn       measurement.MeasuredOn
		wantAge          int
		entryInputs      func() []MeasurementEntryInput
		callFindByIDs    bool
		measurementItems []measurementitem.MeasurementItem
		findByIDsErr     error
		callUpdate       bool
		updateErr        error
	}{
		{
			name:             "success update recomputes the age when the measurement date changes",
			success:          true,
			ctx:              context.Background(),
			userID:           userID,
			callFindByID:     true,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			wantAge:          66,
			entryInputs:      completeEntryInputs,
			callFindByIDs:    true,
			measurementItems: []measurementitem.MeasurementItem{pulseRate},
			callUpdate:       true,
		},
		{
			name:             "success update keeps the stored age when the measurement date is unchanged",
			success:          true,
			ctx:              context.Background(),
			userID:           userID,
			callFindByID:     true,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       storedMeasuredOn,
			wantAge:          70,
			entryInputs:      completeEntryInputs,
			callFindByIDs:    true,
			measurementItems: []measurementitem.MeasurementItem{pulseRate},
			callUpdate:       true,
		},
		{
			name:             "failure list my groups error",
			ctx:              context.Background(),
			userID:           userID,
			callFindByID:     true,
			callFindCustomer: true,
			birthDate:        birthDate,
			listMyGroupsErr:  errors.New("list my groups error"),
			measuredOn:       measuredOn,
			entryInputs:      completeEntryInputs,
		},
		{
			name:             "failure find customer error",
			ctx:              context.Background(),
			userID:           userID,
			callFindByID:     true,
			callFindCustomer: true,
			findCustomerErr:  errors.New("find customer error"),
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs:      completeEntryInputs,
		},
		{
			name:             "failure customer not found is hidden as measurement not found",
			wantErr:          measurement.ErrMeasurementNotFound,
			ctx:              context.Background(),
			userID:           userID,
			callFindByID:     true,
			callFindCustomer: true,
			findCustomerErr:  customer.ErrCustomerNotFound,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs:      completeEntryInputs,
		},
		{
			name:           "failure verify token error",
			ctx:            context.Background(),
			userID:         "",
			verifyTokenErr: errors.New("verify token error"),
			measuredOn:     measuredOn,
			entryInputs:    completeEntryInputs,
		},
		{
			name:         "failure measurement not found",
			wantErr:      measurement.ErrMeasurementNotFound,
			ctx:          context.Background(),
			userID:       userID,
			callFindByID: true,
			findByIDErr:  measurement.ErrMeasurementNotFound,
			myTenantIDs:  []string{tenantID.String()},
			measuredOn:   measuredOn,
			entryInputs:  completeEntryInputs,
		},
		{
			name:             "failure other tenant is hidden as measurement not found",
			wantErr:          measurement.ErrMeasurementNotFound,
			ctx:              context.Background(),
			userID:           userID,
			callFindByID:     true,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{otherTenantID},
			measuredOn:       measuredOn,
			entryInputs:      completeEntryInputs,
		},
		{
			name:             "failure unknown measurement item",
			wantErr:          measurementitem.ErrMeasurementItemNotFound,
			ctx:              context.Background(),
			userID:           userID,
			callFindByID:     true,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs:      completeEntryInputs,
			callFindByIDs:    true,
			measurementItems: nil,
		},
		{
			name:             "failure partial entries on a confirmed measurement",
			wantErr:          measurement.ErrInvalidValueCount,
			ctx:              context.Background(),
			userID:           userID,
			callFindByID:     true,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs: func() []MeasurementEntryInput {
				return []MeasurementEntryInput{{MeasurementItemID: pulseRate.ID()}}
			},
			callFindByIDs:    true,
			measurementItems: []measurementitem.MeasurementItem{pulseRate},
		},
		{
			name:             "failure update error",
			ctx:              context.Background(),
			userID:           userID,
			callFindByID:     true,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs:      completeEntryInputs,
			callFindByIDs:    true,
			measurementItems: []measurementitem.MeasurementItem{pulseRate},
			callUpdate:       true,
			updateErr:        errors.New("update error"),
		},
		{
			name:             "failure measured on before the birth date",
			wantErr:          measurement.ErrInvalidAgeAtMeasurement,
			ctx:              context.Background(),
			userID:           userID,
			callFindByID:     true,
			callFindCustomer: true,
			birthDate:        laterBirthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       beforeBirthMeasuredOn,
			entryInputs:      completeEntryInputs,
		},
		{
			name:             "failure find measurement items error",
			ctx:              context.Background(),
			userID:           userID,
			callFindByID:     true,
			callFindCustomer: true,
			birthDate:        birthDate,
			myTenantIDs:      []string{tenantID.String()},
			measuredOn:       measuredOn,
			entryInputs:      completeEntryInputs,
			callFindByIDs:    true,
			findByIDsErr:     errors.New("find by ids error"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			measurementID := measurement.NewMeasurementID()
			customerID := customer.NewCustomerID()
			createdAt := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
			storedAge, _ := measurement.NewAgeAtMeasurement(70)
			storedMeasurement := measurement.ReconstructMeasurement(measurementID, customerID, storedMeasuredOn, measuredBy, storedAge, measuredBy, false, nil, createdAt, createdAt)

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockMeasurementRepository := mocksmeasurement.NewMockMeasurementRepository(ctrl)
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockMeasurementItemRepository := mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindByID {
				var foundMeasurement measurement.Measurement
				if tt.findByIDErr == nil {
					foundMeasurement = storedMeasurement
				}
				mockMeasurementRepository.EXPECT().FindByID(tt.ctx, measurementID).Return(foundMeasurement, tt.findByIDErr).Times(1)
			}
			if tt.callFindCustomer {
				var foundCustomer customer.Customer
				if tt.findCustomerErr == nil {
					mockCustomer := mockscustomer.NewMockCustomer(ctrl)
					mockCustomer.EXPECT().TenantID().Return(tenantID).AnyTimes()
					mockCustomer.EXPECT().BirthDate().Return(tt.birthDate).AnyTimes()
					foundCustomer = mockCustomer
				}
				mockCustomerRepository.EXPECT().FindByID(tt.ctx, customerID).Return(foundCustomer, tt.findCustomerErr).Times(1)
			}
			if tt.callFindByIDs {
				mockMeasurementItemRepository.EXPECT().FindByIDs(tt.ctx, gomock.Any()).Return(tt.measurementItems, tt.findByIDsErr).Times(1)
			}
			if tt.callUpdate {
				mockMeasurementRepository.EXPECT().Update(tt.ctx, gomock.Any()).Return(tt.updateErr).Times(1)
			}

			u := NewMeasurementUsecase(mockAuthService, mockUserService, mockMeasurementRepository, mockCustomerRepository, mockMeasurementItemRepository)

			entryInputs := tt.entryInputs()

			updatedMeasurement, err := u.UpdateMeasurement(tt.ctx, measurementID, tt.measuredOn, measuredBy, false, entryInputs)
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
				if !updatedMeasurement.MeasuredOn().Equal(tt.measuredOn.Time) {
					t.Errorf("MeasuredOn() = %v, want %v", updatedMeasurement.MeasuredOn(), tt.measuredOn)
				}
				if updatedMeasurement.UpdatedBy().String() != tt.userID {
					t.Errorf("UpdatedBy() = %v, want %v", updatedMeasurement.UpdatedBy(), tt.userID)
				}
				if updatedMeasurement.AgeAtMeasurement().Int() != tt.wantAge {
					t.Errorf("AgeAtMeasurement() = %v, want %v", updatedMeasurement.AgeAtMeasurement().Int(), tt.wantAge)
				}
				if len(updatedMeasurement.Entries()) != len(entryInputs) {
					t.Errorf("len(Entries()) = %v, want %v", len(updatedMeasurement.Entries()), len(entryInputs))
				}
				if !updatedMeasurement.UpdatedAt().After(createdAt) {
					t.Errorf("UpdatedAt() = %v, want after %v", updatedMeasurement.UpdatedAt(), createdAt)
				}
			}
		})
	}
}

func TestDeleteMeasurement(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	otherTenantID := "9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	userID := "google-oauth2|000000000000000000000"

	tests := []struct {
		name             string
		success          bool
		wantErr          error
		ctx              context.Context
		userID           string
		verifyTokenErr   error
		callFindByID     bool
		findByIDErr      error
		callFindCustomer bool
		findCustomerErr  error
		myTenantIDs      []string
		listMyGroupsErr  error
		callDelete       bool
		deleteErr        error
	}{
		{"success delete measurement", true, nil, context.Background(), userID, nil, true, nil, true, nil, []string{tenantID.String()}, nil, true, nil},
		{"failure verify token error", false, nil, context.Background(), "", errors.New("verify token error"), false, nil, false, nil, []string{tenantID.String()}, nil, false, nil},
		{"failure measurement not found", false, measurement.ErrMeasurementNotFound, context.Background(), userID, nil, true, measurement.ErrMeasurementNotFound, false, nil, []string{tenantID.String()}, nil, false, nil},
		{"failure other tenant is hidden as measurement not found", false, measurement.ErrMeasurementNotFound, context.Background(), userID, nil, true, nil, true, nil, []string{otherTenantID}, nil, false, nil},
		{"failure customer not found is hidden as measurement not found", false, measurement.ErrMeasurementNotFound, context.Background(), userID, nil, true, nil, true, customer.ErrCustomerNotFound, []string{tenantID.String()}, nil, false, nil},
		{"failure find customer error", false, nil, context.Background(), userID, nil, true, nil, true, errors.New("find customer error"), []string{tenantID.String()}, nil, false, nil},
		{"failure list my groups error", false, nil, context.Background(), userID, nil, true, nil, true, nil, nil, errors.New("list my groups error"), false, nil},
		{"failure delete error", false, nil, context.Background(), userID, nil, true, nil, true, nil, []string{tenantID.String()}, nil, true, errors.New("delete error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			measurementID := measurement.NewMeasurementID()
			customerID := customer.NewCustomerID()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockMeasurementRepository := mocksmeasurement.NewMockMeasurementRepository(ctrl)
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindByID {
				var foundMeasurement measurement.Measurement
				if tt.findByIDErr == nil {
					mockMeasurement := mocksmeasurement.NewMockMeasurement(ctrl)
					mockMeasurement.EXPECT().CustomerID().Return(customerID).AnyTimes()
					foundMeasurement = mockMeasurement
				}
				mockMeasurementRepository.EXPECT().FindByID(tt.ctx, measurementID).Return(foundMeasurement, tt.findByIDErr).Times(1)
			}
			if tt.callFindCustomer {
				var foundCustomer customer.Customer
				if tt.findCustomerErr == nil {
					mockCustomer := mockscustomer.NewMockCustomer(ctrl)
					mockCustomer.EXPECT().TenantID().Return(tenantID).AnyTimes()
					foundCustomer = mockCustomer
				}
				mockCustomerRepository.EXPECT().FindByID(tt.ctx, customerID).Return(foundCustomer, tt.findCustomerErr).Times(1)
			}
			if tt.callDelete {
				mockMeasurementRepository.EXPECT().Delete(tt.ctx, measurementID).Return(tt.deleteErr).Times(1)
			}

			u := NewMeasurementUsecase(mockAuthService, mockUserService, mockMeasurementRepository, mockCustomerRepository, mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl))

			err := u.DeleteMeasurement(tt.ctx, measurementID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
