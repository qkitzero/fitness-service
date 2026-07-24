package customer

import (
	"testing"
	"time"
)

func TestNewCustomer(t *testing.T) {
	t.Parallel()
	id, err := NewCustomerIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new customer id: %v", err)
	}
	groupID, err := NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	if err != nil {
		t.Errorf("failed to new group id: %v", err)
	}
	name, err := NewName("test customer")
	if err != nil {
		t.Errorf("failed to new name: %v", err)
	}
	nameKana, err := NewNameKana("テストカナ")
	if err != nil {
		t.Errorf("failed to new name kana: %v", err)
	}
	gender, err := NewGender("male")
	if err != nil {
		t.Errorf("failed to new gender: %v", err)
	}
	birthDate, err := NewBirthDate(2000, 1, 1)
	if err != nil {
		t.Errorf("failed to new birth date: %v", err)
	}
	tests := []struct {
		name      string
		success   bool
		id        CustomerID
		groupID   GroupID
		customer  Name
		nameKana  NameKana
		gender    Gender
		birthDate BirthDate
		createdAt time.Time
		updatedAt time.Time
	}{
		{"success new customer", true, id, groupID, name, nameKana, gender, birthDate, time.Now(), time.Now()},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			customer := NewCustomer(tt.id, tt.groupID, tt.customer, tt.nameKana, tt.gender, tt.birthDate, tt.createdAt, tt.updatedAt)
			if tt.success && customer.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", customer.ID(), tt.id)
			}
			if tt.success && customer.GroupID() != tt.groupID {
				t.Errorf("GroupID() = %v, want %v", customer.GroupID(), tt.groupID)
			}
			if tt.success && customer.Name() != tt.customer {
				t.Errorf("Name() = %v, want %v", customer.Name(), tt.customer)
			}
			if tt.success && customer.NameKana() != tt.nameKana {
				t.Errorf("NameKana() = %v, want %v", customer.NameKana(), tt.nameKana)
			}
			if tt.success && customer.Gender() != tt.gender {
				t.Errorf("Gender() = %v, want %v", customer.Gender(), tt.gender)
			}
			if tt.success && !customer.BirthDate().Equal(tt.birthDate.Time) {
				t.Errorf("BirthDate() = %v, want %v", customer.BirthDate(), tt.birthDate)
			}
			if tt.success && !customer.CreatedAt().Equal(tt.createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", customer.CreatedAt(), tt.createdAt)
			}
			if tt.success && !customer.UpdatedAt().Equal(tt.updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", customer.UpdatedAt(), tt.updatedAt)
			}
		})
	}
}

func TestUpdateCustomer(t *testing.T) {
	t.Parallel()
	id, err := NewCustomerIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new customer id: %v", err)
	}
	groupID, err := NewGroupID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	if err != nil {
		t.Errorf("failed to new group id: %v", err)
	}
	name, err := NewName("test customer")
	if err != nil {
		t.Errorf("failed to new name: %v", err)
	}
	nameKana, err := NewNameKana("テストカナ")
	if err != nil {
		t.Errorf("failed to new name kana: %v", err)
	}
	gender, err := NewGender("male")
	if err != nil {
		t.Errorf("failed to new gender: %v", err)
	}
	birthDate, err := NewBirthDate(2000, 1, 1)
	if err != nil {
		t.Errorf("failed to new birth date: %v", err)
	}
	updatedName, err := NewName("updated test customer")
	if err != nil {
		t.Errorf("failed to new name: %v", err)
	}
	updatedNameKana, err := NewNameKana("コウシンカナ")
	if err != nil {
		t.Errorf("failed to new name kana: %v", err)
	}
	updatedGender, err := NewGender("female")
	if err != nil {
		t.Errorf("failed to new gender: %v", err)
	}
	updatedBirthDate, err := NewBirthDate(1999, 12, 31)
	if err != nil {
		t.Errorf("failed to new birth date: %v", err)
	}
	customer := NewCustomer(id, groupID, name, nameKana, gender, birthDate, time.Now(), time.Now())
	tests := []struct {
		name             string
		success          bool
		customer         Customer
		updatedName      Name
		updatedNameKana  NameKana
		updatedGender    Gender
		updatedBirthDate BirthDate
	}{
		{"success update customer", true, customer, updatedName, updatedNameKana, updatedGender, updatedBirthDate},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.customer.Update(tt.updatedName, tt.updatedNameKana, tt.updatedGender, tt.updatedBirthDate)
			if tt.success && tt.customer.Name() != tt.updatedName {
				t.Errorf("Name() = %v, want %v", tt.customer.Name(), tt.updatedName)
			}
			if tt.success && tt.customer.NameKana() != tt.updatedNameKana {
				t.Errorf("NameKana() = %v, want %v", tt.customer.NameKana(), tt.updatedNameKana)
			}
			if tt.success && tt.customer.Gender() != tt.updatedGender {
				t.Errorf("Gender() = %v, want %v", tt.customer.Gender(), tt.updatedGender)
			}
			if tt.success && !tt.customer.BirthDate().Equal(tt.updatedBirthDate.Time) {
				t.Errorf("BirthDate() = %v, want %v", tt.customer.BirthDate(), tt.updatedBirthDate)
			}
			if tt.success && !tt.customer.CreatedAt().Before(tt.customer.UpdatedAt()) {
				t.Errorf("CreatedAt() = %v, UpdatedAt() = %v, want CreatedAt < UpdatedAt", tt.customer.CreatedAt(), tt.customer.UpdatedAt())
			}
		})
	}
}
