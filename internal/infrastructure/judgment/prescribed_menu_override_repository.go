package judgment

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/judgment"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/training"
)

type prescribedMenuOverrideRepository struct {
	db *gorm.DB
}

func NewPrescribedMenuOverrideRepository(db *gorm.DB) judgment.PrescribedMenuOverrideRepository {
	return &prescribedMenuOverrideRepository{db: db}
}

func toPrescribedMenuOverrideModel(o judgment.PrescribedMenuOverride) PrescribedMenuOverrideModel {
	return PrescribedMenuOverrideModel{
		ID:             o.ID(),
		MeasurementID:  o.MeasurementID(),
		SortOrder:      o.SortOrder(),
		Element:        o.Element(),
		Part:           o.Part(),
		TrainingMenuID: o.TrainingMenuID(),
		Amount:         o.Amount(),
		Unit:           o.Unit(),
		Sets:           o.Sets(),
		CreatedAt:      o.CreatedAt(),
		UpdatedAt:      o.UpdatedAt(),
	}
}

func toPrescribedMenuOverride(m PrescribedMenuOverrideModel) (judgment.PrescribedMenuOverride, error) {
	sortOrder, err := training.NewSortOrder(m.SortOrder.Int())
	if err != nil {
		return nil, fmt.Errorf("prescribed menu override %s: %w", m.ID, err)
	}
	var element *measurementitem.Element
	if m.Element != nil {
		parsed, err := measurementitem.NewElement(m.Element.String())
		if err != nil {
			return nil, fmt.Errorf("prescribed menu override %s: %w", m.ID, err)
		}
		element = &parsed
	}
	var part *training.Part
	if m.Part != nil {
		parsed, err := training.NewPart(m.Part.String())
		if err != nil {
			return nil, fmt.Errorf("prescribed menu override %s: %w", m.ID, err)
		}
		part = &parsed
	}
	amount, err := training.NewAmount(m.Amount.Int())
	if err != nil {
		return nil, fmt.Errorf("prescribed menu override %s: %w", m.ID, err)
	}
	unit, err := training.NewUnit(m.Unit.String())
	if err != nil {
		return nil, fmt.Errorf("prescribed menu override %s: %w", m.ID, err)
	}
	sets, err := training.NewSets(m.Sets.Int())
	if err != nil {
		return nil, fmt.Errorf("prescribed menu override %s: %w", m.ID, err)
	}

	return judgment.NewPrescribedMenuOverride(
		m.ID,
		m.MeasurementID,
		sortOrder,
		element,
		part,
		m.TrainingMenuID,
		amount,
		unit,
		sets,
		m.CreatedAt,
		m.UpdatedAt,
	), nil
}

func (r *prescribedMenuOverrideRepository) ListByMeasurementID(ctx context.Context, measurementID measurement.MeasurementID) ([]judgment.PrescribedMenuOverride, error) {
	var prescribedMenuOverrideModels []PrescribedMenuOverrideModel
	if err := r.db.WithContext(ctx).Where("measurement_id = ?", measurementID).Order("sort_order").Find(&prescribedMenuOverrideModels).Error; err != nil {
		return nil, err
	}

	overrides := make([]judgment.PrescribedMenuOverride, 0, len(prescribedMenuOverrideModels))
	for _, prescribedMenuOverrideModel := range prescribedMenuOverrideModels {
		override, err := toPrescribedMenuOverride(prescribedMenuOverrideModel)
		if err != nil {
			return nil, err
		}
		overrides = append(overrides, override)
	}

	return overrides, nil
}

func (r *prescribedMenuOverrideRepository) ReplaceByMeasurementID(ctx context.Context, measurementID measurement.MeasurementID, overrides []judgment.PrescribedMenuOverride) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lockedMeasurementID measurement.MeasurementID
		if err := tx.Raw(`SELECT id FROM measurements WHERE id = ? FOR UPDATE`, measurementID).Scan(&lockedMeasurementID).Error; err != nil {
			return err
		}

		if err := tx.Where("measurement_id = ?", measurementID).Delete(&PrescribedMenuOverrideModel{}).Error; err != nil {
			return err
		}

		if len(overrides) == 0 {
			return nil
		}

		prescribedMenuOverrideModels := make([]PrescribedMenuOverrideModel, 0, len(overrides))
		for _, override := range overrides {
			prescribedMenuOverrideModels = append(prescribedMenuOverrideModels, toPrescribedMenuOverrideModel(override))
		}

		return tx.Create(&prescribedMenuOverrideModels).Error
	})
}

func (r *prescribedMenuOverrideRepository) DeleteByMeasurementID(ctx context.Context, measurementID measurement.MeasurementID) error {
	return r.db.WithContext(ctx).Where("measurement_id = ?", measurementID).Delete(&PrescribedMenuOverrideModel{}).Error
}
