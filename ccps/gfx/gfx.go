package gfx

import (
	"ccps/boards"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math/bits"
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

func Build(v bool, b *boards.Board) []byte {
	verbose = v
	board = *b

	// TOOD Figure out Mame region size (e.g: STF29)
	// See https://github.com/mamedev/mame/blob/master/src/mame/video/cps1.cpp#L1679
	// https://github.com/mamedev/mame/blob/master/src/mame/video/cps1.cpp#L1748
	// sf2 = mapper_STF29 (https://github.com/mamedev/mame/blob/master/src/mame/video/cps1.cpp#L1085)

	// TODO take this value from board
	// Hardcoding it for now
	var regions = []gfxRegion{
		{
			start: b.GFXAreas[0].Start,
			end:   b.GFXAreas[0].Start + b.GFXAreas[0].Size,
			sort:  OBJ,
		}, {
			start: b.GFXAreas[1].Start,
			end:   b.GFXAreas[1].Start + b.GFXAreas[1].Size,
			sort:  SCR1,
		}, {
			start: b.GFXAreas[2].Start,
			end:   b.GFXAreas[2].Start + b.GFXAreas[2].Size,
			sort:  SCR2,
		}, {
			start: b.GFXAreas[3].Start,
			end:   b.GFXAreas[3].Start + b.GFXAreas[3].Size,
			sort:  SCR3,
		},
	}

	// Test if there is a gfx src folder. If not, return null
	if _, err := os.Stat(gfxSrcPath); os.IsNotExist(err) {
		return nil
	}

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

	rom, err := os.ReadFile(romPath)
	if err != nil {
		println("Cannot read generated GFX ROM", err)
		os.Exit(1)
	}

	return rom
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
	for i := 0; i < len(rom); i++ {
		rom[i] = 0xFF
	}

	if verbose {
		println("Created ROM size", len(rom), " for region ", sort)
	}

	tileDim := getTileDim(sort)
	numTiles := len(rom) / tileDim
	allocator := makeAllocator(numTiles, tileDim)

	files, err := ioutil.ReadDir(srcsPath)
	if err != nil {
		if verbose {
			println("Unable to open gfx dir", srcsPath)
		}
		return rom
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
	transparentIndex := uint8(0)
	for i, c := range i.Palette {
		_, _, _, a := c.RGBA()
		if a == 0 {
			transparentIndex = uint8(i)
			break
		}
	}

	if transparentIndex != 15 {
		makeTransparent15(i, transparentIndex)
	}

	// Round up dimension so it perfectly matches tiles layout
	adjustRectToTile(i, tileDim)

	// Image is ready. Write it to ROM
	filename := filepath.Base(src)
	var allocations []Allocation
	if unicode.IsUpper(rune(filename[0])) {
		// This is a sprite (rectangular shape)
		allocations = allocateSprite(allocator, i, tileDim)
	} else {
		// This is a shape (collection of tiles)
		allocations = allocateShape(allocator, i, tileDim)
	}

	// Write tiles according to allocated tiles destinations
	writeTiles(i, allocations, rom, tileDim)
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

// Extract all ith bit of each byte in the array
func mask(mask byte, bytes []byte) byte {
	if len(bytes) != 8 {
		println("Requested masking of array len != 8")
		os.Exit(1)
	}

	r := uint8(0)
	for _, b := range bytes {
		r <<= 1
		r |= (b & mask) >> bits.TrailingZeros8(mask)
	}
	return r
}

func writeTileLine(img *image.Paletted, width int, x int, y int, dst []byte) {
	cursor := 0
	indexes := make([]byte, 8)
	for i := 0; i < width/8; i++ {
		indexes[0] = img.ColorIndexAt(x+i*8+0, y)
		indexes[1] = img.ColorIndexAt(x+i*8+1, y)
		indexes[2] = img.ColorIndexAt(x+i*8+2, y)
		indexes[3] = img.ColorIndexAt(x+i*8+3, y)
		indexes[4] = img.ColorIndexAt(x+i*8+4, y)
		indexes[5] = img.ColorIndexAt(x+i*8+5, y)
		indexes[6] = img.ColorIndexAt(x+i*8+6, y)
		indexes[7] = img.ColorIndexAt(x+i*8+7, y)

		dst[cursor+0] = mask(0x1, indexes)
		dst[cursor+1] = mask(0x2, indexes)
		dst[cursor+2] = mask(0x4, indexes)
		dst[cursor+3] = mask(0x8, indexes)
		cursor += 4
	}
}

// xTile = coordinate of src
func writeTile(x int, y int, i *image.Paletted, rom []byte, tileDim int, tileID int) {
	bytesPerTile := tileDim * tileDim / 2 // 4 bit per pixel, always
	bytesPerLine := tileDim / 2           // 4 bit per pixel, always
	romOffset := tileID * bytesPerTile
	tileDst := rom[romOffset : romOffset+bytesPerTile]
	for h := 0; h < tileDim; h++ {
		lineOffset := bytesPerLine * h
		writeTileLine(i, tileDim, x, y+h, tileDst[lineOffset:lineOffset+bytesPerLine])
	}
}

// Image i (src of colors)
// dsts Allocated tile IDs
// rom , the ROM
// tileID 8,16, or 32
func writeTiles(i *image.Paletted, dsts []Allocation, rom []byte, tileDim int) {
	for _, a := range dsts {
		writeTile(a.srcXTile*tileDim, a.srcYTile*tileDim, i, rom, tileDim, a.dst)
	}
}

type Allocation struct {
	srcXTile int // Img x tile coordinate
	srcYTile int // Img y tile coordinate
	dst      int // ROM tile destiantion
}

func allocateShape(allocator *allocator, img *image.Paletted, tileDim int) []Allocation {
	var tiles []Allocation

	// Dimension (in tiles) of this image
	width := img.Rect.Max.X / tileDim
	height := img.Rect.Max.Y / tileDim

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// If tile empty, skip it.
			if tileIsTransparent(img, x, y, tileDim) {
				continue
			}
			tileId, err := allocator.any()
			if err != nil {
				println("Out of GFX memory (dim=", tileDim, ")")
				os.Exit(1)
			}
			allocation := Allocation{x, y, tileId}
			tiles = append(tiles, allocation)
		}
	}

	return tiles
}

func allocateSprite(allocator *allocator, img *image.Paletted, tileDim int) []Allocation {

	// Dimension (in tiles) of this image
	width := img.Rect.Max.X / tileDim
	height := img.Rect.Max.Y / tileDim

	allocated, err := allocator.allocSprite(width, height)
	if err != nil {
		// TODO
		println("Unable to allocate block")
		os.Exit(1)
	}

	var allocations []Allocation
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			allocation := Allocation{x, y, allocated[x+y*width]}
			allocations = append(allocations, allocation)
		}
	}
	return allocations
}

func tileIsTransparent(img *image.Paletted, x int, y int, dim int) bool {
	for h := 0; h < dim; h++ {
		for w := 0; w < dim; w++ {
			if img.ColorIndexAt(x*dim+w, y*dim+h) != 15 {
				return false
			}
		}
	}
	return true
}

// This function makes sure the image uses 15 as transparent color.
func makeTransparent15(img *image.Paletted, transpIndex uint8) {
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			index := img.ColorIndexAt(x, y)
			if index == transpIndex {
				img.SetColorIndex(x, y, 15)
			} else if index == 15 {
				img.SetColorIndex(x, y, transpIndex)
			}
		}
	}

	palette := make([]color.Color, 16)
	for i, _ := range palette {
		palette[i] = color.RGBA{R: 255, G: 255, B: 255, A: 0}
	}
	for i, _ := range img.Palette {
		palette[i] = img.Palette[i]
	}

	palette[15] = color.RGBA{R: 255, G: 255, B: 255, A: 0}
	img.Palette = palette
}
