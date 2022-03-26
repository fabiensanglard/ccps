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

	// Create output folder
	outputDir := "out/"
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to create dir '%s'", outputDir))
		os.Exit(1)
	}

	// OKI generates oki.rom and oki.h
	okyRom := oki.Build(*verbose, *dryRun, board)
	board.Oki.Epromer(okyRom, "out/")

	// MUS generates mus.c
	mus.Build(*verbose, *dryRun, board)
	z80Rom := z80.Build(*verbose, *dryRun, board)
	board.Z80.Epromer(z80Rom, "out/")

	gfxromPath := gfx.Build(*verbose, *dryRun, board)
	board.GFX.Epromer(gfxromPath, "out/")

	// Needs oki.h, mus.c, gfx.c
	m68kRom := m68k.Build(*verbose, *dryRun, board)
	board.M68k.Epromer(m68kRom, "out/")
}
