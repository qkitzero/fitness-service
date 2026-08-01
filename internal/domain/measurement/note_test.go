package measurement

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
		{"success", true, "膝の痛みあり", "膝の痛みあり"},
		{"success blank", true, "  ", ""},
		{"success with newline", true, "1回目: ふらつき\n2回目: 安定", "1回目: ふらつき\n2回目: 安定"},
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
