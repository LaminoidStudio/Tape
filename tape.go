package main

type Tape struct {
	Cells   []uint8
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
	t.Cells[t.Pointer] = uint8(v)
}

func (t *Tape) Get() uint8 {
	return t.Cells[t.Pointer]
}

func (t *Tape) Set(value uint8) {
	t.Cells[t.Pointer] = value
}
