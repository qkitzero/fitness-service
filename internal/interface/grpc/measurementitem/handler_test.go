package measurementitem

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	measurementitemv1 "github.com/qkitzero/fitness-service/gen/go/measurementitem/v1"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	mocksappmeasurementitem "github.com/qkitzero/fitness-service/mocks/application/measurementitem"
	mocksmeasurementitem "github.com/qkitzero/fitness-service/mocks/domain/measurementitem"
)

func TestListMeasurementItems(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		codes          []string
		names          []string
		categories     []measurementitem.Category
		units          []measurementitem.Unit
		trialCounts    []int
		sideModes      []measurementitem.SideMode
		valueTypes     []measurementitem.ValueType
		listErr        error
		wantCode       codes.Code
		wantCategories []measurementitemv1.Category
		wantUnits      []measurementitemv1.Unit
		wantValueTypes []measurementitemv1.ValueType
		wantSideModes  []measurementitemv1.SideMode
		wantBilaterals []bool
	}{
		{
			name:           "success list measurement items keeps every item in order",
			codes:          []string{"blood_pressure", "height", "body_fat_percentage", "grip_strength"},
			names:          []string{"血圧", "身長", "体脂肪率", "握力"},
			categories:     []measurementitem.Category{measurementitem.CategoryVital, measurementitem.CategoryPhysique, measurementitem.CategoryBodyComposition, measurementitem.CategoryMotorFunction},
			units:          []measurementitem.Unit{measurementitem.UnitMmHg, measurementitem.UnitCm, measurementitem.UnitPercent, measurementitem.UnitKg},
			trialCounts:    []int{1, 1, 1, 2},
			sideModes:      []measurementitem.SideMode{measurementitem.SideModeNone, measurementitem.SideModeNone, measurementitem.SideModeNone, measurementitem.SideModeBilateral},
			valueTypes:     []measurementitem.ValueType{measurementitem.ValueTypePaired, measurementitem.ValueTypeNumeric, measurementitem.ValueTypeNumeric, measurementitem.ValueTypeNumeric},
			wantCode:       codes.OK,
			wantCategories: []measurementitemv1.Category{measurementitemv1.Category_CATEGORY_VITAL, measurementitemv1.Category_CATEGORY_PHYSIQUE, measurementitemv1.Category_CATEGORY_BODY_COMPOSITION, measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION},
			wantUnits:      []measurementitemv1.Unit{measurementitemv1.Unit_UNIT_MMHG, measurementitemv1.Unit_UNIT_CM, measurementitemv1.Unit_UNIT_PERCENT, measurementitemv1.Unit_UNIT_KG},
			wantValueTypes: []measurementitemv1.ValueType{measurementitemv1.ValueType_VALUE_TYPE_PAIRED, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC},
			wantSideModes:  []measurementitemv1.SideMode{measurementitemv1.SideMode_SIDE_MODE_NONE, measurementitemv1.SideMode_SIDE_MODE_NONE, measurementitemv1.SideMode_SIDE_MODE_NONE, measurementitemv1.SideMode_SIDE_MODE_BILATERAL},
			wantBilaterals: []bool{false, false, false, true},
		},
		{
			name:           "success maps the remaining units and value types",
			codes:          []string{"pulse_rate", "eyes_closed_one_leg_stand", "cs30", "stand_up_test", "choice_item"},
			names:          []string{"脈拍", "閉眼片足立ち", "CS-30（30秒立ち座り）", "立ち上がり", "選択式項目"},
			categories:     []measurementitem.Category{measurementitem.CategoryVital, measurementitem.CategoryMotorFunction, measurementitem.CategoryMotorFunction, measurementitem.CategoryMotorFunction, measurementitem.CategoryMotorFunction},
			units:          []measurementitem.Unit{measurementitem.UnitBpm, measurementitem.UnitSec, measurementitem.UnitCount, measurementitem.UnitLevel, measurementitem.UnitCm},
			trialCounts:    []int{1, 2, 1, 1, 1},
			sideModes:      []measurementitem.SideMode{measurementitem.SideModeNone, measurementitem.SideModeBilateral, measurementitem.SideModeNone, measurementitem.SideModeOptionalBilateral, measurementitem.SideModeNone},
			valueTypes:     []measurementitem.ValueType{measurementitem.ValueTypeNumeric, measurementitem.ValueTypeNumeric, measurementitem.ValueTypeNumeric, measurementitem.ValueTypeNumeric, measurementitem.ValueTypeChoice},
			wantCode:       codes.OK,
			wantCategories: []measurementitemv1.Category{measurementitemv1.Category_CATEGORY_VITAL, measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION, measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION, measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION, measurementitemv1.Category_CATEGORY_MOTOR_FUNCTION},
			wantUnits:      []measurementitemv1.Unit{measurementitemv1.Unit_UNIT_BPM, measurementitemv1.Unit_UNIT_SEC, measurementitemv1.Unit_UNIT_COUNT, measurementitemv1.Unit_UNIT_LEVEL, measurementitemv1.Unit_UNIT_CM},
			wantValueTypes: []measurementitemv1.ValueType{measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_NUMERIC, measurementitemv1.ValueType_VALUE_TYPE_CHOICE},
			wantSideModes:  []measurementitemv1.SideMode{measurementitemv1.SideMode_SIDE_MODE_NONE, measurementitemv1.SideMode_SIDE_MODE_BILATERAL, measurementitemv1.SideMode_SIDE_MODE_NONE, measurementitemv1.SideMode_SIDE_MODE_OPTIONAL_BILATERAL, measurementitemv1.SideMode_SIDE_MODE_NONE},
			wantBilaterals: []bool{false, true, false, false, false},
		},
		{
			name:     "success list no measurement items",
			wantCode: codes.OK,
		},
		{
			name:     "failure usecase error",
			listErr:  fmt.Errorf("list measurement items error"),
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
			name:        "failure unmapped category",
			codes:       []string{"grip_strength"},
			names:       []string{"握力"},
			categories:  []measurementitem.Category{measurementitem.Category("flexibility")},
			units:       []measurementitem.Unit{measurementitem.UnitKg},
			trialCounts: []int{2},
			sideModes:   []measurementitem.SideMode{measurementitem.SideModeBilateral},
			valueTypes:  []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			wantCode:    codes.Internal,
		},
		{
			name:        "failure unmapped unit",
			codes:       []string{"grip_strength"},
			names:       []string{"握力"},
			categories:  []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:       []measurementitem.Unit{measurementitem.Unit("newton")},
			trialCounts: []int{2},
			sideModes:   []measurementitem.SideMode{measurementitem.SideModeBilateral},
			valueTypes:  []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			wantCode:    codes.Internal,
		},
		{
			name:        "failure unmapped side mode",
			codes:       []string{"grip_strength"},
			names:       []string{"握力"},
			categories:  []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:       []measurementitem.Unit{measurementitem.UnitKg},
			trialCounts: []int{2},
			sideModes:   []measurementitem.SideMode{measurementitem.SideMode("unilateral")},
			valueTypes:  []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			wantCode:    codes.Internal,
		},
		{
			name:        "failure unmapped value type",
			codes:       []string{"grip_strength"},
			names:       []string{"握力"},
			categories:  []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:       []measurementitem.Unit{measurementitem.UnitKg},
			trialCounts: []int{2},
			sideModes:   []measurementitem.SideMode{measurementitem.SideModeBilateral},
			valueTypes:  []measurementitem.ValueType{measurementitem.ValueType("range")},
			wantCode:    codes.Internal,
		},
		{
			name:        "failure trial count out of proto range",
			codes:       []string{"grip_strength"},
			names:       []string{"握力"},
			categories:  []measurementitem.Category{measurementitem.CategoryMotorFunction},
			units:       []measurementitem.Unit{measurementitem.UnitKg},
			trialCounts: []int{-1},
			sideModes:   []measurementitem.SideMode{measurementitem.SideModeBilateral},
			valueTypes:  []measurementitem.ValueType{measurementitem.ValueTypeNumeric},
			wantCode:    codes.Internal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			wantIDs := make([]measurementitem.MeasurementItemID, 0, len(tt.codes))
			measurementItems := make([]measurementitem.MeasurementItem, 0, len(tt.codes))
			for i := range tt.codes {
				id := measurementitem.NewMeasurementItemID()
				wantIDs = append(wantIDs, id)
				mockMeasurementItem := mocksmeasurementitem.NewMockMeasurementItem(ctrl)
				mockMeasurementItem.EXPECT().ID().Return(id).AnyTimes()
				mockMeasurementItem.EXPECT().Code().Return(measurementitem.Code(tt.codes[i])).AnyTimes()
				mockMeasurementItem.EXPECT().Name().Return(measurementitem.Name(tt.names[i])).AnyTimes()
				mockMeasurementItem.EXPECT().Category().Return(tt.categories[i]).AnyTimes()
				mockMeasurementItem.EXPECT().Unit().Return(tt.units[i]).AnyTimes()
				mockMeasurementItem.EXPECT().TrialCount().Return(measurementitem.TrialCount(tt.trialCounts[i])).AnyTimes()
				mockMeasurementItem.EXPECT().SideMode().Return(tt.sideModes[i]).AnyTimes()
				mockMeasurementItem.EXPECT().ValueType().Return(tt.valueTypes[i]).AnyTimes()
				measurementItems = append(measurementItems, mockMeasurementItem)
			}

			mockUsecase := mocksappmeasurementitem.NewMockMeasurementItemUsecase(ctrl)
			if tt.listErr != nil {
				mockUsecase.EXPECT().ListMeasurementItems(gomock.Any()).Return(nil, tt.listErr).Times(1)
			} else {
				mockUsecase.EXPECT().ListMeasurementItems(gomock.Any()).Return(measurementItems, nil).Times(1)
			}

			handler := NewMeasurementItemHandler(mockUsecase)

			res, err := handler.ListMeasurementItems(ctx, &measurementitemv1.ListMeasurementItemsRequest{})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode != codes.OK {
				return
			}
			if len(res.GetMeasurementItems()) != len(tt.codes) {
				t.Fatalf("len(MeasurementItems) = %v, want %v", len(res.GetMeasurementItems()), len(tt.codes))
			}
			for i, msg := range res.GetMeasurementItems() {
				if msg.GetMeasurementItemId() != wantIDs[i].String() {
					t.Errorf("MeasurementItems[%d].MeasurementItemId = %v, want %v", i, msg.GetMeasurementItemId(), wantIDs[i].String())
				}
				if msg.GetCode() != tt.codes[i] {
					t.Errorf("MeasurementItems[%d].Code = %v, want %v", i, msg.GetCode(), tt.codes[i])
				}
				if msg.GetName() != tt.names[i] {
					t.Errorf("MeasurementItems[%d].Name = %v, want %v", i, msg.GetName(), tt.names[i])
				}
				if msg.GetCategory() != tt.wantCategories[i] {
					t.Errorf("MeasurementItems[%d].Category = %v, want %v", i, msg.GetCategory(), tt.wantCategories[i])
				}
				if msg.GetUnit() != tt.wantUnits[i] {
					t.Errorf("MeasurementItems[%d].Unit = %v, want %v", i, msg.GetUnit(), tt.wantUnits[i])
				}
				if msg.GetTrialCount() != uint32(tt.trialCounts[i]) {
					t.Errorf("MeasurementItems[%d].TrialCount = %v, want %v", i, msg.GetTrialCount(), tt.trialCounts[i])
				}
				if msg.GetSideMode() != tt.wantSideModes[i] {
					t.Errorf("MeasurementItems[%d].SideMode = %v, want %v", i, msg.GetSideMode(), tt.wantSideModes[i])
				}
				if msg.GetBilateral() != tt.wantBilaterals[i] {
					t.Errorf("MeasurementItems[%d].Bilateral = %v, want %v", i, msg.GetBilateral(), tt.wantBilaterals[i])
				}
				if msg.GetValueType() != tt.wantValueTypes[i] {
					t.Errorf("MeasurementItems[%d].ValueType = %v, want %v", i, msg.GetValueType(), tt.wantValueTypes[i])
				}
			}
		})
	}
}
