package gfx

import (
	"ccps/boards"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var verbose bool
var board boards.Board

const outDir = ".tmp/gfx/"

type gfxRegionType int

const (
	OBJ  gfxRegionType = 0
	SCR1 gfxRegionType = 1
	SCR2 gfxRegionType = 2
	SCR3 gfxRegionType = 3
)

type tileDim int

const (
	dimOBJ  tileDim = 16
	dimSRC1 tileDim = 8
	dimSRC2 tileDim = 16
	dimSRC3 tileDim = 32
)

type gfxRegion struct {
	start int
	end   int
	sort  gfxRegionType
}

const gfxSrcPath = "gfx/"

var sortPath = [4]string{
	gfxSrcPath + "obj",
	gfxSrcPath + "scr1",
	gfxSrcPath + "scr2",
	gfxSrcPath + "scr3",
}

func Build(v bool, dryRun bool, b *boards.Board) string {
	verbose = v
	board = *b

	// TOOD Figure out Mame region size (e.g: STF29)
	// See https://github.com/mamedev/mame/blob/master/src/mame/video/cps1.cpp#L1679
	// https://github.com/mamedev/mame/blob/master/src/mame/video/cps1.cpp#L1748
	// sf2 = mapper_STF29 (https://github.com/mamedev/mame/blob/master/src/mame/video/cps1.cpp#L1085)

	// Hardcoding it for now
	var regions = []gfxRegion{
		{
			start: 0x00000,
			end:   144*0x8000 - 1,
			sort:  OBJ,
		}, {
			start: 0x00000,
			end:   0x00000,
			sort:  SCR1,
		}, {
			start: 0x00000,
			end:   0x00000,
			sort:  SCR2,
		}, {
			start: 0x00000,
			end:   0x00000,
			sort:  SCR3,
		},
	}

	// TODO: Check if there is nothing to do

	var sizes [4]int
	for _, region := range regions {
		size := region.end - region.start
		sizes[region.sort] += size
	}

	gfxRom := make([]byte, board.GFX.Size)
	cursor := 0
	for i, path := range sortPath {
		// For every type of GFX assets (OBJ, SCR1, SCR2, SCR3)
		// create a "sort rom".
		rom := createGFX(path, sizes[i], gfxRegionType(i))
		// Add "sort rom" to "everything" GFX ROM
		copy(gfxRom[cursor:], rom)
		cursor += len(rom)
	}

	// Write gfxrom to storage
	// write the whole body at once
	err := os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		fmt.Println("Unable to create dir", outDir)
		os.Exit(1)
	}
	romPath := outDir + "gfx.rom"
	err = ioutil.WriteFile(romPath, gfxRom, 0644)
	if err != nil {
		fmt.Println("Unable to write GFX rom to", romPath)
		os.Exit(1)
	}

	// Return everything rom path
	return romPath
}

func getTileDim(sort gfxRegionType) int {
	switch sort {
	case OBJ:
		return int(dimOBJ)
	case SCR1:
		return int(dimSRC1)
	case SCR2:
		return int(dimSRC2)
	case SCR3:
		return int(dimSRC3)
	}

	println("Requested tile dimension for unknown sort:", sort)
	os.Exit(1)
	return 0
}

// Visit all PNG in folder, find a free location and write them in rom
func createGFX(srcsPath string, size int, sort gfxRegionType) []byte {
	var rom = make([]byte, size)

	if verbose {
		println("Created ROM size", len(rom), " for region ", sort)
	}

	tileDim := getTileDim(sort)
	numTiles := len(rom) / int(tileDim)
	allocator := makeAllocator(numTiles)

	files, err := ioutil.ReadDir(srcsPath)
	if err != nil {
		println("Unable to open gfx dir", srcsPath)
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".png") {
			if verbose {
				println("Skipping non-png '", file.Name(), "'")
			}
			continue
		}

		if verbose {
			println("Processing image '", file.Name(), "'")
		}

		addGFX(srcsPath+"/"+file.Name(), rom, tileDim, allocator)
		// TODO write palette to .h so 68000 can use it.
		// TODO write either a sprite or a shape
	}
	return rom
}

func addGFX(src string, rom []byte, tileDim int, allocator *allocator) {

	file, err := os.Open(src)
	if err != nil {
		println("Unable to open file '", src, "'")
		os.Exit(1)
	}
	defer file.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	img, _, err := image.Decode(file)
	if err != nil {
		println("Unable to decode image'", src, "'")
		os.Exit(1)
	}

	_, ok := img.(image.PalettedImage)
	if !ok {
		if verbose {
			println("Image '", src, ", is not a paletted PNG")
			os.Exit(0)
		}
	}

	i, _ := img.(*image.Paletted)
	if len(i.Palette) > 16 {
		println("Image '", src, "' has more than 16 colors (", len(i.Palette), "'")
		os.Exit(1)
	}

	// Make sure transparency if properly set (index is 15).
	transparentIndex := -1
	for i, c := range i.Palette {
		_, _, _, a := c.RGBA()
		if a == 0 {
			if transparentIndex == -1 {
				transparentIndex = i
			} else {
				println("Image '", src, "' must have exacltly one transparent colors (found ", transparentIndex, ")")
				os.Exit(1)
			}
		}
	}

	if transparentIndex != 15 {
		makeTransparent15(i, uint8(transparentIndex))
	}

	// Round up dimension so it perfectly matches tiles layout
	adjustRectToTile(i, int(tileDim))

	// Image is ready. Write it to ROM
	filename := filepath.Base(src)
	var tileDsts []int
	if unicode.IsUpper(rune(filename[0])) {
		// This is a sprite (rectangular shape)
		tileDsts = allocateSprite(allocator, &i.Rect, tileDim)
	} else {
		// This is a shape (collection of tiles)
		tileDsts = allocateShape(allocator, &i.Rect, tileDim)
	}

	// Write tiles according to allocated tiles destinations
	writeTiles(i, tileDsts, rom, tileDim)
}

