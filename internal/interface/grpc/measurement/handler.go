package measurement

import (
	"context"
	"errors"
	"fmt"
	"log"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	measurementv1 "github.com/qkitzero/fitness-service/gen/go/measurement/v1"
	appmeasurement "github.com/qkitzero/fitness-service/internal/application/measurement"
	domaincustomer "github.com/qkitzero/fitness-service/internal/domain/customer"
	domainmeasurement "github.com/qkitzero/fitness-service/internal/domain/measurement"
	domainmeasurementitem "github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	domainstaff "github.com/qkitzero/fitness-service/internal/domain/staff"
)

type MeasurementHandler struct {
	measurementv1.UnimplementedMeasurementServiceServer
	measurementUsecase appmeasurement.MeasurementUsecase
}

func NewMeasurementHandler(
	measurementUsecase appmeasurement.MeasurementUsecase,
) *MeasurementHandler {
	return &MeasurementHandler{
		measurementUsecase: measurementUsecase,
	}
}

func toDomainSide(s measurementv1.Side) (domainmeasurement.Side, error) {
	switch s {
	case measurementv1.Side_SIDE_NONE:
		return domainmeasurement.SideNone, nil
	case measurementv1.Side_SIDE_LEFT:
		return domainmeasurement.SideLeft, nil
	case measurementv1.Side_SIDE_RIGHT:
		return domainmeasurement.SideRight, nil
	default:
		return domainmeasurement.NewSide("")
	}
}

func toProtoSide(s domainmeasurement.Side) (measurementv1.Side, error) {
	switch s {
	case domainmeasurement.SideNone:
		return measurementv1.Side_SIDE_NONE, nil
	case domainmeasurement.SideLeft:
		return measurementv1.Side_SIDE_LEFT, nil
	case domainmeasurement.SideRight:
		return measurementv1.Side_SIDE_RIGHT, nil
	default:
		return measurementv1.Side_SIDE_UNSPECIFIED, fmt.Errorf("unmapped side %q", s)
	}
}

func toProtoMeasuredOn(m domainmeasurement.MeasuredOn) *date.Date {
	return &date.Date{
		Year:  int32(m.Year()),
		Month: int32(m.Month()),
		Day:   int32(m.Day()),
	}
}

func toProtoMeasurementValue(v domainmeasurement.MeasurementValue) (*measurementv1.MeasurementValue, error) {
	side, err := toProtoSide(v.Side())
	if err != nil {
		return nil, err
	}

	msg := &measurementv1.MeasurementValue{
		TrialIndex: uint32(v.TrialIndex().Int()),
		Side:       side,
	}
	if value := v.Value(); value != nil {
		f := value.Float64()
		msg.Value = &f
	}
	if valueSecondary := v.ValueSecondary(); valueSecondary != nil {
		f := valueSecondary.Float64()
		msg.ValueSecondary = &f
	}
	if valueChoice := v.ValueChoice(); valueChoice != nil {
		s := valueChoice.String()
		msg.ValueChoice = &s
	}

	return msg, nil
}

func toProtoMeasurementEntry(e domainmeasurement.MeasurementEntry) (*measurementv1.MeasurementEntry, error) {
	values := e.Values()
	valueMessages := make([]*measurementv1.MeasurementValue, 0, len(values))
	for _, v := range values {
		valueMessage, err := toProtoMeasurementValue(v)
		if err != nil {
			return nil, err
		}
		valueMessages = append(valueMessages, valueMessage)
	}

	msg := &measurementv1.MeasurementEntry{
		MeasurementItemId: e.MeasurementItemID().String(),
		Unmeasurable:      e.Unmeasurable(),
		Values:            valueMessages,
	}
	if note := e.Note(); note != nil {
		s := note.String()
		msg.Note = &s
	}

	return msg, nil
}

func toProtoMeasurement(m domainmeasurement.Measurement) (*measurementv1.Measurement, error) {
	entries := m.Entries()
	entryMessages := make([]*measurementv1.MeasurementEntry, 0, len(entries))
	for _, e := range entries {
		entryMessage, err := toProtoMeasurementEntry(e)
		if err != nil {
			return nil, err
		}
		entryMessages = append(entryMessages, entryMessage)
	}

	return &measurementv1.Measurement{
		MeasurementId:    m.ID().String(),
		CustomerId:       m.CustomerID().String(),
		MeasuredOn:       toProtoMeasuredOn(m.MeasuredOn()),
		MeasuredBy:       m.MeasuredBy().String(),
		AgeAtMeasurement: uint32(m.AgeAtMeasurement().Int()),
		UpdatedBy:        m.UpdatedBy().String(),
		IsDraft:          m.IsDraft(),
		Entries:          entryMessages,
	}, nil
}

