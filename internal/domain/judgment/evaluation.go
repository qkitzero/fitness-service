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

func representativeValue(entry measurement.MeasurementEntry, item measurementitem.MeasurementItem, scoreDirection measurementitem.ScoreDirection) (measurement.Value, bool) {
	bestBySide := make(map[measurement.Side]int64)
	sides := make([]measurement.Side, 0, 3)
	for _, measurementValue := range entry.Values() {
		value := measurementValue.Value()
		if value == nil {
			continue
		}
		hundredths := toHundredths(value.Float64())
		side := measurementValue.Side()
		best, ok := bestBySide[side]
		if !ok {
			bestBySide[side] = hundredths
			sides = append(sides, side)
			continue
		}
		if isBetter(hundredths, best, scoreDirection) {
			bestBySide[side] = hundredths
		}
	}
	if len(sides) == 0 {
		return measurement.Value(0), false
	}

	representative := bestBySide[sides[0]]
	switch item.SideAggregation() {
	case measurementitem.SideAggregationBest:
		for _, side := range sides[1:] {
			if isBetter(bestBySide[side], representative, scoreDirection) {
				representative = bestBySide[side]
			}
		}
	case measurementitem.SideAggregationWorst:
		for _, side := range sides[1:] {
			if isBetter(representative, bestBySide[side], scoreDirection) {
				representative = bestBySide[side]
			}
		}
	case measurementitem.SideAggregationMean:
		sum := int64(0)
		for _, side := range sides {
			sum += bestBySide[side]
		}
		representative = divideRounded(sum, int64(len(sides)))
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

func zScoreHundredths(value measurement.Value, ageGroupStandard standard.AgeGroupStandard, scoreDirection measurementitem.ScoreDirection) int64 {
	deviation := toHundredths(value.Float64()) - toHundredths(ageGroupStandard.Mean().Float64())
	zScore := divideRounded(evaluationScale*deviation, toHundredths(ageGroupStandard.StandardDeviation().Float64()))
	if scoreDirection == measurementitem.ScoreDirectionLowerIsBetter {
		return -zScore
	}
	return zScore
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

func isPreferredAgeRanges(candidate, current []standard.AgeRange) bool {
	if len(candidate) != len(current) {
		return len(candidate) > len(current)
	}
	for i := range candidate {
		if candidate[i] != current[i] {
			return isYoungerAgeRange(candidate[i], current[i])
		}
	}
	return false
}

func newMotorAge(
	judgedItems []judgedItem,
	ageGroupStandardsByItemID map[measurementitem.MeasurementItemID][]standard.AgeGroupStandard,
) *MotorAge {
	standardsByItemAndAgeRange := make(map[measurementitem.MeasurementItemID]map[standard.AgeRange]standard.AgeGroupStandard, len(judgedItems))
	for _, judged := range judgedItems {
		standardsByAgeRange := make(map[standard.AgeRange]standard.AgeGroupStandard)
		for _, ageGroupStandard := range ageGroupStandardsByItemID[judged.item.ID()] {
			standardsByAgeRange[ageGroupStandard.AgeRange()] = ageGroupStandard
		}
		standardsByItemAndAgeRange[judged.item.ID()] = standardsByAgeRange
	}

	ageRanges := make([]standard.AgeRange, 0)
	for _, judged := range judgedItems {
		candidate := make([]standard.AgeRange, 0, len(standardsByItemAndAgeRange[judged.item.ID()]))
		for ageRange := range standardsByItemAndAgeRange[judged.item.ID()] {
			candidate = append(candidate, ageRange)
		}
		sort.Slice(candidate, func(i, j int) bool {
			return isYoungerAgeRange(candidate[i], candidate[j])
		})
		if isPreferredAgeRanges(candidate, ageRanges) {
			ageRanges = candidate
		}
	}
	if len(ageRanges) == 0 {
		return nil
	}

	coveringItems := make([]judgedItem, 0, len(judgedItems))
	for _, judged := range judgedItems {
		covers := true
		for _, ageRange := range ageRanges {
			if _, ok := standardsByItemAndAgeRange[judged.item.ID()][ageRange]; !ok {
				covers = false
				break
			}
		}
		if covers {
			coveringItems = append(coveringItems, judged)
		}
	}

	best := ageRanges[0]
	bestDistance := int64(-1)
	for _, ageRange := range ageRanges {
		distance := int64(0)
		for _, judged := range coveringItems {
			zScore := zScoreHundredths(judged.value, standardsByItemAndAgeRange[judged.item.ID()][ageRange], *judged.item.ScoreDirection())
			if zScore < 0 {
				zScore = -zScore
			}
			distance += zScore
		}
		if bestDistance < 0 || distance < bestDistance {
			best = ageRange
			bestDistance = distance
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
