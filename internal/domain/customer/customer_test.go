package customer

import (
	"reflect"
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/contact"
	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

func TestNewCustomer(t *testing.T) {
	t.Parallel()
	id, _ := NewCustomerIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := NewName("test customer")
	nameKana, _ := NewNameKana("テストカナ")
	gender, _ := NewGender("male")
	birthDate, _ := NewBirthDate(2000, 1, 1)
	phone, _ := contact.NewPhone("03-1234-5678")
	email, _ := contact.NewEmail("test@example.com")
	postalCode, _ := address.NewPostalCode("123-4567")
	prefecture, _ := address.NewPrefecture("東京都")
	city, _ := address.NewCity("千代田区")
	street, _ := address.NewStreet("1-1-1")
	building, _ := address.NewBuilding("テストビル")
	addr := address.NewAddress(postalCode, prefecture, city, street, building)
	emergencyContactName, _ := NewEmergencyContactName("緊急 太郎")
	emergencyContactRelationship, _ := NewEmergencyContactRelationship("父")
	emergencyContactPhone, _ := contact.NewPhone("090-1234-5678")
	organizationID, _ := organization.NewOrganizationIDFromString("3f2b6c1d-4e5f-6a7b-8c9d-0e1f2a3b4c5d")

	tests := []struct {
		name                         string
		phone                        *contact.Phone
		email                        *contact.Email
		address                      address.Address
		emergencyContactName         *EmergencyContactName
		emergencyContactRelationship *EmergencyContactRelationship
		emergencyContactPhone        *contact.Phone
		organizationID               *organization.OrganizationID
		active                       bool
	}{
		{"success new customer", phone, email, addr, emergencyContactName, emergencyContactRelationship, emergencyContactPhone, &organizationID, true},
		{"success new customer without optional fields", nil, nil, address.Address{}, nil, nil, nil, nil, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			createdAt := time.Now()
			updatedAt := time.Now()
			c := NewCustomer(id, tenantID, name, nameKana, gender, birthDate, tt.phone, tt.email, tt.address, tt.emergencyContactName, tt.emergencyContactRelationship, tt.emergencyContactPhone, tt.organizationID, tt.active, createdAt, updatedAt)

			if c.ID() != id {
				t.Errorf("ID() = %v, want %v", c.ID(), id)
			}
			if c.TenantID() != tenantID {
				t.Errorf("TenantID() = %v, want %v", c.TenantID(), tenantID)
			}
			if c.Name() != name {
				t.Errorf("Name() = %v, want %v", c.Name(), name)
			}
			if c.NameKana() != nameKana {
				t.Errorf("NameKana() = %v, want %v", c.NameKana(), nameKana)
			}
			if c.Gender() != gender {
				t.Errorf("Gender() = %v, want %v", c.Gender(), gender)
			}
			if !c.BirthDate().Equal(birthDate.Time) {
				t.Errorf("BirthDate() = %v, want %v", c.BirthDate(), birthDate)
			}
			if !reflect.DeepEqual(c.Phone(), tt.phone) {
				t.Errorf("Phone() = %v, want %v", c.Phone(), tt.phone)
			}
			if !reflect.DeepEqual(c.Email(), tt.email) {
				t.Errorf("Email() = %v, want %v", c.Email(), tt.email)
			}
			if !reflect.DeepEqual(c.Address(), tt.address) {
				t.Errorf("Address() = %v, want %v", c.Address(), tt.address)
			}
			if !reflect.DeepEqual(c.EmergencyContactName(), tt.emergencyContactName) {
				t.Errorf("EmergencyContactName() = %v, want %v", c.EmergencyContactName(), tt.emergencyContactName)
			}
			if !reflect.DeepEqual(c.EmergencyContactRelationship(), tt.emergencyContactRelationship) {
				t.Errorf("EmergencyContactRelationship() = %v, want %v", c.EmergencyContactRelationship(), tt.emergencyContactRelationship)
			}
			if !reflect.DeepEqual(c.EmergencyContactPhone(), tt.emergencyContactPhone) {
				t.Errorf("EmergencyContactPhone() = %v, want %v", c.EmergencyContactPhone(), tt.emergencyContactPhone)
			}
			if !reflect.DeepEqual(c.OrganizationID(), tt.organizationID) {
				t.Errorf("OrganizationID() = %v, want %v", c.OrganizationID(), tt.organizationID)
			}
			if c.IsActive() != tt.active {
				t.Errorf("IsActive() = %v, want %v", c.IsActive(), tt.active)
			}
			if !c.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", c.CreatedAt(), createdAt)
			}
			if !c.UpdatedAt().Equal(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", c.UpdatedAt(), updatedAt)
			}
		})
	}
}

