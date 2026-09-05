package judgment

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	judgmentv1 "github.com/qkitzero/fitness-service/gen/go/judgment/v1"
	appjudgment "github.com/qkitzero/fitness-service/internal/application/judgment"
	domaincustomer "github.com/qkitzero/fitness-service/internal/domain/customer"
	domainjudgment "github.com/qkitzero/fitness-service/internal/domain/judgment"
	domainmeasurement "github.com/qkitzero/fitness-service/internal/domain/measurement"
	domainmeasurementitem "github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	domainorganization "github.com/qkitzero/fitness-service/internal/domain/organization"
	domainstaff "github.com/qkitzero/fitness-service/internal/domain/staff"
	domainstandard "github.com/qkitzero/fitness-service/internal/domain/standard"
	domaintraining "github.com/qkitzero/fitness-service/internal/domain/training"
	mocksappjudgment "github.com/qkitzero/fitness-service/mocks/application/judgment"
	mocksjudgment "github.com/qkitzero/fitness-service/mocks/domain/judgment"
)

const (
	sampleMeasurementID  = "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"
	sampleCustomerID     = "9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	sampleOrganizationID = "3f2b6c1d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	otherMeasurementID   = "1b2c3d4e-5f6a-7b8c-9d0e-1f2a3b4c5d6e"
	otherCustomerID      = "2c3d4e5f-6a7b-8c9d-0e1f-2a3b4c5d6e7f"
)

