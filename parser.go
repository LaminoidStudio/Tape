package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"
)

const (
	TokenOutput      = '.'
	TokenLeft        = '<'
	TokenRight       = '>'
	TokenDivide      = '/'
	TokenDecrement   = '-'
	TokenInput       = ','
	TokenIncrement   = '+'
	TokenMultiply    = '*'
	TokenRepeatOne   = ']'
	TokenRepeatTwo   = ')'
	TokenRepeatThree = '}'
	TokenSkipOne     = '['
	TokenSkipTwo     = '('
	TokenSkipThree   = '{'
)

func Read(reader io.Reader, memory int) (p *Program, err error) {
	// Read the entire program
	data, err := io.ReadAll(reader)
	if err != nil {
		return
	}

	// And turn it into opcodes
	var opcodes []Opcode
	for _, b := range data {
		opcodes = append(opcodes, Opcode(b&0xf), Opcode(b>>4))
	}

	// Finally create the program
	p = &Program{
		Opcodes: opcodes,
		Tape:    Tape{Cells: make([]uint8, memory)},
	}
	return
}

func Parse(reader io.Reader, memory int, original bool) (p *Program, err error) {
	var buffered = bufio.NewReader(reader)
	var line, column = 1, 1
	var depth int
	var callback func(next rune) bool
	var opcodes []Opcode

	for {
		// Read the next character
		var curr rune
		if curr, _, err = buffered.ReadRune(); err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		} else if curr == unicode.ReplacementChar {
			continue
		}

		// Advance the column for all printing characters
		if unicode.IsPrint(curr) {
			column++
		}

		// If there is a callback, run it and see if it handled the next character
		if callback != nil {
			handled := callback(curr)
			callback = nil
			if handled {
				continue
			}
		}

		// For unhandled characters, decide what to do with them
		switch curr {
		case '\r':
			// Register a callback to consume the next line-feed after a carriage return
			callback = func(next rune) bool {
				return next == '\n'
			}
			fallthrough

		case '\n':
			// Always advance the line for carriage returns or unhandled line feeds
			line++
			column = 1

		case TokenOutput: // .
			// works the same in both, emitted as-is
			opcodes = append(opcodes, OpcodeOutput)

		case TokenLeft: // <
			// works the same in both, emitted as-is
			opcodes = append(opcodes, OpcodeLeft)

		case TokenRight: // >
			// works the same in both, emitted as-is
			opcodes = append(opcodes, OpcodeRight)

		case TokenDivide:
			// ignored in original, emitted as-is else
			if !original {
				opcodes = append(opcodes, OpcodeDivide)
			}

		case TokenDecrement: // -
			// works the same in both and pairs of two are collated
			callback = func(next rune) bool {
				if next == TokenDecrement {
					opcodes = append(opcodes, OpcodeDecrementTwo)
					return true
				}

				opcodes = append(opcodes, OpcodeDecrementOne)
				return false
			}

		case TokenInput: // ,
			// works the same in both, emitted as-is
			opcodes = append(opcodes, OpcodeInput)

		case TokenIncrement: // +
			// works the same in both and pairs of two are collated
			callback = func(next rune) bool {
				if next == TokenIncrement {
					opcodes = append(opcodes, OpcodeIncrementTwo)
					return true
				}

				opcodes = append(opcodes, OpcodeIncrementOne)
				return false
			}

		case TokenMultiply:
			// ignored in original, emitted as-is else
			if !original {
				opcodes = append(opcodes, OpcodeMultiply)
			}

		case TokenRepeatOne: // ]
			// search existing half or crash if none in original, emitted as-is else
			if original {
				switch depth {
				case 1:
					opcodes = append(opcodes, OpcodeRepeatOne)
				case 2:
					opcodes = append(opcodes, OpcodeRepeatTwo)
				case 3:
					opcodes = append(opcodes, OpcodeRepeatThree)
				default:
					err = fmt.Errorf("unmatched parenthesis at line %d, column %d", line, column)
					return
				}
				depth--
			} else {
				opcodes = append(opcodes, OpcodeRepeatOne)
			}

		case TokenRepeatTwo: // )
			// ignored in original, emitted as-is else
			if !original {
				opcodes = append(opcodes, OpcodeRepeatTwo)
			}

		case TokenRepeatThree: // }
			// ignored in original, emitted as-is else
			if !original {
				opcodes = append(opcodes, OpcodeRepeatThree)
			}

		case TokenSkipOne: // [
			// search free pair or crash if none in original, emitted as-is else
			if original {
				switch depth {
				case 0:
					opcodes = append(opcodes, OpcodeSkipOne)
				case 1:
					opcodes = append(opcodes, OpcodeSkipTwo)
				case 2:
					opcodes = append(opcodes, OpcodeSkipThree)
				default:
					err = fmt.Errorf("maximum nesting depth of 3 exceeded at line %d, column %d", line, column)
					return
				}
				depth++
			} else {
				opcodes = append(opcodes, OpcodeSkipOne)
			}

		case TokenSkipTwo: // (
			// ignored in original, emitted as-is else
			if !original {
				opcodes = append(opcodes, OpcodeSkipTwo)
			}

		case TokenSkipThree: // {
			// ignored in original, emitted as-is else
			if !original {
				opcodes = append(opcodes, OpcodeSkipThree)
			}
		default:
			if !original && (unicode.IsPunct(curr) || unicode.IsSymbol(curr)) {
				err = fmt.Errorf("unexpected character %c at line %d, column %d", curr, line, column)
				return
			}
		}
	}

	// At the end, handle any remaining callback
	if callback != nil {
		callback(0)
	}

	// And make sure all parentheses are closed
	if original && depth != 0 {
		err = errors.New("open parentheses at the end of input")
		return
	} else if !original {
		var foundBlocks = map[uint]bool{}
		for _, o := range opcodes {
			marker := o.Marker()
			if o.Skips() {
				foundBlocks[marker] = true
			} else if o.Repeats() {
				unmatched, ok := foundBlocks[marker]
				if !unmatched && !ok {
					err = fmt.Errorf("missing skipping parentheses for %d", marker)
					return
				}

				foundBlocks[marker] = false
			}
		}
		for marker, unmatched := range foundBlocks {
			if unmatched {
				err = fmt.Errorf("missing repeating parentheses for %d", marker)
				return
			}
		}
	}

	// And create the program
	p = &Program{
		Opcodes: opcodes,
		Tape:    Tape{Cells: make([]uint8, memory)},
	}
	return
}
