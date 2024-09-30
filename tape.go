package main

import (
	"strconv"
	"strings"
)

type Tape struct {
	Cells   []int8
	Pointer int
}

func (t *Tape) Move(distance int) {
	s := len(t.Cells)
	t.Pointer += distance
	t.Pointer = (t.Pointer%s + s) % s
}

func (t *Tape) Adjust(offset int) {
	v := int(t.Cells[t.Pointer])
	v += offset
	t.Cells[t.Pointer] = int8(v)
}

func (t *Tape) Get() int8 {
	return t.Cells[t.Pointer]
}

func (t *Tape) Set(value int8) {
	t.Cells[t.Pointer] = value
}

func (t *Tape) String(signed bool) string {
	var b strings.Builder
	b.WriteRune('(')
	b.WriteString(strconv.Itoa(t.Pointer))
	b.WriteRune(')')

	b.WriteRune('[')

	for i, c := range t.Cells {
		if i > 0 {
			b.WriteRune(',')
		}

		b.WriteString(strconv.Itoa(func() int {
			if signed {
				return int(c)
			}

			return int(uint8(c))
		}()))
	}

	b.WriteRune(']')
	return b.String()
}
