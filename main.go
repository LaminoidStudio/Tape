package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const version = "v1.7w"

func main() {
	// Handle panics
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic: %+v\n", err)
		}
	}()

	// Handle command-line arguments
	var (
		code                  string
		memory                = 8
		seed                  int64
		timeout               time.Duration
		run, original, signed bool
	)
	flag.StringVar(&code, "code", code, "the source code to compile")
	flag.IntVar(&memory, "memory", memory, "the length of the tape in bytes")
	flag.DurationVar(&timeout, "timeout", timeout, "the longest duration after which to stop running")
	flag.Int64Var(&seed, "seed", seed, "predictable random number seed over randomness")
	flag.BoolVar(&run, "run", run, "whether to run the code after compilation instead of just decompiling it")
	flag.BoolVar(&original, "original", original, "whether to use original brainfuck syntax")
	flag.BoolVar(&signed, "signed", signed, "whether the tape cells are signed integers")
	flag.Usage = func() {
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Laminoid Tape Compiler & VM", version)
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), `(c) Laminoid Studio (Muessig & Muessig GbR), 2024

Description:
 Laminoid Tape compiles efficient bytecode programs from brainfuck-like source.
 It can also run and single-step, as well as decompile and explain them.
 Both original brainfuck, as well as an enhanced version can be used.
 The enhanced version doesn't support unknown punctuation or symbol characters.
 Only three levels of nesting can be used in original brainfuck mode.
 No strictly matched nesting is required in the enhanced mode.
 Regular numbers, letters and whitespace are ignored in both modes.
 All tape cells are by default unsigned 8-bit integers that wrap.
 If signed mode is disabled, the conditionals behave as usual.

Instructions:
. stop the program
< move left on the tape
> move right on the tape
/ divide by 2 (only enhanced)
- decrement cell value
, input random value into cell
+ increment cell value
* multiply by 2 (only enhanced)
[ if cell is == 0, skip until 1 (enhanced) or to matching repeat (original)
( if cell is == 0, skip until 2 (only enhanced)
{ if cell is == 0, skip until 3 (only enhanced)
] if cell is > 0, repeat until 1 (enhanced) or to matching skip (original)
) if cell is > 0, repeat until 2 (only enhanced)
} if cell is > 0, repeat until 3 (only enhanced)

Encoding:
 Two opcodes encoded per byte, where the low nibble precedes the high nibble.
 At the end of the program, padding can be added with the output opcode (zero).

Opcodes:
 Output (0), Left (1), Right (2), Divide (3),
 DecrementTwo (4), DecrementOne (5), Input (6),
 IncrementOne (7), IncrementTwo (8), Multiply (9),
 RepeatOne (10), RepeatTwo (11), RepeatThree (12),
 SkipOne (13), SkipTwo (14), SkipThree (15)

Usage:`)
		flag.PrintDefaults()
	}
	flag.Parse()

	// Verify memory size
	if memory < 1 {
		panic(errors.New("at least 1 byte of memory must be allocated"))
	}

	// Verify code input
	if len(code) < 1 {
		panic(errors.New("there must be code to compile"))
	}

	// Load or compile the program
	var p *Program
	var err error
	p, err = Parse(strings.NewReader(code), memory, original, signed)
	if err != nil {
		panic(err)
	}

	// Output the decompiled program if not running
	if !run {
		err = p.Explain(os.Stdout)
		if err != nil {
			panic(err)
		}
		return
	}

	// Initialize the seed
	if seed != 0 {
		p.Random = rand.New(rand.NewSource(seed))
	}

	// And actually run the program
	var last int
	var start = time.Now()
	for p.Running() {
		last = p.Pointer
		p.Run()

		_, _ = fmt.Fprintf(os.Stdout, "%x%c %s: %s\n", last/2, func() rune {
			if last%2 == 0 {
				return 'l'
			}
			return 'h'
		}(), p.Opcodes[last].Description(), p.Tape.String(signed))

		if timeout > 0 && time.Now().After(start.Add(timeout)) {
			panic(errors.New("timeout"))
		}
	}
}
