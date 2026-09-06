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
	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/standard"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
	"github.com/qkitzero/fitness-service/internal/domain/training"
	mocksappauth "github.com/qkitzero/fitness-service/mocks/application/auth"
	mocksappuser "github.com/qkitzero/fitness-service/mocks/application/user"
	mockscustomer "github.com/qkitzero/fitness-service/mocks/domain/customer"
	mocksjudgment "github.com/qkitzero/fitness-service/mocks/domain/judgment"
	mocksmeasurement "github.com/qkitzero/fitness-service/mocks/domain/measurement"
	mocksmeasurementitem "github.com/qkitzero/fitness-service/mocks/domain/measurementitem"
	mocksorganization "github.com/qkitzero/fitness-service/mocks/domain/organization"
	mocksstandard "github.com/qkitzero/fitness-service/mocks/domain/standard"
	mockstraining "github.com/qkitzero/fitness-service/mocks/domain/training"
)

func TestItemCodesWithoutAgeGroupStandards(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	motorFunction, _ := measurementitem.NewCategory("motor_function")
	physique, _ := measurementitem.NewCategory("physique")
	kg, _ := measurementitem.NewUnit("kg")
	cm, _ := measurementitem.NewUnit("cm")
	level, _ := measurementitem.NewUnit("level")
	oneTrial, _ := measurementitem.NewTrialCount(1)
	twoTrials, _ := measurementitem.NewTrialCount(2)
	higherIsBetter := measurementitem.ScoreDirectionHigherIsBetter

	gripStrengthID := measurementitem.NewMeasurementItemID()
	gripStrengthCode, _ := measurementitem.NewCode("grip_strength")
	gripStrengthName, _ := measurementitem.NewName("握力")
	gripStrength := measurementitem.NewMeasurementItem(gripStrengthID, gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, measurementitem.SideModeBilateral, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMuscleStrength}, createdAt, updatedAt)

	standUpTestID := measurementitem.NewMeasurementItemID()
	standUpTestCode, _ := measurementitem.NewCode("stand_up_test")
	standUpTestName, _ := measurementitem.NewName("立ち上がり")
	standUpTest := measurementitem.NewMeasurementItem(standUpTestID, standUpTestCode, standUpTestName, motorFunction, level, oneTrial, measurementitem.SideModeOptionalBilateral, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.SideAggregationWorst, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMuscleStrength}, createdAt, updatedAt)

	heightID := measurementitem.NewMeasurementItemID()
	heightCode, _ := measurementitem.NewCode("height")
	heightName, _ := measurementitem.NewName("身長")
	height := measurementitem.NewMeasurementItem(heightID, heightCode, heightName, physique, cm, oneTrial, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, nil, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, nil, createdAt, updatedAt)

	ageRange6064, _ := standard.NewAgeRange(60, 64)
	mean, _ := standard.NewMean(38)
	standardDeviation, _ := standard.NewStandardDeviation(5)
	gripStrengthStandard := standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), gripStrengthID, standard.GenderMale, ageRange6064, mean, standardDeviation, createdAt, updatedAt)

	tests := []struct {
		name              string
		items             []measurementitem.MeasurementItem
		ageGroupStandards []standard.AgeGroupStandard
		want              []string
	}{
		{
			name:              "an item without any age group standard is listed",
			items:             []measurementitem.MeasurementItem{gripStrength, standUpTest},
			ageGroupStandards: []standard.AgeGroupStandard{gripStrengthStandard},
			want:              []string{"stand_up_test"},
		},
		{
			name:              "an item without a score direction is not listed",
			items:             []measurementitem.MeasurementItem{height},
			ageGroupStandards: []standard.AgeGroupStandard{},
			want:              []string{},
		},
		{
			name:              "an item with an age group standard is not listed",
			items:             []measurementitem.MeasurementItem{gripStrength},
			ageGroupStandards: []standard.AgeGroupStandard{gripStrengthStandard},
			want:              []string{},
		},
		{
			name:              "listed codes are sorted",
			items:             []measurementitem.MeasurementItem{standUpTest, gripStrength, height},
			ageGroupStandards: []standard.AgeGroupStandard{},
			want:              []string{"grip_strength", "stand_up_test"},
		},
		{
			name:              "nothing is listed without any item",
			items:             []measurementitem.MeasurementItem{},
			ageGroupStandards: []standard.AgeGroupStandard{gripStrengthStandard},
			want:              []string{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := itemCodesWithoutAgeGroupStandards(tt.items, tt.ageGroupStandards)
			if len(got) != len(tt.want) {
				t.Fatalf("len(itemCodesWithoutAgeGroupStandards()) = %v, want %v", len(got), len(tt.want))
			}
			for i, want := range tt.want {
				if got[i] != want {
					t.Errorf("itemCodesWithoutAgeGroupStandards()[%d] = %v, want %v", i, got[i], want)
				}
			}
		})
	}
}