func TestUpdateCustomer(t *testing.T) {
	t.Parallel()
	id, _ := NewCustomerIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := NewName("test customer")
	nameKana, _ := NewNameKana("テストカナ")
	gender, _ := NewGender("male")
	birthDate, _ := NewBirthDate(2000, 1, 1)
	updatedName, _ := NewName("updated test customer")
	updatedNameKana, _ := NewNameKana("コウシンカナ")
	updatedGender, _ := NewGender("female")
	updatedBirthDate, _ := NewBirthDate(1999, 12, 31)
	phone, _ := contact.NewPhone("090-1234-5678")
	email, _ := contact.NewEmail("updated@example.com")
	postalCode, _ := address.NewPostalCode("543-2100")
	prefecture, _ := address.NewPrefecture("大阪府")
	city, _ := address.NewCity("大阪市")
	street, _ := address.NewStreet("2-2-2")
	building, _ := address.NewBuilding("更新ビル")
	addr := address.NewAddress(postalCode, prefecture, city, street, building)
	emergencyContactName, _ := NewEmergencyContactName("更新 花子")
	emergencyContactRelationship, _ := NewEmergencyContactRelationship("母")
	emergencyContactPhone, _ := contact.NewPhone("080-1234-5678")
	organizationID, _ := organization.NewOrganizationIDFromString("3f2b6c1d-4e5f-6a7b-8c9d-0e1f2a3b4c5d")

	tests := []struct {
		name           string
		organizationID *organization.OrganizationID
	}{
		{"success update customer", &organizationID},
		{"success update customer clearing the organization", nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			past := time.Now().UTC().Add(-time.Hour)
			c := NewCustomer(id, tenantID, name, nameKana, gender, birthDate, nil, nil, address.Address{}, nil, nil, nil, nil, true, past, past)

			initialAddr := c.Address()
			if c.Phone() != nil || c.Email() != nil || initialAddr.PostalCode() != nil || initialAddr.Prefecture() != nil || initialAddr.City() != nil || initialAddr.Street() != nil || initialAddr.Building() != nil || c.EmergencyContactName() != nil || c.EmergencyContactRelationship() != nil || c.EmergencyContactPhone() != nil || c.OrganizationID() != nil {
				t.Errorf("expected nil optional fields before update")
			}

			c.Update(updatedName, updatedNameKana, updatedGender, updatedBirthDate, phone, email, addr, emergencyContactName, emergencyContactRelationship, emergencyContactPhone, tt.organizationID)

			if c.Name() != updatedName {
				t.Errorf("Name() = %v, want %v", c.Name(), updatedName)
			}
			if c.NameKana() != updatedNameKana {
				t.Errorf("NameKana() = %v, want %v", c.NameKana(), updatedNameKana)
			}
			if c.Gender() != updatedGender {
				t.Errorf("Gender() = %v, want %v", c.Gender(), updatedGender)
			}
			if !c.BirthDate().Equal(updatedBirthDate.Time) {
				t.Errorf("BirthDate() = %v, want %v", c.BirthDate(), updatedBirthDate)
			}
			if !reflect.DeepEqual(c.Phone(), phone) {
				t.Errorf("Phone() = %v, want %v", c.Phone(), phone)
			}
			if !reflect.DeepEqual(c.Email(), email) {
				t.Errorf("Email() = %v, want %v", c.Email(), email)
			}
			if !reflect.DeepEqual(c.Address(), addr) {
				t.Errorf("Address() = %v, want %v", c.Address(), addr)
			}
			if !reflect.DeepEqual(c.EmergencyContactName(), emergencyContactName) {
				t.Errorf("EmergencyContactName() = %v, want %v", c.EmergencyContactName(), emergencyContactName)
			}
			if !reflect.DeepEqual(c.EmergencyContactRelationship(), emergencyContactRelationship) {
				t.Errorf("EmergencyContactRelationship() = %v, want %v", c.EmergencyContactRelationship(), emergencyContactRelationship)
			}
			if !reflect.DeepEqual(c.EmergencyContactPhone(), emergencyContactPhone) {
				t.Errorf("EmergencyContactPhone() = %v, want %v", c.EmergencyContactPhone(), emergencyContactPhone)
			}
			if !reflect.DeepEqual(c.OrganizationID(), tt.organizationID) {
				t.Errorf("OrganizationID() = %v, want %v", c.OrganizationID(), tt.organizationID)
			}
			if !c.IsActive() {
				t.Errorf("IsActive() = %v, want %v", c.IsActive(), true)
			}
			if !c.CreatedAt().Equal(past) {
				t.Errorf("CreatedAt() = %v, want %v", c.CreatedAt(), past)
			}
			if !c.UpdatedAt().After(past) {
				t.Errorf("UpdatedAt() = %v, want after %v", c.UpdatedAt(), past)
			}
			if c.UpdatedAt().Location() != time.UTC {
				t.Errorf("UpdatedAt().Location() = %v, want %v", c.UpdatedAt().Location(), time.UTC)
			}
		})
	}
}

func TestSetActiveCustomer(t *testing.T) {
	t.Parallel()
	id, _ := NewCustomerIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := NewName("test customer")
	nameKana, _ := NewNameKana("テストカナ")
	gender, _ := NewGender("male")
	birthDate, _ := NewBirthDate(2000, 1, 1)

	tests := []struct {
		name        string
		initial     bool
		active      bool
		wantUpdated bool
	}{
		{"success deactivate customer", true, false, true},
		{"success activate customer", false, true, true},
		{"success deactivate an already inactive customer", false, false, false},
		{"success activate an already active customer", true, true, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			past := time.Now().UTC().Add(-time.Hour)
			c := NewCustomer(id, tenantID, name, nameKana, gender, birthDate, nil, nil, address.Address{}, nil, nil, nil, nil, tt.initial, past, past)

			c.SetActive(tt.active)

			if c.IsActive() != tt.active {
				t.Errorf("IsActive() = %v, want %v", c.IsActive(), tt.active)
			}
			if !c.CreatedAt().Equal(past) {
				t.Errorf("CreatedAt() = %v, want %v", c.CreatedAt(), past)
			}
			if tt.wantUpdated && !c.UpdatedAt().After(past) {
				t.Errorf("UpdatedAt() = %v, want after %v", c.UpdatedAt(), past)
			}
			if tt.wantUpdated && c.UpdatedAt().Location() != time.UTC {
				t.Errorf("UpdatedAt().Location() = %v, want %v", c.UpdatedAt().Location(), time.UTC)
			}
			if !tt.wantUpdated && !c.UpdatedAt().Equal(past) {
				t.Errorf("UpdatedAt() = %v, want unchanged %v", c.UpdatedAt(), past)
			}
		})
	}
}
