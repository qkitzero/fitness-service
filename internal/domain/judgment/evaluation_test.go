package judgment

import (
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/staff"
	"github.com/qkitzero/fitness-service/internal/domain/standard"
)

type wantItemEvaluation struct {
	measurementItemID measurementitem.MeasurementItemID
	value             float64
	mean              float64
	zScore            float64
	rank              standard.Rank
}

type wantElementEvaluation struct {
	element measurementitem.Element
	zScore  float64
	rank    standard.Rank
}

func TestNewEvaluation(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

	motorFunction, _ := measurementitem.NewCategory("motor_function")
	physique, _ := measurementitem.NewCategory("physique")
	kg, _ := measurementitem.NewUnit("kg")
	cm, _ := measurementitem.NewUnit("cm")
	sec, _ := measurementitem.NewUnit("sec")
	count, _ := measurementitem.NewUnit("count")
	level, _ := measurementitem.NewUnit("level")
	oneTrial, _ := measurementitem.NewTrialCount(1)
	twoTrials, _ := measurementitem.NewTrialCount(2)
	threeTrials, _ := measurementitem.NewTrialCount(3)
	higherIsBetter := measurementitem.ScoreDirectionHigherIsBetter
	lowerIsBetter := measurementitem.ScoreDirectionLowerIsBetter

	gripStrengthID := measurementitem.NewMeasurementItemID()
	gripStrengthCode, _ := measurementitem.NewCode("grip_strength")
	gripStrengthName, _ := measurementitem.NewName("握力")
	gripStrength := measurementitem.NewMeasurementItem(gripStrengthID, gripStrengthCode, gripStrengthName, motorFunction, kg, twoTrials, measurementitem.SideModeBilateral, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMuscleStrength}, createdAt, updatedAt)

	twoStepID := measurementitem.NewMeasurementItemID()
	twoStepCode, _ := measurementitem.NewCode("two_step")
	twoStepName, _ := measurementitem.NewName("2ステップ")
	twoStep := measurementitem.NewMeasurementItem(twoStepID, twoStepCode, twoStepName, motorFunction, cm, twoTrials, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMobility}, createdAt, updatedAt)

	timedUpAndGoID := measurementitem.NewMeasurementItemID()
	timedUpAndGoCode, _ := measurementitem.NewCode("timed_up_and_go")
	timedUpAndGoName, _ := measurementitem.NewName("TUG（タイム・アップ・ゴー）")
	timedUpAndGo := measurementitem.NewMeasurementItem(timedUpAndGoID, timedUpAndGoCode, timedUpAndGoName, motorFunction, sec, twoTrials, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &lowerIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMobility}, createdAt, updatedAt)

	stickReactionID := measurementitem.NewMeasurementItemID()
	stickReactionCode, _ := measurementitem.NewCode("stick_reaction")
	stickReactionName, _ := measurementitem.NewName("棒反応時間")
	stickReaction := measurementitem.NewMeasurementItem(stickReactionID, stickReactionCode, stickReactionName, motorFunction, cm, threeTrials, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &lowerIsBetter, measurementitem.TrialAggregationMean, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementAgility}, createdAt, updatedAt)

	standUpTestID := measurementitem.NewMeasurementItemID()
	standUpTestCode, _ := measurementitem.NewCode("stand_up_test")
	standUpTestName, _ := measurementitem.NewName("立ち上がり")
	standUpTest := measurementitem.NewMeasurementItem(standUpTestID, standUpTestCode, standUpTestName, motorFunction, level, oneTrial, measurementitem.SideModeOptionalBilateral, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationWorst, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMuscleStrength}, createdAt, updatedAt)

	worstSideLowerIsBetterID := measurementitem.NewMeasurementItemID()
	worstSideLowerIsBetterCode, _ := measurementitem.NewCode("worst_side_lower_is_better_item")
	worstSideLowerIsBetterName, _ := measurementitem.NewName("悪い方を採る低いほど良い測定項目")
	worstSideLowerIsBetter := measurementitem.NewMeasurementItem(worstSideLowerIsBetterID, worstSideLowerIsBetterCode, worstSideLowerIsBetterName, motorFunction, sec, twoTrials, measurementitem.SideModeBilateral, measurementitem.ValueTypeNumeric, &lowerIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationWorst, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMobility}, createdAt, updatedAt)

	unknownTrialAggregationID := measurementitem.NewMeasurementItemID()
	unknownTrialAggregationCode, _ := measurementitem.NewCode("unknown_trial_aggregation_item")
	unknownTrialAggregationName, _ := measurementitem.NewName("試行の集約の分からない測定項目")
	unknownTrialAggregation := measurementitem.NewMeasurementItem(unknownTrialAggregationID, unknownTrialAggregationCode, unknownTrialAggregationName, motorFunction, sec, twoTrials, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.TrialAggregation("median"), measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementBalance}, createdAt, updatedAt)

	unknownSideAggregationID := measurementitem.NewMeasurementItemID()
	unknownSideAggregationCode, _ := measurementitem.NewCode("unknown_side_aggregation_item")
	unknownSideAggregationName, _ := measurementitem.NewName("集約の分からない測定項目")
	unknownSideAggregation := measurementitem.NewMeasurementItem(unknownSideAggregationID, unknownSideAggregationCode, unknownSideAggregationName, motorFunction, sec, oneTrial, measurementitem.SideModeBilateral, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregation("median"), measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementBalance}, createdAt, updatedAt)

	heightRatioID := measurementitem.NewMeasurementItemID()
	heightRatioCode, _ := measurementitem.NewCode("height_ratio_item")
	heightRatioName, _ := measurementitem.NewName("身長比で評価する測定項目")
	heightRatioItem := measurementitem.NewMeasurementItem(heightRatioID, heightRatioCode, heightRatioName, motorFunction, cm, twoTrials, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationMean, measurementitem.NormalizationHeightRatio, []measurementitem.Element{measurementitem.ElementMobility}, createdAt, updatedAt)

	unknownNormalizationID := measurementitem.NewMeasurementItemID()
	unknownNormalizationCode, _ := measurementitem.NewCode("unknown_normalization_item")
	unknownNormalizationName, _ := measurementitem.NewName("正規化の分からない測定項目")
	unknownNormalization := measurementitem.NewMeasurementItem(unknownNormalizationID, unknownNormalizationCode, unknownNormalizationName, motorFunction, cm, oneTrial, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationMean, measurementitem.Normalization("weight_ratio"), []measurementitem.Element{measurementitem.ElementMobility}, createdAt, updatedAt)

	heightID := measurementitem.NewMeasurementItemID()
	heightCode, _ := measurementitem.NewCode("height")
	heightName, _ := measurementitem.NewName("身長")
	height := measurementitem.NewMeasurementItem(heightID, heightCode, heightName, physique, cm, oneTrial, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, nil, measurementitem.TrialAggregationBest, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, nil, createdAt, updatedAt)

	cs30ID := measurementitem.NewMeasurementItemID()
	cs30Code, _ := measurementitem.NewCode("cs30")
	cs30Name, _ := measurementitem.NewName("CS-30（30秒立ち座り）")
	cs30 := measurementitem.NewMeasurementItem(cs30ID, cs30Code, cs30Name, motorFunction, count, oneTrial, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMuscleEndurance}, createdAt, updatedAt)

	sitAndReachID := measurementitem.NewMeasurementItemID()
	sitAndReachCode, _ := measurementitem.NewCode("sit_and_reach")
	sitAndReachName, _ := measurementitem.NewName("長座体前屈")
	sitAndReach := measurementitem.NewMeasurementItem(sitAndReachID, sitAndReachCode, sitAndReachName, motorFunction, cm, oneTrial, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementFlexibility}, createdAt, updatedAt)

	walk5mID := measurementitem.NewMeasurementItemID()
	walk5mCode, _ := measurementitem.NewCode("walk_5m")
	walk5mName, _ := measurementitem.NewName("5m歩行")
	walk5m := measurementitem.NewMeasurementItem(walk5mID, walk5mCode, walk5mName, motorFunction, sec, oneTrial, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &lowerIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementMobility}, createdAt, updatedAt)

	seatedStepping20sID := measurementitem.NewMeasurementItemID()
	seatedStepping20sCode, _ := measurementitem.NewCode("seated_stepping_20s")
	seatedStepping20sName, _ := measurementitem.NewName("座位ステップ（20秒）")
	seatedStepping20s := measurementitem.NewMeasurementItem(seatedStepping20sID, seatedStepping20sCode, seatedStepping20sName, motorFunction, count, oneTrial, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementAgility}, createdAt, updatedAt)

	oneLegStandID := measurementitem.NewMeasurementItemID()
	oneLegStandCode, _ := measurementitem.NewCode("eyes_open_one_leg_stand")
	oneLegStandName, _ := measurementitem.NewName("開眼片足立ち")
	oneLegStand := measurementitem.NewMeasurementItem(oneLegStandID, oneLegStandCode, oneLegStandName, motorFunction, sec, twoTrials, measurementitem.SideModeBilateral, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationBest, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementBalance}, createdAt, updatedAt)

	functionalReachID := measurementitem.NewMeasurementItemID()
	functionalReachCode, _ := measurementitem.NewCode("functional_reach")
	functionalReachName, _ := measurementitem.NewName("ファンクショナルリーチ")
	functionalReach := measurementitem.NewMeasurementItem(functionalReachID, functionalReachCode, functionalReachName, motorFunction, cm, oneTrial, measurementitem.SideModeNone, measurementitem.ValueTypeNumeric, &higherIsBetter, measurementitem.TrialAggregationBest, measurementitem.SideAggregationMean, measurementitem.NormalizationNone, []measurementitem.Element{measurementitem.ElementBalance}, createdAt, updatedAt)

	items := []measurementitem.MeasurementItem{gripStrength, twoStep, timedUpAndGo, stickReaction, standUpTest, worstSideLowerIsBetter, unknownTrialAggregation, unknownSideAggregation, heightRatioItem, unknownNormalization, height, cs30, sitAndReach, walk5m, seatedStepping20s, oneLegStand, functionalReach}

	ageRange4044, _ := standard.NewAgeRange(40, 44)
	ageRange5054, _ := standard.NewAgeRange(50, 54)
	ageRange4049, _ := standard.NewAgeRange(40, 49)
	ageRange6064, _ := standard.NewAgeRange(60, 64)
	gripMean4044, _ := standard.NewMean(46)
	gripMean5054, _ := standard.NewMean(42)
	gripMean6064, _ := standard.NewMean(38)
	gripDeviation, _ := standard.NewStandardDeviation(5)
	heightRatioMean, _ := standard.NewMean(1.6)
	heightRatioDeviation, _ := standard.NewStandardDeviation(0.15)
	twoStepMean, _ := standard.NewMean(160)
	twoStepMean6064, _ := standard.NewMean(140)
	twoStepDeviation, _ := standard.NewStandardDeviation(15)
	timedUpAndGoMean, _ := standard.NewMean(6.5)
	timedUpAndGoDeviation, _ := standard.NewStandardDeviation(1)
	stickReactionMean, _ := standard.NewMean(20)
	stickReactionDeviation, _ := standard.NewStandardDeviation(4)
	cs30Mean, _ := standard.NewMean(0)
	cs30Deviation, _ := standard.NewStandardDeviation(0.01)
	sitAndReachMean4044, _ := standard.NewMean(38)
	sitAndReachMean6064, _ := standard.NewMean(30)
	sitAndReachDeviation, _ := standard.NewStandardDeviation(5)
	walk5mMean4044, _ := standard.NewMean(5.2)
	walk5mMean6064, _ := standard.NewMean(4)
	walk5mDeviation, _ := standard.NewStandardDeviation(0.5)
	seatedStepping20sMean4044, _ := standard.NewMean(40)
	seatedStepping20sMean5054, _ := standard.NewMean(35)
	seatedStepping20sDeviation, _ := standard.NewStandardDeviation(5)
	oneLegStandMean, _ := standard.NewMean(30)
	oneLegStandDeviation, _ := standard.NewStandardDeviation(10)
	ageRange4549, _ := standard.NewAgeRange(45, 49)
	functionalReachMean4044, _ := standard.NewMean(10)
	functionalReachMean4049, _ := standard.NewMean(30)
	functionalReachMean4549, _ := standard.NewMean(20)
	functionalReachDeviation, _ := standard.NewStandardDeviation(5)
	standUpTestMean, _ := standard.NewMean(6)
	standUpTestDeviation, _ := standard.NewStandardDeviation(1.5)
	worstSideLowerIsBetterMean, _ := standard.NewMean(8)
	worstSideLowerIsBetterDeviation, _ := standard.NewStandardDeviation(1)
	unknownTrialAggregationMean, _ := standard.NewMean(20)
	unknownTrialAggregationDeviation, _ := standard.NewStandardDeviation(5)
	unknownSideAggregationMean, _ := standard.NewMean(20)
	unknownSideAggregationDeviation, _ := standard.NewStandardDeviation(5)

	ageGroupStandards := []standard.AgeGroupStandard{
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), gripStrengthID, standard.GenderMale, ageRange4044, gripMean4044, gripDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), gripStrengthID, standard.GenderMale, ageRange5054, gripMean5054, gripDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), gripStrengthID, standard.GenderMale, ageRange6064, gripMean6064, gripDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), twoStepID, standard.GenderMale, ageRange4044, twoStepMean, twoStepDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), heightRatioID, standard.GenderMale, ageRange4044, heightRatioMean, heightRatioDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), unknownNormalizationID, standard.GenderMale, ageRange4044, heightRatioMean, heightRatioDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), timedUpAndGoID, standard.GenderMale, ageRange4044, timedUpAndGoMean, timedUpAndGoDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), stickReactionID, standard.GenderMale, ageRange4044, stickReactionMean, stickReactionDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), cs30ID, standard.GenderMale, ageRange4044, cs30Mean, cs30Deviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), sitAndReachID, standard.GenderMale, ageRange4044, sitAndReachMean4044, sitAndReachDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), sitAndReachID, standard.GenderMale, ageRange6064, sitAndReachMean6064, sitAndReachDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), walk5mID, standard.GenderMale, ageRange4044, walk5mMean4044, walk5mDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), walk5mID, standard.GenderMale, ageRange6064, walk5mMean6064, walk5mDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), seatedStepping20sID, standard.GenderMale, ageRange4044, seatedStepping20sMean4044, seatedStepping20sDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), seatedStepping20sID, standard.GenderMale, ageRange5054, seatedStepping20sMean5054, seatedStepping20sDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), oneLegStandID, standard.GenderMale, ageRange4044, oneLegStandMean, oneLegStandDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), oneLegStandID, standard.GenderMale, ageRange4049, oneLegStandMean, oneLegStandDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), functionalReachID, standard.GenderMale, ageRange4049, functionalReachMean4049, functionalReachDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), functionalReachID, standard.GenderMale, ageRange4044, functionalReachMean4044, functionalReachDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), functionalReachID, standard.GenderMale, ageRange4549, functionalReachMean4549, functionalReachDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), standUpTestID, standard.GenderMale, ageRange4044, standUpTestMean, standUpTestDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), worstSideLowerIsBetterID, standard.GenderMale, ageRange4044, worstSideLowerIsBetterMean, worstSideLowerIsBetterDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), unknownTrialAggregationID, standard.GenderMale, ageRange4044, unknownTrialAggregationMean, unknownTrialAggregationDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), unknownSideAggregationID, standard.GenderMale, ageRange4044, unknownSideAggregationMean, unknownSideAggregationDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), twoStepID, standard.GenderFemale, ageRange4044, twoStepMean, twoStepDeviation, createdAt, updatedAt),
		standard.NewAgeGroupStandard(standard.NewAgeGroupStandardID(), twoStepID, standard.GenderFemale, ageRange6064, twoStepMean6064, twoStepDeviation, createdAt, updatedAt),
	}

	zScoreA := standard.ZScore(1.5)
	zScoreBMin := standard.ZScore(0.5)
	zScoreBMax := standard.ZScore(1.5)
	zScoreCMin := standard.ZScore(-0.5)
	zScoreCMax := standard.ZScore(0.5)
	zScoreDMin := standard.ZScore(-1.5)
	zScoreDMax := standard.ZScore(-0.5)
	zScoreE := standard.ZScore(-1.5)
	rankStandardA, _ := standard.NewRankStandard(standard.RankA, &zScoreA, nil, createdAt, updatedAt)
	rankStandardB, _ := standard.NewRankStandard(standard.RankB, &zScoreBMin, &zScoreBMax, createdAt, updatedAt)
	rankStandardC, _ := standard.NewRankStandard(standard.RankC, &zScoreCMin, &zScoreCMax, createdAt, updatedAt)
	rankStandardD, _ := standard.NewRankStandard(standard.RankD, &zScoreDMin, &zScoreDMax, createdAt, updatedAt)
	rankStandardE, _ := standard.NewRankStandard(standard.RankE, nil, &zScoreE, createdAt, updatedAt)
	rankStandards := []standard.RankStandard{rankStandardA, rankStandardB, rankStandardC, rankStandardD, rankStandardE}
	extremeRankStandards := []standard.RankStandard{rankStandardA, rankStandardE}

	motorAge42 := 42
	motorAge45 := 45
	motorAge52 := 52
	motorAge62 := 62

	tests := []struct {
		name                   string
		gender                 standard.Gender
		age                    int
		entries                func() []measurement.MeasurementEntry
		rankStandards          func() []standard.RankStandard
		wantItemEvaluations    []wantItemEvaluation
		wantElementEvaluations []wantElementEvaluation
		wantMotorAge           *int
	}{
		{
			name:   "success a bilateral item is judged by the mean of the best value of each side",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				firstTrial, _ := measurement.NewTrialIndex(1)
				secondTrial, _ := measurement.NewTrialIndex(2)
				firstLeft, _ := measurement.NewValue(32)
				secondLeft, _ := measurement.NewValue(34)
				firstRight, _ := measurement.NewValue(36)
				secondRight, _ := measurement.NewValue(38)
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(firstTrial, measurement.SideLeft, &firstLeft, nil, nil),
					measurement.NewMeasurementValue(secondTrial, measurement.SideLeft, &secondLeft, nil, nil),
					measurement.NewMeasurementValue(firstTrial, measurement.SideRight, &firstRight, nil, nil),
					measurement.NewMeasurementValue(secondTrial, measurement.SideRight, &secondRight, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{gripStrengthID, 36, 46, -2, standard.RankE}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMuscleStrength, -2, standard.RankE}},
			wantMotorAge:           &motorAge62,
		},
		{
			name:   "success a multi trial item is judged by the best trial",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				values := make([]measurement.MeasurementValue, 0, 2)
				for i, f := range []float64{6, 7.5} {
					trialIndex, _ := measurement.NewTrialIndex(i + 1)
					value, _ := measurement.NewValue(f)
					values = append(values, measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil))
				}
				entry, _ := measurement.NewMeasurementEntry(timedUpAndGo, false, nil, values)
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{timedUpAndGoID, 6, 6.5, 0.5, standard.RankB}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMobility, 0.5, standard.RankB}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success a multi trial item is judged by the mean of the trials",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				values := make([]measurement.MeasurementValue, 0, 3)
				for i, f := range []float64{18.2, 19.5, 20.1} {
					trialIndex, _ := measurement.NewTrialIndex(i + 1)
					value, _ := measurement.NewValue(f)
					values = append(values, measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil))
				}
				entry, _ := measurement.NewMeasurementEntry(stickReaction, false, nil, values)
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{stickReactionID, 19.27, 20, 0.18, standard.RankC}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementAgility, 0.18, standard.RankC}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success a multi trial item is judged by the mean of the recorded trials",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				firstTrial, _ := measurement.NewTrialIndex(1)
				thirdTrial, _ := measurement.NewTrialIndex(3)
				first, _ := measurement.NewValue(18)
				third, _ := measurement.NewValue(20)
				entry, _ := measurement.NewMeasurementEntry(stickReaction, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(firstTrial, measurement.SideNone, &first, nil, nil),
					measurement.NewMeasurementValue(thirdTrial, measurement.SideNone, &third, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{stickReactionID, 19, 20, 0.25, standard.RankC}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementAgility, 0.25, standard.RankC}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success an item of an unknown trial aggregation is excluded",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				firstTrial, _ := measurement.NewTrialIndex(1)
				secondTrial, _ := measurement.NewTrialIndex(2)
				first, _ := measurement.NewValue(20)
				second, _ := measurement.NewValue(30)
				entry, _ := measurement.NewMeasurementEntry(unknownTrialAggregation, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(firstTrial, measurement.SideNone, &first, nil, nil),
					measurement.NewMeasurementValue(secondTrial, measurement.SideNone, &second, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
		},
		{
			name:   "success a z score on a rank boundary belongs to the upper rank",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(53.5)
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{gripStrengthID, 53.5, 46, 1.5, standard.RankA}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMuscleStrength, 1.5, standard.RankA}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success a lower is better item ranks higher when the value is smaller",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				firstTrial, _ := measurement.NewTrialIndex(1)
				secondTrial, _ := measurement.NewTrialIndex(2)
				first, _ := measurement.NewValue(6)
				second, _ := measurement.NewValue(5.5)
				entry, _ := measurement.NewMeasurementEntry(timedUpAndGo, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(firstTrial, measurement.SideNone, &first, nil, nil),
					measurement.NewMeasurementValue(secondTrial, measurement.SideNone, &second, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{timedUpAndGoID, 5.5, 6.5, 1, standard.RankB}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMobility, 1, standard.RankB}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success an element averages the z scores of the items behind it",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				firstTrial, _ := measurement.NewTrialIndex(1)
				secondTrial, _ := measurement.NewTrialIndex(2)
				firstTwoStep, _ := measurement.NewValue(160)
				secondTwoStep, _ := measurement.NewValue(155)
				twoStepEntry, _ := measurement.NewMeasurementEntry(twoStep, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(firstTrial, measurement.SideNone, &firstTwoStep, nil, nil),
					measurement.NewMeasurementValue(secondTrial, measurement.SideNone, &secondTwoStep, nil, nil),
				})
				firstTimedUpAndGo, _ := measurement.NewValue(5.5)
				secondTimedUpAndGo, _ := measurement.NewValue(6)
				timedUpAndGoEntry, _ := measurement.NewMeasurementEntry(timedUpAndGo, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(firstTrial, measurement.SideNone, &firstTimedUpAndGo, nil, nil),
					measurement.NewMeasurementValue(secondTrial, measurement.SideNone, &secondTimedUpAndGo, nil, nil),
				})
				return []measurement.MeasurementEntry{twoStepEntry, timedUpAndGoEntry}
			},
			wantItemEvaluations: []wantItemEvaluation{
				{twoStepID, 160, 160, 0, standard.RankC},
				{timedUpAndGoID, 5.5, 6.5, 1, standard.RankB},
			},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMobility, 0.5, standard.RankB}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success a z score is rounded to two decimals",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				firstTrial, _ := measurement.NewTrialIndex(1)
				secondTrial, _ := measurement.NewTrialIndex(2)
				first, _ := measurement.NewValue(170)
				second, _ := measurement.NewValue(165)
				entry, _ := measurement.NewMeasurementEntry(twoStep, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(firstTrial, measurement.SideNone, &first, nil, nil),
					measurement.NewMeasurementValue(secondTrial, measurement.SideNone, &second, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{twoStepID, 170, 160, 0.67, standard.RankB}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMobility, 0.67, standard.RankB}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success a partially entered bilateral item is judged by the entered side",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(46)
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{gripStrengthID, 46, 46, 0, standard.RankC}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMuscleStrength, 0, standard.RankC}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success an element is excluded when none of the items behind it are judged",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(46)
				gripStrengthEntry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &value, nil, nil),
				})
				note, _ := measurement.NewNote("膝の痛みのため測定不可")
				timedUpAndGoEntry, _ := measurement.NewMeasurementEntry(timedUpAndGo, true, note, nil)
				return []measurement.MeasurementEntry{gripStrengthEntry, timedUpAndGoEntry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{gripStrengthID, 46, 46, 0, standard.RankC}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMuscleStrength, 0, standard.RankC}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success the motor age is younger than the age at measurement",
			gender: standard.GenderMale,
			age:    62,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(46)
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{gripStrengthID, 46, 38, 1.6, standard.RankA}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMuscleStrength, 1.6, standard.RankA}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success the motor age is older than the age at measurement",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(38)
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{gripStrengthID, 38, 46, -1.6, standard.RankE}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMuscleStrength, -1.6, standard.RankE}},
			wantMotorAge:           &motorAge62,
		},
		{
			name:   "success a normalized item is judged by the value divided by the height",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				firstTrial, _ := measurement.NewTrialIndex(1)
				secondTrial, _ := measurement.NewTrialIndex(2)
				firstStride, _ := measurement.NewValue(245)
				secondStride, _ := measurement.NewValue(240)
				strideEntry, _ := measurement.NewMeasurementEntry(heightRatioItem, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(firstTrial, measurement.SideNone, &firstStride, nil, nil),
					measurement.NewMeasurementValue(secondTrial, measurement.SideNone, &secondStride, nil, nil),
				})
				heightValue, _ := measurement.NewValue(165)
				heightEntry, _ := measurement.NewMeasurementEntry(height, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(firstTrial, measurement.SideNone, &heightValue, nil, nil),
				})
				return []measurement.MeasurementEntry{strideEntry, heightEntry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{heightRatioID, 1.48, 1.6, -0.8, standard.RankD}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMobility, -0.8, standard.RankD}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success a normalized item is excluded without the height",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				stride, _ := measurement.NewValue(245)
				strideEntry, _ := measurement.NewMeasurementEntry(heightRatioItem, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &stride, nil, nil),
				})
				return []measurement.MeasurementEntry{strideEntry}
			},
		},
		{
			name:   "success a normalized item is excluded when the height is unmeasurable",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				stride, _ := measurement.NewValue(245)
				strideEntry, _ := measurement.NewMeasurementEntry(heightRatioItem, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &stride, nil, nil),
				})
				note, _ := measurement.NewNote("車椅子のため測定不可")
				heightEntry, _ := measurement.NewMeasurementEntry(height, true, note, nil)
				return []measurement.MeasurementEntry{strideEntry, heightEntry}
			},
		},
		{
			name:   "success a normalized item is excluded when the height has no value",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				stride, _ := measurement.NewValue(245)
				strideEntry, _ := measurement.NewMeasurementEntry(heightRatioItem, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &stride, nil, nil),
				})
				heightEntry, _ := measurement.NewMeasurementEntry(height, false, nil, nil)
				return []measurement.MeasurementEntry{strideEntry, heightEntry}
			},
		},
		{
			name:   "success a normalized item is excluded when the recorded height is empty",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				stride, _ := measurement.NewValue(245)
				strideEntry, _ := measurement.NewMeasurementEntry(heightRatioItem, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &stride, nil, nil),
				})
				heightEntry := measurement.ReconstructMeasurementEntry(heightID, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, nil, nil, nil),
				})
				return []measurement.MeasurementEntry{strideEntry, heightEntry}
			},
		},
		{
			name:   "success a normalized item is excluded when the height is zero",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				stride, _ := measurement.NewValue(245)
				strideEntry, _ := measurement.NewMeasurementEntry(heightRatioItem, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &stride, nil, nil),
				})
				heightValue, _ := measurement.NewValue(0)
				heightEntry, _ := measurement.NewMeasurementEntry(height, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &heightValue, nil, nil),
				})
				return []measurement.MeasurementEntry{strideEntry, heightEntry}
			},
		},
		{
			name:   "success a normalized item is excluded when the ratio is out of the value range",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				stride, _ := measurement.NewValue(9999.99)
				strideEntry, _ := measurement.NewMeasurementEntry(heightRatioItem, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &stride, nil, nil),
				})
				heightValue, _ := measurement.NewValue(0.01)
				heightEntry, _ := measurement.NewMeasurementEntry(height, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &heightValue, nil, nil),
				})
				return []measurement.MeasurementEntry{strideEntry, heightEntry}
			},
		},
		{
			name:   "success an item of an unknown normalization is excluded",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(245)
				entry, _ := measurement.NewMeasurementEntry(unknownNormalization, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil),
				})
				heightValue, _ := measurement.NewValue(165)
				heightEntry, _ := measurement.NewMeasurementEntry(height, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &heightValue, nil, nil),
				})
				return []measurement.MeasurementEntry{entry, heightEntry}
			},
		},
		{
			name:   "success an item without a score direction is excluded",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(170)
				heightEntry, _ := measurement.NewMeasurementEntry(height, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{heightEntry}
			},
		},
		{
			name:   "success an age just below every age group is judged by the youngest age group",
			gender: standard.GenderMale,
			age:    38,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(46)
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{gripStrengthID, 46, 46, 0, standard.RankC}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMuscleStrength, 0, standard.RankC}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success an age too far below every age group is excluded",
			gender: standard.GenderMale,
			age:    37,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(46)
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
		},
		{
			name:   "success an age just above every age group is judged by the oldest age group",
			gender: standard.GenderMale,
			age:    84,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(30)
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{gripStrengthID, 30, 38, -1.6, standard.RankE}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMuscleStrength, -1.6, standard.RankE}},
			wantMotorAge:           &motorAge62,
		},
		{
			name:   "success an age too far above every age group is excluded",
			gender: standard.GenderMale,
			age:    85,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(30)
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
		},
		{
			name:   "success an age between two age groups is excluded",
			gender: standard.GenderMale,
			age:    47,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(46)
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
		},
		{
			name:   "success the narrower age group is used when two age groups start at the same age",
			gender: standard.GenderMale,
			age:    38,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(30)
				entry, _ := measurement.NewMeasurementEntry(functionalReach, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{functionalReachID, 30, 10, 4, standard.RankA}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementBalance, 4, standard.RankA}},
			wantMotorAge:           &motorAge45,
		},
		{
			name:   "success the narrower age group is used when two age groups end at the same age",
			gender: standard.GenderMale,
			age:    52,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(30)
				entry, _ := measurement.NewMeasurementEntry(functionalReach, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{functionalReachID, 30, 20, 2, standard.RankA}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementBalance, 2, standard.RankA}},
			wantMotorAge:           &motorAge45,
		},
		{
			name:   "success an item judged by the nearest age group does not drop the motor age",
			gender: standard.GenderMale,
			age:    62,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				gripStrengthValue, _ := measurement.NewValue(38)
				gripStrengthEntry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &gripStrengthValue, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &gripStrengthValue, nil, nil),
				})
				oneLegStandValue, _ := measurement.NewValue(60)
				oneLegStandEntry, _ := measurement.NewMeasurementEntry(oneLegStand, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &oneLegStandValue, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &oneLegStandValue, nil, nil),
				})
				return []measurement.MeasurementEntry{gripStrengthEntry, oneLegStandEntry}
			},
			wantItemEvaluations: []wantItemEvaluation{
				{gripStrengthID, 38, 38, 0, standard.RankC},
				{oneLegStandID, 60, 30, 3, standard.RankA},
			},
			wantElementEvaluations: []wantElementEvaluation{
				{measurementitem.ElementMuscleStrength, 0, standard.RankC},
				{measurementitem.ElementBalance, 3, standard.RankA},
			},
			wantMotorAge: &motorAge62,
		},
		{
			name:   "success an item without any age group standard is excluded",
			gender: standard.GenderFemale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				gripStrengthValue, _ := measurement.NewValue(46)
				gripStrengthEntry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &gripStrengthValue, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &gripStrengthValue, nil, nil),
				})
				twoStepValue, _ := measurement.NewValue(190)
				twoStepEntry, _ := measurement.NewMeasurementEntry(twoStep, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &twoStepValue, nil, nil),
				})
				standUpTestValue, _ := measurement.NewValue(6)
				standUpTestEntry, _ := measurement.NewMeasurementEntry(standUpTest, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &standUpTestValue, nil, nil),
				})
				return []measurement.MeasurementEntry{gripStrengthEntry, twoStepEntry, standUpTestEntry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{twoStepID, 190, 160, 2, standard.RankA}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMobility, 2, standard.RankA}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success every item is excluded when the gender has no age group standards",
			gender: standard.Gender(""),
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(46)
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
		},
		{
			name:   "success an unmeasurable item is excluded",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				note, _ := measurement.NewNote("肩の痛みのため測定不可")
				entry, _ := measurement.NewMeasurementEntry(gripStrength, true, note, nil)
				return []measurement.MeasurementEntry{entry}
			},
		},
		{
			name:   "success an item without values is excluded",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, nil)
				return []measurement.MeasurementEntry{entry}
			},
		},
		{
			name:   "success an item without a numeric value is excluded",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				choice, _ := measurement.NewChoice("選択肢A")
				entry := measurement.ReconstructMeasurementEntry(gripStrengthID, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, nil, nil, &choice),
				})
				return []measurement.MeasurementEntry{entry}
			},
		},
		{
			name:   "success an entry of an unknown measurement item is excluded",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(46)
				entry := measurement.ReconstructMeasurementEntry(measurementitem.NewMeasurementItemID(), false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
		},
		{
			name:    "success a measurement without entries is evaluated as empty",
			gender:  standard.GenderMale,
			age:     42,
			entries: func() []measurement.MeasurementEntry { return nil },
		},
		{
			name:   "success an item whose representative value is out of range is excluded",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value := measurement.Value(10000)
				entry := measurement.ReconstructMeasurementEntry(gripStrengthID, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
		},
		{
			name:   "success an item whose z score is out of range is excluded",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(46)
				entry, _ := measurement.NewMeasurementEntry(cs30, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
		},
		{
			name:   "success an item without a matching rank standard is still counted for the motor age",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				value, _ := measurement.NewValue(46)
				entry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &value, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &value, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			rankStandards: func() []standard.RankStandard { return nil },
			wantMotorAge:  &motorAge42,
		},
		{
			name:   "success an element is excluded when no rank standard matches its average z score",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				twoStepValue, _ := measurement.NewValue(190)
				twoStepEntry, _ := measurement.NewMeasurementEntry(twoStep, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &twoStepValue, nil, nil),
				})
				timedUpAndGoValue, _ := measurement.NewValue(8.5)
				timedUpAndGoEntry, _ := measurement.NewMeasurementEntry(timedUpAndGo, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &timedUpAndGoValue, nil, nil),
				})
				return []measurement.MeasurementEntry{twoStepEntry, timedUpAndGoEntry}
			},
			rankStandards: func() []standard.RankStandard { return extremeRankStandards },
			wantItemEvaluations: []wantItemEvaluation{
				{twoStepID, 190, 160, 2, standard.RankA},
				{timedUpAndGoID, 8.5, 6.5, -2, standard.RankE},
			},
			wantMotorAge: &motorAge42,
		},
		{
			name:   "success a bilateral item aggregated by the best side ignores the weaker side",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				left, _ := measurement.NewValue(10)
				right, _ := measurement.NewValue(60)
				entry, _ := measurement.NewMeasurementEntry(oneLegStand, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &left, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &right, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{oneLegStandID, 60, 30, 3, standard.RankA}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementBalance, 3, standard.RankA}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success an item aggregated by the worst side ignores the stronger side",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				left, _ := measurement.NewValue(8)
				right, _ := measurement.NewValue(7)
				entry, _ := measurement.NewMeasurementEntry(standUpTest, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &left, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &right, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{standUpTestID, 7, 6, 0.67, standard.RankB}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMuscleStrength, 0.67, standard.RankB}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success an item aggregated by the worst side takes the value without a side over the stronger side",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				none, _ := measurement.NewValue(5)
				right, _ := measurement.NewValue(6)
				entry, _ := measurement.NewMeasurementEntry(standUpTest, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &none, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &right, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{standUpTestID, 5, 6, -0.67, standard.RankD}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMuscleStrength, -0.67, standard.RankD}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success an item aggregated by the worst side takes the best trial of each side before the sides",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				firstTrial, _ := measurement.NewTrialIndex(1)
				secondTrial, _ := measurement.NewTrialIndex(2)
				firstLeft, _ := measurement.NewValue(7)
				secondLeft, _ := measurement.NewValue(9)
				firstRight, _ := measurement.NewValue(8)
				secondRight, _ := measurement.NewValue(10)
				entry, _ := measurement.NewMeasurementEntry(worstSideLowerIsBetter, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(firstTrial, measurement.SideLeft, &firstLeft, nil, nil),
					measurement.NewMeasurementValue(secondTrial, measurement.SideLeft, &secondLeft, nil, nil),
					measurement.NewMeasurementValue(firstTrial, measurement.SideRight, &firstRight, nil, nil),
					measurement.NewMeasurementValue(secondTrial, measurement.SideRight, &secondRight, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{worstSideLowerIsBetterID, 8, 8, 0, standard.RankC}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMobility, 0, standard.RankC}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success an item of an unknown side aggregation is excluded",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				left, _ := measurement.NewValue(20)
				right, _ := measurement.NewValue(30)
				entry, _ := measurement.NewMeasurementEntry(unknownSideAggregation, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &left, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &right, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
		},
		{
			name:   "success an item scored lower is better aggregated by the worst side takes the higher side",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				left, _ := measurement.NewValue(7)
				right, _ := measurement.NewValue(9)
				entry, _ := measurement.NewMeasurementEntry(worstSideLowerIsBetter, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &left, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &right, nil, nil),
				})
				return []measurement.MeasurementEntry{entry}
			},
			wantItemEvaluations:    []wantItemEvaluation{{worstSideLowerIsBetterID, 9, 8, -1, standard.RankD}},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMobility, -1, standard.RankD}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success an exact half is rounded away from zero",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				twoStepValue, _ := measurement.NewValue(149.8)
				twoStepEntry, _ := measurement.NewMeasurementEntry(twoStep, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &twoStepValue, nil, nil),
				})
				timedUpAndGoValue, _ := measurement.NewValue(4.83)
				timedUpAndGoEntry, _ := measurement.NewMeasurementEntry(timedUpAndGo, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &timedUpAndGoValue, nil, nil),
				})
				return []measurement.MeasurementEntry{twoStepEntry, timedUpAndGoEntry}
			},
			wantItemEvaluations: []wantItemEvaluation{
				{twoStepID, 149.8, 160, -0.68, standard.RankD},
				{timedUpAndGoID, 4.83, 6.5, 1.67, standard.RankA},
			},
			wantElementEvaluations: []wantElementEvaluation{{measurementitem.ElementMobility, 0.5, standard.RankB}},
			wantMotorAge:           &motorAge42,
		},
		{
			name:   "success opposite deviations within an age group do not cancel out",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				sitAndReachValue, _ := measurement.NewValue(40)
				sitAndReachEntry, _ := measurement.NewMeasurementEntry(sitAndReach, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &sitAndReachValue, nil, nil),
				})
				walk5mValue, _ := measurement.NewValue(5)
				walk5mEntry, _ := measurement.NewMeasurementEntry(walk5m, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &walk5mValue, nil, nil),
				})
				return []measurement.MeasurementEntry{sitAndReachEntry, walk5mEntry}
			},
			wantItemEvaluations: []wantItemEvaluation{
				{sitAndReachID, 40, 38, 0.4, standard.RankC},
				{walk5mID, 5, 5.2, 0.4, standard.RankC},
			},
			wantElementEvaluations: []wantElementEvaluation{
				{measurementitem.ElementFlexibility, 0.4, standard.RankC},
				{measurementitem.ElementMobility, 0.4, standard.RankC},
			},
			wantMotorAge: &motorAge42,
		},
		{
			name:   "success an item that does not cover every age group is left out of the motor age",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				sitAndReachValue, _ := measurement.NewValue(40)
				sitAndReachEntry, _ := measurement.NewMeasurementEntry(sitAndReach, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &sitAndReachValue, nil, nil),
				})
				gripStrengthValue, _ := measurement.NewValue(38)
				gripStrengthEntry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &gripStrengthValue, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &gripStrengthValue, nil, nil),
				})
				return []measurement.MeasurementEntry{sitAndReachEntry, gripStrengthEntry}
			},
			wantItemEvaluations: []wantItemEvaluation{
				{sitAndReachID, 40, 38, 0.4, standard.RankC},
				{gripStrengthID, 38, 46, -1.6, standard.RankE},
			},
			wantElementEvaluations: []wantElementEvaluation{
				{measurementitem.ElementMuscleStrength, -1.6, standard.RankE},
				{measurementitem.ElementFlexibility, 0.4, standard.RankC},
			},
			wantMotorAge: &motorAge62,
		},
		{
			name:   "success the motor age is computed on the age groups of a single item when no item covers all of them",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				sitAndReachValue, _ := measurement.NewValue(40)
				sitAndReachEntry, _ := measurement.NewMeasurementEntry(sitAndReach, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &sitAndReachValue, nil, nil),
				})
				seatedStepping20sValue, _ := measurement.NewValue(35)
				seatedStepping20sEntry, _ := measurement.NewMeasurementEntry(seatedStepping20s, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &seatedStepping20sValue, nil, nil),
				})
				return []measurement.MeasurementEntry{sitAndReachEntry, seatedStepping20sEntry}
			},
			wantItemEvaluations: []wantItemEvaluation{
				{sitAndReachID, 40, 38, 0.4, standard.RankC},
				{seatedStepping20sID, 35, 40, -1, standard.RankD},
			},
			wantElementEvaluations: []wantElementEvaluation{
				{measurementitem.ElementFlexibility, 0.4, standard.RankC},
				{measurementitem.ElementAgility, -1, standard.RankD},
			},
			wantMotorAge: &motorAge52,
		},
		{
			name:   "success elements are ordered by the element order",
			gender: standard.GenderMale,
			age:    42,
			entries: func() []measurement.MeasurementEntry {
				trialIndex, _ := measurement.NewTrialIndex(1)
				stickReactionValue, _ := measurement.NewValue(18)
				stickReactionEntry, _ := measurement.NewMeasurementEntry(stickReaction, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideNone, &stickReactionValue, nil, nil),
				})
				gripStrengthValue, _ := measurement.NewValue(46)
				gripStrengthEntry, _ := measurement.NewMeasurementEntry(gripStrength, false, nil, []measurement.MeasurementValue{
					measurement.NewMeasurementValue(trialIndex, measurement.SideLeft, &gripStrengthValue, nil, nil),
					measurement.NewMeasurementValue(trialIndex, measurement.SideRight, &gripStrengthValue, nil, nil),
				})
				return []measurement.MeasurementEntry{stickReactionEntry, gripStrengthEntry}
			},
			wantItemEvaluations: []wantItemEvaluation{
				{stickReactionID, 18, 20, 0.5, standard.RankB},
				{gripStrengthID, 46, 46, 0, standard.RankC},
			},
			wantElementEvaluations: []wantElementEvaluation{
				{measurementitem.ElementMuscleStrength, 0, standard.RankC},
				{measurementitem.ElementAgility, 0.5, standard.RankB},
			},
			wantMotorAge: &motorAge42,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			measuredOn, _ := measurement.NewMeasuredOn(2026, 8, 1)
			measuredBy, _ := staff.NewStaffID("google-oauth2|000000000000000000000")
			ageAtMeasurement, _ := measurement.NewAgeAtMeasurement(tt.age)
			m := measurement.ReconstructMeasurement(measurement.NewMeasurementID(), customer.NewCustomerID(), measuredOn, measuredBy, ageAtMeasurement, measuredBy, true, tt.entries(), createdAt, updatedAt)

			targetRankStandards := rankStandards
			if tt.rankStandards != nil {
				targetRankStandards = tt.rankStandards()
			}

			e := NewEvaluation(m, items, tt.gender, tt.age, ageGroupStandards, targetRankStandards)

			itemEvaluations := e.ItemEvaluations()
			if len(itemEvaluations) != len(tt.wantItemEvaluations) {
				t.Fatalf("len(ItemEvaluations()) = %v, want %v", len(itemEvaluations), len(tt.wantItemEvaluations))
			}
			for i, want := range tt.wantItemEvaluations {
				got := itemEvaluations[i]
				if got.MeasurementItemID() != want.measurementItemID {
					t.Errorf("ItemEvaluations()[%d].MeasurementItemID() = %v, want %v", i, got.MeasurementItemID(), want.measurementItemID)
				}
				if got.Value().Float64() != want.value {
					t.Errorf("ItemEvaluations()[%d].Value() = %v, want %v", i, got.Value().Float64(), want.value)
				}
				if got.Mean().Float64() != want.mean {
					t.Errorf("ItemEvaluations()[%d].Mean() = %v, want %v", i, got.Mean().Float64(), want.mean)
				}
				if got.ZScore().Float64() != want.zScore {
					t.Errorf("ItemEvaluations()[%d].ZScore() = %v, want %v", i, got.ZScore().Float64(), want.zScore)
				}
				if got.Rank() != want.rank {
					t.Errorf("ItemEvaluations()[%d].Rank() = %v, want %v", i, got.Rank(), want.rank)
				}
			}

			elementEvaluations := e.ElementEvaluations()
			if len(elementEvaluations) != len(tt.wantElementEvaluations) {
				t.Fatalf("len(ElementEvaluations()) = %v, want %v", len(elementEvaluations), len(tt.wantElementEvaluations))
			}
			for i, want := range tt.wantElementEvaluations {
				got := elementEvaluations[i]
				if got.Element() != want.element {
					t.Errorf("ElementEvaluations()[%d].Element() = %v, want %v", i, got.Element(), want.element)
				}
				if got.ZScore().Float64() != want.zScore {
					t.Errorf("ElementEvaluations()[%d].ZScore() = %v, want %v", i, got.ZScore().Float64(), want.zScore)
				}
				if got.Rank() != want.rank {
					t.Errorf("ElementEvaluations()[%d].Rank() = %v, want %v", i, got.Rank(), want.rank)
				}
			}

			motorAge := e.MotorAge()
			if tt.wantMotorAge == nil {
				if motorAge != nil {
					t.Errorf("MotorAge() = %v, want nil", motorAge.Int())
				}
				return
			}
			if motorAge == nil {
				t.Fatalf("MotorAge() = nil, want %v", *tt.wantMotorAge)
			}
			if motorAge.Int() != *tt.wantMotorAge {
				t.Errorf("MotorAge() = %v, want %v", motorAge.Int(), *tt.wantMotorAge)
			}
		})
	}
}
