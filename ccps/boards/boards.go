package boards

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Rom2Eprom func([]byte, string)

type CpsBRegisters struct {
	palette int
}

type ROM struct {
	Filename  string
	WordSize  int // If the ROM accessed in 1 byte, 2 byte or more?
	Size      int // Chip size
	Offset    int // Offset from where to start reading/writing the ROM
	DstOffset int // Offset in memory RAM
	Skip      int // How much to skip after each read/write when converting from ROM to memory space
}

type ROMSet struct {
	Size    int64
	Epromer Rom2Eprom
	Roms    []ROM
}

type GFXArea struct {
	Start int
	Size  int
	Dim   int
}

type Board struct {
	CpsBReg  CpsBRegisters
	GFXAreas [4]GFXArea
	GFX      ROMSet
	Z80      ROMSet
	M68k     ROMSet
	Oki      ROMSet
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
	sf2.Z80.Size = 65_536
	sf2.Z80.Epromer = sf2Z80EPromer

	sf2.M68k.Size = 1_048_576
	sf2.M68k.Epromer = sf2M68kEPromer

	sf2.GFX.Size = 12 * 0x80000 // 0x600000 = 6 MiB
	sf2.GFX.Epromer = sf2GFXEpromer
	// Fiding size can be done via PAL bank mapper (size * 64)
	sf2.GFXAreas = [4]GFXArea{
		{0, 0x480000, 16},
		{0x480000, 0x80000, 32},
		{0x500000, 0x40000, 8},
		{0x540000, 0x80000, 16}}

	sf2.GFX.Roms = []ROM{
		{Filename: "sf2-5m.4a", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x0000000, Skip: 8},
		{Filename: "sf2-7m.6a", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x0000002, Skip: 8},
		{Filename: "sf2-1m.3a", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x0000004, Skip: 8},
		{Filename: "sf2-3m.5a", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x0000006, Skip: 8},

		{Filename: "sf2-6m.4c", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x200000, Skip: 8},
		{Filename: "sf2-8m.6c", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x200002, Skip: 8},
		{Filename: "sf2-2m.3c", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x200004, Skip: 8},
		{Filename: "sf2-4m.5c", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x200006, Skip: 8},

		{Filename: "sf2-13m.4d", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x400000, Skip: 8},
		{Filename: "sf2-15m.6d", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x400002, Skip: 8},
		{Filename: "sf2-9m.3d", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x400004, Skip: 8},
		{Filename: "sf2-11m.5d", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x400006, Skip: 8},
	}

	sf2.Oki.Size = 0x40000
	sf2.Oki.Epromer = okiEpromer
	sf2.Oki.Roms = []ROM{
		{Filename: "sf2_18.11c", Size: 0x20000},
		{Filename: "sf2_19.12c", Size: 0x20000},
	}

	return &sf2
}

func okiEpromer(rom []byte, outDir string) {
	if rom == nil {
		return
	}
	const romSize = 0x20000
	rom1 := rom[0:romSize]
	rom2 := rom[romSize:]

	err := ioutil.WriteFile(outDir+"sf2_18.11c", rom1, 0644)
	if err != nil {
		fmt.Println("Unable to write Oki EPROM 'sf2_18.11c'")
		os.Exit(1)
	}

	err = ioutil.WriteFile(outDir+"sf2_19.12c", rom2, 0644)
	if err != nil {
		fmt.Println("Unable to write Oki EPROM 'sf2_19.12c'")
		os.Exit(1)
	}

}

func sf2GFXEpromer(gfxrom []byte, outDir string) {
	if gfxrom == nil {
		return
	}
	// Split it
	const BANK_SIZE = 0x200000
	bank0 := gfxrom[0*BANK_SIZE : 0*BANK_SIZE+BANK_SIZE]
	bank1 := gfxrom[1*BANK_SIZE : 1*BANK_SIZE+BANK_SIZE]
	bank2 := gfxrom[2*BANK_SIZE : 2*BANK_SIZE+BANK_SIZE]

	const ROM_SIZE = 0x80000

	writeToFile(bank0[0:0x80000], 2, 8, ROM_SIZE, outDir+"sf2-5m.4a")
	writeToFile(bank0[2:0x80000], 2, 8, ROM_SIZE, outDir+"sf2-7m.6a")
	writeToFile(bank0[4:0x80000], 2, 8, ROM_SIZE, outDir+"sf2-1m.3a")
	writeToFile(bank0[6:0x80000], 2, 8, ROM_SIZE, outDir+"sf2-3m.5a")

	writeToFile(bank1[0:0x80000], 2, 8, ROM_SIZE, outDir+"sf2-6m.4c")
	writeToFile(bank1[2:0x80000], 2, 8, ROM_SIZE, outDir+"sf2-8m.6c")
	writeToFile(bank1[4:0x80000], 2, 8, ROM_SIZE, outDir+"sf2-2m.3c")
	writeToFile(bank1[6:0x80000], 2, 8, ROM_SIZE, outDir+"sf2-4m.5c")

	writeToFile(bank2[0:0x80000], 2, 8, ROM_SIZE, outDir+"sf2-13m.4d")
	writeToFile(bank2[2:0x80000], 2, 8, ROM_SIZE, outDir+"sf2-15m.6d")
	writeToFile(bank2[4:0x80000], 2, 8, ROM_SIZE, outDir+"sf2-9m.3d")
	writeToFile(bank2[6:0x80000], 2, 8, ROM_SIZE, outDir+"sf2-11m.5d")
}

func sf2Z80EPromer(rom []byte, outputDir string) {
	path := outputDir + "sf2_9.12a"
	err := ioutil.WriteFile(path, rom, 0644)
	if err != nil {
		fmt.Println("Unable to write EPROM '", path, "'")
		os.Exit(1)
	}
}

func sf2M68kEPromer(r []byte, outputDir string) {
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
