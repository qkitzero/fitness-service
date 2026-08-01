package measurement

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	measurementv1 "github.com/qkitzero/fitness-service/gen/go/measurement/v1"
	appmeasurement "github.com/qkitzero/fitness-service/internal/application/measurement"
	domaincustomer "github.com/qkitzero/fitness-service/internal/domain/customer"
	domainmeasurement "github.com/qkitzero/fitness-service/internal/domain/measurement"
	domainmeasurementitem "github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	domainstaff "github.com/qkitzero/fitness-service/internal/domain/staff"
	mocksappmeasurement "github.com/qkitzero/fitness-service/mocks/application/measurement"
)

const (
	sampleMeasurementID     = "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"
	sampleCustomerID        = "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"
	sampleMeasurementItemID = "3f2b6c1d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"
	sampleStaffID           = "google-oauth2|000000000000000000000"
)

func TestCreateMeasurement(t *testing.T) {
	t.Parallel()
	value := 32.4
	valueSecondary := 82.0
	negativeValue := -1.0
	choice := "片足20cm"
	emptyChoice := " "
	note := "ふらつきあり"
	tooLongNote := strings.Repeat("あ", 256)
	measuredOn := &date.Date{Year: 2026, Month: 8, Day: 1}
	futureMeasuredOn := &date.Date{Year: 2500, Month: 1, Day: 1}

	tests := []struct {
		name            string
		req             *measurementv1.CreateMeasurementRequest
		callUsecase     bool
		wantEntryInputs func() []appmeasurement.MeasurementEntryInput
		createErr       error
		wantCode        codes.Code
	}{
		{
			name: "success create measurement",
			req: &measurementv1.CreateMeasurementRequest{
				CustomerId: sampleCustomerID,
				MeasuredOn: measuredOn,
				MeasuredBy: sampleStaffID,
				Entries: []*measurementv1.MeasurementEntry{{
					MeasurementItemId: sampleMeasurementItemID,
					Note:              &note,
					Values: []*measurementv1.MeasurementValue{
						{TrialIndex: 1, Side: measurementv1.Side_SIDE_LEFT, Value: &value},
						{TrialIndex: 1, Side: measurementv1.Side_SIDE_RIGHT, Value: &value, ValueSecondary: &valueSecondary},
						{TrialIndex: 2, Side: measurementv1.Side_SIDE_NONE, ValueChoice: &choice},
					},
				}},
			},
			callUsecase: true,
			wantEntryInputs: func() []appmeasurement.MeasurementEntryInput {
				measurementItemID, _ := domainmeasurementitem.NewMeasurementItemIDFromString(sampleMeasurementItemID)
				domainNote, _ := domainmeasurement.NewNote("ふらつきあり")
				firstTrial, _ := domainmeasurement.NewTrialIndex(1)
				secondTrial, _ := domainmeasurement.NewTrialIndex(2)
				domainValue, _ := domainmeasurement.NewValue(32.4)
				domainValueSecondary, _ := domainmeasurement.NewValue(82)
				domainChoice, _ := domainmeasurement.NewChoice("片足20cm")
				return []appmeasurement.MeasurementEntryInput{{
					MeasurementItemID: measurementItemID,
					Note:              domainNote,
					Values: []domainmeasurement.MeasurementValue{
						domainmeasurement.NewMeasurementValue(firstTrial, domainmeasurement.SideLeft, &domainValue, nil, nil),
						domainmeasurement.NewMeasurementValue(firstTrial, domainmeasurement.SideRight, &domainValue, &domainValueSecondary, nil),
						domainmeasurement.NewMeasurementValue(secondTrial, domainmeasurement.SideNone, nil, nil, &domainChoice),
					},
				}}
			},
			wantCode: codes.OK,
		},
		{
			name: "success create draft measurement without entries",
			req: &measurementv1.CreateMeasurementRequest{
				CustomerId: sampleCustomerID,
				MeasuredOn: measuredOn,
				MeasuredBy: sampleStaffID,
				IsDraft:    true,
			},
			callUsecase:     true,
			wantEntryInputs: func() []appmeasurement.MeasurementEntryInput { return nil },
			wantCode:        codes.OK,
		},
		{
			name:     "failure invalid customer id",
			req:      &measurementv1.CreateMeasurementRequest{CustomerId: "", MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "failure future measured on",
			req:      &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: futureMeasuredOn, MeasuredBy: sampleStaffID},
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "failure missing measured on",
			req:      &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredBy: sampleStaffID},
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "failure invalid measured by",
			req:      &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: measuredOn, MeasuredBy: ""},
			wantCode: codes.InvalidArgument,
		},
		{
			name: "failure invalid measurement item id",
			req: &measurementv1.CreateMeasurementRequest{
				CustomerId: sampleCustomerID,
				MeasuredOn: measuredOn,
				MeasuredBy: sampleStaffID,
				Entries:    []*measurementv1.MeasurementEntry{{MeasurementItemId: "invalid"}},
			},
			wantCode: codes.InvalidArgument,
		},
		{
			name: "failure too long note",
			req: &measurementv1.CreateMeasurementRequest{
				CustomerId: sampleCustomerID,
				MeasuredOn: measuredOn,
				MeasuredBy: sampleStaffID,
				Entries:    []*measurementv1.MeasurementEntry{{MeasurementItemId: sampleMeasurementItemID, Note: &tooLongNote}},
			},
			wantCode: codes.InvalidArgument,
		},
		{
			name: "failure invalid trial index",
			req: &measurementv1.CreateMeasurementRequest{
				CustomerId: sampleCustomerID,
				MeasuredOn: measuredOn,
				MeasuredBy: sampleStaffID,
				Entries: []*measurementv1.MeasurementEntry{{
					MeasurementItemId: sampleMeasurementItemID,
					Values:            []*measurementv1.MeasurementValue{{TrialIndex: 0, Side: measurementv1.Side_SIDE_NONE, Value: &value}},
				}},
			},
			wantCode: codes.InvalidArgument,
		},
		{
			name: "failure unspecified side",
			req: &measurementv1.CreateMeasurementRequest{
				CustomerId: sampleCustomerID,
				MeasuredOn: measuredOn,
				MeasuredBy: sampleStaffID,
				Entries: []*measurementv1.MeasurementEntry{{
					MeasurementItemId: sampleMeasurementItemID,
					Values:            []*measurementv1.MeasurementValue{{TrialIndex: 1, Side: measurementv1.Side_SIDE_UNSPECIFIED, Value: &value}},
				}},
			},
			wantCode: codes.InvalidArgument,
		},
		{
			name: "failure negative value",
			req: &measurementv1.CreateMeasurementRequest{
				CustomerId: sampleCustomerID,
				MeasuredOn: measuredOn,
				MeasuredBy: sampleStaffID,
				Entries: []*measurementv1.MeasurementEntry{{
					MeasurementItemId: sampleMeasurementItemID,
					Values:            []*measurementv1.MeasurementValue{{TrialIndex: 1, Side: measurementv1.Side_SIDE_NONE, Value: &negativeValue}},
				}},
			},
			wantCode: codes.InvalidArgument,
		},
		{
			name: "failure negative secondary value",
			req: &measurementv1.CreateMeasurementRequest{
				CustomerId: sampleCustomerID,
				MeasuredOn: measuredOn,
				MeasuredBy: sampleStaffID,
				Entries: []*measurementv1.MeasurementEntry{{
					MeasurementItemId: sampleMeasurementItemID,
					Values:            []*measurementv1.MeasurementValue{{TrialIndex: 1, Side: measurementv1.Side_SIDE_NONE, Value: &value, ValueSecondary: &negativeValue}},
				}},
			},
			wantCode: codes.InvalidArgument,
		},
		{
			name: "failure blank choice",
			req: &measurementv1.CreateMeasurementRequest{
				CustomerId: sampleCustomerID,
				MeasuredOn: measuredOn,
				MeasuredBy: sampleStaffID,
				Entries: []*measurementv1.MeasurementEntry{{
					MeasurementItemId: sampleMeasurementItemID,
					Values:            []*measurementv1.MeasurementValue{{TrialIndex: 1, Side: measurementv1.Side_SIDE_NONE, ValueChoice: &emptyChoice}},
				}},
			},
			wantCode: codes.InvalidArgument,
		},
		{
			name:        "failure measurement item not found",
			req:         &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			createErr:   domainmeasurementitem.ErrMeasurementItemNotFound,
			wantCode:    codes.InvalidArgument,
		},
		{
			name:        "failure invalid value count",
			req:         &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			createErr:   domainmeasurement.ErrInvalidValueCount,
			wantCode:    codes.InvalidArgument,
		},
		{
			name:        "failure invalid value for type",
			req:         &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			createErr:   domainmeasurement.ErrInvalidValueForType,
			wantCode:    codes.InvalidArgument,
		},
		{
			name:        "failure invalid side",
			req:         &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			createErr:   domainmeasurement.ErrInvalidSide,
			wantCode:    codes.InvalidArgument,
		},
		{
			name:        "failure duplicate value",
			req:         &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			createErr:   domainmeasurement.ErrDuplicateValue,
			wantCode:    codes.InvalidArgument,
		},
		{
			name:        "failure duplicate entry",
			req:         &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			createErr:   domainmeasurement.ErrDuplicateEntry,
			wantCode:    codes.InvalidArgument,
		},
		{
			name:        "failure invalid age at measurement",
			req:         &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			createErr:   domainmeasurement.ErrInvalidAgeAtMeasurement,
			wantCode:    codes.InvalidArgument,
		},
		{
			name:        "failure customer not found",
			req:         &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			createErr:   domaincustomer.ErrCustomerNotFound,
			wantCode:    codes.NotFound,
		},
		{
			name:        "failure usecase error",
			req:         &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			createErr:   fmt.Errorf("create measurement error"),
			wantCode:    codes.Internal,
		},
		{
			name:        "failure unauthenticated is preserved",
			req:         &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			createErr:   status.Error(codes.Unauthenticated, "auth"),
			wantCode:    codes.Unauthenticated,
		},
		{
			name:        "failure downstream code is not forwarded",
			req:         &measurementv1.CreateMeasurementRequest{CustomerId: sampleCustomerID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			createErr:   status.Error(codes.NotFound, "user not found"),
			wantCode:    codes.Internal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			measurementID, _ := domainmeasurement.NewMeasurementIDFromString(sampleMeasurementID)
			customerID, _ := domaincustomer.NewCustomerIDFromString(sampleCustomerID)
			measuredBy, _ := domainstaff.NewStaffID(sampleStaffID)
			domainMeasuredOn, _ := domainmeasurement.NewMeasuredOn(2026, 8, 1)
			ageAtMeasurement, _ := domainmeasurement.NewAgeAtMeasurement(65)
			createdAt := time.Date(2026, 8, 1, 1, 2, 3, 0, time.UTC)
			createdMeasurement := domainmeasurement.ReconstructMeasurement(measurementID, customerID, domainMeasuredOn, measuredBy, ageAtMeasurement, measuredBy, false, nil, createdAt, createdAt)

			ctx := context.Background()
			mockUsecase := mocksappmeasurement.NewMockMeasurementUsecase(ctrl)
			if tt.callUsecase {
				var returned domainmeasurement.Measurement
				if tt.createErr == nil {
					returned = createdMeasurement
				}
				var wantEntryInputs any = gomock.Any()
				if tt.wantEntryInputs != nil {
					entryInputs := tt.wantEntryInputs()
					if entryInputs == nil {
						entryInputs = []appmeasurement.MeasurementEntryInput{}
					}
					wantEntryInputs = entryInputs
				}
				mockUsecase.EXPECT().CreateMeasurement(gomock.Any(), customerID, domainMeasuredOn, measuredBy, tt.req.GetIsDraft(), wantEntryInputs).Return(returned, tt.createErr).Times(1)
			}

			handler := NewMeasurementHandler(mockUsecase)

			res, err := handler.CreateMeasurement(ctx, tt.req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK && res.GetMeasurementId() != sampleMeasurementID {
				t.Errorf("MeasurementId = %v, want %v", res.GetMeasurementId(), sampleMeasurementID)
			}
		})
	}
}

