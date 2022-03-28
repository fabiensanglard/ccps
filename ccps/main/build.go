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

func build(args []string) {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	verbose := fs.Bool("v", false, "Verbose mode")
	target := fs.String("b", "", "Target board")
	if err := fs.Parse(args); err != nil {
		println(fmt.Sprintln("Cmd parsing error '%s'", err))
	}

	if *target == "" {
		println("Error: No board specified (-b)")
		os.Exit(1)
	}

	board := boards.Get(*target)

	// Create output folder
	outputDir := "out/"
	err := os.RemoveAll(outputDir)
	err = os.MkdirAll(outputDir, 0777)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to create dir '%s'", outputDir))
		os.Exit(1)
	}

	// OKI generates oki.rom and oki.h
	okyRom := oki.Build(*verbose, board)
	board.Oki.Epromer(okyRom, "out/")

	// MUS generates mus.c
	mus.Build(*verbose, board)
	z80Rom := z80.Build(*verbose, board)
	board.Z80.Epromer(z80Rom, "out/")

	gfxromPath := gfx.Build(*verbose, board)
	board.GFX.Epromer(gfxromPath, "out/")

	// Needs oki.h, mus.c, gfx.c
	m68kRom := m68k.Build(*verbose, board)
	board.M68k.Epromer(m68kRom, "out/")
}
