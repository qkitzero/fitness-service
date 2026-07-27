package organization

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/organization"
)

type OrganizationModel struct {
	ID        organization.OrganizationID
	GroupID   organization.GroupID
	Name      organization.Name
	CreatedAt time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}

func (OrganizationModel) TableName() string {
	return "organizations"
}