func TestGetJudgment(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

	motorFunction, _ := domainmeasurementitem.NewCategory("motor_function")
	kg, _ := domainmeasurementitem.NewUnit("kg")
	twoTrials, _ := domainmeasurementitem.NewTrialCount(2)
	higherIsBetter := domainmeasurementitem.ScoreDirectionHigherIsBetter
	gripStrengthID := domainmeasurementitem.NewMeasurementItemID()
	gripStrengthCode, _ := domainmeasurementitem.NewCode("grip_strength")
	gripStrengthName, _ := domainmeasurementitem.NewName("握力")
	gripStrength := domainmeasurementitem.NewMeasurementItem(gripStrengthID, gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, domainmeasurementitem.SideModeBilateral, domainmeasurementitem.ValueTypeNumeric, &higherIsBetter, domainmeasurementitem.SideAggregationMean, []domainmeasurementitem.Element{domainmeasurementitem.ElementMuscleStrength}, createdAt, updatedAt)
	unmappedElementItem := domainmeasurementitem.NewMeasurementItem(gripStrengthID, gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, domainmeasurementitem.SideModeBilateral, domainmeasurementitem.ValueTypeNumeric, &higherIsBetter, domainmeasurementitem.SideAggregationMean, []domainmeasurementitem.Element{domainmeasurementitem.Element("unknown")}, createdAt, updatedAt)

	oneTrial, _ := domainmeasurementitem.NewTrialCount(1)
	count, _ := domainmeasurementitem.NewUnit("count")
	cs30ID := domainmeasurementitem.NewMeasurementItemID()
	cs30Code, _ := domainmeasurementitem.NewCode("cs30")
	cs30Name, _ := domainmeasurementitem.NewName("CS-30（30秒立ち座り）")
	cs30 := domainmeasurementitem.NewMeasurementItem(cs30ID, cs30Code, cs30Name, motorFunction, count, oneTrial, domainmeasurementitem.SideModeNone, domainmeasurementitem.ValueTypeNumeric, &higherIsBetter, domainmeasurementitem.SideAggregationMean, []domainmeasurementitem.Element{domainmeasurementitem.ElementMuscleStrength}, createdAt, updatedAt)
	items := []domainmeasurementitem.MeasurementItem{gripStrength, cs30}

	ageRange6064, _ := domainstandard.NewAgeRange(60, 64)
	gripStrengthMean, _ := domainstandard.NewMean(38)
	cs30Mean, _ := domainstandard.NewMean(20)
	standardDeviation, _ := domainstandard.NewStandardDeviation(5)
	ageGroupStandards := []domainstandard.AgeGroupStandard{
		domainstandard.NewAgeGroupStandard(domainstandard.NewAgeGroupStandardID(), gripStrengthID, domainstandard.GenderMale, ageRange6064, gripStrengthMean, standardDeviation, createdAt, updatedAt),
		domainstandard.NewAgeGroupStandard(domainstandard.NewAgeGroupStandardID(), cs30ID, domainstandard.GenderMale, ageRange6064, cs30Mean, standardDeviation, createdAt, updatedAt),
	}

	zScoreA := domainstandard.ZScore(1.5)
	zScoreBMin := domainstandard.ZScore(0.5)
	zScoreBMax := domainstandard.ZScore(1.5)
	zScoreCMin := domainstandard.ZScore(-0.5)
	zScoreCMax := domainstandard.ZScore(0.5)
	rankStandardA, _ := domainstandard.NewRankStandard(domainstandard.RankA, &zScoreA, nil, createdAt, updatedAt)
	rankStandardB, _ := domainstandard.NewRankStandard(domainstandard.RankB, &zScoreBMin, &zScoreBMax, createdAt, updatedAt)
	rankStandardC, _ := domainstandard.NewRankStandard(domainstandard.RankC, &zScoreCMin, &zScoreCMax, createdAt, updatedAt)
	rankStandards := []domainstandard.RankStandard{rankStandardA, rankStandardB, rankStandardC}
	unmappedItemRankStandard, _ := domainstandard.NewRankStandard(domainstandard.Rank("F"), &zScoreA, nil, createdAt, updatedAt)
	unmappedItemRankStandards := []domainstandard.RankStandard{unmappedItemRankStandard, rankStandardC}
	unmappedElementRankStandard, _ := domainstandard.NewRankStandard(domainstandard.Rank("F"), &zScoreBMin, &zScoreBMax, createdAt, updatedAt)
	unmappedElementRankStandards := []domainstandard.RankStandard{rankStandardA, unmappedElementRankStandard, rankStandardC}

	evaluationOf := func(evaluatedItems []domainmeasurementitem.MeasurementItem, rs []domainstandard.RankStandard) domainjudgment.Evaluation {
		trialIndex, _ := domainmeasurement.NewTrialIndex(1)
		gripStrengthValue, _ := domainmeasurement.NewValue(46)
		gripStrengthEntry := domainmeasurement.ReconstructMeasurementEntry(gripStrengthID, false, nil, []domainmeasurement.MeasurementValue{
			domainmeasurement.NewMeasurementValue(trialIndex, domainmeasurement.SideLeft, &gripStrengthValue, nil, nil),
			domainmeasurement.NewMeasurementValue(trialIndex, domainmeasurement.SideRight, &gripStrengthValue, nil, nil),
		})
		cs30Value, _ := domainmeasurement.NewValue(20)
		cs30Entry := domainmeasurement.ReconstructMeasurementEntry(cs30ID, false, nil, []domainmeasurement.MeasurementValue{
			domainmeasurement.NewMeasurementValue(trialIndex, domainmeasurement.SideNone, &cs30Value, nil, nil),
		})
		measuredOn, _ := domainmeasurement.NewMeasuredOn(2026, 8, 1)
		measuredBy, _ := domainstaff.NewStaffID("google-oauth2|000000000000000000000")
		ageAtMeasurement, _ := domainmeasurement.NewAgeAtMeasurement(62)
		m := domainmeasurement.ReconstructMeasurement(domainmeasurement.NewMeasurementID(), domaincustomer.NewCustomerID(), measuredOn, measuredBy, ageAtMeasurement, measuredBy, false, []domainmeasurement.MeasurementEntry{gripStrengthEntry, cs30Entry}, createdAt, updatedAt)
		return domainjudgment.NewEvaluation(m, evaluatedItems, domainstandard.GenderMale, 62, ageGroupStandards, rs)
	}

	advice, _ := domainjudgment.NewAdvice("週2回のスクワットを継続してください")

	wallPushCode, _ := domaintraining.NewCode("wall_push")
	wallPushName, _ := domaintraining.NewName("壁押し")
	wallPushInstruction, _ := domaintraining.NewInstruction("肘をゆっくり曲げ伸ばしする")
	wallPushAmount, _ := domaintraining.NewAmount(10)
	wallPushSets, _ := domaintraining.NewSets(3)
	wallPushMenu := domaintraining.NewTrainingMenu(domaintraining.NewTrainingMenuID(), wallPushCode, wallPushName, domainmeasurementitem.ElementMuscleStrength, domaintraining.PartUpperLimb, wallPushAmount, domaintraining.UnitReps, wallPushSets, wallPushInstruction, createdAt, updatedAt)
	trainingMenus := []domaintraining.TrainingMenu{wallPushMenu}
	firstSortOrder, _ := domaintraining.NewSortOrder(1)

	prescriptionOf := func(element *domainmeasurementitem.Element, part *domaintraining.Part, unit domaintraining.Unit) domainjudgment.Prescription {
		override := domainjudgment.NewPrescribedMenuOverride(domainjudgment.NewPrescribedMenuOverrideID(), domainmeasurement.NewMeasurementID(), firstSortOrder, element, part, wallPushMenu.ID(), wallPushAmount, unit, wallPushSets, createdAt, updatedAt)
		return domainjudgment.NewPrescriptionFromOverrides([]domainjudgment.PrescribedMenuOverride{override}, trainingMenus)
	}

	muscleStrength := domainmeasurementitem.ElementMuscleStrength
	upperLimb := domaintraining.PartUpperLimb
	unknownElement := domainmeasurementitem.Element("unknown")
	unknownPart := domaintraining.Part("unknown")

	tests := []struct {
		name                   string
		measurementID          string
		callUsecase            bool
		evaluation             func() domainjudgment.Evaluation
		advice                 *domainjudgment.Advice
		isDraft                bool
		getErr                 error
		prescription           func() domainjudgment.Prescription
		unmappedSource         bool
		wantPrescribedMenus    int
		wantPrescribedLabels   bool
		wantPrescribedUnit     judgmentv1.PrescribedUnit
		wantCode               codes.Code
		wantItemEvaluations    int
		wantElementEvaluations int
		wantMotorAge           uint32
		wantAdvice             string
	}{
		{
			name:          "success get judgment",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			evaluation: func() domainjudgment.Evaluation {
				return evaluationOf(items, rankStandards)
			},
			advice: advice,
			prescription: func() domainjudgment.Prescription {
				return prescriptionOf(&muscleStrength, &upperLimb, domaintraining.UnitReps)
			},
			wantPrescribedMenus:    1,
			wantPrescribedLabels:   true,
			wantPrescribedUnit:     judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_REPS,
			wantCode:               codes.OK,
			wantItemEvaluations:    2,
			wantElementEvaluations: 1,
			wantMotorAge:           62,
			wantAdvice:             "週2回のスクワットを継続してください",
		},
		{
			name:          "success get judgment without advice",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			evaluation: func() domainjudgment.Evaluation {
				return evaluationOf(items, rankStandards)
			},
			wantCode:               codes.OK,
			wantItemEvaluations:    2,
			wantElementEvaluations: 1,
			wantMotorAge:           62,
		},
		{
			name:          "success get judgment of a draft measurement",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			evaluation: func() domainjudgment.Evaluation {
				return evaluationOf(items, rankStandards)
			},
			isDraft:                true,
			wantCode:               codes.OK,
			wantItemEvaluations:    2,
			wantElementEvaluations: 1,
			wantMotorAge:           62,
		},
		{
			name:          "success get an empty judgment",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			evaluation: func() domainjudgment.Evaluation {
				return evaluationOf(nil, rankStandards)
			},
			wantCode: codes.OK,
		},
		{
			name:          "success get judgment of a prescription without labels",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			evaluation: func() domainjudgment.Evaluation {
				return evaluationOf(nil, rankStandards)
			},
			prescription: func() domainjudgment.Prescription {
				return prescriptionOf(nil, nil, domaintraining.UnitMinutes)
			},
			wantPrescribedMenus: 1,
			wantPrescribedUnit:  judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_MINUTES,
			wantCode:            codes.OK,
		},
		{
			name:          "failure invalid measurement id",
			measurementID: "",
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure unmapped prescription source",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			evaluation: func() domainjudgment.Evaluation {
				return evaluationOf(nil, rankStandards)
			},
			unmappedSource: true,
			wantCode:       codes.Internal,
		},
		{
			name:          "failure unmapped prescribed element",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			evaluation: func() domainjudgment.Evaluation {
				return evaluationOf(nil, rankStandards)
			},
			prescription: func() domainjudgment.Prescription {
				return prescriptionOf(&unknownElement, &upperLimb, domaintraining.UnitReps)
			},
			wantCode: codes.Internal,
		},
		{
			name:          "failure unmapped prescribed part",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			evaluation: func() domainjudgment.Evaluation {
				return evaluationOf(nil, rankStandards)
			},
			prescription: func() domainjudgment.Prescription {
				return prescriptionOf(&muscleStrength, &unknownPart, domaintraining.UnitReps)
			},
			wantCode: codes.Internal,
		},
		{
			name:          "failure unmapped prescribed unit",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			evaluation: func() domainjudgment.Evaluation {
				return evaluationOf(nil, rankStandards)
			},
			prescription: func() domainjudgment.Prescription {
				return prescriptionOf(&muscleStrength, &upperLimb, domaintraining.Unit("hours"))
			},
			wantCode: codes.Internal,
		},
		{
			name:          "failure judgment not found",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			getErr:        domainjudgment.ErrJudgmentNotFound,
			wantCode:      codes.NotFound,
		},
		{
			name:          "failure measurement not found",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			getErr:        domainmeasurement.ErrMeasurementNotFound,
			wantCode:      codes.NotFound,
		},
		{
			name:          "failure customer not found",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			getErr:        domaincustomer.ErrCustomerNotFound,
			wantCode:      codes.NotFound,
		},
		{
			name:          "failure unauthenticated",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			getErr:        status.Error(codes.Unauthenticated, "unauthenticated"),
			wantCode:      codes.Unauthenticated,
		},
		{
			name:          "failure internal error",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			getErr:        errors.New("get judgment error"),
			wantCode:      codes.Internal,
		},
		{
			name:          "failure unmapped rank of an item evaluation",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			evaluation: func() domainjudgment.Evaluation {
				return evaluationOf(items, unmappedItemRankStandards)
			},
			wantCode: codes.Internal,
		},
		{
			name:          "failure unmapped rank of an element evaluation",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			evaluation: func() domainjudgment.Evaluation {
				return evaluationOf(items, unmappedElementRankStandards)
			},
			wantCode: codes.Internal,
		},
		{
			name:          "failure unmapped element",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			evaluation: func() domainjudgment.Evaluation {
				return evaluationOf([]domainmeasurementitem.MeasurementItem{unmappedElementItem}, rankStandards)
			},
			wantCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappjudgment.NewMockJudgmentUsecase(ctrl)
			if tt.callUsecase {
				var result appjudgment.JudgmentResult
				if tt.getErr == nil {
					measurementID, _ := domainmeasurement.NewMeasurementIDFromString(sampleMeasurementID)
					result = appjudgment.JudgmentResult{
						MeasurementID: measurementID,
						IsDraft:       tt.isDraft,
						Evaluation:    tt.evaluation(),
						Advice:        tt.advice,
					}
					if tt.prescription != nil {
						result.Prescription = tt.prescription()
					}
					if tt.unmappedSource {
						mockPrescribedMenu := mocksjudgment.NewMockPrescribedMenu(ctrl)
						mockPrescribedMenu.EXPECT().Source().Return(domainjudgment.PrescriptionSource("unknown")).AnyTimes()
						mockPrescription := mocksjudgment.NewMockPrescription(ctrl)
						mockPrescription.EXPECT().PrescribedMenus().Return([]domainjudgment.PrescribedMenu{mockPrescribedMenu}).AnyTimes()
						result.Prescription = mockPrescription
					}
				}
				mockUsecase.EXPECT().GetJudgment(gomock.Any(), gomock.Any()).Return(result, tt.getErr).Times(1)
			}

			handler := NewJudgmentHandler(mockUsecase)

			res, err := handler.GetJudgment(context.Background(), &judgmentv1.GetJudgmentRequest{MeasurementId: tt.measurementID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode != codes.OK {
				return
			}

			judgmentMessage := res.GetJudgment()
			if judgmentMessage.GetMeasurementId() != sampleMeasurementID {
				t.Errorf("MeasurementId = %v, want %v", judgmentMessage.GetMeasurementId(), sampleMeasurementID)
			}
			if judgmentMessage.GetIsDraft() != tt.isDraft {
				t.Errorf("IsDraft = %v, want %v", judgmentMessage.GetIsDraft(), tt.isDraft)
			}
			if len(judgmentMessage.GetItemEvaluations()) != tt.wantItemEvaluations {
				t.Errorf("len(ItemEvaluations) = %v, want %v", len(judgmentMessage.GetItemEvaluations()), tt.wantItemEvaluations)
			}
			if len(judgmentMessage.GetElementEvaluations()) != tt.wantElementEvaluations {
				t.Errorf("len(ElementEvaluations) = %v, want %v", len(judgmentMessage.GetElementEvaluations()), tt.wantElementEvaluations)
			}
			if len(judgmentMessage.GetPrescribedMenus()) != tt.wantPrescribedMenus {
				t.Errorf("len(PrescribedMenus) = %v, want %v", len(judgmentMessage.GetPrescribedMenus()), tt.wantPrescribedMenus)
			}
			if tt.wantPrescribedMenus > 0 {
				prescribedMenu := judgmentMessage.GetPrescribedMenus()[0]
				if prescribedMenu.GetSource() != judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_MANUAL {
					t.Errorf("Source = %v, want %v", prescribedMenu.GetSource(), judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_MANUAL)
				}
				if tt.wantPrescribedLabels {
					if prescribedMenu.GetElement() != judgmentv1.Element_ELEMENT_MUSCLE_STRENGTH {
						t.Errorf("Element = %v, want %v", prescribedMenu.GetElement(), judgmentv1.Element_ELEMENT_MUSCLE_STRENGTH)
					}
					if prescribedMenu.GetPart() != judgmentv1.PrescribedPart_PRESCRIBED_PART_UPPER_LIMB {
						t.Errorf("Part = %v, want %v", prescribedMenu.GetPart(), judgmentv1.PrescribedPart_PRESCRIBED_PART_UPPER_LIMB)
					}
				}
				if !tt.wantPrescribedLabels {
					if prescribedMenu.Element != nil {
						t.Errorf("Element = %v, want nil", prescribedMenu.GetElement())
					}
					if prescribedMenu.Part != nil {
						t.Errorf("Part = %v, want nil", prescribedMenu.GetPart())
					}
				}
				if prescribedMenu.GetTrainingMenuId() != wallPushMenu.ID().String() {
					t.Errorf("TrainingMenuId = %v, want %v", prescribedMenu.GetTrainingMenuId(), wallPushMenu.ID().String())
				}
				if prescribedMenu.GetTrainingMenuName() != wallPushName.String() {
					t.Errorf("TrainingMenuName = %v, want %v", prescribedMenu.GetTrainingMenuName(), wallPushName.String())
				}
				if prescribedMenu.GetAmount() != uint32(wallPushAmount.Int()) {
					t.Errorf("Amount = %v, want %v", prescribedMenu.GetAmount(), wallPushAmount.Int())
				}
				if prescribedMenu.GetUnit() != tt.wantPrescribedUnit {
					t.Errorf("Unit = %v, want %v", prescribedMenu.GetUnit(), tt.wantPrescribedUnit)
				}
				if prescribedMenu.GetSets() != uint32(wallPushSets.Int()) {
					t.Errorf("Sets = %v, want %v", prescribedMenu.GetSets(), wallPushSets.Int())
				}
			}
			if tt.wantMotorAge == 0 && judgmentMessage.MotorAge != nil {
				t.Errorf("MotorAge = %v, want nil", judgmentMessage.GetMotorAge())
			}
			if tt.wantMotorAge != 0 && judgmentMessage.GetMotorAge() != tt.wantMotorAge {
				t.Errorf("MotorAge = %v, want %v", judgmentMessage.GetMotorAge(), tt.wantMotorAge)
			}
			if tt.wantAdvice == "" && judgmentMessage.Advice != nil {
				t.Errorf("Advice = %v, want nil", judgmentMessage.GetAdvice())
			}
			if tt.wantAdvice != "" && judgmentMessage.GetAdvice() != tt.wantAdvice {
				t.Errorf("Advice = %v, want %v", judgmentMessage.GetAdvice(), tt.wantAdvice)
			}
			if tt.wantItemEvaluations > 0 {
				itemEvaluation := judgmentMessage.GetItemEvaluations()[0]
				if itemEvaluation.GetMeasurementItemId() != gripStrengthID.String() {
					t.Errorf("MeasurementItemId = %v, want %v", itemEvaluation.GetMeasurementItemId(), gripStrengthID.String())
				}
				if itemEvaluation.GetValue() != 46 {
					t.Errorf("Value = %v, want %v", itemEvaluation.GetValue(), 46)
				}
				if itemEvaluation.GetMean() != 38 {
					t.Errorf("Mean = %v, want %v", itemEvaluation.GetMean(), 38)
				}
				if itemEvaluation.GetZScore() != 1.6 {
					t.Errorf("ZScore = %v, want %v", itemEvaluation.GetZScore(), 1.6)
				}
				if itemEvaluation.GetRank() != judgmentv1.Rank_RANK_A {
					t.Errorf("Rank = %v, want %v", itemEvaluation.GetRank(), judgmentv1.Rank_RANK_A)
				}
			}
			if tt.wantElementEvaluations > 0 {
				elementEvaluation := judgmentMessage.GetElementEvaluations()[0]
				if elementEvaluation.GetElement() != judgmentv1.Element_ELEMENT_MUSCLE_STRENGTH {
					t.Errorf("Element = %v, want %v", elementEvaluation.GetElement(), judgmentv1.Element_ELEMENT_MUSCLE_STRENGTH)
				}
				if elementEvaluation.GetZScore() != 0.8 {
					t.Errorf("ZScore = %v, want %v", elementEvaluation.GetZScore(), 0.8)
				}
				if elementEvaluation.GetRank() != judgmentv1.Rank_RANK_B {
					t.Errorf("Rank = %v, want %v", elementEvaluation.GetRank(), judgmentv1.Rank_RANK_B)
				}
			}
		})
	}
}

