package customer

import "testing"

func TestNewNameKana(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		success  bool
		nameKana string
		want     string
	}{
		{"success katakana", true, "テストカナ", "テストカナ"},
		{"success with space", true, "ヤマダ タロウ", "ヤマダ タロウ"},
		{"success with middle dot", true, "マイケル・ジャクソン", "マイケル・ジャクソン"},
		{"success with long mark", true, "コーヒー", "コーヒー"},
		{"success trim outer space", true, "  テストカナ  ", "テストカナ"},
		{"failure empty", false, "", ""},
		{"failure whitespace only", false, "   ", ""},
		{"failure kanji", false, "山田太郎", ""},
		{"failure romaji", false, "Yamada", ""},
		{"failure hiragana", false, "やまだ", ""},
		{"failure digits and emoji", false, "123😀", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			nameKana, err := NewNameKana(tt.nameKana)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && nameKana.String() != tt.want {
				t.Errorf("String() = %v, want %v", nameKana.String(), tt.want)
			}
		})
	}
}
