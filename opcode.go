package main

type Opcode uint

const (
	// OpcodeOutput stops the program on the current cell and blanks the rest
	OpcodeOutput Opcode = iota
	// OpcodeLeft moves left on the tape
	OpcodeLeft
	// OpcodeRight moves right on the tape
	OpcodeRight
	// OpcodeDivide divides the current cell value by 2
	OpcodeDivide
	// OpcodeDecrementTwo decrements the current cell value by 2
	OpcodeDecrementTwo
	// OpcodeDecrementOne decrements the current cell value by 1
	OpcodeDecrementOne
	// OpcodeInput stores a random value into the current cell
	OpcodeInput
	// OpcodeIncrementOne increments the current cell value by 1
	OpcodeIncrementOne
	// OpcodeIncrementTwo increments the current cell value by 2
	OpcodeIncrementTwo
	// OpcodeMultiply multiplies the current cell value by 2
	OpcodeMultiply
	// OpcodeRepeatOne seeks left to 0
	OpcodeRepeatOne
	// OpcodeRepeatTwo seeks left to 1
	OpcodeRepeatTwo
	// OpcodeRepeatThree seeks left to 2
	OpcodeRepeatThree
	// OpcodeSkipOne seeks right to 0
	OpcodeSkipOne
	// OpcodeSkipTwo seeks right to 1
	OpcodeSkipTwo
	// OpcodeSkipThree seeks right to 2
	OpcodeSkipThree
)

func (o Opcode) Token() string {
	switch o {
	case OpcodeOutput:
		return string(TokenOutput)
	case OpcodeLeft:
		return string(TokenLeft)
	case OpcodeRight:
		return string(TokenRight)
	case OpcodeDivide:
		return string(TokenDivide)
	case OpcodeDecrementTwo:
		return string(TokenDecrement) + string(TokenDecrement)
	case OpcodeDecrementOne:
		return string(TokenDecrement)
	case OpcodeInput:
		return string(TokenInput)
	case OpcodeIncrementOne:
		return string(TokenIncrement)
	case OpcodeIncrementTwo:
		return string(TokenIncrement) + string(TokenIncrement)
	case OpcodeMultiply:
		return string(TokenMultiply)
	case OpcodeRepeatOne:
		return string(TokenRepeatOne)
	case OpcodeRepeatTwo:
		return string(TokenRepeatTwo)
	case OpcodeRepeatThree:
		return string(TokenRepeatThree)
	case OpcodeSkipOne:
		return string(TokenSkipOne)
	case OpcodeSkipTwo:
		return string(TokenSkipTwo)
	case OpcodeSkipThree:
		return string(TokenSkipThree)
	default:
		return ""
	}
}

func (o Opcode) Description() string {
	switch o {
	case OpcodeOutput:
		return "output and stop"
	case OpcodeLeft:
		return "move left"
	case OpcodeRight:
		return "move right"
	case OpcodeDivide:
		return "divide by two"
	case OpcodeDecrementTwo:
		return "decrement twice"
	case OpcodeDecrementOne:
		return "decrement"
	case OpcodeInput:
		return "input random"
	case OpcodeIncrementOne:
		return "increment"
	case OpcodeIncrementTwo:
		return "increment twice"
	case OpcodeMultiply:
		return "multiply by two"
	case OpcodeRepeatOne:
		return "repeat to one"
	case OpcodeRepeatTwo:
		return "repeat to two"
	case OpcodeRepeatThree:
		return "repeat to three"
	case OpcodeSkipOne:
		return "skip to one"
	case OpcodeSkipTwo:
		return "skip to two"
	case OpcodeSkipThree:
		return "skip to three"
	default:
		return ""
	}
}

func (o Opcode) Marker() uint {
	switch o {
	case OpcodeRepeatOne, OpcodeSkipOne:
		return 1
	case OpcodeRepeatTwo, OpcodeSkipTwo:
		return 2
	case OpcodeRepeatThree, OpcodeSkipThree:
		return 3
	default:
		return 0
	}
}

func (o Opcode) Repeats() bool {
	switch o {
	case OpcodeRepeatOne, OpcodeRepeatTwo, OpcodeRepeatThree:
		return true
	default:
		return false
	}
}

func (o Opcode) Skips() bool {
	switch o {
	case OpcodeSkipOne, OpcodeSkipTwo, OpcodeSkipThree:
		return true
	default:
		return false
	}
}