func TestListOrganizationJudgments(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

	motorFunction, _ := domainmeasurementitem.NewCategory("motor_function")
	kg, _ := domainmeasurementitem.NewUnit("kg")
	twoTrials, _ := domainmeasurementitem.NewTrialCount(2)
	higherIsBetter := domainmeasurementitem.ScoreDirectionHigherIsBetter
	gripStrengthID := domainmeasurementitem.NewMeasurementItemID()
	gripStrengthCode, _ := domainmeasurementitem.NewCode("grip_strength")
	gripStrengthName, _ := domainmeasurementitem.NewName("握力")
	gripStrength := domainmeasurementitem.NewMeasurementItem(gripStrengthID, gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, domainmeasurementitem.SideModeBilateral, domainmeasurementitem.ValueTypeNumeric, &higherIsBetter, domainmeasurementitem.SideAggregationMean, []domainmeasurementitem.Element{domainmeasurementitem.ElementMuscleStrength}, createdAt, updatedAt)
	unmappedElementItem := domainmeasurementitem.NewMeasurementItem(gripStrengthID, gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, domainmeasurementitem.SideModeBilateral, domainmeasurementitem.ValueTypeNumeric, &higherIsBetter, domainmeasurementitem.SideAggregationMean, []domainmeasurementitem.Element{domainmeasurementitem.Element("unknown")}, createdAt, updatedAt)

	ageRange6064, _ := domainstandard.NewAgeRange(60, 64)
	gripStrengthMean, _ := domainstandard.NewMean(38)
	standardDeviation, _ := domainstandard.NewStandardDeviation(5)
	ageGroupStandards := []domainstandard.AgeGroupStandard{
		domainstandard.NewAgeGroupStandard(domainstandard.NewAgeGroupStandardID(), gripStrengthID, domainstandard.GenderMale, ageRange6064, gripStrengthMean, standardDeviation, createdAt, updatedAt),
	}

	zScoreA := domainstandard.ZScore(1.5)
	zScoreCMin := domainstandard.ZScore(-0.5)
	zScoreCMax := domainstandard.ZScore(0.5)
	rankStandardA, _ := domainstandard.NewRankStandard(domainstandard.RankA, &zScoreA, nil, createdAt, updatedAt)
	rankStandardC, _ := domainstandard.NewRankStandard(domainstandard.RankC, &zScoreCMin, &zScoreCMax, createdAt, updatedAt)
	rankStandards := []domainstandard.RankStandard{rankStandardA, rankStandardC}
	unmappedRankStandard, _ := domainstandard.NewRankStandard(domainstandard.Rank("F"), &zScoreA, nil, createdAt, updatedAt)
	unmappedRankStandards := []domainstandard.RankStandard{unmappedRankStandard, rankStandardC}

	measuredOn, _ := domainmeasurement.NewMeasuredOn(2026, 8, 1)
	ageAtMeasurement, _ := domainmeasurement.NewAgeAtMeasurement(62)

	evaluationOf := func(evaluatedItems []domainmeasurementitem.MeasurementItem, rs []domainstandard.RankStandard) domainjudgment.Evaluation {
		trialIndex, _ := domainmeasurement.NewTrialIndex(1)
		gripStrengthValue, _ := domainmeasurement.NewValue(46)
		gripStrengthEntry := domainmeasurement.ReconstructMeasurementEntry(gripStrengthID, false, nil, []domainmeasurement.MeasurementValue{
			domainmeasurement.NewMeasurementValue(trialIndex, domainmeasurement.SideLeft, &gripStrengthValue, nil, nil),
			domainmeasurement.NewMeasurementValue(trialIndex, domainmeasurement.SideRight, &gripStrengthValue, nil, nil),
		})
		measuredBy, _ := domainstaff.NewStaffID("google-oauth2|000000000000000000000")
		m := domainmeasurement.ReconstructMeasurement(domainmeasurement.NewMeasurementID(), domaincustomer.NewCustomerID(), measuredOn, measuredBy, ageAtMeasurement, measuredBy, false, []domainmeasurement.MeasurementEntry{gripStrengthEntry}, createdAt, updatedAt)
		return domainjudgment.NewEvaluation(m, evaluatedItems, domainstandard.GenderMale, 62, ageGroupStandards, rs)
	}

	motorAge62 := uint32(62)

	type judgmentSpec struct {
		customerID             string
		measurementID          string
		year                   int32
		month                  int32
		day                    int32
		age                    int
		isDraft                bool
		evaluation             func() domainjudgment.Evaluation
		wantItemEvaluations    int
		wantElementEvaluations int
		wantMotorAge           *uint32
	}

	tests := []struct {
		name            string
		organizationID  string
		includeInactive bool
		callUsecase     bool
		judgments       []judgmentSpec
		listErr         error
		wantCode        codes.Code
	}{
		{
			name:           "success list organization judgments",
			organizationID: sampleOrganizationID,
			callUsecase:    true,
			judgments: []judgmentSpec{
				{
					customerID: sampleCustomerID, measurementID: sampleMeasurementID, year: 2026, month: 8, day: 1, age: 62,
					evaluation: func() domainjudgment.Evaluation {
						return evaluationOf([]domainmeasurementitem.MeasurementItem{gripStrength}, rankStandards)
					},
					wantItemEvaluations: 1, wantElementEvaluations: 1, wantMotorAge: &motorAge62,
				},
				{
					customerID: otherCustomerID, measurementID: otherMeasurementID, year: 2025, month: 12, day: 24, age: 48, isDraft: true,
					evaluation: func() domainjudgment.Evaluation {
						return evaluationOf([]domainmeasurementitem.MeasurementItem{gripStrength}, rankStandards)
					},
					wantItemEvaluations: 1, wantElementEvaluations: 1, wantMotorAge: &motorAge62,
				},
			},
			wantCode: codes.OK,
		},
		{
			name:            "success list organization judgments including inactive customers",
			organizationID:  sampleOrganizationID,
			includeInactive: true,
			callUsecase:     true,
			judgments: []judgmentSpec{
				{
					customerID: sampleCustomerID, measurementID: sampleMeasurementID, year: 2026, month: 8, day: 1, age: 62,
					evaluation: func() domainjudgment.Evaluation {
						return evaluationOf([]domainmeasurementitem.MeasurementItem{gripStrength}, rankStandards)
					},
					wantItemEvaluations: 1, wantElementEvaluations: 1, wantMotorAge: &motorAge62,
				},
			},
			wantCode: codes.OK,
		},
		{
			name:           "success list an empty judgment",
			organizationID: sampleOrganizationID,
			callUsecase:    true,
			judgments: []judgmentSpec{
				{
					customerID: sampleCustomerID, measurementID: sampleMeasurementID, year: 2026, month: 8, day: 1, age: 62,
					evaluation: func() domainjudgment.Evaluation {
						return evaluationOf(nil, rankStandards)
					},
				},
			},
			wantCode: codes.OK,
		},
		{
			name:           "success list no judgments",
			organizationID: sampleOrganizationID,
			callUsecase:    true,
			wantCode:       codes.OK,
		},
		{
			name:           "failure invalid organization id",
			organizationID: "",
			wantCode:       codes.InvalidArgument,
		},
		{
			name:           "failure organization not found",
			organizationID: sampleOrganizationID,
			callUsecase:    true,
			listErr:        domainorganization.ErrOrganizationNotFound,
			wantCode:       codes.NotFound,
		},
		{
			name:           "failure unauthenticated",
			organizationID: sampleOrganizationID,
			callUsecase:    true,
			listErr:        status.Error(codes.Unauthenticated, "unauthenticated"),
			wantCode:       codes.Unauthenticated,
		},
		{
			name:           "failure internal error",
			organizationID: sampleOrganizationID,
			callUsecase:    true,
			listErr:        errors.New("list organization judgments error"),
			wantCode:       codes.Internal,
		},
		{
			name:           "failure unmapped rank of an item evaluation",
			organizationID: sampleOrganizationID,
			callUsecase:    true,
			judgments: []judgmentSpec{
				{
					customerID: sampleCustomerID, measurementID: sampleMeasurementID, year: 2026, month: 8, day: 1, age: 62,
					evaluation: func() domainjudgment.Evaluation {
						return evaluationOf([]domainmeasurementitem.MeasurementItem{gripStrength}, unmappedRankStandards)
					},
				},
			},
			wantCode: codes.Internal,
		},
		{
			name:           "failure unmapped element of an element evaluation",
			organizationID: sampleOrganizationID,
			callUsecase:    true,
			judgments: []judgmentSpec{
				{
					customerID: sampleCustomerID, measurementID: sampleMeasurementID, year: 2026, month: 8, day: 1, age: 62,
					evaluation: func() domainjudgment.Evaluation {
						return evaluationOf([]domainmeasurementitem.MeasurementItem{unmappedElementItem}, rankStandards)
					},
				},
			},
			wantCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappjudgment.NewMockJudgmentUsecase(ctrl)
			if tt.callUsecase {
				var results []appjudgment.OrganizationJudgmentResult
				if tt.listErr == nil {
					results = make([]appjudgment.OrganizationJudgmentResult, 0, len(tt.judgments))
					for _, spec := range tt.judgments {
						customerID, _ := domaincustomer.NewCustomerIDFromString(spec.customerID)
						measurementID, _ := domainmeasurement.NewMeasurementIDFromString(spec.measurementID)
						specMeasuredOn, _ := domainmeasurement.NewMeasuredOn(spec.year, spec.month, spec.day)
						specAgeAtMeasurement, _ := domainmeasurement.NewAgeAtMeasurement(spec.age)
						results = append(results, appjudgment.OrganizationJudgmentResult{
							CustomerID:       customerID,
							MeasurementID:    measurementID,
							MeasuredOn:       specMeasuredOn,
							AgeAtMeasurement: specAgeAtMeasurement,
							IsDraft:          spec.isDraft,
							Evaluation:       spec.evaluation(),
						})
					}
				}
				organizationID, _ := domainorganization.NewOrganizationIDFromString(tt.organizationID)
				mockUsecase.EXPECT().ListOrganizationJudgments(gomock.Any(), organizationID, tt.includeInactive).Return(results, tt.listErr).Times(1)
			}

			handler := NewJudgmentHandler(mockUsecase)

			res, err := handler.ListOrganizationJudgments(context.Background(), &judgmentv1.ListOrganizationJudgmentsRequest{OrganizationId: tt.organizationID, IncludeInactive: tt.includeInactive})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode != codes.OK {
				return
			}

			judgmentMessages := res.GetJudgments()
			if len(judgmentMessages) != len(tt.judgments) {
				t.Fatalf("len(Judgments) = %v, want %v", len(judgmentMessages), len(tt.judgments))
			}
			for i, spec := range tt.judgments {
				judgmentMessage := judgmentMessages[i]
				if judgmentMessage.GetCustomerId() != spec.customerID {
					t.Errorf("Judgments[%d].CustomerId = %v, want %v", i, judgmentMessage.GetCustomerId(), spec.customerID)
				}
				if judgmentMessage.GetMeasurementId() != spec.measurementID {
					t.Errorf("Judgments[%d].MeasurementId = %v, want %v", i, judgmentMessage.GetMeasurementId(), spec.measurementID)
				}
				measuredOnMessage := judgmentMessage.GetMeasuredOn()
				if measuredOnMessage.GetYear() != spec.year || measuredOnMessage.GetMonth() != spec.month || measuredOnMessage.GetDay() != spec.day {
					t.Errorf("Judgments[%d].MeasuredOn = %v, want %04d-%02d-%02d", i, measuredOnMessage, spec.year, spec.month, spec.day)
				}
				if judgmentMessage.GetAgeAtMeasurement() != uint32(spec.age) {
					t.Errorf("Judgments[%d].AgeAtMeasurement = %v, want %v", i, judgmentMessage.GetAgeAtMeasurement(), spec.age)
				}
				if judgmentMessage.GetIsDraft() != spec.isDraft {
					t.Errorf("Judgments[%d].IsDraft = %v, want %v", i, judgmentMessage.GetIsDraft(), spec.isDraft)
				}
				if len(judgmentMessage.GetItemEvaluations()) != spec.wantItemEvaluations {
					t.Fatalf("len(Judgments[%d].ItemEvaluations) = %v, want %v", i, len(judgmentMessage.GetItemEvaluations()), spec.wantItemEvaluations)
				}
				if len(judgmentMessage.GetElementEvaluations()) != spec.wantElementEvaluations {
					t.Fatalf("len(Judgments[%d].ElementEvaluations) = %v, want %v", i, len(judgmentMessage.GetElementEvaluations()), spec.wantElementEvaluations)
				}
				switch {
				case spec.wantMotorAge == nil && judgmentMessage.MotorAge != nil:
					t.Errorf("Judgments[%d].MotorAge = %v, want nil", i, judgmentMessage.GetMotorAge())
				case spec.wantMotorAge != nil && judgmentMessage.MotorAge == nil:
					t.Errorf("Judgments[%d].MotorAge = nil, want %v", i, *spec.wantMotorAge)
				case spec.wantMotorAge != nil && *judgmentMessage.MotorAge != *spec.wantMotorAge:
					t.Errorf("Judgments[%d].MotorAge = %v, want %v", i, *judgmentMessage.MotorAge, *spec.wantMotorAge)
				}
				if spec.wantItemEvaluations > 0 {
					itemEvaluation := judgmentMessage.GetItemEvaluations()[0]
					if itemEvaluation.GetMeasurementItemId() != gripStrengthID.String() {
						t.Errorf("Judgments[%d].ItemEvaluations[0].MeasurementItemId = %v, want %v", i, itemEvaluation.GetMeasurementItemId(), gripStrengthID.String())
					}
					if itemEvaluation.GetZScore() != 1.6 {
						t.Errorf("Judgments[%d].ItemEvaluations[0].ZScore = %v, want %v", i, itemEvaluation.GetZScore(), 1.6)
					}
					if itemEvaluation.GetRank() != judgmentv1.Rank_RANK_A {
						t.Errorf("Judgments[%d].ItemEvaluations[0].Rank = %v, want %v", i, itemEvaluation.GetRank(), judgmentv1.Rank_RANK_A)
					}
				}
				if spec.wantElementEvaluations > 0 {
					elementEvaluation := judgmentMessage.GetElementEvaluations()[0]
					if elementEvaluation.GetElement() != judgmentv1.Element_ELEMENT_MUSCLE_STRENGTH {
						t.Errorf("Judgments[%d].ElementEvaluations[0].Element = %v, want %v", i, elementEvaluation.GetElement(), judgmentv1.Element_ELEMENT_MUSCLE_STRENGTH)
					}
					if elementEvaluation.GetRank() != judgmentv1.Rank_RANK_A {
						t.Errorf("Judgments[%d].ElementEvaluations[0].Rank = %v, want %v", i, elementEvaluation.GetRank(), judgmentv1.Rank_RANK_A)
					}
				}
			}
		})
	}
}

