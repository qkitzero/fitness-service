package judgment

import (
	"context"
	"errors"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	judgmentv1 "github.com/qkitzero/fitness-service/gen/go/judgment/v1"
	appjudgment "github.com/qkitzero/fitness-service/internal/application/judgment"
	domaincustomer "github.com/qkitzero/fitness-service/internal/domain/customer"
	domainjudgment "github.com/qkitzero/fitness-service/internal/domain/judgment"
	domainmeasurement "github.com/qkitzero/fitness-service/internal/domain/measurement"
	domainmeasurementitem "github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	domainstandard "github.com/qkitzero/fitness-service/internal/domain/standard"
	domaintraining "github.com/qkitzero/fitness-service/internal/domain/training"
)

type JudgmentHandler struct {
	judgmentv1.UnimplementedJudgmentServiceServer
	judgmentUsecase appjudgment.JudgmentUsecase
}

func NewJudgmentHandler(
	judgmentUsecase appjudgment.JudgmentUsecase,
) *JudgmentHandler {
	return &JudgmentHandler{
		judgmentUsecase: judgmentUsecase,
	}
}

func toProtoRank(r domainstandard.Rank) (judgmentv1.Rank, error) {
	switch r {
	case domainstandard.RankA:
		return judgmentv1.Rank_RANK_A, nil
	case domainstandard.RankB:
		return judgmentv1.Rank_RANK_B, nil
	case domainstandard.RankC:
		return judgmentv1.Rank_RANK_C, nil
	case domainstandard.RankD:
		return judgmentv1.Rank_RANK_D, nil
	case domainstandard.RankE:
		return judgmentv1.Rank_RANK_E, nil
	default:
		return judgmentv1.Rank_RANK_UNSPECIFIED, fmt.Errorf("unmapped rank %q", r)
	}
}

func toProtoElement(e domainmeasurementitem.Element) (judgmentv1.Element, error) {
	switch e {
	case domainmeasurementitem.ElementMuscleStrength:
		return judgmentv1.Element_ELEMENT_MUSCLE_STRENGTH, nil
	case domainmeasurementitem.ElementMuscleEndurance:
		return judgmentv1.Element_ELEMENT_MUSCLE_ENDURANCE, nil
	case domainmeasurementitem.ElementFlexibility:
		return judgmentv1.Element_ELEMENT_FLEXIBILITY, nil
	case domainmeasurementitem.ElementAgility:
		return judgmentv1.Element_ELEMENT_AGILITY, nil
	case domainmeasurementitem.ElementBalance:
		return judgmentv1.Element_ELEMENT_BALANCE, nil
	case domainmeasurementitem.ElementMobility:
		return judgmentv1.Element_ELEMENT_MOBILITY, nil
	default:
		return judgmentv1.Element_ELEMENT_UNSPECIFIED, fmt.Errorf("unmapped element %q", e)
	}
}

func toDomainElement(e judgmentv1.Element) (domainmeasurementitem.Element, error) {
	switch e {
	case judgmentv1.Element_ELEMENT_MUSCLE_STRENGTH:
		return domainmeasurementitem.ElementMuscleStrength, nil
	case judgmentv1.Element_ELEMENT_MUSCLE_ENDURANCE:
		return domainmeasurementitem.ElementMuscleEndurance, nil
	case judgmentv1.Element_ELEMENT_FLEXIBILITY:
		return domainmeasurementitem.ElementFlexibility, nil
	case judgmentv1.Element_ELEMENT_AGILITY:
		return domainmeasurementitem.ElementAgility, nil
	case judgmentv1.Element_ELEMENT_BALANCE:
		return domainmeasurementitem.ElementBalance, nil
	case judgmentv1.Element_ELEMENT_MOBILITY:
		return domainmeasurementitem.ElementMobility, nil
	default:
		return domainmeasurementitem.Element(""), fmt.Errorf("unmapped element %q", e)
	}
}

func toDomainPart(p judgmentv1.PrescribedPart) (domaintraining.Part, error) {
	switch p {
	case judgmentv1.PrescribedPart_PRESCRIBED_PART_UPPER_LIMB:
		return domaintraining.PartUpperLimb, nil
	case judgmentv1.PrescribedPart_PRESCRIBED_PART_LOWER_LIMB:
		return domaintraining.PartLowerLimb, nil
	case judgmentv1.PrescribedPart_PRESCRIBED_PART_WHOLE_BODY:
		return domaintraining.PartWholeBody, nil
	default:
		return domaintraining.Part(""), fmt.Errorf("unmapped part %q", p)
	}
}

