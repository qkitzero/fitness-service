package training

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/fitness-service/internal/domain/training"
	mocksappauth "github.com/qkitzero/fitness-service/mocks/application/auth"
	mockstraining "github.com/qkitzero/fitness-service/mocks/domain/training"
)

func TestListTrainingMenus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		callList       bool
		ctx            context.Context
		userID         string
		verifyTokenErr error
		menuCount      int
		listErr        error
	}{
		{"success list training menus", true, true, context.Background(), "google-oauth2|000000000000000000000", nil, 3, nil},
		{"success list no training menus", true, true, context.Background(), "google-oauth2|000000000000000000000", nil, 0, nil},
		{"failure verify token error", false, false, context.Background(), "", errors.New("verify token error"), 0, nil},
		{"failure list error", false, true, context.Background(), "google-oauth2|000000000000000000000", nil, 0, errors.New("list error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wantTrainingMenus := make([]training.TrainingMenu, 0, tt.menuCount)
			for range tt.menuCount {
				wantTrainingMenus = append(wantTrainingMenus, mockstraining.NewMockTrainingMenu(ctrl))
			}

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockTrainingMenuRepository := mockstraining.NewMockTrainingMenuRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			if tt.callList {
				if tt.listErr != nil {
					mockTrainingMenuRepository.EXPECT().List(tt.ctx).Return(nil, tt.listErr).Times(1)
				} else {
					mockTrainingMenuRepository.EXPECT().List(tt.ctx).Return(wantTrainingMenus, nil).Times(1)
				}
			}

			u := NewTrainingMenuUsecase(mockAuthService, mockTrainingMenuRepository)

			trainingMenus, err := u.ListTrainingMenus(tt.ctx)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && trainingMenus != nil {
				t.Errorf("expected no training menus on failure, but got %v", len(trainingMenus))
			}
			if tt.success {
				if len(trainingMenus) != tt.menuCount {
					t.Errorf("len(trainingMenus) = %v, want %v", len(trainingMenus), tt.menuCount)
				}
				for i := range wantTrainingMenus {
					if i >= len(trainingMenus) {
						break
					}
					if trainingMenus[i] != wantTrainingMenus[i] {
						t.Errorf("trainingMenus[%d] is not the menu returned by the repository", i)
					}
				}
			}
		})
	}
}
