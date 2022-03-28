package main

import (
	"ccps/m68k"
	"ccps/z80"
	_ "embed"
	"flag"
	"fmt"
	"os"
)

//go:embed postSrcs/m68k/crt0.s
var srcM68kCrt0 []byte

//go:embed postSrcs/m68k/main.c
var srcM68kMain []byte

//go:embed postSrcs/z80/crt0.s
var srcZ80Crt0 []byte

//go:embed postSrcs/z80/main.c
var srcZ80Main []byte

func post(args []string) {
	fs := flag.NewFlagSet("gfx", flag.ContinueOnError)
	verbose := fs.Bool("v", false, "Verbose mode")

	if err := fs.Parse(args); err != nil {
		println(fmt.Sprintf("Cmd parsing error '%s'", err))
		os.Exit(1)
	}

	err := os.MkdirAll(m68k.SrcsPath, 0777)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to create m68K src dir '%s'", m68k.SrcsPath))
		os.Exit(1)
	}

	err = os.MkdirAll(z80.SrcsPath, 0777)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to create z80 src dir '%s'", m68k.SrcsPath))
		os.Exit(1)
	}

	m68kMain := m68k.SrcsPath + "main.c"
	if *verbose {
		println("Creating", m68kMain)
	}
	os.WriteFile(m68kMain, srcM68kMain, 0677)

	m68kCrt0 := m68k.SrcsPath + "crt0.s"
	if *verbose {
		println("Creating", m68kCrt0)
	}
	os.WriteFile(m68kCrt0, srcM68kCrt0, 0677)

	z80Main := z80.SrcsPath + "main.c"
	if *verbose {
		println("Creating", z80Main)
	}
	os.WriteFile(z80Main, srcZ80Main, 0677)

	z80Crt0 := z80.SrcsPath + "crt0.s"
	if *verbose {
		println("Creating", z80Crt0)
	}
	os.WriteFile(z80Crt0, srcZ80Crt0, 0677)
}
