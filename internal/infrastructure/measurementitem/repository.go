package measurementitem

import (
	"context"
	"fmt"
	"sort"

	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
)

type measurementItemRepository struct {
	db *gorm.DB
}

func NewMeasurementItemRepository(db *gorm.DB) measurementitem.MeasurementItemRepository {
	return &measurementItemRepository{db: db}
}

func toElements(m MeasurementItemModel) ([]measurementitem.Element, error) {
	elements := make([]measurementitem.Element, 0, len(m.Elements))
	for _, elementModel := range m.Elements {
		element, err := measurementitem.NewElement(elementModel.Element.String())
		if err != nil {
			return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
		}
		elements = append(elements, element)
	}

	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Order() < elements[j].Order()
	})

	return elements, nil
}

func toDomain(m MeasurementItemModel) (measurementitem.MeasurementItem, error) {
	code, err := measurementitem.NewCode(m.Code.String())
	if err != nil {
		return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
	}
	name, err := measurementitem.NewName(m.Name.String())
	if err != nil {
		return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
	}
	category, err := measurementitem.NewCategory(m.Category.String())
	if err != nil {
		return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
	}
	unit, err := measurementitem.NewUnit(m.Unit.String())
	if err != nil {
		return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
	}
	trialCount, err := measurementitem.NewTrialCount(m.TrialCount.Int())
	if err != nil {
		return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
	}
	sideMode, err := measurementitem.NewSideMode(m.SideMode.String())
	if err != nil {
		return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
	}
	valueType, err := measurementitem.NewValueType(m.ValueType.String())
	if err != nil {
		return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
	}
	var scoreDirection *measurementitem.ScoreDirection
	if m.ScoreDirection != nil {
		s, err := measurementitem.NewScoreDirection(m.ScoreDirection.String())
		if err != nil {
			return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
		}
		scoreDirection = &s
	}
	trialAggregation, err := measurementitem.NewTrialAggregation(m.TrialAggregation.String())
	if err != nil {
		return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
	}
	sideAggregation, err := measurementitem.NewSideAggregation(m.SideAggregation.String())
	if err != nil {
		return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
	}
	normalization, err := measurementitem.NewNormalization(m.Normalization.String())
	if err != nil {
		return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
	}
	elements, err := toElements(m)
	if err != nil {
		return nil, err
	}

	return measurementitem.NewMeasurementItem(
		m.ID,
		code,
		name,
		category,
		unit,
		trialCount,
		sideMode,
		valueType,
		scoreDirection,
		trialAggregation,
		sideAggregation,
		normalization,
		elements,
		m.CreatedAt,
		m.UpdatedAt,
	), nil
}

func withElements(db *gorm.DB) *gorm.DB {
	return db.Preload("Elements")
}

func (r *measurementItemRepository) FindByIDs(ctx context.Context, measurementItemIDs []measurementitem.MeasurementItemID) ([]measurementitem.MeasurementItem, error) {
	if len(measurementItemIDs) == 0 {
		return []measurementitem.MeasurementItem{}, nil
	}

	var measurementItemModels []MeasurementItemModel
	if err := withElements(r.db.WithContext(ctx)).Where("id IN ?", measurementItemIDs).Find(&measurementItemModels).Error; err != nil {
		return nil, err
	}

	measurementItems := make([]measurementitem.MeasurementItem, 0, len(measurementItemModels))
	for _, measurementItemModel := range measurementItemModels {
		measurementItem, err := toDomain(measurementItemModel)
		if err != nil {
			return nil, err
		}
		measurementItems = append(measurementItems, measurementItem)
	}

	return measurementItems, nil
}

func (r *measurementItemRepository) List(ctx context.Context) ([]measurementitem.MeasurementItem, error) {
	var measurementItemModels []MeasurementItemModel
	if err := withElements(r.db.WithContext(ctx)).Find(&measurementItemModels).Error; err != nil {
		return nil, err
	}

	measurementItems := make([]measurementitem.MeasurementItem, 0, len(measurementItemModels))
	for _, measurementItemModel := range measurementItemModels {
		measurementItem, err := toDomain(measurementItemModel)
		if err != nil {
			return nil, err
		}
		measurementItems = append(measurementItems, measurementItem)
	}

	sort.Slice(measurementItems, func(i, j int) bool {
		iOrder, jOrder := measurementItems[i].Category().Order(), measurementItems[j].Category().Order()
		if iOrder != jOrder {
			return iOrder < jOrder
		}
		return measurementItems[i].Code() < measurementItems[j].Code()
	})

	return measurementItems, nil
}
