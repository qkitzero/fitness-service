package judgment

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/qkitzero/fitness-service/internal/domain/judgment"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
)

type judgmentRepository struct {
	db *gorm.DB
}

func NewJudgmentRepository(db *gorm.DB) judgment.JudgmentRepository {
	return &judgmentRepository{db: db}
}

func toModel(j judgment.Judgment) JudgmentModel {
	return JudgmentModel{
		MeasurementID: j.MeasurementID(),
		Advice:        j.Advice(),
		CreatedAt:     j.CreatedAt(),
		UpdatedAt:     j.UpdatedAt(),
	}
}

func toDomain(m JudgmentModel) judgment.Judgment {
	return judgment.ReconstructJudgment(
		m.MeasurementID,
		m.Advice,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

func (r *judgmentRepository) FindByMeasurementID(ctx context.Context, measurementID measurement.MeasurementID) (judgment.Judgment, error) {
	var judgmentModel JudgmentModel
	err := r.db.WithContext(ctx).Where("measurement_id = ?", measurementID).First(&judgmentModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, judgment.ErrJudgmentNotFound
	}
	if err != nil {
		return nil, err
	}

	return toDomain(judgmentModel), nil
}

func (r *judgmentRepository) Upsert(ctx context.Context, j judgment.Judgment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		judgmentModel := toModel(j)

		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "measurement_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"advice",
				"updated_at",
			}),
		}).Create(&judgmentModel).Error; err != nil {
			return err
		}

		return nil
	})
}
