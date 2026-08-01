package staff

import (
	"strings"
	"testing"
)

func TestNewStaffID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		staffID string
		want    string
	}{
		{"success new staff id", true, "google-oauth2|000000000000000000000", "google-oauth2|000000000000000000000"},
		{"success trim whitespace", true, "  google-oauth2|000000000000000000000  ", "google-oauth2|000000000000000000000"},
		{"success upstream format is not validated", true, "user_00000000", "user_00000000"},
		{"success long identity provider subject", true, "waad|ExXKMSuKtbEE-JhoHiKZOgHnT-tt2rTLSl0RcJ1lVpQ", "waad|ExXKMSuKtbEE-JhoHiKZOgHnT-tt2rTLSl0RcJ1lVpQ"},
		{"success max length staff id", true, strings.Repeat("a", 255), strings.Repeat("a", 255)},
		{"failure empty staff id", false, "", ""},
		{"failure whitespace staff id", false, "   ", ""},
		{"failure too long staff id", false, strings.Repeat("a", 256), ""},
		{"failure control character", false, "google-oauth2\x00", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			staffID, err := NewStaffID(tt.staffID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && staffID.String() != tt.want {
				t.Errorf("String() = %v, want %v", staffID.String(), tt.want)
			}
		})
	}
}
