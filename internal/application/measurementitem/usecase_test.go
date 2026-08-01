package measurementitem

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	mocksappauth "github.com/qkitzero/fitness-service/mocks/application/auth"
	mocksmeasurementitem "github.com/qkitzero/fitness-service/mocks/domain/measurementitem"
)

func TestListMeasurementItems(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		callList       bool
		ctx            context.Context
		userID         string
		verifyTokenErr error
		itemCount      int
		listErr        error
	}{
		{"success list measurement items", true, true, context.Background(), "google-oauth2|000000000000000000000", nil, 3, nil},
		{"success list no measurement items", true, true, context.Background(), "google-oauth2|000000000000000000000", nil, 0, nil},
		{"failure verify token error", false, false, context.Background(), "", errors.New("verify token error"), 0, nil},
		{"failure list error", false, true, context.Background(), "google-oauth2|000000000000000000000", nil, 0, errors.New("list error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wantMeasurementItems := make([]measurementitem.MeasurementItem, 0, tt.itemCount)
			for range tt.itemCount {
				wantMeasurementItems = append(wantMeasurementItems, mocksmeasurementitem.NewMockMeasurementItem(ctrl))
			}

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockMeasurementItemRepository := mocksmeasurementitem.NewMockMeasurementItemRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			if tt.callList {
				if tt.listErr != nil {
					mockMeasurementItemRepository.EXPECT().List(tt.ctx).Return(nil, tt.listErr).Times(1)
				} else {
					mockMeasurementItemRepository.EXPECT().List(tt.ctx).Return(wantMeasurementItems, nil).Times(1)
				}
			}

			u := NewMeasurementItemUsecase(mockAuthService, mockMeasurementItemRepository)

			measurementItems, err := u.ListMeasurementItems(tt.ctx)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && measurementItems != nil {
				t.Errorf("expected no measurement items on failure, but got %v", len(measurementItems))
			}
			if tt.success {
				if len(measurementItems) != tt.itemCount {
					t.Errorf("len(measurementItems) = %v, want %v", len(measurementItems), tt.itemCount)
				}
				for i := range wantMeasurementItems {
					if i >= len(measurementItems) {
						break
					}
					if measurementItems[i] != wantMeasurementItems[i] {
						t.Errorf("measurementItems[%d] is not the item returned by the repository", i)
					}
				}
			}
		})
	}
}
