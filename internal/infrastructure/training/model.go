package training

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/training"
)

type TrainingMenuModel struct {
	ID          training.TrainingMenuID
	Code        training.Code
	Name        training.Name
	Element     measurementitem.Element
	Part        training.Part
	Amount      training.Amount
	Unit        training.Unit
	Sets        training.Sets
	Instruction training.Instruction
	CreatedAt   time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:false"`
}

func (TrainingMenuModel) TableName() string {
	return "training_menus"
}

type ElementMenuModel struct {
	Element        measurementitem.Element `gorm:"primaryKey"`
	Part           training.Part           `gorm:"primaryKey"`
	Level          training.Level          `gorm:"primaryKey"`
	TrainingMenuID training.TrainingMenuID
	CreatedAt      time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime:false"`
}

func (ElementMenuModel) TableName() string {
	return "element_menus"
}

type FixedMenuModel struct {
	SortOrder      training.SortOrder `gorm:"primaryKey"`
	TrainingMenuID training.TrainingMenuID
	CreatedAt      time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime:false"`
}

func (FixedMenuModel) TableName() string {
	return "fixed_menus"
}

type AgeDecadeMenuModel struct {
	Decade         training.Decade    `gorm:"primaryKey"`
	SortOrder      training.SortOrder `gorm:"primaryKey"`
	TrainingMenuID training.TrainingMenuID
	CreatedAt      time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime:false"`
}

func (AgeDecadeMenuModel) TableName() string {
	return "age_decade_menus"
}