func TestGetMeasurement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		measurementID string
		callUsecase   bool
		unmappedSide  bool
		getErr        error
		wantCode      codes.Code
	}{
		{"success get measurement", sampleMeasurementID, true, false, nil, codes.OK},
		{"failure invalid measurement id", "", false, false, nil, codes.InvalidArgument},
		{"failure measurement not found", sampleMeasurementID, true, false, domainmeasurement.ErrMeasurementNotFound, codes.NotFound},
		{"failure usecase error", sampleMeasurementID, true, false, fmt.Errorf("get measurement error"), codes.Internal},
		{"failure unauthenticated is preserved", sampleMeasurementID, true, false, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure downstream code is not forwarded", sampleMeasurementID, true, false, status.Error(codes.InvalidArgument, "user service"), codes.Internal},
		{"failure unmapped side in the response", sampleMeasurementID, true, true, nil, codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			measurementID, _ := domainmeasurement.NewMeasurementIDFromString(sampleMeasurementID)
			customerID, _ := domaincustomer.NewCustomerIDFromString(sampleCustomerID)
			measurementItemID, _ := domainmeasurementitem.NewMeasurementItemIDFromString(sampleMeasurementItemID)
			measuredBy, _ := domainstaff.NewStaffID(sampleStaffID)
			measuredOn, _ := domainmeasurement.NewMeasuredOn(2026, 8, 1)
			ageAtMeasurement, _ := domainmeasurement.NewAgeAtMeasurement(65)
			trialIndex, _ := domainmeasurement.NewTrialIndex(1)
			value, _ := domainmeasurement.NewValue(128)
			valueSecondary, _ := domainmeasurement.NewValue(82)
			note, _ := domainmeasurement.NewNote("ふらつきあり")
			createdAt := time.Date(2026, 8, 1, 1, 2, 3, 0, time.UTC)
			side := domainmeasurement.SideNone
			if tt.unmappedSide {
				side = domainmeasurement.Side("both")
			}
			entry := domainmeasurement.ReconstructMeasurementEntry(measurementItemID, false, note, []domainmeasurement.MeasurementValue{
				domainmeasurement.NewMeasurementValue(trialIndex, side, &value, &valueSecondary, nil),
			})
			foundMeasurement := domainmeasurement.ReconstructMeasurement(measurementID, customerID, measuredOn, measuredBy, ageAtMeasurement, measuredBy, true, []domainmeasurement.MeasurementEntry{entry}, createdAt, createdAt)

			ctx := context.Background()
			mockUsecase := mocksappmeasurement.NewMockMeasurementUsecase(ctrl)
			if tt.callUsecase {
				var returned domainmeasurement.Measurement
				if tt.getErr == nil {
					returned = foundMeasurement
				}
				mockUsecase.EXPECT().GetMeasurement(gomock.Any(), measurementID).Return(returned, tt.getErr).Times(1)
			}

			handler := NewMeasurementHandler(mockUsecase)

			res, err := handler.GetMeasurement(ctx, &measurementv1.GetMeasurementRequest{MeasurementId: tt.measurementID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode != codes.OK {
				return
			}

			got := res.GetMeasurement()
			if got.GetMeasurementId() != sampleMeasurementID {
				t.Errorf("MeasurementId = %v, want %v", got.GetMeasurementId(), sampleMeasurementID)
			}
			if got.GetCustomerId() != sampleCustomerID {
				t.Errorf("CustomerId = %v, want %v", got.GetCustomerId(), sampleCustomerID)
			}
			if got.GetMeasuredOn().GetYear() != 2026 || got.GetMeasuredOn().GetMonth() != 8 || got.GetMeasuredOn().GetDay() != 1 {
				t.Errorf("MeasuredOn = %v, want 2026-08-01", got.GetMeasuredOn())
			}
			if got.GetMeasuredBy() != sampleStaffID {
				t.Errorf("MeasuredBy = %v, want %v", got.GetMeasuredBy(), sampleStaffID)
			}
			if got.GetUpdatedBy() != sampleStaffID {
				t.Errorf("UpdatedBy = %v, want %v", got.GetUpdatedBy(), sampleStaffID)
			}
			if got.GetAgeAtMeasurement() != 65 {
				t.Errorf("AgeAtMeasurement = %v, want %v", got.GetAgeAtMeasurement(), 65)
			}
			if !got.GetIsDraft() {
				t.Errorf("IsDraft = %v, want %v", got.GetIsDraft(), true)
			}
			if len(got.GetEntries()) != 1 {
				t.Fatalf("len(Entries) = %v, want %v", len(got.GetEntries()), 1)
			}
			gotEntry := got.GetEntries()[0]
			if gotEntry.GetMeasurementItemId() != sampleMeasurementItemID {
				t.Errorf("MeasurementItemId = %v, want %v", gotEntry.GetMeasurementItemId(), sampleMeasurementItemID)
			}
			if gotEntry.GetUnmeasurable() {
				t.Errorf("Unmeasurable = %v, want %v", gotEntry.GetUnmeasurable(), false)
			}
			if gotEntry.GetNote() != "ふらつきあり" {
				t.Errorf("Note = %v, want %v", gotEntry.GetNote(), "ふらつきあり")
			}
			if len(gotEntry.GetValues()) != 1 {
				t.Fatalf("len(Values) = %v, want %v", len(gotEntry.GetValues()), 1)
			}
			gotValue := gotEntry.GetValues()[0]
			if gotValue.GetTrialIndex() != 1 {
				t.Errorf("TrialIndex = %v, want %v", gotValue.GetTrialIndex(), 1)
			}
			if gotValue.GetSide() != measurementv1.Side_SIDE_NONE {
				t.Errorf("Side = %v, want %v", gotValue.GetSide(), measurementv1.Side_SIDE_NONE)
			}
			if gotValue.GetValue() != 128 {
				t.Errorf("Value = %v, want %v", gotValue.GetValue(), 128)
			}
			if gotValue.GetValueSecondary() != 82 {
				t.Errorf("ValueSecondary = %v, want %v", gotValue.GetValueSecondary(), 82)
			}
			if gotValue.ValueChoice != nil {
				t.Errorf("ValueChoice = %v, want nil", gotValue.GetValueChoice())
			}
		})
	}
}

func TestListMeasurements(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		customerID       string
		callUsecase      bool
		measurementCount int
		unmappedSide     bool
		listErr          error
		wantCode         codes.Code
	}{
		{"success list measurements", sampleCustomerID, true, 2, false, nil, codes.OK},
		{"success list no measurements", sampleCustomerID, true, 0, false, nil, codes.OK},
		{"failure invalid customer id", "", false, 0, false, nil, codes.InvalidArgument},
		{"failure customer not found", sampleCustomerID, true, 0, false, domaincustomer.ErrCustomerNotFound, codes.NotFound},
		{"failure usecase error", sampleCustomerID, true, 0, false, fmt.Errorf("list measurements error"), codes.Internal},
		{"failure unauthenticated is preserved", sampleCustomerID, true, 0, false, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
		{"failure unmapped side in the response", sampleCustomerID, true, 1, true, nil, codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			customerID, _ := domaincustomer.NewCustomerIDFromString(sampleCustomerID)
			measuredBy, _ := domainstaff.NewStaffID(sampleStaffID)
			measuredOn, _ := domainmeasurement.NewMeasuredOn(2026, 8, 1)
			ageAtMeasurement, _ := domainmeasurement.NewAgeAtMeasurement(65)
			createdAt := time.Date(2026, 8, 1, 1, 2, 3, 0, time.UTC)
			var entries []domainmeasurement.MeasurementEntry
			if tt.unmappedSide {
				measurementItemID, _ := domainmeasurementitem.NewMeasurementItemIDFromString(sampleMeasurementItemID)
				trialIndex, _ := domainmeasurement.NewTrialIndex(1)
				value, _ := domainmeasurement.NewValue(72)
				entries = []domainmeasurement.MeasurementEntry{
					domainmeasurement.ReconstructMeasurementEntry(measurementItemID, false, nil, []domainmeasurement.MeasurementValue{
						domainmeasurement.NewMeasurementValue(trialIndex, domainmeasurement.Side("both"), &value, nil, nil),
					}),
				}
			}
			measurements := make([]domainmeasurement.Measurement, 0, tt.measurementCount)
			for range tt.measurementCount {
				measurements = append(measurements, domainmeasurement.ReconstructMeasurement(domainmeasurement.NewMeasurementID(), customerID, measuredOn, measuredBy, ageAtMeasurement, measuredBy, false, entries, createdAt, createdAt))
			}

			ctx := context.Background()
			mockUsecase := mocksappmeasurement.NewMockMeasurementUsecase(ctrl)
			if tt.callUsecase {
				if tt.listErr != nil {
					measurements = nil
				}
				mockUsecase.EXPECT().ListMeasurements(gomock.Any(), customerID).Return(measurements, tt.listErr).Times(1)
			}

			handler := NewMeasurementHandler(mockUsecase)

			res, err := handler.ListMeasurements(ctx, &measurementv1.ListMeasurementsRequest{CustomerId: tt.customerID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK && len(res.GetMeasurements()) != tt.measurementCount {
				t.Errorf("len(Measurements) = %v, want %v", len(res.GetMeasurements()), tt.measurementCount)
			}
		})
	}
}

func TestUpdateMeasurement(t *testing.T) {
	t.Parallel()
	value := 72.0
	measuredOn := &date.Date{Year: 2026, Month: 8, Day: 1}
	futureMeasuredOn := &date.Date{Year: 2500, Month: 1, Day: 1}

	tests := []struct {
		name            string
		req             *measurementv1.UpdateMeasurementRequest
		callUsecase     bool
		unmappedSide    bool
		wantEntryInputs func() []appmeasurement.MeasurementEntryInput
		updateErr       error
		wantCode        codes.Code
	}{
		{
			name: "success update measurement",
			req: &measurementv1.UpdateMeasurementRequest{
				MeasurementId: sampleMeasurementID,
				MeasuredOn:    measuredOn,
				MeasuredBy:    sampleStaffID,
				Entries: []*measurementv1.MeasurementEntry{{
					MeasurementItemId: sampleMeasurementItemID,
					Values:            []*measurementv1.MeasurementValue{{TrialIndex: 1, Side: measurementv1.Side_SIDE_NONE, Value: &value}},
				}},
			},
			callUsecase: true,
			wantEntryInputs: func() []appmeasurement.MeasurementEntryInput {
				measurementItemID, _ := domainmeasurementitem.NewMeasurementItemIDFromString(sampleMeasurementItemID)
				trialIndex, _ := domainmeasurement.NewTrialIndex(1)
				v, _ := domainmeasurement.NewValue(72)
				return []appmeasurement.MeasurementEntryInput{{
					MeasurementItemID: measurementItemID,
					Values: []domainmeasurement.MeasurementValue{
						domainmeasurement.NewMeasurementValue(trialIndex, domainmeasurement.SideNone, &v, nil, nil),
					},
				}}
			},
			wantCode: codes.OK,
		},
		{
			name:     "failure invalid measurement id",
			req:      &measurementv1.UpdateMeasurementRequest{MeasurementId: "", MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "failure future measured on",
			req:      &measurementv1.UpdateMeasurementRequest{MeasurementId: sampleMeasurementID, MeasuredOn: futureMeasuredOn, MeasuredBy: sampleStaffID},
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "failure invalid measured by",
			req:      &measurementv1.UpdateMeasurementRequest{MeasurementId: sampleMeasurementID, MeasuredOn: measuredOn, MeasuredBy: ""},
			wantCode: codes.InvalidArgument,
		},
		{
			name:        "failure measurement not found",
			req:         &measurementv1.UpdateMeasurementRequest{MeasurementId: sampleMeasurementID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			updateErr:   domainmeasurement.ErrMeasurementNotFound,
			wantCode:    codes.NotFound,
		},
		{
			name:        "failure invalid value count",
			req:         &measurementv1.UpdateMeasurementRequest{MeasurementId: sampleMeasurementID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			updateErr:   domainmeasurement.ErrInvalidValueCount,
			wantCode:    codes.InvalidArgument,
		},
		{
			name:        "failure usecase error",
			req:         &measurementv1.UpdateMeasurementRequest{MeasurementId: sampleMeasurementID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			updateErr:   fmt.Errorf("update measurement error"),
			wantCode:    codes.Internal,
		},
		{
			name:        "failure unauthenticated is preserved",
			req:         &measurementv1.UpdateMeasurementRequest{MeasurementId: sampleMeasurementID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase: true,
			updateErr:   status.Error(codes.Unauthenticated, "auth"),
			wantCode:    codes.Unauthenticated,
		},
		{
			name:         "failure unmapped side in the response",
			req:          &measurementv1.UpdateMeasurementRequest{MeasurementId: sampleMeasurementID, MeasuredOn: measuredOn, MeasuredBy: sampleStaffID},
			callUsecase:  true,
			unmappedSide: true,
			wantCode:     codes.Internal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			measurementID, _ := domainmeasurement.NewMeasurementIDFromString(sampleMeasurementID)
			customerID, _ := domaincustomer.NewCustomerIDFromString(sampleCustomerID)
			measurementItemID, _ := domainmeasurementitem.NewMeasurementItemIDFromString(sampleMeasurementItemID)
			measuredBy, _ := domainstaff.NewStaffID(sampleStaffID)
			domainMeasuredOn, _ := domainmeasurement.NewMeasuredOn(2026, 8, 1)
			ageAtMeasurement, _ := domainmeasurement.NewAgeAtMeasurement(65)
			trialIndex, _ := domainmeasurement.NewTrialIndex(1)
			domainValue, _ := domainmeasurement.NewValue(72)
			createdAt := time.Date(2026, 8, 1, 1, 2, 3, 0, time.UTC)
			side := domainmeasurement.SideNone
			if tt.unmappedSide {
				side = domainmeasurement.Side("both")
			}
			entry := domainmeasurement.ReconstructMeasurementEntry(measurementItemID, false, nil, []domainmeasurement.MeasurementValue{
				domainmeasurement.NewMeasurementValue(trialIndex, side, &domainValue, nil, nil),
			})
			updatedMeasurement := domainmeasurement.ReconstructMeasurement(measurementID, customerID, domainMeasuredOn, measuredBy, ageAtMeasurement, measuredBy, false, []domainmeasurement.MeasurementEntry{entry}, createdAt, createdAt)

			ctx := context.Background()
			mockUsecase := mocksappmeasurement.NewMockMeasurementUsecase(ctrl)
			if tt.callUsecase {
				var returned domainmeasurement.Measurement
				if tt.updateErr == nil {
					returned = updatedMeasurement
				}
				var wantEntryInputs any = gomock.Any()
				if tt.wantEntryInputs != nil {
					wantEntryInputs = tt.wantEntryInputs()
				}
				mockUsecase.EXPECT().UpdateMeasurement(gomock.Any(), measurementID, domainMeasuredOn, measuredBy, tt.req.GetIsDraft(), wantEntryInputs).Return(returned, tt.updateErr).Times(1)
			}

			handler := NewMeasurementHandler(mockUsecase)

			res, err := handler.UpdateMeasurement(ctx, tt.req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				if res.GetMeasurement().GetMeasurementId() != sampleMeasurementID {
					t.Errorf("MeasurementId = %v, want %v", res.GetMeasurement().GetMeasurementId(), sampleMeasurementID)
				}
				if len(res.GetMeasurement().GetEntries()) != 1 {
					t.Errorf("len(Entries) = %v, want %v", len(res.GetMeasurement().GetEntries()), 1)
				}
			}
		})
	}
}

func TestDeleteMeasurement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		measurementID string
		callUsecase   bool
		deleteErr     error
		wantCode      codes.Code
	}{
		{"success delete measurement", sampleMeasurementID, true, nil, codes.OK},
		{"failure invalid measurement id", "", false, nil, codes.InvalidArgument},
		{"failure measurement not found", sampleMeasurementID, true, domainmeasurement.ErrMeasurementNotFound, codes.NotFound},
		{"failure usecase error", sampleMeasurementID, true, fmt.Errorf("delete measurement error"), codes.Internal},
		{"failure unauthenticated is preserved", sampleMeasurementID, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			measurementID, _ := domainmeasurement.NewMeasurementIDFromString(sampleMeasurementID)

			ctx := context.Background()
			mockUsecase := mocksappmeasurement.NewMockMeasurementUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().DeleteMeasurement(gomock.Any(), measurementID).Return(tt.deleteErr).Times(1)
			}

			handler := NewMeasurementHandler(mockUsecase)

			_, err := handler.DeleteMeasurement(ctx, &measurementv1.DeleteMeasurementRequest{MeasurementId: tt.measurementID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestToDomainSide(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		side    measurementv1.Side
		want    domainmeasurement.Side
	}{
		{"success none", true, measurementv1.Side_SIDE_NONE, domainmeasurement.SideNone},
		{"success left", true, measurementv1.Side_SIDE_LEFT, domainmeasurement.SideLeft},
		{"success right", true, measurementv1.Side_SIDE_RIGHT, domainmeasurement.SideRight},
		{"failure unspecified", false, measurementv1.Side_SIDE_UNSPECIFIED, domainmeasurement.Side("")},
		{"failure unknown", false, measurementv1.Side(99), domainmeasurement.Side("")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			side, err := toDomainSide(tt.side)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if side != tt.want {
				t.Errorf("toDomainSide(%v) = %v, want %v", tt.side, side, tt.want)
			}
		})
	}
}

func TestToProtoSide(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		side    domainmeasurement.Side
		want    measurementv1.Side
	}{
		{"success none", true, domainmeasurement.SideNone, measurementv1.Side_SIDE_NONE},
		{"success left", true, domainmeasurement.SideLeft, measurementv1.Side_SIDE_LEFT},
		{"success right", true, domainmeasurement.SideRight, measurementv1.Side_SIDE_RIGHT},
		{"failure unmapped side", false, domainmeasurement.Side("both"), measurementv1.Side_SIDE_UNSPECIFIED},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			side, err := toProtoSide(tt.side)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if side != tt.want {
				t.Errorf("toProtoSide(%v) = %v, want %v", tt.side, side, tt.want)
			}
		})
	}
}

func TestToProtoMeasurement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		success          bool
		ageAtMeasurement domainmeasurement.AgeAtMeasurement
		side             domainmeasurement.Side
		trialIndex       domainmeasurement.TrialIndex
	}{
		{"success", true, domainmeasurement.AgeAtMeasurement(65), domainmeasurement.SideLeft, domainmeasurement.TrialIndex(1)},
		{"failure unmapped side", false, domainmeasurement.AgeAtMeasurement(65), domainmeasurement.Side("both"), domainmeasurement.TrialIndex(1)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			measurementID, _ := domainmeasurement.NewMeasurementIDFromString(sampleMeasurementID)
			customerID, _ := domaincustomer.NewCustomerIDFromString(sampleCustomerID)
			measurementItemID, _ := domainmeasurementitem.NewMeasurementItemIDFromString(sampleMeasurementItemID)
			measuredBy, _ := domainstaff.NewStaffID(sampleStaffID)
			measuredOn, _ := domainmeasurement.NewMeasuredOn(2026, 8, 1)
			value, _ := domainmeasurement.NewValue(32.4)
			createdAt := time.Date(2026, 8, 1, 1, 2, 3, 0, time.UTC)
			choice, _ := domainmeasurement.NewChoice("片足20cm")
			entry := domainmeasurement.ReconstructMeasurementEntry(measurementItemID, false, nil, []domainmeasurement.MeasurementValue{
				domainmeasurement.NewMeasurementValue(tt.trialIndex, tt.side, &value, nil, nil),
				domainmeasurement.NewMeasurementValue(tt.trialIndex, tt.side, nil, nil, &choice),
			})
			m := domainmeasurement.ReconstructMeasurement(measurementID, customerID, measuredOn, measuredBy, tt.ageAtMeasurement, measuredBy, false, []domainmeasurement.MeasurementEntry{entry}, createdAt, createdAt)

			msg, err := toProtoMeasurement(m)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && msg.GetMeasurementId() != sampleMeasurementID {
				t.Errorf("MeasurementId = %v, want %v", msg.GetMeasurementId(), sampleMeasurementID)
			}
			if !tt.success && msg != nil {
				t.Errorf("expected no message on failure, but got %v", msg)
			}
		})
	}
}
