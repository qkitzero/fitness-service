package training

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const instructionMaxLen = 2000

type Instruction string

func (i Instruction) String() string {
	return string(i)
}

func NewInstruction(s string) (Instruction, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Instruction(""), fmt.Errorf("invalid instruction")
	}
	if utf8.RuneCountInString(s) > instructionMaxLen {
		return Instruction(""), fmt.Errorf("invalid instruction")
	}
	for _, r := range s {
		if r == '\n' || r == '\r' || r == '\t' {
			continue
		}
		if unicode.IsControl(r) {
			return Instruction(""), fmt.Errorf("invalid instruction")
		}
	}
	return Instruction(s), nil
}
