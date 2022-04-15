package boards

import (
	"fmt"
	"image/color"
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

type PostInfo struct {
	PostPalette    color.Palette
	PostTile       int
	PostTileWidth  int
	PostTileHeight int
}

type Board struct {
	CpsBReg  CpsBRegisters
	GFXAreas [4]GFXArea
	GFX      ROMSet
	Z80      ROMSet
	M68k     ROMSet
	Oki      ROMSet
	Cpsb     int // Version of CPS-B
	Post     PostInfo
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
