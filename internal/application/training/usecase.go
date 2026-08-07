package training

import (
	"context"

	"github.com/qkitzero/fitness-service/internal/application/auth"
	"github.com/qkitzero/fitness-service/internal/domain/training"
)

type TrainingMenuUsecase interface {
	ListTrainingMenus(ctx context.Context) ([]training.TrainingMenu, error)
}

type trainingMenuUsecase struct {
	authService      auth.AuthService
	trainingMenuRepo training.TrainingMenuRepository
}

func NewTrainingMenuUsecase(authService auth.AuthService, trainingMenuRepo training.TrainingMenuRepository) TrainingMenuUsecase {
	return &trainingMenuUsecase{authService: authService, trainingMenuRepo: trainingMenuRepo}
}

func (u *trainingMenuUsecase) ListTrainingMenus(ctx context.Context) ([]training.TrainingMenu, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	trainingMenus, err := u.trainingMenuRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	return trainingMenus, nil
}