func parseMeasurementValue(msg *measurementv1.MeasurementValue) (domainmeasurement.MeasurementValue, error) {
	trialIndex, err := domainmeasurement.NewTrialIndex(int(msg.GetTrialIndex()))
	if err != nil {
		return nil, err
	}
	side, err := toDomainSide(msg.GetSide())
	if err != nil {
		return nil, err
	}

	var value *domainmeasurement.Value
	if msg.Value != nil {
		v, err := domainmeasurement.NewValue(msg.GetValue())
		if err != nil {
			return nil, err
		}
		value = &v
	}

	var valueSecondary *domainmeasurement.Value
	if msg.ValueSecondary != nil {
		v, err := domainmeasurement.NewValue(msg.GetValueSecondary())
		if err != nil {
			return nil, err
		}
		valueSecondary = &v
	}

	var valueChoice *domainmeasurement.Choice
	if msg.ValueChoice != nil {
		c, err := domainmeasurement.NewChoice(msg.GetValueChoice())
		if err != nil {
			return nil, err
		}
		valueChoice = &c
	}

	return domainmeasurement.NewMeasurementValue(trialIndex, side, value, valueSecondary, valueChoice), nil
}

func parseEntryInput(msg *measurementv1.MeasurementEntry) (appmeasurement.MeasurementEntryInput, error) {
	var entryInput appmeasurement.MeasurementEntryInput

	measurementItemID, err := domainmeasurementitem.NewMeasurementItemIDFromString(msg.GetMeasurementItemId())
	if err != nil {
		return entryInput, err
	}
	note, err := domainmeasurement.NewNote(msg.GetNote())
	if err != nil {
		return entryInput, err
	}

	valueMessages := msg.GetValues()
	values := make([]domainmeasurement.MeasurementValue, 0, len(valueMessages))
	for _, valueMessage := range valueMessages {
		value, err := parseMeasurementValue(valueMessage)
		if err != nil {
			return entryInput, err
		}
		values = append(values, value)
	}

	entryInput.MeasurementItemID = measurementItemID
	entryInput.Unmeasurable = msg.GetUnmeasurable()
	entryInput.Note = note
	entryInput.Values = values

	return entryInput, nil
}

type measurementFieldsRequest interface {
	GetMeasuredOn() *date.Date
	GetMeasuredBy() string
	GetIsDraft() bool
	GetEntries() []*measurementv1.MeasurementEntry
}

type measurementFields struct {
	measuredOn  domainmeasurement.MeasuredOn
	measuredBy  domainstaff.StaffID
	isDraft     bool
	entryInputs []appmeasurement.MeasurementEntryInput
}

func parseMeasurementFields(req measurementFieldsRequest) (measurementFields, error) {
	var f measurementFields
	var err error
	if f.measuredOn, err = domainmeasurement.NewMeasuredOn(req.GetMeasuredOn().GetYear(), req.GetMeasuredOn().GetMonth(), req.GetMeasuredOn().GetDay()); err != nil {
		return f, err
	}
	if f.measuredBy, err = domainstaff.NewStaffID(req.GetMeasuredBy()); err != nil {
		return f, err
	}
	f.isDraft = req.GetIsDraft()

	entryMessages := req.GetEntries()
	f.entryInputs = make([]appmeasurement.MeasurementEntryInput, 0, len(entryMessages))
	for _, entryMessage := range entryMessages {
		entryInput, err := parseEntryInput(entryMessage)
		if err != nil {
			return f, err
		}
		f.entryInputs = append(f.entryInputs, entryInput)
	}

	return f, nil
}

