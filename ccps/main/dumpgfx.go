package main

import (
	"ccps/boards"
	"flag"
	"fmt"
)

func dumpGFX(args []string) {
	fs := flag.NewFlagSet("gfx", flag.ContinueOnError)
	boardName := fs.String("board", "sf2", "Target board")

	if err := fs.Parse(args); err != nil {
		println(fmt.Sprintf("Cmd parsing error '%s'", err))
	}

	board := boards.Get(*boardName)

	println(board.GFXSizes[0])

	romSize := 0
	for _, size := range board.GFXSizes {
		romSize += size
	}

	// Desinterleave
	rom := desinterleave(boardName, romSize)

	// Dump ROM
	dims := []int{16, 8, 16, 32}
	cursor := 0
	for i := 0; i < 4; i++ {
		dim := dims[i]
		dumpSheets(dim, rom[cursor:cursor+board.GFXSizes[i]])
		cursor += board.GFXSizes[i]
	}
}

func dumpSheets(dim int, rom []byte) {
}

func desinterleave(name *string, size int) []byte {
	return nil
}
