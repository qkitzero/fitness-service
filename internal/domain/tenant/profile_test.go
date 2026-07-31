package tenant

import (
	"reflect"
	"testing"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/contact"
)

func TestNewProfile(t *testing.T) {
	t.Parallel()
	tenantID, _ := NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	postalCode, _ := address.NewPostalCode("123-4567")
	prefecture, _ := address.NewPrefecture("東京都")
	city, _ := address.NewCity("千代田区")
	street, _ := address.NewStreet("1-1-1")
	building, _ := address.NewBuilding("テストビル")
	addr := address.NewAddress(postalCode, prefecture, city, street, building)
	phone, _ := contact.NewPhone("03-1234-5678")
	email, _ := contact.NewEmail("test@example.com")
	homepageURL, _ := NewHomepageURL("https://example.com")
	note, _ := NewNote("テスト備考")

	tests := []struct {
		name        string
		address     address.Address
		phone       *contact.Phone
		email       *contact.Email
		homepageURL *HomepageURL
		note        *Note
	}{
		{"success new profile", addr, phone, email, homepageURL, note},
		{"success new profile without optional fields", address.Address{}, nil, nil, nil, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			createdAt := time.Now()
			updatedAt := time.Now()
			p := NewProfile(tenantID, tt.address, tt.phone, tt.email, tt.homepageURL, tt.note, createdAt, updatedAt)

			if p.TenantID() != tenantID {
				t.Errorf("TenantID() = %v, want %v", p.TenantID(), tenantID)
			}
			if !reflect.DeepEqual(p.Address(), tt.address) {
				t.Errorf("Address() = %v, want %v", p.Address(), tt.address)
			}
			if !reflect.DeepEqual(p.Phone(), tt.phone) {
				t.Errorf("Phone() = %v, want %v", p.Phone(), tt.phone)
			}
			if !reflect.DeepEqual(p.Email(), tt.email) {
				t.Errorf("Email() = %v, want %v", p.Email(), tt.email)
			}
			if !reflect.DeepEqual(p.HomepageURL(), tt.homepageURL) {
				t.Errorf("HomepageURL() = %v, want %v", p.HomepageURL(), tt.homepageURL)
			}
			if !reflect.DeepEqual(p.Note(), tt.note) {
				t.Errorf("Note() = %v, want %v", p.Note(), tt.note)
			}
			if !p.CreatedAt().Equal(createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", p.CreatedAt(), createdAt)
			}
			if !p.UpdatedAt().Equal(updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", p.UpdatedAt(), updatedAt)
			}
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	t.Parallel()
	tenantID, _ := NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	postalCode, _ := address.NewPostalCode("543-2100")
	prefecture, _ := address.NewPrefecture("大阪府")
	city, _ := address.NewCity("大阪市")
	street, _ := address.NewStreet("2-2-2")
	building, _ := address.NewBuilding("更新ビル")
	addr := address.NewAddress(postalCode, prefecture, city, street, building)
	phone, _ := contact.NewPhone("080-1234-5678")
	email, _ := contact.NewEmail("updated@example.com")
	homepageURL, _ := NewHomepageURL("https://updated.example.com")
	note, _ := NewNote("更新備考")

	tests := []struct {
		name        string
		address     address.Address
		phone       *contact.Phone
		email       *contact.Email
		homepageURL *HomepageURL
		note        *Note
	}{
		{"success update profile", addr, phone, email, homepageURL, note},
		{"success update profile clearing optional fields", address.Address{}, nil, nil, nil, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			past := time.Now().UTC().Add(-time.Hour)
			p := NewProfile(tenantID, address.Address{}, nil, nil, nil, nil, past, past)

			initialAddr := p.Address()
			if initialAddr.PostalCode() != nil || initialAddr.Prefecture() != nil || initialAddr.City() != nil || initialAddr.Street() != nil || initialAddr.Building() != nil || p.Phone() != nil || p.Email() != nil || p.HomepageURL() != nil || p.Note() != nil {
				t.Errorf("expected nil optional fields before update")
			}

			p.Update(tt.address, tt.phone, tt.email, tt.homepageURL, tt.note)

			if p.TenantID() != tenantID {
				t.Errorf("TenantID() = %v, want %v", p.TenantID(), tenantID)
			}
			if !reflect.DeepEqual(p.Address(), tt.address) {
				t.Errorf("Address() = %v, want %v", p.Address(), tt.address)
			}
			if !reflect.DeepEqual(p.Phone(), tt.phone) {
				t.Errorf("Phone() = %v, want %v", p.Phone(), tt.phone)
			}
			if !reflect.DeepEqual(p.Email(), tt.email) {
				t.Errorf("Email() = %v, want %v", p.Email(), tt.email)
			}
			if !reflect.DeepEqual(p.HomepageURL(), tt.homepageURL) {
				t.Errorf("HomepageURL() = %v, want %v", p.HomepageURL(), tt.homepageURL)
			}
			if !reflect.DeepEqual(p.Note(), tt.note) {
				t.Errorf("Note() = %v, want %v", p.Note(), tt.note)
			}
			if !p.CreatedAt().Equal(past) {
				t.Errorf("CreatedAt() = %v, want %v", p.CreatedAt(), past)
			}
			if !p.UpdatedAt().After(past) {
				t.Errorf("UpdatedAt() = %v, want after %v", p.UpdatedAt(), past)
			}
			if p.UpdatedAt().Location() != time.UTC {
				t.Errorf("UpdatedAt().Location() = %v, want %v", p.UpdatedAt().Location(), time.UTC)
			}
		})
	}
}

func TestProfileIsImmutable(t *testing.T) {
	t.Parallel()
	tenantID, _ := NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")

	tests := []struct {
		name   string
		mutate func(p Profile, phone *contact.Phone, email *contact.Email, homepageURL *HomepageURL, note *Note)
	}{
		{
			"mutating the constructor arguments",
			func(p Profile, phone *contact.Phone, email *contact.Email, homepageURL *HomepageURL, note *Note) {
				*phone = contact.Phone("09099999999")
				*email = contact.Email("other@example.com")
				*homepageURL = HomepageURL("https://other.example.com")
				*note = Note("別の備考")
			},
		},
		{
			"mutating the getter results",
			func(p Profile, phone *contact.Phone, email *contact.Email, homepageURL *HomepageURL, note *Note) {
				*p.Phone() = contact.Phone("09099999999")
				*p.Email() = contact.Email("other@example.com")
				*p.HomepageURL() = HomepageURL("https://other.example.com")
				*p.Note() = Note("別の備考")
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			phone := contact.Phone("0312345678")
			email := contact.Email("test@example.com")
			homepageURL := HomepageURL("https://example.com")
			note := Note("テスト備考")
			now := time.Now().UTC()

			p := NewProfile(tenantID, address.Address{}, &phone, &email, &homepageURL, &note, now, now)

			tt.mutate(p, &phone, &email, &homepageURL, &note)

			if p.Phone().String() != "0312345678" {
				t.Errorf("Phone() = %v, want %v", p.Phone().String(), "0312345678")
			}
			if p.Email().String() != "test@example.com" {
				t.Errorf("Email() = %v, want %v", p.Email().String(), "test@example.com")
			}
			if p.HomepageURL().String() != "https://example.com" {
				t.Errorf("HomepageURL() = %v, want %v", p.HomepageURL().String(), "https://example.com")
			}
			if p.Note().String() != "テスト備考" {
				t.Errorf("Note() = %v, want %v", p.Note().String(), "テスト備考")
			}
		})
	}
}

func TestUpdateProfileIsImmutable(t *testing.T) {
	t.Parallel()
	tenantID, _ := NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")

	tests := []struct {
		name   string
		mutate func(p Profile, phone *contact.Phone, email *contact.Email, homepageURL *HomepageURL, note *Note)
	}{
		{
			"mutating the update arguments",
			func(p Profile, phone *contact.Phone, email *contact.Email, homepageURL *HomepageURL, note *Note) {
				*phone = contact.Phone("09099999999")
				*email = contact.Email("other@example.com")
				*homepageURL = HomepageURL("https://other.example.com")
				*note = Note("別の備考")
			},
		},
		{
			"mutating the getter results",
			func(p Profile, phone *contact.Phone, email *contact.Email, homepageURL *HomepageURL, note *Note) {
				*p.Phone() = contact.Phone("09099999999")
				*p.Email() = contact.Email("other@example.com")
				*p.HomepageURL() = HomepageURL("https://other.example.com")
				*p.Note() = Note("別の備考")
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			past := time.Now().UTC().Add(-time.Hour)
			p := NewProfile(tenantID, address.Address{}, nil, nil, nil, nil, past, past)

			phone := contact.Phone("0312345678")
			email := contact.Email("test@example.com")
			homepageURL := HomepageURL("https://example.com")
			note := Note("テスト備考")

			p.Update(address.Address{}, &phone, &email, &homepageURL, &note)

			tt.mutate(p, &phone, &email, &homepageURL, &note)

			if p.Phone().String() != "0312345678" {
				t.Errorf("Phone() = %v, want %v", p.Phone().String(), "0312345678")
			}
			if p.Email().String() != "test@example.com" {
				t.Errorf("Email() = %v, want %v", p.Email().String(), "test@example.com")
			}
			if p.HomepageURL().String() != "https://example.com" {
				t.Errorf("HomepageURL() = %v, want %v", p.HomepageURL().String(), "https://example.com")
			}
			if p.Note().String() != "テスト備考" {
				t.Errorf("Note() = %v, want %v", p.Note().String(), "テスト備考")
			}
		})
	}
}
