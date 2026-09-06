package judgment

import (
	"math"
	"sort"

	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/standard"
)

const (
	evaluationScale                 = 100
	maxYoungerAgeGroupFallbackYears = 2
	maxOlderAgeGroupFallbackYears   = 20
)

type Evaluation interface {
	ItemEvaluations() []ItemEvaluation
	ElementEvaluations() []ElementEvaluation
	MotorAge() *MotorAge
}

type evaluation struct {
	itemEvaluations    []ItemEvaluation
	elementEvaluations []ElementEvaluation
	motorAge           *MotorAge
}

func (e evaluation) ItemEvaluations() []ItemEvaluation {
	itemEvaluations := make([]ItemEvaluation, len(e.itemEvaluations))
	copy(itemEvaluations, e.itemEvaluations)
	return itemEvaluations
}

func (e evaluation) ElementEvaluations() []ElementEvaluation {
	elementEvaluations := make([]ElementEvaluation, len(e.elementEvaluations))
	copy(elementEvaluations, e.elementEvaluations)
	return elementEvaluations
}

func (e evaluation) MotorAge() *MotorAge {
	if e.motorAge == nil {
		return nil
	}
	m := *e.motorAge
	return &m
}

type judgedItem struct {
	item  measurementitem.MeasurementItem
	value measurement.Value
}

func toHundredths(f float64) int64 {
	return int64(math.Round(f * evaluationScale))
}

func fromHundredths(hundredths int64) float64 {
	return float64(hundredths) / evaluationScale
}

func divideRounded(numerator, denominator int64) int64 {
	if numerator < 0 {
		return -((-2*numerator + denominator) / (2 * denominator))
	}
	return (2*numerator + denominator) / (2 * denominator)
}

func isBetter(candidate, current int64, scoreDirection measurementitem.ScoreDirection) bool {
	if scoreDirection == measurementitem.ScoreDirectionLowerIsBetter {
		return candidate < current
	}
	return candidate > current
}

func aggregateTrials(trials []int64, trialAggregation measurementitem.TrialAggregation, scoreDirection measurementitem.ScoreDirection) (int64, bool) {
	switch trialAggregation {
	case measurementitem.TrialAggregationBest:
		best := trials[0]
		for _, trial := range trials[1:] {
			if isBetter(trial, best, scoreDirection) {
				best = trial
			}
		}
		return best, true
	case measurementitem.TrialAggregationMean:
		sum := int64(0)
		for _, trial := range trials {
			sum += trial
		}
		return divideRounded(sum, int64(len(trials))), true
	default:
		return 0, false
	}
}

func representativeValue(entry measurement.MeasurementEntry, item measurementitem.MeasurementItem, scoreDirection measurementitem.ScoreDirection) (measurement.Value, bool) {
	trialsBySide := make(map[measurement.Side][]int64)
	sides := make([]measurement.Side, 0, 3)
	for _, measurementValue := range entry.Values() {
		value := measurementValue.Value()
		if value == nil {
			continue
		}
		side := measurementValue.Side()
		if _, ok := trialsBySide[side]; !ok {
			sides = append(sides, side)
		}
		trialsBySide[side] = append(trialsBySide[side], toHundredths(value.Float64()))
	}
	if len(sides) == 0 {
		return measurement.Value(0), false
	}

	sideValues := make([]int64, 0, len(sides))
	for _, side := range sides {
		aggregated, ok := aggregateTrials(trialsBySide[side], item.TrialAggregation(), scoreDirection)
		if !ok {
			return measurement.Value(0), false
		}
		sideValues = append(sideValues, aggregated)
	}

	representative := sideValues[0]
	switch item.SideAggregation() {
	case measurementitem.SideAggregationBest:
		for _, sideValue := range sideValues[1:] {
			if isBetter(sideValue, representative, scoreDirection) {
				representative = sideValue
			}
		}
	case measurementitem.SideAggregationWorst:
		for _, sideValue := range sideValues[1:] {
			if isBetter(representative, sideValue, scoreDirection) {
				representative = sideValue
			}
		}
	case measurementitem.SideAggregationMean:
		sum := int64(0)
		for _, sideValue := range sideValues {
			sum += sideValue
		}
		representative = divideRounded(sum, int64(len(sideValues)))
	default:
		return measurement.Value(0), false
	}

	value, err := measurement.NewValue(fromHundredths(representative))
	if err != nil {
		return measurement.Value(0), false
	}

	return value, true
}

