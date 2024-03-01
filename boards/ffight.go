package boards

import (
	"fmt"
	"image/color"
	"os"
)

func ffightBoard() *Board {
	ffight := Board{}

	ffight.Cpsb = 4

	ffight.Z80.Size = 65_536
	ffight.Z80.Epromer = ffightZ80EPromer

	ffight.M68k.Size = 1_048_576
	ffight.M68k.Epromer = ffightM68kEPromer

	ffight.GFX.Size = 4 * 0x80000 // 0x200000 = 2 MiB
	ffight.GFX.Epromer = ffightGFXEpromer
	ffight.GFX.Roms = []ROM{
		{Filename: "ff-5m.7a", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x0000000, Skip: 8},
		{Filename: "ff-7m.9a", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x0000002, Skip: 8},
		{Filename: "ff-1m.3a", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x0000004, Skip: 8},
		{Filename: "ff-3m.5a", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x0000006, Skip: 8},
	}

	ffight.Oki.Size = 0x40000
	ffight.Oki.Epromer = ffightOkiEpromer
	ffight.Oki.Roms = []ROM{
		{Filename: "ff_18.11c", Size: 0x20000},
		{Filename: "ff_19.12c", Size: 0x20000},
	}

	// static const struct gfx_range mapper_STF29_table[] =
	// {
		// verified from PAL dump:
		// bank 0 = pin 19 (ROMs 5,6,7,8)
		// bank 1 = pin 14 (ROMs 14,15,16,17)
		// bank 2 = pin 12 (ROMS 24,25,26,27)
	
		/* type 		   start	end 	 bank */
		// { GFXTYPE_SPRITES, 0x00000, 0x07fff, 0 },
	
		// { GFXTYPE_SPRITES, 0x08000, 0x0ffff, 1 },
	
		// { GFXTYPE_SPRITES, 0x10000, 0x11fff, 2 },
		// { GFXTYPE_SCROLL3, 0x02000, 0x03fff, 2 },
		// { GFXTYPE_SCROLL1, 0x04000, 0x04fff, 2 },
		// { GFXTYPE_SCROLL2, 0x05000, 0x07fff, 2 },
		// { 0 }
	// };

	// OBJ	= 0x8000 + 0x2000 + 0x2000 = 0x12000 * 0x40 = 0x480000 (1)
	// SCR1 = 0x1000 * 0x40 = 0x40000 (2)
	// SCR2 = 0x3000 * 0x40 = 0xC0000 (3)
	// SCR3 = 0x2000 * 0x40 = 0x80000 (1)

	// Finding size can be done via PAL bank mapper (size * 64)
	// GFXArea; Start/Size/Dim
	// OBJ, SCR1, SCR2, SCR3
	ffight.GFXAreas = [4]GFXArea{
		{0, 0x480000, 16},
		{0x480000, 0x80000, 32},
		{0x500000, 0x40000, 8},
		{0x540000, 0x80000, 16}}

	// RGBA
	ffight.Post.PostPalette = color.Palette{
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
	ffight.Post.PostTile = 4
	ffight.Post.PostTileWidth = 3
	ffight.Post.PostTileHeight = 5

	return &ffight
}

func ffightGFXEpromer(gfxrom []byte, outDir string) {
	if gfxrom == nil {
		return
	}
	// Split it
	const BANK_SIZE = 0x200000
	bank0 := gfxrom[0*BANK_SIZE : 0*BANK_SIZE+BANK_SIZE]

	const ROM_SIZE = 0x80000

	writeToFile(bank0[0:0x80000], 2, 8, ROM_SIZE, outDir+"ff-5m.7a")
	writeToFile(bank0[2:0x80000], 2, 8, ROM_SIZE, outDir+"ff-7m.9a")
	writeToFile(bank0[4:0x80000], 2, 8, ROM_SIZE, outDir+"ff-1m.3a")
	writeToFile(bank0[6:0x80000], 2, 8, ROM_SIZE, outDir+"ff-3m.5a")
}

func ffightZ80EPromer(rom []byte, outputDir string) {
	path := outputDir + "ff_09.12b"
	err := os.WriteFile(path, rom, 0644)
	if err != nil {
		panic(fmt.Sprintf("Unable to write EPROM '%s'", path))
	}
}

func ffightM68kEPromer(r []byte, outputDir string) {
	const ROM_SIZE = 131072
	writeToFile(r[0*ROM_SIZE:], 1, 2, ROM_SIZE,   outputDir+"ff_36.11f")
	writeToFile(r[0*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"ff_42.11h")
	writeToFile(r[1*ROM_SIZE:], 1, 2, ROM_SIZE,   outputDir+"ff_37.12f")
	writeToFile(r[1*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"ffe_43.12h")
    // WORD SWAPPED ROM ??
	writeToFile(r[2*ROM_SIZE:], 1, 2, 2*ROM_SIZE, outputDir+"ff-32m.8h")
}

func ffightOkiEpromer(rom []byte, outDir string) {
	if rom == nil {
		return
	}
	const romSize = 0x20000
	rom1 := rom[0:romSize]
	rom2 := rom[romSize:]

	err := os.WriteFile(outDir+"ff_18.11c", rom1, 0644)
	if err != nil {
		panic("Unable to write Oki EPROM 'ff_18.11c'")
	}

	err = os.WriteFile(outDir+"ff_19.12c", rom2, 0644)
	if err != nil {
		panic("Unable to write Oki EPROM 'ff_19.12c'")
	}

}
