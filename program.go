package main

import (
	"fmt"
	"io"
	"math"
	"math/rand"
)

type Program struct {
	Opcodes []Opcode
	Tape    Tape
	Random  *rand.Rand
	Pointer int
	Stopped bool
	Signed  bool
}

func (p *Program) Skip(marker uint) {
	for p.Pointer+1 < len(p.Opcodes) {
		p.Pointer++
		op := p.Opcodes[p.Pointer]
		if op.Marker() != marker || op.Skips() {
			continue
		}
		break
	}
}

func (p *Program) Repeat(marker uint) {
	for p.Pointer > 0 {
		p.Pointer--
		op := p.Opcodes[p.Pointer]
		if op.Marker() != marker || op.Repeats() {
			continue
		}
		break
	}
}

func (p *Program) Write(w io.Writer) (err error) {
	b := make([]byte, (len(p.Opcodes)+1)/2)
	for i, o := range p.Opcodes {
		b[i/2] |= byte(o << ((i % 2) * 4))
	}

	_, err = w.Write(b)
	return
}

func (p *Program) Explain(w io.Writer) (err error) {
	for i, o := range p.Opcodes {
		_, err = fmt.Fprintf(w, "%x%c%s %s\n", i/2, func() rune {
			if i%2 == 0 {
				return 'l'
			}
			return 'h'
		}(), o.Token(), o.Description())
		if err != nil {
			return
		}
	}
	return
}

func (p *Program) Running() bool {
	return !p.Stopped && p.Pointer >= 0 && p.Pointer < len(p.Opcodes)
}

func (p *Program) Run() {
	if !p.Running() {
		return
	}

	switch p.Opcodes[p.Pointer] {
	case OpcodeOutput:
		p.Stopped = true
	case OpcodeInput:
		var r int
		if p.Random != nil {
			r = p.Random.Int()
		} else {
			r = rand.Int()
		}
		p.Tape.Set(int8(r & math.MaxUint8))
	case OpcodeLeft:
		p.Tape.Move(-1)
	case OpcodeRight:
		p.Tape.Move(1)
	case OpcodeDivide:
		r := p.Tape.Get() >> 1
		if !p.Signed {
			r &= 0x7f
		}
		p.Tape.Set(r)
	case OpcodeDecrementTwo:
		p.Tape.Adjust(-2)
	case OpcodeDecrementOne:
		p.Tape.Adjust(-1)
	case OpcodeIncrementOne:
		p.Tape.Adjust(1)
	case OpcodeIncrementTwo:
		p.Tape.Adjust(2)
	case OpcodeMultiply:
		p.Tape.Set(p.Tape.Get() << 1)
	case OpcodeSkipOne:
		if p.Tape.Get() == 0 || p.Signed && p.Tape.Get() < 0 {
			p.Skip(1)
		}
	case OpcodeSkipTwo:
		if p.Tape.Get() == 0 || p.Signed && p.Tape.Get() < 0 {
			p.Skip(2)
		}
	case OpcodeSkipThree:
		if p.Tape.Get() == 0 || p.Signed && p.Tape.Get() < 0 {
			p.Skip(3)
		}
	case OpcodeRepeatOne:
		if p.Tape.Get() != 0 {
			p.Repeat(1)
		}
	case OpcodeRepeatTwo:
		if p.Tape.Get() != 0 {
			p.Repeat(2)
		}
	case OpcodeRepeatThree:
		if p.Tape.Get() != 0 {
			p.Repeat(3)
		}
	}

	p.Pointer++
}
