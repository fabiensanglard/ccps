package main

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"

	"github.com/fabiensanglard/ccps/boards"
	"github.com/fabiensanglard/ccps/sites"
	"github.com/spf13/cobra"
)

type Palette struct {
	colors [16]color.RGBA
}

var greyPalette = Palette{[16]color.RGBA{
	{0x00, 0x00, 0x00, 0xff},
	{0x22, 0x22, 0x22, 0xff},
	{0x33, 0x33, 0x33, 0xff},
	{0x44, 0x44, 0x44, 0xff},
	{0x55, 0x55, 0x55, 0xff},
	{0x66, 0x66, 0x66, 0xff},
	{0x77, 0x77, 0x77, 0xff},
	{0x88, 0x88, 0x88, 0xff},
	{0x99, 0x99, 0x99, 0xff},
	{0xaa, 0xaa, 0xaa, 0xff},
	{0xbb, 0xbb, 0xbb, 0xff},
	{0xcc, 0xcc, 0xcc, 0xff},
	{0xdd, 0xdd, 0xdd, 0xff},
	{0xee, 0xee, 0xee, 0xff},
	{0xff, 0xff, 0xff, 0xff},
	{0x00, 0x00, 0x00, 0x00},
}}

func dumpGFX(cmd *cobra.Command, args []string) {
	board := boards.Get(targetBoard)

	dumpFolder := "dump/gfx/"
	if err := os.RemoveAll(dumpFolder); err != nil {
		cmd.PrintErr(err)
		os.Exit(1)
	}
	if err := os.MkdirAll(dumpFolder, 0777); err != nil {
		cmd.Printf("Unable to create GFX dump folder '%s' : '%s'\n", dumpFolder, err.Error())
		os.Exit(1)
	}

	// Desinterleave
	rom := make([]byte, board.GFX.Size)
	desinterleave(cmd, board.GFX.Roms, rom)

	// Dump ROM
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		//for i := 2; i < 3; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			area := board.GFXAreas[i]
			if verbose {
				cmd.Println("Dumping GFX type", area.Dim)
			}
			dumpSheets(cmd, i, dumpFolder, area.Dim, rom[area.Start:area.Start+area.Size])
		}(i)
	}
	wg.Wait()
}

func dumpSheets(cmd *cobra.Command, prefix int, toDir string, dim int, rom []byte) {
	bytesPerSheet := 256 * 256 / 2
	numSheets := len(rom) / bytesPerSheet

	var wg sync.WaitGroup
	for i := 0; i < numSheets; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			offset := i * bytesPerSheet
			path := fmt.Sprintf("%s%d-%d.svg", toDir, prefix, i)
			dumpsheet(cmd, path, dim, rom[offset:offset+bytesPerSheet])
		}(i)
	}
	wg.Wait()
}

func drawLine(cmd *cobra.Command, line []byte, x int, y int, img *image.RGBA, dim int) {
	// 4 bytes -> 8 pixels
	// 8 bytes -> 16 pixels
	// 16 bytes -> 32 pixels
	if len(line) != dim/2 {
		cmd.Println("Unexpected line length for dim ", dim, ". Expected", dim/8, "but got", len(line))
	}
	// 8 pixels -> read  4 bytes
	//16 pixels -> read  8 bytes
	//32 pixels -> read 16 bytes
	cursor := 0
	for i := 0; i < dim/8; i++ {
		// Read four bytes
		bytes := make([]byte, 4)
		for j := 0; j < 4; j++ {
			bytes[j] = line[cursor]
			cursor += 1
		}

		// Write eight indices
		var bits = []byte{128, 64, 32, 16, 8, 4, 2, 1}
		for j := 7; j >= 0; j-- {
			var b1, b2, b3, b4 byte
			if bytes[0]&bits[j] != 0 {
				b1 = 1
			} else {
				b1 = 0
			}
			if bytes[1]&bits[j] != 0 {
				b2 = 1
			} else {
				b2 = 0
			}
			if bytes[2]&bits[j] != 0 {
				b3 = 1
			} else {
				b3 = 0
			}
			if bytes[3]&bits[j] != 0 {
				b4 = 1
			} else {
				b4 = 0
			}
			var value = b4<<3 | b3<<2 | b2<<1 | b1
			// Write
			xCoord := x + j + i*8
			img.Set(xCoord, y, greyPalette.colors[value])
		}
	}
}

func drawTile(cmd *cobra.Command, tile []byte, imgX int, imgY int, img *image.RGBA, dim int) {
	bytesPerLine := dim / 2
	for i := 0; i < dim; i++ {
		offset := i * bytesPerLine
		drawLine(cmd, tile[offset:offset+bytesPerLine], imgX, imgY+i, img, dim)
	}
}

func dumpsheet(cmd *cobra.Command, path string, dim int, sheet []byte) {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{256, 256}})
	tilePerAxis := 256 / dim
	bytesPerTile := dim * dim / 2
	for y := 0; y < tilePerAxis; y++ {
		for x := 0; x < tilePerAxis; x++ {
			offset := (x + y*tilePerAxis) * bytesPerTile
			drawTile(cmd, sheet[offset:offset+bytesPerTile], x*dim, y*dim, img, dim)
		}
	}

	var pngPayload bytes.Buffer
	err := png.Encode(&pngPayload, img)
	if err != nil {
		cmd.Printf("Unable to dump GFX '%s'\n", err.Error())
		os.Exit(1)
	}
	png2svg(cmd, &pngPayload, path, 16)
}

// Template for 16x16 tile sheets
//
//go:embed svgParts/16_top.txt
var svgTop16 []byte

//go:embed svgParts/16_mid.txt
var svgMid16 []byte

//go:embed svgParts/16_bot.txt
var svgBot16 []byte

func png2svg(cmd *cobra.Command, payload *bytes.Buffer, out string, bank int) {
	f, err := os.Create(out)
	if err != nil {
		cmd.PrintErr(err)
		return
	}
	defer f.Close()

	f.WriteString(string(svgTop16))
	f.WriteString(base64.StdEncoding.EncodeToString(payload.Bytes()))
	f.WriteString(string(svgMid16))
	f.WriteString(fmt.Sprintf("%04x", bank<<8))
	f.WriteString(string(svgBot16))
}

func desinterleave(cmd *cobra.Command, srcs []boards.ROM, dst []byte) {
	for _, rom := range srcs {
		path := sites.OutDir + rom.Filename
		content, err := os.ReadFile(path)
		if err != nil {
			cmd.Printf("Unable to open '%s'\n", path)
			os.Exit(1)
		}

		for j := 0; j < rom.Size/rom.WordSize; j++ {
			srcOffset := rom.Offset + j*rom.WordSize
			src := content[srcOffset : srcOffset+rom.WordSize]

			dstOffset := rom.DstOffset + j*rom.Skip
			dst := dst[dstOffset : dstOffset+rom.WordSize]

			for w := 0; w < rom.WordSize; w++ {
				dst[w] = src[w]
			}
		}
	}
}
