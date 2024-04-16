package main

type Tape struct {
	Cells   []uint8
	Pointer int
}

func (t *Tape) Move(distance int) {
	t.Pointer += distance
	t.Pointer %= len(t.Cells)
}

func (t *Tape) Adjust(offset int) {
	v := int(t.Cells[t.Pointer])
	v += offset
	t.Cells[t.Pointer] = uint8(v)
}

func (t *Tape) Get() uint8 {
	return t.Cells[t.Pointer]
}

func (t *Tape) Set(value uint8) {
	t.Cells[t.Pointer] = value
}
