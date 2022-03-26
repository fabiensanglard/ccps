package oki

import (
	"fmt"
	"io/ioutil"
	"os"
)

type OkiRomEntry struct {
	start uint32
	end   uint32
}

func (e *OkiRomEntry) size() uint32 {
	return e.end - e.start + 1
}

type OkiRom struct {
	phrases [][]byte
}

func NewOkiROM() OkiRom {
	rom := OkiRom{}
	rom.phrases = make([][]byte, 0)
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
	oki.phrases = make([][]byte, 0)
	return oki
}

const maxPhrases int = 127
const headerSize = 0x3FF
const indexEntrySize int = 8

func OpenOKI(path string) (*OkiRom, error) {
	rom, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	o := createOKI()

	// Parse
	for i := 1; i <= maxPhrases; i++ {
		addrSlice := rom[i*indexEntrySize : i*indexEntrySize+indexEntrySize]
		entry := o.readEntry(addrSlice)
		if entry.size() == 0 {
			continue
		}
		o.phrases = append(o.phrases, rom[entry.start:entry.end])
	}

	return &o, nil
}

func (o *OkiRom) AddPhrase(phrase []byte) {
	o.phrases = append(o.phrases, phrase)
}

func (o *OkiRom) writeEntry(addr OkiRomEntry, dst []byte) {
	if len(dst) != indexEntrySize {
		println("Entry slice is wrong size. Expected", indexEntrySize, "but got", len(dst))
		os.Exit(1)
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

func (o *OkiRom) genROM() []byte {
	if len(o.phrases) == 0 {
		return nil
	}

	if len(o.phrases) > maxPhrases {
		println(fmt.Sprintf("Too many phrases %d, max=127"))
		os.Exit(1)
	}

	var totalSize uint32 = 0
	for _, phase := range o.phrases {
		totalSize += uint32(len(phase))
	}

	rom := make([]byte, 65536)
	o.writeHeader(rom[0:headerSize])
	o.writePhrases(rom[headerSize:])
	return rom
}

func (o *OkiRom) writeHeader(header []byte) {
	if len(header) != headerSize {
		println("Unexpected oki header size. Got", len(header), "but expected", headerSize)
		os.Exit(1)
	}

	var cursor uint32 = 0
	for i, phrase := range o.phrases {
		entry := OkiRomEntry{}
		entry.start = cursor
		entry.end = cursor + uint32(len(phrase))
		headerOffset := i * indexEntrySize
		o.writeEntry(entry, header[headerOffset:headerOffset+indexEntrySize])
	}
}

func (o *OkiRom) writePhrases(phrases []byte) {

}
