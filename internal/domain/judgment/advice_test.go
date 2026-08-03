package judgment

import (
	"strings"
	"testing"
)

func TestNewAdvice(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		advice  string
		wantNil bool
		want    string
	}{
		{"success new advice", true, "週2回のスクワットを継続してください", false, "週2回のスクワットを継続してください"},
		{"success multiline advice", true, "下肢筋力\n・スクワット 10回×3セット\n・階段昇降", false, "下肢筋力\n・スクワット 10回×3セット\n・階段昇降"},
		{"success tab in advice", true, "握力\tB", false, "握力\tB"},
		{"success trim surrounding space", true, "  週2回のスクワット  ", false, "週2回のスクワット"},
		{"success max length advice", true, strings.Repeat("あ", 2000), false, strings.Repeat("あ", 2000)},
		{"success empty advice is nil", true, "", true, ""},
		{"success whitespace advice is nil", true, "   ", true, ""},
		{"success newline only advice is nil", true, "\n\n", true, ""},
		{"failure too long advice", false, strings.Repeat("あ", 2001), true, ""},
		{"failure control character advice", false, "週2回の\x00スクワット", true, ""},
		{"failure bell character advice", false, "週2回の\aスクワット", true, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			advice, err := NewAdvice(tt.advice)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.wantNil {
				if advice != nil {
					t.Errorf("NewAdvice() = %v, want nil", *advice)
				}
				return
			}

			if advice == nil {
				t.Fatalf("NewAdvice() = nil, want %v", tt.want)
			}
			if advice.String() != tt.want {
				t.Errorf("String() = %v, want %v", advice.String(), tt.want)
			}
		})
	}
}
