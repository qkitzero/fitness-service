package tenant

import (
	"strings"
	"testing"
)

func TestNewNote(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		input   string
		want    string
	}{
		{"success", true, "テスト備考", "テスト備考"},
		{"success blank", true, "  ", ""},
		{"success with newline", true, "営業時間: 10-22\n定休日: 月曜", "営業時間: 10-22\n定休日: 月曜"},
		{"success with carriage return", true, "line1\r\nline2", "line1\r\nline2"},
		{"success with tab", true, "項目\t値", "項目\t値"},
		{"success max length", true, strings.Repeat("あ", 255), strings.Repeat("あ", 255)},
		{"failure too long", false, strings.Repeat("あ", 256), ""},
		{"failure null character", false, "テスト\x00備考", ""},
		{"failure bell character", false, "テスト\a備考", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewNote(tt.input)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && tt.want == "" && got != nil {
				t.Errorf("expected nil, but got %v", *got)
			}
			if tt.success && tt.want != "" && (got == nil || got.String() != tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
