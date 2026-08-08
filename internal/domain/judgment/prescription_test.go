package judgment

import (
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/standard"
	"github.com/qkitzero/fitness-service/internal/domain/training"
)

type trainingMenuSpec struct {
	code    string
	element measurementitem.Element
	part    training.Part
	amount  int
	unit    training.Unit
	sets    int
}

type elementMenuSpec struct {
	element measurementitem.Element
	part    training.Part
	level   int
	code    string
}

type sortedMenuSpec struct {
	decade    int
	sortOrder int
	code      string
}

type rankedElement struct {
	element measurementitem.Element
	rank    standard.Rank
}

type prescribedMenuSpec struct {
	source  PrescriptionSource
	element *measurementitem.Element
	part    *training.Part
	code    string
}

var trainingMenuSpecs = []trainingMenuSpec{
	{"wall_push", measurementitem.ElementMuscleStrength, training.PartUpperLimb, 10, training.UnitReps, 3},
	{"chair_squat", measurementitem.ElementMuscleStrength, training.PartLowerLimb, 12, training.UnitReps, 2},
	{"core_stabilize", measurementitem.ElementMuscleStrength, training.PartWholeBody, 30, training.UnitSeconds, 1},
	{"burpee_jump", measurementitem.ElementMuscleStrength, training.PartWholeBody, 8, training.UnitReps, 3},
	{"heel_raise", measurementitem.ElementMuscleEndurance, training.PartLowerLimb, 15, training.UnitReps, 2},
	{"whole_body_stretch", measurementitem.ElementFlexibility, training.PartWholeBody, 5, training.UnitMinutes, 1},
	{"single_leg_stand", measurementitem.ElementBalance, training.PartLowerLimb, 30, training.UnitSeconds, 2},
	{"radio_exercise", measurementitem.ElementMobility, training.PartWholeBody, 3, training.UnitMinutes, 1},
	{"walking", measurementitem.ElementMuscleEndurance, training.PartWholeBody, 20, training.UnitMinutes, 1},
	{"chair_yoga", measurementitem.ElementFlexibility, training.PartWholeBody, 10, training.UnitMinutes, 1},
}

