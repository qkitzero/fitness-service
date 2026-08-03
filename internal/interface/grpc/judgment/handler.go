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

	msg := &judgmentv1.Judgment{
		MeasurementId:      result.MeasurementID.String(),
		ItemEvaluations:    itemEvaluationMessages,
		ElementEvaluations: elementEvaluationMessages,
		IsDraft:            result.IsDraft,
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
