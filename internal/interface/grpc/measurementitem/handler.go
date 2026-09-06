package measurementitem

import (
	"context"
	"fmt"
	"log"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	measurementitemv1 "github.com/qkitzero/fitness-service/gen/go/measurementitem/v1"
	appmeasurementitem "github.com/qkitzero/fitness-service/internal/application/measurementitem"
	domainmeasurementitem "github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

type MeasurementItemHandler struct {
	measurementitemv1.UnimplementedMeasurementItemServiceServer
	measurementItemUsecase appmeasurementitem.MeasurementItemUsecase
}

func NewMeasurementItemHandler(
	measurementItemUsecase appmeasurementitem.MeasurementItemUsecase,
) *MeasurementItemHandler {
	return &MeasurementItemHandler{
		measurementItemUsecase: measurementItemUsecase,
	}
}

func toProtoCategory(c domainmeasurementitem.Category) (measurementitemv1.Category, error) {
	switch c {
	case domainmeasurementitem.CategoryVital:
		return measurementitemv1.Category_CATEGORY_VITAL, nil
	case domainmeasurementitem.CategoryPhysique:
		return measurementitemv1.Category_CATEGORY_PHYSIQUE, nil
	case domainmeasurementitem.CategoryBodyComposition:
		return measurementitemv1.Category_CATEGORY_BODY_COMPOSITION, nil
	case domainmeasurementitem.CategoryMotorFunction:
		return measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION, nil
	default:
		return measurementitemv1.Category_CATEGORY_UNSPECIFIED, fmt.Errorf("unmapped category %q", c)
	}
}

func toProtoUnit(u domainmeasurementitem.Unit) (measurementitemv1.Unit, error) {
	switch u {
	case domainmeasurementitem.UnitKg:
		return measurementitemv1.Unit_UNIT_KG, nil
	case domainmeasurementitem.UnitCm:
		return measurementitemv1.Unit_UNIT_CM, nil
	case domainmeasurementitem.UnitSec:
		return measurementitemv1.Unit_UNIT_SEC, nil
	case domainmeasurementitem.UnitCount:
		return measurementitemv1.Unit_UNIT_COUNT, nil
	case domainmeasurementitem.UnitMmHg:
		return measurementitemv1.Unit_UNIT_MMHG, nil
	case domainmeasurementitem.UnitPercent:
		return measurementitemv1.Unit_UNIT_PERCENT, nil
	case domainmeasurementitem.UnitBpm:
		return measurementitemv1.Unit_UNIT_BPM, nil
	case domainmeasurementitem.UnitLevel:
		return measurementitemv1.Unit_UNIT_LEVEL, nil
	default:
		return measurementitemv1.Unit_UNIT_UNSPECIFIED, fmt.Errorf("unmapped unit %q", u)
	}
}

func toProtoSideMode(s domainmeasurementitem.SideMode) (measurementitemv1.SideMode, error) {
	switch s {
	case domainmeasurementitem.SideModeNone:
		return measurementitemv1.SideMode_SIDE_MODE_NONE, nil
	case domainmeasurementitem.SideModeBilateral:
		return measurementitemv1.SideMode_SIDE_MODE_BILATERAL, nil
	case domainmeasurementitem.SideModeOptionalBilateral:
		return measurementitemv1.SideMode_SIDE_MODE_OPTIONAL_BILATERAL, nil
	default:
		return measurementitemv1.SideMode_SIDE_MODE_UNSPECIFIED, fmt.Errorf("unmapped side mode %q", s)
	}
}

func toProtoValueType(v domainmeasurementitem.ValueType) (measurementitemv1.ValueType, error) {
	switch v {
	case domainmeasurementitem.ValueTypeNumeric:
		return measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, nil
	case domainmeasurementitem.ValueTypePaired:
		return measurementitemv1.ValueType_VALUE_TYPE_PAIRED, nil
	case domainmeasurementitem.ValueTypeChoice:
		return measurementitemv1.ValueType_VALUE_TYPE_CHOICE, nil
	default:
		return measurementitemv1.ValueType_VALUE_TYPE_UNSPECIFIED, fmt.Errorf("unmapped value type %q", v)
	}
}

func toProtoMeasurementItem(m domainmeasurementitem.MeasurementItem) (*measurementitemv1.MeasurementItem, error) {
	category, err := toProtoCategory(m.Category())
	if err != nil {
		return nil, err
	}
	unit, err := toProtoUnit(m.Unit())
	if err != nil {
		return nil, err
	}
	sideMode, err := toProtoSideMode(m.SideMode())
	if err != nil {
		return nil, err
	}
	valueType, err := toProtoValueType(m.ValueType())
	if err != nil {
		return nil, err
	}
	trialCount := m.TrialCount().Int()
	if trialCount < 0 || int64(trialCount) > math.MaxUint32 {
		return nil, fmt.Errorf("trial count %d out of range", trialCount)
	}

	return &measurementitemv1.MeasurementItem{
		MeasurementItemId: m.ID().String(),
		Code:              m.Code().String(),
		Name:              m.Name().String(),
		Category:          category,
		Unit:              unit,
		TrialCount:        uint32(trialCount),
		ValueType:         valueType,
		SideMode:          sideMode,
	}, nil
}

func mapMeasurementItemError(err error, op string) error {
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.Unauthenticated, codes.PermissionDenied:
			return err
		}
	}
	log.Printf("%s: internal error: %v", op, err)
	return status.Error(codes.Internal, "internal error")
}

func (h *MeasurementItemHandler) ListMeasurementItems(ctx context.Context, _ *measurementitemv1.ListMeasurementItemsRequest) (*measurementitemv1.ListMeasurementItemsResponse, error) {
	measurementItems, err := h.measurementItemUsecase.ListMeasurementItems(ctx)
	if err != nil {
		return nil, mapMeasurementItemError(err, "ListMeasurementItems")
	}

	measurementItemMessages := make([]*measurementitemv1.MeasurementItem, 0, len(measurementItems))
	for _, m := range measurementItems {
		measurementItemMessage, err := toProtoMeasurementItem(m)
		if err != nil {
			return nil, mapMeasurementItemError(err, "ListMeasurementItems")
		}
		measurementItemMessages = append(measurementItemMessages, measurementItemMessage)
	}

	return &measurementitemv1.ListMeasurementItemsResponse{
		MeasurementItems: measurementItemMessages,
	}, nil
}
