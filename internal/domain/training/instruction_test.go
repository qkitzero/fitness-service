package training

import (
	"strings"
	"testing"
)

func TestNewInstruction(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		success     bool
		instruction string
		want        Instruction
	}{
		{"success new instruction", true, "壁に手をついて肘を曲げ伸ばしする。", "壁に手をついて肘を曲げ伸ばしする。"},
		{"success newline instruction", true, "壁に手をつく。\n肘を曲げ伸ばしする。", "壁に手をつく。\n肘を曲げ伸ばしする。"},
		{"success carriage return instruction", true, "壁に手をつく。\r\n肘を曲げ伸ばしする。", "壁に手をつく。\r\n肘を曲げ伸ばしする。"},
		{"success tab instruction", true, "壁に手をつく。\t肘を曲げ伸ばしする。", "壁に手をつく。\t肘を曲げ伸ばしする。"},
		{"success trim space", true, "  肘を曲げ伸ばしする。  ", "肘を曲げ伸ばしする。"},
		{"success max length instruction", true, strings.Repeat("a", 2000), Instruction(strings.Repeat("a", 2000))},
		{"failure empty instruction", false, "", ""},
		{"failure whitespace instruction", false, "   ", ""},
		{"failure too long instruction", false, strings.Repeat("a", 2001), ""},
		{"failure null character instruction", false, "壁に\x00手をつく。", ""},
		{"failure bell character instruction", false, "壁に\a手をつく。", ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			instruction, err := NewInstruction(tt.instruction)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && instruction != tt.want {
				t.Errorf("NewInstruction() = %v, want %v", instruction, tt.want)
			}

			if tt.success && instruction.String() != tt.want.String() {
				t.Errorf("String() = %v, want %v", instruction.String(), tt.want.String())
			}
		})
	}
}