func mapMeasurementError(err error, op string) error {
	if errors.Is(err, domainmeasurement.ErrMeasurementNotFound) || errors.Is(err, domaincustomer.ErrCustomerNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, domainmeasurementitem.ErrMeasurementItemNotFound) ||
		errors.Is(err, domainmeasurement.ErrInvalidValueCount) ||
		errors.Is(err, domainmeasurement.ErrInvalidValueForType) ||
		errors.Is(err, domainmeasurement.ErrInvalidSide) ||
		errors.Is(err, domainmeasurement.ErrDuplicateValue) ||
		errors.Is(err, domainmeasurement.ErrDuplicateEntry) ||
		errors.Is(err, domainmeasurement.ErrInvalidAgeAtMeasurement) {
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

func (h *MeasurementHandler) CreateMeasurement(ctx context.Context, req *measurementv1.CreateMeasurementRequest) (*measurementv1.CreateMeasurementResponse, error) {
	customerID, err := domaincustomer.NewCustomerIDFromString(req.GetCustomerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	fields, err := parseMeasurementFields(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	measurement, err := h.measurementUsecase.CreateMeasurement(ctx, customerID, fields.measuredOn, fields.measuredBy, fields.isDraft, fields.entryInputs)
	if err != nil {
		return nil, mapMeasurementError(err, "CreateMeasurement")
	}

	return &measurementv1.CreateMeasurementResponse{
		MeasurementId: measurement.ID().String(),
	}, nil
}

func (h *MeasurementHandler) GetMeasurement(ctx context.Context, req *measurementv1.GetMeasurementRequest) (*measurementv1.GetMeasurementResponse, error) {
	measurementID, err := domainmeasurement.NewMeasurementIDFromString(req.GetMeasurementId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	measurement, err := h.measurementUsecase.GetMeasurement(ctx, measurementID)
	if err != nil {
		return nil, mapMeasurementError(err, "GetMeasurement")
	}

	measurementMessage, err := toProtoMeasurement(measurement)
	if err != nil {
		return nil, mapMeasurementError(err, "GetMeasurement")
	}

	return &measurementv1.GetMeasurementResponse{
		Measurement: measurementMessage,
	}, nil
}

func (h *MeasurementHandler) ListMeasurements(ctx context.Context, req *measurementv1.ListMeasurementsRequest) (*measurementv1.ListMeasurementsResponse, error) {
	customerID, err := domaincustomer.NewCustomerIDFromString(req.GetCustomerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	measurements, err := h.measurementUsecase.ListMeasurements(ctx, customerID)
	if err != nil {
		return nil, mapMeasurementError(err, "ListMeasurements")
	}

	measurementMessages := make([]*measurementv1.Measurement, 0, len(measurements))
	for _, m := range measurements {
		measurementMessage, err := toProtoMeasurement(m)
		if err != nil {
			return nil, mapMeasurementError(err, "ListMeasurements")
		}
		measurementMessages = append(measurementMessages, measurementMessage)
	}

	return &measurementv1.ListMeasurementsResponse{
		Measurements: measurementMessages,
	}, nil
}

func (h *MeasurementHandler) UpdateMeasurement(ctx context.Context, req *measurementv1.UpdateMeasurementRequest) (*measurementv1.UpdateMeasurementResponse, error) {
	measurementID, err := domainmeasurement.NewMeasurementIDFromString(req.GetMeasurementId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	fields, err := parseMeasurementFields(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	measurement, err := h.measurementUsecase.UpdateMeasurement(ctx, measurementID, fields.measuredOn, fields.measuredBy, fields.isDraft, fields.entryInputs)
	if err != nil {
		return nil, mapMeasurementError(err, "UpdateMeasurement")
	}

	measurementMessage, err := toProtoMeasurement(measurement)
	if err != nil {
		return nil, mapMeasurementError(err, "UpdateMeasurement")
	}

	return &measurementv1.UpdateMeasurementResponse{
		Measurement: measurementMessage,
	}, nil
}

func (h *MeasurementHandler) DeleteMeasurement(ctx context.Context, req *measurementv1.DeleteMeasurementRequest) (*measurementv1.DeleteMeasurementResponse, error) {
	measurementID, err := domainmeasurement.NewMeasurementIDFromString(req.GetMeasurementId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.measurementUsecase.DeleteMeasurement(ctx, measurementID); err != nil {
		return nil, mapMeasurementError(err, "DeleteMeasurement")
	}

	return &measurementv1.DeleteMeasurementResponse{}, nil
}