func TestToProtoRank(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		rank    domainstandard.Rank
		want    judgmentv1.Rank
	}{
		{"success rank a", true, domainstandard.RankA, judgmentv1.Rank_RANK_A},
		{"success rank b", true, domainstandard.RankB, judgmentv1.Rank_RANK_B},
		{"success rank c", true, domainstandard.RankC, judgmentv1.Rank_RANK_C},
		{"success rank d", true, domainstandard.RankD, judgmentv1.Rank_RANK_D},
		{"success rank e", true, domainstandard.RankE, judgmentv1.Rank_RANK_E},
		{"failure unmapped rank", false, domainstandard.Rank("F"), judgmentv1.Rank_RANK_UNSPECIFIED},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rank, err := toProtoRank(tt.rank)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if rank != tt.want {
				t.Errorf("toProtoRank() = %v, want %v", rank, tt.want)
			}
		})
	}
}

func TestToProtoElement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		element domainmeasurementitem.Element
		want    judgmentv1.Element
	}{
		{"success muscle strength", true, domainmeasurementitem.ElementMuscleStrength, judgmentv1.Element_ELEMENT_MUSCLE_STRENGTH},
		{"success muscle endurance", true, domainmeasurementitem.ElementMuscleEndurance, judgmentv1.Element_ELEMENT_MUSCLE_ENDURANCE},
		{"success flexibility", true, domainmeasurementitem.ElementFlexibility, judgmentv1.Element_ELEMENT_FLEXIBILITY},
		{"success agility", true, domainmeasurementitem.ElementAgility, judgmentv1.Element_ELEMENT_AGILITY},
		{"success balance", true, domainmeasurementitem.ElementBalance, judgmentv1.Element_ELEMENT_BALANCE},
		{"success mobility", true, domainmeasurementitem.ElementMobility, judgmentv1.Element_ELEMENT_MOBILITY},
		{"failure unmapped element", false, domainmeasurementitem.Element("unknown"), judgmentv1.Element_ELEMENT_UNSPECIFIED},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			element, err := toProtoElement(tt.element)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if element != tt.want {
				t.Errorf("toProtoElement() = %v, want %v", element, tt.want)
			}
		})
	}
}

