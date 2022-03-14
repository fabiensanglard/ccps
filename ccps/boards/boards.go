package boards

import (
	"fmt"
	"os"
)

type CpsBRegisters struct {
	palette int
}

type ROMSet struct {
	Size      int64
	Filenames []string
}

type Board struct {
	CpsBReg CpsBRegisters
	Gfx     ROMSet
	Z80     ROMSet
	M68k    ROMSet
	Oki     ROMSet
}

func Get(name string) *Board {
	if name == "sf2" {
		return sf2Board()
	}
	println(fmt.Sprintln("Unknown board '%s'", name))
	os.Exit(1)
	return &Board{}
}

var boards []Board

func sf2Board() *Board {
	return &Board{}
}
