package tenant

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tenantv1 "github.com/qkitzero/fitness-service/gen/go/tenant/v1"
	apptenant "github.com/qkitzero/fitness-service/internal/application/tenant"
	domainaddress "github.com/qkitzero/fitness-service/internal/domain/address"
	domaincontact "github.com/qkitzero/fitness-service/internal/domain/contact"
	domaintenant "github.com/qkitzero/fitness-service/internal/domain/tenant"
	mocksapptenant "github.com/qkitzero/fitness-service/mocks/application/tenant"
	mockstenant "github.com/qkitzero/fitness-service/mocks/domain/tenant"
)

const sampleTenantID = "0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b"

func TestGetProfile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		tenantID    string
		callUsecase bool
		getErr      error
		wantCode    codes.Code
	}{
		{"success get profile", sampleTenantID, true, nil, codes.OK},
		{"failure invalid tenant id", "", false, nil, codes.InvalidArgument},
		{"failure profile not found", sampleTenantID, true, domaintenant.ErrProfileNotFound, codes.NotFound},
		{"failure not tenant member", sampleTenantID, true, domaintenant.ErrNotMember, codes.PermissionDenied},
		{"failure unauthenticated", sampleTenantID, true, status.Error(codes.Unauthenticated, "unauthenticated"), codes.Unauthenticated},
		{"failure internal error", sampleTenantID, true, errors.New("get profile error"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksapptenant.NewMockProfileUsecase(ctrl)
			if tt.callUsecase {
				var profile domaintenant.Profile
				if tt.getErr == nil {
					mockProfile := mockstenant.NewMockProfile(ctrl)
					mockProfile.EXPECT().TenantID().Return(domaintenant.TenantID(sampleTenantID)).AnyTimes()
					mockProfile.EXPECT().Address().Return(domainaddress.Address{}).AnyTimes()
					mockProfile.EXPECT().Phone().Return(nil).AnyTimes()
					mockProfile.EXPECT().Email().Return(nil).AnyTimes()
					mockProfile.EXPECT().HomepageURL().Return(nil).AnyTimes()
					mockProfile.EXPECT().Note().Return(nil).AnyTimes()
					profile = mockProfile
				}
				mockUsecase.EXPECT().GetProfile(gomock.Any(), domaintenant.TenantID(tt.tenantID)).Return(profile, tt.getErr).Times(1)
			}

			handler := NewProfileHandler(mockUsecase)

			res, err := handler.GetProfile(context.Background(), &tenantv1.GetProfileRequest{TenantId: tt.tenantID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK && res.GetProfile().GetTenantId() != sampleTenantID {
				t.Errorf("TenantId = %v, want %v", res.GetProfile().GetTenantId(), sampleTenantID)
			}
		})
	}
}

func TestUpsertProfile(t *testing.T) {
	t.Parallel()
	postalCode := "123-4567"
	invalidPostalCode := "12-345"
	prefecture := "東京都"
	invalidPrefecture := "存在しない県"
	city := "千代田区"
	invalidCity := strings.Repeat("あ", 256)
	street := "1-1-1"
	invalidStreet := "1-1\n-1"
	building := "テストビル"
	invalidBuilding := "テスト\x00ビル"
	phone := "03-1234-5678"
	invalidPhone := "invalid"
	email := "test@example.com"
	invalidEmail := "invalid"
	homepageURL := "https://example.com"
	invalidHomepageURL := "ftp://example.com"
	note := "テスト備考"
	invalidNote := "テスト\x00備考"
	blank := ""

	wantPostalCode, _ := domainaddress.NewPostalCode(postalCode)
	wantPrefecture, _ := domainaddress.NewPrefecture(prefecture)
	wantCity, _ := domainaddress.NewCity(city)
	wantStreet, _ := domainaddress.NewStreet(street)
	wantBuilding, _ := domainaddress.NewBuilding(building)
	wantPhone, _ := domaincontact.NewPhone(phone)
	wantEmail, _ := domaincontact.NewEmail(email)
	wantHomepageURL, _ := domaintenant.NewHomepageURL(homepageURL)
	wantNote, _ := domaintenant.NewNote(note)

	tests := []struct {
		name        string
		req         *tenantv1.UpsertProfileRequest
		callUsecase bool
		wantPatch   apptenant.ProfilePatch
		upsertErr   error
		wantCode    codes.Code
	}{
		{
			"success upsert profile",
			&tenantv1.UpsertProfileRequest{
				TenantId:    sampleTenantID,
				PostalCode:  &postalCode,
				Prefecture:  &prefecture,
				City:        &city,
				Street:      &street,
				Building:    &building,
				Phone:       &phone,
				Email:       &email,
				HomepageUrl: &homepageURL,
				Note:        &note,
			},
			true,
			apptenant.ProfilePatch{
				PostalCode: wantPostalCode, HasPostalCode: true,
				Prefecture: wantPrefecture, HasPrefecture: true,
				City: wantCity, HasCity: true,
				Street: wantStreet, HasStreet: true,
				Building: wantBuilding, HasBuilding: true,
				Phone: wantPhone, HasPhone: true,
				Email: wantEmail, HasEmail: true,
				HomepageURL: wantHomepageURL, HasHomepageURL: true,
				Note: wantNote, HasNote: true,
			},
			nil, codes.OK,
		},
		{
			"success upsert profile with an empty patch",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID},
			true,
			apptenant.ProfilePatch{},
			nil, codes.OK,
		},
		{
			"success upsert profile with only the note",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID, Note: &note},
			true,
			apptenant.ProfilePatch{Note: wantNote, HasNote: true},
			nil, codes.OK,
		},
		{
			"success upsert profile clearing the note with a blank value",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID, Note: &blank},
			true,
			apptenant.ProfilePatch{Note: nil, HasNote: true},
			nil, codes.OK,
		},
		{
			"failure invalid tenant id",
			&tenantv1.UpsertProfileRequest{TenantId: ""},
			false,
			apptenant.ProfilePatch{},
			nil, codes.InvalidArgument,
		},
		{
			"failure invalid postal code",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID, PostalCode: &invalidPostalCode},
			false,
			apptenant.ProfilePatch{},
			nil, codes.InvalidArgument,
		},
		{
			"failure invalid prefecture",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID, Prefecture: &invalidPrefecture},
			false,
			apptenant.ProfilePatch{},
			nil, codes.InvalidArgument,
		},
		{
			"failure invalid city",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID, City: &invalidCity},
			false,
			apptenant.ProfilePatch{},
			nil, codes.InvalidArgument,
		},
		{
			"failure invalid street",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID, Street: &invalidStreet},
			false,
			apptenant.ProfilePatch{},
			nil, codes.InvalidArgument,
		},
		{
			"failure invalid building",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID, Building: &invalidBuilding},
			false,
			apptenant.ProfilePatch{},
			nil, codes.InvalidArgument,
		},
		{
			"failure invalid phone",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID, Phone: &invalidPhone},
			false,
			apptenant.ProfilePatch{},
			nil, codes.InvalidArgument,
		},
		{
			"failure invalid email",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID, Email: &invalidEmail},
			false,
			apptenant.ProfilePatch{},
			nil, codes.InvalidArgument,
		},
		{
			"failure invalid homepage url",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID, HomepageUrl: &invalidHomepageURL},
			false,
			apptenant.ProfilePatch{},
			nil, codes.InvalidArgument,
		},
		{
			"failure invalid note",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID, Note: &invalidNote},
			false,
			apptenant.ProfilePatch{},
			nil, codes.InvalidArgument,
		},
		{
			"failure profile not found",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID},
			true,
			apptenant.ProfilePatch{},
			domaintenant.ErrProfileNotFound, codes.NotFound,
		},
		{
			"failure not tenant member",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID},
			true,
			apptenant.ProfilePatch{},
			domaintenant.ErrNotMember, codes.PermissionDenied,
		},
		{
			"failure unauthenticated",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID},
			true,
			apptenant.ProfilePatch{},
			status.Error(codes.Unauthenticated, "unauthenticated"), codes.Unauthenticated,
		},
		{
			"failure internal error",
			&tenantv1.UpsertProfileRequest{TenantId: sampleTenantID},
			true,
			apptenant.ProfilePatch{},
			errors.New("upsert profile error"), codes.Internal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksapptenant.NewMockProfileUsecase(ctrl)
			if tt.callUsecase {
				var profile domaintenant.Profile
				if tt.upsertErr == nil {
					mockProfile := mockstenant.NewMockProfile(ctrl)
					mockProfile.EXPECT().TenantID().Return(domaintenant.TenantID(sampleTenantID)).AnyTimes()
					mockProfile.EXPECT().Address().Return(domainaddress.Address{}).AnyTimes()
					mockProfile.EXPECT().Phone().Return(nil).AnyTimes()
					mockProfile.EXPECT().Email().Return(nil).AnyTimes()
					mockProfile.EXPECT().HomepageURL().Return(nil).AnyTimes()
					mockProfile.EXPECT().Note().Return(nil).AnyTimes()
					profile = mockProfile
				}
				mockUsecase.EXPECT().UpsertProfile(gomock.Any(), domaintenant.TenantID(sampleTenantID), tt.wantPatch).Return(profile, tt.upsertErr).Times(1)
			}

			handler := NewProfileHandler(mockUsecase)

			res, err := handler.UpsertProfile(context.Background(), tt.req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK && res.GetProfile().GetTenantId() != sampleTenantID {
				t.Errorf("TenantId = %v, want %v", res.GetProfile().GetTenantId(), sampleTenantID)
			}
		})
	}
}