func TestUpsertJudgmentAdvice(t *testing.T) {
	t.Parallel()
	adviceText := "週2回のスクワットを継続してください"
	invalidAdvice := "週2回の\x00スクワット"
	tooLongAdvice := strings.Repeat("あ", 2001)
	blank := ""
	parsedAdvice := domainjudgment.Advice(adviceText)

	tests := []struct {
		name          string
		measurementID string
		advice        *string
		callUsecase   bool
		wantPatch     appjudgment.AdvicePatch
		wantAdvice    string
		upsertErr     error
		wantCode      codes.Code
	}{
		{
			name:          "success upsert judgment advice",
			measurementID: sampleMeasurementID,
			advice:        &adviceText,
			callUsecase:   true,
			wantPatch:     appjudgment.AdvicePatch{Advice: &parsedAdvice, HasAdvice: true},
			wantAdvice:    adviceText,
			wantCode:      codes.OK,
		},
		{
			name:          "success clear judgment advice",
			measurementID: sampleMeasurementID,
			advice:        &blank,
			callUsecase:   true,
			wantPatch:     appjudgment.AdvicePatch{HasAdvice: true},
			wantCode:      codes.OK,
		},
		{
			name:          "success omit the advice field",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			wantPatch:     appjudgment.AdvicePatch{},
			wantAdvice:    adviceText,
			wantCode:      codes.OK,
		},
		{
			name:          "failure invalid measurement id",
			measurementID: "",
			advice:        &adviceText,
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure invalid advice",
			measurementID: sampleMeasurementID,
			advice:        &invalidAdvice,
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure too long advice",
			measurementID: sampleMeasurementID,
			advice:        &tooLongAdvice,
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure judgment not found",
			measurementID: sampleMeasurementID,
			advice:        &adviceText,
			callUsecase:   true,
			wantPatch:     appjudgment.AdvicePatch{Advice: &parsedAdvice, HasAdvice: true},
			upsertErr:     domainjudgment.ErrJudgmentNotFound,
			wantCode:      codes.NotFound,
		},
		{
			name:          "failure unauthenticated",
			measurementID: sampleMeasurementID,
			advice:        &adviceText,
			callUsecase:   true,
			wantPatch:     appjudgment.AdvicePatch{Advice: &parsedAdvice, HasAdvice: true},
			upsertErr:     status.Error(codes.Unauthenticated, "unauthenticated"),
			wantCode:      codes.Unauthenticated,
		},
		{
			name:          "failure internal error",
			measurementID: sampleMeasurementID,
			advice:        &adviceText,
			callUsecase:   true,
			wantPatch:     appjudgment.AdvicePatch{Advice: &parsedAdvice, HasAdvice: true},
			upsertErr:     errors.New("upsert judgment advice error"),
			wantCode:      codes.Internal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappjudgment.NewMockJudgmentUsecase(ctrl)
			if tt.callUsecase {
				var upsertedJudgment domainjudgment.Judgment
				if tt.upsertErr == nil {
					var upsertedAdvice *domainjudgment.Advice
					if tt.wantAdvice != "" {
						a := domainjudgment.Advice(tt.wantAdvice)
						upsertedAdvice = &a
					}
					measurementID, _ := domainmeasurement.NewMeasurementIDFromString(sampleMeasurementID)
					mockJudgment := mocksjudgment.NewMockJudgment(ctrl)
					mockJudgment.EXPECT().MeasurementID().Return(measurementID).AnyTimes()
					mockJudgment.EXPECT().Advice().Return(upsertedAdvice).AnyTimes()
					upsertedJudgment = mockJudgment
				}
				mockUsecase.EXPECT().UpsertJudgmentAdvice(gomock.Any(), gomock.Any(), tt.wantPatch).Return(upsertedJudgment, tt.upsertErr).Times(1)
			}

			handler := NewJudgmentHandler(mockUsecase)

			res, err := handler.UpsertJudgmentAdvice(context.Background(), &judgmentv1.UpsertJudgmentAdviceRequest{MeasurementId: tt.measurementID, Advice: tt.advice})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode != codes.OK {
				return
			}

			if res.GetMeasurementId() != sampleMeasurementID {
				t.Errorf("MeasurementId = %v, want %v", res.GetMeasurementId(), sampleMeasurementID)
			}
			if tt.wantAdvice == "" && res.Advice != nil {
				t.Errorf("Advice = %v, want nil", res.GetAdvice())
			}
			if tt.wantAdvice != "" && res.GetAdvice() != tt.wantAdvice {
				t.Errorf("Advice = %v, want %v", res.GetAdvice(), tt.wantAdvice)
			}
		})
	}
}

