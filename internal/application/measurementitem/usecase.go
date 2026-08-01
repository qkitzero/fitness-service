package measurementitem

import (
	"context"

	"github.com/qkitzero/fitness-service/internal/application/auth"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

type MeasurementItemUsecase interface {
	ListMeasurementItems(ctx context.Context) ([]measurementitem.MeasurementItem, error)
}

type measurementItemUsecase struct {
	authService         auth.AuthService
	measurementItemRepo measurementitem.MeasurementItemRepository
}

func NewMeasurementItemUsecase(authService auth.AuthService, measurementItemRepo measurementitem.MeasurementItemRepository) MeasurementItemUsecase {
	return &measurementItemUsecase{authService: authService, measurementItemRepo: measurementItemRepo}
}

func (u *measurementItemUsecase) ListMeasurementItems(ctx context.Context) ([]measurementitem.MeasurementItem, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	measurementItems, err := u.measurementItemRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	return measurementItems, nil
}