func TestToProtoProfile(t *testing.T) {
	t.Parallel()
	tenantID, _ := domaintenant.NewTenantID(sampleTenantID)
	postalCode, _ := domainaddress.NewPostalCode("123-4567")
	prefecture, _ := domainaddress.NewPrefecture("東京都")
	city, _ := domainaddress.NewCity("千代田区")
	street, _ := domainaddress.NewStreet("1-1-1")
	building, _ := domainaddress.NewBuilding("テストビル")
	addr := domainaddress.NewAddress(postalCode, prefecture, city, street, building)
	phone, _ := domaincontact.NewPhone("03-1234-5678")
	email, _ := domaincontact.NewEmail("test@example.com")
	homepageURL, _ := domaintenant.NewHomepageURL("https://example.com")
	note, _ := domaintenant.NewNote("テスト備考")
	now := time.Now().UTC()

	tests := []struct {
		name            string
		profile         domaintenant.Profile
		wantPostalCode  string
		wantPrefecture  string
		wantCity        string
		wantStreet      string
		wantBuilding    string
		wantPhone       string
		wantEmail       string
		wantHomepageURL string
		wantNote        string
	}{
		{
			"success to proto profile",
			domaintenant.NewProfile(tenantID, addr, phone, email, homepageURL, note, now, now),
			"1234567", "東京都", "千代田区", "1-1-1", "テストビル", "0312345678", "test@example.com", "https://example.com", "テスト備考",
		},
		{
			"success to proto profile without optional fields",
			domaintenant.NewProfile(tenantID, domainaddress.Address{}, nil, nil, nil, nil, now, now),
			"", "", "", "", "", "", "", "", "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := toProtoProfile(tt.profile)

			if got.GetTenantId() != sampleTenantID {
				t.Errorf("TenantId = %v, want %v", got.GetTenantId(), sampleTenantID)
			}
			if tt.wantPostalCode == "" && got.PostalCode != nil {
				t.Errorf("PostalCode = %v, want nil", got.PostalCode)
			}
			if tt.wantPostalCode != "" && got.GetPostalCode() != tt.wantPostalCode {
				t.Errorf("PostalCode = %v, want %v", got.GetPostalCode(), tt.wantPostalCode)
			}
			if tt.wantPrefecture == "" && got.Prefecture != nil {
				t.Errorf("Prefecture = %v, want nil", got.Prefecture)
			}
			if tt.wantPrefecture != "" && got.GetPrefecture() != tt.wantPrefecture {
				t.Errorf("Prefecture = %v, want %v", got.GetPrefecture(), tt.wantPrefecture)
			}
			if tt.wantCity == "" && got.City != nil {
				t.Errorf("City = %v, want nil", got.City)
			}
			if tt.wantCity != "" && got.GetCity() != tt.wantCity {
				t.Errorf("City = %v, want %v", got.GetCity(), tt.wantCity)
			}
			if tt.wantStreet == "" && got.Street != nil {
				t.Errorf("Street = %v, want nil", got.Street)
			}
			if tt.wantStreet != "" && got.GetStreet() != tt.wantStreet {
				t.Errorf("Street = %v, want %v", got.GetStreet(), tt.wantStreet)
			}
			if tt.wantBuilding == "" && got.Building != nil {
				t.Errorf("Building = %v, want nil", got.Building)
			}
			if tt.wantBuilding != "" && got.GetBuilding() != tt.wantBuilding {
				t.Errorf("Building = %v, want %v", got.GetBuilding(), tt.wantBuilding)
			}
			if tt.wantPhone == "" && got.Phone != nil {
				t.Errorf("Phone = %v, want nil", got.Phone)
			}
			if tt.wantPhone != "" && got.GetPhone() != tt.wantPhone {
				t.Errorf("Phone = %v, want %v", got.GetPhone(), tt.wantPhone)
			}
			if tt.wantEmail == "" && got.Email != nil {
				t.Errorf("Email = %v, want nil", got.Email)
			}
			if tt.wantEmail != "" && got.GetEmail() != tt.wantEmail {
				t.Errorf("Email = %v, want %v", got.GetEmail(), tt.wantEmail)
			}
			if tt.wantHomepageURL == "" && got.HomepageUrl != nil {
				t.Errorf("HomepageUrl = %v, want nil", got.HomepageUrl)
			}
			if tt.wantHomepageURL != "" && got.GetHomepageUrl() != tt.wantHomepageURL {
				t.Errorf("HomepageUrl = %v, want %v", got.GetHomepageUrl(), tt.wantHomepageURL)
			}
			if tt.wantNote == "" && got.Note != nil {
				t.Errorf("Note = %v, want nil", got.Note)
			}
			if tt.wantNote != "" && got.GetNote() != tt.wantNote {
				t.Errorf("Note = %v, want %v", got.GetNote(), tt.wantNote)
			}
		})
	}
}