func TestToProtoPart(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		part    domaintraining.Part
		want    judgmentv1.PrescribedPart
	}{
		{"success upper limb", true, domaintraining.PartUpperLimb, judgmentv1.PrescribedPart_PRESCRIBED_PART_UPPER_LIMB},
		{"success lower limb", true, domaintraining.PartLowerLimb, judgmentv1.PrescribedPart_PRESCRIBED_PART_LOWER_LIMB},
		{"success whole body", true, domaintraining.PartWholeBody, judgmentv1.PrescribedPart_PRESCRIBED_PART_WHOLE_BODY},
		{"failure unmapped part", false, domaintraining.Part("unknown"), judgmentv1.PrescribedPart_PRESCRIBED_PART_UNSPECIFIED},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			part, err := toProtoPart(tt.part)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if part != tt.want {
				t.Errorf("toProtoPart() = %v, want %v", part, tt.want)
			}
		})
	}
}

func TestToProtoUnit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		unit    domaintraining.Unit
		want    judgmentv1.PrescribedUnit
	}{
		{"success reps", true, domaintraining.UnitReps, judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_REPS},
		{"success seconds", true, domaintraining.UnitSeconds, judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_SECONDS},
		{"success minutes", true, domaintraining.UnitMinutes, judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_MINUTES},
		{"failure unmapped unit", false, domaintraining.Unit("hours"), judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_UNSPECIFIED},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			unit, err := toProtoUnit(tt.unit)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if unit != tt.want {
				t.Errorf("toProtoUnit() = %v, want %v", unit, tt.want)
			}
		})
	}
}

func TestToProtoPrescriptionSource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		source  domainjudgment.PrescriptionSource
		want    judgmentv1.PrescriptionSource
	}{
		{"success element", true, domainjudgment.PrescriptionSourceElement, judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_ELEMENT},
		{"success fixed", true, domainjudgment.PrescriptionSourceFixed, judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_FIXED},
		{"success age decade", true, domainjudgment.PrescriptionSourceAgeDecade, judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_AGE_DECADE},
		{"success manual", true, domainjudgment.PrescriptionSourceManual, judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_MANUAL},
		{"failure unmapped prescription source", false, domainjudgment.PrescriptionSource("unknown"), judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_UNSPECIFIED},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			source, err := toProtoPrescriptionSource(tt.source)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if source != tt.want {
				t.Errorf("toProtoPrescriptionSource() = %v, want %v", source, tt.want)
			}
		})
	}
}

