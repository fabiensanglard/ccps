package main

import (
	"ccps/sites"
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
	postWithBytes(args, srcM68kCrt0, srcM68kMain, srcZ80Crt0, srcZ80Main)
}

func postWithBytes(args []string,
	srcM68kCrt0 []byte,
	srcM68kMain []byte,
	srcZ80Crt0 []byte,
	srcZ80Main []byte) {
	fs := flag.NewFlagSet("postFlags", flag.ContinueOnError)
	verbose := fs.Bool("v", false, "Verbose mode")

	if err := fs.Parse(args); err != nil {
		println(fmt.Sprintf("Cmd parsing error '%s'", err))
		os.Exit(1)
	}

	err := os.MkdirAll(sites.M68kSrcsDir, 0777)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to create m68K src dir '%s'", sites.M68kSrcsDir))
		os.Exit(1)
	}

	err = os.MkdirAll(sites.Z80SrcsDir, 0777)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to create z80 src dir '%s'", sites.Z80SrcsDir))
		os.Exit(1)
	}

	m68kMain := sites.M68kSrcsDir + "main.c"
	if *verbose {
		println("Creating", m68kMain)
	}
	os.WriteFile(m68kMain, srcM68kMain, 0677)

	m68kCrt0 := sites.M68kSrcsDir + "crt0.s"
	if *verbose {
		println("Creating", m68kCrt0)
	}
	os.WriteFile(m68kCrt0, srcM68kCrt0, 0677)

	z80Main := sites.Z80SrcsDir + "main.c"
	if *verbose {
		println("Creating", z80Main)
	}
	os.WriteFile(z80Main, srcZ80Main, 0677)

	z80Crt0 := sites.Z80SrcsDir + "crt0.s"
	if *verbose {
		println("Creating", z80Crt0)
	}
	os.WriteFile(z80Crt0, srcZ80Crt0, 0677)
}
