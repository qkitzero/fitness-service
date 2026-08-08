package training

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	trainingv1 "github.com/qkitzero/fitness-service/gen/go/training/v1"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/training"
	mocksapptraining "github.com/qkitzero/fitness-service/mocks/application/training"
	mockstraining "github.com/qkitzero/fitness-service/mocks/domain/training"
)

func TestListTrainingMenus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		codes        []string
		names        []string
		elements     []measurementitem.Element
		parts        []training.Part
		amounts      []int
		units        []training.Unit
		sets         []int
		instructions []string
		listErr      error
		wantCode     codes.Code
		wantElements []trainingv1.TrainingElement
		wantParts    []trainingv1.Part
		wantUnits    []trainingv1.TrainingUnit
	}{
		{
			name:         "success list training menus keeps every menu in order",
			codes:        []string{"wall_push_up", "squat", "front_plank", "shoulder_rotation"},
			names:        []string{"壁押し", "スクワット", "フロントプランク", "肩回し"},
			elements:     []measurementitem.Element{measurementitem.ElementMuscleStrength, measurementitem.ElementMuscleStrength, measurementitem.ElementMuscleStrength, measurementitem.ElementFlexibility},
			parts:        []training.Part{training.PartUpperLimb, training.PartLowerLimb, training.PartWholeBody, training.PartUpperLimb},
			amounts:      []int{10, 10, 30, 10},
			units:        []training.Unit{training.UnitReps, training.UnitReps, training.UnitSeconds, training.UnitReps},
			sets:         []int{3, 3, 3, 2},
			instructions: []string{"壁に手をついて肘を曲げ伸ばしする。", "腰を落として立ち上がる。", "体を一直線に保つ。", "肘で大きな円を描く。"},
			wantCode:     codes.OK,
			wantElements: []trainingv1.TrainingElement{trainingv1.TrainingElement_TRAINING_ELEMENT_MUSCLE_STRENGTH, trainingv1.TrainingElement_TRAINING_ELEMENT_MUSCLE_STRENGTH, trainingv1.TrainingElement_TRAINING_ELEMENT_MUSCLE_STRENGTH, trainingv1.TrainingElement_TRAINING_ELEMENT_FLEXIBILITY},
			wantParts:    []trainingv1.Part{trainingv1.Part_PART_UPPER_LIMB, trainingv1.Part_PART_LOWER_LIMB, trainingv1.Part_PART_WHOLE_BODY, trainingv1.Part_PART_UPPER_LIMB},
			wantUnits:    []trainingv1.TrainingUnit{trainingv1.TrainingUnit_TRAINING_UNIT_REPS, trainingv1.TrainingUnit_TRAINING_UNIT_REPS, trainingv1.TrainingUnit_TRAINING_UNIT_SECONDS, trainingv1.TrainingUnit_TRAINING_UNIT_REPS},
		},
		{
			name:         "success maps the remaining elements and units",
			codes:        []string{"seated_march", "side_step_touch", "one_leg_stand", "indoor_walk"},
			names:        []string{"座って足踏み", "横ステップタッチ", "片足立ち", "室内歩行"},
			elements:     []measurementitem.Element{measurementitem.ElementMuscleEndurance, measurementitem.ElementAgility, measurementitem.ElementBalance, measurementitem.ElementMobility},
			parts:        []training.Part{training.PartLowerLimb, training.PartLowerLimb, training.PartLowerLimb, training.PartWholeBody},
			amounts:      []int{20, 20, 30, 5},
			units:        []training.Unit{training.UnitReps, training.UnitReps, training.UnitSeconds, training.UnitMinutes},
			sets:         []int{2, 3, 3, 1},
			instructions: []string{"左右の膝を交互に持ち上げる。", "左右に一歩ずつステップする。", "片足を上げて姿勢を保つ。", "無理のない速さで歩く。"},
			wantCode:     codes.OK,
			wantElements: []trainingv1.TrainingElement{trainingv1.TrainingElement_TRAINING_ELEMENT_MUSCLE_ENDURANCE, trainingv1.TrainingElement_TRAINING_ELEMENT_AGILITY, trainingv1.TrainingElement_TRAINING_ELEMENT_BALANCE, trainingv1.TrainingElement_TRAINING_ELEMENT_MOBILITY},
			wantParts:    []trainingv1.Part{trainingv1.Part_PART_LOWER_LIMB, trainingv1.Part_PART_LOWER_LIMB, trainingv1.Part_PART_LOWER_LIMB, trainingv1.Part_PART_WHOLE_BODY},
			wantUnits:    []trainingv1.TrainingUnit{trainingv1.TrainingUnit_TRAINING_UNIT_REPS, trainingv1.TrainingUnit_TRAINING_UNIT_REPS, trainingv1.TrainingUnit_TRAINING_UNIT_SECONDS, trainingv1.TrainingUnit_TRAINING_UNIT_MINUTES},
		},
		{
			name:     "success list no training menus",
			wantCode: codes.OK,
		},
		{
			name:     "failure usecase error",
			listErr:  fmt.Errorf("list training menus error"),
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
			name:         "failure unmapped element",
			codes:        []string{"squat"},
			names:        []string{"スクワット"},
			elements:     []measurementitem.Element{measurementitem.Element("explosive_power")},
			parts:        []training.Part{training.PartLowerLimb},
			amounts:      []int{10},
			units:        []training.Unit{training.UnitReps},
			sets:         []int{3},
			instructions: []string{"腰を落として立ち上がる。"},
			wantCode:     codes.Internal,
		},
		{
			name:         "failure unmapped part",
			codes:        []string{"squat"},
			names:        []string{"スクワット"},
			elements:     []measurementitem.Element{measurementitem.ElementMuscleStrength},
			parts:        []training.Part{training.Part("trunk")},
			amounts:      []int{10},
			units:        []training.Unit{training.UnitReps},
			sets:         []int{3},
			instructions: []string{"腰を落として立ち上がる。"},
			wantCode:     codes.Internal,
		},
		{
			name:         "failure unmapped unit",
			codes:        []string{"squat"},
			names:        []string{"スクワット"},
			elements:     []measurementitem.Element{measurementitem.ElementMuscleStrength},
			parts:        []training.Part{training.PartLowerLimb},
			amounts:      []int{10},
			units:        []training.Unit{training.Unit("hours")},
			sets:         []int{3},
			instructions: []string{"腰を落として立ち上がる。"},
			wantCode:     codes.Internal,
		},
		{
			name:         "failure amount out of proto range",
			codes:        []string{"squat"},
			names:        []string{"スクワット"},
			elements:     []measurementitem.Element{measurementitem.ElementMuscleStrength},
			parts:        []training.Part{training.PartLowerLimb},
			amounts:      []int{-1},
			units:        []training.Unit{training.UnitReps},
			sets:         []int{3},
			instructions: []string{"腰を落として立ち上がる。"},
			wantCode:     codes.Internal,
		},
		{
			name:         "failure sets out of proto range",
			codes:        []string{"squat"},
			names:        []string{"スクワット"},
			elements:     []measurementitem.Element{measurementitem.ElementMuscleStrength},
			parts:        []training.Part{training.PartLowerLimb},
			amounts:      []int{10},
			units:        []training.Unit{training.UnitReps},
			sets:         []int{-1},
			instructions: []string{"腰を落として立ち上がる。"},
			wantCode:     codes.Internal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			wantIDs := make([]training.TrainingMenuID, 0, len(tt.codes))
			trainingMenus := make([]training.TrainingMenu, 0, len(tt.codes))
			for i := range tt.codes {
				id := training.NewTrainingMenuID()
				wantIDs = append(wantIDs, id)
				mockTrainingMenu := mockstraining.NewMockTrainingMenu(ctrl)
				mockTrainingMenu.EXPECT().ID().Return(id).AnyTimes()
				mockTrainingMenu.EXPECT().Code().Return(training.Code(tt.codes[i])).AnyTimes()
				mockTrainingMenu.EXPECT().Name().Return(training.Name(tt.names[i])).AnyTimes()
				mockTrainingMenu.EXPECT().Element().Return(tt.elements[i]).AnyTimes()
				mockTrainingMenu.EXPECT().Part().Return(tt.parts[i]).AnyTimes()
				mockTrainingMenu.EXPECT().Amount().Return(training.Amount(tt.amounts[i])).AnyTimes()
				mockTrainingMenu.EXPECT().Unit().Return(tt.units[i]).AnyTimes()
				mockTrainingMenu.EXPECT().Sets().Return(training.Sets(tt.sets[i])).AnyTimes()
				mockTrainingMenu.EXPECT().Instruction().Return(training.Instruction(tt.instructions[i])).AnyTimes()
				trainingMenus = append(trainingMenus, mockTrainingMenu)
			}

			mockUsecase := mocksapptraining.NewMockTrainingMenuUsecase(ctrl)
			if tt.listErr != nil {
				mockUsecase.EXPECT().ListTrainingMenus(gomock.Any()).Return(nil, tt.listErr).Times(1)
			} else {
				mockUsecase.EXPECT().ListTrainingMenus(gomock.Any()).Return(trainingMenus, nil).Times(1)
			}

			handler := NewTrainingMenuHandler(mockUsecase)

			res, err := handler.ListTrainingMenus(ctx, &trainingv1.ListTrainingMenusRequest{})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode != codes.OK {
				return
			}
			if len(res.GetTrainingMenus()) != len(tt.codes) {
				t.Fatalf("len(TrainingMenus) = %v, want %v", len(res.GetTrainingMenus()), len(tt.codes))
			}
			for i, msg := range res.GetTrainingMenus() {
				if msg.GetTrainingMenuId() != wantIDs[i].String() {
					t.Errorf("TrainingMenus[%d].TrainingMenuId = %v, want %v", i, msg.GetTrainingMenuId(), wantIDs[i].String())
				}
				if msg.GetCode() != tt.codes[i] {
					t.Errorf("TrainingMenus[%d].Code = %v, want %v", i, msg.GetCode(), tt.codes[i])
				}
				if msg.GetName() != tt.names[i] {
					t.Errorf("TrainingMenus[%d].Name = %v, want %v", i, msg.GetName(), tt.names[i])
				}
				if msg.GetElement() != tt.wantElements[i] {
					t.Errorf("TrainingMenus[%d].Element = %v, want %v", i, msg.GetElement(), tt.wantElements[i])
				}
				if msg.GetPart() != tt.wantParts[i] {
					t.Errorf("TrainingMenus[%d].Part = %v, want %v", i, msg.GetPart(), tt.wantParts[i])
				}
				if msg.GetAmount() != uint32(tt.amounts[i]) {
					t.Errorf("TrainingMenus[%d].Amount = %v, want %v", i, msg.GetAmount(), tt.amounts[i])
				}
				if msg.GetUnit() != tt.wantUnits[i] {
					t.Errorf("TrainingMenus[%d].Unit = %v, want %v", i, msg.GetUnit(), tt.wantUnits[i])
				}
				if msg.GetSets() != uint32(tt.sets[i]) {
					t.Errorf("TrainingMenus[%d].Sets = %v, want %v", i, msg.GetSets(), tt.sets[i])
				}
				if msg.GetInstruction() != tt.instructions[i] {
					t.Errorf("TrainingMenus[%d].Instruction = %v, want %v", i, msg.GetInstruction(), tt.instructions[i])
				}
			}
		})
	}
}