func TestItemCodesMissingHeight(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	motorFunction, _ := measurementitem.NewCategory("motor_function")
	physique, _ := measurementitem.NewCategory("physique")
	kg, _ := measurementitem.NewUnit("kg")
	cm, _ := measurementitem.NewUnit("cm")
	oneTrial, _ := measurementitem.NewTrialCount(1)
	twoTrials, _ := measurementitem.NewTrialCount(2)
	higherIsBetter := measurementitem.ScoreDirectionHigherIsBetter

	gripStrengthID := measurementitem.NewMeasurementItemID()
	gripStrengthCode, _ := measurementitem.NewCode("grip_strength")
	gripStrengthName, _ := measurementitem.NewName("握力")
	gripStrength := measurementitem.NewMeasurementItem(gripStrengthID, gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, measurementitem.SideModeBilateral, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMuscleStrength}, createdAt, updatedAt)

	twoStepID := measurementitem.NewMeasurementItemID()
	twoStepCode, _ := measurementitem.NewCode("two_step")
	twoStepName, _ := measurementitem.NewName("2ステップ")
	twoStep := measurementitem.NewMeasurementItem(twoStepID, twoStepCode, twoStepName, motorFunction, cm, twoTrials, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.SideAggregationMean, measurementitem.NormalizationHeightRatio, []measurementitem.Element{measurementitem.ElementMobility}, createdAt, updatedAt)

	anotherRatioID := measurementitem.NewMeasurementItemID()
	anotherRatioCode, _ := measurementitem.NewCode("another_ratio_item")
	anotherRatioName, _ := measurementitem.NewName("身長比で評価する別の測定項目")
	anotherRatio := measurementitem.NewMeasurementItem(anotherRatioID, anotherRatioCode, anotherRatioName, motorFunction, cm, oneTrial, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.SideAggregationMean, measurementitem.NormalizationHeightRatio, []measurementitem.Element{measurementitem.ElementMobility}, createdAt, updatedAt)

	heightID := measurementitem.NewMeasurementItemID()
	heightCode, _ := measurementitem.NewCode("height")
	heightName, _ := measurementitem.NewName("身長")
	height := measurementitem.NewMeasurementItem(heightID, heightCode, heightName, physique, cm, oneTrial, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, nil, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, nil, createdAt, updatedAt)

	trialIndex, _ := measurement.NewTrialIndex(1)
	stride, _ := measurement.NewValue(245)
	heightValue, _ := measurement.NewValue(165)
	zeroValue, _ := measurement.NewValue(0)

	tests := []struct {
		name         string
		items        []measurementitem.MeasurementItem
		measuredIDs  []measurementitem.MeasurementItemID
		unmeasurable []measurementitem.MeasurementItemID
		valueless    []measurementitem.MeasurementItemID
		zeroHeight   bool
		want         []string
	}{
		{
			name:        "a normalized item is listed when the height is not measured",
			items:       []measurementitem.MeasurementItem{twoStep, gripStrength},
			measuredIDs: []measurementitem.MeasurementItemID{twoStepID, gripStrengthID},
			want:        []string{"two_step"},
		},
		{
			name:        "a normalized item is not listed when the height is measured",
			items:       []measurementitem.MeasurementItem{twoStep, height},
			measuredIDs: []measurementitem.MeasurementItemID{twoStepID, heightID},
			want:        []string{},
		},
		{
			name:         "a normalized item is listed when the height is unmeasurable",
			items:        []measurementitem.MeasurementItem{twoStep, height},
			measuredIDs:  []measurementitem.MeasurementItemID{twoStepID, heightID},
			unmeasurable: []measurementitem.MeasurementItemID{heightID},
			want:         []string{"two_step"},
		},
		{
			name:        "a normalized item is listed when the height is zero",
			items:       []measurementitem.MeasurementItem{twoStep, height},
			measuredIDs: []measurementitem.MeasurementItemID{twoStepID, heightID},
			zeroHeight:  true,
			want:        []string{"two_step"},
		},
		{
			name:        "a normalized item is listed when the height has no value",
			items:       []measurementitem.MeasurementItem{twoStep, height},
			measuredIDs: []measurementitem.MeasurementItemID{twoStepID, heightID},
			valueless:   []measurementitem.MeasurementItemID{heightID},
			want:        []string{"two_step"},
		},
		{
			name:        "an item without normalization is not listed",
			items:       []measurementitem.MeasurementItem{gripStrength},
			measuredIDs: []measurementitem.MeasurementItemID{gripStrengthID},
			want:        []string{},
		},
		{
			name:        "listed codes are sorted",
			items:       []measurementitem.MeasurementItem{twoStep, anotherRatio},
			measuredIDs: []measurementitem.MeasurementItemID{twoStepID, anotherRatioID},
			want:        []string{"another_ratio_item", "two_step"},
		},
		{
			name:  "nothing is listed without any item",
			items: []measurementitem.MeasurementItem{},
			want:  []string{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			isListed := func(ids []measurementitem.MeasurementItemID, id measurementitem.MeasurementItemID) bool {
				for _, listed := range ids {
					if listed == id {
						return true
					}
				}
				return false
			}

			entries := make([]measurement.MeasurementEntry, 0, len(tt.measuredIDs))
			for _, measuredID := range tt.measuredIDs {
				mockEntry := mocksmeasurement.NewMockMeasurementEntry(ctrl)
				mockEntry.EXPECT().MeasurementItemID().Return(measuredID).AnyTimes()
				mockEntry.EXPECT().Unmeasurable().Return(isListed(tt.unmeasurable, measuredID)).AnyTimes()
				values := []measurement.MeasurementValue{}
				if !isListed(tt.valueless, measuredID) {
					value := stride
					if measuredID == heightID {
						value = heightValue
						if tt.zeroHeight {
							value = zeroValue
						}
					}
					values = append(values, measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil))
				}
				mockEntry.EXPECT().Values().Return(values).AnyTimes()
				entries = append(entries, mockEntry)
			}

			mockMeasurement := mocksmeasurement.NewMockMeasurement(ctrl)
			mockMeasurement.EXPECT().Entries().Return(entries).AnyTimes()

			got := itemCodesMissingHeight(mockMeasurement, tt.items)
			if len(got) != len(tt.want) {
				t.Fatalf("len(itemCodesMissingHeight()) = %v, want %v", len(got), len(tt.want))
			}
			for i, want := range tt.want {
				if got[i] != want {
					t.Errorf("itemCodesMissingHeight()[%d] = %v, want %v", i, got[i], want)
				}
			}
		})
	}
}

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
	gripStrength := measurementitem.NewMeasurementItem(gripStrengthID, gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, measurementitem.SideModeBilateral, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMuscleStrength}, createdAt, updatedAt)

	cm, _ := measurementitem.NewUnit("cm")
	twoStepID := measurementitem.NewMeasurementItemID()
	twoStepCode, _ := measurementitem.NewCode("two_step")
	twoStepName, _ := measurementitem.NewName("2ステップ")
	twoStep := measurementitem.NewMeasurementItem(twoStepID, twoStepCode, twoStepName, motorFunction, cm, twoTrials, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.SideAggregationMean, measurementitem.NormalizationHeightRatio, []measurementitem.Element{measurementitem.ElementMobility}, createdAt, updatedAt)

	level, _ := measurementitem.NewUnit("level")
	oneTrial, _ := measurementitem.NewTrialCount(1)
	standUpTestID := measurementitem.NewMeasurementItemID()
	standUpTestCode, _ := measurementitem.NewCode("stand_up_test")
	standUpTestName, _ := measurementitem.NewName("立ち上がり")
	standUpTest := measurementitem.NewMeasurementItem(standUpTestID, standUpTestCode, standUpTestName, motorFunction, level, oneTrial, measurementitem.SideModeOptionalBilateral, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.SideAggregationWorst, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMuscleStrength}, createdAt, updatedAt)

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

	standUpTestEntry := func() measurement.MeasurementEntry {
		trialIndex, _ := measurement.NewTrialIndex(1)
		value, _ := measurement.NewValue(6)
		return measurement.ReconstructMeasurementEntry(standUpTestID, false, nil, []measurement.MeasurementValue{
			measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil),
		})
	}

	twoStepEntry := func() measurement.MeasurementEntry {
		trialIndex, _ := measurement.NewTrialIndex(1)
		value, _ := measurement.NewValue(245)
		return measurement.ReconstructMeasurementEntry(twoStepID, false, nil, []measurement.MeasurementValue{
			measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil),
		})
	}

	stretchCode, _ := training.NewCode("whole_body_stretch")
	stretchName, _ := training.NewName("全身ストレッチ")
	stretchInstruction, _ := training.NewInstruction("反動をつけずに伸ばす")
	stretchAmount, _ := training.NewAmount(5)
	stretchSets, _ := training.NewSets(1)
	stretchMenu := training.NewTrainingMenu(training.NewTrainingMenuID(), stretchCode, stretchName, measurementitem.ElementFlexibility, training.PartWholeBody, stretchAmount, training.UnitMinutes, stretchSets, stretchInstruction, createdAt, updatedAt)

	walkingCode, _ := training.NewCode("walking")
	walkingName, _ := training.NewName("ウォーキング")
	walkingInstruction, _ := training.NewInstruction("会話ができる速さで歩く")
	walkingAmount, _ := training.NewAmount(20)
	walkingSets, _ := training.NewSets(1)
	walkingMenu := training.NewTrainingMenu(training.NewTrainingMenuID(), walkingCode, walkingName, measurementitem.ElementMuscleEndurance, training.PartWholeBody, walkingAmount, training.UnitMinutes, walkingSets, walkingInstruction, createdAt, updatedAt)

	trainingMenus := []training.TrainingMenu{stretchMenu, walkingMenu}
	firstSortOrder, _ := training.NewSortOrder(1)
	decade60, _ := training.NewDecade(60)
	fixedMenus := []training.FixedMenu{training.NewFixedMenu(firstSortOrder, stretchMenu.ID(), createdAt, updatedAt)}
	ageDecadeMenus := []training.AgeDecadeMenu{training.NewAgeDecadeMenu(decade60, firstSortOrder, walkingMenu.ID(), createdAt, updatedAt)}

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
		withUnregisteredItem     bool
		withNormalizedItem       bool
		findByIDsErr             error
		listAgeGroupStandardsErr error
		listRankStandardsErr     error
		callPrescribe            bool
		overridden               bool
		listOverridesErr         error
		listTrainingMenusErr     error
		listElementMenusErr      error
		listFixedMenusErr        error
		listAgeDecadeMenusErr    error
		age                      int
		callFindByMeasurementID  bool
		advice                   *judgment.Advice
		findByMeasurementIDErr   error
		wantItemEvaluations      int
		wantPrescribedMenus      []judgment.PrescriptionSource
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
			callPrescribe:           true,
			callFindByMeasurementID: true,
			advice:                  advice,
			wantItemEvaluations:     1,
			wantAdvice:              advice,
			wantPrescribedMenus:     []judgment.PrescriptionSource{judgment.PrescriptionSourceFixed, judgment.PrescriptionSourceAgeDecade},
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
			callPrescribe:           true,
			callFindByMeasurementID: true,
			wantItemEvaluations:     1,
			wantPrescribedMenus:     []judgment.PrescriptionSource{judgment.PrescriptionSourceFixed, judgment.PrescriptionSourceAgeDecade},
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
			callPrescribe:           true,
			callFindByMeasurementID: true,
			advice:                  advice,
			wantAdvice:              advice,
			wantPrescribedMenus:     []judgment.PrescriptionSource{judgment.PrescriptionSourceFixed, judgment.PrescriptionSourceAgeDecade},
		},
		{
			name:                    "success get judgment of a measurement including an item without any age group standard",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			gender:                  customer.GenderMale,
			callListStandards:       true,
			evaluated:               true,
			withUnregisteredItem:    true,
			callPrescribe:           true,
			callFindByMeasurementID: true,
			advice:                  advice,
			wantAdvice:              advice,
			wantItemEvaluations:     1,
			wantPrescribedMenus:     []judgment.PrescriptionSource{judgment.PrescriptionSourceFixed, judgment.PrescriptionSourceAgeDecade},
		},
		{
			name:                    "success get judgment of a measurement missing the height of a normalized item",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			gender:                  customer.GenderMale,
			callListStandards:       true,
			evaluated:               true,
			withNormalizedItem:      true,
			callPrescribe:           true,
			callFindByMeasurementID: true,
			advice:                  advice,
			wantAdvice:              advice,
			wantItemEvaluations:     1,
			wantPrescribedMenus:     []judgment.PrescriptionSource{judgment.PrescriptionSourceFixed, judgment.PrescriptionSourceAgeDecade},
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
			callPrescribe:           true,
			callFindByMeasurementID: true,
			findByMeasurementIDErr:  judgment.ErrJudgmentNotFound,
			wantPrescribedMenus:     []judgment.PrescriptionSource{judgment.PrescriptionSourceFixed, judgment.PrescriptionSourceAgeDecade},
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
			callPrescribe:           true,
			callFindByMeasurementID: true,
			advice:                  advice,
			wantAdvice:              advice,
			wantPrescribedMenus:     []judgment.PrescriptionSource{judgment.PrescriptionSourceFixed, judgment.PrescriptionSourceAgeDecade},
		},
		{
			name:                    "success get judgment of a measurement with an edited prescription",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			gender:                  customer.GenderMale,
			callListStandards:       true,
			evaluated:               true,
			callPrescribe:           true,
			overridden:              true,
			callFindByMeasurementID: true,
			wantItemEvaluations:     1,
			wantPrescribedMenus:     []judgment.PrescriptionSource{judgment.PrescriptionSourceManual},
		},
		{
			name:                    "success get judgment of a measurement above every age group",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			gender:                  customer.GenderMale,
			callListStandards:       true,
			evaluated:               true,
			callPrescribe:           true,
			age:                     70,
			callFindByMeasurementID: true,
			wantItemEvaluations:     1,
			wantPrescribedMenus:     []judgment.PrescriptionSource{judgment.PrescriptionSourceFixed},
		},
		{
			name:                    "success get judgment of a centenarian without age decade menus",
			success:                 true,
			ctx:                     context.Background(),
			userID:                  userID,
			callFindMeasurement:     true,
			callFindCustomer:        true,
			myTenantIDs:             []string{tenantID.String()},
			gender:                  customer.GenderMale,
			callListStandards:       true,
			callPrescribe:           true,
			age:                     105,
			callFindByMeasurementID: true,
			wantPrescribedMenus:     []judgment.PrescriptionSource{judgment.PrescriptionSourceFixed},
		},
		{
			name:                 "failure find training menus of an edited prescription error",
			ctx:                  context.Background(),
			userID:               userID,
			callFindMeasurement:  true,
			callFindCustomer:     true,
			myTenantIDs:          []string{tenantID.String()},
			gender:               customer.GenderMale,
			callListStandards:    true,
			callPrescribe:        true,
			overridden:           true,
			listTrainingMenusErr: errors.New("find training menus error"),
		},
		{
			name:                "failure list prescribed menu overrides error",
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			gender:              customer.GenderMale,
			callListStandards:   true,
			callPrescribe:       true,
			listOverridesErr:    errors.New("list prescribed menu overrides error"),
		},
		{
			name:                 "failure list training menus error",
			ctx:                  context.Background(),
			userID:               userID,
			callFindMeasurement:  true,
			callFindCustomer:     true,
			myTenantIDs:          []string{tenantID.String()},
			gender:               customer.GenderMale,
			callListStandards:    true,
			callPrescribe:        true,
			listTrainingMenusErr: errors.New("list training menus error"),
		},
		{
			name:                "failure list element menus error",
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			gender:              customer.GenderMale,
			callListStandards:   true,
			callPrescribe:       true,
			listElementMenusErr: errors.New("list element menus error"),
		},
		{
			name:                "failure list fixed menus error",
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			gender:              customer.GenderMale,
			callListStandards:   true,
			callPrescribe:       true,
			listFixedMenusErr:   errors.New("list fixed menus error"),
		},
		{
			name:                  "failure list age decade menus error",
			ctx:                   context.Background(),
			userID:                userID,
			callFindMeasurement:   true,
			callFindCustomer:      true,
			myTenantIDs:           []string{tenantID.String()},
			gender:                customer.GenderMale,
			callListStandards:     true,
			callPrescribe:         true,
			listAgeDecadeMenusErr: errors.New("list age decade menus error"),
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
			callPrescribe:           true,
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
			mockOrganizationRepository := mocksorganization.NewMockOrganizationRepository(ctrl)
			mockMeasurementItemRepository := mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl)
			mockAgeGroupStandardRepository := mocksstandard.NewMockAgeGroupStandardRepository(ctrl)
			mockRankStandardRepository := mocksstandard.NewMockRankStandardRepository(ctrl)
			mockPrescribedMenuOverrideRepository := mocksjudgment.NewMockPrescribedMenuOverrideRepository(ctrl)
			mockTrainingMenuRepository := mockstraining.NewMockTrainingMenuRepository(ctrl)
			mockPrescriptionRuleRepository := mockstraining.NewMockPrescriptionRuleRepository(ctrl)

			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()

			if tt.callFindMeasurement {
				var foundMeasurement measurement.Measurement
				if tt.findMeasurementErr == nil {
					age := tt.age
					if age == 0 {
						age = 62
					}
					ageAtMeasurement, _ := measurement.NewAgeAtMeasurement(age)
					mockMeasurement := mocksmeasurement.NewMockMeasurement(ctrl)
					mockMeasurement.EXPECT().ID().Return(measurementID).AnyTimes()
					mockMeasurement.EXPECT().CustomerID().Return(customerID).AnyTimes()
					mockMeasurement.EXPECT().AgeAtMeasurement().Return(ageAtMeasurement).AnyTimes()
					mockMeasurement.EXPECT().IsDraft().Return(tt.isDraft).AnyTimes()
					measurementEntries := entries()
					if tt.withUnregisteredItem {
						measurementEntries = append(measurementEntries, standUpTestEntry())
					}
					if tt.withNormalizedItem {
						measurementEntries = append(measurementEntries, twoStepEntry())
					}
					mockMeasurement.EXPECT().Entries().Return(measurementEntries).AnyTimes()
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
				measurementItemIDs := []measurementitem.MeasurementItemID{gripStrengthID}
				if tt.withUnregisteredItem {
					measurementItemIDs = append(measurementItemIDs, standUpTestID)
				}
				if tt.withNormalizedItem {
					measurementItemIDs = append(measurementItemIDs, twoStepID)
				}
				if tt.evaluated {
					measurementItems = []measurementitem.MeasurementItem{gripStrength}
					if tt.withUnregisteredItem {
						measurementItems = append(measurementItems, standUpTest)
					}
					if tt.withNormalizedItem {
						measurementItems = append(measurementItems, twoStep)
					}
					targetAgeGroupStandards = ageGroupStandards
					targetRankStandards = rankStandards
				}
				mockMeasurementItemRepository.EXPECT().FindByIDs(tt.ctx, measurementItemIDs).Return(measurementItems, tt.findByIDsErr).Times(1)
				if tt.findByIDsErr == nil {
					mockAgeGroupStandardRepository.EXPECT().ListByItemIDsAndGender(tt.ctx, measurementItemIDs, standard.GenderMale).Return(targetAgeGroupStandards, tt.listAgeGroupStandardsErr).Times(1)
				}
				if tt.findByIDsErr == nil && tt.listAgeGroupStandardsErr == nil {
					mockRankStandardRepository.EXPECT().List(tt.ctx).Return(targetRankStandards, tt.listRankStandardsErr).Times(1)
				}
			}

			if tt.callPrescribe {
				overrides := []judgment.PrescribedMenuOverride{}
				if tt.overridden {
					sortOrder, _ := training.NewSortOrder(1)
					overrides = append(overrides, judgment.NewPrescribedMenuOverride(judgment.NewPrescribedMenuOverrideID(), measurementID, sortOrder, nil, nil, walkingMenu.ID(), walkingAmount, training.UnitMinutes, walkingSets, createdAt, updatedAt))
				}
				mockPrescribedMenuOverrideRepository.EXPECT().ListByMeasurementID(tt.ctx, measurementID).Return(overrides, tt.listOverridesErr).Times(1)
				if tt.listOverridesErr == nil && tt.overridden {
					mockTrainingMenuRepository.EXPECT().FindByIDs(tt.ctx, []training.TrainingMenuID{walkingMenu.ID()}).Return(trainingMenus, tt.listTrainingMenusErr).Times(1)
				}
				if tt.listOverridesErr == nil && !tt.overridden {
					mockTrainingMenuRepository.EXPECT().List(tt.ctx).Return(trainingMenus, tt.listTrainingMenusErr).Times(1)
				}
				if tt.listOverridesErr == nil && tt.listTrainingMenusErr == nil && !tt.overridden {
					mockPrescriptionRuleRepository.EXPECT().ListElementMenus(tt.ctx).Return(nil, tt.listElementMenusErr).Times(1)
					if tt.listElementMenusErr == nil {
						mockPrescriptionRuleRepository.EXPECT().ListFixedMenus(tt.ctx).Return(fixedMenus, tt.listFixedMenusErr).Times(1)
					}
					if tt.listElementMenusErr == nil && tt.listFixedMenusErr == nil {
						mockPrescriptionRuleRepository.EXPECT().ListAgeDecadeMenus(tt.ctx).Return(ageDecadeMenus, tt.listAgeDecadeMenusErr).Times(1)
					}
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

			u := NewJudgmentUsecase(mockAuthService, mockUserService, mockJudgmentRepository, mockPrescribedMenuOverrideRepository, mockMeasurementRepository, mockCustomerRepository, mockOrganizationRepository, mockMeasurementItemRepository, mockAgeGroupStandardRepository, mockRankStandardRepository, mockTrainingMenuRepository, mockPrescriptionRuleRepository)

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

			if result.Prescription == nil {
				t.Fatalf("expected a prescription, but got nil")
			}

			prescribedMenus := result.Prescription.PrescribedMenus()
			if len(prescribedMenus) != len(tt.wantPrescribedMenus) {
				t.Fatalf("len(PrescribedMenus()) = %v, want %v", len(prescribedMenus), len(tt.wantPrescribedMenus))
			}
			for i, wantSource := range tt.wantPrescribedMenus {
				if prescribedMenus[i].Source() != wantSource {
					t.Errorf("PrescribedMenus()[%d].Source() = %v, want %v", i, prescribedMenus[i].Source(), wantSource)
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

func TestListOrganizationJudgments(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	otherTenantID := "9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	userID := "google-oauth2|000000000000000000000"

	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	motorFunction, _ := measurementitem.NewCategory("motor_function")
	kg, _ := measurementitem.NewUnit("kg")
	twoTrials, _ := measurementitem.NewTrialCount(2)
	higherIsBetter := measurementitem.ScoreDirectionHigherIsBetter
	gripStrengthID := measurementitem.NewMeasurementItemID()
	gripStrengthCode, _ := measurementitem.NewCode("grip_strength")
	gripStrengthName, _ := measurementitem.NewName("握力")
	gripStrength := measurementitem.NewMeasurementItem(gripStrengthID, gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, measurementitem.SideModeBilateral, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMuscleStrength}, createdAt, updatedAt)

	level, _ := measurementitem.NewUnit("level")
	oneTrial, _ := measurementitem.NewTrialCount(1)
	standUpTestID := measurementitem.NewMeasurementItemID()
	standUpTestCode, _ := measurementitem.NewCode("stand_up_test")
	standUpTestName, _ := measurementitem.NewName("立ち上がり")
	standUpTest := measurementitem.NewMeasurementItem(standUpTestID, standUpTestCode, standUpTestName, motorFunction, level, oneTrial, measurementitem.SideModeOptionalBilateral, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.SideAggregationWorst, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMuscleStrength}, createdAt, updatedAt)

	cm, _ := measurementitem.NewUnit("cm")
	twoStepID := measurementitem.NewMeasurementItemID()
	twoStepCode, _ := measurementitem.NewCode("two_step")
	twoStepName, _ := measurementitem.NewName("2ステップ")
	twoStep := measurementitem.NewMeasurementItem(twoStepID, twoStepCode, twoStepName, motorFunction, cm, twoTrials, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.SideAggregationMean, measurementitem.NormalizationHeightRatio, []measurementitem.Element{measurementitem.ElementMobility}, createdAt, updatedAt)

	ageRange4044, _ := standard.NewAgeRange(40, 44)
	ageRange5054, _ := standard.NewAgeRange(50, 54)
	ageRange6064, _ := standard.NewAgeRange(60, 64)
	maleMean4044, _ := standard.NewMean(46)
	maleMean6064, _ := standard.NewMean(38)
	femaleMean5054, _ := standard.NewMean(45.8)
	femaleMean6064, _ := standard.NewMean(44)
	standardDeviation, _ := standard.NewStandardDeviation(5)
	ageGroupStandards := []standard.AgeGroupStandard{
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), gripStrengthID, standard.GenderMale, ageRange4044, maleMean4044, standardDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), gripStrengthID, standard.GenderMale, ageRange6064, maleMean6064, standardDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), gripStrengthID, standard.GenderFemale, ageRange5054, femaleMean5054, standardDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), gripStrengthID, standard.GenderFemale, ageRange6064, femaleMean6064, standardDeviation, createdAt, updatedAt),
	}
	zScoreA := standard.ZScore(1.5)
	zScoreBMin := standard.ZScore(0.5)
	zScoreBMax := standard.ZScore(1.5)
	zScoreCMin := standard.ZScore(-0.5)
	zScoreCMax := standard.ZScore(0.5)
	rankStandardA, _ := standard.NewRankStandard(standard.RankA, &zScoreA, nil, createdAt, updatedAt)
	rankStandardB, _ := standard.NewRankStandard(standard.RankB, &zScoreBMin, &zScoreBMax, createdAt, updatedAt)
	rankStandardC, _ := standard.NewRankStandard(standard.RankC, &zScoreCMin, &zScoreCMax, createdAt, updatedAt)
	rankStandards := []standard.RankStandard{rankStandardA, rankStandardB, rankStandardC}

	measuredOn, _ := measurement.NewMeasuredOn(2026, 8, 1)
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

	standUpTestEntry := func() measurement.MeasurementEntry {
		trialIndex, _ := measurement.NewTrialIndex(1)
		value, _ := measurement.NewValue(6)
		return measurement.ReconstructMeasurementEntry(standUpTestID, false, nil, []measurement.MeasurementValue{
			measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil),
		})
	}

	twoStepEntry := func() measurement.MeasurementEntry {
		trialIndex, _ := measurement.NewTrialIndex(1)
		value, _ := measurement.NewValue(245)
		return measurement.ReconstructMeasurementEntry(twoStepID, false, nil, []measurement.MeasurementValue{
			measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil),
		})
	}

	type measurementSpec struct {
		customerIndex       int
		isDraft             bool
		age                 int
		wantItemEvaluations int
		wantMean            float64
		wantZScore          float64
		wantRank            standard.Rank
		wantMotorAge        int
	}

	tests := []struct {
		name                     string
		success                  bool
		wantErr                  error
		ctx                      context.Context
		verifyTokenErr           error
		callFindOrganization     bool
		findOrganizationErr      error
		myTenantIDs              []string
		listMyGroupsErr          error
		callListCustomers        bool
		includeInactive          bool
		genders                  []customer.Gender
		listCustomersErr         error
		callListMeasurements     bool
		measurements             []measurementSpec
		listMeasurementsErr      error
		callListMasters          bool
		registeredStandards      bool
		unregisteredItem         bool
		normalizedItem           bool
		listMeasurementItemsErr  error
		listAgeGroupStandardsErr error
		listRankStandardsErr     error
	}{
		{
			name:                 "success list organization judgments",
			success:              true,
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{tenantID.String()},
			callListCustomers:    true,
			genders:              []customer.Gender{customer.GenderMale, customer.GenderFemale},
			callListMeasurements: true,
			measurements: []measurementSpec{
				{customerIndex: 0, isDraft: false, wantItemEvaluations: 1, wantMean: 38, wantZScore: 1.6, wantRank: standard.RankA, wantMotorAge: 42},
				{customerIndex: 0, isDraft: true, wantItemEvaluations: 1, wantMean: 38, wantZScore: 1.6, wantRank: standard.RankA, wantMotorAge: 42},
				{customerIndex: 1, isDraft: false, wantItemEvaluations: 1, wantMean: 44, wantZScore: 0.4, wantRank: standard.RankC, wantMotorAge: 52},
			},
			callListMasters:     true,
			registeredStandards: true,
		},
		{
			name:                 "success list organization judgments including inactive customers",
			success:              true,
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{tenantID.String()},
			callListCustomers:    true,
			includeInactive:      true,
			genders:              []customer.Gender{customer.GenderMale},
			callListMeasurements: true,
			measurements: []measurementSpec{
				{customerIndex: 0, isDraft: false, wantItemEvaluations: 1, wantMean: 38, wantZScore: 1.6, wantRank: standard.RankA, wantMotorAge: 42},
			},
			callListMasters:     true,
			registeredStandards: true,
		},
		{
			name:                 "success list judgments of a customer above every age group",
			success:              true,
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{tenantID.String()},
			callListCustomers:    true,
			genders:              []customer.Gender{customer.GenderMale},
			callListMeasurements: true,
			measurements: []measurementSpec{
				{customerIndex: 0, isDraft: false, age: 70, wantItemEvaluations: 1, wantMean: 38, wantZScore: 1.6, wantRank: standard.RankA, wantMotorAge: 42},
			},
			callListMasters:     true,
			registeredStandards: true,
		},
		{
			name:                 "success list judgments of a customer whose gender has no standards",
			success:              true,
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{tenantID.String()},
			callListCustomers:    true,
			genders:              []customer.Gender{customer.GenderOther},
			callListMeasurements: true,
			measurements:         []measurementSpec{{customerIndex: 0, isDraft: false}},
			callListMasters:      true,
			registeredStandards:  true,
		},
		{
			name:                 "success list judgments without registered standards",
			success:              true,
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{tenantID.String()},
			callListCustomers:    true,
			genders:              []customer.Gender{customer.GenderMale},
			callListMeasurements: true,
			measurements:         []measurementSpec{{customerIndex: 0, isDraft: false}},
			callListMasters:      true,
		},
		{
			name:                 "success list judgments including an item without any age group standard",
			success:              true,
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{tenantID.String()},
			callListCustomers:    true,
			genders:              []customer.Gender{customer.GenderMale},
			callListMeasurements: true,
			measurements: []measurementSpec{
				{customerIndex: 0, isDraft: false, wantItemEvaluations: 1, wantMean: 38, wantZScore: 1.6, wantRank: standard.RankA, wantMotorAge: 42},
			},
			callListMasters:     true,
			registeredStandards: true,
			unregisteredItem:    true,
		},
		{
			name:                 "success list judgments missing the height of a normalized item",
			success:              true,
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{tenantID.String()},
			callListCustomers:    true,
			genders:              []customer.Gender{customer.GenderMale},
			callListMeasurements: true,
			measurements: []measurementSpec{
				{customerIndex: 0, isDraft: false, wantItemEvaluations: 1, wantMean: 38, wantZScore: 1.6, wantRank: standard.RankA, wantMotorAge: 42},
			},
			callListMasters:     true,
			registeredStandards: true,
			normalizedItem:      true,
		},
		{
			name:                 "success list no judgments of an organization without customers",
			success:              true,
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{tenantID.String()},
			callListCustomers:    true,
			callListMeasurements: true,
		},
		{
			name:                 "success list no judgments of customers without measurements",
			success:              true,
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{tenantID.String()},
			callListCustomers:    true,
			genders:              []customer.Gender{customer.GenderMale},
			callListMeasurements: true,
		},
		{
			name:           "failure verify token error",
			ctx:            context.Background(),
			verifyTokenErr: errors.New("verify token error"),
		},
		{
			name:                 "failure organization not found",
			wantErr:              organization.ErrOrganizationNotFound,
			ctx:                  context.Background(),
			callFindOrganization: true,
			findOrganizationErr:  organization.ErrOrganizationNotFound,
		},
		{
			name:                 "failure other tenant is hidden as not found",
			wantErr:              organization.ErrOrganizationNotFound,
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{otherTenantID},
		},
		{
			name:                 "failure list my groups error",
			ctx:                  context.Background(),
			callFindOrganization: true,
			listMyGroupsErr:      errors.New("list my groups error"),
		},
		{
			name:                 "failure list customers error",
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{tenantID.String()},
			callListCustomers:    true,
			listCustomersErr:     errors.New("list customers error"),
		},
		{
			name:                 "failure list measurements error",
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{tenantID.String()},
			callListCustomers:    true,
			genders:              []customer.Gender{customer.GenderMale},
			callListMeasurements: true,
			listMeasurementsErr:  errors.New("list measurements error"),
		},
		{
			name:                    "failure list measurement items error",
			ctx:                     context.Background(),
			callFindOrganization:    true,
			myTenantIDs:             []string{tenantID.String()},
			callListCustomers:       true,
			genders:                 []customer.Gender{customer.GenderMale},
			callListMeasurements:    true,
			measurements:            []measurementSpec{{customerIndex: 0, isDraft: false}},
			callListMasters:         true,
			listMeasurementItemsErr: errors.New("list measurement items error"),
		},
		{
			name:                     "failure list age group standards error",
			ctx:                      context.Background(),
			callFindOrganization:     true,
			myTenantIDs:              []string{tenantID.String()},
			callListCustomers:        true,
			genders:                  []customer.Gender{customer.GenderMale},
			callListMeasurements:     true,
			measurements:             []measurementSpec{{customerIndex: 0, isDraft: false}},
			callListMasters:          true,
			listAgeGroupStandardsErr: errors.New("list age group standards error"),
		},
		{
			name:                 "failure list rank standards error",
			ctx:                  context.Background(),
			callFindOrganization: true,
			myTenantIDs:          []string{tenantID.String()},
			callListCustomers:    true,
			genders:              []customer.Gender{customer.GenderMale},
			callListMeasurements: true,
			measurements:         []measurementSpec{{customerIndex: 0, isDraft: false}},
			callListMasters:      true,
			listRankStandardsErr: errors.New("list rank standards error"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			organizationID := organization.NewOrganizationID()
			customerIDs := make([]customer.CustomerID, 0, len(tt.genders))
			for range tt.genders {
				customerIDs = append(customerIDs, customer.NewCustomerID())
			}

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockJudgmentRepository := mocksjudgment.NewMockJudgmentRepository(ctrl)
			mockMeasurementRepository := mocksmeasurement.NewMockMeasurementRepository(ctrl)
			mockCustomerRepository := mockscustomer.NewMockCustomerRepository(ctrl)
			mockOrganizationRepository := mocksorganization.NewMockOrganizationRepository(ctrl)
			mockMeasurementItemRepository := mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl)
			mockAgeGroupStandardRepository := mocksstandard.NewMockAgeGroupStandardRepository(ctrl)
			mockRankStandardRepository := mocksstandard.NewMockRankStandardRepository(ctrl)
			mockPrescribedMenuOverrideRepository := mocksjudgment.NewMockPrescribedMenuOverrideRepository(ctrl)
			mockTrainingMenuRepository := mockstraining.NewMockTrainingMenuRepository(ctrl)
			mockPrescriptionRuleRepository := mockstraining.NewMockPrescriptionRuleRepository(ctrl)

			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(userID, tt.verifyTokenErr).AnyTimes()

			if tt.callFindOrganization {
				var foundOrganization organization.Organization
				if tt.findOrganizationErr == nil {
					mockOrganization := mocksorganization.NewMockOrganization(ctrl)
					mockOrganization.EXPECT().TenantID().Return(tenantID).AnyTimes()
					foundOrganization = mockOrganization
				}
				mockOrganizationRepository.EXPECT().FindByID(tt.ctx, organizationID).Return(foundOrganization, tt.findOrganizationErr).Times(1)
				if tt.findOrganizationErr == nil {
					mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).Times(1)
				}
			}

			if tt.callListCustomers {
				var customers []customer.Customer
				if tt.listCustomersErr == nil {
					customers = make([]customer.Customer, 0, len(tt.genders))
					for i, gender := range tt.genders {
						mockCustomer := mockscustomer.NewMockCustomer(ctrl)
						mockCustomer.EXPECT().ID().Return(customerIDs[i]).AnyTimes()
						mockCustomer.EXPECT().Gender().Return(gender).AnyTimes()
						customers = append(customers, mockCustomer)
					}
				}
				mockCustomerRepository.EXPECT().ListByOrganizationID(tt.ctx, organizationID, tt.includeInactive).Return(customers, tt.listCustomersErr).Times(1)
			}

			measurementIDs := make([]measurement.MeasurementID, 0, len(tt.measurements))
			ages := make([]int, 0, len(tt.measurements))
			if tt.callListMeasurements {
				var measurements []measurement.Measurement
				if tt.listMeasurementsErr == nil {
					measurements = make([]measurement.Measurement, 0, len(tt.measurements))
					for _, spec := range tt.measurements {
						measurementID := measurement.NewMeasurementID()
						measurementIDs = append(measurementIDs, measurementID)
						age := spec.age
						if age == 0 {
							age = 62
						}
						ages = append(ages, age)
						ageAtMeasurement, _ := measurement.NewAgeAtMeasurement(age)
						mockMeasurement := mocksmeasurement.NewMockMeasurement(ctrl)
						mockMeasurement.EXPECT().ID().Return(measurementID).AnyTimes()
						mockMeasurement.EXPECT().CustomerID().Return(customerIDs[spec.customerIndex]).AnyTimes()
						mockMeasurement.EXPECT().MeasuredOn().Return(measuredOn).AnyTimes()
						mockMeasurement.EXPECT().AgeAtMeasurement().Return(ageAtMeasurement).AnyTimes()
						mockMeasurement.EXPECT().IsDraft().Return(spec.isDraft).AnyTimes()
						measurementEntries := entries()
						if tt.unregisteredItem {
							measurementEntries = append(measurementEntries, standUpTestEntry())
						}
						if tt.normalizedItem {
							measurementEntries = append(measurementEntries, twoStepEntry())
						}
						mockMeasurement.EXPECT().Entries().Return(measurementEntries).AnyTimes()
						measurements = append(measurements, mockMeasurement)
					}
				}
				mockMeasurementRepository.EXPECT().ListByCustomerIDs(tt.ctx, customerIDs).Return(measurements, tt.listMeasurementsErr).Times(1)
			}

			if tt.callListMasters {
				measurementItems := []measurementitem.MeasurementItem{}
				targetAgeGroupStandards := []standard.AgeGroupStandard{}
				targetRankStandards := []standard.RankStandard{}
				if tt.registeredStandards {
					measurementItems = []measurementitem.MeasurementItem{gripStrength}
					if tt.unregisteredItem {
						measurementItems = append(measurementItems, standUpTest)
					}
					if tt.normalizedItem {
						measurementItems = append(measurementItems, twoStep)
					}
					targetAgeGroupStandards = ageGroupStandards
					targetRankStandards = rankStandards
				}
				mockMeasurementItemRepository.EXPECT().List(tt.ctx).Return(measurementItems, tt.listMeasurementItemsErr).Times(1)
				if tt.listMeasurementItemsErr == nil {
					mockAgeGroupStandardRepository.EXPECT().List(tt.ctx).Return(targetAgeGroupStandards, tt.listAgeGroupStandardsErr).Times(1)
				}
				if tt.listMeasurementItemsErr == nil && tt.listAgeGroupStandardsErr == nil {
					mockRankStandardRepository.EXPECT().List(tt.ctx).Return(targetRankStandards, tt.listRankStandardsErr).Times(1)
				}
			}

			u := NewJudgmentUsecase(mockAuthService, mockUserService, mockJudgmentRepository, mockPrescribedMenuOverrideRepository, mockMeasurementRepository, mockCustomerRepository, mockOrganizationRepository, mockMeasurementItemRepository, mockAgeGroupStandardRepository, mockRankStandardRepository, mockTrainingMenuRepository, mockPrescriptionRuleRepository)

			results, err := u.ListOrganizationJudgments(tt.ctx, organizationID, tt.includeInactive)
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

			if len(results) != len(tt.measurements) {
				t.Fatalf("len(results) = %v, want %v", len(results), len(tt.measurements))
			}
			for i, spec := range tt.measurements {
				result := results[i]
				if result.MeasurementID != measurementIDs[i] {
					t.Errorf("results[%d].MeasurementID = %v, want %v", i, result.MeasurementID, measurementIDs[i])
				}
				if result.CustomerID != customerIDs[spec.customerIndex] {
					t.Errorf("results[%d].CustomerID = %v, want %v", i, result.CustomerID, customerIDs[spec.customerIndex])
				}
				if !result.MeasuredOn.Equal(measuredOn.Time) {
					t.Errorf("results[%d].MeasuredOn = %v, want %v", i, result.MeasuredOn, measuredOn)
				}
				if result.AgeAtMeasurement.Int() != ages[i] {
					t.Errorf("results[%d].AgeAtMeasurement = %v, want %v", i, result.AgeAtMeasurement.Int(), ages[i])
				}
				if result.IsDraft != spec.isDraft {
					t.Errorf("results[%d].IsDraft = %v, want %v", i, result.IsDraft, spec.isDraft)
				}

				itemEvaluations := result.Evaluation.ItemEvaluations()
				elementEvaluations := result.Evaluation.ElementEvaluations()
				motorAge := result.Evaluation.MotorAge()
				if len(itemEvaluations) != spec.wantItemEvaluations {
					t.Fatalf("len(results[%d].ItemEvaluations()) = %v, want %v", i, len(itemEvaluations), spec.wantItemEvaluations)
				}
				if spec.wantItemEvaluations == 0 {
					if len(elementEvaluations) != 0 {
						t.Errorf("len(results[%d].ElementEvaluations()) = %v, want %v", i, len(elementEvaluations), 0)
					}
					if motorAge != nil {
						t.Errorf("results[%d].MotorAge() = %v, want nil", i, motorAge)
					}
					continue
				}
				if itemEvaluations[0].Mean().Float64() != spec.wantMean {
					t.Errorf("results[%d].ItemEvaluations()[0].Mean() = %v, want %v", i, itemEvaluations[0].Mean().Float64(), spec.wantMean)
				}
				if itemEvaluations[0].ZScore().Float64() != spec.wantZScore {
					t.Errorf("results[%d].ItemEvaluations()[0].ZScore() = %v, want %v", i, itemEvaluations[0].ZScore().Float64(), spec.wantZScore)
				}
				if itemEvaluations[0].Rank() != spec.wantRank {
					t.Errorf("results[%d].ItemEvaluations()[0].Rank() = %v, want %v", i, itemEvaluations[0].Rank(), spec.wantRank)
				}
				if len(elementEvaluations) != 1 {
					t.Fatalf("len(results[%d].ElementEvaluations()) = %v, want %v", i, len(elementEvaluations), 1)
				}
				if elementEvaluations[0].Rank() != spec.wantRank {
					t.Errorf("results[%d].ElementEvaluations()[0].Rank() = %v, want %v", i, elementEvaluations[0].Rank(), spec.wantRank)
				}
				if motorAge == nil || motorAge.Int() != spec.wantMotorAge {
					t.Errorf("results[%d].MotorAge() = %v, want %v", i, motorAge, spec.wantMotorAge)
				}
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
			mockOrganizationRepository := mocksorganization.NewMockOrganizationRepository(ctrl)
			mockMeasurementItemRepository := mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl)
			mockAgeGroupStandardRepository := mocksstandard.NewMockAgeGroupStandardRepository(ctrl)
			mockRankStandardRepository := mocksstandard.NewMockRankStandardRepository(ctrl)
			mockPrescribedMenuOverrideRepository := mocksjudgment.NewMockPrescribedMenuOverrideRepository(ctrl)
			mockTrainingMenuRepository := mockstraining.NewMockTrainingMenuRepository(ctrl)
			mockPrescriptionRuleRepository := mockstraining.NewMockPrescriptionRuleRepository(ctrl)

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

			u := NewJudgmentUsecase(mockAuthService, mockUserService, mockJudgmentRepository, mockPrescribedMenuOverrideRepository, mockMeasurementRepository, mockCustomerRepository, mockOrganizationRepository, mockMeasurementItemRepository, mockAgeGroupStandardRepository, mockRankStandardRepository, mockTrainingMenuRepository, mockPrescriptionRuleRepository)

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

func TestUpsertPrescription(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	otherTenantID := "9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	userID := "google-oauth2|000000000000000000000"
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

	wallPushCode, _ := training.NewCode("wall_push")
	wallPushName, _ := training.NewName("壁押し")
	wallPushInstruction, _ := training.NewInstruction("肘をゆっくり曲げ伸ばしする")
	wallPushAmount, _ := training.NewAmount(10)
	wallPushSets, _ := training.NewSets(3)
	wallPushMenu := training.NewTrainingMenu(training.NewTrainingMenuID(), wallPushCode, wallPushName, measurementitem.ElementMuscleStrength, training.PartUpperLimb, wallPushAmount, training.UnitReps, wallPushSets, wallPushInstruction, createdAt, updatedAt)

	walkingCode, _ := training.NewCode("walking")
	walkingName, _ := training.NewName("ウォーキング")
	walkingInstruction, _ := training.NewInstruction("会話ができる速さで歩く")
	walkingAmount, _ := training.NewAmount(20)
	walkingSets, _ := training.NewSets(1)
	walkingMenu := training.NewTrainingMenu(training.NewTrainingMenuID(), walkingCode, walkingName, measurementitem.ElementMuscleEndurance, training.PartWholeBody, walkingAmount, training.UnitMinutes, walkingSets, walkingInstruction, createdAt, updatedAt)

	muscleStrength := measurementitem.ElementMuscleStrength
	upperLimb := training.PartUpperLimb
	flexibility := measurementitem.ElementFlexibility
	lowerLimb := training.PartLowerLimb

	editedAmount, _ := training.NewAmount(15)
	editedSets, _ := training.NewSets(5)

	repeatInput := func(menu PrescribedMenuInput, n int) []PrescribedMenuInput {
		menus := make([]PrescribedMenuInput, 0, n)
		for i := 0; i < n; i++ {
			menus = append(menus, menu)
		}
		return menus
	}

	wallPushInput := PrescribedMenuInput{Element: &muscleStrength, Part: &upperLimb, TrainingMenuID: wallPushMenu.ID(), Amount: editedAmount, Unit: training.UnitReps, Sets: editedSets}
	walkingInput := PrescribedMenuInput{TrainingMenuID: walkingMenu.ID(), Amount: walkingAmount, Unit: training.UnitMinutes, Sets: walkingSets}
	unknownInput := PrescribedMenuInput{TrainingMenuID: training.NewTrainingMenuID(), Amount: walkingAmount, Unit: training.UnitMinutes, Sets: walkingSets}

	tests := []struct {
		name                string
		success             bool
		wantErr             error
		ctx                 context.Context
		userID              string
		verifyTokenErr      error
		callFindMeasurement bool
		findMeasurementErr  error
		callFindCustomer    bool
		myTenantIDs         []string
		menus               []PrescribedMenuInput
		callFindByIDs       bool
		foundTrainingMenus  []training.TrainingMenu
		findByIDsErr        error
		callReplace         bool
		replaceErr          error
		wantSortOrders      []int
	}{
		{
			name:                "success upsert a prescription",
			success:             true,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			menus:               []PrescribedMenuInput{wallPushInput, walkingInput},
			callFindByIDs:       true,
			foundTrainingMenus:  []training.TrainingMenu{wallPushMenu, walkingMenu},
			callReplace:         true,
			wantSortOrders:      []int{1, 2},
		},
		{
			name:                "success upsert a prescription repeating a training menu",
			success:             true,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			menus:               []PrescribedMenuInput{walkingInput, walkingInput},
			callFindByIDs:       true,
			foundTrainingMenus:  []training.TrainingMenu{walkingMenu},
			callReplace:         true,
			wantSortOrders:      []int{1, 2},
		},
		{
			name:    "failure empty prescription",
			wantErr: judgment.ErrEmptyPrescription,
			ctx:     context.Background(),
			userID:  userID,
		},
		{
			name:                "failure element without part",
			wantErr:             judgment.ErrInvalidPrescribedMenuLabels,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			menus:               []PrescribedMenuInput{{Element: &muscleStrength, TrainingMenuID: wallPushMenu.ID(), Amount: editedAmount, Unit: training.UnitReps, Sets: editedSets}},
		},
		{
			name:                "failure labels contradicting the training menu",
			wantErr:             judgment.ErrInvalidPrescribedMenuLabels,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			menus:               []PrescribedMenuInput{{Element: &flexibility, Part: &lowerLimb, TrainingMenuID: wallPushMenu.ID(), Amount: editedAmount, Unit: training.UnitReps, Sets: editedSets}},
			callFindByIDs:       true,
			foundTrainingMenus:  []training.TrainingMenu{wallPushMenu},
		},
		{
			name:           "failure verify token error",
			ctx:            context.Background(),
			verifyTokenErr: errors.New("verify token error"),
			menus:          []PrescribedMenuInput{walkingInput},
		},
		{
			name:                "failure measurement not found",
			wantErr:             measurement.ErrMeasurementNotFound,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			findMeasurementErr:  measurement.ErrMeasurementNotFound,
			menus:               []PrescribedMenuInput{walkingInput},
		},
		{
			name:                "failure other tenant is hidden as not found",
			wantErr:             judgment.ErrJudgmentNotFound,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{otherTenantID},
			menus:               []PrescribedMenuInput{walkingInput},
		},
		{
			name:                "failure training menu not found",
			wantErr:             training.ErrTrainingMenuNotFound,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			menus:               []PrescribedMenuInput{walkingInput, unknownInput},
			callFindByIDs:       true,
			foundTrainingMenus:  []training.TrainingMenu{walkingMenu},
		},
		{
			name:                "failure find training menus error",
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			menus:               []PrescribedMenuInput{walkingInput},
			callFindByIDs:       true,
			findByIDsErr:        errors.New("find training menus error"),
		},
		{
			name:                "failure replace prescribed menu overrides error",
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			menus:               []PrescribedMenuInput{walkingInput},
			callFindByIDs:       true,
			foundTrainingMenus:  []training.TrainingMenu{walkingMenu},
			callReplace:         true,
			replaceErr:          errors.New("replace prescribed menu overrides error"),
			wantSortOrders:      []int{1},
		},
		{
			name:                "failure too many prescribed menus",
			wantErr:             training.ErrInvalidSortOrder,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			menus:               repeatInput(walkingInput, 32768),
			callFindByIDs:       true,
			foundTrainingMenus:  []training.TrainingMenu{walkingMenu},
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
			mockOrganizationRepository := mocksorganization.NewMockOrganizationRepository(ctrl)
			mockMeasurementItemRepository := mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl)
			mockAgeGroupStandardRepository := mocksstandard.NewMockAgeGroupStandardRepository(ctrl)
			mockRankStandardRepository := mocksstandard.NewMockRankStandardRepository(ctrl)
			mockPrescribedMenuOverrideRepository := mocksjudgment.NewMockPrescribedMenuOverrideRepository(ctrl)
			mockTrainingMenuRepository := mockstraining.NewMockTrainingMenuRepository(ctrl)
			mockPrescriptionRuleRepository := mockstraining.NewMockPrescriptionRuleRepository(ctrl)

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

			if tt.callFindByIDs {
				mockTrainingMenuRepository.EXPECT().FindByIDs(tt.ctx, gomock.Any()).Return(tt.foundTrainingMenus, tt.findByIDsErr).Times(1)
			}

			if tt.callReplace {
				mockPrescribedMenuOverrideRepository.EXPECT().ReplaceByMeasurementID(tt.ctx, measurementID, gomock.Any()).DoAndReturn(
					func(_ context.Context, _ measurement.MeasurementID, overrides []judgment.PrescribedMenuOverride) error {
						if len(overrides) != len(tt.wantSortOrders) {
							t.Errorf("len(overrides) = %v, want %v", len(overrides), len(tt.wantSortOrders))
							return tt.replaceErr
						}
						for i, wantSortOrder := range tt.wantSortOrders {
							if overrides[i].SortOrder().Int() != wantSortOrder {
								t.Errorf("overrides[%d].SortOrder() = %v, want %v", i, overrides[i].SortOrder().Int(), wantSortOrder)
							}
							if overrides[i].MeasurementID() != measurementID {
								t.Errorf("overrides[%d].MeasurementID() = %v, want %v", i, overrides[i].MeasurementID(), measurementID)
							}
						}
						return tt.replaceErr
					}).Times(1)
			} else {
				mockPrescribedMenuOverrideRepository.EXPECT().ReplaceByMeasurementID(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			}

			u := NewJudgmentUsecase(mockAuthService, mockUserService, mockJudgmentRepository, mockPrescribedMenuOverrideRepository, mockMeasurementRepository, mockCustomerRepository, mockOrganizationRepository, mockMeasurementItemRepository, mockAgeGroupStandardRepository, mockRankStandardRepository, mockTrainingMenuRepository, mockPrescriptionRuleRepository)

			upsertedPrescription, err := u.UpsertPrescription(tt.ctx, measurementID, tt.menus)
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

			prescribedMenus := upsertedPrescription.PrescribedMenus()
			if len(prescribedMenus) != len(tt.menus) {
				t.Fatalf("len(PrescribedMenus()) = %v, want %v", len(prescribedMenus), len(tt.menus))
			}
			for i, menu := range tt.menus {
				if prescribedMenus[i].Source() != judgment.PrescriptionSourceManual {
					t.Errorf("PrescribedMenus()[%d].Source() = %v, want %v", i, prescribedMenus[i].Source(), judgment.PrescriptionSourceManual)
				}
				if prescribedMenus[i].TrainingMenuID() != menu.TrainingMenuID {
					t.Errorf("PrescribedMenus()[%d].TrainingMenuID() = %v, want %v", i, prescribedMenus[i].TrainingMenuID(), menu.TrainingMenuID)
				}
				if prescribedMenus[i].Amount() != menu.Amount {
					t.Errorf("PrescribedMenus()[%d].Amount() = %v, want %v", i, prescribedMenus[i].Amount(), menu.Amount)
				}
				if prescribedMenus[i].Sets() != menu.Sets {
					t.Errorf("PrescribedMenus()[%d].Sets() = %v, want %v", i, prescribedMenus[i].Sets(), menu.Sets)
				}
			}
		})
	}
}

