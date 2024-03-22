package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/fabiensanglard/ccps/boards"
	"github.com/fabiensanglard/ccps/gfx"
	"github.com/fabiensanglard/ccps/sites"
	"github.com/spf13/cobra"
)

//go:embed postSrcs/m68k/crt0.s
var srcM68kCrt0 []byte

//go:embed postSrcs/m68k/main.c
var srcM68kMain []byte

//go:embed postSrcs/z80/crt0.s
var srcZ80Crt0 []byte

//go:embed postSrcs/z80/main.c
var srcZ80Main []byte

func post(cmd *cobra.Command, args []string) {
	postWithBytes(cmd, srcM68kCrt0, srcM68kMain, srcZ80Crt0, srcZ80Main)
}

func postWithBytes(cmd *cobra.Command,
	srcM68kCrt0 []byte,
	srcM68kMain []byte,
	srcZ80Crt0 []byte,
	srcZ80Main []byte) {

	board := boards.Get(targetBoard)
	srcM68kMain = replacePostValues(srcM68kMain, board)

	if err := os.MkdirAll(sites.M68kSrcsDir, 0777); err != nil {
		cmd.Printf("Unable to create m68K src dir '%s'\n", sites.M68kSrcsDir)
		os.Exit(1)
	}

	if err := os.MkdirAll(sites.Z80SrcsDir, 0777); err != nil {
		cmd.Printf("Unable to create z80 src dir '%s'\n", sites.Z80SrcsDir)
		os.Exit(1)
	}

	m68kMain := sites.M68kSrcsDir + "main.c"
	if verbose {
		cmd.Println("Creating", m68kMain)
	}
	os.WriteFile(m68kMain, srcM68kMain, 0677)

	m68kCrt0 := sites.M68kSrcsDir + "crt0.s"
	if verbose {
		cmd.Println("Creating", m68kCrt0)
	}
	os.WriteFile(m68kCrt0, srcM68kCrt0, 0677)

	z80Main := sites.Z80SrcsDir + "main.c"
	if verbose {
		cmd.Println("Creating", z80Main)
	}
	os.WriteFile(z80Main, srcZ80Main, 0677)

	z80Crt0 := sites.Z80SrcsDir + "crt0.s"
	if verbose {
		cmd.Println("Creating", z80Crt0)
	}
	os.WriteFile(z80Crt0, srcZ80Crt0, 0677)
}

func replacePostValues(kMain []byte, board *boards.Board) []byte {
	src := string(kMain)
	src = strings.Replace(src, "<TILE>", fmt.Sprintf("%d", board.Post.PostTile), 1)
	src = strings.Replace(src, "<TILE_HEIGHT>", fmt.Sprintf("%d", board.Post.PostTileHeight), 1)
	src = strings.Replace(src, "<TILE_WIDTH>", fmt.Sprintf("%d", board.Post.PostTileWidth), 1)
	src = strings.Replace(src, "<PALETTE>", fmt.Sprintf("{%s}", gfx.PaletteToString(board.Post.PostPalette)), 1)
	return []byte(src)
}
