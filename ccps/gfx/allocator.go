package gfx

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
)

type allocator struct {
	rover  int
	tiles  []uint64 // 1 tile = 1 bit
	stride int
}

func makeAllocator(numTiles int, tileDim int) *allocator {
	var alloc = allocator{}

	numInts := numTiles / 64
	remainder := numTiles % 64

	if remainder != 0 {
		numInts += 1
	}

	alloc.tiles = make([]uint64, numTiles)
	for i := 0; i < numTiles; i++ {
		alloc.tiles[i] = math.MaxUint64
	}

	// First tile is not used
	//alloc.tiles[0] = math.MaxUint64
	//alloc.tiles[0] <<= 1

	if remainder != 0 {
		alloc.tiles[len(alloc.tiles)-1] = math.MaxUint64 >> (64 - remainder)
	}

	alloc.rover = 0
	//alloc.tileDim = tileDim
	alloc.stride = 256 / tileDim

	return &alloc
}

func (a *allocator) any() (int, error) {
	// Search for first available
	for a.tiles[a.rover] == 0 {
		a.rover++
		if a.rover == len(a.tiles) {
			return 0, errors.New("Unable to allocate: Out of ROM")
		}
	}
	bucket := a.tiles[a.rover]
	value := bits.TrailingZeros64(bucket)

	// Mark value as allocated
	mask := uint64(1) << value
	a.tiles[a.rover] = bucket & ^mask

	return a.rover*64 + value, nil
}

func (a *allocator) allocSprite(w int, h int) ([]int, error) {
	peeker := a.rover

	for peeker != len(a.tiles) {
		bucket := a.tiles[peeker]
		for bucket != 0 {
			tileIndex := peeker + bits.TrailingZeros64(bucket)
			if a.checkBlock(tileIndex, w, h) {
				// Mark block allocated
				return a.markBlock(tileIndex, w, h), nil
			}
		}
		peeker++
	}

	msg := fmt.Sprintf("Unable to allocate sprite (%d,%d)", w, h)
	return nil, errors.New(msg)
}

func (a *allocator) checkBlock(index int, width int, height int) bool {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if !a.has(index + y*a.stride + x) {
				return false
			}
		}
	}
	return true
}

func (a *allocator) has(v int) bool {
	bucket := a.tiles[v/64]
	mask := uint64(1) << v % 64
	return bucket&mask == mask
}

func (a *allocator) mark(v int) {
	bucket := a.tiles[v/64]
	mask := uint64(1) << (v % 64)
	a.tiles[v/64] = bucket & ^mask
}

func (a *allocator) markBlock(index int, w int, h int) []int {
	indexes := make([]int, 0)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			allocatedIndex := index + x + y*a.stride
			a.mark(allocatedIndex)
			indexes = append(indexes, allocatedIndex)
		}
	}
	return indexes
}
