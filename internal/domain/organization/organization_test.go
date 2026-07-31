package organization

import (
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

func TestNewOrganization(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		organizationID   string
		tenantID         string
		organizationName string
	}{
		{"success new organization", "3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f", "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b", "テスト株式会社"},
		{"success new organization of another tenant", "4e2f7b6d-8c9a-5d3e-0b1f-2a3c4d5e6f70", "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d", "別の会社"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, _ := NewOrganizationIDFromString(tt.organizationID)
			tenantID, _ := tenant.NewTenantID(tt.tenantID)
			organizationName, _ := NewName(tt.organizationName)

			createdAt := time.Now()
			updatedAt := time.Now()
			o := NewOrganization(id, tenantID, organizationName, createdAt, updatedAt)

			if o.ID() != id {
				t.Errorf("ID() = %v, want %v", o.ID(), id)
			}
			if o.TenantID() != tenantID {
				t.Errorf("TenantID() = %v, want %v", o.TenantID(), tenantID)
			}
			if o.Name() != organizationName {
				t.Errorf("Name() = %v, want %v", o.Name(), organizationName)
			}
			if !o.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", o.CreatedAt(), createdAt)
			}
			if !o.UpdatedAt().Equal(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", o.UpdatedAt(), updatedAt)
			}
		})
	}
}

func TestUpdateOrganization(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		organizationName string
		updatedName      string
	}{
		{"success update organization", "テスト株式会社", "更新株式会社"},
		{"success update organization to the same name", "テスト株式会社", "テスト株式会社"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, _ := NewOrganizationIDFromString("3d1e6a5c-7b8f-4c2d-9a0e-1f2b3c4d5e6f")
			tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
			organizationName, _ := NewName(tt.organizationName)
			past := time.Now().UTC().Add(-time.Hour)
			o := NewOrganization(id, tenantID, organizationName, past, past)

			updatedName, _ := NewName(tt.updatedName)

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
		})
	}
}