func toDomainUnit(u judgmentv1.PrescribedUnit) (domaintraining.Unit, error) {
	switch u {
	case judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_REPS:
		return domaintraining.UnitReps, nil
	case judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_SECONDS:
		return domaintraining.UnitSeconds, nil
	case judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_MINUTES:
		return domaintraining.UnitMinutes, nil
	default:
		return domaintraining.Unit(""), fmt.Errorf("unmapped unit %q", u)
	}
}

func toProtoPart(p domaintraining.Part) (judgmentv1.PrescribedPart, error) {
	switch p {
	case domaintraining.PartUpperLimb:
		return judgmentv1.PrescribedPart_PRESCRIBED_PART_UPPER_LIMB, nil
	case domaintraining.PartLowerLimb:
		return judgmentv1.PrescribedPart_PRESCRIBED_PART_LOWER_LIMB, nil
	case domaintraining.PartWholeBody:
		return judgmentv1.PrescribedPart_PRESCRIBED_PART_WHOLE_BODY, nil
	default:
		return judgmentv1.PrescribedPart_PRESCRIBED_PART_UNSPECIFIED, fmt.Errorf("unmapped part %q", p)
	}
}

func toProtoUnit(u domaintraining.Unit) (judgmentv1.PrescribedUnit, error) {
	switch u {
	case domaintraining.UnitReps:
		return judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_REPS, nil
	case domaintraining.UnitSeconds:
		return judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_SECONDS, nil
	case domaintraining.UnitMinutes:
		return judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_MINUTES, nil
	default:
		return judgmentv1.PrescribedUnit_PRESCRIBED_UNIT_UNSPECIFIED, fmt.Errorf("unmapped unit %q", u)
	}
}

func toProtoPrescriptionSource(s domainjudgment.PrescriptionSource) (judgmentv1.PrescriptionSource, error) {
	switch s {
	case domainjudgment.PrescriptionSourceElement:
		return judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_ELEMENT, nil
	case domainjudgment.PrescriptionSourceFixed:
		return judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_FIXED, nil
	case domainjudgment.PrescriptionSourceAgeDecade:
		return judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_AGE_DECADE, nil
	case domainjudgment.PrescriptionSourceManual:
		return judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_MANUAL, nil
	default:
		return judgmentv1.PrescriptionSource_PRESCRIPTION_SOURCE_UNSPECIFIED, fmt.Errorf("unmapped prescription source %q", s)
	}
}

func toProtoItemEvaluation(i domainjudgment.ItemEvaluation) (*judgmentv1.ItemEvaluation, error) {
	rank, err := toProtoRank(i.Rank())
	if err != nil {
		return nil, err
	}

	return &judgmentv1.ItemEvaluation{
		MeasurementItemId: i.MeasurementItemID().String(),
		Value:             i.Value().Float64(),
		Mean:              i.Mean().Float64(),
		ZScore:            i.ZScore().Float64(),
		Rank:              rank,
	}, nil
}

func toProtoElementEvaluation(e domainjudgment.ElementEvaluation) (*judgmentv1.ElementEvaluation, error) {
	element, err := toProtoElement(e.Element())
	if err != nil {
		return nil, err
	}
	rank, err := toProtoRank(e.Rank())
	if err != nil {
		return nil, err
	}

	return &judgmentv1.ElementEvaluation{
		Element: element,
		ZScore:  e.ZScore().Float64(),
		Rank:    rank,
	}, nil
}

func toProtoPrescribedMenu(p domainjudgment.PrescribedMenu) (*judgmentv1.PrescribedMenu, error) {
	source, err := toProtoPrescriptionSource(p.Source())
	if err != nil {
		return nil, err
	}
	unit, err := toProtoUnit(p.Unit())
	if err != nil {
		return nil, err
	}

	msg := &judgmentv1.PrescribedMenu{
		Source:           source,
		TrainingMenuId:   p.TrainingMenuID().String(),
		TrainingMenuName: p.TrainingMenuName().String(),
		Amount:           uint32(p.Amount().Int()),
		Unit:             unit,
		Sets:             uint32(p.Sets().Int()),
	}
	if e := p.Element(); e != nil {
		element, err := toProtoElement(*e)
		if err != nil {
			return nil, err
		}
		msg.Element = &element
	}
	if pt := p.Part(); pt != nil {
		part, err := toProtoPart(*pt)
		if err != nil {
			return nil, err
		}
		msg.Part = &part
	}

	return msg, nil
}