func heightHundredths(entries []measurement.MeasurementEntry, itemByID map[measurementitem.MeasurementItemID]measurementitem.MeasurementItem) (int64, bool) {
	for _, entry := range entries {
		if entry.Unmeasurable() {
			continue
		}
		item, ok := itemByID[entry.MeasurementItemID()]
		if !ok || item.Code() != measurementitem.CodeHeight {
			continue
		}
		sum := int64(0)
		count := int64(0)
		for _, measurementValue := range entry.Values() {
			value := measurementValue.Value()
			if value == nil {
				continue
			}
			sum += toHundredths(value.Float64())
			count++
		}
		if count == 0 {
			return 0, false
		}
		return divideRounded(sum, count), true
	}
	return 0, false
}

func normalizedValue(value measurement.Value, normalization measurementitem.Normalization, baseHundredths int64, hasBase bool) (measurement.Value, bool) {
	switch normalization {
	case measurementitem.NormalizationNone:
		return value, true
	case measurementitem.NormalizationHeightRatio:
		if !hasBase || baseHundredths <= 0 {
			return measurement.Value(0), false
		}
		normalized, err := measurement.NewValue(fromHundredths(divideRounded(evaluationScale*toHundredths(value.Float64()), baseHundredths)))
		if err != nil {
			return measurement.Value(0), false
		}
		return normalized, true
	default:
		return measurement.Value(0), false
	}
}

func isYoungerAgeRange(a, b standard.AgeRange) bool {
	if a.From() != b.From() {
		return a.From() < b.From()
	}
	return a.To() < b.To()
}

func isOlderAgeRange(a, b standard.AgeRange) bool {
	if a.To() != b.To() {
		return a.To() > b.To()
	}
	return a.From() > b.From()
}

func findAgeGroupStandard(ageGroupStandards []standard.AgeGroupStandard, age int) (standard.AgeGroupStandard, bool) {
	var youngest, oldest standard.AgeGroupStandard
	for _, ageGroupStandard := range ageGroupStandards {
		ageRange := ageGroupStandard.AgeRange()
		if ageRange.Contains(age) {
			return ageGroupStandard, true
		}
		if youngest == nil || isYoungerAgeRange(ageRange, youngest.AgeRange()) {
			youngest = ageGroupStandard
		}
		if oldest == nil || isOlderAgeRange(ageRange, oldest.AgeRange()) {
			oldest = ageGroupStandard
		}
	}
	if youngest == nil {
		return nil, false
	}
	if age < youngest.AgeRange().From() {
		if youngest.AgeRange().Distance(age) > maxYoungerAgeGroupFallbackYears {
			return nil, false
		}
		return youngest, true
	}
	if age > oldest.AgeRange().To() {
		if oldest.AgeRange().Distance(age) > maxOlderAgeGroupFallbackYears {
			return nil, false
		}
		return oldest, true
	}
	return nil, false
}

func zScoreHundredthsOf(value measurement.Value, meanHundredths, standardDeviationHundredths int64, scoreDirection measurementitem.ScoreDirection) int64 {
	deviation := toHundredths(value.Float64()) - meanHundredths
	zScore := divideRounded(evaluationScale*deviation, standardDeviationHundredths)
	if scoreDirection == measurementitem.ScoreDirectionLowerIsBetter {
		return -zScore
	}
	return zScore
}

func zScoreHundredths(value measurement.Value, ageGroupStandard standard.AgeGroupStandard, scoreDirection measurementitem.ScoreDirection) int64 {
	return zScoreHundredthsOf(
		value,
		toHundredths(ageGroupStandard.Mean().Float64()),
		toHundredths(ageGroupStandard.StandardDeviation().Float64()),
		scoreDirection,
	)
}

func findRank(rankStandards []standard.RankStandard, zScore standard.ZScore) (standard.Rank, bool) {
	for _, rankStandard := range rankStandards {
		if zScoreMin := rankStandard.ZScoreMin(); zScoreMin != nil && zScore < *zScoreMin {
			continue
		}
		if zScoreMax := rankStandard.ZScoreMax(); zScoreMax != nil && zScore >= *zScoreMax {
			continue
		}
		return rankStandard.Rank(), true
	}
	return standard.Rank(""), false
}

