package main

import (
	"slices"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	for num, sub := range []struct {
		in  string
		bf  bool
		ok  bool
		exp []Opcode
	}{
		{"", false, true, nil},
		{"", true, true, nil},
		{"#", false, false, nil},
		{"#", true, true, nil},
		{"+", false, true, []Opcode{OpcodeIncrementOne}},
		{"+", true, true, []Opcode{OpcodeIncrementOne}},
		{"++", false, true, []Opcode{OpcodeIncrementTwo}},
		{"++", true, true, []Opcode{OpcodeIncrementTwo}},
		{"+++", false, true, []Opcode{OpcodeIncrementTwo, OpcodeIncrementOne}},
		{"-", false, true, []Opcode{OpcodeDecrementOne}},
		{"--", false, true, []Opcode{OpcodeDecrementTwo}},
		{"---", false, true, []Opcode{OpcodeDecrementTwo, OpcodeDecrementOne}},
	} {
		p, err := Parse(strings.NewReader(sub.in), 0, sub.bf, false)
		if sub.ok && err != nil {
			t.Errorf("Parse failed in subtest %d: %v", num, err)
			continue
		} else if !sub.ok && err == nil {
			t.Errorf("Parse succeeded unexpectedly in subtest %d", num)
			continue
		}

		if err != nil {
			continue
		}

		if p == nil {
			t.Errorf("Parse succeeded but program is nil in subtest %d", num)
		} else if !slices.Equal(p.Opcodes, sub.exp) {
			t.Errorf("Parse result %v != %v in subtest %d", p.Opcodes, sub.exp, num)
		}
	}
}