func TestNewPrescription(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

	muscleStrength := measurementitem.ElementMuscleStrength
	flexibility := measurementitem.ElementFlexibility
	balance := measurementitem.ElementBalance
	upperLimb := training.PartUpperLimb
	lowerLimb := training.PartLowerLimb
	wholeBody := training.PartWholeBody

	trainingMenus := make([]training.TrainingMenu, 0, len(trainingMenuSpecs))
	menuByCode := make(map[string]training.TrainingMenu, len(trainingMenuSpecs))
	for _, spec := range trainingMenuSpecs {
		code, _ := training.NewCode(spec.code)
		name, _ := training.NewName(spec.code)
		amount, _ := training.NewAmount(spec.amount)
		sets, _ := training.NewSets(spec.sets)
		instruction, _ := training.NewInstruction(spec.code)
		trainingMenu := training.NewTrainingMenu(training.NewTrainingMenuID(), code, name, spec.element, spec.part, amount, spec.unit, sets, instruction, createdAt, updatedAt)
		trainingMenus = append(trainingMenus, trainingMenu)
		menuByCode[spec.code] = trainingMenu
	}

	tests := []struct {
		name             string
		nilEvaluation    bool
		rankedElements   []rankedElement
		age              int
		elementMenuSpecs []elementMenuSpec
		fixedMenuSpecs   []sortedMenuSpec
		decadeMenuSpecs  []sortedMenuSpec
		unknownMenuRule  bool
		want             []prescribedMenuSpec
	}{
		{
			name:             "success prescribe level 1 for rank e",
			rankedElements:   []rankedElement{{measurementitem.ElementMuscleStrength, standard.RankE}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementMuscleStrength, training.PartWholeBody, 1, "core_stabilize"}, {measurementitem.ElementMuscleStrength, training.PartWholeBody, 2, "burpee_jump"}},
			want:             []prescribedMenuSpec{{PrescriptionSourceElement, &muscleStrength, &wholeBody, "core_stabilize"}},
		},
		{
			name:             "success prescribe level 2 for rank d",
			rankedElements:   []rankedElement{{measurementitem.ElementMuscleStrength, standard.RankD}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementMuscleStrength, training.PartWholeBody, 1, "core_stabilize"}, {measurementitem.ElementMuscleStrength, training.PartWholeBody, 2, "burpee_jump"}},
			want:             []prescribedMenuSpec{{PrescriptionSourceElement, &muscleStrength, &wholeBody, "burpee_jump"}},
		},
		{
			name:             "success prescribe nothing for rank c",
			rankedElements:   []rankedElement{{measurementitem.ElementMuscleStrength, standard.RankC}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementMuscleStrength, training.PartWholeBody, 1, "core_stabilize"}, {measurementitem.ElementMuscleStrength, training.PartWholeBody, 3, "burpee_jump"}},
		},
		{
			name:             "success prescribe nothing for rank b",
			rankedElements:   []rankedElement{{measurementitem.ElementMuscleStrength, standard.RankB}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementMuscleStrength, training.PartWholeBody, 4, "core_stabilize"}},
		},
		{
			name:             "success prescribe nothing for rank a",
			rankedElements:   []rankedElement{{measurementitem.ElementMuscleStrength, standard.RankA}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementMuscleStrength, training.PartWholeBody, 5, "core_stabilize"}},
		},
		{
			name:             "success prescribe nothing for an unmapped rank",
			rankedElements:   []rankedElement{{measurementitem.ElementMuscleStrength, standard.Rank("F")}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementMuscleStrength, training.PartWholeBody, 1, "core_stabilize"}},
		},
		{
			name:             "success prescribe every part of a deficient element",
			rankedElements:   []rankedElement{{measurementitem.ElementMuscleStrength, standard.RankD}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementMuscleStrength, training.PartUpperLimb, 2, "wall_push"}, {measurementitem.ElementMuscleStrength, training.PartLowerLimb, 2, "chair_squat"}, {measurementitem.ElementMuscleStrength, training.PartWholeBody, 2, "core_stabilize"}},
			want: []prescribedMenuSpec{
				{PrescriptionSourceElement, &muscleStrength, &upperLimb, "wall_push"},
				{PrescriptionSourceElement, &muscleStrength, &lowerLimb, "chair_squat"},
				{PrescriptionSourceElement, &muscleStrength, &wholeBody, "core_stabilize"},
			},
		},
		{
			name:             "success fall back to the closest lower level",
			rankedElements:   []rankedElement{{measurementitem.ElementMuscleStrength, standard.RankD}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementMuscleStrength, training.PartWholeBody, 1, "core_stabilize"}, {measurementitem.ElementMuscleStrength, training.PartWholeBody, 3, "burpee_jump"}},
			want:             []prescribedMenuSpec{{PrescriptionSourceElement, &muscleStrength, &wholeBody, "core_stabilize"}},
		},
		{
			name:             "success never fall back to a higher level",
			rankedElements:   []rankedElement{{measurementitem.ElementMuscleStrength, standard.RankD}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementMuscleStrength, training.PartWholeBody, 3, "burpee_jump"}},
		},
		{
			name:             "success prefer the exact level over a lower one",
			rankedElements:   []rankedElement{{measurementitem.ElementMuscleStrength, standard.RankD}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementMuscleStrength, training.PartWholeBody, 2, "burpee_jump"}, {measurementitem.ElementMuscleStrength, training.PartWholeBody, 1, "core_stabilize"}},
			want:             []prescribedMenuSpec{{PrescriptionSourceElement, &muscleStrength, &wholeBody, "burpee_jump"}},
		},
		{
			name:            "success compose the fixed and age decade menus",
			age:             45,
			fixedMenuSpecs:  []sortedMenuSpec{{0, 1, "radio_exercise"}, {0, 2, "whole_body_stretch"}},
			decadeMenuSpecs: []sortedMenuSpec{{40, 1, "walking"}, {60, 1, "chair_yoga"}},
			want: []prescribedMenuSpec{
				{PrescriptionSourceFixed, nil, nil, "radio_exercise"},
				{PrescriptionSourceFixed, nil, nil, "whole_body_stretch"},
				{PrescriptionSourceAgeDecade, nil, nil, "walking"},
			},
		},
		{
			name:            "success floor the age of a nonagenarian to ninety",
			age:             95,
			decadeMenuSpecs: []sortedMenuSpec{{90, 1, "walking"}, {40, 1, "chair_yoga"}},
			want:            []prescribedMenuSpec{{PrescriptionSourceAgeDecade, nil, nil, "walking"}},
		},
		{
			name:            "success prescribe no age decade menu for a centenarian",
			age:             105,
			decadeMenuSpecs: []sortedMenuSpec{{90, 1, "walking"}},
		},
		{
			name:            "success prescribe no age decade menu for a child",
			age:             9,
			decadeMenuSpecs: []sortedMenuSpec{{10, 1, "walking"}},
		},
		{
			name:            "success prescribe no age decade menu for an unseeded decade",
			age:             62,
			decadeMenuSpecs: []sortedMenuSpec{{40, 1, "walking"}},
		},
		{
			name:             "success keep the element menu when a fixed menu repeats it",
			rankedElements:   []rankedElement{{measurementitem.ElementFlexibility, standard.RankD}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementFlexibility, training.PartWholeBody, 2, "whole_body_stretch"}},
			fixedMenuSpecs:   []sortedMenuSpec{{0, 1, "whole_body_stretch"}, {0, 2, "radio_exercise"}},
			want: []prescribedMenuSpec{
				{PrescriptionSourceElement, &flexibility, &wholeBody, "whole_body_stretch"},
				{PrescriptionSourceFixed, nil, nil, "radio_exercise"},
			},
		},
		{
			name:            "success keep the fixed menu when an age decade menu repeats it",
			age:             62,
			fixedMenuSpecs:  []sortedMenuSpec{{0, 1, "walking"}},
			decadeMenuSpecs: []sortedMenuSpec{{60, 1, "walking"}, {60, 2, "chair_yoga"}},
			want: []prescribedMenuSpec{
				{PrescriptionSourceFixed, nil, nil, "walking"},
				{PrescriptionSourceAgeDecade, nil, nil, "chair_yoga"},
			},
		},
		{
			name:             "success order the element menus by element then part",
			rankedElements:   []rankedElement{{measurementitem.ElementBalance, standard.RankE}, {measurementitem.ElementMuscleStrength, standard.RankD}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementBalance, training.PartLowerLimb, 1, "single_leg_stand"}, {measurementitem.ElementMuscleStrength, training.PartLowerLimb, 2, "chair_squat"}, {measurementitem.ElementMuscleStrength, training.PartUpperLimb, 2, "wall_push"}},
			fixedMenuSpecs:   []sortedMenuSpec{{0, 2, "radio_exercise"}, {0, 1, "whole_body_stretch"}},
			decadeMenuSpecs:  []sortedMenuSpec{{60, 2, "chair_yoga"}, {60, 1, "walking"}},
			want: []prescribedMenuSpec{
				{PrescriptionSourceElement, &muscleStrength, &upperLimb, "wall_push"},
				{PrescriptionSourceElement, &muscleStrength, &lowerLimb, "chair_squat"},
				{PrescriptionSourceElement, &balance, &lowerLimb, "single_leg_stand"},
				{PrescriptionSourceFixed, nil, nil, "whole_body_stretch"},
				{PrescriptionSourceFixed, nil, nil, "radio_exercise"},
				{PrescriptionSourceAgeDecade, nil, nil, "walking"},
				{PrescriptionSourceAgeDecade, nil, nil, "chair_yoga"},
			},
		},
		{
			name:             "success prescribe the fixed and age decade menus for an empty evaluation",
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementMuscleStrength, training.PartWholeBody, 1, "core_stabilize"}},
			fixedMenuSpecs:   []sortedMenuSpec{{0, 1, "radio_exercise"}},
			decadeMenuSpecs:  []sortedMenuSpec{{60, 1, "walking"}},
			want: []prescribedMenuSpec{
				{PrescriptionSourceFixed, nil, nil, "radio_exercise"},
				{PrescriptionSourceAgeDecade, nil, nil, "walking"},
			},
		},
		{
			name:            "success prescribe the fixed and age decade menus without an evaluation",
			nilEvaluation:   true,
			age:             62,
			fixedMenuSpecs:  []sortedMenuSpec{{0, 1, "radio_exercise"}},
			decadeMenuSpecs: []sortedMenuSpec{{60, 1, "walking"}},
			want: []prescribedMenuSpec{
				{PrescriptionSourceFixed, nil, nil, "radio_exercise"},
				{PrescriptionSourceAgeDecade, nil, nil, "walking"},
			},
		},
		{
			name:           "success prescribe nothing when no rule is seeded",
			rankedElements: []rankedElement{{measurementitem.ElementMuscleStrength, standard.RankE}},
			age:            62,
		},
		{
			name:             "success skip the rules pointing at an unknown training menu",
			rankedElements:   []rankedElement{{measurementitem.ElementMuscleStrength, standard.RankD}},
			age:              62,
			elementMenuSpecs: []elementMenuSpec{{measurementitem.ElementMuscleStrength, training.PartWholeBody, 2, "core_stabilize"}},
			fixedMenuSpecs:   []sortedMenuSpec{{0, 1, "radio_exercise"}},
			decadeMenuSpecs:  []sortedMenuSpec{{60, 1, "walking"}},
			unknownMenuRule:  true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var targetEvaluation Evaluation
			if !tt.nilEvaluation {
				elementEvaluations := make([]ElementEvaluation, 0, len(tt.rankedElements))
				for _, ranked := range tt.rankedElements {
					elementEvaluations = append(elementEvaluations, newElementEvaluation(ranked.element, standard.ZScore(0), ranked.rank))
				}
				targetEvaluation = &evaluation{elementEvaluations: elementEvaluations}
			}

			trainingMenuIDOf := func(code string) training.TrainingMenuID {
				if tt.unknownMenuRule {
					return training.NewTrainingMenuID()
				}
				return menuByCode[code].ID()
			}

			elementMenus := make([]training.ElementMenu, 0, len(tt.elementMenuSpecs))
			for _, spec := range tt.elementMenuSpecs {
				level, _ := training.NewLevel(spec.level)
				elementMenus = append(elementMenus, training.NewElementMenu(spec.element, spec.part, level, trainingMenuIDOf(spec.code), createdAt, updatedAt))
			}

			fixedMenus := make([]training.FixedMenu, 0, len(tt.fixedMenuSpecs))
			for _, spec := range tt.fixedMenuSpecs {
				sortOrder, _ := training.NewSortOrder(spec.sortOrder)
				fixedMenus = append(fixedMenus, training.NewFixedMenu(sortOrder, trainingMenuIDOf(spec.code), createdAt, updatedAt))
			}

			ageDecadeMenus := make([]training.AgeDecadeMenu, 0, len(tt.decadeMenuSpecs))
			for _, spec := range tt.decadeMenuSpecs {
				decade, _ := training.NewDecade(spec.decade)
				sortOrder, _ := training.NewSortOrder(spec.sortOrder)
				ageDecadeMenus = append(ageDecadeMenus, training.NewAgeDecadeMenu(decade, sortOrder, trainingMenuIDOf(spec.code), createdAt, updatedAt))
			}

			p := NewPrescription(targetEvaluation, tt.age, elementMenus, fixedMenus, ageDecadeMenus, trainingMenus)

			prescribedMenus := p.PrescribedMenus()
			if len(prescribedMenus) != len(tt.want) {
				t.Fatalf("len(PrescribedMenus()) = %v, want %v", len(prescribedMenus), len(tt.want))
			}
			for i, want := range tt.want {
				prescribedMenu := prescribedMenus[i]
				trainingMenu := menuByCode[want.code]
				if prescribedMenu.Source() != want.source {
					t.Errorf("PrescribedMenus()[%d].Source() = %v, want %v", i, prescribedMenu.Source(), want.source)
				}
				switch {
				case want.element == nil:
					if prescribedMenu.Element() != nil {
						t.Errorf("PrescribedMenus()[%d].Element() = %v, want nil", i, *prescribedMenu.Element())
					}
				case prescribedMenu.Element() == nil:
					t.Errorf("PrescribedMenus()[%d].Element() = nil, want %v", i, *want.element)
				case *prescribedMenu.Element() != *want.element:
					t.Errorf("PrescribedMenus()[%d].Element() = %v, want %v", i, *prescribedMenu.Element(), *want.element)
				}
				switch {
				case want.part == nil:
					if prescribedMenu.Part() != nil {
						t.Errorf("PrescribedMenus()[%d].Part() = %v, want nil", i, *prescribedMenu.Part())
					}
				case prescribedMenu.Part() == nil:
					t.Errorf("PrescribedMenus()[%d].Part() = nil, want %v", i, *want.part)
				case *prescribedMenu.Part() != *want.part:
					t.Errorf("PrescribedMenus()[%d].Part() = %v, want %v", i, *prescribedMenu.Part(), *want.part)
				}
				if prescribedMenu.TrainingMenuID() != trainingMenu.ID() {
					t.Errorf("PrescribedMenus()[%d].TrainingMenuID() = %v, want %v", i, prescribedMenu.TrainingMenuID(), trainingMenu.ID())
				}
				if prescribedMenu.TrainingMenuName() != trainingMenu.Name() {
					t.Errorf("PrescribedMenus()[%d].TrainingMenuName() = %v, want %v", i, prescribedMenu.TrainingMenuName(), trainingMenu.Name())
				}
				if prescribedMenu.Amount() != trainingMenu.Amount() {
					t.Errorf("PrescribedMenus()[%d].Amount() = %v, want %v", i, prescribedMenu.Amount(), trainingMenu.Amount())
				}
				if prescribedMenu.Unit() != trainingMenu.Unit() {
					t.Errorf("PrescribedMenus()[%d].Unit() = %v, want %v", i, prescribedMenu.Unit(), trainingMenu.Unit())
				}
				if prescribedMenu.Sets() != trainingMenu.Sets() {
					t.Errorf("PrescribedMenus()[%d].Sets() = %v, want %v", i, prescribedMenu.Sets(), trainingMenu.Sets())
				}
			}
		})
	}
}

