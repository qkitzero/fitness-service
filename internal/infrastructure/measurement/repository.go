package measurement

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
)

type measurementRepository struct {
	db *gorm.DB
}

func NewMeasurementRepository(db *gorm.DB) measurement.MeasurementRepository {
	return &measurementRepository{db: db}
}

func toModel(m measurement.Measurement) MeasurementModel {
	return MeasurementModel{
		ID:               m.ID(),
		CustomerID:       m.CustomerID(),
		MeasuredOn:       m.MeasuredOn(),
		MeasuredBy:       m.MeasuredBy(),
		AgeAtMeasurement: m.AgeAtMeasurement(),
		UpdatedBy:        m.UpdatedBy(),
		IsDraft:          m.IsDraft(),
		CreatedAt:        m.CreatedAt(),
		UpdatedAt:        m.UpdatedAt(),
	}
}

func toChildModels(m measurement.Measurement) ([]MeasurementEntryModel, []MeasurementValueModel) {
	entries := m.Entries()
	entryModels := make([]MeasurementEntryModel, 0, len(entries))
	valueModels := make([]MeasurementValueModel, 0, len(entries))
	for _, entry := range entries {
		entryID := uuid.New()
		entryModels = append(entryModels, MeasurementEntryModel{
			ID:                entryID,
			MeasurementID:     m.ID(),
			MeasurementItemID: entry.MeasurementItemID(),
			Unmeasurable:      entry.Unmeasurable(),
			Note:              entry.Note(),
			CreatedAt:         m.UpdatedAt(),
			UpdatedAt:         m.UpdatedAt(),
		})
		for _, value := range entry.Values() {
			valueModels = append(valueModels, MeasurementValueModel{
				ID:                 uuid.New(),
				MeasurementEntryID: entryID,
				TrialIndex:         value.TrialIndex(),
				Side:               value.Side(),
				Value:              value.Value(),
				ValueSecondary:     value.ValueSecondary(),
				ValueChoice:        value.ValueChoice(),
				CreatedAt:          m.UpdatedAt(),
				UpdatedAt:          m.UpdatedAt(),
			})
		}
	}

	return entryModels, valueModels
}

func toDomain(m MeasurementModel) measurement.Measurement {
	entries := make([]measurement.MeasurementEntry, 0, len(m.Entries))
	for _, entryModel := range m.Entries {
		values := make([]measurement.MeasurementValue, 0, len(entryModel.Values))
		for _, valueModel := range entryModel.Values {
			values = append(values, measurement.NewMeasurementValue(
				valueModel.TrialIndex,
				valueModel.Side,
				valueModel.Value,
				valueModel.ValueSecondary,
				valueModel.ValueChoice,
			))
		}
		entries = append(entries, measurement.ReconstructMeasurementEntry(
			entryModel.MeasurementItemID,
			entryModel.Unmeasurable,
			entryModel.Note,
			values,
		))
	}

	return measurement.ReconstructMeasurement(
		m.ID,
		m.CustomerID,
		m.MeasuredOn,
		m.MeasuredBy,
		m.AgeAtMeasurement,
		m.UpdatedBy,
		m.IsDraft,
		entries,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

func createChildren(tx *gorm.DB, m measurement.Measurement) error {
	entryModels, valueModels := toChildModels(m)

	if len(entryModels) > 0 {
		if err := tx.Create(&entryModels).Error; err != nil {
			return err
		}
	}

	if len(valueModels) > 0 {
		if err := tx.Create(&valueModels).Error; err != nil {
			return err
		}
	}

	return nil
}

func withChildren(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Entries", func(db *gorm.DB) *gorm.DB {
			return db.Order("measurement_entries.measurement_item_id")
		}).
		Preload("Entries.Values", func(db *gorm.DB) *gorm.DB {
			return db.Order("measurement_values.trial_index, measurement_values.side")
		})
}

func (r *measurementRepository) Create(ctx context.Context, m measurement.Measurement) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		measurementModel := toModel(m)

		if err := tx.Create(&measurementModel).Error; err != nil {
			return err
		}

		return createChildren(tx, m)
	})
}

func (r *measurementRepository) FindByID(ctx context.Context, id measurement.MeasurementID) (measurement.Measurement, error) {
	var measurementModel MeasurementModel
	err := withChildren(r.db.WithContext(ctx)).Where("id = ?", id).First(&measurementModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, measurement.ErrMeasurementNotFound
	}
	if err != nil {
		return nil, err
	}

	return toDomain(measurementModel), nil
}

func (r *measurementRepository) ListByCustomerID(ctx context.Context, customerID customer.CustomerID) ([]measurement.Measurement, error) {
	var measurementModels []MeasurementModel
	if err := withChildren(r.db.WithContext(ctx)).Where("customer_id = ?", customerID).Order("measured_on DESC, id").Find(&measurementModels).Error; err != nil {
		return nil, err
	}

	measurements := make([]measurement.Measurement, 0, len(measurementModels))
	for _, measurementModel := range measurementModels {
		measurements = append(measurements, toDomain(measurementModel))
	}

	return measurements, nil
}

func (r *measurementRepository) Update(ctx context.Context, m measurement.Measurement) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		measurementModel := toModel(m)

		result := tx.Model(&MeasurementModel{}).
			Where("id = ?", measurementModel.ID).
			Select(
				"measured_on",
				"measured_by",
				"age_at_measurement",
				"updated_by",
				"is_draft",
				"updated_at",
			).
			Updates(measurementModel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return measurement.ErrMeasurementNotFound
		}

		if err := tx.Where("measurement_id = ?", measurementModel.ID).Delete(&MeasurementEntryModel{}).Error; err != nil {
			return err
		}

		return createChildren(tx, m)
	})
}

func (r *measurementRepository) Delete(ctx context.Context, id measurement.MeasurementID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ?", id).Delete(&MeasurementModel{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return measurement.ErrMeasurementNotFound
		}

		return nil
	})
}
