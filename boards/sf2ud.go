package boards

import (
	"fmt"
	"image/color"
	"os"
)

func sf2udBoard() *Board {
	sf2ud := Board{}

	sf2ud.Cpsb = 5

	sf2ud.Z80.Size = 65_536
	sf2ud.Z80.Epromer = sf2udZ80EPromer

	sf2ud.M68k.Size = 1_048_576
	sf2ud.M68k.Epromer = sf2udM68kEPromer

	sf2ud.GFX.Size = 12 * 0x80000 // 0x600000 = 6 MiB
	sf2ud.GFX.Epromer = sf2udGFXEpromer
	sf2ud.GFX.Roms = []ROM{
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

	sf2ud.Oki.Size = 0x40000
	sf2ud.Oki.Epromer = okisf2udEpromer
	sf2ud.Oki.Roms = []ROM{
		{Filename: "sf2_18.11c", Size: 0x20000},
		{Filename: "sf2_19.12c", Size: 0x20000},
	}

	// Finding size can be done via PAL bank mapper (size * 64)
	sf2ud.GFXAreas = [4]GFXArea{
		{0, 0x480000, 16},
		{0x480000, 0x80000, 32},
		{0x500000, 0x40000, 8},
		{0x540000, 0x80000, 16}}

	// RGBA
	sf2ud.Post.PostPalette = color.Palette{
		color.RGBA{0x11, 0x11, 0x11, 0xff},
		color.RGBA{0xFF, 0xDD, 0x99, 0xff},
		color.RGBA{0xFF, 0xBB, 0x99, 0xff},
		color.RGBA{0xEE, 0x99, 0x77, 0xff},
		color.RGBA{0xCC, 0x88, 0x66, 0xff},
		color.RGBA{0x99, 0x66, 0x55, 0xff},
		color.RGBA{0x66, 0x44, 0x33, 0xff},
		color.RGBA{0xBB, 0x00, 0x00, 0xff},
		color.RGBA{0xFF, 0xFF, 0xFF, 0xff},
		color.RGBA{0xEE, 0xEE, 0xCC, 0xff},
		color.RGBA{0xDD, 0xCC, 0xAA, 0xff},
		color.RGBA{0xBB, 0xAA, 0x88, 0xff},
		color.RGBA{0xAA, 0x88, 0x77, 0xff},
		color.RGBA{0x77, 0x66, 0x55, 0xff},
		color.RGBA{0xFF, 0x00, 0x00, 0xff},
		color.RGBA{0x00, 0x00, 0x00, 0x00},
	}
	sf2ud.Post.PostTile = 4
	sf2ud.Post.PostTileWidth = 3
	sf2ud.Post.PostTileHeight = 5

	return &sf2ud
}

func sf2udGFXEpromer(gfxrom []byte, outDir string) {
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

func sf2udZ80EPromer(rom []byte, outputDir string) {
	path := outputDir + "sf2_9.12a"
	err := os.WriteFile(path, rom, 0644)
	if err != nil {
		panic(fmt.Sprintf("Unable to write EPROM '%s'", path))
	}
}

func sf2udM68kEPromer(r []byte, outputDir string) {
	const ROM_SIZE = 131072
	writeToFile(r[0*ROM_SIZE:], 1, 2, ROM_SIZE,   outputDir+"sf2u_30d.11e")
	writeToFile(r[0*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"sf2u_37d.11f")
	writeToFile(r[1*ROM_SIZE:], 1, 2, ROM_SIZE,   outputDir+"sf2u_31d.12e")
	writeToFile(r[1*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"sf2u_38d.12f")
	writeToFile(r[2*ROM_SIZE:], 1, 2, ROM_SIZE,   outputDir+"sf2u_28d.9e")
	writeToFile(r[2*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"sf2u_35d.9f")
	writeToFile(r[3*ROM_SIZE:], 1, 2, ROM_SIZE,   outputDir+"sf2_29b.10e")
	writeToFile(r[3*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"sf2_36b.10f")
}

func okisf2udEpromer(rom []byte, outDir string) {
	if rom == nil {
		return
	}
	const romSize = 0x20000
	rom1 := rom[0:romSize]
	rom2 := rom[romSize:]

	err := os.WriteFile(outDir+"sf2_18.11c", rom1, 0644)
	if err != nil {
		panic("Unable to write Oki EPROM 'sf2_18.11c'")
	}

	err = os.WriteFile(outDir+"sf2_19.12c", rom2, 0644)
	if err != nil {
		panic("Unable to write Oki EPROM 'sf2_19.12c'")
	}

}
