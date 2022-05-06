package oki

import (
	"fmt"
)

type OkiRomEntry struct {
	start uint32
	end   uint32
}

func (e *OkiRomEntry) size() uint32 {
	return e.end - e.start
}

type OkiRom struct {
	Phrases [][]byte
}

func NewOkiROM() OkiRom {
	rom := OkiRom{}
	rom.Phrases = make([][]byte, 0)
	return rom
}

func (o *OkiRom) readEntry(src []byte) *OkiRomEntry {
	entry := OkiRomEntry{}
	entry.start = uint32(src[0]) << 16
	entry.start |= uint32(src[1]) << 8
	entry.start |= uint32(src[2]) << 0
	entry.end = uint32(src[3]) << 16
	entry.end |= uint32(src[4]) << 8
	entry.end |= uint32(src[5]) << 0
	return &entry
}

func createOKI() OkiRom {
	oki := OkiRom{}
	oki.Phrases = make([][]byte, 0)
	return oki
}

const maxPhrases int = 127
const headerSize = 0x3FF
const indexEntrySize int = 8

func OpenOKI(rom []byte) *OkiRom {
	o := createOKI()

	// Parse
	for i := 1; i <= maxPhrases; i++ {
		addrSlice := rom[i*indexEntrySize : i*indexEntrySize+indexEntrySize]
		entry := o.readEntry(addrSlice)
		if entry.start == entry.end {
			break
		}
		o.Phrases = append(o.Phrases, rom[entry.start:entry.end])
	}

	return &o
}

func (o *OkiRom) AddPhrase(phrase []byte) {
	o.Phrases = append(o.Phrases, phrase)
}

func (o *OkiRom) writeEntry(addr OkiRomEntry, dst []byte) {
	if len(dst) != indexEntrySize {
		panic(fmt.Sprintf("Entry slice is wrong size. Expected %d but got %d", indexEntrySize, len(dst)))
	}
	dst[0] = byte(addr.start & 0xFF0000 >> 16)
	dst[1] = byte(addr.start & 0x00FF00 >> 8)
	dst[2] = byte(addr.start & 0x0000FF)
	dst[3] = byte(addr.end & 0xFF0000 >> 16)
	dst[4] = byte(addr.end & 0x00FF00 >> 8)
	dst[5] = byte(addr.end & 0x0000FF)
	dst[6] = 0
	dst[7] = 0
}

func (o *OkiRom) genROM(size int64) []byte {
	if len(o.Phrases) == 0 {
		return nil
	}

	if len(o.Phrases) > maxPhrases {
		panic(fmt.Sprintf("Too many phrases %d, max=%d", len(o.Phrases), maxPhrases))
	}

	var totalSize uint32 = 0
	for _, phase := range o.Phrases {
		totalSize += uint32(len(phase))
	}

	rom := make([]byte, size)

	header := rom[:headerSize]
	// The first entry must be left empty
	o.writeHeader(header[indexEntrySize:])
	o.writePhrases(rom[headerSize:])
	return rom
}

func (o *OkiRom) writeHeader(header []byte) {
	if len(header) != headerSize-indexEntrySize {
		panic(fmt.Sprintf("Unexpected oki header size. Got %d but expected %d", len(header), headerSize))
	}

	var cursor uint32 = 0
	for i, phrase := range o.Phrases {
		entry := OkiRomEntry{}
		entry.start = cursor
		entry.end = cursor + uint32(len(phrase)) - 1
		headerOffset := i * indexEntrySize
		o.writeEntry(entry, header[headerOffset:headerOffset+indexEntrySize])
	}
}

func (o *OkiRom) writePhrases(phrases []byte) {
	for _, p := range o.Phrases {
		copy(phrases, p)
		phrases = phrases[len(p):]
	}
}