func TestToDomainElement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		element judgmentv1.Element
		want    domainmeasurementitem.Element
	}{
		{"success muscle strength", true, judgmentv1.Element_ELEMENT_MUSCLE_STRENGTH, domainmeasurementitem.ElementMuscleStrength},
		{"success muscle endurance", true, judgmentv1.Element_ELEMENT_MUSCLE_ENDURANCE, domainmeasurementitem.ElementMuscleEndurance},
		{"success flexibility", true, judgmentv1.Element_ELEMENT_FLEXIBILITY, domainmeasurementitem.ElementFlexibility},
		{"success agility", true, judgmentv1.Element_ELEMENT_AGILITY, domainmeasurementitem.ElementAgility},
		{"success balance", true, judgmentv1.Element_ELEMENT_BALANCE, domainmeasurementitem.ElementBalance},
		{"success mobility", true, judgmentv1.Element_ELEMENT_MOBILITY, domainmeasurementitem.ElementMobility},
		{"failure unspecified element", false, judgmentv1.Element_ELEMENT_UNSPECIFIED, domainmeasurementitem.Element("")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			element, err := toDomainElement(tt.element)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if element != tt.want {
				t.Errorf("toDomainElement() = %v, want %v", element, tt.want)
			}
		})
	}
}

func TestToDomainPart(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		part    judgmentv1.PrescribedPart
		want    domaintraining.Part
	}{
		{"success upper limb", true, judgmentv1.PrescribedPart_PRESCRIBED_PART_UPPER_LIMB, domaintraining.PartUpperLimb},
		{"success lower limb", true, judgmentv1.PrescribedPart_PRESCRIBED_PART_LOWER_LIMB, domaintraining.PartLowerLimb},
		{"success whole body", true, judgmentv1.PrescribedPart_PRESCRIBED_PART_WHOLE_BODY, domaintraining.PartWholeBody},
		{"failure unspecified part", false, judgmentv1.PrescribedPart_PRESCRIBED_PART_UNSPECIFIED, domaintraining.Part("")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			part, err := toDomainPart(tt.part)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if part != tt.want {
				t.Errorf("toDomainPart() = %v, want %v", part, tt.want)
			}
		})
	}
}

func TestToDomainUnit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		unit    judgmentv1.PrescribedUnit
		want    domaintraining.Unit
	}{
		{"success reps", true, judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_REPS, domaintraining.UnitReps},
		{"success seconds", true, judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_SECONDS, domaintraining.UnitSeconds},
		{"success minutes", true, judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_MINUTES, domaintraining.UnitMinutes},
		{"failure unspecified unit", false, judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_UNSPECIFIED, domaintraining.Unit("")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			unit, err := toDomainUnit(tt.unit)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if unit != tt.want {
				t.Errorf("toDomainUnit() = %v, want %v", unit, tt.want)
			}
		})
	}
}

