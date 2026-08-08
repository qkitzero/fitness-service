package training

import (
	"context"
	"fmt"
	"sort"

	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/training"
)

type prescriptionRuleRepository struct {
	db *gorm.DB
}

func NewPrescriptionRuleRepository(db *gorm.DB) training.PrescriptionRuleRepository {
	return &prescriptionRuleRepository{db: db}
}

func toElementMenu(m ElementMenuModel) (training.ElementMenu, error) {
	element, err := measurementitem.NewElement(m.Element.String())
	if err != nil {
		return nil, fmt.Errorf("element menu %q %q level %d: %w", m.Element, m.Part, m.Level.Int(), err)
	}
	part, err := training.NewPart(m.Part.String())
	if err != nil {
		return nil, fmt.Errorf("element menu %q %q level %d: %w", m.Element, m.Part, m.Level.Int(), err)
	}
	level, err := training.NewLevel(m.Level.Int())
	if err != nil {
		return nil, fmt.Errorf("element menu %q %q level %d: %w", m.Element, m.Part, m.Level.Int(), err)
	}

	return training.NewElementMenu(
		element,
		part,
		level,
		m.TrainingMenuID,
		m.CreatedAt,
		m.UpdatedAt,
	), nil
}

func toFixedMenu(m FixedMenuModel) (training.FixedMenu, error) {
	sortOrder, err := training.NewSortOrder(m.SortOrder.Int())
	if err != nil {
		return nil, fmt.Errorf("fixed menu %d: %w", m.SortOrder.Int(), err)
	}

	return training.NewFixedMenu(
		sortOrder,
		m.TrainingMenuID,
		m.CreatedAt,
		m.UpdatedAt,
	), nil
}

func toAgeDecadeMenu(m AgeDecadeMenuModel) (training.AgeDecadeMenu, error) {
	decade, err := training.NewDecade(m.Decade.Int())
	if err != nil {
		return nil, fmt.Errorf("age decade menu %d %d: %w", m.Decade.Int(), m.SortOrder.Int(), err)
	}
	sortOrder, err := training.NewSortOrder(m.SortOrder.Int())
	if err != nil {
		return nil, fmt.Errorf("age decade menu %d %d: %w", m.Decade.Int(), m.SortOrder.Int(), err)
	}

	return training.NewAgeDecadeMenu(
		decade,
		sortOrder,
		m.TrainingMenuID,
		m.CreatedAt,
		m.UpdatedAt,
	), nil
}

func (r *prescriptionRuleRepository) ListElementMenus(ctx context.Context) ([]training.ElementMenu, error) {
	var elementMenuModels []ElementMenuModel
	if err := r.db.WithContext(ctx).Find(&elementMenuModels).Error; err != nil {
		return nil, err
	}

	elementMenus := make([]training.ElementMenu, 0, len(elementMenuModels))
	for _, elementMenuModel := range elementMenuModels {
		elementMenu, err := toElementMenu(elementMenuModel)
		if err != nil {
			return nil, err
		}
		elementMenus = append(elementMenus, elementMenu)
	}

	sort.Slice(elementMenus, func(i, j int) bool {
		iElementOrder, jElementOrder := elementMenus[i].Element().Order(), elementMenus[j].Element().Order()
		if iElementOrder != jElementOrder {
			return iElementOrder < jElementOrder
		}
		iPartOrder, jPartOrder := elementMenus[i].Part().Order(), elementMenus[j].Part().Order()
		if iPartOrder != jPartOrder {
			return iPartOrder < jPartOrder
		}
		return elementMenus[i].Level() < elementMenus[j].Level()
	})

	return elementMenus, nil
}

func (r *prescriptionRuleRepository) ListFixedMenus(ctx context.Context) ([]training.FixedMenu, error) {
	var fixedMenuModels []FixedMenuModel
	if err := r.db.WithContext(ctx).Order("sort_order").Find(&fixedMenuModels).Error; err != nil {
		return nil, err
	}

	fixedMenus := make([]training.FixedMenu, 0, len(fixedMenuModels))
	for _, fixedMenuModel := range fixedMenuModels {
		fixedMenu, err := toFixedMenu(fixedMenuModel)
		if err != nil {
			return nil, err
		}
		fixedMenus = append(fixedMenus, fixedMenu)
	}

	return fixedMenus, nil
}

func (r *prescriptionRuleRepository) ListAgeDecadeMenus(ctx context.Context) ([]training.AgeDecadeMenu, error) {
	var ageDecadeMenuModels []AgeDecadeMenuModel
	if err := r.db.WithContext(ctx).Order("decade, sort_order").Find(&ageDecadeMenuModels).Error; err != nil {
		return nil, err
	}

	ageDecadeMenus := make([]training.AgeDecadeMenu, 0, len(ageDecadeMenuModels))
	for _, ageDecadeMenuModel := range ageDecadeMenuModels {
		ageDecadeMenu, err := toAgeDecadeMenu(ageDecadeMenuModel)
		if err != nil {
			return nil, err
		}
		ageDecadeMenus = append(ageDecadeMenus, ageDecadeMenu)
	}

	return ageDecadeMenus, nil
}
