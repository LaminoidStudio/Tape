package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
)

const version = "v1.2"

func main() {
	// Handle command-line arguments
	var (
		input, output string
		memory        = 5
		seed          int64
		compile, decompile,
		original, run, step bool
	)
	flag.StringVar(&input, "input", input, "the input file name instead of stdin")
	flag.StringVar(&output, "output", output, "the output file name instead of stdout")
	flag.IntVar(&memory, "memory", memory, "the size of the tape in bytes")
	flag.Int64Var(&seed, "seed", seed, "predictable random number seed over randomness")
	flag.BoolVar(&compile, "compile", compile, "input source code over bytecode")
	flag.BoolVar(&decompile, "decompile", decompile, "output source code over bytecode")
	flag.BoolVar(&original, "original", original, "whether to use original brainfuck syntax")
	flag.BoolVar(&run, "run", run, "whether to run the program")
	flag.BoolVar(&step, "step", step, "whether to log after every step")
	flag.Usage = func() {
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Laminoid Tape Compiler & VM", version)
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), `(c) Laminoid Studio (Muessig & Muessig GbR), 2024

Description:
Laminoid Tape compiler produces byte-code programs from brainfuck-like programs.
It can also run and single-step, as well as decompile and explain them.
Both original brainfuck, as well as an enhanced version can be used.
The new version does not support unknown punctuation or symbol characters.
Only three levels of nesting can be used in original brainfuck mode.

Instructions:
. stop the program
< move left on the tape
> move right on the tape
/ divide by 2 (only new)
- decrement cell value
, input random value into cell
+ increment cell value
* multiply by 2 (only enhanced)
[ skip until 1 (enhanced) or go to to matching repeat (original)
( skip until 2 (only enhanced)
{ skip until 3 (only enhanced)
] repeat until 1 (enhanced) or go to to matching skip (original)
) repeat until 2 (only enhanced)
} repeat until 3 (only enhanced)

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
		p, err = Parse(reader, memory, original)
	} else {
		p, err = Read(reader, memory)
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
			_, _ = fmt.Fprintf(logger, "%x%c %s: %d%v\n", last/2, func() rune {
				if last%2 == 0 {
					return 'l'
				}
				return 'h'
			}(), p.Opcodes[last].Description(), p.Tape.Pointer, p.Tape.Cells)
		}
	}

	// If not logging every step, output the state after the last one
	if !step {
		_, _ = fmt.Fprintf(logger, "%x%c: %d%v\n", last/2, func() rune {
			if last%2 == 0 {
				return 'l'
			}
			return 'h'
		}(), p.Tape.Pointer, p.Tape.Cells)
	}
}