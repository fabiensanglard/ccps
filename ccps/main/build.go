package main

import (
	"ccps/boards"
	"ccps/gfx"
	"ccps/m68k"
	"ccps/mus"
	"ccps/oki"
	"ccps/z80"
	"flag"
	"fmt"
	"os"
)

func build([]string) {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	verbose := fs.Bool("v", false, "Verbose mode")
	target := fs.String("board", "sf2", "Target board")
	dryRun := fs.Bool("d", false, "Dry run (no rom generated)")
	if err := fs.Parse(os.Args[2:]); err != nil {
		println(fmt.Sprintln("Cmd parsing error '%s'", err))
	}

	board := boards.Get(*target)
	oki.Build(*verbose, *dryRun, board)
	mus.Build(*verbose, *dryRun, board)
	gfx.Build(*verbose, *dryRun, board)
	z80.Build(*verbose, *dryRun, board)
	m68k.Build(*verbose, *dryRun, board)
}
