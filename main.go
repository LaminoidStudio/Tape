package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
)

const version = "v1.5"

func main() {
	// Handle command-line arguments
	var (
		input, output string
		memory        = 8
		seed          int64
		compile, decompile,
		original, run, signed, step bool
	)
	flag.StringVar(&input, "input", input, "the input file name instead of stdin")
	flag.StringVar(&output, "output", output, "the output file name instead of stdout")
	flag.IntVar(&memory, "memory", memory, "the length of the tape in bytes")
	flag.Int64Var(&seed, "seed", seed, "predictable random number seed over randomness")
	flag.BoolVar(&compile, "compile", compile, "input source code instead of bytecode")
	flag.BoolVar(&decompile, "decompile", decompile, "output commented source code over bytecode")
	flag.BoolVar(&original, "original", original, "whether to use original brainfuck syntax")
	flag.BoolVar(&run, "run", run, "whether to run the program")
	flag.BoolVar(&signed, "signed", signed, "whether the tape cells are signed integers")
	flag.BoolVar(&step, "step", step, "whether to log every program step")
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
 When running while not decompiling, the output is disabled.

Instructions:
. stop the program
< move left on the tape
> move right on the tape
/ divide by 2 (only enhanced)
- decrement cell value
, input random value into cell
+ increment cell value
* multiply by 2 (only enhanced)
[ if cell is <= 0, skip until 1 (enhanced) or to matching repeat (original)
( if cell is <= 0, skip until 2 (only enhanced)
{ if cell is <= 0, skip until 3 (only enhanced)
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

	// Handle input
	var err error
	var reader = os.Stdin
	if input != "" {
		reader, err = os.Open(input)
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = reader.Close()
		}()
	}

	// Handle output
	var writer = io.Writer(os.Stdout)
	var logger = io.Writer(os.Stdout)
	if output != "" {
		writer, err = os.Create(output)
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = writer.(*os.File).Close()
		}()
	} else if run && !decompile {
		writer = io.Discard
	} else {
		logger = os.Stderr
	}

	// Load or compile the program
	var p *Program
	if compile {
		p, err = Parse(reader, memory, original, signed)
	} else {
		p, err = Read(reader, memory, signed)
	}
	if err != nil {
		panic(err)
	}

	// Output the compiled or decompiled program
	if decompile {
		err = p.Explain(writer)
	} else {
		err = p.Write(writer)
	}
	if err != nil {
		panic(err)
	}

	// Check if the program should run
	if !run {
		return
	}

	// Initialize the seed
	if seed != 0 {
		p.Random = rand.New(rand.NewSource(seed))
	}

	// And actually run the program
	var last int
	for p.Running() {
		last = p.Pointer
		p.Run()

		if step {
			_, _ = fmt.Fprintf(logger, "%x%c %s: %s\n", last/2, func() rune {
				if last%2 == 0 {
					return 'l'
				}
				return 'h'
			}(), p.Opcodes[last].Description(), p.Tape.String(signed))
		}
	}

	// If not logging every step, output the state after the last one
	if !step {
		_, _ = fmt.Fprintf(logger, "%x%c: %s\n", last/2, func() rune {
			if last%2 == 0 {
				return 'l'
			}
			return 'h'
		}(), p.Tape.String(signed))
	}
}
