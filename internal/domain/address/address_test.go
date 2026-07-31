package address

import (
	"reflect"
	"testing"
)

func TestNewAddress(t *testing.T) {
	t.Parallel()
	postalCode := PostalCode("1234567")
	prefecture := Prefecture("東京都")
	city := City("千代田区")
	street := Street("1-1-1")
	building := Building("テストビル")

	tests := []struct {
		name       string
		postalCode *PostalCode
		prefecture *Prefecture
		city       *City
		street     *Street
		building   *Building
	}{
		{"success", &postalCode, &prefecture, &city, &street, &building},
		{"success all nil", nil, nil, nil, nil, nil},
		{"success partial", &postalCode, nil, &city, nil, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewAddress(tt.postalCode, tt.prefecture, tt.city, tt.street, tt.building)

			want := Address{postalCode: tt.postalCode, prefecture: tt.prefecture, city: tt.city, street: tt.street, building: tt.building}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestAddressGettersReturnNilForUnsetFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		isNil func(a Address) bool
	}{
		{"postal code", func(a Address) bool { return a.PostalCode() == nil }},
		{"prefecture", func(a Address) bool { return a.Prefecture() == nil }},
		{"city", func(a Address) bool { return a.City() == nil }},
		{"street", func(a Address) bool { return a.Street() == nil }},
		{"building", func(a Address) bool { return a.Building() == nil }},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got Address

			if !tt.isNil(got) {
				t.Errorf("expected nil, but got a value")
			}
		})
	}
}

func TestAddressIsImmutable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		mutate func(got Address, postalCode *PostalCode, prefecture *Prefecture, city *City, street *Street, building *Building)
	}{
		{
			"mutating the constructor arguments",
			func(got Address, postalCode *PostalCode, prefecture *Prefecture, city *City, street *Street, building *Building) {
				*postalCode = PostalCode("7654321")
				*prefecture = Prefecture("大阪府")
				*city = City("北区")
				*street = Street("2-2-2")
				*building = Building("別のビル")
			},
		},
		{
			"mutating the getter results",
			func(got Address, postalCode *PostalCode, prefecture *Prefecture, city *City, street *Street, building *Building) {
				*got.PostalCode() = PostalCode("7654321")
				*got.Prefecture() = Prefecture("大阪府")
				*got.City() = City("北区")
				*got.Street() = Street("2-2-2")
				*got.Building() = Building("別のビル")
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			postalCode := PostalCode("1234567")
			prefecture := Prefecture("東京都")
			city := City("千代田区")
			street := Street("1-1-1")
			building := Building("テストビル")

			got := NewAddress(&postalCode, &prefecture, &city, &street, &building)

			tt.mutate(got, &postalCode, &prefecture, &city, &street, &building)

			if got.PostalCode().String() != "1234567" {
				t.Errorf("PostalCode() = %v, want %v", got.PostalCode().String(), "1234567")
			}
			if got.Prefecture().String() != "東京都" {
				t.Errorf("Prefecture() = %v, want %v", got.Prefecture().String(), "東京都")
			}
			if got.City().String() != "千代田区" {
				t.Errorf("City() = %v, want %v", got.City().String(), "千代田区")
			}
			if got.Street().String() != "1-1-1" {
				t.Errorf("Street() = %v, want %v", got.Street().String(), "1-1-1")
			}
			if got.Building().String() != "テストビル" {
				t.Errorf("Building() = %v, want %v", got.Building().String(), "テストビル")
			}
		})
	}
}
