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
	domainstaff "github.com/qkitzero/fitness-service/internal/domain/staff"
	domainstandard "github.com/qkitzero/fitness-service/internal/domain/standard"
	mocksappjudgment "github.com/qkitzero/fitness-service/mocks/application/judgment"
	mocksjudgment "github.com/qkitzero/fitness-service/mocks/domain/judgment"
)

const sampleMeasurementID = "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"

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
	gripStrength := domainmeasurementitem.NewMeasurementItem(gripStrengthID, gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, true, domainmeasurementitem.ValueTypeNumeric, &higherIsBetter, domainmeasurementitem.SideAggregationMean, []domainmeasurementitem.Element{domainmeasurementitem.ElementMuscleStrength}, createdAt, updatedAt)
	unmappedElementItem := domainmeasurementitem.NewMeasurementItem(gripStrengthID, gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, true, domainmeasurementitem.ValueTypeNumeric, &higherIsBetter, domainmeasurementitem.SideAggregationMean, []domainmeasurementitem.Element{domainmeasurementitem.Element("unknown")}, createdAt, updatedAt)

	oneTrial, _ := domainmeasurementitem.NewTrialCount(1)
	count, _ := domainmeasurementitem.NewUnit("count")
	cs30ID := domainmeasurementitem.NewMeasurementItemID()
	cs30Code, _ := domainmeasurementitem.NewCode("cs30")
	cs30Name, _ := domainmeasurementitem.NewName("CS-30（30秒立ち座り）")
	cs30 := domainmeasurementitem.NewMeasurementItem(cs30ID, cs30Code, cs30Name, motorFunction, count, oneTrial, false, domainmeasurementitem.ValueTypeNumeric, &higherIsBetter, domainmeasurementitem.SideAggregationMean, []domainmeasurementitem.Element{domainmeasurementitem.ElementMuscleStrength}, createdAt, updatedAt)
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

	tests := []struct {
		name                   string
		measurementID          string
		callUsecase            bool
		evaluation             func() domainjudgment.Evaluation
		advice                 *domainjudgment.Advice
		isDraft                bool
		getErr                 error
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
			advice:                 advice,
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
			name:          "failure invalid measurement id",
			measurementID: "",
			wantCode:      codes.InvalidArgument,
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
