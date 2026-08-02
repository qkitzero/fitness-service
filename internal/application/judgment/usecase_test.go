package judgment

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/judgment"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/standard"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
	mocksappauth "github.com/qkitzero/fitness-service/mocks/application/auth"
	mocksappuser "github.com/qkitzero/fitness-service/mocks/application/user"
	mockscustomer "github.com/qkitzero/fitness-service/mocks/domain/customer"
	mocksjudgment "github.com/qkitzero/fitness-service/mocks/domain/judgment"
	mocksmeasurement "github.com/qkitzero/fitness-service/mocks/domain/measurement"
	mocksmeasurementitem "github.com/qkitzero/fitness-service/mocks/domain/measurementitem"
	mocksstandard "github.com/qkitzero/fitness-service/mocks/domain/standard"
)

func TestGetJudgment(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	otherTenantID := "9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	userID := "google-oauth2|000000000000000000000"
	advice, _ := judgment.NewAdvice("週2回のスクワットを継続してください")

	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	motorFunction, _ := measurementitem.NewCategory("motor_function")
	kg, _ := measurementitem.NewUnit("kg")
	twoTrials, _ := measurementitem.NewTrialCount(2)
	higherIsBetter := measurementitem.ScoreDirectionHigherIsBetter
	gripStrengthID := measurementitem.NewMeasurementItemID()
	gripStrengthCode, _ := measurementitem.NewCode("grip_strength")
	gripStrengthName, _ := measurementitem.NewName("握力")
	gripStrength := measurementitem.NewMeasurementItem(gripStrengthID, gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, true, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.SideAggregationMean, []measurementitem.Element{measurementitem.ElementMuscleStrength}, createdAt, updatedAt)

	ageRange6064, _ := standard.NewAgeRange(60, 64)
	ageRange4044, _ := standard.NewAgeRange(40, 44)
	mean6064, _ := standard.NewMean(38)
	mean4044, _ := standard.NewMean(46)
	standardDeviation, _ := standard.NewStandardDeviation(5)
	ageGroupStandards := []standard.AgeGroupStandard{
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), gripStrengthID, standard.GenderMale, ageRange4044, mean4044, standardDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), gripStrengthID, standard.GenderMale, ageRange6064, mean6064, standardDeviation, createdAt, updatedAt),
	}
	zScoreA := standard.ZScore(1.5)
	rankStandardA, _ := standard.NewRankStandard(standard.RankA, &zScoreA, nil, createdAt, updatedAt)
	rankStandards := []standard.RankStandard{rankStandardA}

	entries := func() []measurement.MeasurementEntry {
		trialIndex, _ := measurement.NewTrialIndex(1)
		value, _ := measurement.NewValue(46)
		return []measurement.MeasurementEntry{
			measurement.ReconstructMeasurementEntry(gripStrengthID, false, nil, []measurement.MeasurementValue{
				measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
				measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &value, nil, nil),
			}),
		}
	}

	tests := []struct {
		name                     string
		success                  bool
		wantErr                  error
		ctx                      context.Context
		userID                   string
		verifyTokenErr           error
		callFindMeasurement      bool
		findMeasurementErr       error
		isDraft                  bool
		callFindCustomer         bool
		findCustomerErr          error
		myTenantIDs              []string
		listMyGroupsErr          error
		gender                   customer.Gender
		callListStandards        bool
		evaluated                bool
		findByIDsErr             error
		listAgeGroupStandardsErr error
		listRankStandardsErr     error
		callFindByMeasurementID  bool
		advice                   *judgment.Advice
		findByMeasurementIDErr   error
		wantItemEvaluations      int
		wantAdvice               *judgment.Advice
	}{
		{
			name:                    "success get judgment of an evaluated measurement",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			gender:                  customer.GenderMale,
			callListStandards:       true,
			evaluated:               true,
			callFindByMeasurementID: true,
			advice:                  advice,
			wantItemEvaluations:     1,
			wantAdvice:              advice,
		},
		{
			name:                    "success get judgment of a draft measurement",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			isDraft:                 true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			gender:                  customer.GenderMale,
			callListStandards:       true,
			evaluated:               true,
			callFindByMeasurementID: true,
			wantItemEvaluations:     1,
		},
		{
			name:                    "success get judgment without registered standards",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			gender:                  customer.GenderMale,
			callListStandards:       true,
			callFindByMeasurementID: true,
			advice:                  advice,
			wantAdvice:              advice,
		},
		{
			name:                    "success get judgment without advice",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			gender:                  customer.GenderMale,
			callListStandards:       true,
			callFindByMeasurementID: true,
			findByMeasurementIDErr:  judgment.ErrJudgmentNotFound,
		},
		{
			name:                    "success get judgment of a customer whose gender has no standards",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			gender:                  customer.GenderOther,
			callFindByMeasurementID: true,
			advice:                  advice,
			wantAdvice:              advice,
		},
		{
			name:           "failure verify token error",
			ctx:            context.Background(),
			verifyTokenErr: errors.New("verify token error"),
		},
		{
			name:                "failure measurement not found",
			wantErr:             measurement.ErrMeasurementNotFound,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			findMeasurementErr:  measurement.ErrMeasurementNotFound,
		},
		{
			name:                "failure customer not found",
			wantErr:             customer.ErrCustomerNotFound,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			findCustomerErr:     customer.ErrCustomerNotFound,
		},
		{
			name:                "failure other tenant is hidden as not found",
			wantErr:             judgment.ErrJudgmentNotFound,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{otherTenantID},
		},
		{
			name:                "failure list my groups error",
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			listMyGroupsErr:     errors.New("list my groups error"),
		},
		{
			name:                "failure find measurement items error",
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			gender:              customer.GenderMale,
			callListStandards:   true,
			findByIDsErr:        errors.New("find measurement items error"),
		},
		{
			name:                     "failure list age group standards error",
			ctx:                      context.Background(),
			userID:                   userID,
			callFindMeasurement:      true,
			callFindCustomer:         true,
			myTenantIDs:              []string{tenantID.String()},
			gender:                   customer.GenderMale,
			callListStandards:        true,
			listAgeGroupStandardsErr: errors.New("list age group standards error"),
		},
		{
			name:                 "failure list rank standards error",
			ctx:                  context.Background(),
			userID:               userID,
			callFindMeasurement:  true,
			callFindCustomer:     true,
			myTenantIDs:          []string{tenantID.String()},
			gender:               customer.GenderMale,
			callListStandards:    true,
			listRankStandardsErr: errors.New("list rank standards error"),
		},
		{
			name:                    "failure find judgment error",
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			gender:                  customer.GenderMale,
			callListStandards:       true,
			callFindByMeasurementID: true,
			findByMeasurementIDErr:  errors.New("find judgment error"),
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

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockJudgmentRepository := mocksjudgment.NewMockJudgmentRepository(ctrl)
			mockMeasurementRepository := mocksmeasurement.NewMockMeasurementRepository(ctrl)
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockMeasurementItemRepository := mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl)
			mockAgeGroupStandardRepository := mocksstandard.NewMockAgeGroupStandardRepository(ctrl)
			mockRankStandardRepository := mocksstandard.NewMockRankStandardRepository(ctrl)

			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()

			if tt.callFindMeasurement {
				var foundMeasurement measurement.Measurement
				if tt.findMeasurementErr == nil {
					ageAtMeasurement, _ := measurement.NewAgeAtMeasurement(62)
					mockMeasurement := mocksmeasurement.NewMockMeasurement(ctrl)
					mockMeasurement.EXPECT().ID().Return(measurementID).AnyTimes()
					mockMeasurement.EXPECT().CustomerID().Return(customerID).AnyTimes()
					mockMeasurement.EXPECT().AgeAtMeasurement().Return(ageAtMeasurement).AnyTimes()
					mockMeasurement.EXPECT().IsDraft().Return(tt.isDraft).AnyTimes()
					mockMeasurement.EXPECT().Entries().Return(entries()).AnyTimes()
					foundMeasurement = mockMeasurement
				}
				mockMeasurementRepository.EXPECT().FindByID(tt.ctx, measurementID).Return(foundMeasurement, tt.findMeasurementErr).Times(1)
			}

			if tt.callFindCustomer {
				var foundCustomer customer.Customer
				if tt.findCustomerErr == nil {
					mockCustomer := mockscustomer.NewMockCustomer(ctrl)
					mockCustomer.EXPECT().TenantID().Return(tenantID).AnyTimes()
					mockCustomer.EXPECT().Gender().Return(tt.gender).AnyTimes()
					foundCustomer = mockCustomer
				}
				mockCustomerRepository.EXPECT().FindByID(tt.ctx, customerID).Return(foundCustomer, tt.findCustomerErr).Times(1)
			}

			if tt.callListStandards {
				measurementItems := []measurementitem.MeasurementItem{}
				targetAgeGroupStandards := []standard.AgeGroupStandard{}
				targetRankStandards := []standard.RankStandard{}
				if tt.evaluated {
					measurementItems = []measurementitem.MeasurementItem{gripStrength}
					targetAgeGroupStandards = ageGroupStandards
					targetRankStandards = rankStandards
				}
				measurementItemIDs := []measurementitem.MeasurementItemID{gripStrengthID}
				mockMeasurementItemRepository.EXPECT().FindByIDs(tt.ctx, measurementItemIDs).Return(measurementItems, tt.findByIDsErr).Times(1)
				if tt.findByIDsErr == nil {
					mockAgeGroupStandardRepository.EXPECT().ListByItemIDsAndGender(tt.ctx, measurementItemIDs, standard.GenderMale).Return(targetAgeGroupStandards, tt.listAgeGroupStandardsErr).Times(1)
				}
				if tt.findByIDsErr == nil && tt.listAgeGroupStandardsErr == nil {
					mockRankStandardRepository.EXPECT().List(tt.ctx).Return(targetRankStandards, tt.listRankStandardsErr).Times(1)
				}
			}

			if tt.callFindByMeasurementID {
				var foundJudgment judgment.Judgment
				if tt.findByMeasurementIDErr == nil {
					mockJudgment := mocksjudgment.NewMockJudgment(ctrl)
					mockJudgment.EXPECT().Advice().Return(tt.advice).AnyTimes()
					foundJudgment = mockJudgment
				}
				mockJudgmentRepository.EXPECT().FindByMeasurementID(tt.ctx, measurementID).Return(foundJudgment, tt.findByMeasurementIDErr).Times(1)
			}

			u := NewJudgmentUsecase(mockAuthService, mockUserService, mockJudgmentRepository, mockMeasurementRepository, mockCustomerRepository, mockMeasurementItemRepository, mockAgeGroupStandardRepository, mockRankStandardRepository)

			result, err := u.GetJudgment(tt.ctx, measurementID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if !tt.success {
				return
			}

			if result.MeasurementID != measurementID {
				t.Errorf("MeasurementID = %v, want %v", result.MeasurementID, measurementID)
			}
			if result.IsDraft != tt.isDraft {
				t.Errorf("IsDraft = %v, want %v", result.IsDraft, tt.isDraft)
			}
			if result.Evaluation == nil {
				t.Fatalf("expected an evaluation, but got nil")
			}

			itemEvaluations := result.Evaluation.ItemEvaluations()
			if len(itemEvaluations) != tt.wantItemEvaluations {
				t.Fatalf("len(ItemEvaluations()) = %v, want %v", len(itemEvaluations), tt.wantItemEvaluations)
			}
			if tt.wantItemEvaluations > 0 {
				itemEvaluation := itemEvaluations[0]
				if itemEvaluation.MeasurementItemID() != gripStrengthID {
					t.Errorf("MeasurementItemID() = %v, want %v", itemEvaluation.MeasurementItemID(), gripStrengthID)
				}
				if itemEvaluation.Mean() != mean6064 {
					t.Errorf("Mean() = %v, want %v", itemEvaluation.Mean(), mean6064)
				}
				if itemEvaluation.ZScore() != standard.ZScore(1.6) {
					t.Errorf("ZScore() = %v, want %v", itemEvaluation.ZScore(), standard.ZScore(1.6))
				}
				if itemEvaluation.Rank() != standard.RankA {
					t.Errorf("Rank() = %v, want %v", itemEvaluation.Rank(), standard.RankA)
				}
				motorAge := result.Evaluation.MotorAge()
				if motorAge == nil || motorAge.Int() != 42 {
					t.Errorf("MotorAge() = %v, want %v", motorAge, 42)
				}
			}

			switch {
			case tt.wantAdvice == nil:
				if result.Advice != nil {
					t.Errorf("Advice = %v, want nil", *result.Advice)
				}
			case result.Advice == nil:
				t.Errorf("Advice = nil, want %v", *tt.wantAdvice)
			case *result.Advice != *tt.wantAdvice:
				t.Errorf("Advice = %v, want %v", *result.Advice, *tt.wantAdvice)
			}
		})
	}
}

