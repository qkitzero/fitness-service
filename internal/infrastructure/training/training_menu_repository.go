package training

import (
	"context"
	"fmt"
	"sort"

	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/training"
)

type trainingMenuRepository struct {
	db *gorm.DB
}

func NewTrainingMenuRepository(db *gorm.DB) training.TrainingMenuRepository {
	return &trainingMenuRepository{db: db}
}

func toTrainingMenu(m TrainingMenuModel) (training.TrainingMenu, error) {
	code, err := training.NewCode(m.Code.String())
	if err != nil {
		return nil, fmt.Errorf("training menu %q: %w", m.Code, err)
	}
	name, err := training.NewName(m.Name.String())
	if err != nil {
		return nil, fmt.Errorf("training menu %q: %w", m.Code, err)
	}
	element, err := measurementitem.NewElement(m.Element.String())
	if err != nil {
		return nil, fmt.Errorf("training menu %q: %w", m.Code, err)
	}
	part, err := training.NewPart(m.Part.String())
	if err != nil {
		return nil, fmt.Errorf("training menu %q: %w", m.Code, err)
	}
	amount, err := training.NewAmount(m.Amount.Int())
	if err != nil {
		return nil, fmt.Errorf("training menu %q: %w", m.Code, err)
	}
	unit, err := training.NewUnit(m.Unit.String())
	if err != nil {
		return nil, fmt.Errorf("training menu %q: %w", m.Code, err)
	}
	sets, err := training.NewSets(m.Sets.Int())
	if err != nil {
		return nil, fmt.Errorf("training menu %q: %w", m.Code, err)
	}
	instruction, err := training.NewInstruction(m.Instruction.String())
	if err != nil {
		return nil, fmt.Errorf("training menu %q: %w", m.Code, err)
	}

	return training.NewTrainingMenu(
		m.ID,
		code,
		name,
		element,
		part,
		amount,
		unit,
		sets,
		instruction,
		m.CreatedAt,
		m.UpdatedAt,
	), nil
}

func toTrainingMenus(trainingMenuModels []TrainingMenuModel) ([]training.TrainingMenu, error) {
	trainingMenus := make([]training.TrainingMenu, 0, len(trainingMenuModels))
	for _, trainingMenuModel := range trainingMenuModels {
		trainingMenu, err := toTrainingMenu(trainingMenuModel)
		if err != nil {
			return nil, err
		}
		trainingMenus = append(trainingMenus, trainingMenu)
	}

	return trainingMenus, nil
}

func (r *trainingMenuRepository) FindByIDs(ctx context.Context, trainingMenuIDs []training.TrainingMenuID) ([]training.TrainingMenu, error) {
	if len(trainingMenuIDs) == 0 {
		return []training.TrainingMenu{}, nil
	}

	var trainingMenuModels []TrainingMenuModel
	if err := r.db.WithContext(ctx).Where("id IN ?", trainingMenuIDs).Find(&trainingMenuModels).Error; err != nil {
		return nil, err
	}

	return toTrainingMenus(trainingMenuModels)
}

func (r *trainingMenuRepository) List(ctx context.Context) ([]training.TrainingMenu, error) {
	var trainingMenuModels []TrainingMenuModel
	if err := r.db.WithContext(ctx).Find(&trainingMenuModels).Error; err != nil {
		return nil, err
	}

	trainingMenus, err := toTrainingMenus(trainingMenuModels)
	if err != nil {
		return nil, err
	}

	sort.Slice(trainingMenus, func(i, j int) bool {
		iElementOrder, jElementOrder := trainingMenus[i].Element().Order(), trainingMenus[j].Element().Order()
		if iElementOrder != jElementOrder {
			return iElementOrder < jElementOrder
		}
		iPartOrder, jPartOrder := trainingMenus[i].Part().Order(), trainingMenus[j].Part().Order()
		if iPartOrder != jPartOrder {
			return iPartOrder < jPartOrder
		}
		return trainingMenus[i].Code() < trainingMenus[j].Code()
	})

	return trainingMenus, nil
}
