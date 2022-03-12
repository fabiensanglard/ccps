package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
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
	o := OkiRom{}
	o.phrases = make([][]byte, 0)
	return o
}

const MAX_PHRASES int = 127
const ADD_STORAGE_SIZE int = 8

func openOKI(path string) (*OkiRom, error) {
	rom, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	o := createOKI()

	// Parse
	for i := 1; i <= MAX_PHRASES; i++ {
		addrSlice := rom[i*ADD_STORAGE_SIZE : i*ADD_STORAGE_SIZE+ADD_STORAGE_SIZE]
		entry := o.readEntry(addrSlice)
		if entry.size() == 0 {
			continue
		}
		o.phrases = append(o.phrases, rom[entry.start:entry.end])
	}

	return &o, nil
}

func (o *OkiRom) addPhrase(phrase []byte) {
	o.phrases = append(o.phrases, phrase)
}

func (o *OkiRom) writeEntry(addr OkiRomEntry, dst []byte) {
	if len(dst) != ADD_STORAGE_SIZE {
		//var panik = fmt.Sprintf("Entry slice is wrong size %d", len(dst))
		panic("Entry slice is wrong size")
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

func (o *OkiRom) write(targetDir string, maxSize uint32) error {
	if len(o.phrases) == 0 {
		return nil
	}

	if len(o.phrases) > MAX_PHRASES {
		return errors.New("Too many phrases %d, max=127")
	}

	var totalSize uint32 = 0
	for _, phase := range o.phrases {
		totalSize += uint32(len(phase))
	}

	if totalSize > maxSize {
		var msg = fmt.Sprintf("Not enough SFX ROM space available (%d)", totalSize)
		return errors.New(msg)
	}

	rom := make([]byte, 65536)

	const HEADER_SIZE = 0x3FF
	// Write header first
	o.writeHeader(rom[0:HEADER_SIZE])
	o.writePhrases(rom[HEADER_SIZE:])

	//TODO Actually write files to disk now

	return nil
}

func (o *OkiRom) writeHeader(header []byte) {
	var cursor uint32 = 0
	for i, phrase := range o.phrases {
		entry := OkiRomEntry{}
		entry.start = cursor
		entry.end = cursor + uint32(len(phrase))
		headerOffset := i * ADD_STORAGE_SIZE
		o.writeEntry(entry, header[headerOffset:headerOffset+ADD_STORAGE_SIZE])
	}
}

func (o *OkiRom) writePhrases(phrases []byte) {

}

func main() {
	path, err2 := exec.LookPath("sdasz80")

	if err2 != nil {
		panic(err2)
	}
	println(path)

	fmt.Println("Starting...")
	cmd := exec.Command("sdasz80")
	//cmd := exec.Command("ls")

	var out bytes.Buffer
	cmd.Stdout = &out

	var err bytes.Buffer
	cmd.Stderr = &err
	cmd.Stdin = strings.NewReader("and old falcon")

	cmd.Run()
	fmt.Printf("out=%s\n", out.String())
	fmt.Printf("err=%s\n", err.String())
}
