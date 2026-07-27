package customer

import (
	"testing"
	"time"
)

func TestNewCustomer(t *testing.T) {
	t.Parallel()
	id, _ := NewCustomerIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	groupID, _ := NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := NewName("test customer")
	nameKana, _ := NewNameKana("テストカナ")
	gender, _ := NewGender("male")
	birthDate, _ := NewBirthDate(2000, 1, 1)
	phone, _ := NewPhone("03-1234-5678")
	email, _ := NewEmail("test@example.com")
	postalCode, _ := NewPostalCode("123-4567")
	prefecture, _ := NewPrefecture("東京都")
	city, _ := NewCity("千代田区")
	street, _ := NewStreet("1-1-1")
	building, _ := NewBuilding("テストビル")
	emergencyContactName, _ := NewEmergencyContactName("緊急 太郎")
	emergencyContactRelationship, _ := NewEmergencyContactRelationship("父")
	emergencyContactPhone, _ := NewPhone("090-1234-5678")

	createdAt := time.Now()
	updatedAt := time.Now()
	c := NewCustomer(id, groupID, name, nameKana, gender, birthDate, phone, email, postalCode, prefecture, city, street, building, emergencyContactName, emergencyContactRelationship, emergencyContactPhone, createdAt, updatedAt)

	if c.ID() != id {
		t.Errorf("ID() = %v, want %v", c.ID(), id)
	}
	if c.GroupID() != groupID {
		t.Errorf("GroupID() = %v, want %v", c.GroupID(), groupID)
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
	if c.Phone() == nil || c.Phone().String() != "0312345678" {
		t.Errorf("Phone() = %v, want 0312345678", c.Phone())
	}
	if c.Email() == nil || c.Email().String() != "test@example.com" {
		t.Errorf("Email() = %v, want test@example.com", c.Email())
	}
	if c.PostalCode() == nil || c.PostalCode().String() != "1234567" {
		t.Errorf("PostalCode() = %v, want 1234567", c.PostalCode())
	}
	if c.Prefecture() == nil || c.Prefecture().String() != "東京都" {
		t.Errorf("Prefecture() = %v, want 東京都", c.Prefecture())
	}
	if c.City() == nil || c.City().String() != "千代田区" {
		t.Errorf("City() = %v, want 千代田区", c.City())
	}
	if c.Street() == nil || c.Street().String() != "1-1-1" {
		t.Errorf("Street() = %v, want 1-1-1", c.Street())
	}
	if c.Building() == nil || c.Building().String() != "テストビル" {
		t.Errorf("Building() = %v, want テストビル", c.Building())
	}
	if c.EmergencyContactName() == nil || c.EmergencyContactName().String() != "緊急 太郎" {
		t.Errorf("EmergencyContactName() = %v, want 緊急 太郎", c.EmergencyContactName())
	}
	if c.EmergencyContactRelationship() == nil || c.EmergencyContactRelationship().String() != "父" {
		t.Errorf("EmergencyContactRelationship() = %v, want 父", c.EmergencyContactRelationship())
	}
	if c.EmergencyContactPhone() == nil || c.EmergencyContactPhone().String() != "09012345678" {
		t.Errorf("EmergencyContactPhone() = %v, want 09012345678", c.EmergencyContactPhone())
	}
	if !c.CreatedAt().Equal(createdAt) {
		t.Errorf("CreatedAt() = %v, want %v", c.CreatedAt(), createdAt)
	}
	if !c.UpdatedAt().Equal(updatedAt) {
		t.Errorf("UpdatedAt() = %v, want %v", c.UpdatedAt(), updatedAt)
	}
}

func TestUpdateCustomer(t *testing.T) {
	t.Parallel()
	id, _ := NewCustomerIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	groupID, _ := NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	name, _ := NewName("test customer")
	nameKana, _ := NewNameKana("テストカナ")
	gender, _ := NewGender("male")
	birthDate, _ := NewBirthDate(2000, 1, 1)
	past := time.Now().UTC().Add(-time.Hour)
	c := NewCustomer(id, groupID, name, nameKana, gender, birthDate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, past, past)

	if c.Phone() != nil || c.Email() != nil || c.PostalCode() != nil || c.Prefecture() != nil || c.City() != nil || c.Street() != nil || c.Building() != nil || c.EmergencyContactName() != nil || c.EmergencyContactRelationship() != nil || c.EmergencyContactPhone() != nil {
		t.Errorf("expected nil optional fields before update")
	}

	updatedName, _ := NewName("updated test customer")
	updatedNameKana, _ := NewNameKana("コウシンカナ")
	updatedGender, _ := NewGender("female")
	updatedBirthDate, _ := NewBirthDate(1999, 12, 31)
	phone, _ := NewPhone("090-1234-5678")
	email, _ := NewEmail("updated@example.com")
	postalCode, _ := NewPostalCode("543-2100")
	prefecture, _ := NewPrefecture("大阪府")
	city, _ := NewCity("大阪市")
	street, _ := NewStreet("2-2-2")
	building, _ := NewBuilding("更新ビル")
	emergencyContactName, _ := NewEmergencyContactName("更新 花子")
	emergencyContactRelationship, _ := NewEmergencyContactRelationship("母")
	emergencyContactPhone, _ := NewPhone("080-1234-5678")

	c.Update(updatedName, updatedNameKana, updatedGender, updatedBirthDate, phone, email, postalCode, prefecture, city, street, building, emergencyContactName, emergencyContactRelationship, emergencyContactPhone)

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
	if c.Phone() == nil || c.Phone().String() != "09012345678" {
		t.Errorf("Phone() = %v, want 09012345678", c.Phone())
	}
	if c.Email() == nil || c.Email().String() != "updated@example.com" {
		t.Errorf("Email() = %v, want updated@example.com", c.Email())
	}
	if c.PostalCode() == nil || c.PostalCode().String() != "5432100" {
		t.Errorf("PostalCode() = %v, want 5432100", c.PostalCode())
	}
	if c.Prefecture() == nil || c.Prefecture().String() != "大阪府" {
		t.Errorf("Prefecture() = %v, want 大阪府", c.Prefecture())
	}
	if c.City() == nil || c.City().String() != "大阪市" {
		t.Errorf("City() = %v, want 大阪市", c.City())
	}
	if c.Street() == nil || c.Street().String() != "2-2-2" {
		t.Errorf("Street() = %v, want 2-2-2", c.Street())
	}
	if c.Building() == nil || c.Building().String() != "更新ビル" {
		t.Errorf("Building() = %v, want 更新ビル", c.Building())
	}
	if c.EmergencyContactName() == nil || c.EmergencyContactName().String() != "更新 花子" {
		t.Errorf("EmergencyContactName() = %v, want 更新 花子", c.EmergencyContactName())
	}
	if c.EmergencyContactRelationship() == nil || c.EmergencyContactRelationship().String() != "母" {
		t.Errorf("EmergencyContactRelationship() = %v, want 母", c.EmergencyContactRelationship())
	}
	if c.EmergencyContactPhone() == nil || c.EmergencyContactPhone().String() != "08012345678" {
		t.Errorf("EmergencyContactPhone() = %v, want 08012345678", c.EmergencyContactPhone())
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
}
