package training

import (
	"context"
)

type TrainingMenuRepository interface {
	List(ctx context.Context) ([]TrainingMenu, error)
	FindByIDs(ctx context.Context, trainingMenuIDs []TrainingMenuID) ([]TrainingMenu, error)
}