func adjustRectToTile(img *image.Paletted, tileDim int) {

	// Round up to next multiple of tileDim
	diffX := (tileDim - img.Rect.Max.X%tileDim) % tileDim
	diffY := (tileDim - img.Rect.Max.Y%tileDim) % tileDim

	// If the image is already in the proper dimensions, nothing to do
	if diffY == 0 && diffX == 0 {
		if verbose {
			println("No adjustment needed")
		}
		return
	}

	width := img.Rect.Max.X + diffX
	height := img.Rect.Max.Y + diffY

	if verbose {
		println("Adjustment needed")
	}

	numPixels := width * height

	newPx := make([]uint8, numPixels)
	for i := 0; i < numPixels; i++ {
		newPx[i] = 15
	}

	// Copy old pic into new pic, line by line
	for h := 0; h < img.Rect.Max.Y; h++ {
		lineSize := img.Rect.Max.X
		newPixLineDst := h * width
		PixLineSrc := h * img.Stride
		copy(newPx[newPixLineDst:newPixLineDst+lineSize], img.Pix[PixLineSrc:PixLineSrc+lineSize])
	}

	img.Rect = image.Rect(0, 0, width, height)
	img.Stride = width
	img.Pix = newPx
}

func writeTiles(i *image.Paletted, dsts []int, rom []byte, tileDim int) {
	width := int(math.Round(float64(i.Rect.Max.X) / float64(tileDim)))
	height := int(math.Round(float64(i.Rect.Max.Y) / float64(tileDim)))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			writeTile(x, y, i, rom, tileDim, dsts[y*width+x])
		}
	}
}

func writeTile(xTile int, yTile int, i *image.Paletted, rom []byte, tileDim int, dstTile int) {
	var width int
	if xTile <= (i.Rect.Max.X / tileDim) {
		width = tileDim
	} else {
		width = i.Rect.Max.X % tileDim
	}

	var height int
	if xTile <= (i.Rect.Max.Y / tileDim) {
		height = tileDim
	} else {
		height = i.Rect.Max.Y % tileDim
	}

	for h := 0; h < height; h++ {
		lineDest := dstTile*tileDim*tileDim/2 + 128*h
		bytes := width / 2
		writeTileLine(i, width, xTile*tileDim, yTile*tileDim, rom[lineDest:lineDest+bytes])
	}
}

func writeTileLine(img *image.Paletted, width int, x int, y int, dst []byte) {
	var acc uint8 = 0
	var cursor = 0
	for i := 0; i < width; i++ {
		index := img.ColorIndexAt(x+i, y)
		acc <<= 4
		acc |= index
		if (i % 2) == 1 { // Write to rom every other pen values.
			dst[cursor] = acc
			cursor += 1
			acc = 0
		}
	}
}

func allocateShape(allocator *allocator, bounds *image.Rectangle, tileDim int) []int {
	var tiles []int

	// Dimension (in tiles) of this image
	width := math.Round(float64(bounds.Max.X) / float64(tileDim))
	height := math.Round(float64(bounds.Max.Y) / float64(tileDim))
	numTiles := int(width * height)

	for !allocator.isEmpty() {
		tileId, err := allocator.any()
		if err != nil {

			os.Exit(1)
		}
		tiles = append(tiles, tileId)
		numTiles--
		if numTiles == 0 {
			break
		}
	}

	if numTiles > 0 {
		println("Out of GFX memory.")
		os.Exit(1)
	}

	return tiles
}

func allocateSprite(allocator *allocator, bounds *image.Rectangle, tileDim int) []int {
	// TODO
	println("Sprites not supported yet")
	os.Exit(1)
	return nil
}

// This function makes sure the image uses 15 as transparent color.
func makeTransparent15(i *image.Paletted, transpIndex uint8) {
	for x := 0; x < i.Bounds().Max.X; x++ {
		for y := 0; y < i.Bounds().Max.Y; y++ {
			if i.ColorIndexAt(x, y) == transpIndex {
				i.SetColorIndex(x, y, 15)
			}
			if i.ColorIndexAt(x, y) == 15 {
				i.SetColorIndex(x, y, transpIndex)
			}
		}
	}
}
