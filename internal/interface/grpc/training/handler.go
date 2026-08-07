package training

import (
	"context"
	"fmt"
	"log"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	trainingv1 "github.com/qkitzero/fitness-service/gen/go/training/v1"
	apptraining "github.com/qkitzero/fitness-service/internal/application/training"
	domainmeasurementitem "github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	domaintraining "github.com/qkitzero/fitness-service/internal/domain/training"
)

type TrainingMenuHandler struct {
	trainingv1.UnimplementedTrainingMenuServiceServer
	trainingMenuUsecase apptraining.TrainingMenuUsecase
}

func NewTrainingMenuHandler(
	trainingMenuUsecase apptraining.TrainingMenuUsecase,
) *TrainingMenuHandler {
	return &TrainingMenuHandler{
		trainingMenuUsecase: trainingMenuUsecase,
	}
}

func toProtoElement(e domainmeasurementitem.Element) (trainingv1.TrainingElement, error) {
	switch e {
	case domainmeasurementitem.ElementMuscleStrength:
		return trainingv1.TrainingElement_TRAINING_ELEMENT_MUSCLE_STRENGTH, nil
	case domainmeasurementitem.ElementMuscleEndurance:
		return trainingv1.TrainingElement_TRAINING_ELEMENT_MUSCLE_ENDURANCE, nil
	case domainmeasurementitem.ElementFlexibility:
		return trainingv1.TrainingElement_TRAINING_ELEMENT_FLEXIBILITY, nil
	case domainmeasurementitem.ElementAgility:
		return trainingv1.TrainingElement_TRAINING_ELEMENT_AGILITY, nil
	case domainmeasurementitem.ElementBalance:
		return trainingv1.TrainingElement_TRAINING_ELEMENT_BALANCE, nil
	case domainmeasurementitem.ElementMobility:
		return trainingv1.TrainingElement_TRAINING_ELEMENT_MOBILITY, nil
	default:
		return trainingv1.TrainingElement_TRAINING_ELEMENT_UNSPECIFIED, fmt.Errorf("unmapped element %q", e)
	}
}

func toProtoPart(p domaintraining.Part) (trainingv1.Part, error) {
	switch p {
	case domaintraining.PartUpperLimb:
		return trainingv1.Part_PART_UPPER_LIMB, nil
	case domaintraining.PartLowerLimb:
		return trainingv1.Part_PART_LOWER_LIMB, nil
	case domaintraining.PartWholeBody:
		return trainingv1.Part_PART_WHOLE_BODY, nil
	default:
		return trainingv1.Part_PART_UNSPECIFIED, fmt.Errorf("unmapped part %q", p)
	}
}

func toProtoUnit(u domaintraining.Unit) (trainingv1.TrainingUnit, error) {
	switch u {
	case domaintraining.UnitReps:
		return trainingv1.TrainingUnit_TRAINING_UNIT_REPS, nil
	case domaintraining.UnitSeconds:
		return trainingv1.TrainingUnit_TRAINING_UNIT_SECONDS, nil
	case domaintraining.UnitMinutes:
		return trainingv1.TrainingUnit_TRAINING_UNIT_MINUTES, nil
	default:
		return trainingv1.TrainingUnit_TRAINING_UNIT_UNSPECIFIED, fmt.Errorf("unmapped unit %q", u)
	}
}

func toProtoTrainingMenu(t domaintraining.TrainingMenu) (*trainingv1.TrainingMenu, error) {
	element, err := toProtoElement(t.Element())
	if err != nil {
		return nil, err
	}
	part, err := toProtoPart(t.Part())
	if err != nil {
		return nil, err
	}
	unit, err := toProtoUnit(t.Unit())
	if err != nil {
		return nil, err
	}
	amount := t.Amount().Int()
	if amount < 0 || int64(amount) > math.MaxUint32 {
		return nil, fmt.Errorf("amount %d out of range", amount)
	}
	sets := t.Sets().Int()
	if sets < 0 || int64(sets) > math.MaxUint32 {
		return nil, fmt.Errorf("sets %d out of range", sets)
	}

	return &trainingv1.TrainingMenu{
		TrainingMenuId: t.ID().String(),
		Code:           t.Code().String(),
		Name:           t.Name().String(),
		Element:        element,
		Part:           part,
		Amount:         uint32(amount),
		Unit:           unit,
		Sets:           uint32(sets),
		Instruction:    t.Instruction().String(),
	}, nil
}

func mapTrainingMenuError(err error, op string) error {
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.Unauthenticated, codes.PermissionDenied:
			return err
		}
	}
	log.Printf("%s: internal error: %v", op, err)
	return status.Error(codes.Internal, "internal error")
}

func (h *TrainingMenuHandler) ListTrainingMenus(ctx context.Context, _ *trainingv1.ListTrainingMenusRequest) (*trainingv1.ListTrainingMenusResponse, error) {
	trainingMenus, err := h.trainingMenuUsecase.ListTrainingMenus(ctx)
	if err != nil {
		return nil, mapTrainingMenuError(err, "ListTrainingMenus")
	}

	trainingMenuMessages := make([]*trainingv1.TrainingMenu, 0, len(trainingMenus))
	for _, t := range trainingMenus {
		trainingMenuMessage, err := toProtoTrainingMenu(t)
		if err != nil {
			return nil, mapTrainingMenuError(err, "ListTrainingMenus")
		}
		trainingMenuMessages = append(trainingMenuMessages, trainingMenuMessage)
	}

	return &trainingv1.ListTrainingMenusResponse{
		TrainingMenus: trainingMenuMessages,
	}, nil
}
