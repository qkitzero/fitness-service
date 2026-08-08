package judgment

import (
	"sort"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/standard"
	"github.com/qkitzero/fitness-service/internal/domain/training"
)

const decadeWidth = 10

var levelByDeficientRank = map[standard.Rank]training.Level{
	standard.RankE: 1,
	standard.RankD: 2,
}

type Prescription interface {
	PrescribedMenus() []PrescribedMenu
}

type prescription struct {
	prescribedMenus []PrescribedMenu
}

func (p prescription) PrescribedMenus() []PrescribedMenu {
	prescribedMenus := make([]PrescribedMenu, len(p.prescribedMenus))
	copy(prescribedMenus, p.prescribedMenus)
	return prescribedMenus
}

type elementPart struct {
	element measurementitem.Element
	part    training.Part
}

func deficientLevels(evaluation Evaluation) map[measurementitem.Element]training.Level {
	levels := make(map[measurementitem.Element]training.Level)
	if evaluation == nil {
		return levels
	}
	for _, elementEvaluation := range evaluation.ElementEvaluations() {
		level, ok := levelByDeficientRank[elementEvaluation.Rank()]
		if !ok {
			continue
		}
		levels[elementEvaluation.Element()] = level
	}
	return levels
}

func AgeDecade(age int) (training.Decade, bool) {
	decade, err := training.NewDecade(age / decadeWidth * decadeWidth)
	if err != nil {
		return training.Decade(0), false
	}
	return decade, true
}

func resolveElementMenu(elementMenus []training.ElementMenu, key elementPart, level training.Level) (training.ElementMenu, bool) {
	var resolved training.ElementMenu
	for _, elementMenu := range elementMenus {
		if elementMenu.Element() != key.element || elementMenu.Part() != key.part || elementMenu.Level() > level {
			continue
		}
		if resolved == nil || elementMenu.Level() > resolved.Level() {
			resolved = elementMenu
		}
	}
	if resolved == nil {
		return nil, false
	}
	return resolved, true
}

func newElementPrescribedMenus(
	evaluation Evaluation,
	elementMenus []training.ElementMenu,
	trainingMenuByID map[training.TrainingMenuID]training.TrainingMenu,
) []PrescribedMenu {
	levels := deficientLevels(evaluation)

	elementParts := make([]elementPart, 0, len(elementMenus))
	seen := make(map[elementPart]struct{}, len(elementMenus))
	for _, elementMenu := range elementMenus {
		if _, ok := levels[elementMenu.Element()]; !ok {
			continue
		}
		key := elementPart{element: elementMenu.Element(), part: elementMenu.Part()}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		elementParts = append(elementParts, key)
	}

	sort.Slice(elementParts, func(i, j int) bool {
		iElementOrder, jElementOrder := elementParts[i].element.Order(), elementParts[j].element.Order()
		if iElementOrder != jElementOrder {
			return iElementOrder < jElementOrder
		}
		return elementParts[i].part.Order() < elementParts[j].part.Order()
	})

	prescribedMenus := make([]PrescribedMenu, 0, len(elementParts))
	for _, key := range elementParts {
		elementMenu, ok := resolveElementMenu(elementMenus, key, levels[key.element])
		if !ok {
			continue
		}
		trainingMenu, ok := trainingMenuByID[elementMenu.TrainingMenuID()]
		if !ok {
			continue
		}
		element, part := key.element, key.part
		prescribedMenus = append(prescribedMenus, newPrescribedMenu(
			PrescriptionSourceElement,
			&element,
			&part,
			trainingMenu.ID(),
			trainingMenu.Name(),
			trainingMenu.Amount(),
			trainingMenu.Unit(),
			trainingMenu.Sets(),
		))
	}

	return prescribedMenus
}

func newFixedPrescribedMenus(
	fixedMenus []training.FixedMenu,
	trainingMenuByID map[training.TrainingMenuID]training.TrainingMenu,
) []PrescribedMenu {
	sorted := make([]training.FixedMenu, len(fixedMenus))
	copy(sorted, fixedMenus)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].SortOrder() < sorted[j].SortOrder()
	})

	prescribedMenus := make([]PrescribedMenu, 0, len(sorted))
	for _, fixedMenu := range sorted {
		trainingMenu, ok := trainingMenuByID[fixedMenu.TrainingMenuID()]
		if !ok {
			continue
		}
		prescribedMenus = append(prescribedMenus, newPrescribedMenu(
			PrescriptionSourceFixed,
			nil,
			nil,
			trainingMenu.ID(),
			trainingMenu.Name(),
			trainingMenu.Amount(),
			trainingMenu.Unit(),
			trainingMenu.Sets(),
		))
	}

	return prescribedMenus
}

