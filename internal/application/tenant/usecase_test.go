package tenant

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/contact"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
	mocksappauth "github.com/qkitzero/fitness-service/mocks/application/auth"
	mocksappuser "github.com/qkitzero/fitness-service/mocks/application/user"
	mockstenant "github.com/qkitzero/fitness-service/mocks/domain/tenant"
)

func TestGetProfile(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")

	tests := []struct {
		name               string
		success            bool
		wantErr            error
		callFindByTenantID bool
		ctx                context.Context
		userID             string
		verifyTokenErr     error
		myTenantIDs        []string
		listMyGroupsErr    error
		findByTenantIDErr  error
	}{
		{"success get profile", true, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{tenantID.String()}, nil, nil},
		{"failure verify token error", false, nil, false, context.Background(), "", errors.New("verify token error"), []string{tenantID.String()}, nil, nil},
		{"failure other tenant is hidden as not found", false, tenant.ErrProfileNotFound, false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil},
		{"failure list my groups error", false, nil, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, errors.New("list my groups error"), nil},
		{"failure profile not found", false, tenant.ErrProfileNotFound, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{tenantID.String()}, nil, tenant.ErrProfileNotFound},
		{"failure find by tenant id error", false, nil, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{tenantID.String()}, nil, errors.New("find by tenant id error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockProfile := mockstenant.NewMockProfile(ctrl)
			mockProfileRepository := mockstenant.NewMockProfileRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindByTenantID {
				var foundProfile tenant.Profile
				if tt.findByTenantIDErr == nil {
					foundProfile = mockProfile
				}
				mockProfileRepository.EXPECT().FindByTenantID(tt.ctx, tenantID).Return(foundProfile, tt.findByTenantIDErr).Times(1)
			}

			u := NewProfileUsecase(mockAuthService, mockUserService, mockProfileRepository)

			foundProfile, err := u.GetProfile(tt.ctx, tenantID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.success && foundProfile != mockProfile {
				t.Errorf("expected the profile returned by the repository")
			}
		})
	}
}

func TestUpsertProfile(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")

	tests := []struct {
		name               string
		success            bool
		wantErr            error
		callFindByTenantID bool
		callUpsert         bool
		ctx                context.Context
		userID             string
		verifyTokenErr     error
		myTenantIDs        []string
		listMyGroupsErr    error
		findByTenantIDErr  error
		upsertErr          error
	}{
		{"success insert profile", true, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{tenantID.String()}, nil, tenant.ErrProfileNotFound, nil},
		{"success update profile", true, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{tenantID.String()}, nil, nil, nil},
		{"failure verify token error", false, nil, false, false, context.Background(), "", errors.New("verify token error"), []string{tenantID.String()}, nil, nil, nil},
		{"failure other tenant is hidden as not found", false, tenant.ErrProfileNotFound, false, false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{"9a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"}, nil, nil, nil},
		{"failure list my groups error", false, nil, false, false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, errors.New("list my groups error"), nil, nil},
		{"failure find by tenant id error", false, nil, true, false, context.Background(), "google-oauth2|000000000000000000000", nil, []string{tenantID.String()}, nil, errors.New("find by tenant id error"), nil},
		{"failure upsert error on insert", false, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{tenantID.String()}, nil, tenant.ErrProfileNotFound, errors.New("upsert error")},
		{"failure upsert error on update", false, nil, true, true, context.Background(), "google-oauth2|000000000000000000000", nil, []string{tenantID.String()}, nil, nil, errors.New("upsert error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			now := time.Now().UTC()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockProfileRepository := mockstenant.NewMockProfileRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(tt.ctx).Return(tt.myTenantIDs, tt.listMyGroupsErr).AnyTimes()
			if tt.callFindByTenantID {
				var foundProfile tenant.Profile
				if tt.findByTenantIDErr == nil {
					foundProfile = tenant.NewProfile(tenantID, address.Address{}, nil, nil, nil, nil, now, now)
				}
				mockProfileRepository.EXPECT().FindByTenantID(tt.ctx, tenantID).Return(foundProfile, tt.findByTenantIDErr).Times(1)
			}
			if tt.callUpsert {
				mockProfileRepository.EXPECT().Upsert(tt.ctx, gomock.Any()).Return(tt.upsertErr).Times(1)
			}

			u := NewProfileUsecase(mockAuthService, mockUserService, mockProfileRepository)

			upsertedProfile, err := u.UpsertProfile(tt.ctx, tenantID, ProfilePatch{})
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.success && upsertedProfile.TenantID() != tenantID {
				t.Errorf("TenantID() = %v, want %v", upsertedProfile.TenantID(), tenantID)
			}
		})
	}
}

func TestUpsertProfileAppliesOnlyTheGivenFields(t *testing.T) {
	t.Parallel()
	tenantID, _ := tenant.NewTenantID("0f4a1a2b-3c4d-5e6f-7a8b-9c0d1e2f3a4b")
	existingPostalCode, _ := address.NewPostalCode("123-4567")
	existingPrefecture, _ := address.NewPrefecture("東京都")
	existingCity, _ := address.NewCity("千代田区")
	existingStreet, _ := address.NewStreet("1-1-1")
	existingBuilding, _ := address.NewBuilding("テストビル")
	existingAddress := address.NewAddress(existingPostalCode, existingPrefecture, existingCity, existingStreet, existingBuilding)
	existingPhone, _ := contact.NewPhone("03-1234-5678")
	existingEmail, _ := contact.NewEmail("test@example.com")
	existingHomepageURL, _ := tenant.NewHomepageURL("https://example.com")
	existingNote, _ := tenant.NewNote("テスト備考")

	updatedPostalCode, _ := address.NewPostalCode("543-2100")
	updatedPrefecture, _ := address.NewPrefecture("大阪府")
	updatedCity, _ := address.NewCity("大阪市")
	updatedStreet, _ := address.NewStreet("2-2-2")
	updatedBuilding, _ := address.NewBuilding("更新ビル")
	updatedPhone, _ := contact.NewPhone("080-1234-5678")
	updatedEmail, _ := contact.NewEmail("updated@example.com")
	updatedHomepageURL, _ := tenant.NewHomepageURL("https://updated.example.com")
	updatedNote, _ := tenant.NewNote("改装中")
	clearedNote, _ := tenant.NewNote("")

	tests := []struct {
		name            string
		patch           ProfilePatch
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
			"an empty patch keeps every existing field",
			ProfilePatch{},
			"1234567", "東京都", "千代田区", "1-1-1", "テストビル", "0312345678", "test@example.com", "https://example.com", "テスト備考",
		},
		{
			"a note-only patch keeps every other field",
			ProfilePatch{Note: updatedNote, HasNote: true},
			"1234567", "東京都", "千代田区", "1-1-1", "テストビル", "0312345678", "test@example.com", "https://example.com", "改装中",
		},
		{
			"a phone-only patch keeps every other field",
			ProfilePatch{Phone: updatedPhone, HasPhone: true},
			"1234567", "東京都", "千代田区", "1-1-1", "テストビル", "08012345678", "test@example.com", "https://example.com", "テスト備考",
		},
		{
			"an explicitly blank note clears only the note",
			ProfilePatch{Note: clearedNote, HasNote: true},
			"1234567", "東京都", "千代田区", "1-1-1", "テストビル", "0312345678", "test@example.com", "https://example.com", "",
		},
		{
			"a city-only patch keeps the other address components",
			ProfilePatch{City: existingCity, HasCity: true},
			"1234567", "東京都", "千代田区", "1-1-1", "テストビル", "0312345678", "test@example.com", "https://example.com", "テスト備考",
		},
		{
			"a full patch replaces every field",
			ProfilePatch{
				PostalCode: updatedPostalCode, HasPostalCode: true,
				Prefecture: updatedPrefecture, HasPrefecture: true,
				City: updatedCity, HasCity: true,
				Street: updatedStreet, HasStreet: true,
				Building: updatedBuilding, HasBuilding: true,
				Phone: updatedPhone, HasPhone: true,
				Email: updatedEmail, HasEmail: true,
				HomepageURL: updatedHomepageURL, HasHomepageURL: true,
				Note: updatedNote, HasNote: true,
			},
			"5432100", "大阪府", "大阪市", "2-2-2", "更新ビル", "08012345678", "updated@example.com", "https://updated.example.com", "改装中",
		},
		{
			"a full patch of blank values clears every field",
			ProfilePatch{
				HasPostalCode: true, HasPrefecture: true, HasCity: true, HasStreet: true, HasBuilding: true,
				HasPhone: true, HasEmail: true, HasHomepageURL: true, HasNote: true,
			},
			"", "", "", "", "", "", "", "", "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			now := time.Now().UTC()
			existingProfile := tenant.NewProfile(tenantID, existingAddress, existingPhone, existingEmail, existingHomepageURL, existingNote, now, now)

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockProfileRepository := mockstenant.NewMockProfileRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(ctx).Return("google-oauth2|000000000000000000000", nil).AnyTimes()
			mockUserService.EXPECT().ListMyGroups(ctx).Return([]string{tenantID.String()}, nil).AnyTimes()
			mockProfileRepository.EXPECT().FindByTenantID(ctx, tenantID).Return(existingProfile, nil).Times(1)
			mockProfileRepository.EXPECT().Upsert(ctx, gomock.Any()).Return(nil).Times(1)

			u := NewProfileUsecase(mockAuthService, mockUserService, mockProfileRepository)

			got, err := u.UpsertProfile(ctx, tenantID, tt.patch)
			if err != nil {
				t.Fatalf("expected no error, but got %v", err)
			}

			addr := got.Address()
			if tt.wantPostalCode == "" && addr.PostalCode() != nil {
				t.Errorf("Address().PostalCode() = %v, want nil", addr.PostalCode())
			}
			if tt.wantPostalCode != "" && (addr.PostalCode() == nil || addr.PostalCode().String() != tt.wantPostalCode) {
				t.Errorf("Address().PostalCode() = %v, want %v", addr.PostalCode(), tt.wantPostalCode)
			}
			if tt.wantPrefecture == "" && addr.Prefecture() != nil {
				t.Errorf("Address().Prefecture() = %v, want nil", addr.Prefecture())
			}
			if tt.wantPrefecture != "" && (addr.Prefecture() == nil || addr.Prefecture().String() != tt.wantPrefecture) {
				t.Errorf("Address().Prefecture() = %v, want %v", addr.Prefecture(), tt.wantPrefecture)
			}
			if tt.wantCity == "" && addr.City() != nil {
				t.Errorf("Address().City() = %v, want nil", addr.City())
			}
			if tt.wantCity != "" && (addr.City() == nil || addr.City().String() != tt.wantCity) {
				t.Errorf("Address().City() = %v, want %v", addr.City(), tt.wantCity)
			}
			if tt.wantStreet == "" && addr.Street() != nil {
				t.Errorf("Address().Street() = %v, want nil", addr.Street())
			}
			if tt.wantStreet != "" && (addr.Street() == nil || addr.Street().String() != tt.wantStreet) {
				t.Errorf("Address().Street() = %v, want %v", addr.Street(), tt.wantStreet)
			}
			if tt.wantBuilding == "" && addr.Building() != nil {
				t.Errorf("Address().Building() = %v, want nil", addr.Building())
			}
			if tt.wantBuilding != "" && (addr.Building() == nil || addr.Building().String() != tt.wantBuilding) {
				t.Errorf("Address().Building() = %v, want %v", addr.Building(), tt.wantBuilding)
			}
			if tt.wantPhone == "" && got.Phone() != nil {
				t.Errorf("Phone() = %v, want nil", got.Phone())
			}
			if tt.wantPhone != "" && (got.Phone() == nil || got.Phone().String() != tt.wantPhone) {
				t.Errorf("Phone() = %v, want %v", got.Phone(), tt.wantPhone)
			}
			if tt.wantEmail == "" && got.Email() != nil {
				t.Errorf("Email() = %v, want nil", got.Email())
			}
			if tt.wantEmail != "" && (got.Email() == nil || got.Email().String() != tt.wantEmail) {
				t.Errorf("Email() = %v, want %v", got.Email(), tt.wantEmail)
			}
			if tt.wantHomepageURL == "" && got.HomepageURL() != nil {
				t.Errorf("HomepageURL() = %v, want nil", got.HomepageURL())
			}
			if tt.wantHomepageURL != "" && (got.HomepageURL() == nil || got.HomepageURL().String() != tt.wantHomepageURL) {
				t.Errorf("HomepageURL() = %v, want %v", got.HomepageURL(), tt.wantHomepageURL)
			}
			if tt.wantNote == "" && got.Note() != nil {
				t.Errorf("Note() = %v, want nil", got.Note())
			}
			if tt.wantNote != "" && (got.Note() == nil || got.Note().String() != tt.wantNote) {
				t.Errorf("Note() = %v, want %v", got.Note(), tt.wantNote)
			}
		})
	}
}
