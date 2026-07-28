package organization

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type OrganizationModel struct {
	ID        organization.OrganizationID
	GroupID   tenant.TenantID
	Name      organization.Name
	CreatedAt time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}

func (OrganizationModel) TableName() string {
	return "organizations"
}