func toProtoPrescribedMenus(prescription domainjudgment.Prescription) ([]*judgmentv1.PrescribedMenu, error) {
	if prescription == nil {
		return nil, nil
	}

	prescribedMenus := prescription.PrescribedMenus()
	prescribedMenuMessages := make([]*judgmentv1.PrescribedMenu, 0, len(prescribedMenus))
	for _, p := range prescribedMenus {
		prescribedMenuMessage, err := toProtoPrescribedMenu(p)
		if err != nil {
			return nil, err
		}
		prescribedMenuMessages = append(prescribedMenuMessages, prescribedMenuMessage)
	}

	return prescribedMenuMessages, nil
}

func toProtoJudgment(result appjudgment.JudgmentResult) (*judgmentv1.Judgment, error) {
	itemEvaluations := result.Evaluation.ItemEvaluations()
	itemEvaluationMessages := make([]*judgmentv1.ItemEvaluation, 0, len(itemEvaluations))
	for _, i := range itemEvaluations {
		itemEvaluationMessage, err := toProtoItemEvaluation(i)
		if err != nil {
			return nil, err
		}
		itemEvaluationMessages = append(itemEvaluationMessages, itemEvaluationMessage)
	}

	elementEvaluations := result.Evaluation.ElementEvaluations()
	elementEvaluationMessages := make([]*judgmentv1.ElementEvaluation, 0, len(elementEvaluations))
	for _, e := range elementEvaluations {
		elementEvaluationMessage, err := toProtoElementEvaluation(e)
		if err != nil {
			return nil, err
		}
		elementEvaluationMessages = append(elementEvaluationMessages, elementEvaluationMessage)
	}

	prescribedMenuMessages, err := toProtoPrescribedMenus(result.Prescription)
	if err != nil {
		return nil, err
	}

	msg := &judgmentv1.Judgment{
		MeasurementId:      result.MeasurementID.String(),
		ItemEvaluations:    itemEvaluationMessages,
		ElementEvaluations: elementEvaluationMessages,
		IsDraft:            result.IsDraft,
		PrescribedMenus:    prescribedMenuMessages,
	}
	if motorAge := result.Evaluation.MotorAge(); motorAge != nil {
		n := uint32(motorAge.Int())
		msg.MotorAge = &n
	}
	if result.Advice != nil {
		s := result.Advice.String()
		msg.Advice = &s
	}

	return msg, nil
}

func mapJudgmentError(err error, op string) error {
	if errors.Is(err, domainjudgment.ErrJudgmentNotFound) ||
		errors.Is(err, domainmeasurement.ErrMeasurementNotFound) ||
		errors.Is(err, domaincustomer.ErrCustomerNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, domaintraining.ErrTrainingMenuNotFound) ||
		errors.Is(err, domaintraining.ErrInvalidSortOrder) ||
		errors.Is(err, domainjudgment.ErrEmptyPrescription) ||
		errors.Is(err, domainjudgment.ErrInvalidPrescribedMenuLabels) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.Unauthenticated, codes.PermissionDenied:
			return err
		}
	}
	log.Printf("%s: internal error: %v", op, err)
	return status.Error(codes.Internal, "internal error")
}