func TestDeletePrescription(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	otherTenantID := "9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	userID := "google-oauth2|000000000000000000000"

	tests := []struct {
		name                string
		success             bool
		wantErr             error
		ctx                 context.Context
		userID              string
		verifyTokenErr      error
		callFindMeasurement bool
		findMeasurementErr  error
		callFindCustomer    bool
		myTenantIDs         []string
		callDelete          bool
		deleteErr           error
	}{
		{
			name:                "success delete a prescription",
			success:             true,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			callDelete:          true,
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
			name:                "failure other tenant is hidden as not found",
			wantErr:             judgment.ErrJudgmentNotFound,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{otherTenantID},
		},
		{
			name:                "success delete a prescription that does not exist",
			success:             true,
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			callDelete:          true,
		},
		{
			name:                "failure delete prescription error",
			ctx:                 context.Background(),
			userID:              userID,
			callFindMeasurement: true,
			callFindCustomer:    true,
			myTenantIDs:         []string{tenantID.String()},
			callDelete:          true,
			deleteErr:           errors.New("delete prescription error"),
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
			mockOrganizationRepository := mocksorganization.NewMockOrganizationRepository(ctrl)
			mockMeasurementItemRepository := mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl)
			mockAgeGroupStandardRepository := mocksstandard.NewMockAgeGroupStandardRepository(ctrl)
			mockRankStandardRepository := mocksstandard.NewMockRankStandardRepository(ctrl)
			mockPrescribedMenuOverrideRepository := mocksjudgment.NewMockPrescribedMenuOverrideRepository(ctrl)
			mockTrainingMenuRepository := mockstraining.NewMockTrainingMenuRepository(ctrl)
			mockPrescriptionRuleRepository := mockstraining.NewMockPrescriptionRuleRepository(ctrl)

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

			if tt.callDelete {
				mockPrescribedMenuOverrideRepository.EXPECT().DeleteByMeasurementID(tt.ctx, measurementID).Return(tt.deleteErr).Times(1)
			} else {
				mockPrescribedMenuOverrideRepository.EXPECT().DeleteByMeasurementID(gomock.Any(), gomock.Any()).Times(0)
			}

			u := NewJudgmentUsecase(mockAuthService, mockUserService, mockJudgmentRepository, mockPrescribedMenuOverrideRepository, mockMeasurementRepository, mockCustomerRepository, mockOrganizationRepository, mockMeasurementItemRepository, mockAgeGroupStandardRepository, mockRankStandardRepository, mockTrainingMenuRepository, mockPrescriptionRuleRepository)

			err := u.DeletePrescription(tt.ctx, measurementID)
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
