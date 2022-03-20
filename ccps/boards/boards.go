package boards

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Rom2Eprom func(string, string)

type CpsBRegisters struct {
	palette int
}

type ROMSet struct {
	Size    int64
	Epromer Rom2Eprom
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
	sf2 := Board{}
	sf2.Z80.Size = 65536
	sf2.Z80.Epromer = sf2Z80EPromer

	sf2.M68k.Size = 1048576
	sf2.M68k.Epromer = sf2M68kEPromer
	return &sf2
}

func cp(src string, dst string) {
	cmd := fmt.Sprintf("cp %s %s", src, dst)
	println(cmd)

	input, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(fmt.Sprintf("cp error: Cannot open src: '%s'", src))
		os.Exit(1)
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		fmt.Println(fmt.Sprintf("cp error: Cannot dst: '%s'", dst))
		os.Exit(1)
	}
}

func sf2Z80EPromer(rom string, outputDir string) {
	cp(rom, outputDir+"sf2_9.12a")
}

func sf2M68kEPromer(rom string, outputDir string) {
	r, err := ioutil.ReadFile(rom)

	if err != nil {
		fmt.Println(fmt.Sprintf("Cannot open Z80 ROM '%s'", rom))
		os.Exit(1)
	}

	const ROM_SIZE = 131072
	writeToFile(r[0*ROM_SIZE:], 2, ROM_SIZE, outputDir+"sf2e_30g.11e")
	writeToFile(r[0*ROM_SIZE+1:], 2, ROM_SIZE, outputDir+"sf2e_37g.11f")
	writeToFile(r[1*ROM_SIZE:], 2, ROM_SIZE, outputDir+"sf2e_31g.12e")
	writeToFile(r[1*ROM_SIZE+1:], 2, ROM_SIZE, outputDir+"sf2e_38g.12f")
	writeToFile(r[2*ROM_SIZE:], 2, ROM_SIZE, outputDir+"sf2e_28g.9e")
	writeToFile(r[2*ROM_SIZE+1:], 2, ROM_SIZE, outputDir+"sf2e_35g.9f")
	writeToFile(r[3*ROM_SIZE:], 2, ROM_SIZE, outputDir+"sf2_29b.10e")
	writeToFile(r[3*ROM_SIZE+1:], 2, ROM_SIZE, outputDir+"sf2_36b.10f")

}

func writeToFile(rom []byte, stride int, size int, filename string) {
	var eprom = make([]byte, size)

	for i := 0; i < size; i++ {
		eprom[i] = rom[i*stride]
	}
	err := ioutil.WriteFile(filename, eprom, 0644)
	if err != nil {
		fmt.Println("Unable to write EPROM '", filename, "'")
		os.Exit(1)
	}
}
