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
	GFX     ROMSet
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

	sf2.GFX.Size = 6291456
	sf2.GFX.Epromer = sf2GFXEpromer
	return &sf2
}

func cp(src string, dst string) {
	//cmd := fmt.Sprintf("cp %s %s", src, dst)
	//if verbose
	//println(cmd)

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

func sf2GFXEpromer(romPath string, outDir string) {
	gfxrom, err := ioutil.ReadFile(romPath)

	if err != nil {
		fmt.Println(fmt.Sprintf("Cannot open GFX ROM '%s'", romPath))
		os.Exit(1)
	}

	// Split it
	const BANK_SIZE = 0x200000
	bank0 := gfxrom[0*BANK_SIZE : 0*BANK_SIZE+BANK_SIZE]
	bank1 := gfxrom[1*BANK_SIZE : 1*BANK_SIZE+BANK_SIZE]
	bank2 := gfxrom[2*BANK_SIZE : 2*BANK_SIZE+BANK_SIZE]

	const ROM_SIZE = 0x80000
	writeToFile(bank0[0:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_06.bin")
	writeToFile(bank0[2:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_08.bin")
	writeToFile(bank0[4:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_05.bin")
	writeToFile(bank0[6:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_07.bin")

	writeToFile(bank1[0:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_15.bin")
	writeToFile(bank1[2:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_17.bin")
	writeToFile(bank1[4:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_14.bin")
	writeToFile(bank1[6:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_16.bin")

	writeToFile(bank2[0:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_25.bin")
	writeToFile(bank2[2:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_27.bin")
	writeToFile(bank2[4:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_24.bin")
	writeToFile(bank2[6:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_26.bin")
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
	writeToFile(r[0*ROM_SIZE:], 1, 2, ROM_SIZE, outputDir+"sf2e_30g.11e")
	writeToFile(r[0*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"sf2e_37g.11f")
	writeToFile(r[1*ROM_SIZE:], 1, 2, ROM_SIZE, outputDir+"sf2e_31g.12e")
	writeToFile(r[1*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"sf2e_38g.12f")
	writeToFile(r[2*ROM_SIZE:], 1, 2, ROM_SIZE, outputDir+"sf2e_28g.9e")
	writeToFile(r[2*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"sf2e_35g.9f")
	writeToFile(r[3*ROM_SIZE:], 1, 2, ROM_SIZE, outputDir+"sf2_29b.10e")
	writeToFile(r[3*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"sf2_36b.10f")
}

// Generate a rom of size [num]. To do so, read from src [rom], batches of [size]
// bytes and skip [skip] bytes on each batch.
func writeToFile(rom []byte, size int, skip int, num int, filename string) {

	if num%size != 0 {
		println("Bad writeToFile. Size", num, "is not evenly divisible by ", size, ".")
		os.Exit(1)
	}

	var eprom = make([]byte, num)
	epromCursor := 0
	for romCursor := 0; romCursor < num/size; {
		copy(eprom[romCursor:romCursor+size], rom[epromCursor:epromCursor+size])
		romCursor += size
		epromCursor += skip
	}
	err := ioutil.WriteFile(filename, eprom, 0644)
	if err != nil {
		fmt.Println("Unable to write EPROM '", filename, "'")
		os.Exit(1)
	}
}