func newAgeDecadePrescribedMenus(
	age int,
	ageDecadeMenus []training.AgeDecadeMenu,
	trainingMenuByID map[training.TrainingMenuID]training.TrainingMenu,
) []PrescribedMenu {
	decade, ok := AgeDecade(age)
	if !ok {
		return nil
	}

	sorted := make([]training.AgeDecadeMenu, 0, len(ageDecadeMenus))
	for _, ageDecadeMenu := range ageDecadeMenus {
		if ageDecadeMenu.Decade() != decade {
			continue
		}
		sorted = append(sorted, ageDecadeMenu)
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].SortOrder() < sorted[j].SortOrder()
	})

	prescribedMenus := make([]PrescribedMenu, 0, len(sorted))
	for _, ageDecadeMenu := range sorted {
		trainingMenu, ok := trainingMenuByID[ageDecadeMenu.TrainingMenuID()]
		if !ok {
			continue
		}
		prescribedMenus = append(prescribedMenus, newPrescribedMenu(
			PrescriptionSourceAgeDecade,
			nil,
			nil,
			trainingMenu.ID(),
			trainingMenu.Name(),
			trainingMenu.Amount(),
			trainingMenu.Unit(),
			trainingMenu.Sets(),
		))
	}

	return prescribedMenus
}

func newTrainingMenuByID(trainingMenus []training.TrainingMenu) map[training.TrainingMenuID]training.TrainingMenu {
	trainingMenuByID := make(map[training.TrainingMenuID]training.TrainingMenu, len(trainingMenus))
	for _, trainingMenu := range trainingMenus {
		trainingMenuByID[trainingMenu.ID()] = trainingMenu
	}
	return trainingMenuByID
}

func NewPrescription(
	evaluation Evaluation,
	age int,
	elementMenus []training.ElementMenu,
	fixedMenus []training.FixedMenu,
	ageDecadeMenus []training.AgeDecadeMenu,
	trainingMenus []training.TrainingMenu,
) Prescription {
	trainingMenuByID := newTrainingMenuByID(trainingMenus)

	groups := [][]PrescribedMenu{
		newElementPrescribedMenus(evaluation, elementMenus, trainingMenuByID),
		newFixedPrescribedMenus(fixedMenus, trainingMenuByID),
		newAgeDecadePrescribedMenus(age, ageDecadeMenus, trainingMenuByID),
	}

	prescribedMenus := make([]PrescribedMenu, 0, len(elementMenus)+len(fixedMenus)+len(ageDecadeMenus))
	seen := make(map[training.TrainingMenuID]struct{})
	for _, group := range groups {
		for _, prescribedMenu := range group {
			if _, ok := seen[prescribedMenu.TrainingMenuID()]; ok {
				continue
			}
			seen[prescribedMenu.TrainingMenuID()] = struct{}{}
			prescribedMenus = append(prescribedMenus, prescribedMenu)
		}
	}

	return &prescription{prescribedMenus: prescribedMenus}
}

func NewPrescriptionFromOverrides(
	overrides []PrescribedMenuOverride,
	trainingMenus []training.TrainingMenu,
) Prescription {
	trainingMenuByID := newTrainingMenuByID(trainingMenus)

	sorted := make([]PrescribedMenuOverride, len(overrides))
	copy(sorted, overrides)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].SortOrder() < sorted[j].SortOrder()
	})

	prescribedMenus := make([]PrescribedMenu, 0, len(sorted))
	for _, override := range sorted {
		trainingMenu, ok := trainingMenuByID[override.TrainingMenuID()]
		if !ok {
			continue
		}
		prescribedMenus = append(prescribedMenus, newPrescribedMenu(
			PrescriptionSourceManual,
			override.Element(),
			override.Part(),
			override.TrainingMenuID(),
			trainingMenu.Name(),
			override.Amount(),
			override.Unit(),
			override.Sets(),
		))
	}

	return &prescription{prescribedMenus: prescribedMenus}
}
