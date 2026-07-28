package organization

import (
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

func TestNewOrganization(t *testing.T) {
	t.Parallel()
	id, _ := NewOrganizationIDFromString("3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f")
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := NewName("テスト株式会社")

	createdAt := time.Now()
	updatedAt := time.Now()
	o := NewOrganization(id, tenantID, name, createdAt, updatedAt)

	if o.ID() != id {
		t.Errorf("ID() = %v, want %v", o.ID(), id)
	}
	if o.TenantID() != tenantID {
		t.Errorf("TenantID() = %v, want %v", o.TenantID(), tenantID)
	}
	if o.Name() != name {
		t.Errorf("Name() = %v, want %v", o.Name(), name)
	}
	if !o.CreatedAt().Equal(createdAt) {
		t.Errorf("CreatedAt() = %v, want %v", o.CreatedAt(), createdAt)
	}
	if !o.UpdatedAt().Equal(updatedAt) {
		t.Errorf("UpdatedAt() = %v, want %v", o.UpdatedAt(), updatedAt)
	}
}

func TestUpdateOrganization(t *testing.T) {
	t.Parallel()
	id, _ := NewOrganizationIDFromString("3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f")
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := NewName("テスト株式会社")
	past := time.Now().UTC().Add(-time.Hour)
	o := NewOrganization(id, tenantID, name, past, past)

	updatedName, _ := NewName("更新株式会社")

	o.Update(updatedName)

	if o.Name() != updatedName {
		t.Errorf("Name() = %v, want %v", o.Name(), updatedName)
	}
	if o.ID() != id {
		t.Errorf("ID() = %v, want %v", o.ID(), id)
	}
	if o.TenantID() != tenantID {
		t.Errorf("TenantID() = %v, want %v", o.TenantID(), tenantID)
	}
	if !o.CreatedAt().Equal(past) {
		t.Errorf("CreatedAt() = %v, want %v", o.CreatedAt(), past)
	}
	if !o.UpdatedAt().After(past) {
		t.Errorf("UpdatedAt() = %v, want after %v", o.UpdatedAt(), past)
	}
	if o.UpdatedAt().Location() != time.UTC {
		t.Errorf("UpdatedAt().Location() = %v, want %v", o.UpdatedAt().Location(), time.UTC)
	}
}