func newElementEvaluations(
	itemEvaluations []ItemEvaluation,
	itemByID map[measurementitem.MeasurementItemID]measurementitem.MeasurementItem,
	rankStandards []standard.RankStandard,
) []ElementEvaluation {
	zScoresByElement := make(map[measurementitem.Element][]standard.ZScore)
	elements := make([]measurementitem.Element, 0, len(itemEvaluations))
	for _, itemEvaluation := range itemEvaluations {
		for _, element := range itemByID[itemEvaluation.MeasurementItemID()].Elements() {
			if _, ok := zScoresByElement[element]; !ok {
				elements = append(elements, element)
			}
			zScoresByElement[element] = append(zScoresByElement[element], itemEvaluation.ZScore())
		}
	}

	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Order() < elements[j].Order()
	})

	elementEvaluations := make([]ElementEvaluation, 0, len(elements))
	for _, element := range elements {
		zScore := standard.MeanZScore(zScoresByElement[element])
		rank, ok := findRank(rankStandards, zScore)
		if !ok {
			continue
		}

		elementEvaluations = append(elementEvaluations, newElementEvaluation(element, zScore, rank))
	}

	return elementEvaluations
}

type standardCurve struct {
	ages               []int
	means              []int64
	standardDeviations []int64
}

func newStandardCurve(ageGroupStandards []standard.AgeGroupStandard) (standardCurve, bool) {
	sorted := make([]standard.AgeGroupStandard, len(ageGroupStandards))
	copy(sorted, ageGroupStandards)
	sort.Slice(sorted, func(i, j int) bool {
		return isYoungerAgeRange(sorted[i].AgeRange(), sorted[j].AgeRange())
	})

	curve := standardCurve{}
	for _, ageGroupStandard := range sorted {
		age := ageGroupStandard.AgeRange().Median()
		if len(curve.ages) > 0 && age <= curve.ages[len(curve.ages)-1] {
			continue
		}
		curve.ages = append(curve.ages, age)
		curve.means = append(curve.means, toHundredths(ageGroupStandard.Mean().Float64()))
		curve.standardDeviations = append(curve.standardDeviations, toHundredths(ageGroupStandard.StandardDeviation().Float64()))
	}

	return curve, len(curve.ages) > 0
}

func extrapolatedMean(meanHundredths, meanDeltaHundredths, ageSpan, ageOffset int64) int64 {
	extrapolated := meanHundredths + divideRounded(meanDeltaHundredths*ageOffset, ageSpan)
	if extrapolated < 0 {
		return 0
	}
	return extrapolated
}

func (c standardCurve) at(age int) (int64, int64) {
	last := len(c.ages) - 1
	if last == 0 {
		return c.means[0], c.standardDeviations[0]
	}
	if age <= c.ages[0] {
		return extrapolatedMean(
			c.means[0],
			c.means[1]-c.means[0],
			int64(c.ages[1]-c.ages[0]),
			int64(age-c.ages[0]),
		), c.standardDeviations[0]
	}
	if age >= c.ages[last] {
		return extrapolatedMean(
			c.means[last],
			c.means[last]-c.means[last-1],
			int64(c.ages[last]-c.ages[last-1]),
			int64(age-c.ages[last]),
		), c.standardDeviations[last]
	}
	for i := 0; i < last; i++ {
		if age < c.ages[i] || age > c.ages[i+1] {
			continue
		}
		ageSpan := int64(c.ages[i+1] - c.ages[i])
		ageOffset := int64(age - c.ages[i])
		mean := c.means[i] + divideRounded((c.means[i+1]-c.means[i])*ageOffset, ageSpan)
		standardDeviation := c.standardDeviations[i] + divideRounded((c.standardDeviations[i+1]-c.standardDeviations[i])*ageOffset, ageSpan)
		return mean, standardDeviation
	}

	return c.means[last], c.standardDeviations[last]
}