func (h *JudgmentHandler) GetJudgment(ctx context.Context, req *judgmentv1.GetJudgmentRequest) (*judgmentv1.GetJudgmentResponse, error) {
	measurementID, err := domainmeasurement.NewMeasurementIDFromString(req.GetMeasurementId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := h.judgmentUsecase.GetJudgment(ctx, measurementID)
	if err != nil {
		return nil, mapJudgmentError(err, "GetJudgment")
	}

	judgmentMessage, err := toProtoJudgment(result)
	if err != nil {
		return nil, mapJudgmentError(err, "GetJudgment")
	}

	return &judgmentv1.GetJudgmentResponse{
		Judgment: judgmentMessage,
	}, nil
}

func parseAdvicePatch(req *judgmentv1.UpsertJudgmentAdviceRequest) (appjudgment.AdvicePatch, error) {
	var patch appjudgment.AdvicePatch
	var err error
	if req.Advice != nil {
		patch.HasAdvice = true
		if patch.Advice, err = domainjudgment.NewAdvice(*req.Advice); err != nil {
			return patch, err
		}
	}
	return patch, nil
}

func (h *JudgmentHandler) UpsertJudgmentAdvice(ctx context.Context, req *judgmentv1.UpsertJudgmentAdviceRequest) (*judgmentv1.UpsertJudgmentAdviceResponse, error) {
	measurementID, err := domainmeasurement.NewMeasurementIDFromString(req.GetMeasurementId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	patch, err := parseAdvicePatch(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	upsertedJudgment, err := h.judgmentUsecase.UpsertJudgmentAdvice(ctx, measurementID, patch)
	if err != nil {
		return nil, mapJudgmentError(err, "UpsertJudgmentAdvice")
	}

	res := &judgmentv1.UpsertJudgmentAdviceResponse{
		MeasurementId: upsertedJudgment.MeasurementID().String(),
	}
	if upsertedAdvice := upsertedJudgment.Advice(); upsertedAdvice != nil {
		s := upsertedAdvice.String()
		res.Advice = &s
	}

	return res, nil
}

func parsePrescribedMenuInput(msg *judgmentv1.PrescribedMenuInput) (appjudgment.PrescribedMenuInput, error) {
	var input appjudgment.PrescribedMenuInput

	if msg.Element != nil {
		element, err := toDomainElement(msg.GetElement())
		if err != nil {
			return input, err
		}
		input.Element = &element
	}
	if msg.Part != nil {
		part, err := toDomainPart(msg.GetPart())
		if err != nil {
			return input, err
		}
		input.Part = &part
	}

	trainingMenuID, err := domaintraining.NewTrainingMenuIDFromString(msg.GetTrainingMenuId())
	if err != nil {
		return input, err
	}
	amount, err := domaintraining.NewAmount(int(msg.GetAmount()))
	if err != nil {
		return input, err
	}
	unit, err := toDomainUnit(msg.GetUnit())
	if err != nil {
		return input, err
	}
	sets, err := domaintraining.NewSets(int(msg.GetSets()))
	if err != nil {
		return input, err
	}

	input.TrainingMenuID = trainingMenuID
	input.Amount = amount
	input.Unit = unit
	input.Sets = sets

	return input, nil
}

func parsePrescribedMenuInputs(msgs []*judgmentv1.PrescribedMenuInput) ([]appjudgment.PrescribedMenuInput, error) {
	inputs := make([]appjudgment.PrescribedMenuInput, 0, len(msgs))
	for _, msg := range msgs {
		input, err := parsePrescribedMenuInput(msg)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, input)
	}
	return inputs, nil
}

func (h *JudgmentHandler) UpsertPrescription(ctx context.Context, req *judgmentv1.UpsertPrescriptionRequest) (*judgmentv1.UpsertPrescriptionResponse, error) {
	measurementID, err := domainmeasurement.NewMeasurementIDFromString(req.GetMeasurementId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	inputs, err := parsePrescribedMenuInputs(req.GetPrescribedMenus())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	upsertedPrescription, err := h.judgmentUsecase.UpsertPrescription(ctx, measurementID, inputs)
	if err != nil {
		return nil, mapJudgmentError(err, "UpsertPrescription")
	}

	prescribedMenuMessages, err := toProtoPrescribedMenus(upsertedPrescription)
	if err != nil {
		return nil, mapJudgmentError(err, "UpsertPrescription")
	}

	return &judgmentv1.UpsertPrescriptionResponse{
		MeasurementId:   measurementID.String(),
		PrescribedMenus: prescribedMenuMessages,
	}, nil
}

func (h *JudgmentHandler) DeletePrescription(ctx context.Context, req *judgmentv1.DeletePrescriptionRequest) (*judgmentv1.DeletePrescriptionResponse, error) {
	measurementID, err := domainmeasurement.NewMeasurementIDFromString(req.GetMeasurementId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.judgmentUsecase.DeletePrescription(ctx, measurementID); err != nil {
		return nil, mapJudgmentError(err, "DeletePrescription")
	}

	return &judgmentv1.DeletePrescriptionResponse{}, nil
}