func TestNewPrescriptionFromOverrides(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)

	muscleStrength := measurementitem.ElementMuscleStrength
	upperLimb := training.PartUpperLimb

	trainingMenus := make([]training.TrainingMenu, 0, len(trainingMenuSpecs))
	menuByCode := make(map[string]training.TrainingMenu, len(trainingMenuSpecs))
	for _, spec := range trainingMenuSpecs {
		code, _ := training.NewCode(spec.code)
		name, _ := training.NewName(spec.code)
		amount, _ := training.NewAmount(spec.amount)
		sets, _ := training.NewSets(spec.sets)
		instruction, _ := training.NewInstruction(spec.code)
		trainingMenu := training.NewTrainingMenu(training.NewTrainingMenuID(), code, name, spec.element, spec.part, amount, spec.unit, sets, instruction, createdAt, updatedAt)
		trainingMenus = append(trainingMenus, trainingMenu)
		menuByCode[spec.code] = trainingMenu
	}

	type overrideSpec struct {
		sortOrder int
		element   *measurementitem.Element
		part      *training.Part
		code      string
		amount    int
		unit      training.Unit
		sets      int
	}

	tests := []struct {
		name            string
		overrideSpecs   []overrideSpec
		unknownMenuRule bool
		want            []overrideSpec
	}{
		{
			name:          "success prescribe the overrides as manual menus",
			overrideSpecs: []overrideSpec{{1, &muscleStrength, &upperLimb, "wall_push", 20, training.UnitReps, 5}},
			want:          []overrideSpec{{1, &muscleStrength, &upperLimb, "wall_push", 20, training.UnitReps, 5}},
		},
		{
			name:          "success prescribe the overrides without labels",
			overrideSpecs: []overrideSpec{{1, nil, nil, "walking", 30, training.UnitMinutes, 1}},
			want:          []overrideSpec{{1, nil, nil, "walking", 30, training.UnitMinutes, 1}},
		},
		{
			name:          "success order the overrides by sort order",
			overrideSpecs: []overrideSpec{{3, nil, nil, "walking", 30, training.UnitMinutes, 1}, {1, nil, nil, "wall_push", 20, training.UnitReps, 5}, {2, nil, nil, "chair_squat", 15, training.UnitReps, 2}},
			want:          []overrideSpec{{1, nil, nil, "wall_push", 20, training.UnitReps, 5}, {2, nil, nil, "chair_squat", 15, training.UnitReps, 2}, {3, nil, nil, "walking", 30, training.UnitMinutes, 1}},
		},
		{
			name:          "success keep the overrides repeating a training menu",
			overrideSpecs: []overrideSpec{{1, nil, nil, "walking", 30, training.UnitMinutes, 1}, {2, nil, nil, "walking", 10, training.UnitMinutes, 2}},
			want:          []overrideSpec{{1, nil, nil, "walking", 30, training.UnitMinutes, 1}, {2, nil, nil, "walking", 10, training.UnitMinutes, 2}},
		},
		{
			name: "success prescribe nothing without overrides",
		},
		{
			name:            "success skip the overrides pointing at an unknown training menu",
			overrideSpecs:   []overrideSpec{{1, nil, nil, "walking", 30, training.UnitMinutes, 1}},
			unknownMenuRule: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			measurementID := measurement.NewMeasurementID()
			overrides := make([]PrescribedMenuOverride, 0, len(tt.overrideSpecs))
			for _, spec := range tt.overrideSpecs {
				sortOrder, _ := training.NewSortOrder(spec.sortOrder)
				amount, _ := training.NewAmount(spec.amount)
				sets, _ := training.NewSets(spec.sets)
				trainingMenuID := menuByCode[spec.code].ID()
				if tt.unknownMenuRule {
					trainingMenuID = training.NewTrainingMenuID()
				}
				overrides = append(overrides, NewPrescribedMenuOverride(NewPrescribedMenuOverrideID(), measurementID, sortOrder, spec.element, spec.part, trainingMenuID, amount, spec.unit, sets, createdAt, updatedAt))
			}

			p := NewPrescriptionFromOverrides(overrides, trainingMenus)

			prescribedMenus := p.PrescribedMenus()
			if len(prescribedMenus) != len(tt.want) {
				t.Fatalf("len(PrescribedMenus()) = %v, want %v", len(prescribedMenus), len(tt.want))
			}
			for i, want := range tt.want {
				prescribedMenu := prescribedMenus[i]
				trainingMenu := menuByCode[want.code]
				if prescribedMenu.Source() != PrescriptionSourceManual {
					t.Errorf("PrescribedMenus()[%d].Source() = %v, want %v", i, prescribedMenu.Source(), PrescriptionSourceManual)
				}
				switch {
				case want.element == nil:
					if prescribedMenu.Element() != nil {
						t.Errorf("PrescribedMenus()[%d].Element() = %v, want nil", i, *prescribedMenu.Element())
					}
				case prescribedMenu.Element() == nil:
					t.Errorf("PrescribedMenus()[%d].Element() = nil, want %v", i, *want.element)
				case *prescribedMenu.Element() != *want.element:
					t.Errorf("PrescribedMenus()[%d].Element() = %v, want %v", i, *prescribedMenu.Element(), *want.element)
				}
				switch {
				case want.part == nil:
					if prescribedMenu.Part() != nil {
						t.Errorf("PrescribedMenus()[%d].Part() = %v, want nil", i, *prescribedMenu.Part())
					}
				case prescribedMenu.Part() == nil:
					t.Errorf("PrescribedMenus()[%d].Part() = nil, want %v", i, *want.part)
				case *prescribedMenu.Part() != *want.part:
					t.Errorf("PrescribedMenus()[%d].Part() = %v, want %v", i, *prescribedMenu.Part(), *want.part)
				}
				if prescribedMenu.TrainingMenuID() != trainingMenu.ID() {
					t.Errorf("PrescribedMenus()[%d].TrainingMenuID() = %v, want %v", i, prescribedMenu.TrainingMenuID(), trainingMenu.ID())
				}
				if prescribedMenu.TrainingMenuName() != trainingMenu.Name() {
					t.Errorf("PrescribedMenus()[%d].TrainingMenuName() = %v, want %v", i, prescribedMenu.TrainingMenuName(), trainingMenu.Name())
				}
				if prescribedMenu.Amount().Int() != want.amount {
					t.Errorf("PrescribedMenus()[%d].Amount() = %v, want %v", i, prescribedMenu.Amount().Int(), want.amount)
				}
				if prescribedMenu.Unit() != want.unit {
					t.Errorf("PrescribedMenus()[%d].Unit() = %v, want %v", i, prescribedMenu.Unit(), want.unit)
				}
				if prescribedMenu.Sets().Int() != want.sets {
					t.Errorf("PrescribedMenus()[%d].Sets() = %v, want %v", i, prescribedMenu.Sets().Int(), want.sets)
				}
			}
		})
	}
}