func TestUpsertJudgmentAdvice(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	otherTenantID := "9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	userID := "google-oauth2|000000000000000000000"
	advice, _ := judgment.NewAdvice("週2回のスクワットを継続してください")
	storedAdvice, _ := judgment.NewAdvice("前回のアドバイス")

	tests := []struct {
		name                    string
		success                 bool
		wantErr                 error
		ctx                     context.Context
		userID                  string
		verifyTokenErr          error
		callFindMeasurement     bool
		findMeasurementErr      error
		callFindCustomer        bool
		myTenantIDs             []string
		callFindByMeasurementID bool
		stored                  *judgment.Advice
		findByMeasurementIDErr  error
		patch                   AdvicePatch
		callUpsert              bool
		upsertErr               error
		wantAdvice              *judgment.Advice
	}{
		{
			name:                    "success update the advice of an existing judgment",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			callFindByMeasurementID: true,
			stored:                  storedAdvice,
			patch:                   AdvicePatch{Advice: advice, HasAdvice: true},
			callUpsert:              true,
			wantAdvice:              advice,
		},
		{
			name:                    "success create a judgment when none exists",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			callFindByMeasurementID: true,
			findByMeasurementIDErr:  judgment.ErrJudgmentNotFound,
			patch:                   AdvicePatch{Advice: advice, HasAdvice: true},
			callUpsert:              true,
			wantAdvice:              advice,
		},
		{
			name:                    "success clear the advice of an existing judgment",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			callFindByMeasurementID: true,
			stored:                  storedAdvice,
			patch:                   AdvicePatch{HasAdvice: true},
			callUpsert:              true,
		},
		{
			name:                    "success keep the stored advice when the advice is omitted",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			callFindByMeasurementID: true,
			stored:                  storedAdvice,
			wantAdvice:              storedAdvice,
		},
		{
			name:                    "success store nothing when the advice is omitted and none exists",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			callFindByMeasurementID: true,
			findByMeasurementIDErr:  judgment.ErrJudgmentNotFound,
		},
		{
			name:                    "success store nothing when clearing an advice that does not exist",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			callFindByMeasurementID: true,
			findByMeasurementIDErr:  judgment.ErrJudgmentNotFound,
			patch:                   AdvicePatch{HasAdvice: true},
		},
		{
			name:           "failure verify token error",
			ctx:            context.Background(),
			verifyTokenErr: errors.New("verify token error"),
			patch:          AdvicePatch{Advice: advice, HasAdvice: true},
		},
		{
			name:                "failure measurement not found",
			wantErr:             measurement.ErrMeasurementNotFound,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			findMeasurementErr:  measurement.ErrMeasurementNotFound,
			patch:               AdvicePatch{Advice: advice, HasAdvice: true},
		},
		{
			name:                "failure other tenant is hidden as not found",
			wantErr:             judgment.ErrJudgmentNotFound,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{otherTenantID},
			patch:               AdvicePatch{Advice: advice, HasAdvice: true},
		},
		{
			name:                    "failure find judgment error",
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			callFindByMeasurementID: true,
			findByMeasurementIDErr:  errors.New("find judgment error"),
			patch:                   AdvicePatch{Advice: advice, HasAdvice: true},
		},
		{
			name:                    "failure upsert error on an existing judgment",
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			callFindByMeasurementID: true,
			stored:                  storedAdvice,
			patch:                   AdvicePatch{Advice: advice, HasAdvice: true},
			callUpsert:              true,
			upsertErr:               errors.New("upsert error"),
		},
		{
			name:                    "failure upsert error on a created judgment",
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			callFindByMeasurementID: true,
			findByMeasurementIDErr:  judgment.ErrJudgmentNotFound,
			patch:                   AdvicePatch{Advice: advice, HasAdvice: true},
			callUpsert:              true,
			upsertErr:               errors.New("upsert error"),
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
			createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
			updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockJudgmentRepository := mocksjudgment.NewMockJudgmentRepository(ctrl)
			mockMeasurementRepository := mocksmeasurement.NewMockMeasurementRepository(ctrl)
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockMeasurementItemRepository := mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl)
			mockAgeGroupStandardRepository := mocksstandard.NewMockAgeGroupStandardRepository(ctrl)
			mockRankStandardRepository := mocksstandard.NewMockRankStandardRepository(ctrl)

			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, nil).AnyTimes()

			if tt.callFindMeasurement {
				var foundMeasurement measurement.Measurement
				if tt.findMeasurementErr == nil {
					mockMeasurement := mocksmeasurement.NewMockMeasurement(ctrl)
					mockMeasurement.EXPECT().CustomerID().Return(customerID).AnyTimes()
					foundMeasurement = mockMeasurement
				}
				mockMeasurementRepository.EXPECT().FindByID(tt.ctx, measurementID).Return(foundMeasurement, tt.findMeasurementErr).Times(1)
			}

			if tt.callFindCustomer {
				mockCustomer := mockscustomer.NewMockCustomer(ctrl)
				mockCustomer.EXPECT().TenantID().Return(tenantID).AnyTimes()
				mockCustomerRepository.EXPECT().FindByID(tt.ctx, customerID).Return(mockCustomer, nil).Times(1)
			}

			if tt.callFindByMeasurementID {
				var foundJudgment judgment.Judgment
				if tt.findByMeasurementIDErr == nil {
					foundJudgment = judgment.ReconstructJudgment(measurementID, tt.stored, createdAt, updatedAt)
				}
				mockJudgmentRepository.EXPECT().FindByMeasurementID(tt.ctx, measurementID).Return(foundJudgment, tt.findByMeasurementIDErr).Times(1)
			}

			if tt.callUpsert {
				mockJudgmentRepository.EXPECT().Upsert(tt.ctx, gomock.Any()).Return(tt.upsertErr).Times(1)
			} else {
				mockJudgmentRepository.EXPECT().Upsert(gomock.Any(), gomock.Any()).Times(0)
			}

			u := NewJudgmentUsecase(mockAuthService, mockUserService, mockJudgmentRepository, mockMeasurementRepository, mockCustomerRepository, mockMeasurementItemRepository, mockAgeGroupStandardRepository, mockRankStandardRepository)

			upsertedJudgment, err := u.UpsertJudgmentAdvice(tt.ctx, measurementID, tt.patch)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if !tt.success {
				return
			}

			if upsertedJudgment.MeasurementID() != measurementID {
				t.Errorf("MeasurementID() = %v, want %v", upsertedJudgment.MeasurementID(), measurementID)
			}
			switch {
			case tt.wantAdvice == nil:
				if upsertedJudgment.Advice() != nil {
					t.Errorf("Advice() = %v, want nil", *upsertedJudgment.Advice())
				}
			case upsertedJudgment.Advice() == nil:
				t.Errorf("Advice() = nil, want %v", *tt.wantAdvice)
			case *upsertedJudgment.Advice() != *tt.wantAdvice:
				t.Errorf("Advice() = %v, want %v", *upsertedJudgment.Advice(), *tt.wantAdvice)
			}
		})
	}
}
