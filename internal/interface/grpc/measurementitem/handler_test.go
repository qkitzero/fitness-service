package measurementitem

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	measurementitemv1 "github.com/qkitzero/fitness-service/gen/go/measurementitem/v1"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	mocksappmeasurementitem "github.com/qkitzero/fitness-service/mocks/application/measurementitem"
	mocksmeasurementitem "github.com/qkitzero/fitness-service/mocks/domain/measurementitem"
)

func TestListMeasurementItems(t *testing.T) {
	t.Parallel()
	higherIsBetter := measurementitem.ScoreDirectionHigherIsBetter
	lowerIsBetter := measurementitem.ScoreDirectionLowerIsBetter
	unmappedScoreDirection := measurementitem.ScoreDirection("unknown")
	protoHigherIsBetter := measurementitemv1.ScoreDirection_SCORE_DIRECTION_HIGHER_IS_BETTER
	protoLowerIsBetter := measurementitemv1.ScoreDirection_SCORE_DIRECTION_LOWER_IS_BETTER

	tests := []struct {
		name                  string
		codes                 []string
		names                 []string
		categories            []measurementitem.Category
		units                 []measurementitem.Unit
		trialCounts           []int
		sideModes             []measurementitem.SideMode
		valueTypes            []measurementitem.ValueType
		normalizations        []measurementitem.Normalization
		scoreDirections       []*measurementitem.ScoreDirection
		trialAggregations     []measurementitem.TrialAggregation
		sideAggregations      []measurementitem.SideAggregation
		elements              [][]measurementitem.Element
		listErr               error
		wantCode              codes.Code
		wantCategories        []measurementitemv1.Category
		wantUnits             []measurementitemv1.Unit
		wantValueTypes        []measurementitemv1.ValueType
		wantSideModes         []measurementitemv1.SideMode
		wantNormalizations    []measurementitemv1.Normalization
		wantScoreDirections   []*measurementitemv1.ScoreDirection
		wantTrialAggregations []measurementitemv1.TrialAggregation
		wantSideAggregations  []measurementitemv1.SideAggregation
		wantElements          [][]measurementitemv1.ItemElement
	}{
		{
			name:                  "success list measurement items keeps every item in order",
			codes:                 []string{"blood_pressure", "height", "body_fat_percentage", "grip_strength"},
			names:                 []string{"血圧", "身長", "体脂肪率", "握力"},
			categories:            []measurementitem.Category{measurementitem.CategoryVital, measurementitem.CategoryPhysique, measurementitem.CategoryBodyComposition, measurementitem.CategoryMotorFunction},
			units:                 []measurementitem.Unit{measurementitem.UnitMmHg, measurementitem.UnitCm, measurementitem.UnitPercent, measurementitem.UnitKg},
			trialCounts:           []int{1, 1, 1, 2},
			sideModes:             []measurementitem.SideMode{measurementitem.SideModeNone, measurementitem.SideModeNone, measurementitem.SideModeNone, measurementitem.SideModeBilateral},
			valueTypes:            []measurementitem.ValueType{measurementitem.ValueTypePaired, measurementitem.ValueTypeNumeric, measurementitem.ValueTypeNumeric, measurementitem.ValueTypeNumeric},
			normalizations:        []measurementitem.Normalization{measurementitem.NormalizationNone, measurementitem.NormalizationNone, measurementitem.NormalizationNone, measurementitem.NormalizationNone},
			wantCode:              codes.OK,
			wantCategories:        []measurementitemv1.Category{measurementitemv1.Category_CATEGORY_VITAL, measurementitemv1.Category_CATEGORY_PHYSIQUE, measurementitemv1.Category_CATEGORY_BODY_COMPOSITION, measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION},
			wantUnits:             []measurementitemv1.Unit{measurementitemv1.Unit_UNIT_MMHG, measurementitemv1.Unit_UNIT_CM, measurementitemv1.Unit_UNIT_PERCENT, measurementitemv1.Unit_UNIT_KG},
			wantValueTypes:        []measurementitemv1.ValueType{measurementitemv1.ValueType_VALUE_TYPE_PAIRED, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC},
			wantSideModes:         []measurementitemv1.SideMode{measurementitemv1.SideMode_SIDE_MODE_NONE, measurementitemv1.SideMode_SIDE_MODE_NONE, measurementitemv1.SideMode_SIDE_MODE_NONE, measurementitemv1.SideMode_SIDE_MODE_BILATERAL},
			wantNormalizations:    []measurementitemv1.Normalization{measurementitemv1.Normalization_NORMALIZATION_NONE, measurementitemv1.Normalization_NORMALIZATION_NONE, measurementitemv1.Normalization_NORMALIZATION_NONE, measurementitemv1.Normalization_NORMALIZATION_NONE},
			scoreDirections:       []*measurementitem.ScoreDirection{nil, nil, nil, &higherIsBetter},
			trialAggregations:     []measurementitem.TrialAggregation{measurementitem.TrialAggregationBest, measurementitem.TrialAggregationBest, measurementitem.TrialAggregationBest, measurementitem.TrialAggregationBest},
			sideAggregations:      []measurementitem.SideAggregation{measurementitem.SideAggregationMean, measurementitem.SideAggregationMean, measurementitem.SideAggregationMean, measurementitem.SideAggregationMean},
			elements:              [][]measurementitem.Element{nil, nil, nil, {measurementitem.ElementMuscleStrength}},
			wantScoreDirections:   []*measurementitemv1.ScoreDirection{nil, nil, nil, &protoHigherIsBetter},
			wantTrialAggregations: []measurementitemv1.TrialAggregation{measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_BEST, measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_BEST, measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_BEST, measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_BEST},
			wantSideAggregations:  []measurementitemv1.SideAggregation{measurementitemv1.SideAggregation_SIDE_AGGREGATION_MEAN, measurementitemv1.SideAggregation_SIDE_AGGREGATION_MEAN, measurementitemv1.SideAggregation_SIDE_AGGREGATION_MEAN, measurementitemv1.SideAggregation_SIDE_AGGREGATION_MEAN},
			wantElements:          [][]measurementitemv1.ItemElement{{}, {}, {}, {measurementitemv1.ItemElement_ITEM_ELEMENT_MUSCLE_STRENGTH}},
		},
		{
			name:                  "success maps the remaining units and value types",
			codes:                 []string{"pulse_rate", "eyes_closed_one_leg_stand", "cs30", "stand_up_test", "choice_item"},
			names:                 []string{"脈拍", "閉眼片足立ち", "CS-30（30秒立ち座り）", "立ち上がり", "選択式項目"},
			categories:            []measurementitem.Category{measurementitem.CategoryVital, measurementitem.CategoryMotorFunction, measurementitem.CategoryMotorFunction, measurementitem.CategoryMotorFunction, measurementitem.CategoryMotorFunction},
			units:                 []measurementitem.Unit{measurementitem.UnitBpm, measurementitem.UnitSec, measurementitem.UnitCount, measurementitem.UnitLevel, measurementitem.UnitCm},
			trialCounts:           []int{1, 2, 1, 1, 1},
			sideModes:             []measurementitem.SideMode{measurementitem.SideModeNone, measurementitem.SideModeBilateral, measurementitem.SideModeNone, measurementitem.SideModeOptionalBilateral, measurementitem.SideModeNone},
			valueTypes:            []measurementitem.ValueType{measurementitem.ValueTypeNumeric, measurementitem.ValueTypeNumeric, measurementitem.ValueTypeNumeric, measurementitem.ValueTypeNumeric, measurementitem.ValueTypeChoice},
			normalizations:        []measurementitem.Normalization{measurementitem.NormalizationNone, measurementitem.NormalizationNone, measurementitem.NormalizationNone, measurementitem.NormalizationNone, measurementitem.NormalizationNone},
			wantCode:              codes.OK,
			wantCategories:        []measurementitemv1.Category{measurementitemv1.Category_CATEGORY_VITAL, measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION, measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION, measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION, measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION},
			wantUnits:             []measurementitemv1.Unit{measurementitemv1.Unit_UNIT_BPM, measurementitemv1.Unit_UNIT_SEC, measurementitemv1.Unit_UNIT_COUNT, measurementitemv1.Unit_UNIT_LEVEL, measurementitemv1.Unit_UNIT_CM},
			wantValueTypes:        []measurementitemv1.ValueType{measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_CHOICE},
			wantSideModes:         []measurementitemv1.SideMode{measurementitemv1.SideMode_SIDE_MODE_NONE, measurementitemv1.SideMode_SIDE_MODE_BILATERAL, measurementitemv1.SideMode_SIDE_MODE_NONE, measurementitemv1.SideMode_SIDE_MODE_OPTIONAL_BILATERAL, measurementitemv1.SideMode_SIDE_MODE_NONE},
			wantNormalizations:    []measurementitemv1.Normalization{measurementitemv1.Normalization_NORMALIZATION_NONE, measurementitemv1.Normalization_NORMALIZATION_NONE, measurementitemv1.Normalization_NORMALIZATION_NONE, measurementitemv1.Normalization_NORMALIZATION_NONE, measurementitemv1.Normalization_NORMALIZATION_NONE},
			scoreDirections:       []*measurementitem.ScoreDirection{nil, &higherIsBetter, &higherIsBetter, &higherIsBetter, nil},
			trialAggregations:     []measurementitem.TrialAggregation{measurementitem.TrialAggregationBest, measurementitem.TrialAggregationBest, measurementitem.TrialAggregationBest, measurementitem.TrialAggregationBest, measurementitem.TrialAggregationBest},
			sideAggregations:      []measurementitem.SideAggregation{measurementitem.SideAggregationMean, measurementitem.SideAggregationBest, measurementitem.SideAggregationMean, measurementitem.SideAggregationWorst, measurementitem.SideAggregationMean},
			elements:              [][]measurementitem.Element{nil, {measurementitem.ElementBalance}, {measurementitem.ElementMuscleEndurance}, {measurementitem.ElementMuscleStrength}, nil},
			wantScoreDirections:   []*measurementitemv1.ScoreDirection{nil, &protoHigherIsBetter, &protoHigherIsBetter, &protoHigherIsBetter, nil},
			wantTrialAggregations: []measurementitemv1.TrialAggregation{measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_BEST, measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_BEST, measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_BEST, measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_BEST, measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_BEST},
			wantSideAggregations:  []measurementitemv1.SideAggregation{measurementitemv1.SideAggregation_SIDE_AGGREGATION_MEAN, measurementitemv1.SideAggregation_SIDE_AGGREGATION_BEST, measurementitemv1.SideAggregation_SIDE_AGGREGATION_MEAN, measurementitemv1.SideAggregation_SIDE_AGGREGATION_WORST, measurementitemv1.SideAggregation_SIDE_AGGREGATION_MEAN},
			wantElements:          [][]measurementitemv1.ItemElement{{}, {measurementitemv1.ItemElement_ITEM_ELEMENT_BALANCE}, {measurementitemv1.ItemElement_ITEM_ELEMENT_MUSCLE_ENDURANCE}, {measurementitemv1.ItemElement_ITEM_ELEMENT_MUSCLE_STRENGTH}, {}},
		},
		{
			name:                  "success maps the height ratio normalization",
			codes:                 []string{"two_step"},
			names:                 []string{"2ステップ"},
			categories:            []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:                 []measurementitem.Unit{measurementitem.UnitCm},
			trialCounts:           []int{2},
			sideModes:             []measurementitem.SideMode{measurementitem.SideModeNone},
			valueTypes:            []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			normalizations:        []measurementitem.Normalization{measurementitem.NormalizationHeightRatio},
			wantCode:              codes.OK,
			wantCategories:        []measurementitemv1.Category{measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION},
			wantUnits:             []measurementitemv1.Unit{measurementitemv1.Unit_UNIT_CM},
			wantValueTypes:        []measurementitemv1.ValueType{measurementitemv1.ValueType_VALUE_TYPE_NUMERIC},
			wantSideModes:         []measurementitemv1.SideMode{measurementitemv1.SideMode_SIDE_MODE_NONE},
			wantNormalizations:    []measurementitemv1.Normalization{measurementitemv1.Normalization_NORMALIZATION_HEIGHT_RATIO},
			scoreDirections:       []*measurementitem.ScoreDirection{&higherIsBetter},
			trialAggregations:     []measurementitem.TrialAggregation{measurementitem.TrialAggregationBest},
			sideAggregations:      []measurementitem.SideAggregation{measurementitem.SideAggregationMean},
			elements:              [][]measurementitem.Element{{measurementitem.ElementMobility}},
			wantScoreDirections:   []*measurementitemv1.ScoreDirection{&protoHigherIsBetter},
			wantTrialAggregations: []measurementitemv1.TrialAggregation{measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_BEST},
			wantSideAggregations:  []measurementitemv1.SideAggregation{measurementitemv1.SideAggregation_SIDE_AGGREGATION_MEAN},
			wantElements:          [][]measurementitemv1.ItemElement{{measurementitemv1.ItemElement_ITEM_ELEMENT_MOBILITY}},
		},
		{
			name:                  "success maps the lower is better, the mean of the trials and the remaining elements",
			codes:                 []string{"stick_reaction", "timed_up_and_go", "sit_and_reach"},
			names:                 []string{"棒反応時間", "TUG（タイム・アップ・ゴー）", "長座体前屈"},
			categories:            []measurementitem.Category{measurementitem.CategoryMotorFunction, measurementitem.CategoryMotorFunction, measurementitem.CategoryMotorFunction},
			units:                 []measurementitem.Unit{measurementitem.UnitCm, measurementitem.UnitSec, measurementitem.UnitCm},
			trialCounts:           []int{3, 2, 2},
			sideModes:             []measurementitem.SideMode{measurementitem.SideModeNone, measurementitem.SideModeNone, measurementitem.SideModeNone},
			valueTypes:            []measurementitem.ValueType{measurementitem.ValueTypeNumeric, measurementitem.ValueTypeNumeric, measurementitem.ValueTypeNumeric},
			normalizations:        []measurementitem.Normalization{measurementitem.NormalizationNone, measurementitem.NormalizationNone, measurementitem.NormalizationNone},
			scoreDirections:       []*measurementitem.ScoreDirection{&lowerIsBetter, &lowerIsBetter, &higherIsBetter},
			trialAggregations:     []measurementitem.TrialAggregation{measurementitem.TrialAggregationMean, measurementitem.TrialAggregationBest, measurementitem.TrialAggregationBest},
			sideAggregations:      []measurementitem.SideAggregation{measurementitem.SideAggregationMean, measurementitem.SideAggregationMean, measurementitem.SideAggregationMean},
			elements:              [][]measurementitem.Element{{measurementitem.ElementAgility}, {measurementitem.ElementMobility, measurementitem.ElementBalance}, {measurementitem.ElementFlexibility}},
			wantCode:              codes.OK,
			wantCategories:        []measurementitemv1.Category{measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION, measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION, measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION},
			wantUnits:             []measurementitemv1.Unit{measurementitemv1.Unit_UNIT_CM, measurementitemv1.Unit_UNIT_SEC, measurementitemv1.Unit_UNIT_CM},
			wantValueTypes:        []measurementitemv1.ValueType{measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC},
			wantSideModes:         []measurementitemv1.SideMode{measurementitemv1.SideMode_SIDE_MODE_NONE, measurementitemv1.SideMode_SIDE_MODE_NONE, measurementitemv1.SideMode_SIDE_MODE_NONE},
			wantNormalizations:    []measurementitemv1.Normalization{measurementitemv1.Normalization_NORMALIZATION_NONE, measurementitemv1.Normalization_NORMALIZATION_NONE, measurementitemv1.Normalization_NORMALIZATION_NONE},
			wantScoreDirections:   []*measurementitemv1.ScoreDirection{&protoLowerIsBetter, &protoLowerIsBetter, &protoHigherIsBetter},
			wantTrialAggregations: []measurementitemv1.TrialAggregation{measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_MEAN, measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_BEST, measurementitemv1.TrialAggregation_TRIAL_AGGREGATION_BEST},
			wantSideAggregations:  []measurementitemv1.SideAggregation{measurementitemv1.SideAggregation_SIDE_AGGREGATION_MEAN, measurementitemv1.SideAggregation_SIDE_AGGREGATION_MEAN, measurementitemv1.SideAggregation_SIDE_AGGREGATION_MEAN},
			wantElements:          [][]measurementitemv1.ItemElement{{measurementitemv1.ItemElement_ITEM_ELEMENT_AGILITY}, {measurementitemv1.ItemElement_ITEM_ELEMENT_MOBILITY, measurementitemv1.ItemElement_ITEM_ELEMENT_BALANCE}, {measurementitemv1.ItemElement_ITEM_ELEMENT_FLEXIBILITY}},
		},
		{
			name:     "success list no measurement items",
			wantCode: codes.OK,
		},
		{
			name:     "failure usecase error",
			listErr:  fmt.Errorf("list measurement items error"),
			wantCode: codes.Internal,
		},
		{
			name:     "failure unauthenticated is preserved",
			listErr:  status.Error(codes.Unauthenticated, "auth"),
			wantCode: codes.Unauthenticated,
		},
		{
			name:     "failure permission denied is preserved",
			listErr:  status.Error(codes.PermissionDenied, "forbidden"),
			wantCode: codes.PermissionDenied,
		},
		{
			name:     "failure downstream code is not forwarded",
			listErr:  status.Error(codes.NotFound, "user not found"),
			wantCode: codes.Internal,
		},
		{
			name:           "failure unmapped category",
			codes:          []string{"grip_strength"},
			names:          []string{"握力"},
			categories:     []measurementitem.Category{measurementitem.Category("flexibility")},
			units:          []measurementitem.Unit{measurementitem.UnitKg},
			trialCounts:    []int{2},
			sideModes:      []measurementitem.SideMode{measurementitem.SideModeBilateral},
			valueTypes:     []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			normalizations: []measurementitem.Normalization{measurementitem.NormalizationNone},
			wantCode:       codes.Internal,
		},
		{
			name:           "failure unmapped unit",
			codes:          []string{"grip_strength"},
			names:          []string{"握力"},
			categories:     []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:          []measurementitem.Unit{measurementitem.Unit("newton")},
			trialCounts:    []int{2},
			sideModes:      []measurementitem.SideMode{measurementitem.SideModeBilateral},
			valueTypes:     []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			normalizations: []measurementitem.Normalization{measurementitem.NormalizationNone},
			wantCode:       codes.Internal,
		},
		{
			name:           "failure unmapped side mode",
			codes:          []string{"grip_strength"},
			names:          []string{"握力"},
			categories:     []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:          []measurementitem.Unit{measurementitem.UnitKg},
			trialCounts:    []int{2},
			sideModes:      []measurementitem.SideMode{measurementitem.SideMode("unilateral")},
			valueTypes:     []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			normalizations: []measurementitem.Normalization{measurementitem.NormalizationNone},
			wantCode:       codes.Internal,
		},
		{
			name:           "failure unmapped value type",
			codes:          []string{"grip_strength"},
			names:          []string{"握力"},
			categories:     []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:          []measurementitem.Unit{measurementitem.UnitKg},
			trialCounts:    []int{2},
			sideModes:      []measurementitem.SideMode{measurementitem.SideModeBilateral},
			valueTypes:     []measurementitem.ValueType{measurementitem.ValueType("range")},
			normalizations: []measurementitem.Normalization{measurementitem.NormalizationNone},
			wantCode:       codes.Internal,
		},
		{
			name:           "failure unmapped normalization",
			codes:          []string{"two_step"},
			names:          []string{"2ステップ"},
			categories:     []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:          []measurementitem.Unit{measurementitem.UnitCm},
			trialCounts:    []int{2},
			sideModes:      []measurementitem.SideMode{measurementitem.SideModeNone},
			valueTypes:     []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			normalizations: []measurementitem.Normalization{measurementitem.Normalization("weight_ratio")},
			wantCode:       codes.Internal,
		},
		{
			name:              "failure unmapped score direction",
			codes:             []string{"grip_strength"},
			names:             []string{"握力"},
			categories:        []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:             []measurementitem.Unit{measurementitem.UnitKg},
			trialCounts:       []int{2},
			sideModes:         []measurementitem.SideMode{measurementitem.SideModeBilateral},
			valueTypes:        []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			normalizations:    []measurementitem.Normalization{measurementitem.NormalizationNone},
			scoreDirections:   []*measurementitem.ScoreDirection{&unmappedScoreDirection},
			trialAggregations: []measurementitem.TrialAggregation{measurementitem.TrialAggregationBest},
			sideAggregations:  []measurementitem.SideAggregation{measurementitem.SideAggregationMean},
			elements:          [][]measurementitem.Element{{measurementitem.ElementMuscleStrength}},
			wantCode:          codes.Internal,
		},
		{
			name:              "failure unmapped trial aggregation",
			codes:             []string{"grip_strength"},
			names:             []string{"握力"},
			categories:        []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:             []measurementitem.Unit{measurementitem.UnitKg},
			trialCounts:       []int{2},
			sideModes:         []measurementitem.SideMode{measurementitem.SideModeBilateral},
			valueTypes:        []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			normalizations:    []measurementitem.Normalization{measurementitem.NormalizationNone},
			scoreDirections:   []*measurementitem.ScoreDirection{&higherIsBetter},
			trialAggregations: []measurementitem.TrialAggregation{measurementitem.TrialAggregation("median")},
			sideAggregations:  []measurementitem.SideAggregation{measurementitem.SideAggregationMean},
			elements:          [][]measurementitem.Element{{measurementitem.ElementMuscleStrength}},
			wantCode:          codes.Internal,
		},
		{
			name:              "failure unmapped side aggregation",
			codes:             []string{"grip_strength"},
			names:             []string{"握力"},
			categories:        []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:             []measurementitem.Unit{measurementitem.UnitKg},
			trialCounts:       []int{2},
			sideModes:         []measurementitem.SideMode{measurementitem.SideModeBilateral},
			valueTypes:        []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			normalizations:    []measurementitem.Normalization{measurementitem.NormalizationNone},
			scoreDirections:   []*measurementitem.ScoreDirection{&higherIsBetter},
			trialAggregations: []measurementitem.TrialAggregation{measurementitem.TrialAggregationBest},
			sideAggregations:  []measurementitem.SideAggregation{measurementitem.SideAggregation("median")},
			elements:          [][]measurementitem.Element{{measurementitem.ElementMuscleStrength}},
			wantCode:          codes.Internal,
		},
		{
			name:              "failure unmapped element",
			codes:             []string{"grip_strength"},
			names:             []string{"握力"},
			categories:        []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:             []measurementitem.Unit{measurementitem.UnitKg},
			trialCounts:       []int{2},
			sideModes:         []measurementitem.SideMode{measurementitem.SideModeBilateral},
			valueTypes:        []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			normalizations:    []measurementitem.Normalization{measurementitem.NormalizationNone},
			scoreDirections:   []*measurementitem.ScoreDirection{&higherIsBetter},
			trialAggregations: []measurementitem.TrialAggregation{measurementitem.TrialAggregationBest},
			sideAggregations:  []measurementitem.SideAggregation{measurementitem.SideAggregationMean},
			elements:          [][]measurementitem.Element{{measurementitem.ElementMuscleStrength, measurementitem.Element("unknown")}},
			wantCode:          codes.Internal,
		},
		{
			name:           "failure trial count out of proto range",
			codes:          []string{"grip_strength"},
			names:          []string{"握力"},
			categories:     []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:          []measurementitem.Unit{measurementitem.UnitKg},
			trialCounts:    []int{-1},
			sideModes:      []measurementitem.SideMode{measurementitem.SideModeBilateral},
			valueTypes:     []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			normalizations: []measurementitem.Normalization{measurementitem.NormalizationNone},
			wantCode:       codes.Internal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			wantIDs := make([]measurementitem.MeasurementItemID, 0, len(tt.codes))
			measurementItems := make([]measurementitem.MeasurementItem, 0, len(tt.codes))
			for i := range tt.codes {
				id := measurementitem.NewMeasurementItemID()
				wantIDs = append(wantIDs, id)
				mockMeasurementItem := mocksmeasurementitem.NewMockMeasurementItem(ctrl)
				mockMeasurementItem.EXPECT().ID().Return(id).AnyTimes()
				mockMeasurementItem.EXPECT().Code().Return(measurementitem.Code(tt.codes[i])).AnyTimes()
				mockMeasurementItem.EXPECT().Name().Return(measurementitem.Name(tt.names[i])).AnyTimes()
				mockMeasurementItem.EXPECT().Category().Return(tt.categories[i]).AnyTimes()
				mockMeasurementItem.EXPECT().Unit().Return(tt.units[i]).AnyTimes()
				mockMeasurementItem.EXPECT().TrialCount().Return(measurementitem.TrialCount(tt.trialCounts[i])).AnyTimes()
				mockMeasurementItem.EXPECT().SideMode().Return(tt.sideModes[i]).AnyTimes()
				mockMeasurementItem.EXPECT().ValueType().Return(tt.valueTypes[i]).AnyTimes()
				mockMeasurementItem.EXPECT().Normalization().Return(tt.normalizations[i]).AnyTimes()
				var scoreDirection *measurementitem.ScoreDirection
				if tt.scoreDirections != nil {
					scoreDirection = tt.scoreDirections[i]
				}
				mockMeasurementItem.EXPECT().ScoreDirection().Return(scoreDirection).AnyTimes()
				trialAggregation := measurementitem.TrialAggregationBest
				if tt.trialAggregations != nil {
					trialAggregation = tt.trialAggregations[i]
				}
				mockMeasurementItem.EXPECT().TrialAggregation().Return(trialAggregation).AnyTimes()
				sideAggregation := measurementitem.SideAggregationMean
				if tt.sideAggregations != nil {
					sideAggregation = tt.sideAggregations[i]
				}
				mockMeasurementItem.EXPECT().SideAggregation().Return(sideAggregation).AnyTimes()
				var elements []measurementitem.Element
				if tt.elements != nil {
					elements = tt.elements[i]
				}
				mockMeasurementItem.EXPECT().Elements().Return(elements).AnyTimes()
				measurementItems = append(measurementItems, mockMeasurementItem)
			}

			mockUsecase := mocksappmeasurementitem.NewMockMeasurementItemUsecase(ctrl)
			if tt.listErr != nil {
				mockUsecase.EXPECT().ListMeasurementItems(gomock.Any()).Return(nil, tt.listErr).Times(1)
			} else {
				mockUsecase.EXPECT().ListMeasurementItems(gomock.Any()).Return(measurementItems, nil).Times(1)
			}

			handler := NewMeasurementItemHandler(mockUsecase)

			res, err := handler.ListMeasurementItems(ctx, &measurementitemv1.ListMeasurementItemsRequest{})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode != codes.OK {
				return
			}
			if len(res.GetMeasurementItems()) != len(tt.codes) {
				t.Fatalf("len(MeasurementItems) = %v, want %v", len(res.GetMeasurementItems()), len(tt.codes))
			}
			for i, msg := range res.GetMeasurementItems() {
				if msg.GetMeasurementItemId() != wantIDs[i].String() {
					t.Errorf("MeasurementItems[%d].MeasurementItemId = %v, want %v", i, msg.GetMeasurementItemId(), wantIDs[i].String())
				}
				if msg.GetCode() != tt.codes[i] {
					t.Errorf("MeasurementItems[%d].Code = %v, want %v", i, msg.GetCode(), tt.codes[i])
				}
				if msg.GetName() != tt.names[i] {
					t.Errorf("MeasurementItems[%d].Name = %v, want %v", i, msg.GetName(), tt.names[i])
				}
				if msg.GetCategory() != tt.wantCategories[i] {
					t.Errorf("MeasurementItems[%d].Category = %v, want %v", i, msg.GetCategory(), tt.wantCategories[i])
				}
				if msg.GetUnit() != tt.wantUnits[i] {
					t.Errorf("MeasurementItems[%d].Unit = %v, want %v", i, msg.GetUnit(), tt.wantUnits[i])
				}
				if msg.GetTrialCount() != uint32(tt.trialCounts[i]) {
					t.Errorf("MeasurementItems[%d].TrialCount = %v, want %v", i, msg.GetTrialCount(), tt.trialCounts[i])
				}
				if msg.GetSideMode() != tt.wantSideModes[i] {
					t.Errorf("MeasurementItems[%d].SideMode = %v, want %v", i, msg.GetSideMode(), tt.wantSideModes[i])
				}
				if msg.GetValueType() != tt.wantValueTypes[i] {
					t.Errorf("MeasurementItems[%d].ValueType = %v, want %v", i, msg.GetValueType(), tt.wantValueTypes[i])
				}
				if msg.GetNormalization() != tt.wantNormalizations[i] {
					t.Errorf("MeasurementItems[%d].Normalization = %v, want %v", i, msg.GetNormalization(), tt.wantNormalizations[i])
				}
				switch want := tt.wantScoreDirections[i]; {
				case want == nil && msg.ScoreDirection != nil:
					t.Errorf("MeasurementItems[%d].ScoreDirection = %v, want nil", i, msg.GetScoreDirection())
				case want != nil && msg.ScoreDirection == nil:
					t.Errorf("MeasurementItems[%d].ScoreDirection = nil, want %v", i, *want)
				case want != nil && msg.GetScoreDirection() != *want:
					t.Errorf("MeasurementItems[%d].ScoreDirection = %v, want %v", i, msg.GetScoreDirection(), *want)
				}
				if msg.GetTrialAggregation() != tt.wantTrialAggregations[i] {
					t.Errorf("MeasurementItems[%d].TrialAggregation = %v, want %v", i, msg.GetTrialAggregation(), tt.wantTrialAggregations[i])
				}
				if msg.GetSideAggregation() != tt.wantSideAggregations[i] {
					t.Errorf("MeasurementItems[%d].SideAggregation = %v, want %v", i, msg.GetSideAggregation(), tt.wantSideAggregations[i])
				}
				if !slices.Equal(msg.GetElements(), tt.wantElements[i]) {
					t.Errorf("MeasurementItems[%d].Elements = %v, want %v", i, msg.GetElements(), tt.wantElements[i])
				}
			}
		})
	}
}
