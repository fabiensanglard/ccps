package main

import (
	"ccps/boards"
	"ccps/gfx"
	"ccps/m68k"
	"ccps/mus"
	"ccps/oki"
	"ccps/sites"
	"ccps/z80"
	"flag"
	"fmt"
)

func build(args []string) {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	verbose := fs.Bool("v", false, "Verbose mode")
	target := fs.String("b", "", "Target board")
	if err := fs.Parse(args); err != nil {
		println(fmt.Sprintln("Cmd parsing error '%s'", err))
	}

	if *target == "" {
		panic("Error: No board specified (-b)")
	}

	board := boards.Get(*target)

	// Create output folders where ROM and GFX/SFX/MFX source code will be generated
	sites.EnsureOutDir()
	sites.EnsureCodeGenDirs()

	// OKI generates oki.rom and oki.h (where the sound IDs are stored).
	// oki.h must be imported not in Z80 but in m68k because this is where
	// sound sample playback is decided.
	okyRom, okiIDsHeader := oki.Build(*verbose, board)
	board.Oki.Epromer(okyRom, sites.OutDir)
	okiIDsHeader.WriteTo(sites.Z80GenDir + "okiIds.h")

	// MUS generates no ROM, only mus.c. This source file is to be added to
	// the list of files to be compiled when creating z80.rom.
	z80Code := mus.Build(*verbose, board)
	z80Code.WriteTo(sites.Z80GenDir + "MFXz80.c")

	//
	z80Rom := z80.Build(*verbose, board)
	board.Z80.Epromer(z80Rom, sites.OutDir)

	// The GFX builder returns a ROM containing the GFX assets
	// but also a source code file containing tile IDs for the
	// shapes and sprites along with palettes values.
	gfxromPath, m68kDec, m68kDef := gfx.Build(*verbose, board)
	board.GFX.Epromer(gfxromPath, sites.OutDir)
	m68kDec.WriteTo(sites.M68kGenDir + "gfx.c")
	m68kDef.WriteTo(sites.M68kGenDir + "gfx.h")

	// Generate cpsa and cpsb header
	cpsAHeader := gfx.GenCpsAHeader(*verbose, board)
	cpsAHeader.WriteTo(sites.M68kGenDir + "cpsa.h")
	cpsBHeader := gfx.GenCpsBHeader(*verbose, board)
	cpsBHeader.WriteTo(sites.M68kGenDir + "cpsb.h")

	// Needs oki.h, mus.c, gfx.c
	m68kRom := m68k.Build(*verbose, board)
	board.M68k.Epromer(m68kRom, sites.OutDir)
}
