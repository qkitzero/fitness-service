package tenant

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/contact"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type ProfileModel struct {
	TenantID    tenant.TenantID `gorm:"primaryKey"`
	PostalCode  *address.PostalCode
	Prefecture  *address.Prefecture
	City        *address.City
	Street      *address.Street
	Building    *address.Building
	Phone       *contact.Phone
	Email       *contact.Email
	HomepageURL *tenant.HomepageURL
	Note        *tenant.Note
	CreatedAt   time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:false"`
}

func (ProfileModel) TableName() string {
	return "tenant_profiles"
}
