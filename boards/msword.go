package boards

import (
	"fmt"
	"image/color"
	"os"
)

func mswordBoard() *Board {
	msword := Board{}

	msword.Cpsb = 13

	msword.Z80.Size = 65_536
	msword.Z80.Epromer = mswordZ80EPromer

	msword.M68k.Size = 1_048_576
	msword.M68k.Epromer = mswordM68kEPromer

	msword.GFX.Size = 4 * 0x80000 // 0x200000 = 2 MiB
	msword.GFX.Epromer = mswordGFXEpromer
	msword.GFX.Roms = []ROM{
		{Filename: "ms-5m.7a", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x0000000, Skip: 8},
		{Filename: "ms-7m.9a", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x0000002, Skip: 8},
		{Filename: "ms-1m.3a", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x0000004, Skip: 8},
		{Filename: "ms-3m.5a", WordSize: 2, Offset: 0, Size: 0x80000, DstOffset: 0x0000006, Skip: 8},
	}

	msword.Oki.Size = 0x40000
	msword.Oki.Epromer = mswordOkiEpromer
	msword.Oki.Roms = []ROM{
		{Filename: "ms_18.11c", Size: 0x20000},
		{Filename: "ms_19.12c", Size: 0x20000},
	}


    // MS24B and MS22B are equivalent, but since we could dump both PALs we are
    // documenting both.
    
    // #define mapper_MS24B    { 0x8000, 0, 0, 0 }, mapper_MS24B_table
    // static const struct gfx_range mapper_MS24B_table[] =
    // {
	    // verified from PAL dump:
	    // bank 0 = pin 16 (ROMs 1,3,5,7)
	    // pin 14 duplicates pin 16 allowing to populate the 8-bit ROM sockets
	    // instead of the 16-bit ones.
	    // pin 12 is enabled only for sprites:
	    // 0 0000-3fff
	    // pin 19 is never enabled
    
	    // type            start   end     bank
	    // { GFXTYPE_SPRITES, 0x0000, 0x3fff, 0 },
	    // { GFXTYPE_SCROLL1, 0x4000, 0x4fff, 0 },
	    // { GFXTYPE_SCROLL2, 0x5000, 0x6fff, 0 },
	    // { GFXTYPE_SCROLL3, 0x7000, 0x7fff, 0 },
	    //{ 0 }
    // };

    // OBJ	= 0x4000 * 0x40 = 0x100000 (1)
	// SCR1 = 0x1000 * 0x40 = 0x40000 (2)
	// SCR2 = 0x2000 * 0x40 = 0x80000 (3)
	// SCR3 = 0x1000 * 0x40 = 0x40000 (1)


	// Finding size can be done via PAL bank mapper (size * 64)
	msword.GFXAreas = [4]GFXArea{
		{0, 0x100000, 16},
		{0x100000, 0x80000, 32},
		{0x120000, 0x40000, 8},
		{0x160000, 0x80000, 16}}

	// RGBA
	msword.Post.PostPalette = color.Palette{
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
	msword.Post.PostTile = 4
	msword.Post.PostTileWidth = 3
	msword.Post.PostTileHeight = 5

	return &msword
}

func mswordGFXEpromer(gfxrom []byte, outDir string) {
	if gfxrom == nil {
		return
	}
	// Split it
	const BANK_SIZE = 0x200000
	bank0 := gfxrom[0*BANK_SIZE : 0*BANK_SIZE+BANK_SIZE]

	const ROM_SIZE = 0x80000

	writeToFile(bank0[0:0x80000], 2, 8, ROM_SIZE, outDir+"ms-5m.7a")
	writeToFile(bank0[2:0x80000], 2, 8, ROM_SIZE, outDir+"ms-7m.9a")
	writeToFile(bank0[4:0x80000], 2, 8, ROM_SIZE, outDir+"ms-1m.3a")
	writeToFile(bank0[6:0x80000], 2, 8, ROM_SIZE, outDir+"ms-3m.5a")
}

func mswordZ80EPromer(rom []byte, outputDir string) {
	path := outputDir + "ms_09.12b"
	err := os.WriteFile(path, rom, 0644)
	if err != nil {
		panic(fmt.Sprintf("Unable to write EPROM '%s'", path))
	}
}

func mswordM68kEPromer(r []byte, outputDir string) {
	const ROM_SIZE = 131072
	writeToFile(r[0*ROM_SIZE:], 1, 2, ROM_SIZE,   outputDir+"mse_30.11f")
	writeToFile(r[0*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"mse_35.11h")
	writeToFile(r[1*ROM_SIZE:], 1, 2, ROM_SIZE,   outputDir+"mse_31.12f")
	writeToFile(r[1*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"mse_36.12h")
    // WORD SWAPPED ROM
	// ROM_LOAD16_WORD_SWAP( "ms-32m.8h", 0x80000, 0x80000, CRC(2475ddfc) SHA1(cc34dfae8124aa781320be6870a1929495eee456) )
}

func mswordOkiEpromer(rom []byte, outDir string) {
	if rom == nil {
		return
	}
	const romSize = 0x20000
	rom1 := rom[0:romSize]
	rom2 := rom[romSize:]

	err := os.WriteFile(outDir+"ms_18.11c", rom1, 0644)
	if err != nil {
		panic("Unable to write Oki EPROM 'ms_18.11c'")
	}

	err = os.WriteFile(outDir+"ms_19.12c", rom2, 0644)
	if err != nil {
		panic("Unable to write Oki EPROM 'ms_19.12c'")
	}

}
