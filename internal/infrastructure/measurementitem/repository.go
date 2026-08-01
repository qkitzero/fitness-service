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
	valueType, err := measurementitem.NewValueType(m.ValueType.String())
	if err != nil {
		return nil, fmt.Errorf("measurement item %q: %w", m.Code, err)
	}

	return measurementitem.NewMeasurementItem(
		m.ID,
		code,
		name,
		category,
		unit,
		trialCount,
		m.Bilateral,
		valueType,
		m.CreatedAt,
		m.UpdatedAt,
	), nil
}

func (r *measurementItemRepository) List(ctx context.Context) ([]measurementitem.MeasurementItem, error) {
	var measurementItemModels []MeasurementItemModel
	if err := r.db.WithContext(ctx).Find(&measurementItemModels).Error; err != nil {
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