func newMotorAge(
	judgedItems []judgedItem,
	ageGroupStandardsByItemID map[measurementitem.MeasurementItemID][]standard.AgeGroupStandard,
) *MotorAge {
	curves := make([]standardCurve, 0, len(judgedItems))
	values := make([]measurement.Value, 0, len(judgedItems))
	scoreDirections := make([]measurementitem.ScoreDirection, 0, len(judgedItems))
	youngestFrom, oldestTo, youngestNode, oldestNode := 0, 0, 0, 0
	for _, judged := range judgedItems {
		ageGroupStandards := ageGroupStandardsByItemID[judged.item.ID()]
		curve, ok := newStandardCurve(ageGroupStandards)
		if !ok {
			continue
		}
		if len(curves) == 0 {
			youngestFrom, oldestTo = ageGroupStandards[0].AgeRange().From(), ageGroupStandards[0].AgeRange().To()
			youngestNode, oldestNode = curve.ages[0], curve.ages[len(curve.ages)-1]
		}
		for _, ageGroupStandard := range ageGroupStandards {
			ageRange := ageGroupStandard.AgeRange()
			if ageRange.From() < youngestFrom {
				youngestFrom = ageRange.From()
			}
			if ageRange.To() > oldestTo {
				oldestTo = ageRange.To()
			}
		}
		if curve.ages[0] < youngestNode {
			youngestNode = curve.ages[0]
		}
		if node := curve.ages[len(curve.ages)-1]; node > oldestNode {
			oldestNode = node
		}
		curves = append(curves, curve)
		values = append(values, judged.value)
		scoreDirections = append(scoreDirections, *judged.item.ScoreDirection())
	}
	if len(curves) == 0 {
		return nil
	}

	lowestAge := youngestFrom - maxYoungerAgeGroupFallbackYears
	if lowestAge < 0 {
		lowestAge = 0
	}
	highestAge := oldestTo + maxOlderAgeGroupFallbackYears
	nodeSpanDoubled := int64(youngestNode + oldestNode)

	best := lowestAge
	bestDistance := int64(-1)
	bestNodeSpanOffset := int64(-1)
	for age := lowestAge; age <= highestAge; age++ {
		distance := int64(0)
		for i, curve := range curves {
			mean, standardDeviation := curve.at(age)
			zScore := zScoreHundredthsOf(values[i], mean, standardDeviation, scoreDirections[i])
			if zScore < 0 {
				zScore = -zScore
			}
			distance += zScore
		}
		nodeSpanOffset := int64(2*age) - nodeSpanDoubled
		if nodeSpanOffset < 0 {
			nodeSpanOffset = -nodeSpanOffset
		}
		if bestDistance < 0 || distance < bestDistance || (distance == bestDistance && nodeSpanOffset < bestNodeSpanOffset) {
			best, bestDistance, bestNodeSpanOffset = age, distance, nodeSpanOffset
		}
	}

	motorAge := NewMotorAge(best)

	return &motorAge
}

func NewEvaluation(
	m measurement.Measurement,
	items []measurementitem.MeasurementItem,
	gender standard.Gender,
	age int,
	ageGroupStandards []standard.AgeGroupStandard,
	rankStandards []standard.RankStandard,
) Evaluation {
	itemByID := make(map[measurementitem.MeasurementItemID]measurementitem.MeasurementItem, len(items))
	for _, item := range items {
		itemByID[item.ID()] = item
	}

	ageGroupStandardsByItemID := make(map[measurementitem.MeasurementItemID][]standard.AgeGroupStandard)
	for _, ageGroupStandard := range ageGroupStandards {
		if ageGroupStandard.Gender() != gender {
			continue
		}
		measurementItemID := ageGroupStandard.MeasurementItemID()
		ageGroupStandardsByItemID[measurementItemID] = append(ageGroupStandardsByItemID[measurementItemID], ageGroupStandard)
	}

	entries := m.Entries()
	baseHundredths, hasBase := heightHundredths(entries, itemByID)
	itemEvaluations := make([]ItemEvaluation, 0, len(entries))
	judgedItems := make([]judgedItem, 0, len(entries))
	for _, entry := range entries {
		if entry.Unmeasurable() {
			continue
		}
		item, ok := itemByID[entry.MeasurementItemID()]
		if !ok {
			continue
		}
		scoreDirection := item.ScoreDirection()
		if scoreDirection == nil {
			continue
		}
		value, ok := representativeValue(entry, item, *scoreDirection)
		if !ok {
			continue
		}
		value, ok = normalizedValue(value, item.Normalization(), baseHundredths, hasBase)
		if !ok {
			continue
		}
		ageGroupStandard, ok := findAgeGroupStandard(ageGroupStandardsByItemID[item.ID()], age)
		if !ok {
			continue
		}
		zScore, err := standard.NewZScore(fromHundredths(zScoreHundredths(value, ageGroupStandard, *scoreDirection)))
		if err != nil {
			continue
		}

		judgedItems = append(judgedItems, judgedItem{item: item, value: value})

		rank, ok := findRank(rankStandards, zScore)
		if !ok {
			continue
		}

		itemEvaluations = append(itemEvaluations, newItemEvaluation(item.ID(), value, ageGroupStandard.Mean(), zScore, rank))
	}

	return &evaluation{
		itemEvaluations:    itemEvaluations,
		elementEvaluations: newElementEvaluations(itemEvaluations, itemByID, rankStandards),
		motorAge:           newMotorAge(judgedItems, ageGroupStandardsByItemID),
	}
}