func TestUpsertPrescription(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

	wallPushCode, _ := domaintraining.NewCode("wall_push")
	wallPushName, _ := domaintraining.NewName("壁押し")
	wallPushInstruction, _ := domaintraining.NewInstruction("肘をゆっくり曲げ伸ばしする")
	wallPushAmount, _ := domaintraining.NewAmount(10)
	wallPushSets, _ := domaintraining.NewSets(3)
	wallPushMenu := domaintraining.NewTrainingMenu(domaintraining.NewTrainingMenuID(), wallPushCode, wallPushName, domainmeasurementitem.ElementMuscleStrength, domaintraining.PartUpperLimb, wallPushAmount, domaintraining.UnitReps, wallPushSets, wallPushInstruction, createdAt, updatedAt)
	firstSortOrder, _ := domaintraining.NewSortOrder(1)
	muscleStrength := domainmeasurementitem.ElementMuscleStrength
	upperLimb := domaintraining.PartUpperLimb

	prescription := func() domainjudgment.Prescription {
		override := domainjudgment.NewPrescribedMenuOverride(domainjudgment.NewPrescribedMenuOverrideID(), domainmeasurement.NewMeasurementID(), firstSortOrder, &muscleStrength, &upperLimb, wallPushMenu.ID(), wallPushAmount, domaintraining.UnitReps, wallPushSets, createdAt, updatedAt)
		return domainjudgment.NewPrescriptionFromOverrides([]domainjudgment.PrescribedMenuOverride{override}, []domaintraining.TrainingMenu{wallPushMenu})
	}

	validMenu := &judgmentv1.PrescribedMenuInput{
		Element:        judgmentv1.Element_ELEMENT_MUSCLE_STRENGTH.Enum(),
		Part:           judgmentv1.PrescribedPart_PRESCRIBED_PART_UPPER_LIMB.Enum(),
		TrainingMenuId: wallPushMenu.ID().String(),
		Amount:         10,
		Unit:           judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_REPS,
		Sets:           3,
	}

	tests := []struct {
		name                string
		measurementID       string
		menus               []*judgmentv1.PrescribedMenuInput
		callUsecase         bool
		unmappedSource      bool
		upsertErr           error
		wantCode            codes.Code
		wantPrescribedMenus int
	}{
		{
			name:                "success upsert prescription",
			measurementID:       sampleMeasurementID,
			menus:               []*judgmentv1.PrescribedMenuInput{validMenu},
			callUsecase:         true,
			wantCode:            codes.OK,
			wantPrescribedMenus: 1,
		},
		{
			name:          "success upsert a prescription without labels",
			measurementID: sampleMeasurementID,
			menus: []*judgmentv1.PrescribedMenuInput{{
				TrainingMenuId: wallPushMenu.ID().String(),
				Amount:         20,
				Unit:           judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_MINUTES,
				Sets:           1,
			}},
			callUsecase:         true,
			wantCode:            codes.OK,
			wantPrescribedMenus: 1,
		},
		{
			name:          "failure empty prescription",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			upsertErr:     domainjudgment.ErrEmptyPrescription,
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure invalid prescribed menu labels",
			measurementID: sampleMeasurementID,
			menus:         []*judgmentv1.PrescribedMenuInput{validMenu},
			callUsecase:   true,
			upsertErr:     domainjudgment.ErrInvalidPrescribedMenuLabels,
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure too many prescribed menus",
			measurementID: sampleMeasurementID,
			menus:         []*judgmentv1.PrescribedMenuInput{validMenu},
			callUsecase:   true,
			upsertErr:     domaintraining.ErrInvalidSortOrder,
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure invalid measurement id",
			measurementID: "",
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure invalid training menu id",
			measurementID: sampleMeasurementID,
			menus:         []*judgmentv1.PrescribedMenuInput{{TrainingMenuId: "", Amount: 10, Unit: judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_REPS, Sets: 3}},
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure unspecified element",
			measurementID: sampleMeasurementID,
			menus:         []*judgmentv1.PrescribedMenuInput{{Element: judgmentv1.Element_ELEMENT_UNSPECIFIED.Enum(), TrainingMenuId: wallPushMenu.ID().String(), Amount: 10, Unit: judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_REPS, Sets: 3}},
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure unspecified part",
			measurementID: sampleMeasurementID,
			menus:         []*judgmentv1.PrescribedMenuInput{{Part: judgmentv1.PrescribedPart_PRESCRIBED_PART_UNSPECIFIED.Enum(), TrainingMenuId: wallPushMenu.ID().String(), Amount: 10, Unit: judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_REPS, Sets: 3}},
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure invalid amount",
			measurementID: sampleMeasurementID,
			menus:         []*judgmentv1.PrescribedMenuInput{{TrainingMenuId: wallPushMenu.ID().String(), Amount: 0, Unit: judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_REPS, Sets: 3}},
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure unspecified unit",
			measurementID: sampleMeasurementID,
			menus:         []*judgmentv1.PrescribedMenuInput{{TrainingMenuId: wallPushMenu.ID().String(), Amount: 10, Unit: judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_UNSPECIFIED, Sets: 3}},
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure invalid sets",
			measurementID: sampleMeasurementID,
			menus:         []*judgmentv1.PrescribedMenuInput{{TrainingMenuId: wallPushMenu.ID().String(), Amount: 10, Unit: judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_REPS, Sets: 0}},
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure training menu not found",
			measurementID: sampleMeasurementID,
			menus:         []*judgmentv1.PrescribedMenuInput{validMenu},
			callUsecase:   true,
			upsertErr:     domaintraining.ErrTrainingMenuNotFound,
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure judgment not found",
			measurementID: sampleMeasurementID,
			menus:         []*judgmentv1.PrescribedMenuInput{validMenu},
			callUsecase:   true,
			upsertErr:     domainjudgment.ErrJudgmentNotFound,
			wantCode:      codes.NotFound,
		},
		{
			name:           "failure unmapped prescription source",
			measurementID:  sampleMeasurementID,
			menus:          []*judgmentv1.PrescribedMenuInput{validMenu},
			callUsecase:    true,
			unmappedSource: true,
			wantCode:       codes.Internal,
		},
		{
			name:          "failure upsert prescription error",
			measurementID: sampleMeasurementID,
			menus:         []*judgmentv1.PrescribedMenuInput{validMenu},
			callUsecase:   true,
			upsertErr:     errors.New("upsert prescription error"),
			wantCode:      codes.Internal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappjudgment.NewMockJudgmentUsecase(ctrl)
			if tt.callUsecase {
				var upsertedPrescription domainjudgment.Prescription
				if tt.upsertErr == nil {
					upsertedPrescription = prescription()
				}
				if tt.unmappedSource {
					mockPrescribedMenu := mocksjudgment.NewMockPrescribedMenu(ctrl)
					mockPrescribedMenu.EXPECT().Source().Return(domainjudgment.PrescriptionSource("unknown")).AnyTimes()
					mockPrescription := mocksjudgment.NewMockPrescription(ctrl)
					mockPrescription.EXPECT().PrescribedMenus().Return([]domainjudgment.PrescribedMenu{mockPrescribedMenu}).AnyTimes()
					upsertedPrescription = mockPrescription
				}
				mockUsecase.EXPECT().UpsertPrescription(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(_ context.Context, gotMeasurementID domainmeasurement.MeasurementID, gotMenus []appjudgment.PrescribedMenuInput) (domainjudgment.Prescription, error) {
						if gotMeasurementID.String() != tt.measurementID {
							t.Errorf("MeasurementID = %v, want %v", gotMeasurementID.String(), tt.measurementID)
						}
						if len(gotMenus) != len(tt.menus) {
							t.Fatalf("len(menus) = %v, want %v", len(gotMenus), len(tt.menus))
						}
						for i, want := range tt.menus {
							got := gotMenus[i]
							switch {
							case want.Element == nil:
								if got.Element != nil {
									t.Errorf("menus[%d].Element = %v, want nil", i, *got.Element)
								}
							case got.Element == nil:
								t.Errorf("menus[%d].Element = nil, want %v", i, want.GetElement())
							case *got.Element != domainmeasurementitem.ElementMuscleStrength:
								t.Errorf("menus[%d].Element = %v, want %v", i, *got.Element, domainmeasurementitem.ElementMuscleStrength)
							}
							switch {
							case want.Part == nil:
								if got.Part != nil {
									t.Errorf("menus[%d].Part = %v, want nil", i, *got.Part)
								}
							case got.Part == nil:
								t.Errorf("menus[%d].Part = nil, want %v", i, want.GetPart())
							case *got.Part != domaintraining.PartUpperLimb:
								t.Errorf("menus[%d].Part = %v, want %v", i, *got.Part, domaintraining.PartUpperLimb)
							}
							if got.TrainingMenuID.String() != want.GetTrainingMenuId() {
								t.Errorf("menus[%d].TrainingMenuID = %v, want %v", i, got.TrainingMenuID, want.GetTrainingMenuId())
							}
							if got.Amount.Int() != int(want.GetAmount()) {
								t.Errorf("menus[%d].Amount = %v, want %v", i, got.Amount.Int(), want.GetAmount())
							}
							if got.Sets.Int() != int(want.GetSets()) {
								t.Errorf("menus[%d].Sets = %v, want %v", i, got.Sets.Int(), want.GetSets())
							}
							wantUnit, _ := toDomainUnit(want.GetUnit())
							if got.Unit != wantUnit {
								t.Errorf("menus[%d].Unit = %v, want %v", i, got.Unit, wantUnit)
							}
						}
						return upsertedPrescription, tt.upsertErr
					}).Times(1)
			} else {
				mockUsecase.EXPECT().UpsertPrescription(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			}

			handler := NewJudgmentHandler(mockUsecase)

			res, err := handler.UpsertPrescription(context.Background(), &judgmentv1.UpsertPrescriptionRequest{MeasurementId: tt.measurementID, PrescribedMenus: tt.menus})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode != codes.OK {
				return
			}

			if res.GetMeasurementId() != sampleMeasurementID {
				t.Errorf("MeasurementId = %v, want %v", res.GetMeasurementId(), sampleMeasurementID)
			}
			if len(res.GetPrescribedMenus()) != tt.wantPrescribedMenus {
				t.Fatalf("len(PrescribedMenus) = %v, want %v", len(res.GetPrescribedMenus()), tt.wantPrescribedMenus)
			}
			if tt.wantPrescribedMenus > 0 {
				prescribedMenu := res.GetPrescribedMenus()[0]
				if prescribedMenu.GetSource() != judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_MANUAL {
					t.Errorf("Source = %v, want %v", prescribedMenu.GetSource(), judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_MANUAL)
				}
				if prescribedMenu.GetTrainingMenuId() != wallPushMenu.ID().String() {
					t.Errorf("TrainingMenuId = %v, want %v", prescribedMenu.GetTrainingMenuId(), wallPushMenu.ID().String())
				}
				if prescribedMenu.GetTrainingMenuName() != wallPushName.String() {
					t.Errorf("TrainingMenuName = %v, want %v", prescribedMenu.GetTrainingMenuName(), wallPushName.String())
				}
			}
		})
	}
}

func TestDeletePrescription(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		measurementID string
		callUsecase   bool
		deleteErr     error
		wantCode      codes.Code
	}{
		{
			name:          "success delete prescription",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			wantCode:      codes.OK,
		},
		{
			name:          "failure invalid measurement id",
			measurementID: "",
			wantCode:      codes.InvalidArgument,
		},
		{
			name:          "failure judgment not found",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			deleteErr:     domainjudgment.ErrJudgmentNotFound,
			wantCode:      codes.NotFound,
		},
		{
			name:          "failure delete prescription error",
			measurementID: sampleMeasurementID,
			callUsecase:   true,
			deleteErr:     errors.New("delete prescription error"),
			wantCode:      codes.Internal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappjudgment.NewMockJudgmentUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().DeletePrescription(gomock.Any(), gomock.Any()).Return(tt.deleteErr).Times(1)
			} else {
				mockUsecase.EXPECT().DeletePrescription(gomock.Any(), gomock.Any()).Times(0)
			}

			handler := NewJudgmentHandler(mockUsecase)

			_, err := handler.DeletePrescription(context.Background(), &judgmentv1.DeletePrescriptionRequest{MeasurementId: tt.measurementID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}
