package standard

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/standard"
)

type ageGroupStandardRepository struct {
	db *gorm.DB
}

func NewAgeGroupStandardRepository(db *gorm.DB) standard.AgeGroupStandardRepository {
	return &ageGroupStandardRepository{db: db}
}

func toAgeGroupStandard(m AgeGroupStandardModel) (standard.AgeGroupStandard, error) {
	gender, err := standard.NewGender(m.Gender.String())
	if err != nil {
		return nil, fmt.Errorf("age group standard %q: %w", m.ID, err)
	}
	ageRange, err := standard.NewAgeRange(m.AgeFrom, m.AgeTo)
	if err != nil {
		return nil, fmt.Errorf("age group standard %q: %w", m.ID, err)
	}
	mean, err := standard.NewMean(m.Mean.Float64())
	if err != nil {
		return nil, fmt.Errorf("age group standard %q: %w", m.ID, err)
	}
	standardDeviation, err := standard.NewStandardDeviation(m.StandardDeviation.Float64())
	if err != nil {
		return nil, fmt.Errorf("age group standard %q: %w", m.ID, err)
	}

	return standard.NewAgeGroupStandard(
		m.ID,
		m.MeasurementItemID,
		gender,
		ageRange,
		mean,
		standardDeviation,
		m.CreatedAt,
		m.UpdatedAt,
	), nil
}

func toAgeGroupStandards(models []AgeGroupStandardModel) ([]standard.AgeGroupStandard, error) {
	ageGroupStandards := make([]standard.AgeGroupStandard, 0, len(models))
	for _, model := range models {
		ageGroupStandard, err := toAgeGroupStandard(model)
		if err != nil {
			return nil, err
		}
		ageGroupStandards = append(ageGroupStandards, ageGroupStandard)
	}

	return ageGroupStandards, nil
}

func (r *ageGroupStandardRepository) ListByItemID(ctx context.Context, measurementItemID measurementitem.MeasurementItemID) ([]standard.AgeGroupStandard, error) {
	var ageGroupStandardModels []AgeGroupStandardModel
	if err := r.db.WithContext(ctx).Where("measurement_item_id = ?", measurementItemID).Order("gender, age_from").Find(&ageGroupStandardModels).Error; err != nil {
		return nil, err
	}

	return toAgeGroupStandards(ageGroupStandardModels)
}

func (r *ageGroupStandardRepository) ListByItemIDsAndGender(ctx context.Context, measurementItemIDs []measurementitem.MeasurementItemID, gender standard.Gender) ([]standard.AgeGroupStandard, error) {
	if len(measurementItemIDs) == 0 {
		return []standard.AgeGroupStandard{}, nil
	}

	var ageGroupStandardModels []AgeGroupStandardModel
	if err := r.db.WithContext(ctx).Where("measurement_item_id IN ? AND gender = ?", measurementItemIDs, gender).Order("measurement_item_id, age_from").Find(&ageGroupStandardModels).Error; err != nil {
		return nil, err
	}

	return toAgeGroupStandards(ageGroupStandardModels)
}
